# TUI Pane Help Tranche

Date: 2026-04-04

## Summary

This tranche adds a focused **`/tree-pane-help`** command so operators can see the tree pane/browser control surface without scanning the full global `/help` output.

## What was added

### New slash command
- `/tree-pane-help`

### Behavior
It emits a compact pane/browser-specific help block covering:
- pane visibility controls
- pane focus controls
- size controls
- position controls
- preview controls
- grouped controls
- presets
- browser controls
- refresh/reset/clear controls

## Verification added

Added focused regression coverage verifying:
- `/tree-pane-help` emits the expected help header
- the output includes pane quick-toggle controls such as `/tree-pane-preview-toggle`

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because the pane/browser subsystem now has enough commands that a focused local help surface materially improves usability.

## Design insight

The key principle is simple:

> subsystem-local help reduces operator friction without changing any underlying runtime behavior.

## Recommended next move

The strongest next step remains continued operator polish where command discoverability and repeated-use efficiency improve meaningfully.

## Bottom line

This tranche gives the tree pane subsystem its own compact discoverability surface, which is increasingly justified as the command family grows.
