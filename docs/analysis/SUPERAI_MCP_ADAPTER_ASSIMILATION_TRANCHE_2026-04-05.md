# SuperAI MCP Adapter Assimilation Tranche

Date: 2026-04-05

## Summary

This tranche continues the `../superai` → `hyperharness` assimilation effort at another important runtime/integration boundary:
- `foundation/adapters/mcp.go`

This is a high-value seam because it sits where:
- MCP server config discovery
- tool-hint exposure
- route hints
- process-start behavior
- model-facing adapter context

come together.

## Key finding

`foundation/adapters/mcp.go` was already present in `hyperharness` and effectively parity-aligned with `../superai`.

The highest-value work here was therefore not redundant copying. It was to:
- preserve the same MCP-facing surface
- harden nil/config/discovery edge handling
- extract server-name/status/start helpers
- add regression coverage for the helper seams

## What changed

### 1. Nil-safe adapter behavior
`(*MCPAdapter).Status()` now tolerates a nil receiver and records a warning rather than depending on a fully initialized adapter.

`LookupServer()` also now resolves homeDir safely even when called through a nil adapter.

### 2. Helper extraction for status assembly
Added:
- `sortedMCPServerNames(...)`
- `buildMCPServerStatus(...)`
- `routeHintForAdapter(...)`
- `mcpAdapterWorkingDir(...)`

These helpers make server-name sorting, status construction, and adapter fallback behavior directly testable.

### 3. Helper extraction for process start behavior
Added:
- `startMCPServer(server, name, workingDir)`

This isolates server-start validation and working-directory application from the surrounding adapter method.

### 4. Behavior clarified, not changed away from parity
The intended behavior remains the same:
- MCP server names are sorted
- tool hints are generated deterministically
- route hints still use the Borg adapter when present
- configured server start still uses command/args/env/working-dir values

The difference is that the behavior is now explicit and regression-tested.

## Regression coverage added

Added:
- `foundation/adapters/mcp_assimilation_test.go`

Coverage includes:
- sorted MCP server-name generation
- MCP server status building
- nil-adapter status behavior
- nil-adapter route hint fallback
- nil-adapter lookup behavior
- missing-command validation for server start
- deterministic tool hint generation
- command resolvability checks
- working-directory propagation into server-start behavior

## Validation

Verified successfully:

```bash
go test ./foundation/adapters ./tools ./agent ./agents ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because MCP is one of the core integration lanes for the unified harness vision.

By making the config/status/start seams explicit and tested, `hyperharness` gains:
- safer MCP adapter behavior
- clearer discovery semantics
- better confidence in process-start contract behavior
- stronger foundations for later MCP aggregation/router work

## Design insight

A useful lesson from this tranche is:

> process-launch and config-discovery adapters are much easier to trust once their sorting, fallback, and validation logic are represented by small helper seams instead of being embedded directly inside one large method.

That is what this tranche accomplished.

## Recommended next move

The strongest next assimilation step is now:

1. continue through remaining runtime/integration seams where adapters, tool exposure, and CLI/runtime flow meet
2. prioritize boundary files that either:
   - assemble model-facing context
   - bridge external control planes
   - or wrap exact-name tool behavior for model use
3. keep validating with:
   - `go test ./foundation/adapters ./tools ./agent ./agents ./tui ./cmd ./foundation/...`

## Bottom line

This tranche confirms that `hyperharness` already contained the `superai` MCP adapter surface, and upgrades it into a clearer, safer, regression-tested base for the broader Go-native assimilation effort.
