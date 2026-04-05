# TUI Pre-Switch Summary Preview Tranche

Date: 2026-04-04

## Summary

This tranche adds **pre-switch branch-summary preparation preview** into the cursor-driven TUI tree browser.

Before this tranche, the browser could:
- show a filtered entry list
- show a basic preview of the selected entry
- switch branches through the canonical runtime on Enter

That was already useful, but the next natural ergonomic improvement was:

> show what the branch-summary machinery is likely to do **before** the switch happens.

This tranche does exactly that.

## What was added

### 1. Enriched browser item model
Expanded `TreeBrowserItem` to carry summary-preparation hints:
- `SummaryEntries`
- `CommonAncestorID`
- `ReadFilesCount`
- `ModifiedFilesCount`

These fields are still derived from the canonical runtime, not invented by the UI.

### 2. Browser item generation now consults canonical summary preparation
`buildFoundationTreeBrowser` now, for non-leaf entries:
- calls `PrepareBranchSummaryWithBudget(sessionID, entry.ID, 128)`
- extracts summary preparation facts into the browser item

That means the browser preview is now informed by the same verified branch-summary preparation engine that powers actual switching.

### 3. Preview pane upgraded
The browser preview now includes:
- selected entry id
- kind
- label
- preview text
- branch summary entry count estimate
- common ancestor id (when present)
- read file count
- modified file count

If the selected item is already the active leaf, the preview shows that branch summary is unnecessary.

This is a valuable operator-facing affordance because it helps users understand the likely cost/context impact of switching before they commit to it.

## Verification added

Updated the browser-mode test in `tui/slash_test.go` to assert that the browser preview includes branch-summary preparation details.

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche is important because it improves decision support rather than just navigation.

The browser is no longer only a way to move around the tree. It is becoming a way to understand:
- what branch you are looking at,
- what switching will imply,
- and how much abandoned context will likely be summarized.

That is much closer to the kind of thoughtful branch UX that serious long-session agent work needs.

## Design insight

The key architectural rule remained intact:

> **the browser preview does not guess branch semantics on its own; it asks the canonical runtime for preparation details.**

This is exactly the right pattern.

The browser gains richer guidance while still remaining only a thin UI layer over verified foundation behavior.

## What is still missing

This tranche gives the browser better previews, but some next-level ergonomics still remain:

1. **folding/grouping**
   - especially useful once the tree gets larger

2. **multi-line or richer visual branch layout**
   - current view is still list-oriented

3. **explicit confirm/no-confirm switch mode**
   - now that the preview is richer, a confirmation workflow could make sense

## Recommended next move

The strongest next step is now:

> **add one more layer of browser ergonomics, likely folding/grouping or confirmation-before-switch, on top of the same preview-enhanced runtime-backed browser.**

## Bottom line

This tranche upgrades the browser from:
- entry preview

to:
- entry preview plus branch-summary impact preview.

That is exactly the right kind of ergonomics improvement for a truth-first harness.
