# TUI Tree Layout Cues Tranche

Date: 2026-04-04

## Summary

This tranche adds richer **tree layout cues** to the cursor-driven TUI browser so larger session histories are easier to scan.

Before this tranche, the browser already supported:
- cursor navigation
- filtering
- preview-before-switch
- branch-summary preparation preview
- confirmation-before-switch

But the list rendering itself was still fairly flat.

This tranche improves structural readability by surfacing:
- entry depth
- child counts
- indentation and simple branch glyphs

## What was added

### 1. Enriched browser item structure
Expanded `TreeBrowserItem` with:
- `ParentID`
- `Depth`
- `ChildCount`

These are still derived from the canonical runtime/session data; they do not create any new source of truth.

### 2. Depth computation in browser item generation
`buildFoundationTreeBrowser` now computes:
- depth by parent chain order
- child count via canonical `GetChildren`

This gives each browser item a structural position in the session tree.

### 3. Richer browser rendering
The browser list now renders entries with:
- indentation by depth
- a simple branch glyph (`•` for root-ish entries, `└─` for deeper ones)
- child-count display when present
- all prior cues retained:
  - active leaf marker
  - cursor marker
  - kind
  - label
  - preview text

### 4. Preview pane enriched
The preview pane now includes:
- depth
- child count

This helps the operator understand not just what an entry is, but where it sits structurally in the session tree.

## Verification added

Expanded the browser-mode test to verify that the rendered browser view includes child-count/tree-layout cues.

### Full validation
Verified successfully:

```bash
go test ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche is important because the browser is now better at communicating structural information, not just selection state.

That matters as session trees grow larger and branching becomes more common.

Without stronger layout cues, a browser can technically work but still be mentally expensive to use.

With this change, the browser becomes easier to scan and reason about while staying faithful to the canonical runtime.

## Design insight

The most important design principle preserved here is:

> **layout cues are rendered from runtime-derived structure, not invented UI metadata.**

That is exactly the correct direction.

The browser is becoming visually richer, but still only by exposing more of the real session structure already present underneath.

## What is still missing

Even with these cues, the browser remains list-oriented rather than a fully graphical tree.

Likely future improvements include:

1. **folding/grouping**
   - collapse subtrees or group by parent branch

2. **more graph-like tree rendering**
   - especially for complex branching sessions

3. **branch preview summaries in a dedicated pane**
   - separating structural preview from switch preview even more clearly

## Recommended next move

The strongest next step is now:

> **add folding/grouping to the browser so large trees become more navigable without losing the canonical runtime grounding.**

That would continue the current pattern:
- first make the semantics truthful
- then make the UX progressively easier to use at scale

## Bottom line

This tranche improves the browser’s readability by surfacing structural depth and child-count cues.

It is another incremental but meaningful step from:
- functional tree tooling

toward:
- a genuinely usable long-session branch browser.
