package mcp

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func containsStr(s, sub string) bool {
	return strings.Contains(s, sub)
}

// MockStdioClient is a test client that doesn't actually spawn a process.
type MockStdioClient struct {
	Name      string
	Tools     []MCPTool
	Resources []MCPResource
	Prompts   []MCPPrompt
	running   bool
}

func (m *MockStdioClient) Start() error                     { m.running = true; return nil }
func (m *MockStdioClient) Stop() error                      { m.running = false; return nil }
func (m *MockStdioClient) IsRunning() bool                  { return m.running }
func (m *MockStdioClient) ToolCount() int                   { return len(m.Tools) }
func (m *MockStdioClient) Initialize(ctx context.Context) error { return nil }

func TestAggregatorCreation(t *testing.T) {
	agg := NewAggregator()
	if agg == nil {
		t.Fatal("aggregator should not be nil")
	}
	if agg.ServerCount() != 0 {
		t.Error("new aggregator should have no servers")
	}
	if agg.TotalToolCount() != 0 {
		t.Error("new aggregator should have no tools")
	}
}

func TestAggregatorStatus(t *testing.T) {
	agg := NewAggregator()
	status := agg.Status()
	if status == nil {
		t.Fatal("status should not be nil")
	}
	if status["serverCount"] != 0 {
		t.Error("serverCount should be 0")
	}
}

func TestAggregatorStatusJSON(t *testing.T) {
	agg := NewAggregator()
	data, err := agg.StatusJSON()
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("status JSON should not be empty")
	}
}

func TestAggregatorHasToolEmpty(t *testing.T) {
	agg := NewAggregator()
	if agg.HasTool("anything") {
		t.Error("empty aggregator should not have any tools")
	}
}

func TestAggregatorServerNamesEmpty(t *testing.T) {
	agg := NewAggregator()
	names := agg.ServerNames()
	if len(names) != 0 {
		t.Error("empty aggregator should have no server names")
	}
}

func TestAggregatorListToolsEmpty(t *testing.T) {
	agg := NewAggregator()
	tools, err := agg.ListTools(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(tools) != 0 {
		t.Error("empty aggregator should list no tools")
	}
}

func TestJSONRPCRequestSerialization(t *testing.T) {
	req := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/list",
		Params:  nil,
	}
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	if !containsStr(string(data), `"jsonrpc":"2.0"`) {
		t.Errorf("missing jsonrpc version: %s", data)
	}
}

func TestJSONRPCResponseDeserialization(t *testing.T) {
	raw := `{"jsonrpc":"2.0","id":1,"result":{"tools":[{"name":"test","description":"test tool"}]}}`
	var resp JSONRPCResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.JSONRPC != "2.0" {
		t.Error("wrong jsonrpc version")
	}
	if resp.Error != nil {
		t.Error("should not have error")
	}
}

func TestRPCError(t *testing.T) {
	err := &RPCError{Code: -32600, Message: "Invalid Request"}
	if err.Error() != "RPC error -32600: Invalid Request" {
		t.Errorf("unexpected error string: %s", err.Error())
	}
}

func TestMCPToolSerialization(t *testing.T) {
	tool := MCPTool{
		Name:        "read_file",
		Description: "Read a file",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"path": map[string]interface{}{"type": "string"},
			},
		},
	}
	data, err := json.Marshal(tool)
	if err != nil {
		t.Fatal(err)
	}
	if !containsStr(string(data), "read_file") {
		t.Errorf("missing tool name: %s", data)
	}
}

func TestStdioClientCreation(t *testing.T) {
	client := NewStdioClient("test", "echo", []string{}, nil)
	if client.Name != "test" {
		t.Error("name mismatch")
	}
	if client.Command != "echo" {
		t.Error("command mismatch")
	}
	if client.IsRunning() {
		t.Error("should not be running before start")
	}
}

func TestStdioClientToolCount(t *testing.T) {
	client := NewStdioClient("test", "echo", nil, nil)
	client.Tools = []MCPTool{
		{Name: "tool1"},
		{Name: "tool2"},
	}
	if client.ToolCount() != 2 {
		t.Errorf("expected 2 tools, got %d", client.ToolCount())
	}
}

func TestAggregatorShutdown(t *testing.T) {
	agg := NewAggregator()
	agg.Shutdown()
	// Should be safe to call multiple times
	agg.Shutdown()
}

func TestContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	<-ctx.Done()
	if ctx.Err() == nil {
		t.Error("context should have timed out")
	}
}
