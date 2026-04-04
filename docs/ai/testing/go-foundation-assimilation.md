# Go Foundation Assimilation Testing Strategy

## Testing goals
The new foundation should be tested in layers.

## Layer 1: Contract tests
For every exact-name tool contract:
- validate parameter schema shape,
- validate name stability,
- validate result envelope invariants,
- compare native behavior against bridged/reference expectations.

Initial contract targets:
- `read`
- `write`
- `edit`
- `bash`

## Layer 2: Agent-loop tests
For the Pi-derived harness contract:
- event order (`agent_start` ... `agent_end`)
- tool execution sequencing
- parallel vs sequential tool modes
- steering/follow-up delivery modes
- session persistence and restore
- compaction hooks

## Layer 3: Integration tests
- HyperCode provider routing integration
- HyperCode MCP inventory integration
- memory retrieval integration
- imported session continuity
- TUI/JSON/RPC mode smoke tests

## Layer 4: Parity verification tests
For each assimilated upstream family:
- feature checklist coverage
- tool contract parity snapshots
- prompt/command UX parity snapshots where appropriate
- migration tests ensuring native replacement preserves behavior

## Current baseline run
Observed before any targeted fixes:

```text
go test ./...
- fixture compile error in aider/tests/fixtures/languages/go/test.go
- API mismatch in mcp/mcphost_test.go
- orchestrator panic from duplicate sqlite registration
```

## Tests added in this phase
- `foundation/compat/catalog_test.go`
- `foundation/assimilation/inventory_test.go`
- `foundation/pi/foundation_test.go`
- `foundation/pi/runtime_test.go`
- `foundation/pi/session_test.go`

## Tests that should be added next
1. `cmd/foundation` smoke tests
2. tool contract schema/result snapshot tests
3. truncation edge-case tests for `read` and `bash`
4. JSON/RPC transport tests
5. compatibility tests against bridged HyperCode-backed adapters

## Exit criteria for the next milestone
- foundation packages compile cleanly,
- contract registry tests pass,
- CLI inspection and execution commands are covered,
- default tool contracts are backed by real implementations,
- session persistence is stable under create/list/load/fork flows,
- maturity labels are truthful and enforced by tests.
