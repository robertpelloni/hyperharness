# SuperAI Agents Provider Assimilation Tranche

Date: 2026-04-05

## Summary

This tranche continues the `../superai` → `hyperharness` assimilation effort with the remaining provider-adjacent `agents/` surfaces:
- `agents/provider.go`
- `agents/provider_stub.go`

## Key finding

These files were already present in `hyperharness` and effectively parity-aligned with `../superai`.

The highest-value work here was therefore not redundant copying. It was to:
- preserve the same provider-facing surface
- harden nil/stream edge handling
- extract stable default/helper seams
- add regression coverage for provider contract behavior

## What changed

### 1. `agents/provider.go` hardening
Added explicit defaults:
- `defaultGeminiBaseURL`
- `defaultGeminiModel`

Added helper seams:
- `defaultLegacyTools()`
- `lastProviderMessageContent(...)`

Behavioral improvements:
- nil provider validation in `Chat`
- nil provider validation in `Stream`
- nil chunk-channel validation in `Stream`
- default model fallback in `GetModelName()` when unset
- `FetchLegacyToolArray()` now returns a fresh copy via helper defaults
- `NewGeminiBorgProvider()` now uses explicit constants for base URL and model

These changes make the provider contract more explicit and easier to test without altering the intended high-level behavior.

### 2. `agents/provider_stub.go` hardening
Behavioral improvements:
- nil provider validation in `Chat`
- nil provider validation in `Stream`
- nil chunk-channel validation in `Stream`
- shared `lastProviderMessageContent(...)` helper used for prompt extraction

This turns the local stub provider into a clearer and more defensive contract surface.

## Regression coverage added

Added:
- `agents/provider_assimilation_test.go`

Coverage includes:
- `lastProviderMessageContent(...)` behavior
- fresh-copy behavior for legacy tool defaults
- Gemini provider validation and default-model behavior
- Gemini provider streaming behavior
- `NewGeminiBorgProvider()` disclosure-wrapper behavior
- default provider validation and normal response behavior

## Validation

Verified successfully:

```bash
go test ./agent ./agents ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because the provider surface sits underneath much of the recently-hardened `agents/` coordination and autonomy behavior.

By making the provider contracts safer and more explicit, `hyperharness` gains:
- more predictable behavior under nil/misconfigured paths
- clearer default configuration semantics
- better preservation of the disclosure-wrapper model
- stronger confidence before moving into deeper provider or model routing work

## Design insight

A useful lesson from this tranche is:

> provider contracts become much easier to evolve safely once defaults, prompt extraction, and disclosure wrapping are represented by explicit helper seams rather than implicit assumptions.

That matters for long-term exact-contract compatibility.

## Recommended next move

The strongest next assimilation step is now:

1. continue through any remaining richer `agents/` surfaces that have not yet been hardened
2. then widen back out to adjacent runtime/provider integration layers if needed
3. keep validating with:
   - `go test ./agent ./agents ./tui ./cmd ./foundation/...`

## Bottom line

This tranche confirms that `hyperharness` already contained the `superai` provider-adjacent `agents/` surface, and upgrades it into a clearer, safer, regression-tested base for continued Go-native assimilation.
