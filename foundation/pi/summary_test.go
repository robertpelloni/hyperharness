package pi

import (
	"context"
	"strings"
	"testing"
)

func TestDeterministicSummaryGeneratorProducesStructuredSummary(t *testing.T) {
	prep := &BranchSummaryPreparation{
		TargetID:               "target-1",
		OldLeafID:              "old-1",
		CommonAncestorID:       "ancestor-1",
		EntriesToSummarize:     []SessionEntry{{Kind: "message", Role: "user", Text: "Refactor auth flow safely"}},
		SerializedConversation: "[User]: Refactor auth flow safely",
		FileOps: BranchSummaryFileOps{
			ReadFiles:     []string{"auth.go", "session.go"},
			ModifiedFiles: []string{"auth.go"},
		},
		EstimatedTokens: 42,
		MaxTokens:       128,
	}
	summary, err := NewDeterministicSummaryGenerator().GenerateBranchSummary(context.Background(), prep)
	if err != nil {
		t.Fatal(err)
	}
	for _, needle := range []string{
		"## Goal",
		"## Constraints & Preferences",
		"## Progress",
		"## Key Decisions",
		"## Next Steps",
		"## Critical Context",
		"<read-files>",
		"auth.go",
		"session.go",
		"<modified-files>",
	} {
		if !strings.Contains(summary, needle) {
			t.Fatalf("expected summary to contain %q, got:\n%s", needle, summary)
		}
	}
}

func TestDeterministicCompactionGeneratorProducesStructuredSummary(t *testing.T) {
	prep := &CompactionPreparation{
		LeafID:                 "leaf-1",
		EntriesToSummarize:     []SessionEntry{{Kind: "message", Role: "user", Text: "Older conversation"}},
		SerializedConversation: "[User]: Older conversation",
		FileOps: BranchSummaryFileOps{
			ReadFiles:     []string{"old.go"},
			ModifiedFiles: []string{"old.go", "state.go"},
		},
		EstimatedTokens:  64,
		TokensBefore:     512,
		FirstKeptEntryID: "kept-1",
		KeepRecentTokens: 128,
	}
	summary, err := DeterministicSummaryGenerator{}.GenerateCompactionSummary(context.Background(), prep)
	if err != nil {
		t.Fatal(err)
	}
	for _, needle := range []string{"## Goal", "## Progress", "## Key Decisions", "<read-files>", "old.go", "state.go", "kept-1"} {
		if !strings.Contains(summary, needle) {
			t.Fatalf("expected compaction summary to contain %q, got:\n%s", needle, summary)
		}
	}
}

func TestBranchWithGeneratedSummaryAppendsSummaryEntry(t *testing.T) {
	dir := t.TempDir()
	runtime := NewRuntime(dir, DefaultSessionStore(dir))
	session, err := runtime.CreateSession("alpha")
	if err != nil {
		t.Fatal(err)
	}
	id := session.Metadata.SessionID

	session, err = runtime.AppendUserText(id, "A")
	if err != nil {
		t.Fatal(err)
	}
	aID := session.Metadata.LeafID

	session, err = runtime.AppendUserText(id, "B")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.AppendUserText(id, "C"); err != nil {
		t.Fatal(err)
	}
	currentLeaf, err := runtime.GetLeafID(id)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := runtime.BranchSession(id, aID); err != nil {
		t.Fatal(err)
	}
	session, err = runtime.AppendUserText(id, "E")
	if err != nil {
		t.Fatal(err)
	}
	targetID := session.Metadata.LeafID

	if _, err := runtime.BranchSession(id, currentLeaf); err != nil {
		t.Fatal(err)
	}
	session, summary, err := runtime.BranchWithGeneratedSummary(context.Background(), id, targetID, 128, nil, map[string]any{"readFiles": []string{"a.go"}})
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(summary) == "" {
		t.Fatal("expected generated summary text")
	}
	last := session.Entries[len(session.Entries)-1]
	if last.Kind != "branch_summary" {
		t.Fatalf("expected branch_summary entry, got %#v", last)
	}
	if last.ParentID != targetID {
		t.Fatalf("expected summary to attach to target %q, got %q", targetID, last.ParentID)
	}
	if last.FromID != currentLeaf {
		t.Fatalf("expected fromId %q, got %q", currentLeaf, last.FromID)
	}
	if !strings.Contains(last.Summary, "## Goal") {
		t.Fatalf("expected structured generated summary, got:\n%s", last.Summary)
	}
}

func TestCompactWithGeneratedSummaryAppendsCompactionEntry(t *testing.T) {
	dir := t.TempDir()
	runtime := NewRuntime(dir, DefaultSessionStore(dir))
	session, err := runtime.CreateSession("alpha")
	if err != nil {
		t.Fatal(err)
	}
	id := session.Metadata.SessionID

	if _, err := runtime.AppendUserText(id, strings.Repeat("older ", 40)); err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.AppendUserText(id, strings.Repeat("middle ", 40)); err != nil {
		t.Fatal(err)
	}
	session, err = runtime.AppendUserText(id, strings.Repeat("recent ", 40))
	if err != nil {
		t.Fatal(err)
	}

	session, summary, err := runtime.CompactWithGeneratedSummary(context.Background(), id, 60, nil, map[string]any{"readFiles": []string{"ctx.go"}})
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(summary) == "" {
		t.Fatal("expected generated compaction summary")
	}
	last := session.Entries[len(session.Entries)-1]
	if last.Kind != "compaction" {
		t.Fatalf("expected compaction entry, got %#v", last)
	}
	if !strings.Contains(last.Summary, "## Goal") {
		t.Fatalf("expected structured compaction summary, got:\n%s", last.Summary)
	}
	if last.FirstKeptID == "" {
		t.Fatal("expected compaction to record first kept entry")
	}
	if last.TokensBefore == 0 {
		t.Fatal("expected compaction tokensBefore to be populated")
	}
	ctx, err := runtime.BuildSessionContext(id, "")
	if err != nil {
		t.Fatal(err)
	}
	if !ctx.CompactionUsed {
		t.Fatal("expected context to mark compaction use")
	}
}
