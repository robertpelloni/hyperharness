// Package mcp provides the MCP server that exposes HyperHarness tools
// as MCP tool endpoints, allowing any MCP client to invoke our 109+ tools.
//
// This implements the "bidirectional routing" requirement: HyperHarness acts as
// both MCP client (connecting to external servers) and MCP server (exposing its
// own tools to external clients like Claude Code, Cursor, etc.).
package mcp

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

// MCPServer exposes HyperHarness tools via MCP protocol over stdio.
type MCPServer struct {
	tools      map[string]MCPToolDef
	handler    func(name string, args map[string]interface{}) (string, error)
	mu         sync.RWMutex
	clientConn net.Conn
}

// MCPToolDef defines a tool exposed via MCP.
type MCPToolDef struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Schema      json.RawMessage `json:"inputSchema"`
}

// NewMCPServer creates a new MCP server with the given tool handler.
func NewMCPServer(handler func(name string, args map[string]interface{}) (string, error)) *MCPServer {
	return &MCPServer{
		tools:   make(map[string]MCPToolDef),
		handler: handler,
	}
}

// RegisterTool adds a tool to the MCP server's catalog.
func (s *MCPServer) RegisterTool(def MCPToolDef) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tools[def.Name] = def
}

// RegisterTools bulk-registers tools.
func (s *MCPServer) RegisterTools(tools []MCPToolDef) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, t := range tools {
		s.tools[t.Name] = t
	}
}

// ToolCount returns the number of registered tools.
func (s *MCPServer) ToolCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.tools)
}

// ListTools returns all registered tool definitions.
func (s *MCPServer) ListTools() []MCPToolDef {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]MCPToolDef, 0, len(s.tools))
	for _, t := range s.tools {
		result = append(result, t)
	}
	return result
}

// ServeStdio starts the MCP server reading from stdin and writing to stdout.
// This is the standard MCP transport mode used by Claude Code, Cursor, etc.
func (s *MCPServer) ServeStdio(stdin io.Reader, stdout io.Writer) error {
	decoder := json.NewDecoder(stdin)
	encoder := json.NewEncoder(stdout)

	for {
		var msg map[string]interface{}
		if err := decoder.Decode(&msg); err != nil {
			if err == io.EOF {
				return nil
			}
			continue
		}

		response := s.handleMessage(msg)
		if response != nil {
			encoder.Encode(response)
		}
	}
}

// handleMessage processes a single JSON-RPC message.
func (s *MCPServer) handleMessage(msg map[string]interface{}) map[string]interface{} {
	method, _ := msg["method"].(string)
	id := msg["id"]

	// Handle notifications (no response)
	if id == nil {
		switch method {
		case "notifications/cancelled":
			log.Printf("[MCP Server] Received cancellation notification")
		}
		return nil
	}

	switch method {
	case "initialize":
		return s.handleInitialize(id)
	case "ping":
		return s.handlePing(id)
	case "tools/list":
		return s.handleToolsList(id)
	case "tools/call":
		return s.handleToolsCall(id, msg)
	case "resources/list":
		return s.handleResourcesList(id)
	case "prompts/list":
		return s.handlePromptsList(id)
	default:
		return s.makeError(id, -32601, fmt.Sprintf("Method not found: %s", method))
	}
}

// handleInitialize responds to the MCP initialize request.
func (s *MCPServer) handleInitialize(id interface{}) map[string]interface{} {
	return map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"result": map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{
					"listChanged": false,
				},
			},
			"serverInfo": map[string]interface{}{
				"name":    "hyperharness",
				"version": "0.2.0",
			},
		},
	}
}

// handlePing responds to ping requests.
func (s *MCPServer) handlePing(id interface{}) map[string]interface{} {
	return map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"result":  map[string]interface{}{},
	}
}

// handleToolsList returns all registered tools.
func (s *MCPServer) handleToolsList(id interface{}) map[string]interface{} {
	s.mu.RLock()
	tools := make([]interface{}, 0, len(s.tools))
	for _, t := range s.tools {
		schema := t.Schema
		if schema == nil {
			schema = json.RawMessage(`{"type":"object","properties":{}}`)
		}
		tools = append(tools, map[string]interface{}{
			"name":        t.Name,
			"description": t.Description,
			"inputSchema": json.RawMessage(schema),
		})
	}
	s.mu.RUnlock()

	return map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"result": map[string]interface{}{
			"tools": tools,
		},
	}
}

// handleToolsCall executes a tool call.
func (s *MCPServer) handleToolsCall(id interface{}, msg map[string]interface{}) map[string]interface{} {
	params, _ := msg["params"].(map[string]interface{})
	toolName, _ := params["name"].(string)
	args, _ := params["arguments"].(map[string]interface{})

	if toolName == "" {
		return s.makeError(id, -32602, "Missing tool name")
	}

	// Check if tool exists
	s.mu.RLock()
	_, exists := s.tools[toolName]
	s.mu.RUnlock()

	if !exists {
		return s.makeError(id, -32602, fmt.Sprintf("Unknown tool: %s", toolName))
	}

	// Execute the tool
	result, err := s.handler(toolName, args)
	if err != nil {
		return map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      id,
			"result": map[string]interface{}{
				"content": []interface{}{
					map[string]interface{}{
						"type": "text",
						"text": fmt.Sprintf("Error: %v", err),
					},
				},
				"isError": true,
			},
		}
	}

	return map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"result": map[string]interface{}{
			"content": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": result,
				},
			},
			"isError": false,
		},
	}
}

// handleResourcesList returns available resources (empty for now).
func (s *MCPServer) handleResourcesList(id interface{}) map[string]interface{} {
	return map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"result": map[string]interface{}{
			"resources": []interface{}{},
		},
	}
}

// handlePromptsList returns available prompts (empty for now).
func (s *MCPServer) handlePromptsList(id interface{}) map[string]interface{} {
	return map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"result": map[string]interface{}{
			"prompts": []interface{}{},
		},
	}
}

// makeError creates a JSON-RPC error response.
func (s *MCPServer) makeError(id interface{}, code int, message string) map[string]interface{} {
	return map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
		},
	}
}

// ExportToolCatalog exports all tools as an MCP tool catalog JSON.
func (s *MCPServer) ExportToolCatalog() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	catalog := map[string]interface{}{
		"version":  "2024-11-05",
		"server":   "hyperharness",
		"toolCount": len(s.tools),
		"tools":    s.tools,
	}

	data, _ := json.MarshalIndent(catalog, "", "  ")
	return string(data)
}

// BuildToolCatalogFromRegistry builds MCP tool definitions from the tools registry.
// This is the bridge between the internal tool registry and the MCP server.
func BuildToolCatalogFromRegistry(toolList []struct {
	Name        string
	Description string
	Schema      string
}) []MCPToolDef {
	defs := make([]MCPToolDef, 0, len(toolList))
	for _, t := range toolList {
		schema := json.RawMessage(t.Schema)
		if len(schema) == 0 {
			schema = json.RawMessage(`{"type":"object","properties":{}}`)
		}
		// Normalize schema - ensure it's valid JSON
		if !json.Valid(schema) {
			schema = json.RawMessage(`{"type":"object","properties":{}}`)
		}
		defs = append(defs, MCPToolDef{
			Name:        t.Name,
			Description: strings.TrimSpace(t.Description),
			Schema:      schema,
		})
	}
	return defs
}

// Ensure os and log are used
var _ = os.Getenv
var _ = log.Printf
