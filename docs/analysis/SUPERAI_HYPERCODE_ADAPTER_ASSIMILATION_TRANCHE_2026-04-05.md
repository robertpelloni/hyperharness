# SuperAI HyperCode Adapter Assimilation Tranche

Date: 2026-04-05

## Summary

This tranche widens the `../superai` → `hyperharness` assimilation effort into the HyperCode/Borg adapter boundary inside:
- `foundation/adapters/hypercode.go`
- related provider-visibility helpers in `foundation/adapters/providers.go`

This is a high-value seam because it translates real runtime/provider/MCP/discovery state into model-facing system context.

## Key finding

`foundation/adapters/hypercode.go` was already present in `hyperharness` and effectively parity-aligned with `../superai`.

The highest-value work here was therefore not redundant copying. It was to:
- preserve the same adapter-facing behavior
- harden nil-adapter and discovery-edge behavior
- extract repo-candidate and context-rendering seams
- add regression coverage for provider-visibility ordering and system-context composition

## What changed

### 1. Nil-safe adapter status behavior
`(*HyperCodeAdapter).Status()` now handles a nil receiver explicitly and returns a warning instead of panicking.

This improves safety at an important integration boundary.

### 2. System-context rendering extraction
Added:
- `renderHyperCodeSystemContext(...)`

This isolates the final system-context rendering logic into a directly testable helper.

### 3. HyperCode repo candidate extraction
Added:
- `hyperCodeRepoCandidates(workingDir, homeDir)`

Behavioral improvements:
- candidate generation is explicit
- duplicate candidates are removed
- candidate ordering is deterministic

This makes adjacent-repo discovery easier to reason about and test.

### 4. Safer MCP config listing path
`listMCPServers()` now tolerates a nil adapter by resolving from an empty/default home path rather than assuming a fully initialized receiver.

### 5. Provider helper regression coverage
Added direct coverage for:
- `detectAvailableProviders(...)`
- `sortStrings(...)`

This matters because provider visibility becomes part of the model-facing adapter context.

## Regression coverage added

Added:
- `foundation/adapters/hypercode_assimilation_test.go`
- `foundation/adapters/providers_assimilation_test.go`

Coverage includes:
- repo-candidate deduplication and ordering
- rendered system-context content
- nil-adapter status/context behavior
- adjacent repo discovery behavior
- provider detection inclusion and sorted ordering
- string sorting helper behavior

## Validation

Verified successfully:

```bash
go test ./foundation/adapters ./tools ./agent ./agents ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because the adapter boundary is where multiple important truths are assembled together:
- Borg assimilation state
- provider visibility
- MCP configuration visibility
- HyperCode repo discovery
- model-facing system context

Strengthening that seam improves the reliability of everything built on top of it.

## Design insight

A useful lesson from this tranche is:

> integration adapters become much easier to trust once discovery, rendering, and fallback behavior are represented by explicit helper seams instead of being embedded implicitly inside one status method.

That is exactly what this tranche improved.

## Recommended next move

The strongest next assimilation step is now:

1. continue through remaining runtime/integration seams where tool exposure, adapters, and CLI flows meet
2. keep prioritizing boundary files that assemble exact-name tool context for model-facing use
3. keep validating with:
   - `go test ./foundation/adapters ./tools ./agent ./agents ./tui ./cmd ./foundation/...`

## Bottom line

This tranche confirms that `hyperharness` already contained the `superai` HyperCode/Borg adapter surface, and upgrades it into a clearer, safer, regression-tested base for the broader Go-native assimilation effort.
