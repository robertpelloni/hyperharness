# TUI Summary Slash Tranche

Date: 2026-04-04

## Summary

This tranche surfaces the previously verified foundation summary-generation workflows through the TUI slash-command layer.

Before this tranche, the canonical foundation already had:
- native branch-summary preparation and generation
- native compaction preparation and generation
- CLI exposure through `foundation summary branch` and `foundation summary compact`

The TUI, however, still had no truthful bridge to those capabilities.

This tranche fixes that by:
- adding a minimal foundation session bridge into the TUI model
- persisting prompt traffic into the canonical foundation session model
- and exposing summary flows through real slash commands

## What was added

### 1. Minimal foundation session bridge in the TUI model
Added to `tui/chat.go`:
- `foundationSessionID` on the TUI model

This is the smallest truthful bridge needed to let the TUI interact with the canonical foundation session substrate.

### 2. Lazy foundation session creation
Added in `tui/foundation_bridge.go`:
- `ensureFoundationSession`

Behavior:
- if the TUI already has a foundation session, reuse it
- otherwise create a new native foundation session for the TUI working directory

This keeps the TUI simple while making the summary workflows real instead of simulated.

### 3. Prompt persistence into the foundation session model
Added helpers:
- `appendFoundationUserText`
- `appendFoundationAssistantText`

Then wired prompt flow so that:
- user prompts get appended to the foundation session
- assistant display responses get appended as assistant messages

This is an important step because the summary workflows now operate on **actual TUI-created session history**, not on a fake or detached history.

### 4. TUI bridge helpers for summary display
Added in `tui/foundation_bridge.go`:
- `buildFoundationCompactionDisplay`
- `buildFoundationBranchSummaryDisplay`
- `parseSummaryArgs`

These are thin, truthful adapters over the already-verified foundation runtime.

### 5. New slash commands
Added in `tui/slash.go`:
- `/fsession`
- `/summary-compact [keepRecentTokens]`
- `/summary-branch <targetEntryId> [maxTokens]`

#### `/fsession`
Shows (or lazily creates) the active foundation session ID for the TUI.

#### `/summary-compact`
Runs the verified native compaction-generation flow against the active foundation session and appends the result to TUI history.

#### `/summary-branch`
Runs the verified native branch-summary generation flow against a target entry in the active foundation session and appends the result to TUI history.

### 6. Help menu updated
The TUI `/help` output now documents the new summary-related slash primitives.

## Verification added

Expanded `tui/slash_test.go` to verify:
- `/summary-compact` produces foundation compaction summary output
- `/summary-branch` produces foundation branch summary output
- the TUI can create/use a foundation session and execute the summary flows on top of it

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche is important because it is the first time the TUI is wired to the real summary-generation substrate rather than just sitting alongside it.

In practical terms, the TUI now has a truthful bridge to:
- session persistence
- compaction summary generation
- branch summary generation

That means future `/tree`-style behavior can now build on real session state and real summary generation rather than hypothetical future infrastructure.

## Design insight

The most important design choice here was:

> **Use a minimal, lazy session bridge rather than inventing a parallel TUI-only session model.**

That keeps the TUI aligned with the canonical foundation instead of fragmenting behavior.

This is exactly the correct approach for a project that is trying to consolidate multiple harness behaviors into one truthful Go substrate.

## What is still missing

This tranche meaningfully improves the TUI, but several next steps remain:

1. **Actual `/tree` navigation UX**
   - The summary generation is now accessible from slash commands.
   - A real tree browser/switcher is not yet built.

2. **Automatic branch summary prompting**
   - The TUI can trigger summaries explicitly.
   - It does not yet guide users through branch transitions the way Pi’s `/tree` flow does.

3. **Compaction command ergonomics**
   - `/summary-compact` exists.
   - A more Pi-like `/compact` surface could be layered on top later.

## Recommended next move

The strongest next step is now:

> **begin implementing a truthful `/tree`-style TUI flow backed directly by the canonical session/branch-summary foundation.**

At this point, the core substrate is strong enough that a tree-navigation UI would be building on real behavior.

## Bottom line

This tranche promoted summary generation from:
- verified foundation internals
- and CLI-only command surfaces

to:
- **actual TUI slash functionality backed by the canonical foundation runtime**.

That is a real step toward an end-to-end truthful Go harness.
