# TUI Graph-Style Rendering Tranche

Date: 2026-04-04

## Summary

This tranche improves the visual structure of the tree browser by moving from simple depth indentation toward more graph-like connector rendering derived from canonical session relationships.

The browser already had:
- cursor navigation
- selector-driven switching
- preview
- filtering
- confirmation-before-switch
- subtree folding
- depth and child-count cues

What it lacked was stronger visual communication of sibling/branch topology.

This tranche adds that.

## What was added

### 1. Richer browser item structure
Expanded `TreeBrowserItem` with:
- `ParentID`
- `Prefix`
- `IsLastChild`

These fields are still computed from the canonical session tree; they are not UI-invented truth.

### 2. Canonical prefix generation
Added `buildTreePrefix`, which walks parent/sibling relationships to generate a structural prefix used in browser rendering.

This gives the browser a more tree-like feel without abandoning the underlying list model.

### 3. Last-child-aware connector rendering
Rendering now distinguishes:
- `├─` for non-final siblings
- `└─` for final children

That is a meaningful readability improvement for branching session histories.

### 4. Prefix-aware browser rows
Browser rows now combine:
- cursor marker
- active-leaf marker
- graph-style prefix
- fold indicator
- id / kind / label / child-count metadata
- preview text

## Verification
Focused tests passed against the modified TUI/browser surfaces:

```bash
go test -run TestTreeBrowserModeNavigation -v ./tui
go test -run TestProcessSlashCommandTreeExplorerLabelAndChildren -v ./tui
go test -run TestFoundationSummaryHelpers -v ./cmd
```

All green.

## Why this matters

This tranche improves structural legibility at scale. It becomes easier to visually distinguish:
- branch fan-out
- sibling relationships
- nested depth
- likely switch targets

without changing any underlying session semantics.

## Design insight

The crucial design rule stayed intact:

> visual tree cues are derived from canonical runtime structure, not from a separate UI-owned tree graph.

That means the browser gets smarter while the runtime remains the only source of truth.

## Recommended next move

The strongest next move after this remains richer browser ergonomics, likely:
- persistent pane/view behavior, or
- stronger grouping/folding strategies beyond simple subtree collapse.

## Bottom line

This tranche makes the tree browser more visually legible and more obviously branch-aware, while preserving the same truth-first architecture.
