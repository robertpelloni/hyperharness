package ingest

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/robertpelloni/hyperharness/internal/ast"
	"github.com/robertpelloni/hyperharness/internal/memory"
)

var (
	htmlTagRe     = regexp.MustCompile("<[^>]*>")
	markdownLinkRe = regexp.MustCompile(`\[([^\]]+)\]\([^\)]+\)`)
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
	// Strip HTML tags if present (heuristic)
	if strings.Contains(input, "<") && strings.Contains(input, ">") {
		input = htmlTagRe.ReplaceAllString(input, "")
	}

	// Strip Markdown links: [text](url) -> text
	if strings.Contains(input, "[") && strings.Contains(input, "]") {
		input = markdownLinkRe.ReplaceAllString(input, "$1")
	}

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

// ProcessFile reads and normalizes a file, optionally summarizing it if it's Go code.
func (p *DataProcessor) ProcessFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", path, err)
	}

	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".go" {
		summary, err := ast.SummarizeGoFile(path, data)
		if err == nil {
			return p.Normalize(summary), nil
		}
		// Fallback to raw content if summarization fails
	}

	return p.Normalize(string(data)), nil
}

// IngestFile processes a file and stores it in the knowledge base.
func (p *DataProcessor) IngestFile(path string, tags []string, scope memory.KnowledgeScope) error {
	content, err := p.ProcessFile(path)
	if err != nil {
		return err
	}

	title := filepath.Base(path)
	return p.IngestText(title, content, append(tags, "file", filepath.Ext(path)), scope)
}
