# TUI Pane Preview Quick Toggle Tranche

Date: 2026-04-04

## Summary

This tranche adds a **quick preview toggle** for the persistent tree pane.

Before this tranche, operators could control pane preview visibility with:
- `/tree-pane-preview on`
- `/tree-pane-preview off`

That was explicit but slightly verbose for repeated toggling.

This tranche adds a one-shot convenience toggle.

## What was added

### 1. New slash command
Added:
- `/tree-pane-preview-toggle`

### 2. Preview toggle behavior
The underlying preview handler now also supports:
- `toggle`
- empty arg as toggle behavior

This allows both:
- explicit on/off control
- and fast toggle control

## Verification added

Added focused regression coverage verifying:
- first quick toggle turns preview off
- second quick toggle turns preview back on

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because small operator convenience controls compound in value once a pane/browser subsystem becomes rich enough to adjust frequently.

A quick preview toggle reduces friction during long sessions when switching between compact and detailed inspection modes.

## Design insight

The key principle remains unchanged:

> quick toggles are convenience layers over the same pane layout state, not new semantics.

## Recommended next move

The strongest next step is now:

> continue adding small operator-polish controls where they reduce repeated command friction, especially around pane/browser display state.

## Bottom line

This tranche is a small but useful ergonomics refinement that makes preview visibility faster to control during real interactive use.
