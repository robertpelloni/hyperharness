package tui

// ═══════════════════════════════════════════════════════════════════════
// Chat entries — structured message history
// Mirrors pi-mono's session entries + goose's Turn/ResponseItem + opencode's message types
// ═══════════════════════════════════════════════════════════════════════

import "time"

type EntryType int

const (
	EntryUser              EntryType = iota
	EntryAssistant
	EntryTool
	EntrySystem
	EntryThinking
	EntryShellProposal
	EntryDiff
	EntryCompactionSummary
	EntryBashMode
	EntryPermission        // goose-style permission request dialog
	EntryImage             // pi-mono clipboard image paste
	EntryCustom            // pi-mono custom/extension messages
	EntryError             // goose-style error display
	EntryQueue             // goose-style queued message indicator
)

// ChatEntry represents a single entry in the chat history.
// Unifies pi-mono's session entry types, goose's Turn/ResponseItem,
// opencode's message types, and claude-code's transcript entries.
type ChatEntry struct {
	Type          EntryType
	Content       string
	ToolName      string
	ToolArgs      string
	ToolOut       string
	ToolDur       time.Duration
	ToolErr       bool
	ToolKind      string // goose-style: "read", "edit", "delete", "move", "search", "execute", "think", "fetch", "other"
	ToolStatus    string // goose-style: "pending", "in_progress", "completed", "failed"
	ToolLocations []ToolLocation
	Expanded      bool
	Provider      string
	Model         string
	Timestamp     time.Time

	// Thinking
	ThinkingLevel string // "off", "minimal", "low", "medium", "high", "xhigh"
	Hidden        bool

	// Streaming
	Streaming bool

	// Compaction
	TokensBefore int
	TokensAfter  int

	// Permission (goose-style)
	PermissionOptions []PermissionOption
	PermissionIdx     int
	PermissionResolved bool

	// Image (pi-mono clipboard paste)
	ImagePath string
	ImageMime string

	// Custom/extension (pi-mono)
	CustomLabel string
	CustomBg    string

	// Error (goose-style)
	ErrorMessage string

	// Queue (goose-style)
	QueuePosition int
}

// ToolLocation represents a file:line reference (goose-style)
type ToolLocation struct {
	Path string
	Line int
}

// PermissionOption represents a goose-style permission choice
type PermissionOption struct {
	ID     string
	Name   string
	Kind   string // "allow_once", "allow_always", "reject_once", "reject_always"
	Active bool
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

// ToolKindIcons maps goose-style tool kinds to emoji icons
var ToolKindIcons = map[string]string{
	"read":     "📖",
	"edit":     "✏️",
	"delete":   "🗑",
	"move":     "📦",
	"search":   "🔍",
	"execute":  "▶",
	"think":    "💭",
	"fetch":    "🌐",
	"switch_mode": "🔀",
	"other":    "⚙",
}

// ToolStatusIndicators maps goose-style tool statuses to visual indicators
var ToolStatusIndicators = map[string]struct {
	Icon  string
	Color string
}{
	"pending":     {"○", "#5A6D84"},
	"in_progress": {"◑", "#C4883A"},
	"completed":   {"●", "#3A7D7B"},
	"failed":      {"✗", "#C0354A"},
}

// DetectToolKind infers a goose-style tool kind from the tool name
func DetectToolKind(name string) string {
	kinds := map[string]string{
		"read": "read", "write": "edit", "edit": "edit", "bash": "execute",
		"shell": "execute", "grep": "search", "find": "search", "search": "search",
		"glob": "search", "ls": "read", "tree": "read", "fetch": "fetch",
		"web": "fetch", "delete": "delete", "move": "move", "rename": "move",
		"think": "think", "thinking": "think", "plan": "think",
	}
	if kind, ok := kinds[name]; ok {
		return kind
	}
	return "other"
}
