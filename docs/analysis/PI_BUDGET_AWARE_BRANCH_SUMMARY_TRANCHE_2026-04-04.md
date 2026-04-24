# Pi Budget-Aware Branch Summary Tranche

Date: 2026-04-04

## Summary

This tranche built the next layer on top of the native branch-summary groundwork:

- **budget-aware branch-summary preparation**
- **conversation serialization for summarization**
- **cumulative file-operation extraction**
- **a structured summary template hook**

This is the first point where the canonical Go Pi foundation begins to look like a real summarization substrate rather than only a tree-navigation substrate.

## What was added

### 1. Richer `BranchSummaryPreparation`
The preparation object now carries much more than just structural branch information.

Added fields:
- `SerializedConversation`
- `FileOps`
- `EstimatedTokens`
- `MaxTokens`

This matters because future summary generation needs two things at once:
1. a bounded conversation payload to summarize
2. preserved operational context such as files read and modified

### 2. `BranchSummaryFileOps`
Added a dedicated struct for cumulative file-operation tracking:
- `ReadFiles`
- `ModifiedFiles`

This follows the same broad direction as Pi’s documented summary details model.

### 3. `PrepareBranchSummaryWithBudget`
Added a budget-aware preparation API:

```go
PrepareBranchSummaryWithBudget(sessionID, targetID string, maxTokens int)
```

Behavior:
- computes the abandoned branch path as before
- trims the summary conversation window using a token budget
- serializes the selected entries for future summary generation
- computes cumulative file operations from the entire abandoned branch path
- returns estimated token usage for the serialized preparation

This is a very important progression.

The branch-summary pipeline now distinguishes between:
- **what text is small enough to summarize right now**
- and **what operational context should still be preserved cumulatively**

That is a stronger design than only trimming everything blindly.

### 4. Text serialization helpers
Added helpers to turn session entries into summary-safe text:
- `serializeConversation`
- `serializeSingleEntry`
- `flattenToolResultText`
- `truncateSummaryText`
- `estimateSerializedTokens`

The serialization format now captures branch-relevant history in a summary-friendly way using labels such as:
- `[User]: ...`
- `[Assistant tool calls]: ...`
- `[Tool result]: ...`
- `[Branch summary]: ...`
- `[Compaction summary]: ...`

This is directly aligned with Pi’s documented conversation serialization concept.

### 5. Budget trimming helper
Added:
- `trimEntriesToBudget`

Behavior:
- walks the abandoned branch entries from newest to oldest
- keeps entries until the budget is reached
- preserves chronological order in the final output

This gives the native foundation a truthful and testable approximation of summary-window budgeting.

### 6. Cumulative file-op extraction helpers
Added:
- `collectBranchFileOps`
- `collectFileOpsFromEntry`

These inspect:
- direct tool calls such as `read`, `write`, and `edit`
- previous `branch_summary` and `compaction` details

This means the system can preserve read/modified file context even when the summary text itself is budget-trimmed.

That is a strong design choice because it separates:
- text-token limits
- from long-term operational context retention

### 7. Structured summary template hook
Added:
- `DefaultStructuredSummaryTemplate(prep *BranchSummaryPreparation) string`

This produces a summary skeleton with:
- Goal
- Constraints & Preferences
- Progress
- Key Decisions
- Next Steps
- Critical Context
- `<read-files>`
- `<modified-files>`

That does not yet generate the summary itself, but it provides the exact hook shape needed for:
- later LLM generation
- custom summary providers
- or extension interception

### 8. Runtime wrapper
Added runtime-level wrapper:
- `PrepareBranchSummaryWithBudget`

That exposes the new preparation layer through the canonical runtime surface.

## Verification added

### New test coverage
Extended `foundation/pi/session_test.go` to verify:
- budget-aware branch-summary preparation
- serialized conversation generation
- estimated token tracking
- cumulative file-op extraction
- structured summary template content

### Full validation
Verified successfully:

```bash
go test ./foundation/pi/...
go test ./foundation/...
```

Everything remained green.

## Design insight

A crucial design decision in this tranche was:

> **Use budget trimming for the summary text window, but keep file-op tracking cumulative across the whole abandoned branch path.**

This is better than a naive design where file context disappears simply because token budget forced a shorter text window.

That makes the branch-summary substrate both:
- more useful, and
- more robust for future agent reasoning.

## What is still missing

This tranche significantly improves summary preparation, but several things remain for fuller Pi-style summarization parity:

1. **Actual LLM summary generation**
   - The preparation layer is ready.
   - The native foundation does not yet generate structured summaries automatically.

2. **Better token estimation**
   - Current estimation is intentionally simple.
   - Future work can use provider-aware token estimators if needed.

3. **Compaction preparation parity**
   - This tranche focused on branch summarization.
   - Similar budget-aware preparation should eventually be mirrored for compaction.

4. **Full `/tree` UX integration**
   - The substrate is now much stronger.
   - Operator-facing tree workflows still need to be built on top.

## Recommended next move

The next strongest move is now:

> **add native summary-generation hooks and/or provider-backed summary generation for branch summaries and compaction preparation**

That would turn the current preparation substrate into a full summarization workflow.

After that, the system will be in a strong position for:
- truthful `/tree` behavior,
- branch-summary UX,
- and compaction parity expansion.

## Bottom line

This tranche gave the canonical Go Pi foundation its first genuinely useful **summary preparation engine**.

It now knows:
- what branch work is being abandoned,
- how to trim that work into a bounded summary window,
- how to serialize it for summarization,
- and how to preserve cumulative file context beyond the text window.

That is the right substrate for the next stage of real parity work.
