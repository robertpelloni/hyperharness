# HANDOFF.md - Cross-Model Session Continuity

## Current State: v0.3.0
- **Deep Tool Wiring Complete**: TodoWrite, Agent, LSP, WebSearch, WebFetch, Config, Skill are now wired to real backends.
- **MCP Transports Complete**: Both stdio and SSE client transports are fully implemented.
- **509 tests across 33 packages, ALL PASSING**
- **136+ unique tool surfaces across 15+ harnesses**

## Architecture
```
hyperharness/
├── agent/          - Agent loop with tool calling
├── agents/         - Multi-agent provider registry
├── cmd/            - CLI commands (root, serve, tui)
├── config/         - Configuration management
├── foundation/     - Core: pi tools, adapters, repomap, orchestration
├── internal/
│   ├── borg/       - Borg core engine with adapters
│   ├── buildinfo/  - Version from VERSION file
│   ├── cache/      - TTL+LRU cache with events
│   ├── context/    - Context manager with compaction
│   ├── controlplane/ - Tool detector (30+ CLI harnesses)
│   ├── ctxharvester/ - Semantic chunking, token budgets
│   ├── eventbus/   - Pub/sub with wildcard patterns
│   ├── extensions/ - MCP extension manager
│   ├── fs/         - Language detection, gitignore walking
│   ├── git/        - Git CLI wrapper + submodule mgmt
│   ├── healer/     - LLM-powered error diagnosis
│   ├── mcp/        - MCP protocol client/server
│   ├── memory/     - Knowledge store + SQLite FTS5
│   ├── providers/  - LLM provider catalog + registry
│   ├── sessions/   - Session management
│   ├── skills/     - Skill discovery/execution
│   └── subagents/  - Subagent lifecycle
├── llm/            - Multi-provider LLM routing (8 providers)
├── mcp/            - MCP server, aggregator, stdio client
├── orchestrator/   - A2A broker, council, priority queue, workflow DAG
├── repl/           - Read-eval-print loop
├── rpc/            - JSON-RPC 2.0 server/client
├── security/       - Permission system
├── tools/          - 136+ tool implementations (all harnesses)
└── tui/            - Terminal UI (BubbleTea)
```

## Key Decisions
- All tool names match source harnesses EXACTLY (models trained on them)
- Package `internal/context` aliased as `contextmgr` to avoid stdlib conflict
- `CouncilDelegate` (not `CouncilMember`) for role constant to avoid struct name clash
- `PriorityQueue` (not `TaskQueue`) for in-memory queue; existing `TaskQueue` in queue.go uses SQLite
- `glebarez/sqlite` driver (not modernc) to avoid double-registration with orchestrator/database.go
- `SessionTodoStore` (not `TodoStore`) to avoid clash with crush_parity.go

## Build Instructions
```bash
export PATH="/c/Program Files/Go/bin:$PATH"
go build -buildvcs=false ./...
go test -buildvcs=false ./... -count=1 -timeout 180s
```

## Next Steps
1. Wire tool detectors into actual tool dispatch (use detected tools as backends)
2. Memory SQLite backend with vector embeddings
3. Integration tests with actual AI model tool calling
4. Port remaining submodule agent loop patterns (streaming, context budgets)
5. Performance benchmarks
6. MCP deep integration: Bidirectional routing (expose internal tools via MCP server)
