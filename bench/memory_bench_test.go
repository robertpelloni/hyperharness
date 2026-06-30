package bench

import (
	"fmt"
	"os"
	"testing"

	"github.com/robertpelloni/hyperharness/internal/memory"
)

func BenchmarkMemorySearch(b *testing.B) {
	b.StopTimer()
	dbPath := "bench_memory.db"
	os.Remove(dbPath)
	defer os.Remove(dbPath)

	store, err := memory.NewSQLiteMemoryStore(dbPath)
	if err != nil {
		b.Fatalf("Failed to create memory store: %v", err)
	}
	defer store.Close()

	// Populate data
	for i := 0; i < 1000; i++ {
		err := store.Store(memory.MemoryEntry{
			Content: fmt.Sprintf("Benchmark data content number %d for testing memory search", i),
			Scope: "bench_scope",
		})
		if err != nil {
			b.Fatalf("Failed to save memory: %v", err)
		}
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = store.Search("content number", 10)
	}
}

func BenchmarkMemorySearchScoped(b *testing.B) {
	b.StopTimer()
	dbPath := "bench_memory_scoped.db"
	os.Remove(dbPath)
	defer os.Remove(dbPath)

	store, err := memory.NewSQLiteMemoryStore(dbPath)
	if err != nil {
		b.Fatalf("Failed to create memory store: %v", err)
	}
	defer store.Close()

	// Populate data
	for i := 0; i < 1000; i++ {
		scope := "bench_scope"
		if i%2 == 0 {
			scope = "other_scope"
		}
		err := store.Store(memory.MemoryEntry{
			Content: fmt.Sprintf("Benchmark data content number %d for testing memory search", i),
			Scope: scope,
		})
		if err != nil {
			b.Fatalf("Failed to save memory: %v", err)
		}
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = store.SearchScoped("content number", "bench_scope", 10)
	}
}
