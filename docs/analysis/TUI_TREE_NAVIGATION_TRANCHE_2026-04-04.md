# TUI Tree Navigation Tranche

Date: 2026-04-04

## Summary

This tranche introduces a minimal but truthful `/tree`-style navigation surface into the TUI, backed directly by the canonical `foundation/pi` session and summary runtime.

Before this tranche, the TUI could:
- create/use a canonical foundation session
- append prompts and assistant messages into that session
- trigger native branch-summary generation via `/summary-branch`
- trigger native compaction generation via `/summary-compact`

But there was still no operator-facing tree navigation surface.

This tranche adds the first truthful step in that direction.

## What was added

### 1. Foundation tree display helper
Added in `tui/foundation_bridge.go`:
- `buildFoundationTreeDisplay`

This renders a simple textual tree-oriented session listing showing:
- session id
- active leaf id
- each entry id
- parent id
- entry kind
- a short preview of entry text/tool name
- marker for the active leaf

This is intentionally simple, but truthful.
It directly reflects the real canonical foundation session state.

### 2. Foundation tree switch helper
Added in `tui/foundation_bridge.go`:
- `switchFoundationTreeDisplay`

Behavior:
- resolves the active foundation session
- compares the current leaf to the requested target
- if already on target, emits a stable informational display
- otherwise performs a real branch transition using:
  - `BranchWithGeneratedSummary`
- returns a structured summary display for the switch

This is important because the TUI is no longer simulating a branch switch. It is executing the same canonical runtime flow already verified in the foundation layer.

### 3. TUI `/tree` slash command
Added in `tui/slash.go`:
- `/tree`
- `/tree <targetEntryId> [maxTokens]`

#### `/tree`
Shows the active canonical foundation session tree.

#### `/tree <targetEntryId> [maxTokens]`
Switches to a target entry using the native generated branch-summary flow and displays the resulting summary.

This is the first real `/tree`-style operator flow in the TUI built on the truthful foundation substrate.

### 4. Help output updated
The TUI `/help` text now documents:
- `/tree`
- `/tree <targetEntryId> [maxTokens]`

## Verification added

Expanded `tui/slash_test.go` to verify:
- `/tree` produces a foundation tree display
- `/tree <targetEntryId> [maxTokens]` produces a branch-switch display
- the branch-switch display includes structured summary content
- the active leaf actually changes after the switch

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because it turns the summary-generation and branch semantics work into a visible TUI navigation surface.

Until now, the project had:
- a truthful foundation,
- CLI exposure,
- and TUI summary commands,

but not yet an actual tree-oriented navigation primitive.

This tranche begins to close that gap.

It is not the full Pi `/tree` UX yet, but it is the first truthful bridge from:
- canonical branch semantics
- into operator-facing TUI navigation.

## Design insight

The key design choice here was:

> **Expose a minimal tree surface first, but make it fully truthful.**

That is better than building a flashy UI over unverified behavior.

The current `/tree` implementation is intentionally plain-text and conservative, but it is directly backed by the canonical runtime.

That means future richer tree UIs can be layered on top without having to rewrite the semantics.

## What is still missing

This tranche is a real step forward, but several pieces remain for fuller Pi-style tree parity:

1. **Interactive branch browser**
   - Current `/tree` is command-driven, not an interactive explorer.

2. **Fold/unfold and filtering modes**
   - Pi’s tree UI supports richer navigation and filtering.

3. **Explicit branch labels/bookmarks in TUI UX**
   - The canonical foundation supports labels, but the TUI does not yet present them richly.

4. **User-prompted branch summary choice**
   - Current behavior is deterministic through the explicit command.
   - A richer UX could offer summary/no-summary choice during branch switching.

## Recommended next move

The strongest next step is now:

> **expand `/tree` from a minimal command surface into a richer TUI session/branch explorer while keeping the canonical foundation runtime as the only source of truth.**

That would include:
- richer entry formatting
- label visibility
- branch-child visibility
- perhaps next-step keyboard-driven branch selection

## Bottom line

This tranche gives the TUI its first real `/tree`-style navigation capability.

It is not a cosmetic mock.
It is a truthful operator-facing surface over the verified canonical branch/session/summary engine.

That is exactly the right way to continue the port.
