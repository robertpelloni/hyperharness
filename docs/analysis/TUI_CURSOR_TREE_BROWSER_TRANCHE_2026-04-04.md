# TUI Cursor Tree Browser Tranche

Date: 2026-04-04

## Summary

This tranche takes the selector-style tree workflow one step further by adding a **cursor-driven browser mode** to the TUI.

Before this tranche, the TUI had:
- `/tree`
- `/tree-children`
- `/label`
- `/tree-select`
- `/tree-go <index> [maxTokens]`

That was already a meaningful ergonomic improvement over raw entry IDs, but it still required explicit slash-command selection.

This tranche adds a more interactive browser-style flow with keyboard navigation.

## What was added

### 1. Browser state in the TUI model
Added to the TUI model:
- `browserActive bool`
- `browserItems []TreeBrowserItem`
- `browserIndex int`

This is still intentionally lightweight UI state. It does not replace the canonical session/runtime model.

### 2. `TreeBrowserItem`
Added a dedicated browser item struct in `tui/foundation_bridge.go` with:
- entry id
- kind
- label
- preview text
- active-leaf marker

This gives the browser mode a cleaner representation than reusing raw strings.

### 3. Browser data builder
Added:
- `buildFoundationTreeBrowser`

This loads canonical session entries through the foundation runtime and turns them into browser items.

### 4. Browser renderer
Added:
- `renderTreeBrowser`

This provides a textual browser view showing:
- cursor position
- active leaf marker
- numeric index
- entry id
- entry kind
- label
- preview text
- usage hints (`↑/↓`, `Enter`, `Esc`)

### 5. Browser selection opener
Added:
- `openSelectedTreeBrowser`

This resolves the currently selected browser item into a real canonical branch switch using the existing summary-preserving runtime flow.

### 6. New slash command
Added:
- `/tree-browser`

Behavior:
- opens the cursor-driven tree browser for the active foundation session
- loads canonical items into browser state
- activates keyboard navigation mode

### 7. Browser-aware keyboard handling
Updated `tui/chat.go` so that when browser mode is active:
- `Up` moves selection upward
- `Down` moves selection downward
- `Enter` switches branches through the canonical runtime and closes the browser
- `Esc` closes the browser instead of quitting the whole app

This is the first true keyboard-driven tree navigation mode in the TUI.

### 8. View integration
Updated `model.View()` so the TUI now renders the browser interface when browser mode is active rather than the normal input prompt.

## Verification added

Expanded `tui/slash_test.go` to verify:
- `/tree-browser` opens browser mode
- browser item list is populated
- `Down` changes the browser index
- `Enter` performs a real canonical branch switch
- browser mode closes after selection
- branch-switch output is appended to history
- the active leaf actually changes

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche is important because it is the first time the TUI tree explorer moves from:
- command-driven navigation

to:
- keyboard-driven interactive browsing.

That is a real UX step toward a richer Pi-like tree experience.

And critically, it was done without violating the project’s main architectural discipline:
- the browser owns only transient UI selection state
- the canonical runtime still owns branch/session truth

## Design insight

The most important design principle preserved here is:

> **interactive browser state is allowed, but only as a transient UI view over the canonical runtime model.**

That is exactly the right way to keep the port honest.

The TUI becomes more powerful, but the semantics still come from the same verified `foundation/pi` runtime.

## What is still missing

This tranche creates a true cursor-driven browser mode, but a fuller tree UX still remains possible.

Natural next steps include:

1. **better visual tree layout**
   - current browser is still entry-list oriented rather than a more graph-like tree presentation

2. **filtering/folding**
   - especially useful for larger sessions

3. **preview-before-switch UX**
   - showing the target branch path or summary preparation before committing a switch

4. **branch-summary opt-in/opt-out flow**
   - current browser Enter uses the summary-preserving switch path directly
   - richer UX could allow explicit summary/no-summary choice

## Recommended next move

The strongest next step is now:

> **add richer browser ergonomics — previews, filtering, and possibly a more structured tree rendering — while keeping the same canonical runtime as the only source of truth.**

At this point, the project has moved from speculative future tree ideas to a real, test-backed, keyboard-driven tree browser substrate.

## Bottom line

This tranche delivers the first true interactive tree browser mode in the TUI.

It is not just a pretty wrapper over static data.
It drives real branch switching through the canonical runtime and keeps the entire semantics chain truthful.
