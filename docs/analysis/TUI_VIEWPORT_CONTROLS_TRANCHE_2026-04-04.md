# TUI Viewport Controls Tranche

Date: 2026-04-04

## Summary

This tranche strengthens pane/browser navigation for large trees by adding **Home/End/PgUp/PgDn viewport controls**.

Before this tranche, the pane/browser supported:
- cursor navigation
- filtering
- confirmation-before-switch
- grouping
- folding
- persistent pane behavior
- focusable pane behavior
- viewport-aware rendering

But moving across larger trees still required repeated line-by-line navigation.

This tranche fixes that.

## What was added

### 1. Home / End controls
For both modal browser mode and focused pinned pane mode:
- `Home` jumps selection to the first visible item
- `End` jumps selection to the last visible item

### 2. Page-style movement
For modal browser mode:
- `PgUp` moves selection upward by a coarse page step
- `PgDn` moves selection downward by a coarse page step

For focused pinned pane mode:
- `PgUp` moves upward by the current pane height
- `PgDn` moves downward by the current pane height

This ties pane-focused page movement directly to the configured viewport size.

### 3. Shared numeric helpers
Added a simple `min` helper alongside `max` so index clamping remains safe and explicit.

### 4. Browser instructions updated
The browser renderer now documents:
- `PgUp`
- `PgDn`
- `Home`
- `End`

in its navigation guidance.

## Verification added

Expanded `tui/slash_test.go` with `TestTreePaneViewportControls`, verifying:
- `End` moves selection to the end
- `Home` returns selection to zero
- `PgDn` advances by at least one viewport step
- `PgUp` returns toward the top
- viewport metadata is visible when appropriate

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche is important because once the pane/browser became persistent, focusable, and viewport-aware, the next usability bottleneck was long-distance movement.

Page-scale controls make the browser much more practical for deeper or longer session trees.

## Design insight

The key principle stayed intact:

> viewport controls move selection over canonical runtime-derived visible items; they do not introduce any new tree semantics.

That is exactly the right boundary.

## Recommended next move

The strongest next step is now:

> **continue strengthening split-view ergonomics**, likely with stronger pane-specific layout controls or clearer active-pane affordances.

## Bottom line

This tranche makes large tree navigation faster and more comfortable without compromising the same truth-first architecture.
