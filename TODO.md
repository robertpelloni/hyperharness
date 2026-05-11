# HyperHarness TODO

# HyperHarness TODO

## Immediate (Next Session)

- [x] Project-scoped memory isolation
- [x] Memory export/import
- [x] Fuzz test: tool parameter parsing
- [x] Handle Windows path separators consistently across all tools
- [x] Add graceful context cancellation in agent loop
- [x] Fix `foundation/pi` TestLineScanner potential timeout
- [x] Fix `rpc` TestServerMultipleClients race condition

## Completed Deep Wiring & Integration
*Moved from Immediate since Phase 4 is now completed*



### Deep Tool Wiring
- [x] Wire `TodoWrite` tool to session-level state persistence
- [x] Wire `Agent` tool to actual `internal/subagents` Manager
- [x] Wire `LSP` tool to gopls/other language servers via stdio
- [x] Wire `WebSearch` to Exa or Brave Search API
- [x] Wire `WebFetch` to actual HTTP client
- [x] Wire `PowerShell` to actual PowerShell execution on Windows
- [x] Wire `Config` tool to actual `internal/config` Manager
- [x] Wire `Skill` tool to actual `internal/skills` Manager
- [x] Wire `platform__manage_schedule` to actual cron system

### MCP Deep Integration
- [x] Implement actual stdio transport for MCP server connections
- [x] Add SSE transport support
- [x] Tool discovery from connected MCP servers
- [x] Bidirectional routing (expose internal tools via MCP server)
- [x] Smithery registry API integration

### Memory System Enhancement
- [x] SQLite backend with FTS5 for knowledge base
- [x] Vector embeddings for semantic search
- [x] Memory decay weighting (older = less relevant)
- [x] Project-scoped memory isolation
- [x] Memory export/import

### Tests
- [x] Integration test: full agent loop with mock provider
- [x] Integration test: tool execution through Harness.ExecuteTool()
- [x] Integration test: MCP server connection lifecycle
- [x] Benchmark: tool dispatch latency
- [ ] Benchmark: memory search performance
- [x] Fuzz test: tool parameter parsing

## Short Term

### Bug Fixes & Robustness
- [x] Fix `foundation/pi` TestLineScanner potential timeout
- [x] Fix `rpc` TestServerMultipleClients race condition
- [x] Handle Windows path separators consistently across all tools
- [x] Add error recovery for crashed MCP server processes
- [x] Add graceful context cancellation in agent loop

### Code Quality
- [x] Add inline documentation to all tool Execute functions
- [ ] Add godoc comments to all exported types
- [x] Refactor duplicate helper functions across parity files
- [x] Consolidate truncateString, getStr, getInt into shared package
- [x] Add input validation to all tool parameters

### Feature Parity Gaps
- [x] OpenCode `plan_enter`/`plan_exit` — wire to actual planning mode
- [x] Crush `batch` — wire to actual parallel tool execution
- [x] Crush `delegate` — wire to actual subagent delegation
- [x] Crush `job_*` — wire to actual background process management
- [x] Smithery `smithery_install` — wire to actual registry API
- [x] Hypercode `context_manager` — wire to actual context compaction

## Medium Term

### Advanced Memory & UI
- [x] Build a Bubbletea-driven Terminal Dashboard with split panes
- [x] Implement Git-Aware Context (auto-inject git diffs and commits)
- [x] Implement AST-Aware Chunking (summarize files by function signatures instead of lines)

### Typescript Parity
- [x] Analyze `pi-cli/` plugins from npm
- [x] Sync parity schemas (PascalCase/snake_case) to `src/core/tools/parity`
- [x] Port FTS5 SQLite Memory with Age Decay to TS `core/memory`
- [x] Integrate Plugin hooks for Taskplane, Babysitter, Pompom

### Architecture
- [x] Implement multi-agent "Council" architecture (Director-Worker paradigm)
- [x] Implement "Debate Mode" for failing tests between agents
- [x] Auto-Discovery of local MCP Servers (scan VSCode/NPM global)
- [x] Implement Containerized Bash execution for sandboxing
- [x] Create permission manifests for subagents (restrict `write`/`bash`)
- [x] Auto-generate TypeScript and Go code (`Tool` objects) directly from upstream JSON schemas



### UI
- [ ] Full TUI with bubbletea
- [ ] Syntax-highlighted tool output
- [ ] Session tree visualization
- [ ] Model selector UI
- [ ] Token/cost tracking footer
- [ ] Real-time streaming display

### Performance
- [x] Memory-mapped JSONL reading for large sessions
- [x] Concurrent tool execution for independent calls
- [ ] Context window pre-compaction
- [x] Tool result caching
- [x] Connection pooling for LLM providers

## Ongoing
- [ ] Keep tool parity up to date with upstream harnesses
- [ ] Monitor new harness releases for new tool surfaces
- [ ] Maintain test coverage above 80%
- [ ] Keep documentation synchronized with code
