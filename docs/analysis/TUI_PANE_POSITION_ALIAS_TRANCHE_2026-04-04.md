# TUI Pane Position Alias Tranche

Date: 2026-04-04

## Summary

This tranche adds **direct top/bottom aliases** for persistent tree pane placement.

Before this tranche, pane position could already be controlled with:
- `/tree-pane-position top`
- `/tree-pane-position bottom`
- `/tree-pane-position-toggle`

That was already functional, but direct aliases reduce friction for the most common explicit placement requests.

## What was added

### New slash commands
- `/tree-pane-top`
- `/tree-pane-bottom`

These route directly into the existing pane-position handler.

## Verification added

Added focused regression coverage verifying:
- `/tree-pane-top` sets the pane position to `top`
- `/tree-pane-bottom` sets the pane position to `bottom`

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because the pane/browser subsystem is now rich enough that quick, explicit commands for the most common states improve confidence and repeated-use speed.

## Design insight

The key principle remains unchanged:

> direct aliases are just ergonomic wrappers over the same underlying position state transitions.

## Recommended next move

The strongest next step is now:

> continue adding explicit aliases only when they clearly improve confidence or reduce friction compared to generic parameterized commands.

## Bottom line

This tranche is a small but useful polish improvement that makes pane placement faster and more explicit in daily use.
