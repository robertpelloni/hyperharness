# Go Foundation Bootstrap Implementation Notes

## What was added

### Code
- `foundation/pi/foundation.go`
  - Pi-derived foundation specification
  - thinking levels
  - transport and delivery modes
  - run-event vocabulary
  - built-in tool descriptors

- `foundation/compat/types.go`
  - exact tool contract types
  - parity maturity model

- `foundation/compat/catalog.go`
  - thread-safe contract registry

- `foundation/compat/default_catalog.go`
  - initial Pi-compatible default tool contract set (`read`, `write`, `edit`, `bash`)

- `foundation/assimilation/inventory.go`
  - upstream assimilation inventory covering imported toolchains and HyperCode

- `foundation/assimilation/summary.go`
  - category summarization helpers

- `cmd/foundation.go`
  - CLI inspection for foundation spec, inventory, and tools

### Documentation
- requirements, design, planning, implementation, and testing documents under `docs/ai/`

## Why this phase is useful
This phase does not pretend to have completed the full port. It creates the structures required to do that port in a disciplined way:
- one place for exact tool contracts,
- one place for the Pi-derived harness contract,
- one place for the upstream assimilation inventory,
- one place to inspect those decisions from the CLI,
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

## Recommended next implementation sequence
1. replace placeholder agent loop and session handling with `foundation/pi`-backed runtime packages,
2. add native implementations and contract tests for `read`, `write`, `edit`, and `bash`,
3. port repo-map and edit strategies,
4. add HyperCode/Borg provider and MCP adapters,
5. layer in delegation, verification, and detached/background runs.
