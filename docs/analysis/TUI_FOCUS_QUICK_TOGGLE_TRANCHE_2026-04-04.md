# TUI Focus Quick Toggle Tranche

Date: 2026-04-04

## Summary

This tranche adds a **quick focus-toggle alias** for the persistent tree pane.

Before this tranche, keyboard focus for the pinned pane could already be controlled with:
- `/tree-pane-focus`

This tranche adds a more explicit convenience alias in the same style as other quick pane controls.

## What was added

### 1. New slash command
Added:
- `/tree-pane-focus-toggle`

### 2. Behavior
The command delegates directly to the existing focus-toggle behavior, so it remains a thin ergonomic alias rather than a separate state path.

## Verification added

Added focused regression coverage verifying:
- first quick focus toggle enables pane focus
- second quick focus toggle disables pane focus

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because the pane/browser subsystem now has a family of quick toggles, and consistency across those controls improves usability.

## Design insight

The key principle remains unchanged:

> quick-toggle aliases are ergonomic wrappers over existing state transitions, not new semantics.

## Recommended next move

The strongest next step is now:

> continue only with small convenience additions that clearly reduce real operator friction.

## Bottom line

This is a small but consistent improvement that keeps pane controls feeling coherent as the subsystem grows.
