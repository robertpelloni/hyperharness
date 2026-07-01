# Handoff Document

## Recent Progress
- Synchronized TS tools/parity tools (Claude, Crush, Goose, OpenCode, Kimi, Cursor, Windsurf) to maintain alignment with the Go backend, fixing build errors inside pi-cli/packages/coding-agent.
- Enhanced the TUI Dashboard (tui/dashboard.go, tui/chat.go) to represent core features (Memory Search, Agent Delegation, Code Execution, MCP Servers) visually with detailed tooltip [i] descriptions.
- Implemented and tuned a performance benchmark suite (bench/ folder) measuring Tool Dispatch Overhead, Context Compaction, Memory Search (SQLite FTS5), and Subagent Spawning, and documented baselines into BENCHMARKS.md.
- Implemented TS client parity for Memory Isolation and SQLite FTS5 in `pi-cli/packages/coding-agent/src/core/memory/sqlite-store.ts`.
- Synced remaining parity tools (Goose, OpenCode, Kimi, Cursor, Windsurf) by updating the Python schema extraction targets and regenerating TS tool facades.
- Began stubbing Phase 6 Daemon Hardening in `cmd/serve.go`.
- Created `parity-validation.test.ts` to test parity tool generation and unified output structure.

## Ongoing Work
- Missing end-to-end integration tests between TS Client and Go Core need to be created in `pi-cli/packages/coding-agent/test/` to fully clear the Phase 5 roadmap.
- Moving into Phase 6 readiness focusing on actual daemon hardening, crash recovery, auto-restart logic, and observability hooks using a formal daemon framework.
- The Go plugin system implementation is still missing.

## Technical Notes
- We use an automated Python script (scripts/generate_tools.py) to keep the Go parity tools and TypeScript parity bindings identical. The execute signature in the TS files must match exactly AgentTool[any]['execute'].
- Tests should always exclude submodules like aider and kilocode using 'go test $(go list ./... | grep -v aider\\|kilocode) -v' to avoid timeouts or unmanaged code builds.
