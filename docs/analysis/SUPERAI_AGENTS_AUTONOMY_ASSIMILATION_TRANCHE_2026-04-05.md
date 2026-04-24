# SuperAI Agents Autonomy Assimilation Tranche

Date: 2026-04-05

## Summary

This tranche continues the `../superai` → `hyperharness` assimilation effort with richer `agents/`-level autonomy and shell coordination files:
- `agents/council.go`
- `agents/auton.go`
- `agents/shell_assistant.go`

## Key finding

These files were already present in `hyperharness` and effectively parity-aligned with `../superai`.

The highest-value work here was therefore not redundant copying. It was to:
- preserve the same surface and intended role
- improve determinism where needed
- harden autonomy/control-loop edge cases
- extract testable helper seams
- add regression coverage

## What changed

### 1. `agents/council.go` hardening and determinism
Added helper seams:
- `defaultCouncilPersonas()`
- `councilApprovedFromReason(...)`
- `councilFinalVerdict(...)`

Behavioral improvements:
- approval detection now trims leading/trailing whitespace before checking `VOTE: APPROVED`
- council results are sorted by persona after concurrent evaluation, making output deterministic

This is a real quality improvement because goroutine completion order no longer leaks into returned vote ordering.

### 2. `agents/auton.go` hardening and helper extraction
Added helper seams:
- `autoDriveStart(...)`
- `buildAutoDrivePrompt(...)`
- `autoDriveToolResultMessage(...)`

Added validation for:
- nil autodrive receiver
- missing director
- missing director provider
- empty objective
- empty sandbox directory
- non-positive max iterations

Behavioral improvement:
- aborts that happen during the wait window are now respected before the next chat round proceeds
- `Abort()` is now nil-safe
- `IsRunning` is reliably reset via `defer`

This makes AutoDrive safer and more truthful as an autonomy loop.

### 3. `agents/shell_assistant.go` hardening
Added:
- nil/missing provider validation
- empty intent validation
- `buildShellTranslationMessages(...)` helper extraction

This makes the shell translation contract explicit and testable.

## Regression coverage added

Added:
- `agents/autonomy_assimilation_test.go`

Coverage includes:
- council approval detection and consensus calculation
- council persona defaults
- AutoDrive validation failures
- AutoDrive prompt/tool-result history behavior
- AutoDrive abort propagation
- AutoDrive chat error propagation
- shell translator validation
- shell translation message construction

## Validation

Verified successfully:

```bash
go test ./agent ./agents ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because it strengthens some of the highest-risk shared surfaces in the current sibling-repo assimilation work:
- concurrent review voting
- autonomous execution loops
- natural-language-to-shell translation

By making these surfaces more deterministic and better validated, `hyperharness` becomes a stronger canonical Go-native base for deeper autonomy work.

## Design insight

A useful lesson from this tranche is:

> once higher-autonomy code is present, the most important improvements are often not new features but stronger control-loop contracts and deterministic outputs.

That makes future expansion safer.

## Recommended next move

The strongest next assimilation step is now:

1. continue through the remaining richer `agents/` surfaces
2. inspect next files such as:
   - `agents/codemode.go`
   - `agents/rag.go`
   - `agents/provider_hypercode.go`
3. keep validating with:
   - `go test ./agent ./agents ./tui ./cmd ./foundation/...`

## Bottom line

This tranche confirms that `hyperharness` already contained the `superai` `agents/` autonomy surface, and upgrades it into a more deterministic, safer, regression-tested base for continued Go-native assimilation.
