package agent

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/robertpelloni/hyperharness/tools"
	"github.com/sashabaranov/go-openai"
)

// MockOpenAIClient overrides standard client behavior for E2E testing
type MockOpenAIClient struct {
	ToolCallsToReturn []openai.ToolCall
	FinalResponse     string
	CallCount         int
}

func (m *MockOpenAIClient) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	m.CallCount++

	// If it's the first call, return tool calls
	if m.CallCount == 1 && len(m.ToolCallsToReturn) > 0 {
		return openai.ChatCompletionResponse{
			Choices: []openai.ChatCompletionChoice{
				{
					Message: openai.ChatCompletionMessage{
						Role:      openai.ChatMessageRoleAssistant,
						ToolCalls: m.ToolCallsToReturn,
					},
				},
			},
		}, nil
	}

	// Otherwise return final response
	return openai.ChatCompletionResponse{
		Choices: []openai.ChatCompletionChoice{
			{
				Message: openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleAssistant,
					Content: m.FinalResponse,
				},
			},
		},
	}, nil
}

// TestHarnessE2E verifies that the Agent loop correctly executes multi-harness tools
func TestHarnessE2E(t *testing.T) {
	registry := tools.NewRegistry()

	testCases := []struct {
		name       string
		toolName   string
		toolArgs   map[string]interface{}
		expectedIn string
	}{
		{
			name:     "Tabby Completion Workflow",
			toolName: "tabby_completion",
			toolArgs: map[string]interface{}{
				"segments": map[string]interface{}{"prefix": "func main() {"},
			},
			expectedIn: "Tabby Completion",
		},
		{
			name:     "Warp Command Workflow",
			toolName: "RequestCommandOutput",
			toolArgs: map[string]interface{}{
				"command": "echo 'Hello Warp'",
			},
			expectedIn: "Hello Warp",
		},
		{
			name:     "Wave Scrollback Workflow",
			toolName: "TermGetScrollback",
			toolArgs: map[string]interface{}{
				"blockId": "test-block",
			},
			expectedIn: "No session found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			argsJSON, _ := json.Marshal(tc.toolArgs)

			// Setup Agent with manual override (since NewAgent uses env vars and real client)
			a := &Agent{
				tools: registry,
				messages: []openai.ChatCompletionMessage{
					{Role: openai.ChatMessageRoleSystem, Content: "System prompt"},
				},
			}

			// In a real scenario we'd use an interface for openai.Client to make this cleaner,
			// but for this implementation we simulate the handleToolCalls directly.

			toolCalls := []openai.ToolCall{
				{
					ID:   "call_123",
					Type: openai.ToolTypeFunction,
					Function: openai.FunctionCall{
						Name:      tc.toolName,
						Arguments: string(argsJSON),
					},
				},
			}

			// We need to verify handleToolCalls logic without a real network call
			// For testing purposes, we use executeToolCall directly which is the core logic
			result := executeToolCall(a.tools, toolCalls[0], nil)

			if !strings.Contains(strings.ToLower(result), strings.ToLower(tc.expectedIn)) {
				t.Errorf("Workflow %s: expected result containing %q, got: %s", tc.name, tc.expectedIn, result)
			}
		})
	}
}
