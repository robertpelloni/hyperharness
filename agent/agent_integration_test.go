package agent

import (
	"encoding/json"
	"strings"
	"testing"
	"github.com/sashabaranov/go-openai"
	"github.com/robertpelloni/hyperharness/tools"
)

// TestExecuteToolCallIntegration tests if the tool execution pipeline
// wires correctly from an OpenAI ToolCall object through the tools.Registry.
func TestExecuteToolCallIntegration(t *testing.T) {
	// Setup a clean registry
	registry := tools.NewRegistry()

	// Ensure the "read" tool is registered (which it should be by default)
	tool, ok := registry.Find("read")
	if !ok {
		t.Fatalf("Expected to find the 'read' tool in the registry")
	}

	if tool.Name != "read" {
		t.Errorf("Expected tool name 'read', got %s", tool.Name)
	}

	// Create a mock ToolCall for the read tool
	// Note: We use a non-existent file to ensure it fails predictably
	// but the failure proves the tool was actually executed.
	argsJSON, _ := json.Marshal(map[string]interface{}{
		"path": "nonexistent_integration_test_file.txt",
	})

	tc := openai.ToolCall{
		ID: "call_abc123",
		Type: openai.ToolTypeFunction,
		Function: openai.FunctionCall{
			Name: "read",
			Arguments: string(argsJSON),
		},
	}

	// Execute it
	result := executeToolCall(registry, tc)

	// Verify output
	if !strings.Contains(result, "no such file or directory") && !strings.Contains(result, "cannot find") {
		t.Errorf("Expected result to indicate file not found, but got: %s", result)
	}
}

// TestChatFlowRejectsMissingComponents tests basic validations in the Chat flow.
func TestChatFlowRejectsMissingComponents(t *testing.T) {
	var a *Agent
	_, err := a.Chat("hello")
	if err == nil || !strings.Contains(err.Error(), "agent is required") {
		t.Errorf("Expected error for nil agent")
	}

	a = &Agent{}
	_, err = a.Chat("hello")
	if err == nil || !strings.Contains(err.Error(), "client is required") {
		t.Errorf("Expected error for nil client")
	}
}
