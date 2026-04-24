# TUI Split-View Ergonomics Notes

Date: 2026-04-04

## Summary

The current TUI now supports a persistent, live, focusable, viewport-aware, and size-adjustable tree pane backed by the canonical foundation runtime.

This note records the current state after the pane-size tranche and highlights the next likely ergonomics frontier: stronger split-view coexistence.

## Current pane capabilities
- persistent visibility via `/tree-pane`
- focus toggle via `/tree-pane-focus`
- adjustable height via `/tree-pane-size <n>`
- auto-refresh after session mutations
- viewport-aware rendering
- shared runtime semantics with the browser mode

## Why split-view ergonomics are the next step
At this point the pane is no longer experimental or passive. It is:
- persistent
- live
- interactive
- and configurable

That means the next UX gains are more about layout quality than raw capability.

## Likely next directions
1. Better pane/history/prompt balance in the TUI view
2. More explicit visual separation between pane mode and normal prompt mode
3. Potential independent pane viewport controls
4. Optional alternate pane positions or stronger split-view composition

## Key architectural rule
Any future split-view work should continue to treat the pane as:
- a rendering/controller surface
- over canonical runtime-backed browser state

and not as a new owner of session/branch truth.
