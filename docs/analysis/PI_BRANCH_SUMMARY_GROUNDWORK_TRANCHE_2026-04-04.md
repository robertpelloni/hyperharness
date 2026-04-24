# Pi Branch Summary Groundwork Tranche

Date: 2026-04-04

## Summary

This tranche built the next logical layer on top of the newly explicit branch/leaf semantics in the canonical `foundation/pi` track:

- **common ancestor computation**
- **branch-summary preparation**
- **native branch-with-summary append behavior**

This does **not** yet implement full Pi `/tree` UX or LLM-generated branch summaries automatically, but it establishes the exact substrate needed for that future work.

## Why this mattered

After the previous tranche, the native Go foundation could:
- track the active leaf explicitly,
- branch to earlier entries,
- append new entries to the active branch,
- and reconstruct active-branch context.

That was enough to represent a tree.

But Pi’s real branch behavior involves something more subtle:

> when moving from one branch to another, the harness needs to understand **what work is being abandoned**, determine the **common ancestor**, and create the substrate for a summary of the abandoned path.

Without that, later `/tree` parity would still be shallow.

## What was added

### 1. `BranchSummaryPreparation`
Added a new native struct:
- `TargetID`
- `OldLeafID`
- `CommonAncestorID`
- `EntriesToSummarize`

This captures the exact information the runtime/UI will need before generating a branch summary.

### 2. `GetCommonAncestor(sessionID, firstLeafID, secondLeafID)`
Added native common-ancestor computation based on the two branch paths.

Behavior:
- reconstruct both branches from root to leaf
- walk forward until divergence
- return the deepest shared node

This matches the conceptual foundation of Pi’s documented branch summarization flow.

### 3. `PrepareBranchSummary(sessionID, targetID)`
Added native preparation logic that:
- discovers the current/old leaf
- computes the common ancestor between old leaf and target
- collects the entries on the abandoned branch path beyond that ancestor
- returns them as `EntriesToSummarize`

This is the main bridge from raw tree structure to summarization workflow.

### 4. `BranchWithSummary(sessionID, targetID, summary, details)`
Added a helper that:
1. prepares branch-summary context
2. branches the session to the target
3. appends a `branch_summary` entry on the new branch

This is a very important semantic move.

It means the native foundation can now preserve abandoned-branch context in a branch-summary entry attached to the destination branch, which is the right structural direction for future `/tree` parity.

### 5. Runtime wrappers
Added runtime-level wrappers in `foundation/pi/runtime.go`:
- `GetCommonAncestor`
- `PrepareBranchSummary`
- `BranchWithSummary`

That exposes the new branch-summary substrate through the canonical runtime surface.

## Verification added

### New branch-summary test coverage
Added a new test in `foundation/pi/session_test.go` that verifies:
- creation of a branching tree
- common-ancestor correctness
- preparation of the abandoned branch path
- exact entries selected for summary
- branch-summary append behavior
- branch-summary attachment to the destination branch
- correct `fromId` preservation

### Full validation
Verified successfully:

```bash
go test ./foundation/pi/...
go test ./foundation/...
```

Everything remained green.

## What this means architecturally

This tranche upgrades the native foundation from:
- branch-aware storage

into:
- **branch-transition-aware storage**.

That is an important distinction.

Now the foundation can reason about:
- where the user was,
- where the user is going,
- what branch path is being left behind,
- what the common ancestor is,
- and what entries should be summarized.

That is the correct substrate for future parity in:
- `/tree`
- branch summaries
- branch-aware context preservation
- future compaction/branch interplay

## What is still missing

This is groundwork, not the final branch-summary system.

Remaining gaps include:

1. **Automatic summary generation**
   - We now compute what must be summarized.
   - We do not yet invoke an LLM to generate the structured branch summary.

2. **Budget-aware branch-summary preparation**
   - Pi includes token-budget-aware preparation.
   - The native foundation currently returns the structurally correct abandoned path, but not yet a token-budget-constrained subset.

3. **File-operation accumulation across summaries**
   - Pi’s branch summaries track cumulative read/modified files.
   - The native foundation can carry `details`, but does not yet compute those cumulative file-operation sets.

4. **CLI/TUI `/tree` flow**
   - The low-level semantics now exist.
   - The operator-facing branch navigation workflow still needs to be built.

## Recommended next move

With explicit leaf semantics and branch-summary preparation now in place, the next strongest move is:

> **implement budget-aware branch-summary preparation and structured summary generation hooks**

That would make the branch-summary substrate not just structurally correct, but operationally useful.

After that, the system will be ready for a truthful `/tree`-style interactive flow in the canonical foundation path.

## Bottom line

This tranche did not add UI flash.
It added a missing layer of **tree intelligence**.

The canonical Go Pi foundation can now:
- understand branch divergence,
- compute common ancestors,
- identify abandoned branch work,
- and attach preserved branch context to a new branch.

That is real progress toward truthful Pi branch behavior.
