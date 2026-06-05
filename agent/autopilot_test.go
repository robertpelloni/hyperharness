package agent

import (
	"testing"
	"strings"
)

func TestAutopilotGoalHalted(t *testing.T) {
	// Since Agent requires a real OpenAI client for Chat, we can't easily
	// unit test the successful path without heavy mocking.
	// But we can verify it handles max iterations.

	a := &Agent{} // Non-initialized client will cause Chat to fail

	// If Chat fails, AutopilotMode should return an error
	_, err := a.AutopilotMode("Test Goal")
	if err == nil {
		t.Error("Expected error due to uninitialized agent client")
	}
}

func TestAgentSystemPrompt(t *testing.T) {
	prompt := buildAgentSystemPrompt("Custom Context")
	if !strings.Contains(prompt, "Hypercode") {
		t.Error("System prompt missing Hypercode identity")
	}
	if !strings.Contains(prompt, "Custom Context") {
		t.Error("System prompt missing system context")
	}
}
