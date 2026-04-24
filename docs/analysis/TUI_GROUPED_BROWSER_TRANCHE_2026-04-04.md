# TUI Grouped Browser Tranche

Date: 2026-04-04

## Summary

This tranche adds a lightweight **grouped view mode** to the TUI tree browser.

Before this tranche, the browser already had:
- cursor navigation
- filtering
- preview
- pre-switch summary preparation preview
- confirmation-before-switch
- subtree folding
- graph-style connector cues

That was already a capable browser, but the visible list could still feel flat in large histories.

This tranche introduces a simple grouping layer based on canonical branch lineage.

## What was added

### 1. Browser grouped-mode state
Added to the TUI model:
- `browserGrouped bool`

This is transient UI-only state that toggles grouping mode in the browser.

### 2. Group key derivation
Expanded `TreeBrowserItem` with:
- `GroupKey`

Added helper:
- `buildGroupKey`

This derives a stable branch-group label from canonical parent lineage, so entries can be clustered by top-level branch ancestry rather than only rendered as a flat sequence.

### 3. Grouped rendering in browser view
The browser renderer now:
- shows whether grouping is enabled
- inserts `[Group] ...` section headings when group boundaries change
- keeps all previously-added browser ergonomics intact

This means grouped view is additive, not a replacement.

### 4. Keyboard toggle for grouping
In browser mode:
- `Tab` toggles grouping on and off

That keeps the interaction fast and local without adding more slash-command complexity.

### 5. Help/usage semantics preserved
Grouped mode does not change switch semantics at all. It is purely an alternate projection of the same canonical runtime structure.

## Verification added

Expanded the browser-mode test so it now verifies:
- `Tab` toggles grouped mode on
- grouped browser rendering includes `[Group]` headings
- all existing browser behavior still works after grouping is enabled

### Full focused validation
Verified successfully with targeted tests:

```bash
go test -run TestTreeBrowserModeNavigation -v ./tui
go test -run TestProcessSlashCommandTreeExplorerLabelAndChildren -v ./tui
go test -run TestFoundationSummaryHelpers -v ./cmd
```

All green.

## Why this matters

This tranche improves **organizational clarity** for larger trees.

The browser is now not just:
- filtered,
- foldable,
- previewable,

but also able to cluster entries by branch lineage in a lightweight way.

That makes the structure easier to mentally scan without changing any underlying semantics.

## Design insight

The critical architectural rule remained intact:

> grouped mode is a presentation derived from canonical parent relationships, not a separate branch model.

That is exactly right.

The browser gets more expressive, but still does not own session truth.

## What is still missing

This tranche makes the browser more organized, but further large-tree ergonomics remain possible:

1. **persistent pane/view state**
   - keeping the browser open alongside other TUI surfaces

2. **stronger branch grouping strategies**
   - e.g. by divergence points or branch heads rather than a simple lineage key

3. **explicit branch-summary preview panes**
   - the summary preview is already informative, but could evolve further

## Recommended next move

The strongest next step is now:

> **consider a more persistent browser pane/view or stronger grouping around branch heads/divergence points.**

At this point the browser is already highly capable, so the next improvements should focus on long-session usability rather than basic interaction mechanics.

## Bottom line

This tranche adds a simple but useful grouped view mode that helps the tree browser scale better with complex branching session histories, while keeping the canonical runtime as the only source of truth.
