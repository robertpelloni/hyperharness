# TUI Confirm-Before-Switch Tranche

Date: 2026-04-04

## Summary

This tranche adds **confirmation-before-switch** behavior to the cursor-driven TUI tree browser.

Before this tranche, the browser already had:
- cursor navigation
- type-to-filter
- selected-entry preview
- branch-summary preparation preview
- direct Enter-to-switch behavior

That meant users could see useful context before switching, but the actual branch switch still happened immediately on Enter.

This tranche introduces a lightweight confirmation step that better matches the increasing richness of the preview surface.

## What was added

### 1. Confirm-pending UI state
Added to the TUI model:
- `browserConfirmPending bool`

This is transient UI-only state used to represent that the currently selected browser item has been armed for switching but not yet executed.

### 2. Browser keyboard behavior updated
While browser mode is active:

#### First `Enter`
- arms the currently selected entry for switching
- does **not** switch yet
- enters confirm-pending mode

#### Second `Enter`
- confirms the switch
- executes the canonical runtime-backed switch
- closes the browser

#### `Y`
- confirms the switch while in confirm-pending mode

#### `N`, `Esc`, or `Backspace`
- cancels the pending confirmation without closing the whole app

#### Navigation keys during confirm-pending mode
- are intentionally ignored so the pending selection remains stable until the user confirms or cancels

### 3. Browser view updated
The browser rendering now shows different guidance depending on state:
- normal browser mode instructions
- confirm-pending instructions

When confirm-pending is active, the preview pane includes:
- `[Confirm]`
- explicit instructions about how to confirm or cancel

This is important because the browser is now not just a selector, but a small guided transition workflow.

## Verification added

Updated the browser-mode test in `tui/slash_test.go` to verify:
- first `Enter` enters confirm-pending mode
- confirm prompt appears in the browser view
- second `Enter` performs the actual switch
- the browser closes only after confirmed switching
- the canonical leaf actually changes after confirmation

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche is important because the browser now supports a better decision boundary.

At this point the browser preview is rich enough to inform a choice, so requiring a second confirmation step is a natural upgrade:
- first inspect
- then arm
- then confirm

That gives users a safer branch-switching workflow without weakening the underlying semantics.

## Design insight

The key design principle preserved here is:

> **confirmation is a UI workflow layer only; the actual switch still happens through the same canonical runtime path as before.**

This is exactly the right pattern.

The project continues to add ergonomics without introducing alternate sources of session/branch truth.

## What is still missing

This tranche improves safety and ergonomics, but some richer browser features are still open:

1. **folding/grouping**
   - useful for large session trees

2. **more graph-like tree presentation**
   - current browser is still list-oriented

3. **richer branch previews**
   - perhaps showing a fuller summary snippet or child grouping before confirmation

## Recommended next move

The strongest next step is now:

> **add folding/grouping or richer tree layout in the browser, now that navigation, preview, filtering, and confirmation are all in place.**

At this point the TUI tree browser has a solid interaction backbone, and the next improvements can focus on scale and clarity.

## Bottom line

This tranche turns the tree browser from:
- immediate switch on selection

into:
- a small, guided confirm-before-switch workflow.

That is the right next ergonomics step for a truth-first branch navigation surface.
