# HyperHarness TODO

## Phase 5: Production Refinement
- [x] Implement actual LLM summarization for context compaction
- [x] Add budget enforcement and cost tracking to DirectorAgent
- [x] Fix budget timestamp handling (After → !Before)
- [x] Feature Registry scaffolding (internal/features)
- [x] Register all extensions as discoverable features
- [x] Markdown renderer (tui/markdown_renderer.go)
- [x] Tool Execution Dispatcher (agents/tool_dispatcher.go)
- [x] Adapter interfaces (internal/adapters/adapters.go)
- [x] Integrate ToolDispatcher + Budget into Director
- [x] Implement Debate Mode for failing code edits
- [x] Optimize subagent tool permissions based on role
- [x] Add real-time streaming status updates to TUI
- [x] Replace stub tool handlers with production implementations
- [x] Add /debate slash command wiring council into TUI
- [x] Add /save and /load session persistence commands
- [x] Resolve legacy merge conflict markers across web UI and package manifests
- [x] Sanitize build scripts (`build.sh`, `start.sh`) for dual-binary targets (`hypercode` and `hyperharness`)

## Phase 6: Ecosystem & Desktop Integration
- [ ] Desktop System Tray integration (`internal/tray` via systray) to manage daemon lifecycle, inspect status, and quit servers cleanly
- [ ] Complete consolidated single-page unified control plane view for operator dashboard
- [ ] MCP marketplace client connector and discovery integration
- [ ] VS Code and JetBrains IDE extensions packaging and distribution pipeline
