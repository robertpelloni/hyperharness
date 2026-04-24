# TUI Pane Position Quick Toggle Tranche

Date: 2026-04-04

## Summary

This tranche adds a **quick top/bottom position toggle** for the persistent tree pane.

Before this tranche, pane position could already be controlled explicitly with:
- `/tree-pane-position top`
- `/tree-pane-position bottom`

That was flexible but slightly verbose for frequent adjustment.

This tranche adds a one-shot convenience toggle.

## What was added

### 1. New slash command
Added:
- `/tree-pane-position-toggle`

### 2. Position handler toggle support
The underlying pane-position handler now also supports:
- `toggle`
- empty argument as toggle behavior

This allows both:
- explicit top/bottom placement
- and fast placement flipping

## Verification added

Added focused regression coverage verifying:
- first quick toggle moves the pane from top to bottom
- second quick toggle moves it back from bottom to top

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because small layout toggles become valuable once the pane is already highly configurable and frequently used during long sessions.

The command reduces friction for one of the most common split-view adjustments.

## Design insight

The key principle remains intact:

> quick position toggling is just a convenience layer over the same pane layout state.

No runtime/session semantics change.

## Recommended next move

The strongest next step is now:

> continue adding small convenience controls where they reduce repeated pane/browser friction during long interactive sessions.

## Bottom line

This tranche is a compact but useful operator-polish improvement that makes split-view adjustment faster in practice.
