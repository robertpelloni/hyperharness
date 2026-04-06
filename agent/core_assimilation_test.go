package agent

import (
	"errors"
	"strings"
	"testing"

	"github.com/robertpelloni/hyperharness/tools"
	"github.com/sashabaranov/go-openai"
)

func TestBuildAgentSystemPromptIncludesCoreGuidance(t *testing.T) {
	prompt := buildAgentSystemPrompt("CTX")
	for _, needle := range []string{
		"You are Hypercode, a Go-native coding and terminal assistant integrated with Borg and HyperCode.",
		"Prefer the exact-name Pi-compatible tools read, write, edit, and bash when solving coding tasks.",
		"Use repomap for repository-wide context when a condensed map would help.",
		"CTX",
	} {
		if !strings.Contains(prompt, needle) {
			t.Fatalf("expected prompt to contain %q, got %q", needle, prompt)
		}
	}
}

func TestFirstChoiceMessageValidation(t *testing.T) {
	if _, err := firstChoiceMessage(openai.ChatCompletionResponse{}); err == nil || !strings.Contains(err.Error(), "no completion choices returned") {
		t.Fatalf("expected empty choice error, got %v", err)
	}
	msg, err := firstChoiceMessage(openai.ChatCompletionResponse{Choices: []openai.ChatCompletionChoice{{Message: openai.ChatCompletionMessage{Content: "ok"}}}})
	if err != nil || msg.Content != "ok" {
		t.Fatalf("unexpected first choice result: %#v %v", msg, err)
	}
}

func TestExecuteToolCallHandlesUnknownAndErrorPaths(t *testing.T) {
	unknown := executeToolCall(nil, openai.ToolCall{Function: openai.FunctionCall{Name: "missing"}})
	if unknown != "Unknown tool: missing" {
		t.Fatalf("unexpected unknown tool result: %q", unknown)
	}

	registry := &tools.Registry{Tools: []tools.Tool{{
		Name: "demo",
		Execute: func(args map[string]interface{}) (string, error) {
			return "partial", errors.New("boom")
		},
	}}}
	result := executeToolCall(registry, openai.ToolCall{Function: openai.FunctionCall{Name: "demo", Arguments: `{"x":1}`}})
	if result != "partial" {
		t.Fatalf("expected partial output to win on error, got %q", result)
	}

	registry = &tools.Registry{Tools: []tools.Tool{{
		Name: "demo",
		Execute: func(args map[string]interface{}) (string, error) {
			return "", errors.New("boom")
		},
	}}}
	result = executeToolCall(registry, openai.ToolCall{Function: openai.FunctionCall{Name: "demo", Arguments: `not-json`}})
	if !strings.Contains(result, "Error executing demo: boom") {
		t.Fatalf("unexpected error result: %q", result)
	}
}

func TestAgentValidationGuards(t *testing.T) {
	var nilAgent *Agent
	if _, err := nilAgent.Chat("hello"); err == nil || !strings.Contains(err.Error(), "agent is required") {
		t.Fatalf("expected nil agent error, got %v", err)
	}
	if _, err := nilAgent.handleToolCalls(nil); err == nil || !strings.Contains(err.Error(), "agent is required") {
		t.Fatalf("expected nil agent tool-call error, got %v", err)
	}

	a := &Agent{}
	if _, err := a.Chat("hello"); err == nil || !strings.Contains(err.Error(), "openai client is required") {
		t.Fatalf("expected client error, got %v", err)
	}
	a.client = &openai.Client{}
	if _, err := a.Chat("hello"); err == nil || !strings.Contains(err.Error(), "tool registry is required") {
		t.Fatalf("expected tool registry error, got %v", err)
	}
	if _, err := a.handleToolCalls(nil); err == nil || !strings.Contains(err.Error(), "tool registry is required") {
		t.Fatalf("expected tool registry error in tool calls, got %v", err)
	}
}
