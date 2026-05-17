package tui

// ═══════════════════════════════════════════════════════════════════════
// Slash commands — mirrors pi-mono's BUILTIN_SLASH_COMMANDS
// ═══════════════════════════════════════════════════════════════════════

type SlashCommand struct {
	Name        string
	Description string
	Category    string // "session", "tools", "navigation", "settings", "model"
}

var BuiltinSlashCommands = []SlashCommand{
	// Session management
	{Name: "help", Description: "Show all slash commands", Category: "session"},
	{Name: "hotkeys", Description: "Show keyboard shortcuts", Category: "session"},
	{Name: "clear", Description: "Clear chat history", Category: "session"},
	{Name: "compact", Description: "Compact session context", Category: "session"},
	{Name: "session", Description: "Show session info and stats", Category: "session"},
	{Name: "name", Description: "Set session display name", Category: "session"},
	{Name: "fork", Description: "Create a fork from a previous message", Category: "session"},
	{Name: "export", Description: "Export session to file", Category: "session"},
	{Name: "import", Description: "Import a session from file", Category: "session"},
	{Name: "new", Description: "Start a new session", Category: "session"},
	{Name: "resume", Description: "Resume a different session", Category: "session"},
	{Name: "quit", Description: "Quit HyperHarness", Category: "session"},
	{Name: "exit", Description: "Quit HyperHarness", Category: "session"},

	// Navigation
	{Name: "tree", Description: "Navigate session tree", Category: "navigation"},
	{Name: "tree-select", Description: "Show numbered entry selector", Category: "navigation"},
	{Name: "tree-browser", Description: "Open file tree browser (modal)", Category: "navigation"},
	{Name: "tree-pane", Description: "Toggle persistent file tree pane", Category: "navigation"},
	{Name: "tree-pane-help", Description: "Show tree pane controls", Category: "navigation"},

	// Tools
	{Name: "commit", Description: "Generate git commit message", Category: "tools"},
	{Name: "plan", Description: "Build orchestration plan", Category: "tools"},
	{Name: "repomap", Description: "Generate repository map", Category: "tools"},
	{Name: "tools", Description: "List all available tools", Category: "tools"},
	{Name: "mcp", Description: "Show MCP tool listing", Category: "tools"},
	{Name: "adapters", Description: "Show adapter status", Category: "tools"},

	// Settings / Model
	{Name: "dashboard", Description: "Toggle split-pane dashboard", Category: "settings"},
	{Name: "model", Description: "Show or select model", Category: "model"},
	{Name: "providers", Description: "Show LLM provider status", Category: "model"},
	{Name: "settings", Description: "Open settings menu", Category: "settings"},
	{Name: "login", Description: "Login with OAuth provider", Category: "model"},
	{Name: "logout", Description: "Logout from OAuth provider", Category: "model"},
	{Name: "reload", Description: "Reload configuration", Category: "settings"},

	// Foundation
	{Name: "fsession", Description: "Show foundation session info", Category: "tools"},
	{Name: "summary-compact", Description: "Compact session context", Category: "session"},
	{Name: "summary-branch", Description: "Branch summary toward target", Category: "navigation"},
	{Name: "label", Description: "Set a label on an entry", Category: "navigation"},

	// Tree pane controls
	{Name: "tree-pane-show", Description: "Show file tree pane", Category: "navigation"},
	{Name: "tree-pane-hide", Description: "Hide file tree pane", Category: "navigation"},
	{Name: "tree-pane-focus", Description: "Toggle pane focus", Category: "navigation"},
	{Name: "tree-pane-size", Description: "Set pane height", Category: "navigation"},
	{Name: "tree-pane-preview", Description: "Toggle preview", Category: "navigation"},
	{Name: "tree-pane-grouped", Description: "Toggle grouped rendering", Category: "navigation"},
	{Name: "tree-pane-preset", Description: "Apply named preset", Category: "navigation"},
	{Name: "tree-pane-cycle", Description: "Cycle pane presets", Category: "navigation"},
	{Name: "tree-pane-status", Description: "Show pane state", Category: "navigation"},
	{Name: "tree-pane-reset", Description: "Reset to defaults", Category: "navigation"},
	{Name: "tree-pane-refresh", Description: "Refresh from runtime", Category: "navigation"},
	{Name: "tree-browser-clear", Description: "Clear browser state", Category: "navigation"},
	{Name: "changelog", Description: "Show changelog", Category: "session"},
}
