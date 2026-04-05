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
