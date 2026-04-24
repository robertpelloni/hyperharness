# SuperAI Tools Registry Assimilation Tranche

Date: 2026-04-05

## Summary

This tranche widens the `../superai` → `hyperharness` assimilation effort from the hardened `agent/` seam into the central tool-registry surface:
- `tools/registry.go`

This is a high-value seam because it is the practical boundary between:
- exact-name model-facing tool exposure
- Pi-compatible foundation tool wrapping
- the core agent tool-call loop
- legacy compatibility tools

## Key finding

`tools/registry.go` was already present in `hyperharness` and effectively parity-aligned with `../superai`.

The highest-value work here was therefore not redundant copying. It was to:
- preserve the same exposed tool surface
- make exact-name lookup explicit in the registry itself
- extract foundation-tool wrapping into a testable helper
- add regression coverage for result formatting and lookup behavior

## What changed

### 1. Explicit exact-name registry lookup
Added:
- `(*Registry).Find(name string) (Tool, bool)`

Behavioral characteristics:
- nil-safe
- blank-name safe
- exact-name lookup only

This is a meaningful improvement because the registry now exposes its own exact-name lookup contract rather than requiring callers to manually loop over `Registry.Tools`.

### 2. Foundation-tool wrapper extraction
Added:
- `newFoundationTool(runtime, contract)`

This isolates the Pi/foundation tool wrapping behavior into a single helper.

Benefits:
- contract cloning becomes explicit
- execution/result formatting behavior is directly testable
- registration code is simpler and clearer

### 3. Core agent loop now uses registry lookup
Updated:
- `agent/agent.go`

The central tool-call execution path now uses:
- `registry.Find(...)`

instead of open-coded manual iteration over the tool list.

This makes the exact-name lookup path more explicit and aligned with the registry contract.

## Regression coverage added

Added:
- `tools/registry_assimilation_test.go`

Coverage includes:
- exact-name registry lookup behavior
- nil/blank lookup handling
- foundation-tool result formatting behavior
- cloned-contract behavior in foundation tool wrapping
- runtime-backed foundation tool execution through the extracted helper
- read-tool discovery through the new `Find(...)` path

## Validation

Verified successfully:

```bash
go test ./tools ./agent ./agents ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because tool lookup and tool wrapping are among the most important exact-contract surfaces in the entire harness.

By making them more explicit and testable, `hyperharness` gains:
- clearer exact-name tool resolution
- safer foundation-tool wrapping behavior
- stronger confidence that the registry remains compatible with model expectations

## Design insight

A useful lesson from this tranche is:

> exact-name tool parity is not just about the tool schemas; it also depends on the lookup and wrapping seams that expose those tools to the model loop.

This tranche strengthens those seams directly.

## Recommended next move

The strongest next assimilation step is now:

1. continue reviewing the remaining runtime/integration seams around tool exposure and execution
2. prefer boundary files where exact-name tool contracts meet providers, adapters, or CLI flows
3. keep validating with:
   - `go test ./tools ./agent ./agents ./tui ./cmd ./foundation/...`

## Bottom line

This tranche confirms that `hyperharness` already contained the `superai` tool-registry surface, and upgrades it into a clearer, safer, regression-tested base for exact-name tool-contract assimilation.
