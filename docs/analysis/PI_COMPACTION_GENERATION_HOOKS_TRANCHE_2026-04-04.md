# Pi Compaction Generation Hooks Tranche

Date: 2026-04-04

## Summary

This tranche mirrors the branch-summary generator architecture for **compaction** in the canonical `foundation/pi` path.

Before this tranche, the foundation had:
- explicit branch/leaf semantics,
- common-ancestor computation,
- budget-aware branch-summary preparation,
- a deterministic branch-summary generator,
- and generate-and-append branch-summary APIs.

Compaction still lagged behind that maturity.

This tranche closes that gap by adding:
- compaction preparation structs,
- budget-aware compaction preparation,
- deterministic compaction generation,
- and runtime/session APIs for generate-and-append compaction summaries.

## What was added

### 1. `CompactionPreparation`
Added a new native preparation struct carrying:
- `LeafID`
- `EntriesToSummarize`
- `SerializedConversation`
- `FileOps`
- `EstimatedTokens`
- `TokensBefore`
- `FirstKeptEntryID`
- `KeepRecentTokens`

This gives compaction the same kind of rich preparation substrate that branch summaries already had.

### 2. Budget-aware compaction preparation
Added on `SessionStore`:
- `PrepareCompaction(sessionID)`
- `PrepareCompactionWithBudget(sessionID, keepRecentTokens)`

Behavior:
- reconstructs the active branch
- filters out metadata-only entries from the compaction context
- estimates total pre-compaction context size
- trims the recent tail to the requested keep budget
- identifies `FirstKeptEntryID`
- produces the compacted span as `EntriesToSummarize`
- serializes that span for summary generation
- computes cumulative file operations across the summarized span

This is the compaction analogue of the branch-summary preparation layer.

### 3. `CompactionSummaryGenerator`
Added a separate compaction-generation interface:

```go
type CompactionSummaryGenerator interface {
    GenerateCompactionSummary(ctx context.Context, prep *CompactionPreparation) (string, error)
}
```

That gives compaction the same pluggable generation seam as branch summaries.

### 4. Deterministic compaction generation
`DeterministicSummaryGenerator` now also implements:
- `GenerateCompactionSummary`

It generates a structured compaction summary with:
- Goal
- Constraints & Preferences
- Progress
- Key Decisions
- Next Steps
- Critical Context
- read-files
- modified-files

This is fully local, deterministic, and testable.

### 5. Session-store generation APIs
Added:
- `GenerateCompactionSummary(ctx, prep, generator)`
- `CompactWithGeneratedSummary(ctx, sessionID, keepRecentTokens, generator, details)`

This yields a complete native compaction generation flow:
1. prepare compaction span
2. generate structured summary
3. append a `compaction` entry with summary, `firstKeptEntryId`, and `tokensBefore`

### 6. Runtime wrappers
Added on `foundation/pi/runtime.go`:
- `PrepareCompaction`
- `PrepareCompactionWithBudget`
- `GenerateCompactionSummary`
- `CompactWithGeneratedSummary`

This exposes the compaction-generation workflow through the canonical runtime API.

## Verification added

### New test coverage
Expanded `foundation/pi/summary_test.go` to verify:
- deterministic compaction summary generation
- generated compaction entry append behavior
- `firstKeptEntryId` preservation
- `tokensBefore` population
- `BuildSessionContext()` recognizing compaction use after generated compaction

### Full validation
Verified successfully:

```bash
go test ./foundation/pi/...
go test ./foundation/...
```

Everything remained green.

## Why this matters

This tranche is important because the canonical foundation now has a **unified summary-generation architecture** for both:
- branch summaries
- compaction summaries

That is a major architectural milestone.

The summarization story is no longer one-off or asymmetrical.

Instead, the foundation now has a coherent pattern:
1. prepare context
2. trim to a budget
3. serialize for generation
4. preserve cumulative file ops
5. generate summary through a pluggable generator
6. append the summary entry back into session state

That is exactly the kind of truthful, extensible substrate a long-term harness needs.

## Design insight

The most important design decision here is:

> **Branch summaries and compaction summaries now share the same architectural pattern while remaining distinct in semantics.**

That is better than forcing one mechanism to pretend to be the other.

Branch summaries care about:
- divergence,
- common ancestors,
- abandoned branch context.

Compaction cares about:
- context-window pressure,
- keeping recent context alive,
- summarizing older context spans.

The generator architecture now respects both.

## What is still missing

This tranche is a major step, but a few important follow-ons remain:

1. **Provider-backed compaction generation**
   - The interface is ready.
   - The generation is still deterministic/local by default.

2. **More exact Pi cut-point behavior**
   - Current compaction preparation is truthful and useful.
   - It does not yet fully reproduce every nuance of Pi’s split-turn logic.

3. **Tree/CLI/TUI integration**
   - The summary-generation substrate now exists for both branch and compaction workflows.
   - The operator-facing `/tree` and compaction UX still needs to be surfaced.

## Recommended next move

The strongest next step is now:

> **start surfacing the truthful summary-generation workflows in CLI/TUI behavior, beginning with `/tree`-style branch workflows and then compaction command flows.**

At this point, the foundation substrate is strong enough that UI work would be building on something real rather than speculative.

## Bottom line

This tranche completed the second half of the summarization architecture.

The canonical Go Pi foundation now has:
- budget-aware preparation,
- deterministic generation,
- and generate-and-append workflows

for **both branch summaries and compaction summaries**.

That is a real foundation milestone.
