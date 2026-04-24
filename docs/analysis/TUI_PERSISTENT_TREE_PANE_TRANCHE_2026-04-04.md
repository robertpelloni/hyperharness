# TUI Persistent Tree Pane Tranche

Date: 2026-04-04

## Summary

This tranche adds a **persistent tree pane/view** to the TUI so the canonical tree browser can remain visible while normal prompt interaction continues.

Before this tranche, tree browsing was available through:
- transient browser mode
- selector-based tree switching
- slash-command-driven tree inspection

Those were useful, but inherently modal. Once you exited browser mode, the tree disappeared.

This tranche adds a longer-lived browser view.

## What was added

### 1. Persistent pane state
Added to the TUI model:
- `browserPinned bool`

This is transient UI state indicating that a non-modal tree pane should remain visible while the user continues normal prompt interaction.

### 2. Pane management helpers
Added in `tui/foundation_bridge.go`:
- `pinFoundationTreeBrowser`
- `unpinFoundationTreeBrowser`

`pinFoundationTreeBrowser` refreshes browser items from the canonical runtime and keeps them available for rendering as a persistent pane.

### 3. New TUI slash command
Added in `tui/slash.go`:
- `/tree-pane`

Behavior:
- first invocation pins the tree pane
- second invocation hides it

This is intentionally simple and operator-friendly.

### 4. View integration
Updated `model.View()` so that when:
- `browserActive == false`
- and `browserPinned == true`

the tree browser is rendered as a persistent pane above the normal prompt input area.

That means the user can keep tree context visible while still interacting with the normal TUI prompt.

### 5. Help updated
The TUI help output now documents `/tree-pane`.

## Verification added

Expanded `tui/slash_test.go` to verify:
- `/tree-pane` pins the browser
- pinned browser content appears in the TUI view
- toggling `/tree-pane` again hides it

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche is important because it shifts the browser from a purely transient mode into a more persistent operator aid.

That is a significant usability improvement for long sessions, because the user no longer has to repeatedly reopen the browser to maintain session-tree awareness.

## Design insight

The key architectural principle preserved here is:

> **the persistent pane is still just a view over canonical runtime-derived browser items; it does not own branch/session truth.**

That is exactly right.

The pane improves operator continuity without fragmenting the source of truth.

## What is still missing

This tranche adds persistence, but some richer pane ergonomics remain possible:

1. **auto-refresh of the pinned pane after more session mutations**
2. **independent pane scrolling/navigation**
3. **split-view or side-panel layout work**

## Recommended next move

The strongest next step is now:

> **add refresh and stronger pane ergonomics, or move toward a more dedicated split-view browser layout.**

At this point the TUI tree experience has enough capability that the next frontier is better multi-surface coexistence, not just more commands.

## Bottom line

This tranche makes the tree browser persistent, which is exactly the next right step for long-running interactive coding sessions.
