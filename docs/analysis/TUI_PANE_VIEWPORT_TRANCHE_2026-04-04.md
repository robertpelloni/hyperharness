# TUI Pane Viewport Tranche

Date: 2026-04-04

## Summary

This tranche adds **viewport behavior** to the persistent/focusable tree pane so large session trees remain readable while the selected entry stays in view.

Before this tranche, the persistent pane was:
- visible
- live
- and focusable

But it still rendered the full visible item list at once.

This tranche gives the pinned pane a bounded window while keeping the active/selected area centered when possible.

## What was added

### 1. Viewport-aware tree rendering
The tree browser renderer now accepts:
- `maxVisible`
- `title`

This allows the same renderer to serve both:
- full browser mode
- bounded pane mode

### 2. Window calculation
The renderer now computes a window over the visible item list:
- center around the selected item when possible
- clamp safely at start/end boundaries

### 3. Pane-specific viewport metadata
When the visible set exceeds the pane window, the pane now shows:
- `showing=X-Y of Z`

That gives the operator explicit awareness that the pane is showing a viewport rather than the full tree.

### 4. Distinct pane/browser headings
The same renderer can now label itself appropriately:
- `[Foundation Tree Browser]`
- `[Foundation Tree Pane]`

This keeps the interaction modes clearer.

### 5. Pane rendering now uses a bounded window
The pinned pane now renders with a limited viewport size instead of the entire visible tree list.

That is important for larger session trees because it avoids the pane taking over the whole screen while still remaining useful.

## Verification
Validated successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because a persistent pane should help the user, not overwhelm the screen.

Once the browser became:
- persistent
- live
- focusable

viewport behavior became the next natural ergonomics requirement.

Now the pane is better suited to larger sessions while still keeping selection context visible.

## Design insight

The key design rule remained intact:

> viewporting is a rendering concern only; the visible window is still derived from the same canonical runtime-backed browser item list.

That keeps the pane useful without changing any underlying session/branch semantics.

## Recommended next move

The strongest next step is now:

> **continue improving pane/browser coexistence, or move toward stronger split-view layout behavior.**

At this point, the pane has become a real long-session companion surface.

## Bottom line

This tranche makes the pinned pane scale better with larger trees by giving it a bounded, selection-aware viewport.
