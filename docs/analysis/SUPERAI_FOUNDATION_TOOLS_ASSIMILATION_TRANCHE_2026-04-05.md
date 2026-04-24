# SuperAI Foundation Tools Assimilation Tranche

Date: 2026-04-05

## Summary

This tranche ports a safe, high-value slice of `../superai` into the canonical `hyperharness` foundation layer.

The focus was deliberately narrow:
- keep `foundation/pi` as the canonical implementation truth
- import richer native-tool behavior where the semantics are clearly compatible
- keep the core validation target green throughout

## What was assimilated

### 1. Richer `foundation/pi` native tool behavior
Ported the useful `superai` behavior for:
- `grep`
- `find`
- `ls`

into `hyperharness/foundation/pi`.

This added support for:
- explicit no-result messages
  - `grep`: `No matches found`
  - `find`: `No files found matching pattern`
  - `ls`: `(empty directory)`
- result-limit details
  - `GrepToolDetails.MatchLimitReached`
  - `FindToolDetails.ResultLimitReached`
  - `LsToolDetails.EntryLimitReached`
- truncation-aware result payloads for these tools
- sorted `ls` output
- relative subdirectory matching in `find`
- line truncation helpers for `grep`
- helper functions for path normalization and file-line reading

### 2. New runtime detail types
Added detail payload types in `foundation/pi/runtime_types.go`:
- `GrepToolDetails`
- `FindToolDetails`
- `LsToolDetails`

These keep the result structure explicit and testable.

### 3. Buildability fix during module identity migration
While continuing the broader assimilation work, an import-path mismatch was discovered in:
- `cmd/hyperharness/main.go`
- `internal/agent/runtime.go`

Those files still referenced:
- `github.com/hyperharness/hyperharness/...`

They were corrected to:
- `github.com/robertpelloni/hyperharness/...`

This was necessary to keep the repository buildable while the larger identity normalization remains in progress.

## Verification added

Added focused regression tests in:
- `foundation/pi/tools_extra_port_test.go`

Coverage includes:
- grep no-match messaging
- find no-result messaging
- ls empty-directory messaging
- grep limit detail reporting
- find limit detail reporting
- ls limit detail reporting
- find matching relative subdirectory patterns like `src/*.go`

## Validation

Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

All targeted core packages remained green.

## Design insight

This tranche is important because it demonstrates the correct assimilation pattern:

> import the strongest compatible behavior from `superai` into the canonical `hyperharness` foundation layer, while preserving existing contracts and keeping the test surface green.

In other words, this is not parallel-track drift.
It is controlled convergence.

## Why this was the right next step

Compared with larger agent/autonomy ports, this tranche had the best risk-adjusted value:
- low blast radius
- directly improves canonical tool behavior
- easy to validate
- improves parity and operator confidence
- does not entangle unfinished agent-side architecture yet

## Recommended next move

The highest-value next assimilation step is now:

1. port `../superai/foundation/pi/tools_extra.go` behavior only where it still exceeds `hyperharness`
2. then move one level outward into `agent/` assimilation, starting with the safest utilities:
   - `context.go`
   - `compare.go`
   - `diff.go`
3. keep validation strict after each tranche:
   - `go test ./tui ./cmd ./foundation/...`

## Bottom line

This tranche successfully imported a meaningful, testable slice of `superai` into `hyperharness` without destabilizing the canonical foundation track.
