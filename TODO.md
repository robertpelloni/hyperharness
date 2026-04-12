# HyperHarness TODO

## Immediate (Next Session)

### Deep Tool Wiring
- [x] Wire `TodoWrite` tool to session-level state persistence
- [x] Wire `Agent` tool to actual `internal/subagents` Manager
- [x] Wire `LSP` tool to gopls/other language servers via stdio
- [x] Wire `WebSearch` to Exa or Brave Search API
- [x] Wire `WebFetch` to actual HTTP client
- [ ] Wire `PowerShell` to actual PowerShell execution on Windows
- [x] Wire `Config` tool to actual `internal/config` Manager
- [x] Wire `Skill` tool to actual `internal/skills` Manager
- [ ] Wire `platform__manage_schedule` to actual cron system

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
- [ ] Project-scoped memory isolation
- [ ] Memory export/import

### Tests
- [x] Integration test: full agent loop with mock provider
- [x] Integration test: tool execution through Harness.ExecuteTool()
- [x] Integration test: MCP server connection lifecycle
- [ ] Benchmark: tool dispatch latency
- [ ] Benchmark: memory search performance
- [ ] Fuzz test: tool parameter parsing

## Short Term

### Bug Fixes & Robustness
- [ ] Fix `foundation/pi` TestLineScanner potential timeout
- [ ] Fix `rpc` TestServerMultipleClients race condition
- [ ] Handle Windows path separators consistently across all tools
- [ ] Add error recovery for crashed MCP server processes
- [ ] Add graceful context cancellation in agent loop

### Code Quality
- [ ] Add inline documentation to all tool Execute functions
- [ ] Add godoc comments to all exported types
- [ ] Refactor duplicate helper functions across parity files
- [x] Consolidate truncateString, getStr, getInt into shared package
- [ ] Add input validation to all tool parameters

### Feature Parity Gaps
- [x] OpenCode `plan_enter`/`plan_exit` — wire to actual planning mode
- [ ] Crush `batch` — wire to actual parallel tool execution
- [x] Crush `delegate` — wire to actual subagent delegation
- [ ] Crush `job_*` — wire to actual background process management
- [ ] Smithery `smithery_install` — wire to actual registry API
- [ ] Hypercode `context_manager` — wire to actual context compaction

## Medium Term

### Architecture
- [ ] Plugin system for custom tool providers
- [ ] Event bus for inter-subsystem communication
- [ ] Middleware chain for tool execution (pre/post hooks)
- [ ] Rate limiting per provider
- [ ] Circuit breaker for failing providers
- [ ] Retry with exponential backoff

### UI
- [ ] Full TUI with bubbletea
- [ ] Syntax-highlighted tool output
- [ ] Session tree visualization
- [ ] Model selector UI
- [ ] Token/cost tracking footer
- [ ] Real-time streaming display

### Performance
- [ ] Memory-mapped JSONL reading for large sessions
- [ ] Concurrent tool execution for independent calls
- [ ] Context window pre-compaction
- [ ] Tool result caching
- [ ] Connection pooling for LLM providers

## Ongoing
- [ ] Keep tool parity up to date with upstream harnesses
- [ ] Monitor new harness releases for new tool surfaces
- [ ] Maintain test coverage above 80%
- [ ] Keep documentation synchronized with code
