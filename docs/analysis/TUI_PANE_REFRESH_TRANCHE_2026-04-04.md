# TUI Pane Refresh Tranche

Date: 2026-04-04

## Summary

This tranche adds an explicit **manual refresh command** for the persistent tree pane.

Before this tranche, the pane already auto-refreshed after normal prompt/session mutations. That was good, but a manual operator-triggered refresh is still valuable as a predictable control.

## What was added

### 1. New slash command
Added:
- `/tree-pane-refresh`

### 2. Behavior
The command:
- requires that the pane is pinned
- re-runs the canonical pane refresh logic
- emits a confirmation message when successful

### 3. Verification added
Added focused regression coverage verifying:
- the view changes after a session mutation + manual refresh
- the refresh confirmation message is emitted
- the refreshed pane reflects the new session content

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because a manual refresh affordance complements the automatic refresh path and gives the operator a simple “sync now” control.

## Design insight

The key principle remains intact:

> manual refresh still just re-reads canonical runtime-backed session state; it does not invent or mutate alternate pane truth.

## Recommended next move

The strongest next step is now:

> continue adding small, explicit operator controls where they reduce uncertainty and friction in long sessions.

## Bottom line

This tranche gives the pane a straightforward manual sync control, improving operator confidence and recovery ergonomics.
