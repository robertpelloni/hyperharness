# TUI Pane Size Tranche

Date: 2026-04-04

## Summary

This tranche adds **adjustable persistent tree pane height** to the TUI.

Before this tranche, the pinned pane had a fixed viewport height. That was useful, but too rigid for different session sizes and operator preferences.

This tranche makes the pane height explicitly configurable from the TUI command surface.

## What was added

### 1. Pane height state
Added to the TUI model:
- `browserPaneHeight int`

Default behavior:
- initialized to `8`
- safely falls back to `8` when unset or invalid

### 2. Pane height command
Added in `tui/slash.go`:
- `/tree-pane-size <n>`

Behavior:
- accepts a positive integer
- updates the persistent pane viewport height
- emits a status message confirming the new height

### 3. View integration
Updated the pinned pane render path so it now uses:
- `browserPaneHeight`
- with a fallback to `8`

That means the same viewport-aware renderer can be reused with a user-adjustable size.

## Verification added

Added a focused TUI test:
- `TestProcessSlashCommandTreePaneSize`

This verifies:
- pane height updates through the slash command
- the status message is emitted correctly

### Validation
Verified successfully:

```bash
go test ./tui
go test ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche improves long-session usability by making the persistent pane tunable instead of fixed.

That is a meaningful ergonomics improvement because different repositories and session shapes need different viewport sizes.

## Design insight

The key architectural principle remains unchanged:

> pane size is purely a rendering/view concern layered over canonical runtime-derived browser state.

The operator can tune visibility without changing any underlying branch/session semantics.

## Recommended next move

The strongest next step is now:

> **continue strengthening pane/browser coexistence ergonomics**, likely with refresh controls, layout toggles, or additional pane-focused navigation affordances.

## Bottom line

This tranche makes the persistent tree pane more practical by letting the operator control how much of the session tree is visible at once.
