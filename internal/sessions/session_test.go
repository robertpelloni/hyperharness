package sessions

import (
	"encoding/json"
	"testing"
	"time"
)

func TestManagerCreateSession(t *testing.T) {
	dir := t.TempDir()
	mgr, err := NewManager(dir)
	if err != nil {
		t.Fatalf("new manager failed: %v", err)
	}

	session, err := mgr.CreateSession("/tmp", "openai", "gpt-4o")
	if err != nil {
		t.Fatalf("create session failed: %v", err)
	}

	if session.ID == "" {
		t.Error("session ID should not be empty")
	}
	if session.CWD != "/tmp" {
		t.Errorf("cwd: got %q", session.CWD)
	}
}

func TestManagerGetSession(t *testing.T) {
	dir := t.TempDir()
	mgr, _ := NewManager(dir)

	session, _ := mgr.CreateSession("/tmp", "openai", "gpt-4o")

	found, ok := mgr.GetSession(session.ID)
	if !ok {
		t.Fatal("session should be found")
	}
	if found.ID != session.ID {
		t.Errorf("ID mismatch: got %s, want %s", found.ID, session.ID)
	}

	_, ok = mgr.GetSession("nonexistent")
	if ok {
		t.Error("nonexistent should not be found")
	}
}

func TestManagerListSessions(t *testing.T) {
	dir := t.TempDir()
	mgr, _ := NewManager(dir)

	mgr.CreateSession("/tmp", "openai", "gpt-4o")
	mgr.CreateSession("/tmp", "anthropic", "claude-3")

	sessions := mgr.ListSessions()
	if len(sessions) != 2 {
		t.Errorf("expected 2 sessions, got %d", len(sessions))
	}
}

func TestManagerLoadSession(t *testing.T) {
	dir := t.TempDir()
	mgr, _ := NewManager(dir)

	session, _ := mgr.CreateSession("/tmp", "openai", "gpt-4o")

	loaded, err := mgr.LoadSession(session.ID)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if loaded.ID != session.ID {
		t.Errorf("loaded ID: %s", loaded.ID)
	}
}

func TestSessionJSONRoundtrip(t *testing.T) {
	session := &Session{
		ID:        "test_123",
		CWD:       "/tmp/project",
		Provider:  "openai",
		ModelID:   "gpt-4o",
		CreatedAt: time.Now(),
	}

	data, err := json.Marshal(session)
	if err != nil {
		t.Fatal(err)
	}

	var decoded Session
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}

	if decoded.ID != "test_123" {
		t.Errorf("ID: %s", decoded.ID)
	}
	if decoded.Provider != "openai" {
		t.Errorf("provider: %s", decoded.Provider)
	}
}

func TestManagerSessionDir(t *testing.T) {
	dir := t.TempDir()
	mgr, _ := NewManager(dir)

	if mgr.SessionDir() != dir {
		t.Errorf("session dir: got %q, want %q", mgr.SessionDir(), dir)
	}
}

func TestManagerSave(t *testing.T) {
	dir := t.TempDir()
	mgr, _ := NewManager(dir)

	mgr.CreateSession("/tmp", "openai", "gpt-4o")

	if err := mgr.Save(); err != nil {
		t.Fatalf("save failed: %v", err)
	}
}

func TestManagerEmptyDir(t *testing.T) {
	dir := t.TempDir()
	mgr, _ := NewManager(dir)

	sessions := mgr.ListSessions()
	if len(sessions) != 0 {
		t.Errorf("expected 0 sessions, got %d", len(sessions))
	}
}

func TestEntryTypeValues(t *testing.T) {
	types := map[EntryType]string{
		EntryUser:       "user",
		EntryAssistant:  "assistant",
		EntryToolCall:   "tool_call",
		EntryToolResult: "tool_result",
		EntrySystem:     "system",
		EntryCompaction: "compaction",
		EntryBranch:     "branch",
		EntryLabel:      "label",
	}

	for et, expected := range types {
		if string(et) != expected {
			t.Errorf("entry type %q != %q", et, expected)
		}
	}
}

func TestSessionEntryJSONRoundtrip(t *testing.T) {
	now := time.Now()
	entry := &SessionEntry{
		ID:        "entry_123",
		Type:      EntryUser,
		Content:   "Test content with special chars: <>&\"'",
		Timestamp: now,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		t.Fatal(err)
	}

	var decoded SessionEntry
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}

	if decoded.ID != "entry_123" {
		t.Errorf("ID: %s", decoded.ID)
	}
	if decoded.Type != EntryUser {
		t.Errorf("type: %s", decoded.Type)
	}
	if decoded.Content != "Test content with special chars: <>&\"'" {
		t.Errorf("content: %s", decoded.Content)
	}
}

func TestManagerMultipleSessions(t *testing.T) {
	dir := t.TempDir()
	mgr, _ := NewManager(dir)

	for i := 0; i < 10; i++ {
		_, err := mgr.CreateSession("/tmp", "openai", "gpt-4o")
		if err != nil {
			t.Fatalf("create %d failed: %v", i, err)
		}
	}

	sessions := mgr.ListSessions()
	if len(sessions) != 10 {
		t.Errorf("expected 10 sessions, got %d", len(sessions))
	}
}
