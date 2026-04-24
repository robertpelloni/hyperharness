// Package mcp provides MCP (Model Context Protocol) client and aggregation
// capabilities. Ported from hypercode/go/internal/mcp/client.go and
// hypercode/go/internal/mcp/aggregator.go
//
// This implements the MCP client role: connecting to external MCP servers
// via stdio transport and aggregating their tools for unified access.
//
// WHAT: Stdio MCP client that communicates with MCP servers via JSON-RPC 2.0
// WHY: Enables hyperharness to connect to any MCP server and aggregate tools
// HOW: Spawns server processes, communicates over stdin/stdout pipes
package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"time"
)

// JSON-RPC 2.0 types for MCP communication.

// JSONRPCRequest is a JSON-RPC 2.0 request.
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

// JSONRPCResponse is a JSON-RPC 2.0 response.
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}

// RPCError is a JSON-RPC 2.0 error.
type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (e *RPCError) Error() string {
	return fmt.Sprintf("RPC error %d: %s", e.Code, e.Message)
}

// MCPTool represents a tool from an MCP server.
type MCPTool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

// MCPResource represents a resource from an MCP server.
type MCPResource struct {
	URI         string      `json:"uri"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	MimeType    string      `json:"mimeType"`
}

// MCPPrompt represents a prompt template from an MCP server.
type MCPPrompt struct {
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	Arguments   []MCPPromptArgument       `json:"arguments"`
}

// MCPPromptArgument represents a prompt template argument.
type MCPPromptArgument struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

// ServerInfo holds server metadata from the initialize response.
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// InitializeResult is the response from MCP initialize.
type InitializeResult struct {
	ProtocolVersion string     `json:"protocolVersion"`
	Capabilities    Capabilities `json:"capabilities"`
	ServerInfo      ServerInfo `json:"serverInfo"`
}

// Capabilities describes what an MCP server supports.
type Capabilities struct {
	Tools     *ToolCapabilities     `json:"tools,omitempty"`
	Resources *ResourceCapabilities `json:"resources,omitempty"`
	Prompts   *PromptCapabilities   `json:"prompts,omitempty"`
}

// ToolCapabilities describes tool capabilities.
type ToolCapabilities struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// ResourceCapabilities describes resource capabilities.
type ResourceCapabilities struct {
	Subscribe   bool `json:"subscribe,omitempty"`
	ListChanged bool `json:"listChanged,omitempty"`
}

// PromptCapabilities describes prompt capabilities.
type PromptCapabilities struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// StdioClient connects to an MCP server via stdin/stdout pipes.
// It handles the full MCP lifecycle: initialize → tool discovery → tool calls.
type StdioClient struct {
	Name    string
	Command string
	Args    []string
	Env     map[string]string

	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser

	pendingMu sync.Mutex
	pending   map[interface{}]chan *JSONRPCResponse
	nextID    int64

	// Server state after initialization
	Info         ServerInfo
	Capabilities Capabilities
	Tools        []MCPTool
	Resources    []MCPResource
	Prompts      []MCPPrompt
	initialized  bool
}

// NewStdioClient creates a new MCP stdio client.
func NewStdioClient(name, command string, args []string, env map[string]string) *StdioClient {
	return &StdioClient{
		Name:    name,
		Command: command,
		Args:    args,
		Env:     env,
		pending: make(map[interface{}]chan *JSONRPCResponse),
	}
}

// Start spawns the server process and begins reading responses.
func (c *StdioClient) Start() error {
	c.cmd = exec.Command(c.Command, c.Args...)

	// Extend current process environment
	c.cmd.Env = append([]string{}, os.Environ()...)
	for k, v := range c.Env {
		c.cmd.Env = append(c.cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	stdin, err := c.cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}
	c.stdin = stdin

	stdout, err := c.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	c.stdout = stdout

	// Capture stderr for debugging
	c.cmd.Stderr = os.Stderr

	if err := c.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	// Start reading responses in background
	go c.readLoop()

	return nil
}

// Initialize sends the MCP initialize handshake.
func (c *StdioClient) Initialize(ctx context.Context) error {
	params := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"roots": map[string]interface{}{
				"listChanged": true,
			},
		},
		"clientInfo": map[string]interface{}{
			"name":    "HyperHarness",
			"version": "0.2.0",
		},
	}

	resp, err := c.Call(ctx, "initialize", params)
	if err != nil {
		return fmt.Errorf("initialize failed: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("initialize error: %v", resp.Error)
	}

	// Parse initialize result
	resultBytes, err := json.Marshal(resp.Result)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	var initResult InitializeResult
	if err := json.Unmarshal(resultBytes, &initResult); err != nil {
		return fmt.Errorf("failed to parse initialize result: %w", err)
	}

	c.Info = initResult.ServerInfo
	c.Capabilities = initResult.Capabilities

	// Send initialized notification
	_, _ = c.Notify(ctx, "notifications/initialized", nil)

	c.initialized = true

	// Discover tools if supported
	if c.Capabilities.Tools != nil {
		if err := c.discoverTools(ctx); err != nil {
			// Non-fatal: tools discovery failure shouldn't prevent usage
			fmt.Fprintf(os.Stderr, "[MCP %s] Tool discovery warning: %v\n", c.Name, err)
		}
	}

	// Discover resources if supported
	if c.Capabilities.Resources != nil {
		_ = c.discoverResources(ctx)
	}

	// Discover prompts if supported
	if c.Capabilities.Prompts != nil {
		_ = c.discoverPrompts(ctx)
	}

	return nil
}

// discoverTools fetches the tool list from the server.
func (c *StdioClient) discoverTools(ctx context.Context) error {
	resp, err := c.Call(ctx, "tools/list", nil)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return resp.Error
	}

	resultBytes, _ := json.Marshal(resp.Result)
	var listResult struct {
		Tools []MCPTool `json:"tools"`
	}
	if err := json.Unmarshal(resultBytes, &listResult); err != nil {
		return err
	}

	c.Tools = listResult.Tools
	return nil
}

// discoverResources fetches the resource list from the server.
func (c *StdioClient) discoverResources(ctx context.Context) error {
	resp, err := c.Call(ctx, "resources/list", nil)
	if err != nil {
		return err
	}
	if resp.Result != nil {
		resultBytes, _ := json.Marshal(resp.Result)
		var listResult struct {
			Resources []MCPResource `json:"resources"`
		}
		if err := json.Unmarshal(resultBytes, &listResult); err == nil {
			c.Resources = listResult.Resources
		}
	}
	return nil
}

// discoverPrompts fetches the prompt template list from the server.
func (c *StdioClient) discoverPrompts(ctx context.Context) error {
	resp, err := c.Call(ctx, "prompts/list", nil)
	if err != nil {
		return err
	}
	if resp.Result != nil {
		resultBytes, _ := json.Marshal(resp.Result)
		var listResult struct {
			Prompts []MCPPrompt `json:"prompts"`
		}
		if err := json.Unmarshal(resultBytes, &listResult); err == nil {
			c.Prompts = listResult.Prompts
		}
	}
	return nil
}

// CallTool invokes a tool on the MCP server.
func (c *StdioClient) CallTool(ctx context.Context, toolName string, arguments map[string]interface{}) (interface{}, error) {
	params := map[string]interface{}{
		"name":      toolName,
		"arguments": arguments,
	}

	resp, err := c.Call(ctx, "tools/call", params)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Result, nil
}

// GetPrompt retrieves a prompt template from the server.
func (c *StdioClient) GetPrompt(ctx context.Context, promptName string, args map[string]string) (interface{}, error) {
	params := map[string]interface{}{
		"name":      promptName,
		"arguments": args,
	}
	resp, err := c.Call(ctx, "prompts/get", params)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}

// ReadResource reads a resource from the server.
func (c *StdioClient) ReadResource(ctx context.Context, uri string) (interface{}, error) {
	params := map[string]interface{}{
		"uri": uri,
	}
	resp, err := c.Call(ctx, "resources/read", params)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}

// Call sends a JSON-RPC request and waits for a response.
func (c *StdioClient) Call(ctx context.Context, method string, params interface{}) (*JSONRPCResponse, error) {
	id := atomic.AddInt64(&c.nextID, 1)

	c.pendingMu.Lock()
	ch := make(chan *JSONRPCResponse, 1)
	c.pending[id] = ch
	c.pendingMu.Unlock()

	defer func() {
		c.pendingMu.Lock()
		delete(c.pending, id)
		c.pendingMu.Unlock()
	}()

	req := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  params,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	if _, err := fmt.Fprintln(c.stdin, string(data)); err != nil {
		return nil, fmt.Errorf("failed to write to server: %w", err)
	}

	select {
	case resp := <-ch:
		return resp, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("request timed out: %w", ctx.Err())
	}
}

// Notify sends a JSON-RPC notification (no ID, no response expected).
func (c *StdioClient) Notify(ctx context.Context, method string, params interface{}) (*JSONRPCResponse, error) {
	req := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	if _, err := fmt.Fprintln(c.stdin, string(data)); err != nil {
		return nil, err
	}
	return nil, nil
}

// readLoop reads JSON-RPC responses from the server's stdout.
func (c *StdioClient) readLoop() {
	scanner := bufio.NewScanner(c.stdout)
	// Increase scanner buffer for large responses
	scanner.Buffer(make([]byte, 0, 1024*1024), 10*1024*1024)

	for scanner.Scan() {
		var resp JSONRPCResponse
		if err := json.Unmarshal(scanner.Bytes(), &resp); err != nil {
			continue // Skip non-JSON lines
		}

		if resp.ID != nil {
			c.pendingMu.Lock()
			if ch, ok := c.pending[resp.ID]; ok {
				ch <- &resp
				delete(c.pending, resp.ID)
			}
			c.pendingMu.Unlock()
		}
		// Notifications (no ID) are silently consumed
	}
}

// Stop terminates the server process.
func (c *StdioClient) Stop() error {
	if c.cmd != nil && c.cmd.Process != nil {
		// Try graceful shutdown first
		c.stdin.Close()
		done := make(chan error, 1)
		go func() {
			done <- c.cmd.Wait()
		}()

		select {
		case <-done:
			return nil
		case <-time.After(5 * time.Second):
			return c.cmd.Process.Kill()
		}
	}
	return nil
}

// IsRunning returns true if the server process is alive.
func (c *StdioClient) IsRunning() bool {
	if c.cmd == nil || c.cmd.Process == nil {
		return false
	}
	return c.cmd.ProcessState == nil
}

// ToolCount returns the number of discovered tools.
func (c *StdioClient) ToolCount() int {
	return len(c.Tools)
}
