# TUI Reset Controls Tranche

Date: 2026-04-04

## Summary

This tranche adds **reset and quick-clear controls** for the tree pane/browser subsystem.

Before this tranche, the pane/browser had become quite configurable:
- height
- position
- preview toggle
- presets
- grouping
- filter state
- collapse state
- confirm-before-switch state
- focus state

That made the subsystem powerful, but it also increased the value of fast recovery commands.

This tranche adds those recovery controls.

## What was added

### 1. Browser transient-state clear command
Added:
- `/tree-browser-clear`

This clears transient browser state:
- filter
- pending confirmation
- collapse state
- selection index

This is useful when the browser has become locally messy and the operator wants a quick clean slate without touching the canonical session.

### 2. Pane reset command
Added:
- `/tree-pane-reset`

This restores pane configuration defaults:
- height → `8`
- position → `top`
- preview → `true`
- grouped → `false`
- focus → `false`
- confirm pending → `false`
- filter → cleared
- collapsed state → cleared

This gives the operator a fast recovery path back to the default pane behavior.

## Verification added

Expanded `tui/slash_test.go` with focused regression coverage verifying:
- browser transient state clears correctly
- pane configuration resets correctly

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because as the pane/browser subsystem becomes more configurable, recovery ergonomics become important.

Operators need a fast way to return to a known-good state without restarting the TUI or manually undoing multiple toggles.

## Design insight

The key principle remains intact:

> reset/clear commands only affect UI/layout/control state; they do not alter canonical session or branch truth.

That is exactly the right boundary.

## Recommended next move

The strongest next step is now:

> continue improving operator control around the pane/browser subsystem, likely with additional quick toggles or higher-level workflow presets.

## Bottom line

This tranche gives the tree pane/browser subsystem a better recovery story, which is important once a UI grows rich enough to accumulate meaningful local state.
