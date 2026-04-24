# TUI Pane Preview Toggle Tranche

Date: 2026-04-04

## Summary

This tranche adds a **preview toggle** for the persistent tree pane so operators can choose between:
- a denser compact pane
- or a richer preview-enabled pane

Before this tranche, the pane inherited the same detailed preview behavior as the browser view. That was powerful, but sometimes too visually heavy for long-running sessions.

This tranche makes that behavior operator-configurable.

## What was added

### 1. Pane preview state
Added to the TUI model:
- `browserPanePreview bool`

Default:
- `true`

### 2. Shared renderer support
The shared tree renderer now accepts a `showPreview` flag, allowing:
- full browser mode to always keep previews
- persistent pane mode to respect the pane-preview setting

### 3. New slash command
Added:
- `/tree-pane-preview <on|off>`

Behavior:
- sets preview visibility for the persistent pane
- confirms the new state in TUI history

### 4. View integration
The pinned pane now renders preview details only when `browserPanePreview` is enabled.

This gives the operator direct control over how much detail the pane should occupy during long sessions.

## Verification
A focused pane-preview toggle test is now present in `tui/slash_test.go`.

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because the persistent pane is now flexible not just in size and position, but also in information density.

That is an important coexistence improvement for long sessions.

## Design insight

The key principle remains intact:

> preview visibility is a pure rendering concern layered over canonical runtime-derived browser state.

No session or branch semantics change when preview is hidden.

## Recommended next move

The strongest next step is now:

> **continue refining pane/browser coexistence**, likely with stronger pane-focused interaction affordances or alternate layout presets.

## Bottom line

This tranche makes the persistent pane more adaptable by letting the operator choose a compact or preview-rich presentation without changing any underlying runtime truth.
