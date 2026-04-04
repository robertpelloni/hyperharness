# Foundation Reconciliation and Verified Status

Date: 2026-04-04

## Executive Summary

The repository now contains **two visible Go tracks**:

1. An **existing native foundation track** already present in the repo under:
   - `foundation/pi`
   - `foundation/compat`
   - `foundation/assimilation`
   - `foundation/repomap`
   - `foundation/adapters`
   - `foundation/orchestration`
2. A newer `internal/*` and `cmd/hyperharness/*` track added during the initial Go-port push.

After inspecting the repo comprehensively, the correct strategic decision is:

> **Treat `foundation/pi` + companion foundation packages as the canonical Pi-to-Go foundation, and evolve that path first.**

This is the right move because the existing foundation packages are already:
- test-backed,
- integrated into the repo’s existing Cobra CLI surfaces,
- aligned with `foundation/compat` and `foundation/assimilation`, and
- conceptually closer to the clean-room compatibility layer we need for Borg/HyperCode integration.

## What I Verified

### Module path and dependency state
- Restored the module path to the repo’s real import graph:
  - `module github.com/robertpelloni/hypercode`
- Ran dependency reconciliation successfully.
- Verified that the foundation packages build and test successfully.

### Verified passing test target
The following command now passes:

```bash
go test ./foundation/...
```

That means these packages are currently green:
- `foundation/adapters`
- `foundation/assimilation`
- `foundation/compat`
- `foundation/orchestration`
- `foundation/pi`
- `foundation/repomap`

## Native Pi parity expanded

Before this tranche, the native Pi foundation and compatibility catalog were centered on **4 tools**:
- `read`
- `write`
- `edit`
- `bash`

I expanded the verified native Pi surface to **7 tools**:
- `read`
- `write`
- `edit`
- `bash`
- `grep`
- `find`
- `ls`

### Changes made

#### `foundation/pi/runtime_types.go`
Added native tool input types for:
- `GrepToolInput`
- `FindToolInput`
- `LSToolInput`

#### `foundation/pi/foundation.go`
Expanded the model-facing builtin tool descriptors and JSON schema contracts to include:
- `grep`
- `find`
- `ls`

#### `foundation/pi/tools_native.go`
Expanded `DefaultToolHandlers()` to include:
- `executeGrepTool`
- `executeFindTool`
- `executeLSTool`

Implemented native handlers with the following behavior:
- `grep`: recursive file search with simple regex/literal matching and optional glob filter
- `find`: recursive file discovery by glob pattern
- `ls`: directory listing with limit support

#### `foundation/compat/default_catalog.go`
Expanded the exact-name compatibility catalog so Pi-native contracts now include:
- `grep`
- `find`
- `ls`

#### Tests updated
- `foundation/pi/foundation_test.go`
- `foundation/pi/runtime_test.go`
- `foundation/pi/tool_snapshot_test.go`
- `foundation/compat/catalog_test.go`

All updated tests now pass.

## Architectural Insight: which Go track should win?

### Canonical near-term foundation
The canonical near-term foundation should be:

- `foundation/pi` for exact Pi-derived runtime/tool/session contracts
- `foundation/compat` for model-facing compatibility inventory
- `foundation/assimilation` for upstream inventory and sequencing
- `foundation/repomap` for Aider-style repo condensation
- `foundation/adapters` for Borg/HyperCode/MCP/provider bridging
- `foundation/orchestration` for higher-order planning/daemon layers

### Why not center on `internal/*` first?
The newer `internal/*` tree contains valuable ideas and a lot of implementation mass, but today it is:
- less integrated into the repo’s existing CLI shape,
- more speculative,
- more duplicate with the already-present foundation packages,
- and not yet the cleanest path to verifiable parity.

Recommendation:

> Continue using `foundation/*` as the implementation truth and selectively merge ideas from `internal/*` into that track where they improve the design.

## HyperCode alignment insights

After inspecting `../hypercode`, the most important design principle is:

> **Do not reimplement HyperCode’s control-plane truth inside the harness.**

HyperCode should remain the owner of:
- MCP aggregation,
- provider routing,
- memory/session substrate at the control-plane level,
- operator visibility and diagnostics,
- mesh/service discovery.

The harness should own:
- exact model-facing tool contracts,
- agent UX and terminal workflows,
- local session ergonomics,
- Pi/Claude/Aider/OpenCode-compatible harness behavior,
- and clean adapters into HyperCode.

### Correct responsibility split

#### HyperCode owns
- orchestration truth
- MCP runtime management
- provider health/routing/fallback
- mesh and inventory visibility
- operator dashboard and daemon surfaces

#### Hyperharness owns
- agent-facing shell
- exact tool contracts
- tool-call ergonomics
- prompt/session/workflow UX
- native tool execution layer
- compatibility shims and harness semantics

This is the cleanest architecture for Borg as well.

## Current verified state

### Verified green
- Foundation packages compile and pass tests.
- Native Pi parity exists and is now expanded to 7 tools.
- Compatibility catalog reflects those 7 native Pi tools.
- Assimilation inventory already documents upstream systems and intended strategy.

### Not yet fully reconciled
- The whole repository does **not** yet build cleanly as one unified binary path.
- The newer `internal/*` track and the older root/cmd/tui/orchestrator track are not yet fully harmonized.
- Full-feature parity with every upstream tool is still a long program, not a single-step completion.

## Recommended next implementation sequence

### Phase A — stabilize the real foundation
1. Keep `foundation/*` green at all times.
2. Continue expanding `foundation/pi` exact native parity:
   - remaining Pi-native tools and modes
   - session semantics
   - compaction and branching features
3. Treat `foundation/compat` as the source of model-facing contract truth.

### Phase B — assimilate Aider and OpenCode strengths
1. Port repo map ranking and context condensation deeper into `foundation/repomap`.
2. Add Aider-style edit strategies and git-aware flows.
3. Add OpenCode/Crush/Gemini-style TUI ergonomics to the UI layer.

### Phase C — bridge, do not duplicate, HyperCode
1. Use `foundation/adapters` to expose HyperCode provider routing.
2. Use HyperCode as the default MCP control plane.
3. Use HyperCode memory/session APIs where daemon truth matters.
4. Keep the harness UX local-first while delegating control-plane truth.

### Phase D — exact compatibility inventory growth
Expand `foundation/compat` to include model-facing tool contracts for:
- Claude Code expectations
- Aider-facing behaviors
- Codex/Code CLI contract expectations
- Goose/Open Interpreter style tool surfaces
- MCP-exposed external tool catalogs

## Specific recommendations

1. **Promote `foundation/compat` to the single source of truth for tool contracts.**
2. **Promote `foundation/pi` to the single source of truth for Pi-native behavior.**
3. **Merge useful ideas from `internal/*` into `foundation/*`, not vice versa.**
4. **Avoid overclaiming parity until verified by tests/snapshots.**
5. **Treat HyperCode as substrate, not competitor, inside the same local ecosystem.**

## Result of this tranche

This tranche did not just add code. It clarified the architecture and reduced strategic ambiguity.

The project is now in a better state because:
- the module path is aligned,
- the verified foundation is green,
- the native Pi contract surface is broader,
- and the canonical path forward is clearer.

That is the right kind of progress for a project this ambitious.
