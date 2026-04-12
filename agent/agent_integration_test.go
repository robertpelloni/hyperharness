package agent

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/robertpelloni/hyperharness/internal/providers"
	"github.com/robertpelloni/hyperharness/internal/sessions"
	"github.com/robertpelloni/hyperharness/llm"
)

// MockProvider implements a mock LLM provider for testing the agent loop.
type MockProvider struct {
	Response        string
	ToolCalls       []llm.ToolCall
	CallCount       int
	LastSystem      string
	LastMessages    []llm.Message
	LastTools       []llm.Tool
	StreamingBlocks bool
}

func (m *MockProvider) GenerateResponse(ctx context.Context, req llm.Request) (*llm.Response, error) {
	m.CallCount++
	m.LastSystem = req.SystemPrompt
	m.LastMessages = req.Messages
	m.LastTools = req.Tools

	if m.StreamingBlocks {
		// Just a stub for non-streaming behavior when streaming is requested,
		// in a real mock we'd push to channels, but for simplicity of the interface
		// we return the final response
	}

	return &llm.Response{
		Text:      m.Response,
		ToolCalls: m.ToolCalls,
		ModelInfo: map[string]interface{}{"model": "mock-model"},
		Usage: llm.TokenUsage{
			PromptTokens:     10,
			CompletionTokens: 20,
			TotalTokens:      30,
		},
	}, nil
}

func (m *MockProvider) StreamResponse(ctx context.Context, req llm.Request) (<-chan llm.StreamChunk, error) {
	m.CallCount++
	m.LastSystem = req.SystemPrompt
	m.LastMessages = req.Messages
	m.LastTools = req.Tools

	ch := make(chan llm.StreamChunk, 2)
	go func() {
		defer close(ch)

		// Send text
		ch <- llm.StreamChunk{Text: m.Response}

		// Send tools
		for _, tc := range m.ToolCalls {
			ch <- llm.StreamChunk{ToolCall: &tc}
		}
	}()

	return ch, nil
}

func (m *MockProvider) GetModels() []providers.ModelInfo {
	return []providers.ModelInfo{{ID: "mock-model"}}
}
func (m *MockProvider) CheckHealth() error { return nil }

func TestAgentLoopIntegration(t *testing.T) {
	// Setup mock provider
	mockProv := &MockProvider{
		Response: "I will use the read tool to check that file.",
		ToolCalls: []llm.ToolCall{
			{
				ID:   "call_123",
				Name: "read",
				Args: `{"path": "test.txt"}`,
			},
		},
	}

	// Register the mock provider
	llm.ProviderRegistry.RegisterProvider("mock", func() providers.Provider { return mockProv })

	// Create a session
	sess := sessions.NewSession(".", "mock", "mock-model")

	// Create the agent
	agt := NewAgent(sess, AgentConfig{
		MaxTurns:    3,
		MaxTokens:   1000,
		Temperature: 0.7,
	})

	// Run the agent loop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := agt.Run(ctx, "Read the test file")
	if err != nil {
		t.Fatalf("Agent loop failed: %v", err)
	}

	// Verify interactions
	if mockProv.CallCount == 0 {
		t.Errorf("Expected provider to be called, got 0 calls")
	}

	// Since the read tool isn't actually executing against real test files in this isolated test,
	// it might fail or return a not found error. We just care that the loop wired correctly and called tools.

	foundToolResult := false
	for _, entry := range sess.Entries {
		if entry.Type == sessions.EntryToolResult && strings.Contains(entry.Content, "read") {
			foundToolResult = true
			break
		}
	}

	// Tool execution happens via the pi foundation tools, so the "read" tool should have been invoked
	// and logged a tool result entry in the session.
	if !foundToolResult {
		t.Errorf("Expected to find a tool result entry for the read tool in the session history")
	}
}
