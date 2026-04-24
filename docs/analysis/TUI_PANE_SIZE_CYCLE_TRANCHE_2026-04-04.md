# TUI Pane Size Cycle Tranche

Date: 2026-04-04

## Summary

This tranche adds a **quick pane-height cycle** command for the persistent tree pane.

Before this tranche, pane height could be controlled explicitly via `/tree-pane-size <n>`. That was flexible but slightly verbose for repeated adjustments.

This tranche adds a simple convenience cycle.

## What was added

### 1. New slash command
Added:
- `/tree-pane-size-cycle`

### 2. Cycle behavior
The command cycles through a set of common heights:
- `6`
- `8`
- `10`
- `12`
- `14`
- back to `6`

This provides a quick way to tune pane density during long sessions.

## Verification added

Added focused regression coverage verifying that repeated size cycling advances through the expected values.

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because small repeated layout adjustments become expensive if they require explicit argument entry every time.

A simple cycle command reduces that friction.

## Design insight

The key principle remains the same:

> the size cycle is a convenience layer over pane layout state, not a new semantic concept.

## Recommended next move

The strongest next step is now:

> continue adding small operator-polish controls where they reduce repeated pane/browser adjustment overhead.

## Bottom line

This tranche makes pane height tuning faster and more ergonomic in practice.
