# TUI Pane Preset Extension Tranche

Date: 2026-04-04

## Summary

This tranche extends the existing pane preset system with more workflow-shaped modes.

Before this tranche, the pane presets supported:
- `compact`
- `detailed`

Those were useful, but still a small preset vocabulary.

This tranche adds:
- `navigation`
- `review`

so the persistent tree pane can better match different operator intents.

## What was added

### 1. Extended preset command vocabulary
`/tree-pane-preset` now accepts:
- `compact`
- `detailed`
- `navigation`
- `review`

### 2. `navigation` preset
Applies:
- `browserPaneHeight = 10`
- `browserPanePreview = false`
- `browserPanePosition = "bottom"`
- `browserGrouped = true`

This favors a navigation-first pane with more structure and less preview density.

### 3. `review` preset
Applies:
- `browserPaneHeight = 14`
- `browserPanePreview = true`
- `browserPanePosition = "top"`
- `browserGrouped = true`

This favors a more context-rich, inspection-oriented pane.

### 4. Regression coverage expanded
The existing pane preset test now verifies:
- `compact`
- `detailed`
- `navigation`
- `review`

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because the pane is no longer just configurable — it is now more clearly adapted to distinct operator workflows.

That is a more useful preset layer than a simple binary compact/detailed split.

## Design insight

The key principle still holds:

> presets remain bundles of rendering/layout state layered over canonical runtime-backed browser state.

No session or branch semantics change when a preset changes.

## Recommended next move

The strongest next step is now:

> **continue refining split-view ergonomics with either more operator presets or stronger pane/browser mode distinctions.**

## Bottom line

This tranche makes pane presets more expressive and more aligned with real workflow differences, while preserving the same truth-first architecture.
