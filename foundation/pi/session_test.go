package pi

import "testing"

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
