# Handoff Document

## Recent Progress
- Synchronized TS tools/parity tools (Claude, Crush, Goose, OpenCode, Kimi, Cursor, Windsurf) to maintain alignment with the Go backend, fixing build errors inside pi-cli/packages/coding-agent.
- Enhanced the TUI Dashboard (tui/dashboard.go, tui/chat.go) to represent core features (Memory Search, Agent Delegation, Code Execution, MCP Servers) visually with detailed tooltip [i] descriptions, now natively wired to `internal/sessions` and `internal/subagents` GlobalManagers to observe state in real-time.
- Implemented and tuned a performance benchmark suite (bench/ folder) measuring Tool Dispatch Overhead, Context Compaction, Memory Search (SQLite FTS5), and Subagent Spawning, and documented baselines into BENCHMARKS.md.
- Implemented TS client parity for Memory Isolation and SQLite FTS5 in `pi-cli/packages/coding-agent/src/core/memory/sqlite-store.ts`.
- Synced remaining parity tools (Goose, OpenCode, Kimi, Cursor, Windsurf) by updating the Python schema extraction targets and regenerating TS tool facades.
- Began stubbing Phase 6 Daemon Hardening in `cmd/serve.go` (Catching SIGTERM / SIGINT and panic wrapper deferrals).
- Updated `DEPLOY.md` to note new TUI telemetry and debug flag modes.
- Completed Phase 5 Production Readiness audit via `AUDIT.md` revealing false claim completions in Subagent execution, missing Test Coverages for parity tools, and missing integration test suites.

## Ongoing Work
- [RESOLVED] Wired actual LLM logic and StreamChat hooks to subagents.GlobalManager, replacing mock blocks identified in AUDIT.md.
- Need to fully implement dynamic plugin system loader rather than stubbing it.
- Configured `e2e_parity_test.go` to explicitly validate the full call chain of `Goose delegate`, `OpenCode task`, and `Claude Code Agent` parity tools through the GlobalManager spawning framework.

## Technical Notes
- We use an automated Python script (scripts/generate_tools.py) to keep the Go parity tools and TypeScript parity bindings identical. The execute signature in the TS files must match exactly AgentTool[any]['execute'].
- Tests should always exclude submodules like aider and kilocode using 'go test $(go list ./... | grep -v aider\\|kilocode) -v' to avoid timeouts or unmanaged code builds.

- Created `internal/subagents/manager_integration_test.go` to explicitly assert subagent tool permissions, task assignments, and context handover, fulfilling Phase 5 end-to-end integration constraints.
