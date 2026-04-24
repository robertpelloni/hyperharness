# Pi Summary Generation Hooks Tranche

Date: 2026-04-04

## Summary

This tranche converts the native branch-summary preparation substrate into an actual **summary-generation workflow**.

Until now, the canonical Go Pi foundation could:
- compute branch divergence,
- prepare abandoned-branch entries,
- trim the summary text window to a budget,
- serialize the branch for summarization,
- track cumulative file operations,
- and append branch-summary entries if a caller supplied summary text.

That was strong groundwork, but it still required an external caller to invent the summary itself.

This tranche adds:
- a pluggable summary-generation interface,
- a deterministic built-in generator,
- and a full runtime/session path for **prepare → generate → append**.

## What was added

### 1. `SummaryGenerator` interface
Added a small pluggable interface in `foundation/pi/summary.go`:

```go
type SummaryGenerator interface {
    GenerateBranchSummary(ctx context.Context, prep *BranchSummaryPreparation) (string, error)
}
```

This is intentionally minimal.

It is the correct seam for future evolution because it allows:
- deterministic local generation,
- provider-backed generation,
- HyperCode-routed generation,
- or extension-injected generation,

without changing the rest of the session and runtime APIs.

### 2. `SummaryGeneratorFunc`
Added a function adapter so simple closures can act as generators.

This will make future experimentation easier.

### 3. `DeterministicSummaryGenerator`
Added the default built-in generator.

This generator does **not** call an LLM.
Instead, it turns the branch-summary preparation into a predictable structured summary.

It includes sections for:
- Goal
- Constraints & Preferences
- Progress
- Key Decisions
- Next Steps
- Critical Context
- read-files
- modified-files

This is a very important design choice.

It means the canonical foundation now has a **fully local, deterministic, testable summary path** that does not depend on network access or provider availability.

That is exactly the kind of truthful default a foundation layer should have.

### 4. Structured summary inference helpers
Added helper functions that extract structured information from a preparation object:
- `inferGoal`
- `inferCompletedItems`
- `inferNextSteps`
- `inferCriticalContext`

These make the deterministic generator more useful than a static template dump.

### 5. Session-store generation API
Added on `SessionStore`:
- `GenerateBranchSummary(ctx, prep, generator)`
- `BranchWithGeneratedSummary(ctx, sessionID, targetID, maxTokens, generator, details)`

Behavior:
1. prepare the abandoned branch with budget
2. generate a summary using the provided generator or the deterministic fallback
3. branch to the destination
4. append the generated branch summary entry

That is the first complete native branch-summary generation flow in the canonical Go foundation.

### 6. Runtime wrappers
Added on `foundation/pi/runtime.go`:
- `GenerateBranchSummary`
- `BranchWithGeneratedSummary`

This exposes the new generation workflow through the canonical runtime surface.

## Verification added

### New tests
Added `foundation/pi/summary_test.go` with verification for:
- deterministic structured-summary generation
- generated summary attachment to the destination branch
- preservation of `fromId`
- structured summary content presence

### Full validation
Verified successfully:

```bash
go test ./foundation/pi/...
go test ./foundation/...
```

Everything remained green.

## Why this matters

This tranche is important because it upgrades the native foundation from:
- “summary preparation exists”

to:
- “summary generation is now a real first-class part of the foundation API.”

That is a meaningful milestone.

It means the system can now support future higher-level workflows such as:
- truthful `/tree` branch transitions,
- deterministic offline branch summaries,
- provider-backed branch summaries,
- extension-supplied branch summaries,
- and eventually compaction generation via the same seam.

## Design insight

The most important design decision in this tranche is:

> **The default generator is deterministic and local, while the API remains open for stronger provider-backed generation later.**

That is the right layering:
- default behavior is reliable and testable
- future behavior can be more intelligent without destabilizing the foundation

## What is still missing

This tranche is a major step, but several obvious follow-ons remain:

1. **Provider-backed summary generation**
   - The interface is ready.
   - A HyperCode/provider adapter implementation does not yet exist.

2. **Compaction generation hooks**
   - The new interface pattern should be mirrored for compaction preparation.

3. **More exact Pi-style summary content behavior**
   - The deterministic generator is structurally useful and truthful.
   - It is not yet trying to emulate a real LLM-generated Pi summary verbatim.

4. **UI integration**
   - The generation path exists.
   - A `/tree` or branch-navigation UI using this path is not yet built.

## Recommended next move

The next strongest move is now:

> **apply the same generator/hook architecture to compaction preparation, then start surfacing this summary workflow in CLI/TUI branch navigation.**

That would unify the summarization story across:
- branch summaries
- compaction summaries
- and future extension/provider-driven summary strategies

## Bottom line

This tranche turned the branch-summary pipeline into a real, usable generation system.

The canonical Go Pi foundation can now:
- prepare branch summaries,
- generate structured summaries locally,
- and append them to the correct destination branch.

That is a real parity-enabling milestone.
