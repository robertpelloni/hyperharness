//go:build !nosqlite
// +build !nosqlite

package memory

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSQLiteMemoryStoreCreation(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	store, err := NewSQLiteMemoryStore(dbPath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	if store.Path() != dbPath {
		t.Errorf("path mismatch: %s", store.Path())
	}
}

func TestSQLiteMemoryStoreAndRetrieve(t *testing.T) {
	dir := t.TempDir()
	store, err := NewSQLiteMemoryStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	entry := MemoryEntry{
		ID:      "test-entry-1",
		Content: "Go is a statically typed, compiled programming language",
		Type:    "knowledge",
		Tags:    []string{"go", "programming"},
		Source:  "test",
		Title:   "Go Language",
		Scope:   "global",
	}

	if err := store.Store(entry); err != nil {
		t.Fatal(err)
	}

	// Retrieve by ID
	retrieved, err := store.Get(entry.ID)
	if err != nil {
		t.Fatal(err)
	}
	if retrieved == nil {
		t.Fatal("entry should be found")
	}
	if retrieved.Content != entry.Content {
		t.Errorf("content mismatch: %s", retrieved.Content)
	}
	if retrieved.Title != "Go Language" {
		t.Errorf("title mismatch: %s", retrieved.Title)
	}
}

func TestSQLiteMemorySearch(t *testing.T) {
	dir := t.TempDir()
	store, err := NewSQLiteMemoryStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	// Store some entries
	entries := []MemoryEntry{
		{Content: "HyperHarness is a Go-based CLI harness", Type: "knowledge", Tags: []string{"go", "harness"}, Title: "HyperHarness Overview"},
		{Content: "MCP protocol enables tool discovery", Type: "knowledge", Tags: []string{"mcp", "protocol"}, Title: "MCP Protocol"},
		{Content: "SQLite FTS5 provides fast text search", Type: "technical", Tags: []string{"sqlite", "search"}, Title: "FTS5 Search"},
	}
	for _, e := range entries {
		if err := store.Store(e); err != nil {
			t.Fatal(err)
		}
	}

	// Search
	results, err := store.Search("harness", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Error("should find results for 'harness'")
	}

	// Search for SQLite content
	results2, err := store.Search("SQLite", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(results2) == 0 {
		t.Error("should find results for 'SQLite'")
	}
}

func TestSQLiteMemorySearchByScope(t *testing.T) {
	dir := t.TempDir()
	store, err := NewSQLiteMemoryStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	store.Store(MemoryEntry{Content: "global knowledge", Scope: "global", Title: "Global"})
	store.Store(MemoryEntry{Content: "project knowledge", Scope: "project:test", Title: "Project"})

	results, err := store.SearchByScope("knowledge", "global", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Error("should find global results")
	}
}

func TestSQLiteMemoryDelete(t *testing.T) {
	dir := t.TempDir()
	store, err := NewSQLiteMemoryStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	entry := MemoryEntry{ID: "delete-test", Content: "to be deleted", Title: "Delete Me"}
	store.Store(entry)

	// Verify exists
	retrieved, _ := store.Get(entry.ID)
	if retrieved == nil {
		t.Fatal("entry should exist")
	}

	// Delete
	if err := store.Delete(entry.ID); err != nil {
		t.Fatal(err)
	}

	// Verify gone
	retrieved, _ = store.Get(entry.ID)
	if retrieved != nil {
		t.Error("entry should be deleted")
	}
}

func TestSQLiteMemoryCount(t *testing.T) {
	dir := t.TempDir()
	store, err := NewSQLiteMemoryStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	count, err := store.Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Errorf("initial count should be 0, got %d", count)
	}

	store.Store(MemoryEntry{Content: "entry 1"})
	store.Store(MemoryEntry{Content: "entry 2"})

	count, _ = store.Count()
	if count != 2 {
		t.Errorf("count should be 2, got %d", count)
	}
}

func TestSQLiteMemoryCountByType(t *testing.T) {
	dir := t.TempDir()
	store, err := NewSQLiteMemoryStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	store.Store(MemoryEntry{Content: "k1", Type: "knowledge"})
	store.Store(MemoryEntry{Content: "k2", Type: "knowledge"})
	store.Store(MemoryEntry{Content: "d1", Type: "decision"})

	counts, err := store.CountByType()
	if err != nil {
		t.Fatal(err)
	}
	if counts["knowledge"] != 2 {
		t.Errorf("knowledge count: %d", counts["knowledge"])
	}
	if counts["decision"] != 1 {
		t.Errorf("decision count: %d", counts["decision"])
	}
}

func TestSQLiteMemoryListRecent(t *testing.T) {
	dir := t.TempDir()
	store, err := NewSQLiteMemoryStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	store.Store(MemoryEntry{Content: "first"})
	time.Sleep(10 * time.Millisecond)
	store.Store(MemoryEntry{Content: "second"})

	results, err := store.ListRecent(10)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	// Most recent first
	if results[0].Content != "second" {
		t.Errorf("most recent first: %s", results[0].Content)
	}
}

func TestSQLiteMemoryListByType(t *testing.T) {
	dir := t.TempDir()
	store, err := NewSQLiteMemoryStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	store.Store(MemoryEntry{Content: "k1", Type: "knowledge"})
	store.Store(MemoryEntry{Content: "d1", Type: "decision"})
	store.Store(MemoryEntry{Content: "k2", Type: "knowledge"})

	results, err := store.ListByType("knowledge", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 knowledge entries, got %d", len(results))
	}
}

func TestSQLiteMemoryTags(t *testing.T) {
	dir := t.TempDir()
	store, err := NewSQLiteMemoryStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	store.Store(MemoryEntry{Content: "entry1", Tags: []string{"go", "test"}})
	store.Store(MemoryEntry{Content: "entry2", Tags: []string{"go", "mcp"}})

	tags, err := store.Tags()
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 3 { // go, test, mcp
		t.Errorf("expected 3 unique tags, got %d: %v", len(tags), tags)
	}
}

func TestSQLiteMemorySearchByTags(t *testing.T) {
	dir := t.TempDir()
	store, err := NewSQLiteMemoryStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	store.Store(MemoryEntry{Content: "go entry", Tags: []string{"go", "programming"}})
	store.Store(MemoryEntry{Content: "python entry", Tags: []string{"python"}})
	store.Store(MemoryEntry{Content: "go mcp entry", Tags: []string{"go", "mcp"}})

	results, err := store.SearchByTags([]string{"go"}, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results for 'go' tag, got %d", len(results))
	}
}

func TestSQLiteMemoryUpdate(t *testing.T) {
	dir := t.TempDir()
	store, err := NewSQLiteMemoryStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	entry := MemoryEntry{ID: "update-test", Content: "original"}
	store.Store(entry)

	// Update with same ID
	entry.Content = "updated"
	store.Store(entry)

	retrieved, err := store.Get(entry.ID)
	if err != nil {
		t.Fatal(err)
	}
	if retrieved.Content != "updated" {
		t.Errorf("should be updated: %s", retrieved.Content)
	}

	count, _ := store.Count()
	if count != 1 {
		t.Errorf("should still have 1 entry, got %d", count)
	}
}

func TestSQLiteMemoryCloseAndReopen(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	store, err := NewSQLiteMemoryStore(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	store.Store(MemoryEntry{Content: "persistent"})
	store.Close()

	// Reopen
	store2, err := NewSQLiteMemoryStore(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer store2.Close()

	count, _ := store2.Count()
	if count != 1 {
		t.Errorf("should persist data, got %d entries", count)
	}

	results, _ := store2.Search("persistent", 10)
	if len(results) != 1 {
		t.Errorf("should find persisted entry, got %d", len(results))
	}
}

// Needed for fmt.Errorf in test
var _ = func() error { return fmt.Errorf("") }

// Ensure os and filepath are used
var _ = os.O_RDONLY
var _ = filepath.Join
