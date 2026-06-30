# Handoff Document

## Recent Progress
- Synchronized TS tools/parity tools (Claude, Crush) to maintain alignment with the Go backend, fixing build errors inside pi-cli/packages/coding-agent.
- Enhanced the TUI Dashboard (tui/dashboard.go, tui/chat.go) to represent core features (Memory Search, Agent Delegation, Code Execution, MCP Servers) visually with detailed tooltip [i] descriptions.
- Implemented and tuned a performance benchmark suite (bench/ folder) measuring Tool Dispatch Overhead, Context Compaction, Memory Search (SQLite FTS5), and Subagent Spawning, and documented baselines into BENCHMARKS.md.

## Ongoing Work
- Phase 5 of ROADMAP.md is mostly complete. Next logical steps involve fleshing out 'End-to-end integration tests with real AI models' or 'Plugin system (Go plugins for custom tools)'. Or proceeding to Phase 6 readiness focusing on daemon hardening and crash recovery.

## Technical Notes
- We use an automated Python script (scripts/generate_tools.py) to keep the Go parity tools and TypeScript parity bindings identical. The execute signature in the TS files must match exactly AgentTool[any]['execute'].
- Tests should always exclude submodules like aider and kilocode using 'go test $(go list ./... | grep -v aider\\|kilocode) -v' to avoid timeouts or unmanaged code builds.
