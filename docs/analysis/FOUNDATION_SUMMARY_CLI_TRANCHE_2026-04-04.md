# Foundation Summary CLI Tranche

Date: 2026-04-04

## Summary

This tranche surfaces the previously verified native summary-generation workflows through the existing `foundation` CLI command tree.

Before this tranche, the canonical Go foundation already had truthful native APIs for:
- branch-summary preparation
- branch-summary generation
- compaction preparation
- compaction generation
- generate-and-append workflows for both branch and compaction summaries

However, those workflows were still mostly internal/runtime-level primitives.

This tranche makes them accessible through the real CLI surfaces already used by the repository.

## What was added

### 1. New `foundation summary` command group
Added a new top-level `foundation` subcommand group:

- `foundation summary branch`
- `foundation summary compact`

This is the correct place for them because they are:
- foundation-native,
- directly backed by the canonical `foundation/pi` runtime,
- and conceptually about inspecting/exercising the Go foundation port.

### 2. `foundation summary branch`
Added CLI support for branch-summary generation.

Flags:
- `--session`
- `--target`
- `--max-tokens`

Behavior:
- prepares a branch summary using the native budget-aware preparation flow
- generates a summary using the native generation path
- appends the summary to the destination branch via the canonical runtime
- emits JSON describing:
  - preparation
  - summary
  - generated output
  - resulting session

This surfaces a real, truthful foundation capability instead of a placeholder command.

### 3. `foundation summary compact`
Added CLI support for compaction generation.

Flags:
- `--session`
- `--keep-recent-tokens`

Behavior:
- prepares a compaction span using the native compaction preparation flow
- generates a deterministic compaction summary
- appends the generated compaction entry to the session
- emits JSON describing:
  - preparation
  - summary
  - generated output
  - resulting session

This is an important milestone because compaction is no longer merely a low-level API. It is now visible through the real operator-facing foundation CLI.

### 4. Foundation helper payload functions
Added helper request/response pathways in `cmd/foundation_http.go`:
- `prepareFoundationBranchSummary`
- `generateFoundationBranchSummary`
- `prepareFoundationCompaction`
- `generateFoundationCompaction`

This matters because these helpers are reusable beyond Cobra itself. They can also support:
- future HTTP endpoints,
- future TUI slash surfaces,
- and future integration tests.

That is the right design: the CLI commands are thin shells over reusable helpers, not one-off logic.

## Verification added

### Command/helper verification
Expanded `cmd/foundation_http_test.go` to verify:
- branch-summary helper behavior
- compaction helper behavior
- preparation/summarization/session payloads are all returned

### Full validation
Verified successfully:

```bash
go test ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche is important because it turns the summary-generation work from:
- “verified internal substrate”

into:
- “real operator-accessible capability.”

That is a big step in maturity.

The project now has a stronger continuity between:
- internal foundation APIs,
- command-line operator surfaces,
- and future TUI/HTTP surfaces.

In other words:

> **the truthful foundation is no longer hidden behind code-only APIs; it is now exposed in the actual CLI.**

## Design insight

The most important design decision in this tranche is that the new CLI commands are not inventing new behavior.

They are thin wrappers over the already-verified canonical foundation runtime.

That preserves the project’s current best practice:
- build truth in `foundation/pi`
- expose that truth through command surfaces
- avoid parallel fake implementations in the UI layer

This is the correct layering for a project of this complexity.

## What is still missing

This tranche surfaces summary-generation through CLI, but not yet through all operator surfaces.

Remaining high-value gaps include:

1. **TUI slash support for summary workflows**
   - `/tree`-style workflows are still not surfaced through the TUI.

2. **HTTP endpoint exposure**
   - The helper layer is ready for it, but dedicated HTTP routes were not added in this tranche.

3. **Provider-backed generation**
   - Commands currently use the deterministic local generator path.
   - That is correct as a default, but richer generation can later be added via adapters.

## Recommended next move

The strongest next move is now:

> **surface the same truthful branch/compaction generation flows through TUI slash workflows, starting with `/tree`-style behavior or summary-oriented slash commands.**

The canonical foundation is ready for it.

## Bottom line

This tranche promoted summary generation from a verified internal engine to a real CLI-accessible feature.

That is exactly the right direction:
- first build a truthful foundation,
- then expose it through real interfaces,
- then expand the surfaces further.
