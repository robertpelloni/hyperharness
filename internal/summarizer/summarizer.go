package summarizer

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/robertpelloni/hyperharness/internal/ai"
)

// Config holds summarization configuration
type Config struct {
	Model          string  `json:"model"`
	MaxTokens      int     `json:"maxTokens"`
	TriggerPct     float64 `json:"triggerPct"`     // Trigger summarization at this % of maxTokens
	SummaryRatio   float64 `json:"summaryRatio"`   // Target ratio of summary to original tokens
	IncludeCode    bool    `json:"includeCode"`    // Whether to include code snippets in summary
	PreserveErrors bool    `json:"preserveErrors"` // Always preserve error messages
}

// DefaultConfig returns sensible defaults for summarization
func DefaultConfig() Config {
	return Config{
		Model:          "claude-sonnet-4-20250514",
		MaxTokens:      4000,
		TriggerPct:     80.0,
		SummaryRatio:   0.3,
		IncludeCode:    false,
		PreserveErrors: true,
	}
}

// Summary represents a compressed summary of conversation history
type Summary struct {
	mu            sync.RWMutex
	Content       string   `json:"content"`
	TokenCount    int      `json:"tokenCount"`
	MessageCount  int      `json:"messageCount"` // Number of messages summarized
	KeyPoints     []string `json:"keyPoints"`
	TasksCompleted []string `json:"tasksCompleted"`
	FilesModified  []string `json:"filesModified"`
}

// Summarizer provides LLM-based context summarization
type Summarizer struct {
	mu       sync.RWMutex
	provider ai.Provider
	config   Config
	history  []Summary
}

// NewSummarizer creates a new summarizer with the given provider and config
func NewSummarizer(provider ai.Provider, config Config) *Summarizer {
	if config.Model == "" {
		config = DefaultConfig()
	}
	return &Summarizer{
		provider: provider,
		config:   config,
		history:  make([]Summary, 0),
	}
}

// SummarizeContext compresses a list of messages into a summary using an LLM
func (s *Summarizer) SummarizeContext(ctx context.Context, messages []ai.Message) (*Summary, error) {
	if len(messages) == 0 {
		return &Summary{Content: "", TokenCount: 0, MessageCount: 0}, nil
	}

	// Build the text to summarize
	var sb strings.Builder
	sb.WriteString("Please summarize the following conversation history concisely. Focus on:\n")
	sb.WriteString("1. Key decisions and conclusions\n")
	sb.WriteString("2. Tasks completed and their outcomes\n")
	sb.WriteString("3. Files that were modified or created\n")
	sb.WriteString("4. Any unresolved issues or open questions\n\n")
	sb.WriteString("Keep the summary brief but preserve important technical details.\n\n")
	sb.WriteString("=== CONVERSATION HISTORY ===\n\n")

	for i, msg := range messages {
		sb.WriteString(fmt.Sprintf("[%d] %s: ", i+1, msg.Role))
		if len(msg.Content) > 500 {
			sb.WriteString(msg.Content[:500])
			sb.WriteString("...[truncated]\n\n")
		} else {
			sb.WriteString(msg.Content + "\n\n")
		}
	}

	sb.WriteString("=== END HISTORY ===\n\n")
	sb.WriteString("Provide a structured summary with sections for: Key Points, Tasks Completed, Files Modified, and Open Questions.")

	// Call LLM to generate summary
	resp, err := s.provider.GenerateText(ctx, s.config.Model, []ai.Message{
		{Role: "system", Content: "You are an expert at summarizing technical conversations. Create concise, actionable summaries that preserve key information."},
		{Role: "user", Content: sb.String()},
	})
	if err != nil {
		return nil, fmt.Errorf("summarization failed: %w", err)
	}

	// Parse the summary to extract structured information
	summary := s.parseSummary(resp.Content, len(messages))

	s.mu.Lock()
	defer s.mu.Unlock()
	s.history = append(s.history, *summary)

	return summary, nil
}

// parseSummary extracts structured information from the LLM's summary text
func (s *Summarizer) parseSummary(content string, messageCount int) *Summary {
	summary := &Summary{
		Content:      content,
		TokenCount:   len(content) / 4, // Rough estimate
		MessageCount: messageCount,
		KeyPoints:    make([]string, 0),
		TasksCompleted: make([]string, 0),
		FilesModified:  make([]string, 0),
	}

	// Simple parsing to extract sections (could be improved with better NLP)
	lines := strings.Split(content, "\n")
	currentSection := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Detect section headers
		lower := strings.ToLower(line)
		if strings.Contains(lower, "key point") || strings.Contains(lower, "key decisions") {
			currentSection = "keypoints"
			continue
		} else if strings.Contains(lower, "task") && (strings.Contains(lower, "completed") || strings.Contains(lower, "done")) {
			currentSection = "tasks"
			continue
		} else if strings.Contains(lower, "file") && (strings.Contains(lower, "modified") || strings.Contains(lower, "created")) {
			currentSection = "files"
			continue
		}

		// Extract items from lists
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "• ") || strings.HasPrefix(line, "* ") {
			item := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(line, "- "), "• "), "* "))
			switch currentSection {
			case "keypoints":
				summary.KeyPoints = append(summary.KeyPoints, item)
			case "tasks":
				summary.TasksCompleted = append(summary.TasksCompleted, item)
			case "files":
				summary.FilesModified = append(summary.FilesModified, item)
			}
		}
	}

	// If no structured sections found, just use the full content
	if len(summary.KeyPoints) == 0 && len(summary.TasksCompleted) == 0 && len(summary.FilesModified) == 0 {
		summary.KeyPoints = append(summary.KeyPoints, content)
	}

	return summary
}

// GetHistory returns the history of summaries
func (s *Summarizer) GetHistory() []Summary {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Summary, len(s.history))
	copy(result, s.history)
	return result
}

// ClearHistory clears the summary history
func (s *Summarizer) ClearHistory() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.history = make([]Summary, 0)
}

// ShouldTrigger returns true if summarization should be triggered based on current token usage
func (s *Summarizer) ShouldTrigger(currentTokens, maxTokens int) bool {
	if maxTokens <= 0 {
		return false
	}
	utilPct := float64(currentTokens) / float64(maxTokens) * 100
	return utilPct >= s.config.TriggerPct
}

// CompactMessages uses summarization to compact old messages while preserving recent ones
func (s *Summarizer) CompactMessages(ctx context.Context, messages []ai.Message, keepRecent int) ([]ai.Message, error) {
	if len(messages) <= keepRecent {
		return messages, nil
	}

	// Separate old messages to summarize from recent messages to keep
	oldMessages := messages[:len(messages)-keepRecent]
	recentMessages := messages[len(messages)-keepRecent:]

	// Summarize old messages
	summary, err := s.SummarizeContext(ctx, oldMessages)
	if err != nil {
		// Fallback to simple truncation if summarization fails
		return messages[len(messages)-keepRecent:], nil
	}

	// Create a system message with the summary
	summaryContent := fmt.Sprintf("=== PREVIOUS CONVERSATION SUMMARY ===\n%s\n=== END SUMMARY ===", summary.Content)
	summaryMsg := ai.Message{
		Role:    "system",
		Content: summaryContent,
	}

	// Combine summary with recent messages
	result := append([]ai.Message{summaryMsg}, recentMessages...)
	return result, nil
}