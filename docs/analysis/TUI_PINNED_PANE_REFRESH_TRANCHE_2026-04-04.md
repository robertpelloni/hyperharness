# TUI Pinned Pane Refresh Tranche

Date: 2026-04-04

## Summary

This tranche improves the persistent tree pane by making it **auto-refresh after session mutations**.

Before this tranche, the pane could stay visible while normal prompt interaction continued, but its contents would only update when explicitly reopened or re-rendered through other browser flows.

That meant the pane was persistent, but not fully live.

This tranche fixes that.

## What was added

### 1. Refresh helper for pinned browser panes
Added in `tui/foundation_bridge.go`:
- `refreshPinnedFoundationTreeBrowser`

Behavior:
- no-ops if the pane is not pinned
- ensures the canonical session exists
- rebuilds browser items from the canonical runtime
- clamps the selected index against the refreshed visible result set

### 2. Prompt-flow refresh hooks
Integrated pinned-pane refresh into the normal TUI message flow:
- after appending user text into the canonical foundation session
- after appending assistant text into the canonical foundation session

This means the pinned tree pane now updates as the session evolves during normal use.

### 3. Focused regression coverage
Added test coverage to verify:
- a pinned tree pane exists
- a new session mutation changes the view after refresh
- the refreshed pane reflects the new session content

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because it turns the persistent pane from:
- static snapshot held open

into:
- live context companion for the active session.

That is exactly what a long-running session browser pane should become.

## Design insight

The key principle remained intact:

> the pane refreshes by re-reading canonical runtime state, not by mutating or guessing session structure locally.

That keeps the pane truthful while making it much more useful.

## Recommended next move

The strongest next step is now:

> **move from simple pinned refresh toward stronger pane ergonomics or split-view behavior**, since the pane is now not only persistent but also live.

## Bottom line

This tranche makes the persistent tree pane feel like a real companion surface instead of a static overlay.
