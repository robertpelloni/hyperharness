# TUI Grouped Toggle Tranche

Date: 2026-04-04

## Summary

This tranche adds an explicit **grouped-rendering toggle command** for the persistent tree pane.

Before this tranche, grouped rendering could be reached through:
- keyboard `Tab` in browser/focused-pane contexts
- or indirectly through some presets

That was useful, but not ideal for direct operator control from slash commands.

This tranche adds a simple command-level control.

## What was added

### 1. New slash command
Added:
- `/tree-pane-grouped <on|off|toggle>`

### 2. Behavior
The command supports:
- `on`
- `off`
- `toggle`
- empty arg behaves like toggle

It updates `browserGrouped` directly and emits a confirmation message.

## Verification added

Added focused regression coverage to verify:
- grouped mode turns on
- grouped mode toggles back off

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because it gives operators a faster, more explicit control over one of the more important pane/browser layout modes.

That fits the larger direction of the current work: making the pane subsystem easier to control and recover without sacrificing the canonical runtime architecture.

## Design insight

The key principle remains the same:

> grouped mode is still just a rendering/layout concern over canonical runtime-backed browser state.

The command improves control, not semantics.

## Recommended next move

The strongest next step is now:

> continue adding small operator-control affordances where they reduce friction, especially around pane/browser layout and state recovery.

## Bottom line

This tranche is a small but useful operator-polish improvement that makes grouped rendering easier to control directly from the TUI command surface.
