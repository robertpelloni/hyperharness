# TUI Preview Explicit On/Off Tranche

Date: 2026-04-04

## Summary

This tranche adds **explicit preview on/off aliases** for the persistent tree pane.

Before this tranche, preview visibility could be controlled with:
- `/tree-pane-preview on`
- `/tree-pane-preview off`
- `/tree-pane-preview-toggle`

That was already flexible, but explicit aliases make the command family more consistent with the newer explicit grouped/focus controls.

## What was added

### New slash commands
- `/tree-pane-preview-on`
- `/tree-pane-preview-off`

These route directly into the existing preview-state handler.

## Verification added

Added focused regression coverage verifying:
- explicit preview-on enables preview
- explicit preview-off disables preview

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because consistency across the pane command family reduces mental load and increases confidence in repeated use.

## Design insight

The key principle remains unchanged:

> explicit preview aliases are just direct entry points into the same preview-state transitions.

## Recommended next move

The strongest next step is now:

> continue only where explicit aliases clearly improve confidence or consistency over existing controls.

## Bottom line

This is a small but coherent polish improvement that makes preview control more consistent with the rest of the pane command family.
