package agent

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sashabaranov/go-openai"
)

func TestTrimHistoryKeepsSystemPromptAndMostRecentMessages(t *testing.T) {
	a := &Agent{
		messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: "system"},
			{Role: openai.ChatMessageRoleUser, Content: "u1"},
			{Role: openai.ChatMessageRoleAssistant, Content: "a1"},
			{Role: openai.ChatMessageRoleUser, Content: "u2"},
			{Role: openai.ChatMessageRoleAssistant, Content: "a2"},
		},
	}

	a.TrimHistory(3)

	if len(a.messages) != 3 {
		t.Fatalf("expected 3 messages after trim, got %d", len(a.messages))
	}
	if a.messages[0].Role != openai.ChatMessageRoleSystem || a.messages[0].Content != "system" {
		t.Fatalf("expected system prompt to be preserved, got %#v", a.messages[0])
	}
	if a.messages[1].Content != "u2" || a.messages[2].Content != "a2" {
		t.Fatalf("expected most recent messages to remain, got %#v", a.messages)
	}
}

func TestTrimHistoryWithNonPositiveLimitIsNoOp(t *testing.T) {
	original := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: "system"},
		{Role: openai.ChatMessageRoleUser, Content: "u1"},
	}
	a := &Agent{messages: append([]openai.ChatCompletionMessage(nil), original...)}

	a.TrimHistory(0)

	if len(a.messages) != len(original) {
		t.Fatalf("expected history length %d, got %d", len(original), len(a.messages))
	}
	for i := range original {
		if a.messages[i].Role != original[i].Role || a.messages[i].Content != original[i].Content {
			t.Fatalf("expected history to remain unchanged, got %#v", a.messages)
		}
	}
}

func TestApplyInlineDiffSkipsUnifiedDiffHeaders(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.txt")
	if err := os.WriteFile(path, []byte("alpha\nbeta\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	a := &Agent{}
	diffContent := strings.Join([]string{
		"--- a/sample.txt",
		"+++ b/sample.txt",
		"@@ -1,2 +1,2 @@",
		"-beta",
		"+gamma",
	}, "\n")
	if err := a.ApplyInlineDiff(path, diffContent); err != nil {
		t.Fatalf("apply inline diff failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	if strings.Contains(text, "--- a/sample.txt") || strings.Contains(text, "+++ b/sample.txt") || strings.Contains(text, "@@") {
		t.Fatalf("expected unified diff headers to be ignored, got %q", text)
	}
	if !strings.Contains(text, "alpha\n") || !strings.Contains(text, "gamma\n") {
		t.Fatalf("expected resulting file to contain retained and added content, got %q", text)
	}
	if strings.Contains(text, "beta\n") {
		t.Fatalf("expected removed line to be absent, got %q", text)
	}
}
