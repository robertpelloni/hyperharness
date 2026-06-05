package ingest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/robertpelloni/hyperharness/internal/memory"
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

	t.Run("NormalizeMarkdown", func(t *testing.T) {
		input := "Check [this link](http://example.com) for more info."
		expected := "Check this link for more info."
		result := processor.Normalize(input)
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("NormalizeHTML", func(t *testing.T) {
		input := "<div>  Hello   <b>World</b>  </div>"
		expected := "Hello World"
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

	t.Run("IngestFile", func(t *testing.T) {
		testFile := filepath.Join(tempDir, "test.go")
		os.WriteFile(testFile, []byte("package test\n\nfunc Main() {}"), 0644)

		err := processor.IngestFile(testFile, []string{"source"}, memory.ScopeProject)
		if err != nil {
			t.Fatalf("IngestFile failed: %v", err)
		}

		entries := kb.List(memory.ScopeProject)
		// We expect 2 entries now: IngestText from previous run + IngestFile
		if len(entries) != 2 {
			t.Fatalf("Expected 2 entries, got %d", len(entries))
		}

		// Find the newly added entry
		var fileEntry *memory.KnowledgeEntry
		for _, e := range entries {
			if e.Title == "test.go" {
				fileEntry = e
				break
			}
		}

		if fileEntry == nil {
			t.Fatal("file entry not found")
		}
		if !strings.Contains(fileEntry.Content, "func Main()") {
			t.Errorf("Expected summary content, got %q", fileEntry.Content)
		}
	})
}
