# Go Foundation Bootstrap Implementation Notes

## What was added

### Code
- `foundation/pi/foundation.go`
  - Pi-derived foundation specification
  - thinking levels
  - transport and delivery modes
  - run-event vocabulary
  - built-in tool descriptors

- `foundation/pi/runtime_types.go`
  - tool input/output types
  - run event payloads
  - truncation/detail structures

- `foundation/pi/tools_native.go`
  - native `read`, `write`, `edit`, and `bash` implementations
  - truncation behavior
  - exact-name tool handlers

- `foundation/pi/runtime.go`
  - evented runtime for native tool execution
  - session-linked tool execution
  - session creation/list/load/fork helpers

- `foundation/pi/session.go`
  - JSONL-backed session metadata and tree entries
  - create/load/save/list/fork operations

- `foundation/compat/types.go`
  - exact tool contract types
  - parity maturity model

- `foundation/compat/catalog.go`
  - thread-safe contract registry

- `foundation/compat/default_catalog.go`
  - Pi-compatible default tool contract set (`read`, `write`, `edit`, `bash`)
  - updated to `native` maturity for the default tool set

- `foundation/assimilation/inventory.go`
  - upstream assimilation inventory covering imported toolchains and HyperCode

- `foundation/assimilation/summary.go`
  - category summarization helpers

- `cmd/foundation.go`
  - CLI inspection for foundation spec, inventory, and tools
  - native execution surface for exact-name tools
  - native session create/list/show/fork commands
  - native repo-map generation command

- `foundation/repomap/repomap.go`
  - Aider-inspired native repo map baseline with graph-ranking groundwork
  - ranked file ordering using mentioned files/idents
  - lightweight definition/reference propagation across files
  - symbol extraction for common source forms

- `foundation/repomap/repomap_test.go`
  - repo map output validation
  - ranking validation for mention-based and graph-based prioritization

- `tools/registry.go`
  - top-level tool registry now exposes exact-name Pi-compatible tools via the native foundation runtime
  - repomap is available from the legacy registry surface as a native foundation-backed tool
  - tool schemas are now forwarded instead of using one placeholder schema for every tool

- `tools/repomap.go`
  - legacy wrapper now delegates to `foundation/repomap`

- `mcp/client.go`
  - top-level MCP client now uses the adapter seam for connection status and tool hint listing

- `mcp/manager.go`
  - top-level MCP manager now uses the adapter seam for configured server discovery and startup

- `mcp/config.go`
  - top-level MCP config now wraps adapter-owned parsing instead of duplicating config logic

- `mcp/mcphost.go`
  - defensive guard added for empty MCP binary path to avoid nil-process panics in tests

- `agent/agent.go`
  - top-level agent now advertises the native exact-name tools preferentially
  - OpenAI tool registration now uses per-tool schemas instead of one fake generic schema
  - system prompt now incorporates HyperCode/Borg and provider adapter context

- `foundation/adapters/hypercode.go`
  - first HyperCode/Borg adapter seam for the Go foundation
  - exposes assimilation status, memory context, provider status, MCP config visibility, and adjacent HyperCode repo discovery

- `foundation/adapters/providers.go`
  - provider adapter seam for current provider/model visibility
  - detects available providers from config/env
  - probes Ollama models when relevant

- `foundation/adapters/mcp_config.go`
  - adapter-owned MCP config parsing to avoid circular coupling
  - normalizes server command/env visibility for foundation consumers

- `foundation/adapters/mcp.go`
  - MCP adapter seam for configured server discovery, tool hints, route hints, and configured-server startup

- `foundation/adapters/hypercode_test.go`
  - validates adapter status, routing, and system-context construction

- `foundation/adapters/providers_test.go`
  - validates provider status/context construction

- `foundation/adapters/mcp_test.go`
  - validates MCP adapter status, tool hints, and route calls

- `tools/registry_test.go`
  - verifies exact Pi tools and repomap are present in the registry

- `foundation/pi/tool_snapshot_test.go`
  - snapshot-style verification for baseline tool results

- `agent/agent_test.go`
  - verifies top-level OpenAI tool registration exposes exact-schema tool definitions
  - verifies HyperCode adapter presence on the top-level agent

- `mcp/client_test.go`
  - verifies MCP client tool hint listing through the adapter seam

- `mcp/manager_test.go`
  - verifies MCP manager configured-tool listing and missing-server handling

### Documentation
- requirements, design, planning, implementation, and testing documents under `docs/ai/`

## Why this phase is useful
This phase still does not pretend to have completed the full port, but it now moves beyond pure scaffolding and establishes a truthful native baseline:
- one place for exact tool contracts,
- one place for the Pi-derived harness contract,
- one place for native default tool execution,
- one place for JSONL-backed native sessions,
- one place for the upstream assimilation inventory,
- one CLI surface to inspect and exercise those decisions,
- one documentation trail explaining the chosen architecture.

## Important baseline observations

### Existing codebase truthfulness gap
The current Go code advertises broad parity in some command descriptions, but several implementations are still placeholder-level. The new foundation work is intentionally separating:
- **declared compatibility** from
- **actual native implementation**.

### Existing test baseline issues
Before this phase, `go test ./...` already failed for unrelated reasons:
- `aider/tests/fixtures/languages/go/test.go` has an unused import.
- `mcp/mcphost_test.go` is out of sync with the host API.
- `orchestrator` panics because SQLite is registered twice.

These issues were observed and documented, not silently ignored or misrepresented as introduced by the new foundation work.

## Validation added in this phase
- native runtime tests for `read`, `write`, `edit`, and `bash`
- session persistence/list/fork tests
- ordered runtime event tests
- parity/truncation tests for `read` and `bash`
- snapshot-style tool result verification
- repo map generation and ranking tests
- top-level tool registry tests confirming native exact-name tool exposure
- top-level agent tool-schema registration tests
- HyperCode/Borg adapter seam tests
- provider adapter seam tests
- MCP adapter seam and top-level MCP package tests

## Recommended next implementation sequence
1. continue routing remaining top-level placeholder orchestration surfaces to `foundation/pi` runtime packages,
2. deepen repo-map ranking toward richer Aider-style graph semantics and add edit strategies,
3. expand `foundation/adapters` from status/config visibility into real HyperCode/Borg provider routing and richer MCP execution adapters,
4. expand snapshot/result-shape coverage and CLI smoke coverage,
5. layer in delegation, verification, detached/background runs, and JSON/RPC transport.
