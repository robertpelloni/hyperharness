# Handoff

## What was done
- Added a Go-native foundation bootstrap for a Pi-derived harness:
  - `foundation/pi`
  - `foundation/compat`
  - `foundation/assimilation`
  - `cmd/foundation.go`
- Extended the bootstrap into a **truthful native runtime baseline**:
  - native `read`, `write`, `edit`, `bash` handlers
  - evented runtime execution
  - JSONL-backed session persistence with create/list/load/fork
  - CLI execution and session management commands under `foundation`
- Added the first native **Aider-inspired repo map baseline**:
  - ranked source scanning
  - lightweight symbol extraction
  - deterministic `<repo_map>` output
  - `hypercode foundation repomap`
- Strengthened repo-map quality with first **graph-ranking groundwork**:
  - files that reference symbols defined in other files now push ranking weight toward those defining files
  - mention-based ranking still works, but cross-file code relationships now matter too
- Routed top-level tool registration closer to the new foundation:
  - `tools.Registry` now exposes exact-name native Pi-compatible tools from `foundation/pi`
  - `agent.Agent` now registers per-tool schemas instead of one placeholder schema
  - legacy `tools/repomap.go` now delegates to `foundation/repomap`
- Added the first HyperCode/Borg adapter seam:
  - `foundation/adapters/hypercode.go` exposes assimilation status, memory context, provider status, MCP config visibility, and adjacent HyperCode repo discovery
  - `foundation/adapters/providers.go` exposes current provider/model visibility, detected providers, and Ollama model discovery when relevant
  - `foundation/adapters/mcp_config.go` and `foundation/adapters/mcp.go` expose MCP config parsing, server discovery, tool hints, route hints, and configured-server startup seams
  - `hypercode foundation adapters` now inspects that seam directly
  - `agent.Agent` now incorporates adapter-derived system context into its system prompt
- Migrated more top-level surfaces onto the foundation/adapter layer:
  - `mcp/client.go`, `mcp/manager.go`, and `mcp/config.go` now consume adapter-backed behavior
  - `cmd/foundation_http.go` provides reusable foundation-backed HTTP helpers
  - `cmd/serve.go` now exposes `/api/v1/foundation/*` endpoints and routes `/fs/read` through the foundation `read` tool
- Added deeper verification coverage:
  - snapshot-style tests for baseline tool results
  - top-level agent schema registration test
  - HyperCode/Borg adapter seam test
  - provider adapter seam test
  - MCP adapter seam tests and top-level MCP package tests
  - foundation-backed HTTP helper tests
- Added comprehensive `docs/ai/` documentation for requirements, design, planning, implementation, and testing.
- Added `docs/ai/design/upstream-toolchain-analysis.md` summarizing the imported upstream systems and assimilation strategy.
- Fixed the duplicate SQLite driver registration issue in `orchestrator/queue.go` and `orchestrator/vectors.go` by removing redundant `modernc.org/sqlite` imports.

## Validation completed
- `gofmt -w cmd/foundation.go cmd/foundation_http.go cmd/foundation_http_test.go cmd/serve.go foundation/adapters/hypercode.go foundation/adapters/hypercode_test.go foundation/adapters/providers.go foundation/adapters/providers_test.go foundation/adapters/mcp_config.go foundation/adapters/mcp.go foundation/adapters/mcp_test.go foundation/compat/types.go foundation/compat/catalog.go foundation/compat/default_catalog.go foundation/compat/catalog_test.go foundation/pi/foundation.go foundation/pi/foundation_test.go foundation/pi/runtime_types.go foundation/pi/runtime.go foundation/pi/runtime_test.go foundation/pi/session.go foundation/pi/session_test.go foundation/pi/tool_parity_test.go foundation/pi/tool_snapshot_test.go foundation/pi/tools_native.go foundation/assimilation/inventory.go foundation/assimilation/summary.go foundation/assimilation/inventory_test.go foundation/repomap/repomap.go foundation/repomap/repomap_test.go tools/registry.go tools/repomap.go tools/registry_test.go agent/agent.go agent/agent_test.go mcp/config.go mcp/client.go mcp/manager.go mcp/client_test.go mcp/manager_test.go mcp/mcphost.go mcp/mcphost_test.go orchestrator/vectors.go orchestrator/queue.go`
- `go test ./foundation/... ./cmd ./orchestrator ./tools ./agent ./mcp`
- `go run . foundation adapters`
- `go run . foundation tools`
- `go run . foundation inventory`
- `go run . foundation repomap --dir foundation --max-files 5`
- `go run . foundation session create --name smoke`
- `go run . foundation exec --tool write --input '{"path":"smoke.txt","content":"hello"}'`

## Important repo state notes
- The broader repository still has unrelated baseline failures under full `go test ./...`:
  - `aider/tests/fixtures/languages/go/test.go` has an unused import.
  - `mcp/mcphost_test.go` is out of sync with the host API.
- Many existing Go packages still contain placeholder or aspirational parity logic and should be migrated to the new `foundation/*` packages over time.

## Recommended next steps
1. Continue routing remaining top-level placeholder orchestration/tool surfaces onto the new `foundation/pi` runtime.
2. Expand verified result-shape and snapshot tests for `read`, `write`, `edit`, and `bash` plus CLI/HTTP smoke coverage.
3. Deepen `foundation/repomap` from graph-ranking groundwork toward fuller Aider-style graph ranking and richer edit engines.
4. Expand `foundation/adapters` from visibility seams into provider routing, memory, and richer MCP runtime adapters backed by HyperCode/Borg.
5. Migrate TUI and orchestration code to the new truthful foundation instead of placeholder parity claims.
