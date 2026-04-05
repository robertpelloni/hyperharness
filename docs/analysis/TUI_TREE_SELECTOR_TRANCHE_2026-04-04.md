# TUI Tree Selector Tranche

Date: 2026-04-04

## Summary

This tranche extends the TUI tree explorer with a selector-style workflow so users can switch branches by **numbered selection** instead of manually copying entry IDs.

Before this tranche, the TUI already had:
- `/tree`
- `/tree <targetEntryId> [maxTokens]`
- `/tree-children <entryId>`
- `/label <entryId> <label>`

That was truthful and useful, but still awkward in practice because switching branches required manually typing entry IDs.

This tranche adds the next layer of ergonomics while keeping the UI fully grounded in the canonical runtime.

## What was added

### 1. Selector state in the TUI model
Added to `tui/chat.go`:
- `foundationTreeSelection []string`

This stores the current numbered mapping from selector display rows to canonical foundation entry IDs.

This is intentionally minimal state.
It does not replace the foundation runtime’s truth. It only caches the current selector mapping for the TUI layer.

### 2. Selector builder helper
Added in `tui/foundation_bridge.go`:
- `buildFoundationTreeSelectionDisplay`

Behavior:
- loads the canonical session
- reads the active leaf
- renders a numbered list of session entries
- preserves the entry IDs in returned selection order

The display includes:
- an active-leaf marker
- entry index
- entry id
- entry kind
- optional label
- preview text/tool name

### 3. Selector-driven branch switching helper
Added:
- `switchFoundationTreeSelection`

Behavior:
- validates the requested numeric index
- resolves it to the cached canonical entry ID
- delegates to the same existing runtime-backed tree switch flow

That means branch switching still uses the same truthful summary-preserving runtime semantics as before.

### 4. New TUI slash commands
Added in `tui/slash.go`:
- `/tree-select`
- `/tree-go <index> [maxTokens]`

#### `/tree-select`
Shows a numbered entry list and stores the current selector mapping inside the TUI model.

#### `/tree-go <index> [maxTokens]`
Uses the stored selector mapping to switch to the selected entry via the canonical branch-summary runtime flow.

This is an important usability upgrade because it lowers the friction of branch switching without changing the underlying truth model.

### 5. Help output updated
The TUI `/help` output now documents the selector-style tree commands.

## Verification added

Expanded `tui/slash_test.go` to verify:
- `/tree-select` produces selector output
- selector IDs are stored in the TUI model
- `/tree-go` performs a real branch switch using the selector mapping
- the active leaf actually changes
- the generated switch output still contains structured summary content

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche is important because it is the first meaningful step from:
- command-driven tree inspection

toward:
- a more navigable branch-selection UX.

It still keeps the UI simple, but it improves one of the biggest ergonomic pain points: manually typing entry IDs.

That makes the TUI tree flow more practical while still preserving the core architectural discipline of the project.

## Design insight

The most important design decision here was:

> **store only a lightweight selector mapping in the TUI, while leaving branch truth entirely in the canonical runtime.**

That is the correct compromise.

The UI becomes more ergonomic without starting to own a competing representation of the session tree.

## What is still missing

This tranche is a strong usability improvement, but full tree UX still remains future work.

Remaining gaps include:

1. **True interactive selection widget**
   - Current selector is command-based, not keyboard-driven in-place navigation.

2. **Filtering/folding**
   - Pi’s tree UI supports richer exploration modes.

3. **Better child/branch visualization**
   - The selector is entry-list oriented, not yet a full tree graph/explorer.

## Recommended next move

The strongest next step is now:

> **begin building a more interactive in-TUI tree selector/browser using the now-rich canonical tree APIs and the lightweight selector state model.**

That could include:
- a navigable selector view
- direct child browsing
- label-aware sorting/filtering
- branch previews before switching

## Bottom line

This tranche makes the TUI tree flow substantially more usable without compromising the underlying architecture.

The selector is still thin.
The runtime is still the source of truth.
That is exactly how this project should keep evolving.
