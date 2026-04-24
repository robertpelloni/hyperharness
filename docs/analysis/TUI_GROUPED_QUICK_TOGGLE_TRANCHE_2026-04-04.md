# TUI Grouped Quick Toggle Tranche

Date: 2026-04-04

## Summary

This tranche adds a **quick grouped-mode toggle** for the persistent tree pane.

Before this tranche, grouped rendering could already be controlled explicitly with:
- `/tree-pane-grouped on`
- `/tree-pane-grouped off`
- `/tree-pane-grouped toggle`

That was already flexible, but a dedicated one-shot toggle keeps parity with other quick pane controls.

## What was added

### 1. New slash command
Added:
- `/tree-pane-grouped-toggle`

### 2. Behavior
The command simply toggles grouped mode on/off through the same underlying grouped-state handler.

This keeps the implementation small and consistent while reducing command friction.

## Verification added

Added focused regression coverage verifying:
- first grouped quick toggle enables grouped mode
- second grouped quick toggle disables grouped mode

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because the tree pane/browser subsystem now has enough operator controls that consistency between them becomes valuable.

Preview already had a quick toggle; grouped mode now has one too.

## Design insight

The key principle remains unchanged:

> quick toggles are operator conveniences layered over the same existing pane state, not new semantics.

## Recommended next move

The strongest next step is now:

> continue adding only those small convenience controls that clearly reduce repeated operator friction.

## Bottom line

This is a compact but worthwhile polish improvement that makes grouped rendering faster to control during real use.
