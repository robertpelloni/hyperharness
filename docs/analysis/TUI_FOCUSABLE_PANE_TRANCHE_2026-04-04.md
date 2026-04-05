# TUI Focusable Pane Tranche

Date: 2026-04-04

## Summary

This tranche extends the persistent tree pane so it can be **focused and navigated directly** without switching into full modal browser mode.

Before this tranche, the pane could stay visible and auto-refresh after session mutations, but it remained essentially passive.

This tranche changes that by letting the operator temporarily direct keyboard navigation into the pinned pane itself.

## What was added

### 1. Focus state for the pinned pane
Added to the TUI model:
- `browserPinnedFocus bool`

This indicates that keyboard navigation should be routed into the pinned pane rather than only into the input line or the modal browser.

### 2. New slash command
Added in `tui/slash.go`:
- `/tree-pane-focus`

Behavior:
- toggles focus on the pinned tree pane
- if the pane is not pinned, reports an error
- when disabling focus, clears any pending confirmation state

### 3. Focus-aware key handling
When the pane is pinned **and focused**, the TUI now routes keyboard input into the pane using the same canonical browser semantics:
- `Up` / `Down` move selection
- `Left` / `Right` collapse/expand subtrees
- `Tab` toggles grouping
- type-to-filter still works
- `Enter` arms/executes branch switching
- `Y` confirms a pending switch
- `N`, `Esc`, or `Backspace` cancel pending confirmation

This means the persistent pane is no longer just a passive companion view — it can become an active navigation surface.

### 4. View feedback
When the pinned pane has focus, the TUI view now shows:
- `[Tree Pane Focused]`

That makes the current keyboard-target explicit to the operator.

### 5. Pane lifecycle hygiene
When the pane is hidden via `/tree-pane`, the focus and pending-confirmation state are cleared as part of teardown.

## Verification added

Expanded `tui/slash_test.go` with a focused regression that verifies:
- `/tree-pane` pins the pane
- `/tree-pane-focus` enables focus
- the focused pane marker appears in the view
- keyboard navigation moves the pane selection
- Enter enters confirm-pending mode
- `Y` confirms the switch
- the canonical leaf changes after confirmation

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because the persistent pane is now an active navigation surface rather than only a live readout.

That is a major usability step for long-running sessions:
- the user can keep the pane visible
- shift focus into it when needed
- navigate and switch branches
- then return to normal prompt interaction

All while the canonical runtime remains the sole source of session/branch truth.

## Design insight

The key architectural principle preserved here is:

> the focused pane reuses the same browser semantics and canonical runtime operations rather than inventing a second navigation model.

That is exactly right.

The pane becomes more powerful, but only by becoming another view/controller over the same verified runtime behavior.

## What is still missing

This tranche significantly strengthens pane ergonomics, but a few next-level improvements remain possible:

1. **split-view layout refinements**
   - e.g. more deliberate composition of pane + prompt + history

2. **independent pane scrolling or viewport controls**
   - especially for very large trees

3. **more explicit pane/browser mode separation in UI styling**
   - making the focus state even clearer visually

## Recommended next move

The strongest next step is now:

> **improve split-view ergonomics or add stronger pane-focused viewport behavior**, since the pane is now persistent, live, and focusable.

## Bottom line

This tranche makes the persistent tree pane truly interactive.

It is now not just visible and live, but also directly navigable — all while staying rooted in the same canonical foundation runtime.
