package tools

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/robertpelloni/hyperharness/foundation/compat"
	foundationpi "github.com/robertpelloni/hyperharness/foundation/pi"
)

func TestRegistryFind(t *testing.T) {
	var nilRegistry *Registry
	if _, ok := nilRegistry.Find("read"); ok {
		t.Fatal("expected nil registry lookup to fail")
	}

	registry := &Registry{Tools: []Tool{{Name: "read"}, {Name: "bash"}}}
	if _, ok := registry.Find(" "); ok {
		t.Fatal("expected blank lookup to fail")
	}
	tool, ok := registry.Find("bash")
	if !ok || tool.Name != "bash" {
		t.Fatalf("expected bash lookup, got %#v ok=%v", tool, ok)
	}
	if _, ok := registry.Find("missing"); ok {
		t.Fatal("expected missing lookup to fail")
	}
}

func TestFormatFoundationToolResult(t *testing.T) {
	if got := formatFoundationToolResult(nil); got != "" {
		t.Fatalf("expected empty output for nil result, got %q", got)
	}

	plain := &foundationpi.ToolResult{ToolName: "read", Content: []any{foundationpi.TextContent{Type: "text", Text: "hello"}}}
	if got := formatFoundationToolResult(plain); got != "hello" {
		t.Fatalf("expected plain text formatting, got %q", got)
	}

	detailed := &foundationpi.ToolResult{ToolName: "read", Content: []any{foundationpi.TextContent{Type: "text", Text: "hello"}}, Details: &foundationpi.ReadToolDetails{}}
	if got := formatFoundationToolResult(detailed); !strings.Contains(got, `"toolName": "read"`) {
		t.Fatalf("expected JSON formatting for detailed result, got %q", got)
	}
}

func TestNewFoundationToolUsesContractCopyAndRuntime(t *testing.T) {
	runtime := foundationpi.NewRuntime(t.TempDir(), nil)
	contract := compat.ToolContract{
		Name:        "write",
		Description: "write file",
		Parameters:  json.RawMessage(`{"type":"object"}`),
	}
	tool := newFoundationTool(runtime, contract)
	if tool.Name != "write" || tool.Description != "write file" {
		t.Fatalf("unexpected tool metadata: %#v", tool)
	}
	contract.Name = "mutated"
	contract.Description = "mutated"
	if tool.Name != "write" || tool.Description != "write file" {
		t.Fatalf("expected tool metadata to remain cloned, got %#v", tool)
	}
	out, err := tool.Execute(map[string]interface{}{"path": "demo.txt", "content": "hello"})
	if err != nil {
		t.Fatalf("expected foundation tool execution to succeed, got %v", err)
	}
	if !strings.Contains(out, "Successfully wrote") {
		t.Fatalf("unexpected tool output: %q", out)
	}
}

func TestRegistryFindReadToolUsesFoundationBehavior(t *testing.T) {
	registry := NewRegistry()
	tool, ok := registry.Find("read")
	if !ok {
		t.Fatal("missing read tool")
	}
	output, err := tool.Execute(map[string]interface{}{"path": filepath.Join("..", "go.mod")})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output, "module ") {
		t.Fatalf("unexpected read output: %q", output)
	}
	var schema map[string]any
	if err := json.Unmarshal(tool.Parameters, &schema); err != nil {
		t.Fatalf("invalid parameter schema: %v", err)
	}
}
