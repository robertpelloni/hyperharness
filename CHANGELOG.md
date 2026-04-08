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
