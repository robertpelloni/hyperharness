# TUI Browser Preview and Filter Tranche

Date: 2026-04-04

## Summary

This tranche improves the new cursor-driven TUI tree browser with two important ergonomics:
- **type-to-filter**
- **preview-before-switch**

These enhancements make the browser more usable while preserving the same architecture discipline as before:
- session and branch truth remain in the canonical runtime
- the browser owns only transient UI state

## What was added

### 1. Browser filter state
Added to the TUI model:
- `browserFilter string`

This lets the browser maintain a temporary filter string while active.

### 2. Filter-aware keyboard handling
While browser mode is active:
- typing runes appends to the filter
- `Space` also contributes to the filter
- `Backspace/Delete` removes the last filter character
- the browser selection index is clamped to the visible filtered result set

This makes the tree browser substantially more usable once session histories get longer.

### 3. Filter helper
Added in `tui/foundation_bridge.go`:
- `filterTreeBrowserItems`

This filters browser items using a lowercase contains-match over:
- kind
- label
- preview text

### 4. Preview panel in browser rendering
Expanded `renderTreeBrowser` so the browser now shows:
- filter state
- match count
- a preview block for the currently selected entry

That means selection is no longer blind. The user can see details about the highlighted entry before switching branches.

### 5. Browser rendering now includes
- active leaf marker
- cursor marker
- item index
- id
- kind
- label
- preview text
- filter string
- filtered match count
- preview pane for the selected item

## Verification added

Expanded the browser-mode test in `tui/slash_test.go` to verify:
- browser view includes a preview section
- type-to-filter updates browser state
- filtered view reflects the filter string
- backspace clears filter state correctly enough to continue navigation
- Enter still performs a real canonical branch switch after filtering interactions

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche is important because it makes the browser meaningfully more usable without increasing architectural risk.

Before:
- the browser was keyboard-driven, but still mostly a flat list with selection

After:
- users can search/filter the list
- see a preview before switching
- and still rely on the exact same verified branch-switching flow underneath

## Design insight

The most important design decision here is:

> **filtering and preview are UI-only enrichments; they do not alter the branch/session semantics layer.**

That is the right architecture.

The UI becomes more sophisticated, but the runtime remains the only owner of branch truth.

## What is still missing

The browser is now much more usable, but further improvements remain possible:

1. **folding/grouping**
   - especially by parent/branch segments

2. **more visual tree layout**
   - current browser is still list-oriented rather than graph-oriented

3. **preview of branch-summary plan before switch**
   - currently the preview is entry-focused, not summary-preparation-focused

## Recommended next move

The strongest next step is now:

> **add branch-preview/fold ergonomics to the browser, possibly including preview of children or a pre-switch summary-preparation view.**

That would make the browser more powerful without departing from the verified canonical runtime.

## Bottom line

This tranche upgrades the cursor-driven tree browser from a bare selector into a more usable explorer.

It now supports search-like narrowing and preview, while still preserving the same truth-first architecture.
