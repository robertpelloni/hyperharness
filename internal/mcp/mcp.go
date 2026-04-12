// Package mcp implements MCP (Model Context Protocol) client and server.
// This is the unified tool transport layer that:
//   - Acts as MCP client: connects to external MCP servers (stdio/SSE)
//     and exposes their tools to the hyperharness agent
//   - Acts as MCP server: exposes hyperharness's own tools as MCP endpoints
//   - Tool inventory: maintains a unified catalog of ALL tools across
//     built-in tools, MCP server tools, and extension tools
//
// Built from the MCP specification (mcp.dev) with compatibility for
// all major implementations.
package mcp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"
)

// ProtocolVersion is the MCP protocol version.
const ProtocolVersion = "2024-11-05"

// JSONRPCMessage is the base JSON-RPC message.
type JSONRPCMessage struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      *json.Number    `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
}

// JSONRPCError is a JSON-RPC error.
type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// InitializeRequest is the MCP initialize request.
type InitializeRequest struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ClientInfo      map[string]string      `json:"clientInfo"`
}

// InitializeResult is the MCP initialize response.
type InitializeResult struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ServerInfo      map[string]string      `json:"serverInfo"`
}

// Tool represents an MCP tool.
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// ToolCall represents a request to call a tool.
type ToolCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// ToolResult represents a tool execution result.
type ToolResult struct {
	Content []ToolContent `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

// ToolContent is content within a tool result.
type ToolContent struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// MCPServer defines an external MCP server to connect to.
type MCPServer struct {
	Name      string            `json:"name"`
	Command   string            `json:"command,omitempty"`
	Args      []string          `json:"args,omitempty"`
	Env       map[string]string `json:"env,omitempty"`
	URL       string            `json:"url,omitempty"`
	Transport string            `json:"transport"` // "stdio" or "sse"
	Connected bool              `json:"connected"`
	Tools     []Tool            `json:"tools,omitempty"`

	// Internal runtime fields
	stdin  io.WriteCloser `json:"-"`
	stdout *bufio.Reader  `json:"-"`
	cmd    *exec.Cmd      `json:"-"`
	reqID  int            `json:"-"`
	mu     sync.Mutex     `json:"-"`

	// SSE specific
	endpoint string       `json:"-"`
	client   *http.Client `json:"-"`
}

// Registry manages all MCP servers and their tools.
type Registry struct {
	servers      map[string]*MCPServer
	allTools     map[string]*Tool   // tool name -> Tool (unique across all servers)
	serverTools  map[string][]*Tool // server name -> tools
	mu           sync.RWMutex
	builtInTools map[string]func(args map[string]interface{}) (ToolResult, error)
}

// NewRegistry creates a new MCP registry.
func NewRegistry() *Registry {
	return &Registry{
		servers:      make(map[string]*MCPServer),
		allTools:     make(map[string]*Tool),
		serverTools:  make(map[string][]*Tool),
		builtInTools: make(map[string]func(args map[string]interface{}) (ToolResult, error)),
	}
}

// RegisterServer registers an MCP server.
func (r *Registry) RegisterServer(server *MCPServer) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.servers[server.Name] = server
	return nil
}

// RegisterBuiltInTool registers a built-in tool with the MCP registry.
func (r *Registry) RegisterBuiltInTool(name string, fn func(args map[string]interface{}) (ToolResult, error)) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.builtInTools[name] = fn

	// Also register in allTools
	r.allTools[name] = &Tool{
		Name:        name,
		Description: "Built-in tool",
		InputSchema: map[string]interface{}{},
	}
}

// Connect connects to an MCP server and fetches its tools.
func (r *Registry) Connect(serverName string) error {
	r.mu.Lock()
	server, ok := r.servers[serverName]
	r.mu.Unlock()

	if !ok {
		return fmt.Errorf("server not found: %s", serverName)
	}

	if server.Transport == "stdio" {
		return r.connectStdio(server)
	} else if server.Transport == "sse" {
		return r.connectSSE(server)
	}

	return fmt.Errorf("unsupported transport: %s", server.Transport)
}

// connectStdio spawns an MCP server process and communicates via stdin/stdout.
func (r *Registry) connectStdio(server *MCPServer) error {
	cmd := exec.Command(server.Command, server.Args...)

	// Set environment
	env := os.Environ()
	for k, v := range server.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Env = env

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	// Send initialize
	initReq := InitializeRequest{
		ProtocolVersion: ProtocolVersion,
		ClientInfo:      map[string]string{"name": "hyperharness", "version": "0.1.0"},
		Capabilities:    map[string]interface{}{},
	}

	id := json.Number("1")
	initData, _ := json.Marshal(JSONRPCMessage{
		JSONRPC: "2.0",
		Method:  "initialize",
		ID:      &id,
		Params:  mustMarshal(initReq),
	})

	stdin.Write(append(initData, '\n'))

	// Read response
	decoder := json.NewDecoder(stdout)
	var result JSONRPCMessage
	if err := decoder.Decode(&result); err != nil {
		return err
	}

	// List tools
	listData := []byte(fmt.Sprintf(`{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`))
	stdin.Write(append(listData, '\n'))

	var toolsResult JSONRPCMessage
	if err := decoder.Decode(&toolsResult); err != nil {
		return err
	}

	// Parse tools
	var toolsList struct {
		Tools []Tool `json:"tools"`
	}
	if err := json.Unmarshal(toolsResult.Result, &toolsList); err != nil {
		return err
	}

	r.mu.Lock()
	server.Connected = true
	server.Tools = toolsList.Tools

	for i := range toolsList.Tools {
		tool := &toolsList.Tools[i]
		r.allTools[tool.Name] = tool
		r.serverTools[server.Name] = append(r.serverTools[server.Name], tool)
	}
	r.mu.Unlock()

	_ = stdout
	_ = cmd

	return nil
}

// connectSSE connects to an MCP server via SSE transport.
func (r *Registry) connectSSE(server *MCPServer) error {
	server.client = &http.Client{Timeout: 30 * time.Second}

	// 1. Send initialize
	server.reqID = 1
	initReq := InitializeRequest{
		ProtocolVersion: ProtocolVersion,
		ClientInfo:      map[string]string{"name": "hyperharness", "version": "0.1.0"},
		Capabilities:    map[string]interface{}{},
	}

	id := json.Number(fmt.Sprintf("%d", server.reqID))
	initData := JSONRPCMessage{
		JSONRPC: "2.0",
		Method:  "initialize",
		ID:      &id,
		Params:  mustMarshal(initReq),
	}

	reqBody, _ := json.Marshal(initData)
	resp, err := server.client.Post(server.URL, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to connect to SSE endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("SSE initialize returned status %d", resp.StatusCode)
	}

	var initResp JSONRPCMessage
	if err := json.NewDecoder(resp.Body).Decode(&initResp); err != nil {
		return fmt.Errorf("failed to decode SSE init response: %w", err)
	}

	// Server responds with a specific message endpoint URL we should use for tools/call
	// In a real MCP SSE connection, we get the endpoint header
	endpoint := resp.Header.Get("endpoint")
	if endpoint != "" {
		server.endpoint = endpoint
	} else {
		server.endpoint = server.URL // Fallback to same URL
	}

	// Send initialized notification
	initializedData := JSONRPCMessage{
		JSONRPC: "2.0",
		Method:  "notifications/initialized",
	}
	notifBody, _ := json.Marshal(initializedData)
	server.client.Post(server.endpoint, "application/json", bytes.NewReader(notifBody))

	// 2. Fetch tools
	server.reqID++
	listId := json.Number(fmt.Sprintf("%d", server.reqID))
	listData := JSONRPCMessage{
		JSONRPC: "2.0",
		Method:  "tools/list",
		ID:      &listId,
	}

	listReqBody, _ := json.Marshal(listData)
	listResp, err := server.client.Post(server.endpoint, "application/json", bytes.NewReader(listReqBody))
	if err != nil {
		return fmt.Errorf("failed to list tools via SSE: %w", err)
	}
	defer listResp.Body.Close()

	var toolsResult JSONRPCMessage
	if err := json.NewDecoder(listResp.Body).Decode(&toolsResult); err != nil {
		return fmt.Errorf("failed to decode tools/list response: %w", err)
	}

	var toolsList struct {
		Tools []Tool `json:"tools"`
	}
	if err := json.Unmarshal(toolsResult.Result, &toolsList); err != nil {
		return fmt.Errorf("failed to unmarshal tools list: %w", err)
	}

	r.mu.Lock()
	server.Connected = true
	server.Tools = toolsList.Tools

	for i := range toolsList.Tools {
		tool := &toolsList.Tools[i]
		r.allTools[tool.Name] = tool
		r.serverTools[server.Name] = append(r.serverTools[server.Name], tool)
	}
	r.mu.Unlock()

	return nil
}

// CallTool calls a tool by name across all providers.
func (r *Registry) CallTool(name string, args map[string]interface{}) (ToolResult, error) {
	r.mu.RLock()
	// Check built-in tools first
	if fn, ok := r.builtInTools[name]; ok {
		r.mu.RUnlock()
		return fn(args)
	}

	// Check MCP server tools
	_, ok := r.allTools[name]
	r.mu.RUnlock()

	if !ok {
		return ToolResult{
			Content: []ToolContent{{Type: "text", Text: fmt.Sprintf("Tool not found: %s", name)}},
			IsError: true,
		}, fmt.Errorf("tool not found: %s", name)
	}

	// Find which server owns this tool and call it
	for serverName, tools := range r.serverTools {
		for _, t := range tools {
			if t.Name == name {
				return r.callMCPTool(serverName, name, args)
			}
		}
	}

	return ToolResult{
		Content: []ToolContent{{Type: "text", Text: fmt.Sprintf("No server found for tool: %s", name)}},
		IsError: true,
	}, fmt.Errorf("no server found for tool: %s", name)
}

// callMCPTool calls a tool on a specific MCP server.
func (r *Registry) callMCPTool(serverName, toolName string, args map[string]interface{}) (ToolResult, error) {
	r.mu.RLock()
	server, ok := r.servers[serverName]
	r.mu.RUnlock()

	if !ok {
		return ToolResult{IsError: true}, fmt.Errorf("server %s not found", serverName)
	}

	if server.Transport == "stdio" {
		server.mu.Lock()
		defer server.mu.Unlock()

		server.reqID++
		idStr := fmt.Sprintf("%d", server.reqID)
		id := json.Number(idStr)

		params := map[string]interface{}{
			"name":      toolName,
			"arguments": args,
		}

		req := JSONRPCMessage{
			JSONRPC: "2.0",
			ID:      &id,
			Method:  "tools/call",
			Params:  mustMarshal(params),
		}

		reqData, err := json.Marshal(req)
		if err != nil {
			return ToolResult{IsError: true}, err
		}

		if _, err := server.stdin.Write(append(reqData, '\n')); err != nil {
			return ToolResult{IsError: true}, fmt.Errorf("failed to write tools/call request: %w", err)
		}

		line, err := server.stdout.ReadString('\n')
		if err != nil {
			return ToolResult{IsError: true}, fmt.Errorf("failed to read tools/call response: %w", err)
		}

		var resp JSONRPCMessage
		if err := json.Unmarshal([]byte(line), &resp); err != nil {
			return ToolResult{IsError: true}, fmt.Errorf("failed to decode tools/call response: %w", err)
		}

		if resp.Error != nil {
			return ToolResult{
				IsError: true,
				Content: []ToolContent{{Type: "text", Text: resp.Error.Message}},
			}, fmt.Errorf("MCP tool error %d: %s", resp.Error.Code, resp.Error.Message)
		}

		var toolResult ToolResult
		if err := json.Unmarshal(resp.Result, &toolResult); err != nil {
			return ToolResult{IsError: true}, fmt.Errorf("failed to parse tool result: %w", err)
		}

		return toolResult, nil
	} else if server.Transport == "sse" {
		server.mu.Lock()
		defer server.mu.Unlock()

		server.reqID++
		idStr := fmt.Sprintf("%d", server.reqID)
		id := json.Number(idStr)

		params := map[string]interface{}{
			"name":      toolName,
			"arguments": args,
		}

		req := JSONRPCMessage{
			JSONRPC: "2.0",
			ID:      &id,
			Method:  "tools/call",
			Params:  mustMarshal(params),
		}

		reqData, err := json.Marshal(req)
		if err != nil {
			return ToolResult{IsError: true}, err
		}

		resp, err := server.client.Post(server.endpoint, "application/json", bytes.NewReader(reqData))
		if err != nil {
			return ToolResult{IsError: true}, fmt.Errorf("failed to send tools/call via SSE: %w", err)
		}
		defer resp.Body.Close()

		var jsonResp JSONRPCMessage
		if err := json.NewDecoder(resp.Body).Decode(&jsonResp); err != nil {
			return ToolResult{IsError: true}, fmt.Errorf("failed to decode SSE tools/call response: %w", err)
		}

		if jsonResp.Error != nil {
			return ToolResult{
				IsError: true,
				Content: []ToolContent{{Type: "text", Text: jsonResp.Error.Message}},
			}, fmt.Errorf("MCP tool error %d: %s", jsonResp.Error.Code, jsonResp.Error.Message)
		}

		var toolResult ToolResult
		if err := json.Unmarshal(jsonResp.Result, &toolResult); err != nil {
			return ToolResult{IsError: true}, fmt.Errorf("failed to parse tool result: %w", err)
		}

		return toolResult, nil
	}

	return ToolResult{IsError: true}, fmt.Errorf("unsupported transport %s", server.Transport)
}

// ListTools returns all available tools (built-in + MCP).
func (r *Registry) ListTools() []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var tools []Tool
	for _, t := range r.allTools {
		tools = append(tools, *t)
	}
	return tools
}

// GetToolInventory returns a structured inventory.
func (r *Registry) GetToolInventory() map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	servers := make(map[string]interface{})
	for name, srv := range r.servers {
		servers[name] = map[string]interface{}{
			"connected": srv.Connected,
			"tools":     len(srv.Tools),
			"transport": srv.Transport,
		}
	}

	return map[string]interface{}{
		"totalBuiltIn": len(r.builtInTools),
		"totalMCP":     len(r.allTools) - len(r.builtInTools),
		"total":        len(r.allTools),
		"servers":      servers,
	}
}

// LoadFromFile loads MCP server definitions from a JSON file.
func (r *Registry) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var config struct {
		MCPServers map[string]struct {
			Command   string            `json:"command"`
			Args      []string          `json:"args"`
			Env       map[string]string `json:"env"`
			URL       string            `json:"url"`
			Transport string            `json:"transport"`
		} `json:"mcpServers"`
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	for name, srv := range config.MCPServers {
		transport := srv.Transport
		if transport == "" {
			if srv.URL != "" {
				transport = "sse"
			} else {
				transport = "stdio"
			}
		}

		server := &MCPServer{
			Name:      name,
			Command:   srv.Command,
			Args:      srv.Args,
			Env:       srv.Env,
			URL:       srv.URL,
			Transport: transport,
		}

		if err := r.RegisterServer(server); err != nil {
			return err
		}
	}

	return nil
}

func mustMarshal(v interface{}) json.RawMessage {
	data, _ := json.Marshal(v)
	return data
}
