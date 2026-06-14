package summarizer

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/robertpelloni/hyperharness/internal/ai"
)

// mockProvider simulates an LLM for testing
type mockProvider struct {
	response string
	err      error
}

func (m *mockProvider) GenerateText(ctx context.Context, model string, messages []ai.Message) (*ai.LLMResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &ai.LLMResponse{
		Content:  m.response,
		Provider: "mock",
		Model:    model,
	}, nil
}

func (m *mockProvider) StreamChat(ctx context.Context, model string, messages []ai.Message, callback func(string) error) (*ai.LLMResponse, error) {
	return m.GenerateText(ctx, model, messages)
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Model == "" {
		t.Error("DefaultConfig() returned empty model")
	}
	if cfg.MaxTokens <= 0 {
		t.Errorf("DefaultConfig() MaxTokens = %d, want > 0", cfg.MaxTokens)
	}
	if cfg.TriggerPct <= 0 || cfg.TriggerPct > 100 {
		t.Errorf("DefaultConfig() TriggerPct = %f, want between 0-100", cfg.TriggerPct)
	}
	if cfg.SummaryRatio <= 0 || cfg.SummaryRatio > 1 {
		t.Errorf("DefaultConfig() SummaryRatio = %f, want between 0-1", cfg.SummaryRatio)
	}
}

func TestNewSummarizer(t *testing.T) {
	provider := &mockProvider{response: "test summary"}
	cfg := DefaultConfig()
	s := NewSummarizer(provider, cfg)
	if s == nil {
		t.Fatal("NewSummarizer() returned nil")
	}
	if len(s.GetHistory()) != 0 {
		t.Error("New summarizer should have empty history")
	}
}

func TestSummarizeContext_Empty(t *testing.T) {
	s := NewSummarizer(&mockProvider{}, DefaultConfig())
	summary, err := s.SummarizeContext(context.Background(), []ai.Message{})
	if err != nil {
		t.Fatalf("SummarizeContext() error = %v", err)
	}
	if summary.MessageCount != 0 {
		t.Errorf("summary.MessageCount = %d, want 0", summary.MessageCount)
	}
}

func TestSummarizeContext_Success(t *testing.T) {
	mockResp := `- Decided to rewrite the auth module
- Fixed login timeout issue
- Modified auth/login.go and auth/session.go
- Still need to review the logout flow`

	s := NewSummarizer(&mockProvider{response: mockResp}, DefaultConfig())
	msgs := []ai.Message{
		{Role: "user", Content: "We need to fix the login issue"},
		{Role: "assistant", Content: "Let me look at the auth module"},
		{Role: "user", Content: "The timeout is 30s, should be 60s"},
		{Role: "assistant", Content: "I'll update the configuration and tests"},
	}

	summary, err := s.SummarizeContext(context.Background(), msgs)
	if err != nil {
		t.Fatalf("SummarizeContext() error = %v", err)
	}

	if summary.Content == "" {
		t.Error("summary.Content is empty")
	}
	if summary.MessageCount != len(msgs) {
		t.Errorf("summary.MessageCount = %d, want %d", summary.MessageCount, len(msgs))
	}
	if summary.TokenCount <= 0 {
		t.Errorf("summary.TokenCount = %d, want > 0", summary.TokenCount)
	}
}

func TestParseSummary_Structured(t *testing.T) {
	s := NewSummarizer(&mockProvider{response: ""}, DefaultConfig())
	content := `Here is the summary:

Key Points:
- Decision 1: Use Go modules
- Decision 2: Adopt gRPC

Tasks Completed:
- Implemented the API layer
- Wrote integration tests

Files Modified:
- cmd/server.go
- internal/api/handler.go`

	summary := s.parseSummary(content, 5)
	if len(summary.KeyPoints) != 2 {
		t.Errorf("KeyPoints = %d, want 2", len(summary.KeyPoints))
	}
	if len(summary.TasksCompleted) != 2 {
		t.Errorf("TasksCompleted = %d, want 2", len(summary.TasksCompleted))
	}
	if len(summary.FilesModified) != 2 {
		t.Errorf("FilesModified = %d, want 2", len(summary.FilesModified))
	}
}

func TestParseSummary_Unstructured(t *testing.T) {
	s := NewSummarizer(&mockProvider{response: ""}, DefaultConfig())
	content := "Just a plain text summary without any structured sections."

	summary := s.parseSummary(content, 3)
	if len(summary.KeyPoints) == 0 {
		t.Error("Expected at least one key point with full content fallback")
	}
	if summary.KeyPoints[0] != content {
		t.Errorf("KeyPoints[0] = %q, want full content", summary.KeyPoints[0])
	}
}

func TestShouldTrigger(t *testing.T) {
	s := NewSummarizer(&mockProvider{response: ""}, DefaultConfig())

	tests := []struct {
		name     string
		current  int
		max      int
		want     bool
	}{
		{"below threshold", 50, 100, false}, // 50% < 80% trigger
		{"at threshold", 80, 100, true},    // 80% >= 80% trigger
		{"above threshold", 90, 100, true},  // 90% >= 80% trigger
		{"zero max", 100, 0, false},
		{"negative max", 100, -1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.ShouldTrigger(tt.current, tt.max)
			if got != tt.want {
				t.Errorf("ShouldTrigger(%d, %d) = %v, want %v", tt.current, tt.max, got, tt.want)
			}
		})
	}
}

func TestShouldTrigger_Threshold(t *testing.T) {
	cfg := DefaultConfig()
	cfg.TriggerPct = 80.0
	s := NewSummarizer(&mockProvider{response: ""}, cfg)

	if s.ShouldTrigger(79, 100) {
		t.Error("ShouldTrigger(79, 100) = true, want false (below 80%)")
	}
	if !s.ShouldTrigger(80, 100) {
		t.Error("ShouldTrigger(80, 100) = false, want true (at 80%)")
	}
	if !s.ShouldTrigger(90, 100) {
		t.Error("ShouldTrigger(90, 100) = false, want true (above 80%)")
	}
}

func TestCompactMessages_NoCompactionNeeded(t *testing.T) {
	s := NewSummarizer(&mockProvider{response: "summary"}, DefaultConfig())
	msgs := []ai.Message{
		{Role: "user", Content: "hello"},
		{Role: "assistant", Content: "hi"},
	}

	result, err := s.CompactMessages(context.Background(), msgs, 5)
	if err != nil {
		t.Fatalf("CompactMessages() error = %v", err)
	}
	if len(result) != len(msgs) {
		t.Errorf("result length = %d, want %d (no compaction needed)", len(result), len(msgs))
	}
}

func TestCompactMessages_WithCompaction(t *testing.T) {
	mockResp := "- Summarized the old conversation"
	s := NewSummarizer(&mockProvider{response: mockResp}, DefaultConfig())

	msgs := make([]ai.Message, 10)
	for i := range msgs {
		msgs[i] = ai.Message{
			Role:    "user",
			Content: strings.Repeat("message ", 50),
		}
	}

	result, err := s.CompactMessages(context.Background(), msgs, 3)
	if err != nil {
		t.Fatalf("CompactMessages() error = %v", err)
	}

	// Should have: 1 summary message + 3 recent = 4 total
	if len(result) != 4 {
		t.Errorf("result length = %d, want 4 (1 summary + 3 recent)", len(result))
	}
	if result[0].Role != "system" {
		t.Errorf("first message role = %q, want 'system'", result[0].Role)
	}
	if !strings.Contains(result[0].Content, mockResp) {
		t.Error("summary message doesn't contain mock response")
	}
}

func TestCompactMessages_FallbackOnError(t *testing.T) {
	s := NewSummarizer(&mockProvider{err: fmt.Errorf("API error")}, DefaultConfig())

	msgs := make([]ai.Message, 10)
	for i := range msgs {
		msgs[i] = ai.Message{Role: "user", Content: "test"}
	}

	result, err := s.CompactMessages(context.Background(), msgs, 3)
	if err != nil {
		t.Fatalf("CompactMessages() should not propagate error: %v", err)
	}
	// Should fall back to simple truncation (3 recent)
	if len(result) != 3 {
		t.Errorf("fallback result length = %d, want 3", len(result))
	}
}

func TestGetHistory(t *testing.T) {
	s := NewSummarizer(&mockProvider{response: "summary1"}, DefaultConfig())

	// Create first summary
	_, err := s.SummarizeContext(context.Background(), []ai.Message{
		{Role: "user", Content: "msg1"},
	})
	if err != nil {
		t.Fatalf("first SummarizeContext() error = %v", err)
	}

	if len(s.GetHistory()) != 1 {
		t.Errorf("history length = %d, want 1", len(s.GetHistory()))
	}

	// Clear history
	s.ClearHistory()
	if len(s.GetHistory()) != 0 {
		t.Errorf("after ClearHistory(), length = %d, want 0", len(s.GetHistory()))
	}
}

// TestSummarizeContext_LargeMessages ensures we handle truncated content
func TestSummarizeContext_LargeMessages(t *testing.T) {
	mockResp := "- Summarized large conversation"
	s := NewSummarizer(&mockProvider{response: mockResp}, DefaultConfig())

	// Create a message with content > 500 chars
	longContent := strings.Repeat("This is a long message that should be truncated. ", 20)
	msgs := []ai.Message{
		{Role: "user", Content: longContent},
	}

	summary, err := s.SummarizeContext(context.Background(), msgs)
	if err != nil {
		t.Fatalf("SummarizeContext() error = %v", err)
	}

	if summary.MessageCount != 1 {
		t.Errorf("summary.MessageCount = %d, want 1", summary.MessageCount)
	}
}

func TestSummary_Immutability(t *testing.T) {
	s := NewSummarizer(&mockProvider{response: "- point 1\n- point 2"}, DefaultConfig())

	_, err := s.SummarizeContext(context.Background(), []ai.Message{
		{Role: "user", Content: "test"},
	})
	if err != nil {
		t.Fatalf("SummarizeContext() error = %v", err)
	}

	history := s.GetHistory()
	history[0] = Summary{} // Try to modify the returned slice

	// Original should be unchanged
	originalHistory := s.GetHistory()
	if originalHistory[0].TokenCount == 0 {
		t.Error("original history was modified")
	}
}
