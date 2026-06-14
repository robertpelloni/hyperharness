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
- [ ] Implement Debate Mode for failing code edits
- [ ] Optimize subagent tool permissions based on role
- [ ] Add real-time streaming status updates to TUI
- [ ] Replace stub tool handlers with production implementations
- [ ] Add /debate slash command wiring council into TUI
- [ ] Add /save and /load session persistence commands
