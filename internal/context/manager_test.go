package context

import (
	"strings"
	"testing"
)

func TestManagerAddMessage(t *testing.T) {
	m := NewManager(100, 200000)
	m.AddMessage("user", "hello")
	m.AddMessage("assistant", "hi there")

	msgs := m.GetMessages()
	if len(msgs) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(msgs))
	}
	if msgs[0].Role != "user" || msgs[0].Content != "hello" {
		t.Errorf("first message mismatch: %+v", msgs[0])
	}
	if msgs[1].Role != "assistant" || msgs[1].Content != "hi there" {
		t.Errorf("second message mismatch: %+v", msgs[1])
	}
}

func TestManagerCompact(t *testing.T) {
	m := NewManager(5, 200000)

	for i := 0; i < 10; i++ {
		m.AddMessage("user", "message "+strings.Repeat("x", 100))
	}

	if len(m.GetMessages()) != 10 {
		t.Fatalf("expected 10 messages before compact")
	}

	removed := m.Compact(5)
	if removed != 5 {
		t.Errorf("expected 5 removed, got %d", removed)
	}
	if len(m.GetMessages()) != 5 {
		t.Errorf("expected 5 messages after compact, got %d", len(m.GetMessages()))
	}
}

func TestManagerCompactToFit(t *testing.T) {
	m := NewManager(100, 500) // 500 token budget

	for i := 0; i < 20; i++ {
		m.AddMessage("user", strings.Repeat("word ", 50)) // ~50 tokens each
	}

	removed := m.CompactToFit(500)
	if removed == 0 {
		t.Error("expected some messages to be removed")
	}

	status := m.Status()
	if status.TotalTokens > 500+200 { // allow some margin
		t.Errorf("total tokens %d exceeds budget", status.TotalTokens)
	}
}

func TestManagerInject(t *testing.T) {
	m := NewManager(100, 200000)
	m.SetSystemPrompt("You are helpful.")
	m.Inject("Project: test-project")
	m.Inject("Language: Go")

	providerMsgs := m.GetMessagesForProvider()
	if len(providerMsgs) < 2 {
		t.Fatalf("expected at least 2 messages (system + context), got %d", len(providerMsgs))
	}

	// First should be system prompt
	if providerMsgs[0].Role != "system" || !strings.Contains(providerMsgs[0].Content, "helpful") {
		t.Errorf("first message should be system prompt: %+v", providerMsgs[0])
	}

	// Second should be injected context
	if providerMsgs[1].Role != "system" {
		t.Errorf("second message should be system: %+v", providerMsgs[1])
	}
	if !strings.Contains(providerMsgs[1].Content, "test-project") {
		t.Errorf("injected context missing project: %s", providerMsgs[1].Content)
	}
}

func TestManagerStatus(t *testing.T) {
	m := NewManager(100, 200000)
	m.AddMessage("user", "hello")
	m.AddMessage("assistant", "world")
	m.Inject("extra context")

	status := m.Status()
	if status.TotalMessages != 2 {
		t.Errorf("expected 2 messages, got %d", status.TotalMessages)
	}
	if status.MaxMessages != 100 {
		t.Errorf("expected max 100, got %d", status.MaxMessages)
	}
	if len(status.InjectedContext) != 1 {
		t.Errorf("expected 1 injected context, got %d", len(status.InjectedContext))
	}
}

func TestManagerClear(t *testing.T) {
	m := NewManager(100, 200000)
	m.AddMessage("user", "hello")
	m.Inject("context")

	m.Clear()

	if len(m.GetMessages()) != 0 {
		t.Error("messages should be empty after clear")
	}
	status := m.Status()
	if len(status.InjectedContext) != 0 {
		t.Error("injected context should be empty after clear")
	}
}

func TestManagerSummarize(t *testing.T) {
	m := NewManager(100, 200000)
	m.AddMessage("user", "do something")
	m.AddMessage("assistant", "done")
	m.AddToolMessage("tc1", "bash", "command output")

	summary := m.Summarize()
	if !strings.Contains(summary, "User messages:") {
		t.Errorf("summary missing user count: %s", summary)
	}
	if !strings.Contains(summary, "Tool calls:") {
		t.Errorf("summary missing tool count: %s", summary)
	}
	if !strings.Contains(summary, "done") {
		t.Errorf("summary missing last assistant msg: %s", summary)
	}
}

func TestManagerToolMessage(t *testing.T) {
	m := NewManager(100, 200000)
	m.AddToolMessage("call_123", "bash", "output here")

	msgs := m.GetMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	if msgs[0].Role != "tool" {
		t.Errorf("expected tool role, got %s", msgs[0].Role)
	}
	if msgs[0].ToolCallID != "call_123" {
		t.Errorf("tool call ID mismatch: %s", msgs[0].ToolCallID)
	}
	if msgs[0].Name != "bash" {
		t.Errorf("tool name mismatch: %s", msgs[0].Name)
	}
}

func TestEstimateTokens(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"", 0},
		{"a", 0},     // 1/4 = 0
		{"abcd", 1},  // 4/4 = 1
		{"abcdefgh", 2}, // 8/4 = 2
	}

	for _, tc := range tests {
		result := estimateTokens(tc.input)
		if result != tc.expected {
			t.Errorf("estimateTokens(%q) = %d, want %d", tc.input, result, tc.expected)
		}
	}
}
