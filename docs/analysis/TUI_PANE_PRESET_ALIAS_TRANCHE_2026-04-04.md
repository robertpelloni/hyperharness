# TUI Pane Preset Alias Tranche

Date: 2026-04-04

## Summary

This tranche adds **direct preset aliases** for the persistent tree pane.

Before this tranche, operators could apply named pane presets via:
- `/tree-pane-preset compact`
- `/tree-pane-preset detailed`
- `/tree-pane-preset navigation`
- `/tree-pane-preset review`

That was flexible, but explicit preset aliases make the most common layouts even faster to reach.

## What was added

### New slash commands
Added:
- `/tree-pane-compact`
- `/tree-pane-detailed`
- `/tree-pane-navigation`
- `/tree-pane-review`

Each delegates directly to the existing preset logic.

## Verification added

Added focused regression coverage verifying that each alias applies the expected preset state.

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because once the pane/browser subsystem accumulates enough controls, the difference between “available” and “fast to use” becomes important.

Direct aliases reduce command friction for common workflows.

## Design insight

The key principle remains the same:

> aliases are operator ergonomics only; the underlying pane semantics stay centralized in the existing preset implementation.

## Recommended next move

The strongest next step is now:

> continue adding only those shortcut controls that clearly improve repeated-use ergonomics without fragmenting behavior.

## Bottom line

This tranche gives the pane system a more direct, workflow-friendly command surface for its most important layout presets.
