package ingest

import (
	"testing"
	"github.com/robertpelloni/hyperharness/internal/memory"
	"os"
	"path/filepath"
)

func TestDataProcessor(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "ingest-test")
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "memory.db")
	kb, _ := memory.NewKnowledgeBase(dbPath)
	processor := NewDataProcessor(kb)

	t.Run("Normalize", func(t *testing.T) {
		input := "  Hello   \n\n\nWorld  "
		expected := "Hello\n\nWorld"
		result := processor.Normalize(input)
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("IngestText", func(t *testing.T) {
		err := processor.IngestText("Test Title", "  Test   Content  ", []string{"test"}, memory.ScopeProject)
		if err != nil {
			t.Fatalf("IngestText failed: %v", err)
		}

		entries := kb.List(memory.ScopeProject)
		if len(entries) != 1 {
			t.Fatalf("Expected 1 entry, got %d", len(entries))
		}
		if entries[0].Title != "Test Title" {
			t.Errorf("Expected title 'Test Title', got %q", entries[0].Title)
		}
		if entries[0].Content != "Test Content" {
			t.Errorf("Expected content 'Test Content', got %q", entries[0].Content)
		}
	})
}
