# TUI Focus Explicit On/Off Tranche

Date: 2026-04-04

## Summary

This tranche adds **explicit focus on/off commands** for the persistent tree pane.

Before this tranche, pane focus could be controlled with:
- `/tree-pane-focus`
- `/tree-pane-focus-toggle`

Those worked well, but explicit on/off controls improve operator confidence and align with the growing pattern of both toggle and explicit commands.

## What was added

### New slash commands
- `/tree-pane-focus-on`
- `/tree-pane-focus-off`

### Implementation detail
Added a small helper:
- `handleTreePaneFocusValue(m, enabled bool)`

This centralizes the focus-on/focus-off behavior while letting the existing toggle command remain a thin wrapper.

### Verification added
Added focused regression coverage verifying:
- explicit focus-on enables pane focus
- explicit focus-off disables pane focus

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because explicit commands reduce ambiguity. In complex long-session surfaces, operators often prefer commands that state exactly what they want rather than relying on toggles.

## Design insight

The key principle remains the same:

> explicit on/off controls are just clearer entry points into the same pane-focus state transitions.

No new session or branch semantics are introduced.

## Recommended next move

The strongest next step is now:

> continue selectively adding explicit controls only where they noticeably improve confidence over toggle-only behavior.

## Bottom line

This tranche makes pane focus control more explicit and operator-friendly without complicating the underlying model.
