# TUI Mode Distinction Tranche

Date: 2026-04-04

## Summary

This tranche strengthens the visual distinction between:
- modal browser mode
- passive pinned pane mode
- focused pinned pane mode

Before this tranche, the system already supported all three states functionally, but the rendered titles and visual affordances could still blur together.

This tranche makes those states much more explicit in the UI.

## What was added

### 1. Distinct modal browser title
The modal browser now renders with:
- `[Foundation Tree Browser :: Modal]`

This makes it immediately obvious that the user is in the full modal tree browser rather than looking at a pinned pane.

### 2. Distinct passive pane title
The pinned pane now renders with:
- `[Foundation Tree Pane :: Passive]`

This clarifies that the pane is visible but not currently the keyboard focus target.

### 3. Distinct focused pane title
When the pinned pane has keyboard focus, it now renders with:
- `[Foundation Tree Pane :: Focused]`

This replaces the weaker prior indicator and makes the active interaction mode clearer right at the pane title.

## Verification added

Expanded the TUI regression coverage so tests now verify:
- modal browser title appears in browser mode
- passive pane title appears when the pane is pinned but not focused
- focused pane title appears when pane focus is enabled

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because the tree subsystem now has enough modes that subtle distinctions are no longer enough.

The operator needs to know at a glance:
- whether they are in a modal browser
- whether the pane is passive
- whether the pane owns keyboard focus

Making those distinctions explicit reduces mode confusion and supports longer, denser sessions.

## Design insight

The key principle here is:

> state distinctions should be visible in the UI without requiring the user to infer them from behavior alone.

That is especially important once the tree browser and pane become persistent, focusable, grouped, foldable, and configurable.

## Recommended next move

The strongest next step is now:

> continue refining pane/browser coexistence with clearer layout and interaction affordances, since the system now has a more mature mode model.

## Bottom line

This tranche improves interaction clarity by making browser and pane modes visually explicit rather than merely implied.
