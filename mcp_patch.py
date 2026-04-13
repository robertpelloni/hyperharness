import sys
import re

with open("internal/mcp/mcp.go", "r") as f:
    content = f.read()

# Add stdio and server fields
structs_pattern = """// MCPServer defines an external MCP server to connect to.
type MCPServer struct {
	Name      string            `json:"name"`
	Command   string            `json:"command,omitempty"`
	Args      []string          `json:"args,omitempty"`
	Env       map[string]string `json:"env,omitempty"`
	URL       string            `json:"url,omitempty"`
	Transport string            `json:"transport"` // "stdio" or "sse"
	Connected bool              `json:"connected"`
	Tools     []Tool            `json:"tools,omitempty"`
}"""

structs_replacement = """// MCPServer defines an external MCP server to connect to.
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
	stdin     io.WriteCloser    `json:"-"`
	stdout    *bufio.Reader     `json:"-"`
	cmd       *exec.Cmd         `json:"-"`
	reqID     int               `json:"-"`
	mu        sync.Mutex        `json:"-"`
}"""

content = content.replace(structs_pattern, structs_replacement)

# Fix connectStdio
connect_stdio_pattern = """// connectStdio spawns an MCP server process and communicates via stdin/stdout.
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

	stdin.Write(append(initData, '\\n'))

	// Read response
	decoder := json.NewDecoder(stdout)
	var result JSONRPCMessage
	if err := decoder.Decode(&result); err != nil {
		return err
	}

	// List tools
	listData := []byte(fmt.Sprintf(`{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`))
	stdin.Write(append(listData, '\\n'))

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
}"""

connect_stdio_replacement = """// connectStdio spawns an MCP server process and communicates via stdin/stdout.
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

	server.stdin = stdin
	server.stdout = bufio.NewReader(stdout)
	server.cmd = cmd
	server.reqID = 0

	// Send initialize
	initReq := InitializeRequest{
		ProtocolVersion: ProtocolVersion,
		ClientInfo:      map[string]string{"name": "hyperharness", "version": "0.1.0"},
		Capabilities:    map[string]interface{}{},
	}

	id := json.Number("1")
	server.reqID = 1
	initData, _ := json.Marshal(JSONRPCMessage{
		JSONRPC: "2.0",
		Method:  "initialize",
		ID:      &id,
		Params:  mustMarshal(initReq),
	})

	server.stdin.Write(append(initData, '\\n'))

	// Read response
	line, err := server.stdout.ReadString('\\n')
	if err != nil {
		return fmt.Errorf("failed to read init response: %w", err)
	}

	var result JSONRPCMessage
	if err := json.Unmarshal([]byte(line), &result); err != nil {
		return fmt.Errorf("failed to decode init response: %w", err)
	}

	// Send initialized notification
	initializedData := []byte(`{"jsonrpc":"2.0","method":"notifications/initialized","params":{}}\n`)
	server.stdin.Write(initializedData)

	// List tools
	server.reqID = 2
	listData := []byte(fmt.Sprintf(`{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}\n`))
	server.stdin.Write(listData)

	line, err = server.stdout.ReadString('\\n')
	if err != nil {
		return fmt.Errorf("failed to read tools/list response: %w", err)
	}

	var toolsResult JSONRPCMessage
	if err := json.Unmarshal([]byte(line), &toolsResult); err != nil {
		return fmt.Errorf("failed to decode tools/list response: %w", err)
	}

	if toolsResult.Error != nil {
		return fmt.Errorf("tools/list error %d: %s", toolsResult.Error.Code, toolsResult.Error.Message)
	}

	// Parse tools
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
}"""

content = content.replace(connect_stdio_pattern.replace("\\n", "\n"), connect_stdio_replacement.replace("\\n", "\n"))

# Fix callMCPTool
call_mcp_old = """// callMCPTool calls a tool on a specific MCP server.
func (r *Registry) callMCPTool(serverName, toolName string, args map[string]interface{}) (ToolResult, error) {
	// Implementation would send a tools/call request to the server
	return ToolResult{
		Content: []ToolContent{{Type: "text", Text: fmt.Sprintf("Called tool %s on server %s", toolName, serverName)}},
	}, nil
}"""

call_mcp_new = """// callMCPTool calls a tool on a specific MCP server.
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

		if _, err := server.stdin.Write(append(reqData, '\\n')); err != nil {
			return ToolResult{IsError: true}, fmt.Errorf("failed to write tools/call request: %w", err)
		}

		line, err := server.stdout.ReadString('\\n')
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
		return ToolResult{IsError: true}, fmt.Errorf("SSE transport not implemented for tools/call")
	}

	return ToolResult{IsError: true}, fmt.Errorf("unsupported transport %s", server.Transport)
}"""

content = content.replace(call_mcp_old.replace("\\n", "\n"), call_mcp_new.replace("\\n", "\n"))

# Add IO imports
if '"bufio"' not in content:
    content = content.replace(
        '"encoding/json"\n\t"fmt"\n\t\n\t"os"\n\t"os/exec"\n\t"sync"',
        '"bufio"\n\t"encoding/json"\n\t"fmt"\n\t"io"\n\t"os"\n\t"os/exec"\n\t"sync"'
    )

with open("internal/mcp/mcp.go", "w") as f:
    f.write(content)
print("Updated mcp.go")
