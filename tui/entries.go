package tui

// ═══════════════════════════════════════════════════════════════════════
// Chat entries — structured message history (mirrors pi-mono's session entries)
// ═══════════════════════════════════════════════════════════════════════

import (
	"time"
)

type EntryType int

const (
	EntryUser EntryType = iota
	EntryAssistant
	EntryTool
	EntrySystem
	EntryThinking
	EntryShellProposal
	EntryDiff
	EntryCompactionSummary
	EntryBashMode
)

// ChatEntry represents a single entry in the chat history.
// Mirrors pi-mono's session entry types (user, assistant, tool_call, tool_result, etc.)
type ChatEntry struct {
	Type      EntryType
	Content   string
	ToolName  string
	ToolArgs  string
	ToolOut   string
	ToolDur   time.Duration
	ToolErr   bool
	Expanded  bool
	Provider  string
	Model     string
	Timestamp time.Time

	// For thinking blocks
	ThinkingLevel string // "off", "minimal", "low", "medium", "high", "xhigh"
	Hidden        bool

	// For tool streaming
	Streaming bool

	// For compaction
	TokensBefore int
	TokensAfter  int
}

// ThinkingLevelColors maps thinking levels to border colors (pi-mono style)
var ThinkingLevelColors = map[string]string{
	"off":     "#6B7280",
	"minimal": "#3B82F6",
	"low":     "#8B5CF6",
	"medium":  "#F59E0B",
	"high":    "#EF4444",
	"xhigh":   "#DC2626",
}
