import sys
import re

with open("internal/mcp/mcp.go", "r") as f:
    content = f.read()

# Make sure imports are there
if '"net/http"' not in content:
    content = content.replace(
        '"os/exec"\n\t"sync"',
        '"os/exec"\n\t"sync"\n\t"net/http"\n\t"bytes"\n\t"context"'
    )

# Add SSE variables to MCPServer
server_pattern = """	// Internal runtime fields
	stdin     io.WriteCloser    `json:"-"`
	stdout    *bufio.Reader     `json:"-"`
	cmd       *exec.Cmd         `json:"-"`
	reqID     int               `json:"-"`
	mu        sync.Mutex        `json:"-"`"""

server_replacement = """	// Internal runtime fields
	stdin     io.WriteCloser    `json:"-"`
	stdout    *bufio.Reader     `json:"-"`
	cmd       *exec.Cmd         `json:"-"`
	reqID     int               `json:"-"`
	mu        sync.Mutex        `json:"-"`

	// SSE specific
	endpoint  string            `json:"-"`
	client    *http.Client      `json:"-"`"""

content = content.replace(server_pattern, server_replacement)

sse_pattern = """// connectSSE connects to an MCP server via SSE transport.
func (r *Registry) connectSSE(server *MCPServer) error {
	// SSE implementation would go here
	return fmt.Errorf("SSE transport not yet implemented")
}"""

sse_replacement = """// connectSSE connects to an MCP server via SSE transport.
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
}"""

content = content.replace(sse_pattern, sse_replacement)

call_mcp_old = """	} else if server.Transport == "sse" {
		return ToolResult{IsError: true}, fmt.Errorf("SSE transport not implemented for tools/call")
	}"""

call_mcp_new = """	} else if server.Transport == "sse" {
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
	}"""

content = content.replace(call_mcp_old, call_mcp_new)

with open("internal/mcp/mcp.go", "w") as f:
    f.write(content)
print("Updated mcp.go with SSE transport")
