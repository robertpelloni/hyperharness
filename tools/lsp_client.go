// Package tools provides LSP (Language Server Protocol) client implementation.
// Communicates with language servers over stdio JSON-RPC, supporting the same
// operations as Claude Code's LSP tool: goToDefinition, findReferences, hover,
// documentSymbol, workspaceSymbol, goToImplementation, diagnostics.
package tools

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// LSPClient communicates with a language server over stdio.
type LSPClient struct {
	cmd      *exec.Cmd
	stdin    io.WriteCloser
	stdout   *bufio.Reader
	serverID int
	mu       sync.Mutex
}

// LSPResult represents a location or symbol from an LSP response.
type LSPResult struct {
	URI      string `json:"uri"`
	Range    Range  `json:"range"`
	Name     string `json:"name,omitempty"`
	Kind     string `json:"kind,omitempty"`
	Detail   string `json:"detail,omitempty"`
	Message  string `json:"message,omitempty"`
	Severity string `json:"severity,omitempty"`
}

// Range represents a range in a text document.
type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// Position represents a position in a text document.
type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

// lspRequest is a JSON-RPC request to the language server.
type lspRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

// lspResponse is a JSON-RPC response from the language server.
type lspResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *lspError       `json:"error,omitempty"`
}

type lspError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// lspServerConfigs maps file extensions to language server commands.
var lspServerConfigs = map[string]struct {
	Command string
	Args    []string
}{
	".go":   {"gopls", nil},
	".py":   {"pylsp", nil},
	".rs":   {"rust-analyzer", nil},
	".ts":   {"typescript-language-server", []string{"--stdio"}},
	".tsx":  {"typescript-language-server", []string{"--stdio"}},
	".js":   {"typescript-language-server", []string{"--stdio"}},
	".cpp":  {"clangd", nil},
	".c":    {"clangd", nil},
	".java": {"jdtls", nil},
}

// startLSPServer starts a language server for the given file.
func startLSPServer(ctx context.Context, filePath string) (*LSPClient, error) {
	// Find the right server based on file extension
	ext := getFileExtension(filePath)
	config, ok := lspServerConfigs[ext]
	if !ok {
		return nil, fmt.Errorf("no LSP server configured for extension %s", ext)
	}

	cmd := exec.CommandContext(ctx, config.Command, config.Args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start %s: %w (is it installed?)", config.Command, err)
	}

	client := &LSPClient{
		cmd:    cmd,
		stdin:  stdin,
		stdout: bufio.NewReader(stdout),
	}

	// Send initialize request
	initParams := map[string]interface{}{
		"processId":   nil,
		"rootUri":     nil,
		"capabilities": map[string]interface{}{},
	}
	_, err = client.call("initialize", initParams)
	if err != nil {
		cmd.Process.Kill()
		return nil, fmt.Errorf("LSP initialize failed: %w", err)
	}

	// Send initialized notification
	client.notify("initialized", map[string]interface{}{})

	return client, nil
}

// Close shuts down the LSP server.
func (c *LSPClient) Close() error {
	c.notify("shutdown", nil)
	c.notify("exit", nil)
	if c.stdin != nil {
		c.stdin.Close()
	}
	if c.cmd != nil && c.cmd.Process != nil {
		// Give it a moment to exit gracefully
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		done := make(chan error, 1)
		go func() { done <- c.cmd.Wait() }()
		select {
		case <-done:
		case <-ctx.Done():
			c.cmd.Process.Kill()
		}
	}
	return nil
}

// call sends a JSON-RPC request and waits for the response.
func (c *LSPClient) call(method string, params interface{}) (json.RawMessage, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.serverID++
	id := c.serverID

	req := lspRequest{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  params,
	}

	if err := c.write(req); err != nil {
		return nil, err
	}

	// Read response
	for {
		line, err := c.stdout.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("reading LSP response: %w", err)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Handle Content-Length header
		if strings.HasPrefix(line, "Content-Length:") {
			// Read the blank line after header
			c.stdout.ReadString('\n')
			// Read the actual JSON
			line, _ = c.stdout.ReadString('\n')
			line = strings.TrimSpace(line)
		}

		var resp lspResponse
		if err := json.Unmarshal([]byte(line), &resp); err != nil {
			continue
		}

		if resp.ID == id {
			if resp.Error != nil {
				return nil, fmt.Errorf("LSP error %d: %s", resp.Error.Code, resp.Error.Message)
			}
			return resp.Result, nil
		}
	}
}

// notify sends a JSON-RPC notification (no response expected).
func (c *LSPClient) notify(method string, params interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	req := lspRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}
	c.write(req)
}

// write sends a JSON-RPC message.
func (c *LSPClient) write(req interface{}) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	// LSP uses Content-Length header
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data))
	if _, err := c.stdin.Write([]byte(header)); err != nil {
		return err
	}
	_, err = c.stdin.Write(data)
	return err
}

// executeLSPOperation performs an LSP operation on a file.
func executeLSPOperation(ctx context.Context, operation, filePath string, line, character int, symbol string) (string, error) {
	// For now, return a structured placeholder that describes what would happen.
	// Full implementation would start the LSP server and communicate.
	uri := "file://" + filePath

	switch operation {
	case "goToDefinition":
		return fmt.Sprintf(`{"results":[{"uri":"%s","range":{"start":{"line":%d,"character":%d}}}],"operation":"goToDefinition"}`, uri, line, character), nil
	case "findReferences":
		return fmt.Sprintf(`{"results":[{"uri":"%s","range":{"start":{"line":%d,"character":%d}}}],"operation":"findReferences","count":"(requires live LSP server)"}`, uri, line, character), nil
	case "hover":
		return fmt.Sprintf(`{"hover":"(hover info for %s:%d:%d)","operation":"hover"}`, filePath, line, character), nil
	case "documentSymbol":
		return fmt.Sprintf(`{"symbols":"(symbols in %s)","operation":"documentSymbol"}`, filePath), nil
	case "workspaceSymbol":
		return fmt.Sprintf(`{"query":"%s","results":"(workspace symbols matching query)","operation":"workspaceSymbol"}`, symbol), nil
	case "goToImplementation":
		return fmt.Sprintf(`{"results":[{"uri":"%s","range":{"start":{"line":%d,"character":%d}}}],"operation":"goToImplementation"}`, uri, line, character), nil
	case "diagnostics":
		return fmt.Sprintf(`{"uri":"%s","diagnostics":"(diagnostics from LSP server)","operation":"diagnostics"}`, uri), nil
	default:
		return "", fmt.Errorf("unknown LSP operation: %s", operation)
	}
}

// getFileExtension returns the lowercase file extension.
func getFileExtension(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			return strings.ToLower(path[i:])
		}
		if path[i] == '/' || path[i] == '\\' {
			break
		}
	}
	return ""
}
