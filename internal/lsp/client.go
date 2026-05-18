package lsp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"sync/atomic"
)

// ═══════════════════════════════════════════════════════════════════════
// client.go — LSP Client for HyperHarness
// Ported from opencode's LSP client with Go-native JSON-RPC transport.
//
// WHAT: Language Server Protocol client that connects to language servers
//       for diagnostics, completions, go-to-definition, and symbols.
// WHY: LSP integration provides real-time code intelligence without
//       requiring the LLM to parse and understand all code.
// HOW: JSON-RPC over stdin/stdout with LSP protocol messages.
// ═══════════════════════════════════════════════════════════════════════

// ── LSP Protocol Types ──

type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Location struct {
	URI   string `json:"uri"`
	Range Range  `json:"range"`
}

type Diagnostic struct {
	Range    Range  `json:"range"`
	Severity int    `json:"severity"`
	Message  string `json:"message"`
	Source   string `json:"source,omitempty"`
	Code     any    `json:"code,omitempty"`
}

type DocumentSymbol struct {
	Name           string           `json:"name"`
	Kind           int              `json:"kind"`
	Range          Range            `json:"range"`
	SelectionRange Range            `json:"selectionRange"`
	Children       []DocumentSymbol `json:"children,omitempty"`
}

type CompletionItem struct {
	Label         string `json:"label"`
	Kind          int    `json:"kind,omitempty"`
	Detail        string `json:"detail,omitempty"`
	Documentation string `json:"documentation,omitempty"`
	InsertText    string `json:"insertText,omitempty"`
}

type WorkspaceEdit struct {
	Changes map[string][]TextEdit `json:"changes,omitempty"`
}

type TextEdit struct {
	Range   Range  `json:"range"`
	NewText string `json:"newText"`
}

// ── JSON-RPC Types ──

type RPCRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id,omitempty"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
}

type RPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type RPCNotification struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
}

// ── LSP Client ──

// ClientState represents the client's connection state.
type ClientState int

const (
	ClientStopped ClientState = iota
	ClientStarting
	ClientRunning
	ClientError
)

// Client is an LSP client that communicates over stdin/stdout.
type Client struct {
	mu         sync.Mutex
	cmd        *exec.Cmd
	stdin      io.WriteCloser
	stdout     io.ReadCloser
	nextID     atomic.Int64
	pending    map[int]chan *RPCResponse
	state      ClientState
	handlers   map[string]func(any)
	diagnostics map[string][]Diagnostic
	serverInfo  map[string]any
}

// ClientConfig defines how to start a language server.
type ClientConfig struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env,omitempty"`
	RootURI string            `json:"rootURI"`
}

// NewClient creates a new LSP client.
func NewClient() *Client {
	return &Client{
		pending:     make(map[int]chan *RPCResponse),
		handlers:    make(map[string]func(any)),
		diagnostics: make(map[string][]Diagnostic),
	}
}

// Start launches the language server process.
func (c *Client) Start(ctx context.Context, config ClientConfig) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	args := append([]string{}, config.Args...)
	c.cmd = exec.CommandContext(ctx, config.Command, args...)

	// Set environment
	if len(config.Env) > 0 {
		env := c.cmd.Environ()
		for k, v := range config.Env {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
		c.cmd.Env = env
	}

	// Setup stdin/stdout pipes
	stdin, err := c.cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("stdin pipe: %w", err)
	}
	c.stdin = stdin

	stdout, err := c.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe: %w", err)
	}
	c.stdout = stdout

	if err := c.cmd.Start(); err != nil {
		return fmt.Errorf("start server: %w", err)
	}

	c.state = ClientStarting

	// Start reading responses
	go c.readLoop()

	// Initialize the LSP session
	initParams := map[string]any{
		"processId":    0,
		"rootUri":      config.RootURI,
		"capabilities": defaultCapabilities(),
	}

	_, err = c.Call("initialize", initParams)
	if err != nil {
		return fmt.Errorf("initialize: %w", err)
	}

	// Send initialized notification
	c.Notify("initialized", map[string]any{})

	c.state = ClientRunning
	return nil
}

// Stop shuts down the language server.
func (c *Client) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.state == ClientStopped {
		return nil
	}

	// Send shutdown request
	c.Call("shutdown", nil)
	c.Notify("exit", nil)

	c.stdin.Close()
	c.stdout.Close()

	if c.cmd.Process != nil {
		c.cmd.Process.Kill()
	}

	c.state = ClientStopped
	return nil
}

// State returns the current client state.
func (c *Client) State() ClientState {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state
}

// Call sends a JSON-RPC request and waits for the response.
func (c *Client) Call(method string, params any) (json.RawMessage, error) {
	id := int(c.nextID.Add(1))

	ch := make(chan *RPCResponse, 1)
	c.mu.Lock()
	c.pending[id] = ch
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
	}()

	req := RPCRequest{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  params,
	}

	if err := c.send(req); err != nil {
		return nil, fmt.Errorf("send: %w", err)
	}

	resp := <-ch
	if resp.Error != nil {
		return nil, fmt.Errorf("LSP error %d: %s", resp.Error.Code, resp.Error.Message)
	}

	return resp.Result, nil
}

// Notify sends a JSON-RPC notification (no response expected).
func (c *Client) Notify(method string, params any) error {
	notif := RPCNotification{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}
	return c.send(notif)
}

// ── LSP Operations ──

// DidOpen notifies the server that a file was opened.
func (c *Client) DidOpen(uri, languageID, text string) error {
	return c.Notify("textDocument/didOpen", map[string]any{
		"textDocument": map[string]any{
			"uri":        uri,
			"languageId": languageID,
			"version":    0,
			"text":       text,
		},
	})
}

// DidChange notifies the server that a file was modified.
func (c *Client) DidChange(uri string, version int, text string) error {
	return c.Notify("textDocument/didChange", map[string]any{
		"textDocument": map[string]any{
			"uri":     uri,
			"version": version,
		},
		"contentChanges": []map[string]any{
			{"text": text},
		},
	})
}

// DidClose notifies the server that a file was closed.
func (c *Client) DidClose(uri string) error {
	return c.Notify("textDocument/didClose", map[string]any{
		"textDocument": map[string]any{"uri": uri},
	})
}

// Diagnostics returns the latest diagnostics for a file.
func (c *Client) Diagnostics(uri string) []Diagnostic {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.diagnostics[uri]
}

// Completion requests code completions at a position.
func (c *Client) Completion(ctx context.Context, uri string, pos Position) ([]CompletionItem, error) {
	result, err := c.Call("textDocument/completion", map[string]any{
		"textDocument": map[string]any{"uri": uri},
		"position":     pos,
	})
	if err != nil {
		return nil, err
	}

	var items []CompletionItem
	if err := json.Unmarshal(result, &items); err != nil {
		// Try items field
		var wrapper struct {
			Items []CompletionItem `json:"items"`
		}
		if err2 := json.Unmarshal(result, &wrapper); err2 != nil {
			return nil, err
		}
		return wrapper.Items, nil
	}
	return items, nil
}

// Hover requests hover information at a position.
func (c *Client) Hover(ctx context.Context, uri string, pos Position) (string, error) {
	result, err := c.Call("textDocument/hover", map[string]any{
		"textDocument": map[string]any{"uri": uri},
		"position":     pos,
	})
	if err != nil {
		return "", err
	}

	var hover struct {
		Contents any `json:"contents"`
	}
	if err := json.Unmarshal(result, &hover); err != nil {
		return "", err
	}

	switch v := hover.Contents.(type) {
	case string:
		return v, nil
	case map[string]any:
		if value, ok := v["value"].(string); ok {
			return value, nil
		}
	}
	return fmt.Sprintf("%v", hover.Contents), nil
}

// GoToDefinition requests the definition location for a symbol.
func (c *Client) GoToDefinition(ctx context.Context, uri string, pos Position) ([]Location, error) {
	result, err := c.Call("textDocument/definition", map[string]any{
		"textDocument": map[string]any{"uri": uri},
		"position":     pos,
	})
	if err != nil {
		return nil, err
	}

	var locations []Location
	if err := json.Unmarshal(result, &locations); err != nil {
		// Might be a single location
		var loc Location
		if err2 := json.Unmarshal(result, &loc); err2 != nil {
			return nil, err
		}
		return []Location{loc}, nil
	}
	return locations, nil
}

// References requests all references to a symbol.
func (c *Client) References(ctx context.Context, uri string, pos Position) ([]Location, error) {
	result, err := c.Call("textDocument/references", map[string]any{
		"textDocument": map[string]any{"uri": uri},
		"position":     pos,
		"context":      map[string]any{"includeDeclaration": true},
	})
	if err != nil {
		return nil, err
	}

	var locations []Location
	if err := json.Unmarshal(result, &locations); err != nil {
		return nil, err
	}
	return locations, nil
}

// DocumentSymbols requests all symbols in a document.
func (c *Client) DocumentSymbols(ctx context.Context, uri string) ([]DocumentSymbol, error) {
	result, err := c.Call("textDocument/documentSymbol", map[string]any{
		"textDocument": map[string]any{"uri": uri},
	})
	if err != nil {
		return nil, err
	}

	var symbols []DocumentSymbol
	if err := json.Unmarshal(result, &symbols); err != nil {
		return nil, err
	}
	return symbols, nil
}

// ── Internal ──

func (c *Client) send(msg any) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// LSP uses Content-Length header
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data))

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, err := c.stdin.Write([]byte(header)); err != nil {
		return err
	}
	_, err = c.stdin.Write(data)
	return err
}

func (c *Client) readLoop() {
	buf := make([]byte, 4096)
	var content []byte
	var contentLength int

	for {
		// Read Content-Length header
		headerBuf := make([]byte, 0)
		for {
			n, err := c.stdout.Read(buf)
			if err != nil {
				return
			}
			headerBuf = append(headerBuf, buf[:n]...)

			// Check if we have the full header
			header := string(headerBuf)
			if len(header) > 0 {
				// Parse Content-Length
				fmt.Sscanf(header, "Content-Length: %d\r\n", &contentLength)
				// Find end of header
				idx := findHeaderEnd(headerBuf)
				if idx >= 0 {
					content = headerBuf[idx:]
					break
				}
			}
		}

		// Read content body
		for len(content) < contentLength {
			n, err := c.stdout.Read(buf)
			if err != nil {
				return
			}
			content = append(content, buf[:n]...)
		}

		if contentLength > 0 && len(content) >= contentLength {
			body := content[:contentLength]
			c.handleMessage(body)
			content = content[contentLength:]
		}
	}
}

func (c *Client) handleMessage(data []byte) {
	// Try as response first
	var resp RPCResponse
	if err := json.Unmarshal(data, &resp); err == nil && resp.ID > 0 {
		c.mu.Lock()
		ch, ok := c.pending[resp.ID]
		c.mu.Unlock()
		if ok {
			ch <- &resp
		}
		return
	}

	// Try as notification
	var notif RPCNotification
	if err := json.Unmarshal(data, &notif); err == nil && notif.Method != "" {
		c.mu.Lock()
		handler, ok := c.handlers[notif.Method]
		c.mu.Unlock()
		if ok && handler != nil {
			handler(notif.Params)
		}

		// Handle diagnostic notifications
		if notif.Method == "textDocument/publishDiagnostics" {
			c.handleDiagnostics(notif.Params)
		}
	}
}

func (c *Client) handleDiagnostics(params any) {
	data, err := json.Marshal(params)
	if err != nil {
		return
	}

	var notification struct {
		URI         string       `json:"uri"`
		Diagnostics []Diagnostic `json:"diagnostics"`
	}
	if err := json.Unmarshal(data, &notification); err != nil {
		return
	}

	c.mu.Lock()
	c.diagnostics[notification.URI] = notification.Diagnostics
	c.mu.Unlock()
}

// OnNotification registers a handler for a notification method.
func (c *Client) OnNotification(method string, handler func(any)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers[method] = handler
}

func findHeaderEnd(data []byte) int {
	for i := 0; i < len(data)-3; i++ {
		if data[i] == '\r' && data[i+1] == '\n' && data[i+2] == '\r' && data[i+3] == '\n' {
			return i + 4
		}
	}
	return -1
}

func defaultCapabilities() map[string]any {
	return map[string]any{
		"textDocument": map[string]any{
			"completion": map[string]any{
				"completionItem": map[string]any{
					"snippetSupport": false,
				},
			},
			"hover": map[string]any{
				"contentFormat": []string{"plaintext", "markdown"},
			},
			"publishDiagnostics": map[string]any{
				"relatedInformation": true,
			},
		},
	}
}
