# SuperAI Agent Utilities Assimilation Tranche

Date: 2026-04-05

## Summary

This tranche continues the `../superai` → `hyperharness` assimilation effort at the next-lowest-risk layer after `foundation/pi` native tools.

The target files were:
- `agent/context.go`
- `agent/compare.go`
- `agent/diff.go`

## Key finding

A useful discovery from this tranche is that these utility files were already present in `hyperharness` and were effectively aligned with `../superai`.

That means the highest-value work here was not a blind copy.
It was:
- verify parity
- harden edge cases
- add regression coverage so the assimilated behavior stays preserved

## What changed

### 1. `agent/context.go` hardening
`TrimHistory` now safely handles edge cases that could otherwise produce surprising behavior or indexing problems:
- non-positive `maxMessages`
- empty history
- already-small histories
- explicit `maxMessages == 1`

The core behavior remains the same:
- preserve the system prompt
- preserve the most recent messages

### 2. `agent/diff.go` hardening
`ApplyInlineDiff` now ignores unified-diff metadata lines such as:
- `--- ...`
- `+++ ...`
- `@@ ...`

It also applies a slightly safer removal path when the removed content is not followed by a newline.

This keeps the implementation intentionally simple while avoiding obvious corruption from unified-diff headers being written into the target file.

### 3. Regression coverage added
Added:
- `agent/utilities_assimilation_test.go`

Coverage includes:
- preserving the system prompt and newest messages during trim
- ensuring non-positive trim limits are a no-op
- ensuring inline diff application ignores unified diff headers and performs the expected add/remove behavior

## Validation

Verified successfully:

```bash
go test ./agent ./agents ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because it turns implicit parity into explicit, defended parity.

Without tests, these already-assimilated utilities could drift silently.
With this tranche:
- the behavior is now documented
- the edge cases are safer
- the `superai` parity is pinned by regression coverage

## Design insight

This tranche reinforces an important assimilation principle:

> not every valuable assimilation step is a big code import; sometimes the right move is to prove parity, harden the edges, and lock it in with tests.

That is especially important when consolidating sibling repositories with overlapping implementations.

## Recommended next move

The strongest next agent-side assimilation step is now:

1. inspect `agent/oracle.go` and `agent/autocomplete.go`
2. port or harden only the parts that materially exceed the current `hyperharness` behavior
3. continue validating with:
   - `go test ./agent ./agents ./tui ./cmd ./foundation/...`

## Bottom line

This tranche confirms that `hyperharness` already contained part of the `superai` agent utility surface, and it upgrades that parity from informal to tested and safer.
