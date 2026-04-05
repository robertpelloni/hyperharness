# TUI Grouped Explicit On/Off Tranche

Date: 2026-04-04

## Summary

This tranche adds **explicit grouped on/off aliases** for the persistent tree pane.

Before this tranche, grouped rendering could be controlled with:
- `/tree-pane-grouped on`
- `/tree-pane-grouped off`
- `/tree-pane-grouped toggle`
- `/tree-pane-grouped-toggle`

That was already flexible, but explicit aliases improve consistency with the newer focus on/off controls.

## What was added

### New slash commands
- `/tree-pane-grouped-on`
- `/tree-pane-grouped-off`

These map directly to the existing grouped-state handler.

## Verification added

Added focused regression coverage verifying:
- explicit grouped-on enables grouped mode
- explicit grouped-off disables grouped mode

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because the pane/browser subsystem now has enough controls that consistency between toggle-style and explicit-style commands becomes an ergonomic benefit.

## Design insight

The key principle remains unchanged:

> explicit grouped on/off aliases are just clearer entry points into the same grouped-state transitions.

## Recommended next move

The strongest next step is now:

> continue only where explicit aliases clearly improve confidence or consistency over existing toggle-style controls.

## Bottom line

This tranche is a small consistency/polish improvement that makes grouped-mode control feel more coherent alongside the rest of the pane command family.
