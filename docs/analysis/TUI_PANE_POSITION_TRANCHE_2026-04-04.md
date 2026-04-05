# TUI Pane Position Tranche

Date: 2026-04-04

## Summary

This tranche adds **persistent pane position control** to the TUI so operators can choose whether the live tree pane renders above or below the main prompt/history flow.

Before this tranche, the pane was:
- persistent
- live
- focusable
- viewport-aware
- and size-adjustable

But its placement was fixed.

This tranche makes the layout more adaptable.

## What was added

### 1. Pane position state
Added to the TUI model:
- `browserPanePosition string`

Default:
- `"top"`

### 2. New slash command
Added in `tui/slash.go`:
- `/tree-pane-position <top|bottom>`

Behavior:
- accepts `top` or `bottom`
- updates pane layout preference
- emits a status message confirming the new position

### 3. View integration
Updated `model.View()` so that when the pane is pinned:
- it renders above the main prompt/history flow when `top`
- it renders below the main prompt/history flow when `bottom`

This applies without changing any underlying runtime/browser semantics.

## Verification added

Added a focused regression test:
- `TestProcessSlashCommandTreePanePosition`

This verifies:
- the command updates position state correctly
- the expected status message is emitted

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche is important because the pane is now not only persistent and live, but also more adaptable to different operator preferences and terminal layouts.

That is a meaningful split-view ergonomics improvement.

## Design insight

The key architectural rule remains intact:

> pane position is purely a layout concern layered over the same canonical runtime-backed browser state.

That is the correct boundary.

The operator gets more control over the interface without any duplication of session/branch truth.

## Recommended next move

The strongest next step is now:

> **continue improving split-view ergonomics**, likely with better pane/history separation or stronger viewport/navigation affordances inside the pane.

## Bottom line

This tranche makes the persistent tree pane more flexible by allowing top/bottom placement, improving coexistence with the main prompt/history flow while preserving the same truth-first foundation architecture.
