# TUI Graph-Style Rendering Tranche

Date: 2026-04-04

## Summary

This tranche deepens the tree browser’s visual structure by moving from plain indentation toward **more graph-like rendering cues**.

Before this tranche, the browser already had:
- depth/indentation
- child-count cues
- filtering
- preview
- pre-switch summary preparation preview
- confirmation-before-switch
- subtree folding

That made the browser useful, but the visual representation of sibling/branch structure was still fairly coarse.

This tranche improves that by surfacing:
- parent linkage in the browser item model
- last-child awareness
- connector-like rendering cues (`├─`, `└─`)
- branch prefixes derived from canonical parent relationships

## What was added

### 1. Richer browser item structure
Expanded `TreeBrowserItem` with:
- `ParentID`
- `Prefix`
- `IsLastChild`

This data is still derived from the canonical session tree, not UI invention.

### 2. Canonical prefix generation
Added a helper:
- `buildTreePrefix`

This walks parent relationships and sibling ordering to produce a structural prefix for each item.

That means the visual grouping now reflects the actual runtime-derived tree relationships more closely than simple depth indentation alone.

### 3. Last-child aware rendering
Rendering now distinguishes between:
- `├─` for non-final siblings
- `└─` for last children

This is a meaningful readability improvement because it helps the eye follow branching structure rather than just entry order.

### 4. Prefix-aware browser output
The browser rows now combine:
- cursor marker
- active-leaf marker
- graph-style prefix
- fold indicator
- id/kind/label/child-count
- preview text

This makes the browser closer to a real branch explorer than a simple annotated list.

## Verification added

Expanded the browser-mode test so that the rendered browser view must now contain graph-style connector glyphs (`├─` or `└─`).

### Full validation
Verified successfully:

```bash
go test -v ./tui
go test -v ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because the browser is now better at conveying **branch topology**, not just node metadata.

That is an important distinction for long-session agent work.

A branch browser should not just say “these entries exist.” It should help the user perceive:
- what is a sibling
- what is nested beneath what
- where branching fans out
- and where a selected item sits in that structure

This tranche improves that without abandoning the truth-first architecture.

## Design insight

The key design principle preserved here is:

> **graph-style rendering is still derived from canonical runtime relationships, not from a separate UI-owned tree model.**

This is critical.

The browser gets visually smarter, but the tree semantics still come from the same `foundation/pi` session truth.

## What is still missing

This tranche improves rendering, but some advanced structure-oriented ergonomics remain possible:

1. **richer grouping modes**
   - e.g. grouping by branch heads or recent divergence points

2. **more explicit branch summaries in the browser list**
   - not just in preview

3. **persistent browser panes/views**
   - rather than a transient mode only

## Recommended next move

The strongest next step is now:

> **add richer grouping or persistent browser panes on top of the now graph-like, filterable, foldable, confirmable browser.**

At this point the browser has moved from a raw list into a true structural explorer.

## Bottom line

This tranche improves visual tree clarity in a real way.

The browser is now not only interactive and safe, but also more structurally legible — which is exactly what a branch-heavy session explorer needs.
