package mcp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRegistryNew(t *testing.T) {
	r := NewRegistry()
	if r == nil {
		t.Fatal("registry should not be nil")
	}
}

func TestRegistryRegisterServer(t *testing.T) {
	r := NewRegistry()
	server := &MCPServer{
		Name:      "test-server",
		Command:   "test-cmd",
		Transport: "stdio",
	}

	if err := r.RegisterServer(server); err != nil {
		t.Fatalf("register failed: %v", err)
	}
}

func TestRegistryBuiltInTool(t *testing.T) {
	r := NewRegistry()
	called := false
	r.RegisterBuiltInTool("test-tool", func(args map[string]interface{}) (ToolResult, error) {
		called = true
		return ToolResult{
			Content: []ToolContent{{Type: "text", Text: "result"}},
		}, nil
	})

	tools := r.ListTools()
	if len(tools) != 1 {
		t.Errorf("expected 1 tool, got %d", len(tools))
	}

	result, err := r.CallTool("test-tool", map[string]interface{}{})
	if err != nil {
		t.Fatalf("call failed: %v", err)
	}
	if !called {
		t.Error("tool was not called")
	}
	if result.Content[0].Text != "result" {
		t.Errorf("result text: %s", result.Content[0].Text)
	}
}

func TestRegistryCallUnknownTool(t *testing.T) {
	r := NewRegistry()
	_, err := r.CallTool("nonexistent", nil)
	if err == nil {
		t.Error("expected error for unknown tool")
	}
}

func TestRegistryListTools(t *testing.T) {
	r := NewRegistry()
	r.RegisterBuiltInTool("tool1", func(args map[string]interface{}) (ToolResult, error) {
		return ToolResult{}, nil
	})
	r.RegisterBuiltInTool("tool2", func(args map[string]interface{}) (ToolResult, error) {
		return ToolResult{}, nil
	})

	tools := r.ListTools()
	if len(tools) != 2 {
		t.Errorf("expected 2 tools, got %d", len(tools))
	}
}

func TestRegistryGetInventory(t *testing.T) {
	r := NewRegistry()
	r.RegisterBuiltInTool("built-in", func(args map[string]interface{}) (ToolResult, error) {
		return ToolResult{}, nil
	})
	r.RegisterServer(&MCPServer{Name: "ext-server", Transport: "stdio"})

	inventory := r.GetToolInventory()
	if inventory["totalBuiltIn"] != 1 {
		t.Errorf("built-in count: %v", inventory["totalBuiltIn"])
	}
	if inventory["total"] != 1 {
		t.Errorf("total count: %v", inventory["total"])
	}
}

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "mcp.json")

	config := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"demo": map[string]interface{}{
				"command":   "demo-cmd",
				"args":      []string{"arg1"},
				"transport": "stdio",
			},
		},
	}

	data, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile(configPath, data, 0o644)

	r := NewRegistry()
	if err := r.LoadFromFile(configPath); err != nil {
		t.Fatalf("load failed: %v", err)
	}

	inventory := r.GetToolInventory()
	servers, _ := inventory["servers"].(map[string]interface{})
	if servers == nil {
		t.Fatal("servers should not be nil")
	}
	demo, ok := servers["demo"]
	if !ok {
		t.Error("demo server should be registered")
	}
	_ = demo
}

func TestLoadFromFileNotFound(t *testing.T) {
	r := NewRegistry()
	err := r.LoadFromFile("/nonexistent/path/mcp.json")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestJSONRPCMessage(t *testing.T) {
	msg := JSONRPCMessage{
		JSONRPC: "2.0",
		Method:  "initialize",
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}

	var decoded JSONRPCMessage
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}

	if decoded.JSONRPC != "2.0" {
		t.Errorf("jsonrpc: %s", decoded.JSONRPC)
	}
	if decoded.Method != "initialize" {
		t.Errorf("method: %s", decoded.Method)
	}
}

func TestToolStruct(t *testing.T) {
	tool := Tool{
		Name:        "read_file",
		Description: "Read file contents",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"path": map[string]string{"type": "string"},
			},
		},
	}

	data, _ := json.Marshal(tool)
	var decoded Tool
	json.Unmarshal(data, &decoded)

	if decoded.Name != "read_file" {
		t.Errorf("name: %s", decoded.Name)
	}
}

func TestToolCallStruct(t *testing.T) {
	tc := ToolCall{
		Name: "bash",
		Arguments: map[string]interface{}{
			"command": "echo hello",
		},
	}

	data, _ := json.Marshal(tc)
	var decoded ToolCall
	json.Unmarshal(data, &decoded)

	if decoded.Name != "bash" {
		t.Errorf("name: %s", decoded.Name)
	}
}

func TestToolResultStruct(t *testing.T) {
	result := ToolResult{
		Content: []ToolContent{
			{Type: "text", Text: "hello"},
			{Type: "image", Text: "base64..."},
		},
		IsError: false,
	}

	data, _ := json.Marshal(result)
	var decoded ToolResult
	json.Unmarshal(data, &decoded)

	if len(decoded.Content) != 2 {
		t.Errorf("content count: %d", len(decoded.Content))
	}
	if decoded.Content[0].Text != "hello" {
		t.Errorf("text: %s", decoded.Content[0].Text)
	}
}

func TestProtocolVersion(t *testing.T) {
	if ProtocolVersion == "" {
		t.Error("protocol version should not be empty")
	}
}
