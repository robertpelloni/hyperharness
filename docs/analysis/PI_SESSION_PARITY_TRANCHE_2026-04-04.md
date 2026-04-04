# Pi Session Parity Tranche

Date: 2026-04-04

## Summary

This tranche focused on the next highest-value parity gap after expanding the native Pi tool surface: the **session model**.

The native `foundation/pi` runtime already had a truthful but simplified session store. It could:
- create sessions,
- append basic entries,
- persist them to JSONL,
- list sessions,
- and fork a session.

That was useful, but it was still behind real Pi semantics in important ways. In particular, it lacked richer session entry kinds and the ability to reconstruct session context with metadata like:
- model changes,
- thinking-level changes,
- compaction summaries,
- branch summaries,
- custom extension state,
- custom extension messages,
- labels,
- and session display name overrides.

This tranche upgrades the session layer toward that reality.

## What was implemented

### 1. Versioned session metadata
`foundation/pi/session.go`

The session metadata now tracks:
- `Version`
- `ParentSession`

This moves the native foundation closer to Pi’s documented session-versioned storage model and makes session lineage explicit.

### 2. Richer session entry schema
`SessionEntry` was expanded to support more Pi-style entry kinds and metadata:

Added fields for:
- `ThinkingLevel`
- `Provider`
- `ModelID`
- `Summary`
- `FromID`
- `FirstKeptID`
- `TokensBefore`
- `CustomType`
- `Data`
- `Display`
- `TargetID`
- `Label`
- `Details`

This gives the session layer enough expressive power to represent a much more realistic subset of Pi session events.

### 3. New append helpers for real session semantics
Added native helpers on `SessionStore` for:
- `AppendThinkingLevelChange`
- `AppendModelChange`
- `AppendCompaction`
- `AppendBranchSummary`
- `AppendCustomEntry`
- `AppendCustomMessage`
- `AppendSessionInfo`
- `AppendLabelChange`

This matters because parity is not just about tool names. It is also about preserving enough session structure to support:
- context recovery,
- branch movement,
- compaction semantics,
- and extension state.

### 4. Better parent-chain behavior on append
`AppendEntry()` now assigns a parent automatically when omitted, using the latest entry in the session.

That is a small but important upgrade toward Pi’s tree-oriented session model.

### 5. Session inspection helpers
Added new helpers on `SessionStore`:
- `GetEntry`
- `GetChildren`
- `GetBranch`
- `GetLabel`
- `GetSessionName`

These are core building blocks for later work on:
- tree navigation,
- branch UI,
- bookmarks/labels,
- and session-focused TUI workflows.

### 6. Session context reconstruction
Added:
- `BuildSessionContext(sessionID, leafID)`
- `SessionContext`

This is a major step.

The native foundation can now reconstruct a richer, session-derived context summary including:
- current provider/model inferred from `model_change` entries,
- current thinking level inferred from `thinking_level_change` entries,
- display name inferred from `session_info`,
- branch path derived from parent links,
- compaction presence and summary entries,
- exclusion of metadata-only entries (`label`, `session_info`, `custom`) from context,
- inclusion of `custom_message` entries in context.

This does **not** yet reproduce every nuance of Pi’s full `buildSessionContext()` logic, but it moves the session layer from basic persistence toward genuine model-context reconstruction.

### 7. Runtime convenience wrappers
Added corresponding runtime methods in `foundation/pi/runtime.go`:
- `AppendThinkingLevelChange`
- `AppendModelChange`
- `AppendCompaction`
- `AppendBranchSummary`
- `AppendSessionInfo`
- `AppendLabelChange`
- `BuildSessionContext`

That makes the richer session semantics available through the runtime surface instead of only the store.

## Verification

### Tests kept green
Verified successfully:

```bash
go test ./foundation/pi/...
go test ./foundation/...
```

### New/expanded behavior tests
`foundation/pi/session_test.go` now verifies:
- versioned session creation
- parent session preservation on fork
- model change persistence
- thinking-level change persistence
- session-info override behavior
- labels
- custom entries
- custom messages
- compaction entries
- child lookup
- branch reconstruction
- session context reconstruction

## Architectural significance

This tranche matters because the harness is not just a bag of tools.

A serious coding harness also needs a truthful session model. Without it, the project cannot honestly claim deep parity with Pi-style workflows, because so much of the user experience depends on session semantics:
- continuing work,
- switching branches,
- preserving abandoned context,
- changing models mid-stream,
- capturing compactions,
- and remembering bookmarks.

In other words:

> **Tool parity without session parity is only partial harness parity.**

This tranche closes part of that gap.

## What is still missing for fuller Pi session parity

The native session layer is now meaningfully stronger, but it still does not fully match Pi’s documented model. Remaining gaps include:

1. **Exact header/entry wire format parity**
   - Pi uses a `type`-rich event model and a more specific v2/v3 session format.
   - The native foundation is still compatible in spirit, but not identical in every serialized field.

2. **Explicit leaf tracking and branch mutation semantics**
   - The current implementation reconstructs branches from parent links.
   - It does not yet expose a full leaf-management API equivalent to Pi’s runtime session manager behavior.

3. **More exact compaction semantics**
   - Pi’s compaction logic uses a more structured summary workflow with precise cut-point handling, split turns, and previous-summary reuse.
   - The current native layer stores compaction entries and uses `firstKeptEntryId`, but it does not yet reproduce the complete algorithm.

4. **Full branch summarization workflow**
   - Entries can now represent branch summaries, but automatic branch-summary generation and common-ancestor-aware tree switching remain future work.

5. **Custom entry/message parity**
   - The native layer can now persist custom entries and custom messages, but it does not yet implement a full extension event lifecycle equivalent to Pi.

## Recommended next step after this tranche

The next best move is still the same broad sequence, but now one level deeper:

### Immediate next recommendation
1. Keep `foundation/*` green.
2. Continue Pi-native parity in one of these two directions:
   - **A. Session/tree semantics**: explicit leaf movement, branch switching, branch summary generation
   - **B. Remaining harness-mode parity**: print/json/rpc/session-oriented CLI behavior built on the native foundation

### My recommendation
The strongest next move is:

> **implement explicit branch/leaf management and branch-summary semantics in `foundation/pi`**

Why?
Because now that the session store can hold richer entry types and build a useful context, the next natural step is to make the tree itself more truthful and navigable.

That will unlock later work on:
- a Pi-style `/tree` flow,
- session branch summaries,
- and more faithful compaction + branch context behavior.

## Bottom line

This tranche significantly improved the native Go foundation in a way that is both:
- technically real, and
- strategically correct.

It did not add cosmetic parity.
It added **state-model parity**.

That is exactly the kind of progress this project needs.
