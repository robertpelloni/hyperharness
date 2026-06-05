package ingest

import (
	"strings"
	"time"

	"github.com/robertpelloni/hyperharness/internal/memory"
)

// DataProcessor handles normalization and cleaning of input data.
type DataProcessor struct {
	KnowledgeBase *memory.KnowledgeBase
}

// NewDataProcessor creates a new data processor.
func NewDataProcessor(kb *memory.KnowledgeBase) *DataProcessor {
	return &DataProcessor{KnowledgeBase: kb}
}

// Normalize cleans and standardizes a text input.
func (p *DataProcessor) Normalize(input string) string {
	// Remove leading/trailing whitespace
	cleaned := strings.TrimSpace(input)

	// Normalize spaces within lines
	lines := strings.Split(cleaned, "\n")
	for i, line := range lines {
		lines[i] = strings.Join(strings.Fields(line), " ")
	}
	cleaned = strings.Join(lines, "\n")

	// Replace multiple newlines (3+) with double newlines
	for strings.Contains(cleaned, "\n\n\n") {
		cleaned = strings.ReplaceAll(cleaned, "\n\n\n", "\n\n")
	}
	return cleaned
}

// IngestText processes and stores a text snippet into the knowledge base.
func (p *DataProcessor) IngestText(title, content string, tags []string, scope memory.KnowledgeScope) error {
	normalized := p.Normalize(content)
	entry := &memory.KnowledgeEntry{
		Title:     title,
		Content:   normalized,
		Tags:      tags,
		Scope:     scope,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return p.KnowledgeBase.Store(entry)
}
