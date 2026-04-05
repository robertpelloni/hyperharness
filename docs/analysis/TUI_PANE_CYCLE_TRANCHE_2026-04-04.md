# TUI Pane Cycle Tranche

Date: 2026-04-04

## Summary

This tranche adds a **quick preset cycle command** for the persistent tree pane.

Before this tranche, operators could switch pane layouts with `/tree-pane-preset`, but they still had to type the full preset name.

This tranche adds a faster cycling control for common workflow transitions.

## What was added

### 1. New slash command
Added in `tui/slash.go`:
- `/tree-pane-cycle`

### 2. Cycle behavior
The command cycles through the currently supported pane presets in this order:
1. `compact`
2. `navigation`
3. `detailed`
4. `review`
5. back to `compact`

### 3. State-based detection
The cycle command infers the next preset from the current pane state rather than storing a separate preset enum.

That keeps it lightweight and consistent with the existing preset implementation model.

## Verification added

Expanded `tui/slash_test.go` to verify that repeated `/tree-pane-cycle` invocations move through:
- compact
- navigation
- detailed
- review

with the expected state changes for each.

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because it reduces friction for operators who want to move quickly between pane personalities during long sessions.

It complements the named preset command nicely:
- `/tree-pane-preset` is explicit
- `/tree-pane-cycle` is fast and ergonomic

## Design insight

The key principle remains intact:

> the cycle command is just a convenience layer over the existing pane preset state bundles.

It does not introduce a new semantic model.

## Recommended next move

The strongest next step is now:

> continue building higher-level operator controls and recovery ergonomics around the tree pane/browser subsystem.

## Bottom line

This tranche makes the pane layout system faster to use in practice, which is exactly the kind of operator polish that becomes valuable once the pane subsystem is already powerful.
