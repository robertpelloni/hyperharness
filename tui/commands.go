package tui

// ═══════════════════════════════════════════════════════════════════════
// Slash commands — unified from pi-mono, goose, opencode, claude-code
// ═══════════════════════════════════════════════════════════════════════

type SlashCommand struct {
	Name        string
	Description string
	Category    string // "session", "tools", "navigation", "settings", "model", "debug", "agent"
	Keybind     string // optional keybinding hint (pi-mono style)
}

var BuiltinSlashCommands = []SlashCommand{
	// ── Session management (pi-mono + claude-code) ──────────────
	{Name: "help", Description: "Show all slash commands", Category: "session", Keybind: "Ctrl+Shift+/"},
	{Name: "hotkeys", Description: "Show keyboard shortcuts", Category: "session"},
	{Name: "clear", Description: "Clear chat history", Category: "session"},
	{Name: "compact", Description: "Compact session context (pi-mono)", Category: "session"},
	{Name: "session", Description: "Show session info and stats", Category: "session"},
	{Name: "name", Description: "Set session display name", Category: "session"},
	{Name: "fork", Description: "Create a fork from a previous message", Category: "session", Keybind: "Ctrl+Shift+F"},
	{Name: "export", Description: "Export session to file", Category: "session"},
	{Name: "share", Description: "Share session (claude-code)", Category: "session"},
	{Name: "copy", Description: "Copy session to clipboard (pi-mono)", Category: "session"},
	{Name: "import", Description: "Import a session from file", Category: "session"},
	{Name: "new", Description: "Start a new session", Category: "session", Keybind: "Ctrl+Shift+N"},
	{Name: "resume", Description: "Resume a different session", Category: "session", Keybind: "Ctrl+Shift+R"},
	{Name: "quit", Description: "Quit HyperHarness", Category: "session"},
	{Name: "exit", Description: "Quit HyperHarness", Category: "session"},
	{Name: "rewind", Description: "Rewind to a previous turn (claude-code)", Category: "session"},

	// ── Navigation (pi-mono + opencode) ────────────────────────
	{Name: "tree", Description: "Navigate session tree (pi-mono)", Category: "navigation", Keybind: "Esc Esc"},
	{Name: "tree-select", Description: "Show numbered entry selector", Category: "navigation"},
	{Name: "tree-browser", Description: "Open file tree browser (modal)", Category: "navigation", Keybind: "Ctrl+L"},
	{Name: "tree-pane", Description: "Toggle persistent file tree pane", Category: "navigation"},
	{Name: "tree-pane-help", Description: "Show tree pane controls", Category: "navigation"},

	// ── Tools (pi-mono + goose + claude-code) ───────────────────
	{Name: "commit", Description: "Generate git commit message", Category: "tools", Keybind: "Ctrl+Y"},
	{Name: "plan", Description: "Build orchestration plan", Category: "tools"},
	{Name: "repomap", Description: "Generate repository map", Category: "tools"},
	{Name: "tools", Description: "List all available tools by group", Category: "tools"},
	{Name: "mcp", Description: "Show MCP tool listing", Category: "tools"},
	{Name: "adapters", Description: "Show adapter status", Category: "tools"},
	{Name: "diff", Description: "Show diff of uncommitted changes (claude-code)", Category: "tools"},
	{Name: "permissions", Description: "Manage tool permissions (claude-code)", Category: "tools"},
	{Name: "memory", Description: "Manage persistent memory (claude-code)", Category: "tools"},

	// ── Settings / Model (pi-mono + opencode) ──────────────────
	{Name: "dashboard", Description: "Toggle split-pane dashboard", Category: "settings", Keybind: "Ctrl+D"},
	{Name: "provider", Description: "Select LLM provider", Category: "model"},
	{Name: "model", Description: "Show or select model", Category: "model", Keybind: "Ctrl+Shift+M"},
	{Name: "scoped-models", Description: "Show scoped models (pi-mono)", Category: "model"},
	{Name: "providers", Description: "Show LLM provider status", Category: "model"},
	{Name: "settings", Description: "Open settings menu", Category: "settings"},
	{Name: "login", Description: "Login with OAuth provider", Category: "model"},
	{Name: "logout", Description: "Logout from OAuth provider", Category: "model"},
	{Name: "reload", Description: "Reload configuration", Category: "settings"},
	{Name: "theme", Description: "Change color theme (opencode)", Category: "settings"},
	{Name: "config", Description: "Edit configuration (claude-code)", Category: "settings"},
	{Name: "env", Description: "Show environment info (claude-code)", Category: "settings"},
	{Name: "doctor", Description: "Run diagnostics (claude-code)", Category: "settings"},
	{Name: "output-style", Description: "Set output style (claude-code)", Category: "settings"},
	{Name: "vim", Description: "Toggle vim keybindings (claude-code)", Category: "settings"},
	{Name: "color", Description: "Toggle color output (claude-code)", Category: "settings"},

	// ── Debug / Agent (pi-mono + goose) ────────────────────────
	{Name: "debug", Description: "Toggle debug mode (pi-mono)", Category: "debug"},
	{Name: "stats", Description: "Show usage statistics (claude-code)", Category: "debug"},
	{Name: "cost", Description: "Show cost breakdown (claude-code)", Category: "debug"},
	{Name: "thinkback", Description: "Review thinking history (claude-code)", Category: "debug"},
	{Name: "status", Description: "Show agent status (opencode)", Category: "agent"},
	{Name: "effort", Description: "Set effort level (claude-code)", Category: "agent"},
	{Name: "tag", Description: "Tag the current turn (claude-code)", Category: "agent"},
	{Name: "tasks", Description: "Show running tasks (claude-code)", Category: "agent"},

	// ── Foundation ─────────────────────────────────────────────
	{Name: "fsession", Description: "Show foundation session info", Category: "tools"},
	{Name: "summary-compact", Description: "Compact session context", Category: "session"},
	{Name: "summary-branch", Description: "Branch summary toward target", Category: "navigation"},
	{Name: "label", Description: "Set a label on an entry", Category: "navigation"},

	// ── Tree pane controls ─────────────────────────────────────
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
	{Name: "release-notes", Description: "Show release notes (claude-code)", Category: "session"},
	{Name: "todos", Description: "Show todo list (pi-mono)", Category: "tools"},
	{Name: "bookmark", Description: "Bookmark last message (pi-mono)", Category: "navigation"},
	{Name: "plan-mode", Description: "Toggle read-only plan mode (pi-mono)", Category: "agent"},
	{Name: "handoff", Description: "Transfer context to new session (pi-mono)", Category: "session"},
	{Name: "notify", Description: "Test desktop notification (pi-mono)", Category: "debug"},
}
