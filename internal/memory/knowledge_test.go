package memory

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestKnowledgeBaseStore(t *testing.T) {
	dir := t.TempDir()
	kb, err := NewKnowledgeBase(filepath.Join(dir, "test.json"))
	if err != nil {
		t.Fatal(err)
	}

	entry := &KnowledgeEntry{
		Title:   "Go Testing Patterns",
		Content: "Use table-driven tests for multiple test cases",
		Tags:    []string{"testing", "go", "patterns"},
		Scope:   ScopeGlobal,
	}

	if err := kb.Store(entry); err != nil {
		t.Fatalf("store failed: %v", err)
	}

	if entry.ID == "" {
		t.Error("entry should have an ID after storing")
	}
	if entry.CreatedAt.IsZero() {
		t.Error("entry should have a creation timestamp")
	}
}

func TestKnowledgeBaseSearch(t *testing.T) {
	dir := t.TempDir()
	kb, _ := NewKnowledgeBase(filepath.Join(dir, "test.json"))

	// Store some entries
	kb.Store(&KnowledgeEntry{
		Title: "Go Testing", Content: "Use table-driven tests in Go",
		Tags: []string{"testing", "go"}, Scope: ScopeGlobal,
	})
	kb.Store(&KnowledgeEntry{
		Title: "Python Testing", Content: "Use pytest for Python testing",
		Tags: []string{"testing", "python"}, Scope: ScopeGlobal,
	})
	kb.Store(&KnowledgeEntry{
		Title: "Go Concurrency", Content: "Use goroutines and channels",
		Tags: []string{"concurrency", "go"}, Scope: ScopeProject,
	})

	// Search by keyword
	results := kb.Search("testing", nil, "")
	if len(results) != 2 {
		t.Errorf("search 'testing': expected 2, got %d", len(results))
	}

	// Search by keyword and scope
	results = kb.Search("", nil, ScopeProject)
	if len(results) != 1 {
		t.Errorf("search scope project: expected 1, got %d", len(results))
	}

	// Search by tags
	results = kb.Search("", []string{"go"}, "")
	if len(results) != 2 {
		t.Errorf("search tag 'go': expected 2, got %d", len(results))
	}

	// Search with no matches
	results = kb.Search("nonexistent", nil, "")
	if len(results) != 0 {
		t.Errorf("search 'nonexistent': expected 0, got %d", len(results))
	}
}

func TestKnowledgeBaseGet(t *testing.T) {
	dir := t.TempDir()
	kb, _ := NewKnowledgeBase(filepath.Join(dir, "test.json"))

	entry := &KnowledgeEntry{
		Title: "Test Entry", Content: "content", Tags: []string{"test"}, Scope: ScopeGlobal,
	}
	kb.Store(entry)

	found, ok := kb.Get(entry.ID)
	if !ok {
		t.Fatal("entry should be found")
	}
	if found.Title != "Test Entry" {
		t.Errorf("title: %s", found.Title)
	}

	_, ok = kb.Get("nonexistent")
	if ok {
		t.Error("nonexistent should not be found")
	}
}

func TestKnowledgeBaseDelete(t *testing.T) {
	dir := t.TempDir()
	kb, _ := NewKnowledgeBase(filepath.Join(dir, "test.json"))

	entry := &KnowledgeEntry{
		Title: "Delete Me", Content: "content", Tags: []string{"test"}, Scope: ScopeGlobal,
	}
	kb.Store(entry)

	if err := kb.Delete(entry.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	_, ok := kb.Get(entry.ID)
	if ok {
		t.Error("entry should be gone after delete")
	}

	// Delete nonexistent
	if err := kb.Delete("nonexistent"); err == nil {
		t.Error("expected error for deleting nonexistent")
	}
}

func TestKnowledgeBaseList(t *testing.T) {
	dir := t.TempDir()
	kb, _ := NewKnowledgeBase(filepath.Join(dir, "test.json"))

	kb.Store(&KnowledgeEntry{Title: "Entry 1", Content: "c1", Scope: ScopeGlobal})
	kb.Store(&KnowledgeEntry{Title: "Entry 2", Content: "c2", Scope: ScopeProject})
	kb.Store(&KnowledgeEntry{Title: "Entry 3", Content: "c3", Scope: ScopeSession})

	all := kb.List("")
	if len(all) != 3 {
		t.Errorf("list all: expected 3, got %d", len(all))
	}

	global := kb.List(ScopeGlobal)
	if len(global) != 1 {
		t.Errorf("list global: expected 1, got %d", len(global))
	}
}

func TestKnowledgeBasePersistence(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "persist.json")

	// Create and store
	kb1, _ := NewKnowledgeBase(dbPath)
	kb1.Store(&KnowledgeEntry{
		Title: "Persistent Entry", Content: "this should survive restart",
		Tags: []string{"test"}, Scope: ScopeGlobal,
	})

	// Create new instance pointing to same file
	kb2, _ := NewKnowledgeBase(dbPath)
	kb2.Load()

	results := kb2.Search("Persistent", nil, "")
	if len(results) != 1 {
		t.Fatalf("persistent search: expected 1, got %d", len(results))
	}
	if results[0].Title != "Persistent Entry" {
		t.Errorf("title after reload: %s", results[0].Title)
	}
}

func TestKnowledgeBaseUpdate(t *testing.T) {
	dir := t.TempDir()
	kb, _ := NewKnowledgeBase(filepath.Join(dir, "test.json"))

	entry := &KnowledgeEntry{
		Title: "Original", Content: "original content", Tags: []string{"v1"}, Scope: ScopeGlobal,
	}
	kb.Store(entry)

	// Update by storing again with same ID
	entry.Title = "Updated"
	entry.Content = "updated content"
	entry.Tags = []string{"v2"}
	kb.Store(entry)

	found, _ := kb.Get(entry.ID)
	if found.Title != "Updated" {
		t.Errorf("title after update: %s", found.Title)
	}
	if found.Content != "updated content" {
		t.Errorf("content after update: %s", found.Content)
	}
}

func TestKnowledgeBaseContext(t *testing.T) {
	dir := t.TempDir()
	kb, _ := NewKnowledgeBase(filepath.Join(dir, "test.json"))

	kb.Store(&KnowledgeEntry{
		Title: "Project Config", Content: "This project uses Go 1.21",
		Tags: []string{"config"}, Scope: ScopeProject,
	})
	kb.Store(&KnowledgeEntry{
		Title: "Testing Rules", Content: "Always write table-driven tests",
		Tags: []string{"testing"}, Scope: ScopeGlobal,
	})

	ctx := kb.BuildContextForAgent("/tmp", []string{"Config"})
	if ctx == "" {
		t.Fatal("context should not be empty")
	}
	if !strings.Contains(ctx, "Project Config") {
		t.Errorf("context should contain project config title, got: %s", ctx)
	}
}

func TestKnowledgeBaseStats(t *testing.T) {
	dir := t.TempDir()
	kb, _ := NewKnowledgeBase(filepath.Join(dir, "test.json"))

	kb.Store(&KnowledgeEntry{Title: "E1", Content: "c1", Scope: ScopeGlobal, Tags: []string{"a"}})
	kb.Store(&KnowledgeEntry{Title: "E2", Content: "c2", Scope: ScopeGlobal, Tags: []string{"b"}})
	kb.Store(&KnowledgeEntry{Title: "E3", Content: "c3", Scope: ScopeProject, Tags: []string{"a", "b"}})

	stats := kb.Stats()
	if stats["totalEntries"] != 3 {
		t.Errorf("total entries: %v", stats["totalEntries"])
	}
}

func TestKnowledgeEntryJSONRoundtrip(t *testing.T) {
	entry := &KnowledgeEntry{
		ID:      "test_123",
		Title:   "Test",
		Content: "Content with special chars: <>&\"'",
		Tags:    []string{"tag1", "tag2"},
		Scope:   ScopeProject,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		t.Fatal(err)
	}

	var decoded KnowledgeEntry
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}

	if decoded.Title != "Test" {
		t.Errorf("title: %s", decoded.Title)
	}
	if len(decoded.Tags) != 2 {
		t.Errorf("tags count: %d", len(decoded.Tags))
	}
}

func TestKnowledgeBaseDefaultPath(t *testing.T) {
	kb, err := NewKnowledgeBase("")
	if err != nil {
		t.Fatalf("create with empty path: %v", err)
	}
	if kb.storePath == "" {
		t.Error("store path should default to home dir")
	}
}

func TestKnowledgeBaseEmptyFile(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "empty.json")
	os.WriteFile(dbPath, []byte("[]"), 0o644)

	kb, err := NewKnowledgeBase(dbPath)
	if err != nil {
		t.Fatalf("create from empty: %v", err)
	}

	entries := kb.List("")
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || (len(s) > 0 && contains(s[1:], sub)))
}
