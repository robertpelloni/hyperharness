// Package ui implements the terminal UI layer - the interactive REPL,
// TUI components, and message rendering system.
// Ported from Pi's interactive mode with charmbracelet/bubbletea.
// Features:
// - Rich TUI with syntax-highlighted tool outputs
// - Session tree navigation (/tree command)
// - Branch visualization and switching
// - Model selector with Ctrl+P cycling
// - Thinking level selector (Shift+Tab cycling)
// - File autocomplete (@ to fuzzy-search)
// - Image pasting and rendering
// - Extension UI components
// - Real-time token/cost tracking footer
// - Compaction summaries
// - Message queue indicators
package ui

import (
	"fmt"
	"strings"
	"time"
)

// Theme defines the color scheme and styling for the TUI.
type Theme struct {
	Name string
	
	// Colors
	Accent            string `json:"accent"`
	AccentMuted       string `json:"accentMuted"`
	Background        string `json:"background"`
	Foreground        string `json:"foreground"`
	UserMessage       string `json:"userMessage"`
	AssistantMessage  string `json:"assistantMessage"`
	ToolTitle         string `json:"toolTitle"`
	ToolOutput        string `json:"toolOutput"`
	ToolOutputMuted   string `json:"toolOutputMuted"`
	Error             string `json:"error"`
	Warning           string `json:"warning"`
	Info              string `json:"info"`
	Muted             string `json:"muted"`
	Success           string `json:"success"`
	Border            string `json:"border"`
	EditorBorder      string `json:"editorBorder"`
	FooterBackground  string `json:"footerBackground"`
	FooterForeground  string `json:"footerForeground"`
	ThinkingBlock     string `json:"thinkingBlock"`
	ModelSelector     string `json:"modelSelector"`
}

// DefaultDark returns the dark theme (matching Pi's dark theme).
func DefaultDark() Theme {
	return Theme{
		Name:            "dark",
		Accent:          "#7aa2f7",
		AccentMuted:     "#5577aa",
		Background:      "#1a1b26",
		Foreground:      "#c0caf5",
		UserMessage:     "#9ece6a",
		AssistantMessage: "#c0caf5",
		ToolTitle:       "#7aa2f7",
		ToolOutput:      "#a9b1d6",
		ToolOutputMuted: "#565f89",
		Error:           "#f7768e",
		Warning:         "#e0af68",
		Info:            "#7dcfff",
		Muted:           "#565f89",
		Success:         "#9ece6a",
		Border:          "#5577aa",
		EditorBorder:    "#7aa2f7",
		ThinkingBlock:   "#565f89",
	}
}

// DefaultLight returns the light theme.
func DefaultLight() Theme {
	return Theme{
		Name:            "light",
		Accent:          "#2563eb",
		AccentMuted:     "#64748b",
		Background:      "#ffffff",
		Foreground:      "#1e293b",
		UserMessage:     "#16a34a",
		AssistantMessage: "#1e293b",
		ToolTitle:       "#2563eb",
		ToolOutput:      "#475569",
		ToolOutputMuted: "#94a3b8",
		Error:           "#dc2626",
		Warning:         "#d97706",
		Info:            "#0284c7",
		Muted:           "#94a3b8",
		Success:         "#16a34a",
		Border:          "#cbd5e1",
		EditorBorder:    "#2563eb",
		ThinkingBlock:   "#94a3b8",
	}
}

// Message represents a rendered message in the UI.
type Message struct {
	Type      string    `json:"type"` // "user", "assistant", "tool_call", "tool_result", "notification", "error", "system"
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	Model     string    `json:"model,omitempty"`
	Provider  string    `json:"provider,omitempty"`
	
	// Tool specific
	ToolName   string `json:"toolName,omitempty"`
	ToolStatus string `json:"toolStatus,omitempty"` // "running", "completed", "failed"
	Duration   string `json:"duration,omitempty"`
	ErrorText  string `json:"errorText,omitempty"`
	
	// Thinking
	Thinking    string `json:"thinking,omitempty"`
	ThinkingVisible bool `json:"thinkingVisible"`
	
	// Context
	BranchID    string `json:"branchId,omitempty"`
	ParentID    string `json:"parentId,omitempty"`
	LabelText   string `json:"labelText,omitempty"`
	Labeled     bool   `json:"labeled"`
	
	// Token tracking
	TokenCount  int     `json:"tokenCount,omitempty"`
	CostUSD     float64 `json:"costUSD,omitempty"`
	
	// Render cache
	renderedLines []string
	renderedWidth int
}

// FooterState tracks the footer display data.
type FooterState struct {
	CWD          string
	SessionName  string
	ModelID      string
	Provider     string
	TotalTokens  int
	TotalCost    float64
	ThinkingLevel string
	
	// Message queue
	QueuedMessages int
	QueuedType     string // "steering", "follow-up"
	
	// Agent state
	AgentStatus string // "idle", "thinking", "tool_call", "error"
	Progress    string // Progress indicator for ongoing operations
}

// EditorState tracks the editor status.
type EditorState struct {
	Content    string
	CursorPos  int
	Focused    bool
	Width      int
	AutoComplete struct {
		Visible    bool
		Items      []AutoCompleteItem
		Selected   int
		Type       string // "file", "command", "model", "skill"
		SearchText string
	}
}

// AutoCompleteItem is an item in the autocomplete dropdown.
type AutoCompleteItem struct {
	Display string
	Value   string
	Icon    string
}

// CommandRegistry manages available slash commands.
type CommandRegistry struct {
	commands map[string]*Command
}

// Command represents a registered slash command.
type Command struct {
	Name        string
	Description string
	Handler     func(args string) (string, error)
}

// NewCommandRegistry creates a command registry.
func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		commands: make(map[string]*Command),
	}
}

// Register adds a command to the registry.
func (cr *CommandRegistry) Register(cmd *Command) {
	cr.commands[cmd.Name] = cmd
}

// Execute runs a command by name.
func (cr *CommandRegistry) Execute(name, args string) (string, error) {
	cmd, ok := cr.commands[name]
	if !ok {
		return fmt.Sprintf("Unknown command: /%s", name), fmt.Errorf("unknown command: %s", name)
	}
	return cmd.Handler(args)
}

// ListCommands returns all command names and descriptions.
func (cr *CommandRegistry) ListCommands() []Command {
	var cmds []Command
	for _, cmd := range cr.commands {
		cmds = append(cmds, *cmd)
	}
	return cmds
}

// FormatDuration formats milliseconds to a human-readable string.
func FormatDuration(ms int64) string {
	if ms < 1000 {
		return fmt.Sprintf("%dms", ms)
	}
	return fmt.Sprintf("%.1fs", float64(ms)/1000.0)
}

// FormatCost formats a dollar amount with appropriate precision.
func FormatCost(dollars float64) string {
	if dollars < 0.01 {
		return fmt.Sprintf("$%.4f", dollars)
	}
	return fmt.Sprintf("$%.2f", dollars)
}

// FormatTokenCount formats token count with K/M/B suffixes.
func FormatTokenCount(tokens int) string {
	if tokens < 1000 {
		return fmt.Sprintf("%d", tokens)
	}
	if tokens < 1000000 {
		return fmt.Sprintf("%.1fk", float64(tokens)/1000)
	}
	return fmt.Sprintf("%.1fm", float64(tokens)/1000000)
}

// ToolOutputFormatter formats tool execution output.
type ToolOutputFormatter struct {
	Theme Theme
}

// FormatToolCall formats a tool call for display.
func (f *ToolOutputFormatter) FormatToolCall(toolName string, args map[string]interface{}) string {
	argsStr := formatArgsShort(args)
	return fmt.Sprintf("%s %s", f.Theme.ToolTitle, toolName) + " " + argsStr
}

// FormatToolResult formats a tool result for display.
func (f *ToolOutputFormatter) FormatToolResult(toolName string, result string, duration time.Duration, isError bool) string {
	var parts []string
	
	// Title
	title := fmt.Sprintf("%s %s", f.Theme.ToolTitle, toolName)
	if isError {
		title = fmt.Sprintf("%s %s (error)", f.Theme.Error, toolName)
	}
	parts = append(parts, title)
	
	// Duration
	if duration > 0 {
		parts = append(parts, fmt.Sprintf(" (%s)", FormatDuration(duration.Milliseconds())))
	}
	
	// Output
	if result != "" {
		truncated := truncateResult(result, 2000)
		parts = append(parts, "\n"+f.Theme.ToolOutput+truncated)
	}
	
	return strings.Join(parts, "")
}

// FormatThinking formats thinking block.
func (f *ToolOutputFormatter) FormatThinking(content string) string {
	lines := strings.Split(content, "\n")
	formatted := make([]string, len(lines))
	for i, line := range lines {
		formatted[i] = f.Theme.ThinkingBlock + "│ " + line
	}
	return f.Theme.Muted + "━ Thinking ━━\n" + strings.Join(formatted, "\n")
}

// SystemPromptBuilder builds the complete system prompt.
type SystemPromptBuilder struct {
	basePrompt       string
	skills           []string
	agencies         []string // AGENTS.md contents
	promptTemplate   string
	memoryContext    string
	customPrompts    []string
}

// NewSystemPromptBuilder creates a new builder.
func NewSystemPromptBuilder() *SystemPromptBuilder {
	return &SystemPromptBuilder{
		basePrompt: DefaultBasePrompt(),
	}
}

// WithBasePrompt sets a custom base prompt.
func (b *SystemPromptBuilder) WithBasePrompt(prompt string) *SystemPromptBuilder {
	b.basePrompt = prompt
	return b
}

// WithSkills adds skill descriptions.
func (b *SystemPromptBuilder) WithSkills(skills []string) *SystemPromptBuilder {
	b.skills = skills
	return b
}

// WithAGENTS adds AGENTS.md content.
func (b *SystemPromptBuilder) WithAGENTS(contents []string) *SystemPromptBuilder {
	b.agencies = contents
	return b
}

// WithPromptTemplate adds a prompt template.
func (b *SystemPromptBuilder) WithPromptTemplate(template string) *SystemPromptBuilder {
	b.promptTemplate = template
	return b
}

// WithMemoryContext adds memory context.
func (b *SystemPromptBuilder) WithMemoryContext(context string) *SystemPromptBuilder {
	b.memoryContext = context
	return b
}

// WithCustomPrompt adds a custom prompt.
func (b *SystemPromptBuilder) WithCustomPrompts(prompts []string) *SystemPromptBuilder {
	b.customPrompts = prompts
	return b
}

// Build returns the complete system prompt.
func (b *SystemPromptBuilder) Build() string {
	var parts []string
	
	parts = append(parts, b.basePrompt)
	
	// Add AGENTS.md content
	for _, content := range b.agencies {
		if content != "" {
			parts = append(parts, "\n## Project Instructions (from AGENTS.md)\n\n"+content)
		}
	}
	
	// Add memory context
	if b.memoryContext != "" {
		parts = append(parts, "\n## Relevant Context\n\n"+b.memoryContext)
	}
	
	// Add skill descriptions
	if len(b.skills) > 0 {
		parts = append(parts, "\n## Available Skills\n\n"+strings.Join(b.skills, "\n\n"))
	}
	
	// Add custom prompts
	for _, prompt := range b.customPrompts {
		if prompt != "" {
			parts = append(parts, "\n"+prompt)
		}
	}
	
	return strings.Join(parts, "\n")
}

// DefaultBasePrompt returns the base system prompt.
func DefaultBasePrompt() string {
	return `You are an AI coding assistant working inside a terminal-based coding agent called hyperharness.

You have access to the following tools to interact with the filesystem and execute commands:
- read: Read file contents (supports offset/limit for large files, handles images)
- write: Create or overwrite files (with append mode support)
- edit: Edit files with exact text replacement (supports multiple edits per call)
- bash: Execute shell commands (with timeout support, streaming output, truncation)
- grep: Search file contents for patterns (with glob filtering, context lines)
- find: Find files by glob pattern (recursive, respects .gitignore)
- ls: List directory contents (with file type indicators)
- patch: Apply SEARCH/REPLACE block patches to files

Guidelines:
- Always check your work - read the file after editing to verify changes
- Use the most appropriate tool for each task
- For large files, use offset/limit to avoid reading unnecessary content
- When changing multiple files, batch your operations efficiently
- If you're unsure about file contents, read first before editing
- Handle errors gracefully and provide clear information about failures
- Write tests for code changes when appropriate
- Follow project conventions from AGENTS.md files when they exist
- Keep responses focused and actionable`
}

func truncateResult(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "\n...(truncated)"
}

func formatArgsShort(args map[string]interface{}) string {
	var parts []string
	for k, v := range args {
		s := fmt.Sprintf("%v", v)
		if len(s) > 100 {
			s = s[:100] + "..."
		}
		parts = append(parts, fmt.Sprintf("%s=%s", k, s))
	}
	return strings.Join(parts, " ")
}
