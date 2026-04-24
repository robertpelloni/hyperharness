# TUI Interactive Tree Selector Tranche

Date: 2026-04-04

## Summary

This tranche extends the TUI tree exploration workflow from:
- plain `/tree` inspection,
- direct `/tree <entryId>` switching,
- and richer child/label surfaces,

into a more interactive selector-style experience built on top of the same canonical foundation runtime.

The key improvement is that users no longer need to manually copy full entry IDs just to switch branches.

## What was added

### 1. Selector state in the TUI model
Added to the TUI model:
- `foundationTreeSelection []string`

This stores the current ordered mapping between selector rows and canonical foundation entry IDs.

Importantly, this is only lightweight UI state.
It does **not** create a competing session or branch model.
The runtime remains the source of truth.

### 2. Numbered tree-selection display
Added in `tui/foundation_bridge.go`:
- `buildFoundationTreeSelectionDisplay`

This renders a numbered selection view that includes:
- active-leaf marker
- numeric index
- entry id
- entry kind
- label (if present)
- preview text/tool name

It also returns the canonical ordered entry ID list used by the TUI for follow-up selection.

### 3. Selector-driven branch switching
Added:
- `switchFoundationTreeSelection`

Behavior:
- validates a requested numeric index
- resolves it through the selector mapping
- delegates the actual switch to the same canonical branch-summary runtime flow already verified earlier

This means the new selector workflow improves ergonomics without introducing new semantics.

### 4. New TUI slash commands
Added in `tui/slash.go`:
- `/tree-select`
- `/tree-go <index> [maxTokens]`

#### `/tree-select`
Shows the numbered selector view and caches the entry order in the TUI model.

#### `/tree-go <index> [maxTokens]`
Resolves the chosen index back to the canonical entry id and performs a real summary-preserving branch switch.

### 5. Help output updated
The TUI `/help` output now documents:
- `/tree-select`
- `/tree-go <index> [maxTokens]`

## Verification added

Expanded `tui/slash_test.go` to verify:
- selector display output is produced
- selector IDs are cached in the TUI model
- `/tree-go` performs a real switch through the canonical runtime
- the active leaf changes after selection-driven switching
- structured summary content still appears in the switch output

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche is important because it meaningfully improves the usability of the TUI tree flow without sacrificing architectural correctness.

Before this change, switching branches required manual entry-ID handling.
After this change, users can work through a selector-style numbered interface while still using the exact same runtime-backed summary-preserving switch semantics.

That is a real UX improvement built on truthful foundations.

## Design insight

The most important design principle preserved here is:

> **improve navigation ergonomics with ephemeral UI state only; keep all branch truth in the canonical runtime.**

That is exactly the right architecture for this codebase.

The TUI is allowed to remember a temporary list of selectable IDs.
It is **not** allowed to become the owner of session/branch semantics.

## What is still missing

This tranche makes the tree UX more interactive, but it is still command-driven rather than fully navigable.

Remaining likely next steps include:

1. **cursor-driven tree browser mode**
   - moving beyond slash-command-based selection into a live navigable view

2. **filter/fold modes**
   - especially if the goal is to approach Pi-style richness over time

3. **branch preview before switch**
   - showing a preview of what will be summarized before the switch is executed

## Recommended next move

The strongest next step is now:

> **begin building a cursor-driven in-TUI tree browser mode using the verified selector/state/runtime building blocks already in place.**

At this point the substrate is strong enough, the ergonomics are improving, and the next obvious upgrade is to move from slash-command selection toward interactive navigation.

## Bottom line

This tranche transforms the TUI tree flow from:
- manual id-driven switching

into:
- selector-driven switching backed by the exact same canonical runtime semantics.

That is precisely the kind of UX-on-top-of-truth progression this project needs.
