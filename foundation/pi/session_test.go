package pi

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSessionStoreCreateAppendListAndFork(t *testing.T) {
	dir := t.TempDir()
	store := NewSessionStore(dir)
	session, err := store.Create("alpha", "/workspace/project")
	if err != nil {
		t.Fatal(err)
	}
	if session.Metadata.Version != CurrentSessionVersion {
		t.Fatalf("expected session version %d, got %d", CurrentSessionVersion, session.Metadata.Version)
	}
	if _, err := store.AppendEntry(session.Metadata.SessionID, SessionEntry{Kind: "message", Role: "user", Text: "hello"}); err != nil {
		t.Fatal(err)
	}
	loaded, err := store.Load(session.Metadata.SessionID)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(loaded.Entries))
	}
	listed, err := store.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(listed) != 1 || listed[0].SessionID != session.Metadata.SessionID {
		t.Fatalf("unexpected session list: %#v", listed)
	}
	forked, err := store.Fork(session.Metadata.SessionID, loaded.Entries[0].ID, "beta")
	if err != nil {
		t.Fatal(err)
	}
	if forked.Metadata.SessionID == session.Metadata.SessionID {
		t.Fatal("expected forked session to have a new id")
	}
	if forked.Metadata.ParentSession != session.Metadata.SessionID {
		t.Fatalf("expected parentSession to be preserved, got %#v", forked.Metadata.ParentSession)
	}
	if len(forked.Entries) != 1 {
		t.Fatalf("expected forked session to copy entries, got %d", len(forked.Entries))
	}
}

func TestSessionStoreLeafTrackingAndBranching(t *testing.T) {
	dir := t.TempDir()
	store := NewSessionStore(dir)
	session, err := store.Create("alpha", "/workspace/project")
	if err != nil {
		t.Fatal(err)
	}
	id := session.Metadata.SessionID

	session, err = store.AppendEntry(id, SessionEntry{Kind: "message", Role: "user", Text: "root"})
	if err != nil {
		t.Fatal(err)
	}
	rootID := session.Metadata.LeafID
	if rootID == "" {
		t.Fatal("expected leaf after first append")
	}

	session, err = store.AppendEntry(id, SessionEntry{Kind: "message", Role: "assistant", Text: "main branch"})
	if err != nil {
		t.Fatal(err)
	}
	mainLeaf := session.Metadata.LeafID
	if mainLeaf == rootID {
		t.Fatal("expected leaf to advance")
	}

	if _, err := store.Branch(id, rootID); err != nil {
		t.Fatal(err)
	}
	leaf, err := store.GetLeafID(id)
	if err != nil {
		t.Fatal(err)
	}
	if leaf != rootID {
		t.Fatalf("expected leaf to move to root, got %q", leaf)
	}

	session, err = store.AppendEntry(id, SessionEntry{Kind: "message", Role: "assistant", Text: "alternate branch"})
	if err != nil {
		t.Fatal(err)
	}
	altLeaf := session.Metadata.LeafID
	if altLeaf == mainLeaf || altLeaf == rootID {
		t.Fatal("expected a new alternate leaf")
	}
	if session.Entries[len(session.Entries)-1].ParentID != rootID {
		t.Fatalf("expected alternate branch parent %q, got %q", rootID, session.Entries[len(session.Entries)-1].ParentID)
	}

	ctx, err := store.BuildSessionContext(id, "")
	if err != nil {
		t.Fatal(err)
	}
	if ctx.LeafID != altLeaf {
		t.Fatalf("expected context to follow alternate leaf %q, got %q", altLeaf, ctx.LeafID)
	}
	if len(ctx.Entries) != 2 {
		t.Fatalf("expected 2 entries on active branch, got %d", len(ctx.Entries))
	}
	if ctx.Entries[len(ctx.Entries)-1].Text != "alternate branch" {
		t.Fatalf("unexpected active branch tail: %#v", ctx.Entries[len(ctx.Entries)-1])
	}

	if _, err := store.ResetLeaf(id); err != nil {
		t.Fatal(err)
	}
	leaf, err = store.GetLeafID(id)
	if err != nil {
		t.Fatal(err)
	}
	if leaf != session.Entries[len(session.Entries)-1].ID {
		t.Fatalf("expected fallback leaf to latest entry, got %q", leaf)
	}
}

func TestSessionStoreCommonAncestorAndBranchSummary(t *testing.T) {
	dir := t.TempDir()
	store := NewSessionStore(dir)
	session, err := store.Create("alpha", "/workspace/project")
	if err != nil {
		t.Fatal(err)
	}
	id := session.Metadata.SessionID

	session, err = store.AppendEntry(id, SessionEntry{Kind: "message", Role: "user", Text: "A"})
	if err != nil {
		t.Fatal(err)
	}
	aID := session.Metadata.LeafID

	session, err = store.AppendEntry(id, SessionEntry{Kind: "message", Role: "assistant", Text: "B"})
	if err != nil {
		t.Fatal(err)
	}
	bID := session.Metadata.LeafID

	session, err = store.AppendEntry(id, SessionEntry{Kind: "message", Role: "assistant", Text: "C"})
	if err != nil {
		t.Fatal(err)
	}
	oldLeafID := session.Metadata.LeafID

	if _, err := store.Branch(id, aID); err != nil {
		t.Fatal(err)
	}
	session, err = store.AppendEntry(id, SessionEntry{Kind: "message", Role: "assistant", Text: "E"})
	if err != nil {
		t.Fatal(err)
	}
	session, err = store.AppendEntry(id, SessionEntry{Kind: "message", Role: "assistant", Text: "F"})
	if err != nil {
		t.Fatal(err)
	}
	fID := session.Metadata.LeafID

	ancestor, err := store.GetCommonAncestor(id, oldLeafID, fID)
	if err != nil {
		t.Fatal(err)
	}
	if ancestor != aID {
		t.Fatalf("expected common ancestor %q, got %q", aID, ancestor)
	}

	prep, err := store.PrepareBranchSummary(id, bID)
	if err != nil {
		t.Fatal(err)
	}
	if prep.OldLeafID != fID {
		t.Fatalf("expected old leaf %q, got %q", fID, prep.OldLeafID)
	}
	if prep.CommonAncestorID != aID {
		t.Fatalf("expected common ancestor %q, got %q", aID, prep.CommonAncestorID)
	}
	if len(prep.EntriesToSummarize) != 2 {
		t.Fatalf("expected 2 entries to summarize, got %d", len(prep.EntriesToSummarize))
	}
	if prep.EntriesToSummarize[0].Text != "E" || prep.EntriesToSummarize[1].Text != "F" {
		t.Fatalf("unexpected entries to summarize: %#v", prep.EntriesToSummarize)
	}

	session, err = store.BranchWithSummary(id, bID, "summary of E/F", map[string]any{"readFiles": []string{"x.go"}})
	if err != nil {
		t.Fatal(err)
	}
	if session.Metadata.LeafID == bID {
		t.Fatal("expected branch summary append to advance leaf beyond target")
	}
	last := session.Entries[len(session.Entries)-1]
	if last.Kind != "branch_summary" {
		t.Fatalf("expected last entry to be branch_summary, got %#v", last)
	}
	if last.ParentID != bID {
		t.Fatalf("expected branch summary to attach to target %q, got %q", bID, last.ParentID)
	}
	if last.FromID != fID {
		t.Fatalf("expected branch summary from leaf %q, got %q", fID, last.FromID)
	}
	if last.Summary != "summary of E/F" {
		t.Fatalf("unexpected summary: %q", last.Summary)
	}
}

func TestBranchSummaryPreparationBudgetAndFileOps(t *testing.T) {
	dir := t.TempDir()
	store := NewSessionStore(dir)
	session, err := store.Create("alpha", "/workspace/project")
	if err != nil {
		t.Fatal(err)
	}
	id := session.Metadata.SessionID

	session, err = store.AppendEntry(id, SessionEntry{Kind: "message", Role: "user", Text: strings.Repeat("A", 200)})
	if err != nil {
		t.Fatal(err)
	}
	rootID := session.Metadata.LeafID

	session, err = store.AppendEntry(id, SessionEntry{Kind: "tool_call", ToolName: "read", ToolInput: json.RawMessage(`{"path":"a.go"}`)})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.AppendEntry(id, SessionEntry{Kind: "tool_result", Role: "toolResult", ToolName: "read", Result: &ToolResult{ToolName: "read", Content: []any{TextContent{Type: "text", Text: strings.Repeat("output", 300)}}}}); err != nil {
		t.Fatal(err)
	}
	session, err = store.AppendEntry(id, SessionEntry{Kind: "tool_call", ToolName: "edit", ToolInput: json.RawMessage(`{"path":"b.go"}`)})
	if err != nil {
		t.Fatal(err)
	}
	session, err = store.AppendEntry(id, SessionEntry{Kind: "message", Role: "assistant", Text: "long branch tail"})
	if err != nil {
		t.Fatal(err)
	}
	currentLeaf := session.Metadata.LeafID

	if _, err := store.Branch(id, rootID); err != nil {
		t.Fatal(err)
	}
	session, err = store.AppendEntry(id, SessionEntry{Kind: "message", Role: "assistant", Text: "alternate"})
	if err != nil {
		t.Fatal(err)
	}
	targetID := session.Metadata.LeafID

	if _, err := store.Branch(id, currentLeaf); err != nil {
		t.Fatal(err)
	}

	prep, err := store.PrepareBranchSummaryWithBudget(id, targetID, 400)
	if err != nil {
		t.Fatal(err)
	}
	if prep.OldLeafID != currentLeaf {
		t.Fatalf("expected old leaf %q, got %q", currentLeaf, prep.OldLeafID)
	}
	if prep.MaxTokens != 400 {
		t.Fatalf("unexpected max tokens: %d", prep.MaxTokens)
	}
	if prep.EstimatedTokens > 400 && len(prep.EntriesToSummarize) > 1 {
		t.Fatalf("expected trimmed preparation, got estimated tokens %d", prep.EstimatedTokens)
	}
	if prep.SerializedConversation == "" {
		t.Fatal("expected serialized conversation")
	}
	if len(prep.FileOps.ReadFiles) == 0 || prep.FileOps.ReadFiles[0] != "a.go" {
		t.Fatalf("unexpected read files: %#v", prep.FileOps.ReadFiles)
	}
	if len(prep.FileOps.ModifiedFiles) == 0 || prep.FileOps.ModifiedFiles[0] != "b.go" {
		t.Fatalf("unexpected modified files: %#v", prep.FileOps.ModifiedFiles)
	}
	template := DefaultStructuredSummaryTemplate(prep)
	if !strings.Contains(template, "<read-files>") || !strings.Contains(template, "a.go") {
		t.Fatalf("unexpected summary template: %s", template)
	}
}

func TestSessionStoreExtendedEntriesAndContext(t *testing.T) {
	dir := t.TempDir()
	store := NewSessionStore(dir)
	session, err := store.Create("alpha", "/workspace/project")
	if err != nil {
		t.Fatal(err)
	}
	id := session.Metadata.SessionID

	session, err = store.AppendEntry(id, SessionEntry{Kind: "message", Role: "user", Text: "hello"})
	if err != nil {
		t.Fatal(err)
	}
	firstID := session.Entries[len(session.Entries)-1].ID

	if _, err := store.AppendModelChange(id, "anthropic", "claude-sonnet-4"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.AppendThinkingLevelChange(id, "high"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.AppendSessionInfo(id, "Renamed Session"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.AppendLabelChange(id, firstID, "checkpoint-1"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.AppendCustomEntry(id, "my-extension", map[string]any{"count": 42}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.AppendCustomMessage(id, "my-extension", "injected context", true, map[string]any{"visible": true}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.AppendCompaction(id, "summary of earlier work", firstID, 12345, map[string]any{"readFiles": []string{"a.go"}}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.AppendEntry(id, SessionEntry{Kind: "message", Role: "assistant", Text: "after compaction"}); err != nil {
		t.Fatal(err)
	}

	name, err := store.GetSessionName(id)
	if err != nil {
		t.Fatal(err)
	}
	if name != "Renamed Session" {
		t.Fatalf("unexpected session name: %q", name)
	}

	label, err := store.GetLabel(id, firstID)
	if err != nil {
		t.Fatal(err)
	}
	if label != "checkpoint-1" {
		t.Fatalf("unexpected label: %q", label)
	}

	children, err := store.GetChildren(id, firstID)
	if err != nil {
		t.Fatal(err)
	}
	if len(children) == 0 {
		t.Fatal("expected at least one child entry")
	}

	branch, err := store.GetBranch(id, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(branch) == 0 {
		t.Fatal("expected non-empty branch")
	}

	ctx, err := store.BuildSessionContext(id, "")
	if err != nil {
		t.Fatal(err)
	}
	if ctx.Provider != "anthropic" || ctx.ModelID != "claude-sonnet-4" {
		t.Fatalf("unexpected model context: provider=%q model=%q", ctx.Provider, ctx.ModelID)
	}
	if ctx.ThinkingLevel != "high" {
		t.Fatalf("unexpected thinking level: %q", ctx.ThinkingLevel)
	}
	if !ctx.CompactionUsed {
		t.Fatal("expected compaction to be reflected in context")
	}
	if len(ctx.Entries) == 0 {
		t.Fatal("expected context entries")
	}
}
