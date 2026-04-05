package pi

import (
	"context"
	"fmt"
	"strings"
)

// SummaryGenerator is the pluggable interface for generating branch summaries.
// It is intentionally small so later provider-backed generators can slot in
// without changing the session/runtime APIs.
type SummaryGenerator interface {
	GenerateBranchSummary(ctx context.Context, prep *BranchSummaryPreparation) (string, error)
}

// SummaryGeneratorFunc adapts a function to SummaryGenerator.
type SummaryGeneratorFunc func(ctx context.Context, prep *BranchSummaryPreparation) (string, error)

func (fn SummaryGeneratorFunc) GenerateBranchSummary(ctx context.Context, prep *BranchSummaryPreparation) (string, error) {
	return fn(ctx, prep)
}

// DeterministicSummaryGenerator is the default built-in summarizer.
// It does not call an LLM; instead it creates a structured, predictable
// summary from the prepared branch information.
type DeterministicSummaryGenerator struct{}

func NewDeterministicSummaryGenerator() SummaryGenerator {
	return DeterministicSummaryGenerator{}
}

func (DeterministicSummaryGenerator) GenerateBranchSummary(ctx context.Context, prep *BranchSummaryPreparation) (string, error) {
	if prep == nil {
		return "", fmt.Errorf("branch summary preparation is nil")
	}
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}
	goal := inferGoal(prep)
	done := inferCompletedItems(prep)
	nextSteps := inferNextSteps(prep)
	critical := inferCriticalContext(prep)

	var b strings.Builder
	b.WriteString("## Goal\n")
	b.WriteString(goal)
	b.WriteString("\n\n## Constraints & Preferences\n")
	if prep.MaxTokens > 0 {
		b.WriteString(fmt.Sprintf("- Summary preparation budget: %d estimated tokens\n", prep.MaxTokens))
	}
	if prep.CommonAncestorID != "" {
		b.WriteString(fmt.Sprintf("- Branch diverged from common ancestor `%s`\n", prep.CommonAncestorID))
	}
	b.WriteString("\n## Progress\n### Done\n")
	for _, item := range done {
		b.WriteString("- [x] ")
		b.WriteString(item)
		b.WriteString("\n")
	}
	b.WriteString("\n### In Progress\n")
	b.WriteString("- [ ] Continue work from the destination branch using the preserved summary context\n")
	b.WriteString("\n### Blocked\n")
	b.WriteString("- No explicit blockers recorded in the abandoned branch slice\n")
	b.WriteString("\n## Key Decisions\n")
	b.WriteString("- **Branch switch**: preserve abandoned branch context in a native branch summary entry\n")
	if len(prep.FileOps.ReadFiles) > 0 || len(prep.FileOps.ModifiedFiles) > 0 {
		b.WriteString("- **File context retained**: cumulative read/modified files are preserved even when the summary text window is budget-trimmed\n")
	}
	b.WriteString("\n## Next Steps\n")
	for i, item := range nextSteps {
		b.WriteString(fmt.Sprintf("%d. %s\n", i+1, item))
	}
	b.WriteString("\n## Critical Context\n")
	for _, item := range critical {
		b.WriteString("- ")
		b.WriteString(item)
		b.WriteString("\n")
	}
	b.WriteString("\n<read-files>\n")
	b.WriteString(strings.Join(prep.FileOps.ReadFiles, "\n"))
	b.WriteString("\n</read-files>\n\n<modified-files>\n")
	b.WriteString(strings.Join(prep.FileOps.ModifiedFiles, "\n"))
	b.WriteString("\n</modified-files>")
	return strings.TrimSpace(b.String()), nil
}

func inferGoal(prep *BranchSummaryPreparation) string {
	for _, entry := range prep.EntriesToSummarize {
		if entry.Kind == "message" && entry.Role == "user" && strings.TrimSpace(entry.Text) != "" {
			return truncateSummaryText(entry.Text, 240)
		}
	}
	if prep.TargetID != "" {
		return fmt.Sprintf("Continue work after switching branches toward target `%s`.", prep.TargetID)
	}
	return "Continue work after switching branches using the preserved abandoned-branch context."
}

func inferCompletedItems(prep *BranchSummaryPreparation) []string {
	items := []string{}
	if count := len(prep.EntriesToSummarize); count > 0 {
		items = append(items, fmt.Sprintf("Prepared %d abandoned branch entries for summarization", count))
	}
	if prep.EstimatedTokens > 0 {
		items = append(items, fmt.Sprintf("Serialized abandoned branch context (~%d estimated tokens)", prep.EstimatedTokens))
	}
	if len(prep.FileOps.ReadFiles) > 0 {
		items = append(items, fmt.Sprintf("Tracked %d read file(s)", len(prep.FileOps.ReadFiles)))
	}
	if len(prep.FileOps.ModifiedFiles) > 0 {
		items = append(items, fmt.Sprintf("Tracked %d modified file(s)", len(prep.FileOps.ModifiedFiles)))
	}
	if len(items) == 0 {
		items = append(items, "Prepared a branch summary transition")
	}
	return items
}

func inferNextSteps(prep *BranchSummaryPreparation) []string {
	steps := []string{}
	if prep.TargetID != "" {
		steps = append(steps, fmt.Sprintf("Resume execution from branch target `%s`", prep.TargetID))
	}
	if len(prep.FileOps.ModifiedFiles) > 0 {
		steps = append(steps, "Re-open or inspect the modified files if they remain relevant on the destination branch")
	}
	steps = append(steps, "Use the summary plus preserved file context to continue work without revisiting the full abandoned branch transcript")
	return steps
}

func inferCriticalContext(prep *BranchSummaryPreparation) []string {
	items := []string{}
	if prep.OldLeafID != "" {
		items = append(items, fmt.Sprintf("Abandoned branch leaf: `%s`", prep.OldLeafID))
	}
	if prep.CommonAncestorID != "" {
		items = append(items, fmt.Sprintf("Common ancestor with destination branch: `%s`", prep.CommonAncestorID))
	}
	if prep.MaxTokens > 0 {
		items = append(items, fmt.Sprintf("Prepared within a max summary budget of %d estimated tokens", prep.MaxTokens))
	}
	if prep.SerializedConversation != "" {
		items = append(items, fmt.Sprintf("Serialized abandoned branch payload length: %d chars", len(prep.SerializedConversation)))
	}
	return items
}

// GenerateBranchSummary uses the provided generator or falls back to the deterministic generator.
func (s *SessionStore) GenerateBranchSummary(ctx context.Context, prep *BranchSummaryPreparation, generator SummaryGenerator) (string, error) {
	if generator == nil {
		generator = NewDeterministicSummaryGenerator()
	}
	return generator.GenerateBranchSummary(ctx, prep)
}

// BranchWithGeneratedSummary prepares, generates, and appends a branch summary.
func (s *SessionStore) BranchWithGeneratedSummary(ctx context.Context, sessionID, targetID string, maxTokens int, generator SummaryGenerator, details any) (*SessionFile, string, error) {
	prep, err := s.PrepareBranchSummaryWithBudget(sessionID, targetID, maxTokens)
	if err != nil {
		return nil, "", err
	}
	summary, err := s.GenerateBranchSummary(ctx, prep, generator)
	if err != nil {
		return nil, "", err
	}
	session, err := s.BranchWithSummary(sessionID, targetID, summary, details)
	if err != nil {
		return nil, "", err
	}
	return session, summary, nil
}
