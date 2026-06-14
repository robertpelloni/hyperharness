package features

// RegisterAllExtensions registers all built-in extension features from the
// hyperharness codebase. Call this during application startup (main.go / cmd/)
// to make all features discoverable via the global registry.
func RegisterAllExtensions() {
	registerFn := func(name string) {
		Register(name, func() error { return nil })
	}

	// ── TUI / UI Extensions ──────────────────────────────────────────
	registerFn("tui.markdown-renderer")
	registerFn("tui.split-layout")
	registerFn("tui.status-bar")
	registerFn("tui.provider-selector")
	registerFn("tui.model-selector")
	registerFn("tui.autocomplete")
	registerFn("tui.history-cycling")
	registerFn("tui.tree-browser")
	registerFn("tui.dashboard")
	registerFn("tui.logs-view")
	registerFn("tui.settings-panel")

	// ── Todo Feature ──────────────────────────────────────────────────
	registerFn("todos.add")
	registerFn("todos.list")
	registerFn("todos.toggle")
	registerFn("todos.clear")

	// ── Git Features ──────────────────────────────────────────────────
	registerFn("git.checkpoint")
	registerFn("git.auto-commit")
	registerFn("git.dirty-repo-guard")
	registerFn("git.protected-paths")
	registerFn("git.commit")

	// ── Agent / LLM Features ─────────────────────────────────────────
	registerFn("agent.council-debate")
	registerFn("agent.budget-enforcement")
	registerFn("agent.risk-evaluator")
	registerFn("agent.director-planning")
	registerFn("agent.session-management")
	registerFn("agent.code-mode")
	registerFn("agent.shell-assistant")
	registerFn("agent.context-compaction")
	registerFn("agent.thinking-levels")
	registerFn("agent.file-triggers")

	// ── Notifications ─────────────────────────────────────────────────
	registerFn("notify.osc777")
	registerFn("notify.osc99")
	registerFn("notify.windows-toast")

	// ── Plan Mode ─────────────────────────────────────────────────────
	registerFn("plan-mode.enabled")
	registerFn("plan-mode.steps")
	registerFn("plan-mode.restricted-tools")

	// ── Handoff / Summarize ───────────────────────────────────────────
	registerFn("handoff.prompt")
	registerFn("summarize.latest")
	registerFn("summarize.all")

	// ── Bookmark ──────────────────────────────────────────────────────
	registerFn("bookmark.session")
	registerFn("bookmark.restore")

	// ── Provider / Model ──────────────────────────────────────────────
	registerFn("provider.openai")
	registerFn("provider.anthropic")
	registerFn("provider.gemini")
	registerFn("provider.ollama")
	registerFn("provider.xiaomi")
	registerFn("provider.hypercode")

	// ── Foundation Tools ──────────────────────────────────────────────
	registerFn("pi.read")
	registerFn("pi.write")
	registerFn("pi.edit")
	registerFn("pi.bash")
	registerFn("pi.grep")
	registerFn("pi.find")
	registerFn("pi.ls")
	registerFn("pi.directory")

	// ── Integration / Parity ──────────────────────────────────────────
	registerFn("parity.claude-code")
	registerFn("parity.gemini-cli")
	registerFn("parity.opencode")
	registerFn("parity.goose")
	registerFn("parity.kimi")
	registerFn("parity.aider")
	registerFn("parity.hermes-agent")
	registerFn("parity.warp")
	registerFn("parity.wave")
	registerFn("parity.tabby")
	registerFn("parity.antigravity")
	registerFn("parity.codex-desktop")

	// ── Memory / Context ──────────────────────────────────────────────
	registerFn("memory.session-persistence")
	registerFn("memory.auto-save")
	registerFn("memory.rag-search")
	registerFn("memory.knowledge-base")

	// ── MCP (Model Context Protocol) ──────────────────────────────────
	registerFn("mcp.stdio-transport")
	registerFn("mcp.sse-transport")
	registerFn("mcp.tool-registry")
	registerFn("mcp.resource-access")

	// ── Council / Debate ──────────────────────────────────────────────
	registerFn("council.security-architect")
	registerFn("council.senior-engineer")
	registerFn("council.devops-chief")
	registerFn("council.consensus-voting")
	registerFn("council.parallel-execution")
}
