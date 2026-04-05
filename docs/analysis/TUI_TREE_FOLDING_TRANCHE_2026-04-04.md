# TUI Tree Folding Tranche

Date: 2026-04-04

## Summary

This tranche adds **folding/collapse of subtrees** to the cursor-driven TUI tree browser.

Before this tranche, the browser already had:
- cursor navigation
- selector-driven switching
- type-to-filter
- entry preview
- branch-summary preparation preview
- confirm-before-switch
- structural depth and child-count cues

That was already a strong browsing surface, but large trees could still become visually dense.

This tranche addresses that by allowing subtree collapse/expand directly in the browser.

## What was added

### 1. Browser collapse state
Added transient browser UI state:
- `browserCollapsed map[string]bool`

This tracks which entry IDs are currently folded in the browser.

As with other browser state, this is UI-local only. It does not alter canonical session truth.

### 2. Fold-aware visible-item computation
Added in `tui/foundation_bridge.go`:
- `visibleTreeBrowserItems`

Behavior:
- starts from the current item list
- applies text filtering
- then hides descendants of collapsed entries until depth unwinds

This is the core mechanic that makes folding work in a list-oriented browser.

### 3. Keyboard controls for collapse/expand
In browser mode:
- `Left` collapses the selected item if it has children
- `Right` expands the selected item if it was collapsed

This keeps the interaction lightweight and keyboard-native.

### 4. Collapse/expand rendering cues
The browser rendering now shows:
- `[+]` for collapsed items with children
- `[-]` for expanded items with children

These cues appear alongside the existing depth/indentation layout so the structural state is visible.

### 5. Selection clamping across fold state changes
When collapsing a subtree changes the number of visible items, the browser clamps the selected index so it always remains valid.

That keeps the browser behavior stable and predictable.

## Verification added

Expanded the browser-mode test to verify:
- `Left` produces a collapsed cue (`[+]`)
- `Right` restores an expanded cue (`[-]`)
- folding interacts correctly with the existing browser flow

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche is important because it improves **navigability at scale**.

The browser is no longer just showing structure; it now gives the user a way to reduce visible complexity when branch trees get large.

That is exactly the right next step after adding depth, child counts, preview, and confirmation.

## Design insight

The key architectural rule remained intact:

> **folding is purely a UI projection over canonical runtime-derived structure.**

That is critical.

The browser becomes more usable, but the runtime remains the only source of truth about the actual session tree.

## What is still missing

This tranche significantly improves scale ergonomics, but a few stronger browser features remain possible:

1. **more graph-like tree rendering**
   - current browser is still line-oriented with structural cues, not a full tree visualizer

2. **richer branch grouping**
   - grouping by branch segments rather than only depth/child count

3. **separate pre-switch summary preview mode**
   - the preview is already strong, but could become more structured

## Recommended next move

The strongest next step is now:

> **build richer grouping/graph-style rendering on top of the now-filterable, foldable, confirmable browser.**

At this point, the browser interaction model is strong enough that layout clarity is the next UX frontier.

## Bottom line

This tranche gives the browser real subtree folding.

That makes the TUI tree explorer significantly more practical for larger session histories while preserving the same truth-first architecture.
