# TUI Pane Show/Hide Tranche

Date: 2026-04-04

## Summary

This tranche adds **explicit show/hide commands** for the persistent tree pane.

Before this tranche, pane visibility was controlled via a toggle command (`/tree-pane`). That worked, but explicit show/hide controls are better for operator confidence and scripting-style repeatability.

## What was added

### New slash commands
- `/tree-pane-show`
- `/tree-pane-hide`

### Behavior
- `/tree-pane-show`
  - pins the pane if it is not already visible
  - emits an informational message if it is already visible
- `/tree-pane-hide`
  - hides the pane if it is visible
  - emits an informational message if it is already hidden

Both commands preserve the same canonical runtime-backed pane behavior and clean up local focus/confirm state appropriately on hide.

## Verification added

Added focused regression coverage verifying:
- explicit show makes the pane visible
- explicit hide makes the pane hidden again

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because explicit controls reduce ambiguity. Toggle-only behavior is compact, but show/hide commands are clearer in repeated use and easier to reason about.

## Design insight

The key principle remains the same:

> explicit show/hide are operator-surface affordances layered over the same underlying pane state; they do not change session or branch semantics.

## Recommended next move

The strongest next step is still to continue adding only those small operator controls that clearly reduce repeated friction.

## Bottom line

This tranche improves control clarity around pane visibility without adding any new runtime complexity.
