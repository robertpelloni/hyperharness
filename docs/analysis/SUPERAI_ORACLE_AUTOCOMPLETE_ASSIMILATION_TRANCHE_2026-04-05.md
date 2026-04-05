# SuperAI Oracle & Autocomplete Assimilation Tranche

Date: 2026-04-05

## Summary

This tranche continues the `../superai` → `hyperharness` assimilation effort at the next agent-utility layer:
- `agent/oracle.go`
- `agent/autocomplete.go`

## Key finding

As with the prior utility tranche, these files were already present in `hyperharness` and were effectively parity-aligned with `../superai`.

The right work here was therefore not a wholesale copy.
It was to:
- verify parity
- extract and stabilize the essential prompt-building/orchestration behavior
- add defensive validation for obvious failure edges
- lock it in with regression coverage

## What changed

### 1. `agent/oracle.go` hardening and extraction
`OracleQuery` now delegates through a small internal helper:
- `oracleQueryWithChat(prompt, chat)`

This made the orchestration logic directly testable without requiring network access.

Additional helper functions were introduced for prompt construction:
- `buildOraclePlanPrompt`
- `buildOracleExecutionPrompt`
- `buildOracleSynthesisPrompt`

Behavioral improvements:
- rejects empty prompts
- rejects missing chat function
- preserves the existing three-stage Oracle flow:
  1. plan
  2. execute/gather
  3. synthesize

### 2. `agent/autocomplete.go` hardening
`SuggestCompletion` now:
- rejects missing/nil OpenAI client
- checks for empty completion choices before indexing into the response

A helper was also extracted:
- `buildCompletionPrompt`

This makes the key prompt contract testable and explicit.

## Regression coverage added

Added:
- `agent/oracle_autocomplete_assimilation_test.go`

Coverage includes:
- verifying the full three-stage Oracle prompt flow
- validating empty-prompt rejection in Oracle mode
- validating error propagation from Oracle stages
- validating completion prompt construction
- validating missing-client rejection for autocomplete

## Validation

Verified successfully:

```bash
go test ./agent ./agents ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because it upgrades more of the already-assimilated `superai` surface from:
- present but unguarded

to:
- present
- explicit
- safer
- regression-tested

That is especially valuable in agent-side code, where hidden prompt/orchestration assumptions can drift silently.

## Design insight

This tranche reinforces the same assimilation principle as the previous one:

> when sibling repositories already share code, the most valuable work is often to make parity testable and resilient rather than merely copying files again.

## Recommended next move

The strongest next assimilation step is now:

1. inspect `agent/async.go` and `agent/shell.go`
2. identify where `superai` materially exceeds current `hyperharness` behavior
3. assimilate the safest portions first with regression coverage
4. keep validating with:
   - `go test ./agent ./agents ./tui ./cmd ./foundation/...`

## Bottom line

This tranche confirms that `hyperharness` already contained the `superai` Oracle and autocomplete surface, and it upgrades that parity into explicit, safer, tested behavior.
