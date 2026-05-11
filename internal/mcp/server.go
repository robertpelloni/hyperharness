package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/robertpelloni/hyperharness/tools"
)

// Serve exposes the given tool registry as an MCP server over stdio.
// It blocks reading from stdin and writing to stdout.
func Serve(registry *tools.Registry) error {
	reader := bufio.NewReader(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		var req JSONRPCMessage
		if err := json.Unmarshal(line, &req); err != nil {
			// Invalid JSON, ignore
			continue
		}

		if req.JSONRPC != "2.0" {
			continue
		}

		var res JSONRPCMessage
		res.JSONRPC = "2.0"
		res.ID = req.ID

		switch req.Method {
		case "initialize":
			res.Result = json.RawMessage(`{
				"protocolVersion": "2024-11-05",
				"capabilities": {
					"tools": {}
				},
				"serverInfo": {
					"name": "hyperharness-mcp",
					"version": "1.0.0"
				}
			}`)
		case "tools/list":
			// Convert registry tools to MCP tool format
			type MCPTool struct {
				Name        string          `json:"name"`
				Description string          `json:"description"`
				InputSchema json.RawMessage `json:"inputSchema"`
			}
			type ListResponse struct {
				Tools []MCPTool `json:"tools"`
			}
			
			var mcpTools []MCPTool
			for _, t := range registry.Tools {
				mcpTools = append(mcpTools, MCPTool{
					Name:        t.Name,
					Description: t.Description,
					InputSchema: t.Parameters,
				})
			}
			
			lr := ListResponse{Tools: mcpTools}
			b, _ := json.Marshal(lr)
			res.Result = b
			
		case "tools/call":
			var callReq struct {
				Name      string                 `json:"name"`
				Arguments map[string]interface{} `json:"arguments"`
			}
			if err := json.Unmarshal(req.Params, &callReq); err != nil {
				res.Error = &JSONRPCError{Code: -32602, Message: "Invalid params"}
				break
			}

			tool, ok := registry.Find(callReq.Name)
			if !ok {
				res.Error = &JSONRPCError{Code: -32601, Message: fmt.Sprintf("Tool not found: %s", callReq.Name)}
				break
			}

			// Map arguments
			result, err := tool.Execute(callReq.Arguments)

			type ToolResult struct {
				Content []struct {
					Type string `json:"type"`
					Text string `json:"text"`
				} `json:"content"`
				IsError bool `json:"isError"`
			}

			tr := ToolResult{IsError: err != nil}
			text := result
			if err != nil {
				text = err.Error()
			}
			tr.Content = []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			}{{Type: "text", Text: text}}

			b, _ := json.Marshal(tr)
			res.Result = b

		default:
			res.Error = &JSONRPCError{Code: -32601, Message: "Method not found"}
		}

		if err := encoder.Encode(res); err != nil {
			return err
		}
	}
	return nil
}
