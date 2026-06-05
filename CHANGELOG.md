## [0.4.4-alpha] - 2026-06-05

### Added
- Comprehensive tool signature parity for 6 major AI harnesses: Tabby AI, Warp Terminal, Wave Terminal, Hermes Agent, Antigravity 2.0, and Codex Desktop.
- Registered 18+ new tool surfaces in the central Go `Registry`.
- Detailed architectural analysis for each ported harness in `docs/analysis/`.
- Native Go implementation for Wave's `TermGetScrollback` tool, wired to internal session logs.
- Unit tests for Tabby, Warp, and Wave parity implementations.

### Fixed
- Repaired corrupted `pnpm-lock.yaml` indentation that blocked CI pipelines.
- Removed orphaned submodule gitlinks (`archive/claude-mem`) causing git clone failures.
- Resolved package-level variable name collision between `tools` and `sessions` packages.
- Fixed unused import regression in `hermes_parity.go`.

## [0.4.3] - 2026-04-27

### Added
- Tool JSON-RPC payload input validation wrapping `gojsonschema` dynamically prior to tool executions (`internal/toolregistry/validator.go`).
- Fuzz testing targeting raw JSON boundary conditions (`tools/fuzz_test.go`).
- Memory-Mapped session database IO. Swapped blocking `bufio.Scanner` for `syscall.Mmap` across session boundaries.
- Provider connection pooling explicitly capping TLS/TCP bounds globally.

### Changed
- Hardened all internal tools testing logic for tool parity mapping arrays securely.

## [0.4.1] - 2026-04-22

### Added
- Bubbletea Terminal Dashboard wired into TUI conditionally via `/dashboard`.
- Subagent Permission Manifests (`internal/subagents/manifests.go`) guarding dangerous execution tools.
- Auto-Discovery of Local MCP Servers from `claude_desktop_config.json`.
- HyperHarness internal registry exposed as an MCP Server standard (`internal/mcp/server.go`).
- AST-Aware chunking and Context window pre-compaction logic injected before OpenAI chat creations (`internal/ast`).
- Git-Aware context tracking injecting differential deltas naturally into the prompt (`internal/git/awareness.go`).
- Crush parity job background tools fully connected to `internal/sessions/background`.
- Goose/Claude `platform__manage_schedule` cron tasks heavily integrated via `internal/sessions/cron`.

### Changed
- Refactored `agent.go` pipeline to process tools concurrently inside `handleToolCalls` using `sync.WaitGroup`.
- Replaced redundant read-only tool executions with strict caching boundaries via `internal/cache`.

# HyperHarness Changelog

All notable changes to this project are documented here.
Version numbers follow [Semantic Versioning](https://semver.org/).

## [0.2.0] - 2026-04-08

### Added - Core Infrastructure
- **Unified Harness Integration Layer** (`hypercode/harness.go`) - Single entry point wiring 15+ subsystems
- **Borg Core Engine** (`internal/borg/core.go`) - Adapter pattern with lifecycle, events, hooks
- **RPC System** (`rpc/rpc.go`) - Full JSON-RPC 2.0 with TCP/Unix transport
- **Context Manager** (`internal/context/manager.go`) - Message tracking, compaction, injection, token estimation
- **Subagent Manager** (`internal/subagents/manager.go`) - 10 agent types (code, research, review, plan, build, test, debug, doc, security, devops)
- **Skill System** (`internal/skills/manager.go`) - SKILL.md discovery, parameter substitution, trigger matching, 4 built-in skills
- **Extension Manager** (`internal/extensions/manager.go`) - MCP extension lifecycle, Smithery registry, tool discovery, JSON config
- **File System Utilities** (`internal/fs/util.go`) - Language detection (30+ languages), gitignore walking, project root, size formatting

### Added - Tool Parity (23 new Claude Code tools)
- `TodoWrite` - Session-level task checklist with pending/in_progress/completed states
- `Agent` - Subagent spawning (Explore, Plan, verification, custom types)
- `WebSearch` - Web search with domain filtering
- `WebFetch` - URL content fetching
- `AskUserQuestion` - Interactive multi-choice questions
- `NotebookEdit` - Jupyter notebook cell editing
- `LSP` - Language Server Protocol integration (7 operations)
- `Skill` - Named skill invocation
- `ToolSearch` - Search available tools by capability
- `Config` - Configuration get/set/list/reset
- `EnterPlanMode` / `ExitPlanMode` - Planning mode toggle
- `EnterWorktree` / `ExitWorktree` - Git worktree management
- `SendMessage` - Inter-agent messaging
- `TaskCreate` / `TaskGet` / `TaskList` / `TaskUpdate` - Task management CRUD
- `Sleep` - Async delay utility
- `platform__manage_schedule` - Cron job management
- `PowerShell` - Windows-specific shell execution
- `Brief` - Project overview

### Fixed
- Restored borg/adapter.go (deleted in prior merge but still referenced)
- Removed stale borg import from agent/agent.go
- Fixed `NewGeminiHyperCodeProvider` → `NewGeminiBorgProvider` in orchestrator/queue_workers.go
- Fixed `cmd/foundation_http_test.go` variable redeclaration
- Fixed `agents/provider_assimilation_test.go` model name and content assertions
- Fixed `foundation/pi/tool_snapshot_test.go` diff normalization for DiffPrettyText ANSI output
- Fixed `foundation/pi/tools_extra_test.go` stringReader infinite loop (track read position)
- Fixed `foundation/pi/tools_extra_test.go` platform-aware path tests

### Test Coverage
- **315 tests** across **25 packages** - all passing
- New test files for: borg (7), extensions (11), memory (12), providers (12), sessions (11), mcp (13), skills (8), subagents (8), context (8), fs (11), rpc (6)

### Documentation
- `docs/analysis/INTEGRATION_STATUS_2026-04-08.md` - Comprehensive status report
- `docs/analysis/CLI_HARNESS_TOOL_PARITY_ANALYSIS_2026-04-08.md` - Tool parity analysis

## [0.1.0] - 2026-04-07

### Added - Initial Release
- **Foundation Pi Tools**: read, write, edit, bash, grep, find, ls (exact Pi parity)
- **122 Tool Surfaces** from 15+ CLI harnesses
- **Crush CLI Parity**: 18 tools (multiedit, view, write, glob, bash, webfetch, websearch, diagnostics, todos, tree, scan, download, load, safe, sourcegraph, codesearch, batch, delegate)
- **Gemini CLI Parity**: 7 tools (read_file, write_file, edit_file, shell, grep, list_directory, read_many_files)
- **OpenCode Parity**: 14 tools (file_read, file_write, file_edit, bash, grep, list_files, etc.)
- **Kimi CLI Parity**: 14 tools (file operations, search, task management)
- **Goose Parity**: 4 tools (developer, text_editor, browser, computer)
- **Cursor Parity**: 5 tools (file ops, code search, run command)
- **Windsurf Parity**: 2 tools (read_file, write_file)
- **Copilot CLI Parity**: 1 tool (run_in_terminal)
- **Aider V2 Parity**: 1 tool (aider)
- **Mistral Parity**: 2 tools (shell, editor)
- **Grok Parity**: 6 tools (file ops, search, list_directory)
- **Smithery Parity**: 2 tools (discover, install)
- **Hypercode Parity**: 3 tools (context_manager, memory_store, memory_search)
- **MCP Gateway**: 1 tool (mcp)
- **Agent System**: Director, Council, RAG, Autonomy, Disclosure
- **Provider System**: 15+ LLM providers with failover
- **Session System**: JSONL trees with branching and forking
- **Memory System**: Knowledge base with FTS, tags, scopes
- **MCP Protocol**: Client/server implementation
- **Security**: Permission manager with autonomy levels
- **TUI**: Interactive REPL with slash commands

## [0.3.0] - 2026-04-09

### Added
- Engineered a generic schema extractor and generator (`scripts/extract_schemas.py` and `scripts/generate_tools.py`) that successfully extracts over 1,000 tool schemas from upstream repositories and auto-generates parity code in both Go and TypeScript.
- Built `scripts/extract_schemas.py` and `scripts/generate_tools.py` to automatically extract tool schemas from `hermes-agent`, `claude-code`, `opencode`, and other submodules and auto-generate the Go (`tools/generated_parity.go`) and TypeScript (`core/tools/parity/generated.ts`) boilerplate definitions.
- Added exact parity tool schemas for `hermes-agent` and `ii-agent`.
- Wired up the `PowerShell` tool for direct execution on Windows environments.
- LLM multi-provider routing (8 providers: OpenAI, Anthropic, Gemini, DeepSeek, OpenRouter, Groq, LMStudio, Ollama)
- Provider catalog with status checking, health monitoring, and model selection
- Swarm and subagent system prompts for multi-model orchestration
- A2A message broker for inter-agent communication
- Council multi-model deliberation with voting and consensus
- PriorityQueue with LRU eviction, retry, cancellation, worker pool
- TTL+LRU in-memory cache with singleton registry and events
- EventBus with exact match, wildcard patterns, global listeners
- Healer service with LLM-powered error diagnosis and auto-fix
- Context harvester with semantic chunking, time decay, token budgets
- Git service: full CLI wrapper (log, status, diff, blame, stash, branch, remote, init, clone)
- Submodule management with concurrent update and structured reports
- Tool detector: detects 30+ installed CLI harnesses with versions
- Buildinfo: single source of truth for version (reads VERSION file)

### Metrics
- 509 tests across 33 packages
- 8 LLM providers with auto-routing
- 30+ CLI tool definitions for detection
- 136+ unique tool surfaces across 15+ harnesses
