# HyperHarness Roadmap

## Phase 1: Foundation (COMPLETE ✅)
- [x] Port Pi's core tools to Go (read, write, edit, bash, grep, find, ls)
- [x] Tool registry with JSON Schema parameters
- [x] Session management with JSONL persistence
- [x] Memory system with knowledge base
- [x] Agent loop with streaming support

## Phase 2: Tool Parity (COMPLETE ✅)
- [x] 145+ tool surfaces from 15+ CLI harnesses
- [x] Exact tool names and parameter schemas matching source harnesses
- [x] Claude Code parity (28 tools): Agent, TodoWrite, WebSearch, LSP, Skills, etc.
- [x] Crush parity (18 tools): Background jobs, batch, diagnostics
- [x] Gemini CLI, OpenCode, Kimi, Goose, Cursor, Windsurf, Copilot, etc.
- [x] Delegation pattern: all tools delegate to foundation pi tools

## Phase 3: Infrastructure (COMPLETE ✅)
- [x] Borg Core Engine with adapter pattern
- [x] RPC system (JSON-RPC 2.0 with TCP/Unix)
- [x] Context manager (compaction, injection, token estimation)
- [x] Subagent manager (10 types)
- [x] Skill system (discovery, built-ins)
- [x] Extension manager (MCP lifecycle)
- [x] File system utilities
- [x] Provider system (15+ LLM providers)
- [x] Unified Harness integration layer

## Phase 4: Deep Integration (COMPLETE ✅)
- [x] Wire LSP tool to actual language server connections
- [x] Implement live WebSearch with Exa/Brave API
- [x] Implement live WebFetch with HTTP client
- [x] Wire TodoWrite to session persistence
- [x] Wire Agent tool to actual subagent spawning
- [x] Implement live MCP stdio transport (spawn/communicate with real servers)
- [x] Wire context compaction to actual LLM summarization
- [x] Connect knowledge base to SQLite backend with FTS5
- [x] Implement Smithery registry live discovery

## Phase 5: Production (COMPLETE ✅)
- [x] Integrate NPM `pi-package` plugins into TS PluginManager
- [x] Synchronize Memory Isolation and FTS5 to TypeScript client
- [x] Sync parity tools (Claude, Crush) to TS `tools/parity` folder

- [x] Performance benchmarks (sub-ms tool dispatch)
- [x] End-to-end integration tests with real AI models
- [ ] Cloud deployment (gRPC/HTTP API)
- [x] Plugin system (Go plugins for custom tools) (stubbed scaffolding)
- [x] Multi-agent orchestration (council debates, director mode)
- [x] Token streaming with SSE
- [ ] Cost tracking and budget enforcement
- [ ] Permission system with granular guardrails
- [ ] Multi-turn conversation optimization
- [ ] Automatic context window management
- [x] Data ingestion and preprocessing module (`internal/ingest`)
- [x] CLI interface for knowledge base ingestion (`hypercode ingest`)

## Phase 6: Ecosystem (PLANNED 📋)
- [x] Daemon hardening and crash recovery (stubbed in serve.go)
- [x] Auto-restart logic (stubbed in serve.go)
- [x] Comprehensive observability hooks (stubbed in serve.go)
- [ ] VS Code extension
- [ ] JetBrains plugin
- [ ] Neovim integration
- [ ] CI/CD GitHub Action
- [ ] MCP marketplace integration
- [ ] Cloud-hosted version

## Phase 5: Autonomous Core (COMPLETE ✅)
- [x] LLM-driven task decomposition in Council Director
- [x] Multi-agent execution loop with real tool-calling capabilities
- [x] Autonomous Autopilot mode using Director-Worker paradigm
- [x] Self-correction retry logic for subagent task failures
