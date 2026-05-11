//go:build !nosqlite
// +build !nosqlite

package memory

import (
	"path/filepath"
	"testing"
)

func TestSQLiteMemoryStoreExportImport(t *testing.T) {
	dir := t.TempDir()
	dbPath1 := filepath.Join(dir, "test1.db")
	dbPath2 := filepath.Join(dir, "test2.db")

	store1, err := NewSQLiteMemoryStore(dbPath1)
	if err != nil {
		t.Fatalf("failed to create store1: %v", err)
	}
	defer store1.Close()

	entry1 := MemoryEntry{
		Content: "test export content 1",
		Scope:   "project_a",
		Tags:    []string{"export", "test"},
	}
	entry2 := MemoryEntry{
		Content: "test export content 2",
		Scope:   "project_b",
	}

	store1.Store(entry1)
	store1.Store(entry2)

	// Export all
	exports, err := store1.ExportAll()
	if err != nil {
		t.Fatalf("ExportAll failed: %v", err)
	}

	if len(exports) != 2 {
		t.Fatalf("Expected 2 exports, got %d", len(exports))
	}

	store2, err := NewSQLiteMemoryStore(dbPath2)
	if err != nil {
		t.Fatalf("failed to create store2: %v", err)
	}
	defer store2.Close()

	err = store2.ImportEntries(exports)
	if err != nil {
		t.Fatalf("ImportEntries failed: %v", err)
	}

	count, err := store2.Count()
	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}

	if count != 2 {
		t.Fatalf("Expected 2 entries in store2, got %d", count)
	}
}

func TestSQLiteMemoryStoreScopeIsolation(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test_scope.db")

	store, err := NewSQLiteMemoryStore(dbPath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	store.Store(MemoryEntry{Content: "global knowledge", Scope: "global"})
	store.Store(MemoryEntry{Content: "project A specific knowledge", Scope: "project_A"})
	store.Store(MemoryEntry{Content: "project A specific 2", Scope: "project_A"})
	store.Store(MemoryEntry{Content: "project B knowledge", Scope: "project_B"})

	// FTS5 scoped search
	results, err := store.SearchScoped("knowledge", "project_A", 10)
	if err != nil {
		t.Fatalf("SearchScoped failed: %v", err)
	}

	if len(results) != 1 { // Only "project A specific knowledge" matches "knowledge" in "project_A" scope
		t.Fatalf("Expected 1 result for scoped search, got %d", len(results))
	}

	results, err = store.SearchScoped("specific", "project_A", 10)
	if err != nil {
		t.Fatalf("SearchScoped failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results for scoped search, got %d", len(results))
	}

	// Wait, we need a small delay or test with just search if we test time decay
}
