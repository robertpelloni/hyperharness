# TUI Pane Preset Tranche

Date: 2026-04-04

## Summary

This tranche adds **named pane layout presets** so operators can quickly switch the persistent tree pane between compact and detailed configurations.

Before this tranche, pane ergonomics could already be tuned through separate commands for:
- size
- position
- preview visibility

That was flexible, but verbose.

This tranche adds a faster operator-friendly layer on top of those controls.

## What was added

### 1. New slash command
Added in `tui/slash.go`:
- `/tree-pane-preset <compact|detailed>`

### 2. Compact preset
Applies:
- `browserPaneHeight = 6`
- `browserPanePreview = false`
- `browserPanePosition = "bottom"`

This favors a smaller, denser pane that stays out of the way.

### 3. Detailed preset
Applies:
- `browserPaneHeight = 12`
- `browserPanePreview = true`
- `browserPanePosition = "top"`

This favors a larger, more informative pane with richer preview visibility.

### 4. Help output updated
The TUI help text now documents `/tree-pane-preset`.

## Verification
Added focused coverage in `tui/slash_test.go` to verify:
- the compact preset applies the expected state changes
- the detailed preset applies the expected state changes

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because it reduces friction for operators who want different pane personalities for different workflows.

Instead of adjusting multiple pane settings one by one, they can now switch between a compact and a detailed pane with a single command.

## Design insight

The key principle remains intact:

> pane presets are just operator-facing bundles of existing rendering/layout settings over the same canonical runtime-backed browser state.

No session or branch semantics change.

## Recommended next move

The strongest next step is now:

> **continue refining pane/browser coexistence with higher-level operator presets or split-view controls**, since the pane has now become highly configurable.

## Bottom line

This tranche adds a useful operator shortcut layer on top of the existing pane controls, making the TUI tree pane more practical in day-to-day use.
