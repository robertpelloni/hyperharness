# TUI Tree Explorer Tranche

Date: 2026-04-04

## Summary

This tranche expands the initial TUI `/tree` surface into a richer branch explorer while keeping all semantics anchored to the canonical `foundation/pi` runtime.

Before this tranche, the TUI had:
- `/tree` for textual session/branch inspection
- `/tree <targetEntryId> [maxTokens]` for branch switching with generated summary preservation

That was already truthful, but still minimal.

This tranche extends the explorer with:
- label visibility in tree output
- child-branch inspection
- label-setting from the TUI
- and runtime wrappers for the required foundation session operations

## What was added

### 1. Runtime wrappers for richer tree inspection
Added to `foundation/pi/runtime.go`:
- `GetChildren`
- `GetLabel`
- `GetSessionName`

These are small but important. They let the UI stay thin while depending on the canonical runtime rather than duplicating store logic.

### 2. Richer tree display
Expanded `buildFoundationTreeDisplay` in `tui/foundation_bridge.go`.

It now shows:
- session id
- session name
- active leaf id
- per-entry child counts
- per-entry labels
- preview text/tool name
- active leaf marker

This makes the tree display meaningfully more informative without abandoning the simple, truthful text-first approach.

### 3. Child-branch inspection helper
Added:
- `buildFoundationChildrenDisplay`

This surfaces the direct children of a specific entry, letting the TUI show actual local branch fan-out rather than only a flat list of entries.

### 4. Label-setting helper
Added:
- `setFoundationLabel`

This lets the TUI attach labels/bookmarks to canonical session entries using the verified foundation session model.

### 5. New TUI slash commands
Added in `tui/slash.go`:
- `/tree-children <entryId>`
- `/label <entryId> <label>`

#### `/tree-children`
Shows direct child branches for an entry.

#### `/label`
Sets a label on a specific entry through the canonical runtime.

### 6. Help text updated
The TUI slash help output now documents the richer tree-explorer commands.

## Verification added

Expanded `tui/slash_test.go` to verify:
- label-setting output
- labeled tree output
- tree-children output
- child branch previews appearing in the child explorer view

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche is important because it moves the TUI tree surface from:
- a minimal branch viewer

to:
- a more useful branch explorer.

It is still text-first and intentionally conservative, but it now exposes more of the real structure that matters for branch-oriented session work:
- labels/bookmarks
- child branches
- session naming
- active leaf context

That is a meaningful improvement in operator ergonomics without compromising the architectural discipline of the project.

## Design insight

The key design decision here remains the same as the previous `/tree` tranche:

> **enrich the UI by asking the canonical runtime for more truth, not by inventing a parallel tree model inside the TUI.**

That is exactly the right discipline for this codebase.

The TUI now knows more, but only because the foundation runtime exposes more.

## What is still missing

This is a strong step, but full Pi-style tree UX still remains future work.

Key remaining gaps include:

1. **Interactive branch selector**
   - Current interaction is still command-driven, not an in-place explorer widget.

2. **Filter/fold modes**
   - Pi’s tree UI exposes richer browsing affordances.

3. **Label clearing/editing ergonomics**
   - Current `/label` surface sets labels, but the UX can still become richer.

4. **Better branch path rendering**
   - The current output is still entry-list oriented rather than a fuller visual tree.

## Recommended next move

The strongest next step is now:

> **start building an interactive TUI tree/branch selector on top of the now-richer canonical runtime surfaces.**

That could still begin modestly — for example, a simple explorer mode or selector view — but it would now be building on truthful runtime APIs rather than placeholders.

## Bottom line

This tranche makes the TUI `/tree` surface more useful and more faithful to real branch-oriented work.

It does so without breaking the project’s most important rule:
- the foundation remains the source of truth,
- and the UI remains a view/controller over that truth.
