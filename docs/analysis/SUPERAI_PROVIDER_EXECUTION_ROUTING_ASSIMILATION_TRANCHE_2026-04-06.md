# SuperAI Provider Execution & Routing Assimilation Tranche

Date: 2026-04-06

## Summary

This tranche continues the `../superai` → `hyperharness` assimilation effort at a critical runtime/integration seam:
- `foundation/adapters/provider_execution.go`
- `foundation/adapters/provider_routing.go`

This is a high-value boundary because it sits where:
- prompt intent detection
- cost-preference routing
- task-type classification
- default model selection
- provider availability
- execution hints for agents

all resolve into concrete routing decisions.

## Key finding

Both files were already present in `hyperharness` and effectively parity-aligned with `../superai`.

The highest-value work here was therefore not redundant copying. It was to:
- preserve the same provider-facing behavior
- extract implicit normalization/selection logic into explicit helpers
- make default model constants explicit
- add direct regression coverage for task/cost/model normalization

## What changed

### 1. Provider execution helper extraction
Added:
- `normalizeTaskType(...)`
- `promptPreview(...)`
- `buildExecutionHint(...)`

This isolates intent detection, preview truncation, and hint construction into directly testable seams.

### 2. Provider routing helper extraction
Added explicit default model constants:
- `defaultOpenAIModel`
- `defaultGoogleModel`
- `defaultOllamaModel`

Added routing normalization helpers:
- `normalizedPreference(...)`
- `normalizedTaskType(...)`
- `firstOr(...)`
- `defaultModelForProvider(...)`

This isolates cost-preference classification, task-type mapping, and fallback model resolution into explicit, testable seams.

### 3. Behavior clarified, not forked
The intended routing behavior remains exactly the same:
- local preference routes to ollama when available
- budget preference routes to google/flash
- quality preference routes to openai/gpt-4o
- task type detection remains identical
- fallback defaults remain openai -> gpt-4o

The difference is that the decision logic is now explicit and regression-tested.

## Regression coverage added

Added:
- `foundation/adapters/provider_execution_assimilation_test.go`
- `foundation/adapters/provider_routing_assimilation_test.go`

Coverage includes:
- explicit vs inferred task-type normalization
- prompt preview truncation at 140 characters
- execution hint formatting
- full `inferTaskType` keyword coverage
- cost-preference normalization ("cheap"/"low" -> "budget", "high"/"best" -> "quality")
- task-type normalization ("code"/"refactor" -> "coding", etc.)
- `firstOr` fallback behavior
- per-provider default model resolution
- empty/blank input handling

## Validation

Verified successfully:

```bash
go test ./foundation/adapters ./tools ./agent ./agents ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because provider routing sits exactly at the center of the harness decision loop.

By normalizing preferences, task types, and model selection through explicit helpers, `hyperharness` gains:
- clearer routing semantics
- directly testable preference classification
- explicit default model contracts
- safer fallback behavior
- stronger confidence that provider routing decisions remain stable as the harness expands

## Design insight

A useful lesson from this tranche is:

> routing and cost-preference logic becomes dramatically easier to trust once keyword classification, provider mapping, and default model selection are represented by small normalization helpers instead of being embedded directly inside a large switch/case block.

That is exactly what this tranche accomplished.

## Bottom line

This tranche confirms that `hyperharness` already contained the `superai` provider execution and routing surface, and upgrades it into a clearer, safer, regression-tested base for the broader Go-native assimilation effort.
