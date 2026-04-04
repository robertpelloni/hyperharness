# Go Foundation Assimilation Design

## Design Summary
Build a **Go-native harness core** that uses:
- Go for runtime, orchestration, session control, and TUI integration.
- exact tool-contract manifests for model-facing compatibility,
- HyperCode/Borg as the default control-plane substrate,
- optional adapters for upstream-compatible behavior,
- a modular architecture that allows truthfully staged parity.

## Why Go
Go is the preferred foundation for this repository because it aligns with the product direction:
- easy static distribution,
- fast startup,
- simple concurrency for multi-agent orchestration,
- strong fit with Charmbracelet TUI tooling,
- natural adjacency to Ollama and other Go-native infrastructure,
- lower contributor friction than a Rust-first core.

## High-Level Architecture

```mermaid
flowchart TD
    UI[TUI / CLI / JSON / RPC] --> Harness[Go Harness Core]
    Harness --> Agent[Agent Runtime]
    Harness --> Sessions[Session + Branching + Compaction]
    Harness --> Compat[Exact Tool Contract Registry]
    Harness --> Adapters[Assimilation Adapters]
    Harness --> HC[HyperCode Control Plane]

    Agent --> Context[Context Engine]
    Agent --> Tools[Native Tools]
    Agent --> Delegation[Delegation + Verification]

    Context --> RepoMap[Repo Map / Code Graph]
    Context --> Memory[Memory + Retrieval]

    HC --> MCP[MCP Routing + Discovery]
    HC --> Providers[Provider Routing + Auth]
    HC --> Imports[Session/History Imports]
    HC --> Observability[Operator Telemetry]
```

## Package-Level Design

### 1. `foundation/pi`
Purpose: define the Pi-derived core harness contract in Go.

Current responsibilities:
- thinking levels
- message delivery modes
- transport selection
- tool execution modes
- run-event vocabulary
- built-in tool definitions
- default foundation spec

Future responsibilities:
- actual Go agent runtime
- session model
- extension host API
- prompt template runtime
- settings loader and migration

### 2. `foundation/compat`
Purpose: preserve exact model-facing tool contracts.

Responsibilities:
- exact name/parameter/result contract storage
- maturity labeling (`planned`, `bridged`, `speced`, `native`, `verified`)
- source-specific compatibility catalogs
- later: contract tests against live/native implementations

### 3. `foundation/assimilation`
Purpose: keep the imported tool ecosystem visible and explicit.

Responsibilities:
- inventory of upstream systems
- strengths and unique traits
- assimilation strategy per source
- future: per-feature port map and parity scoreboard

### 4. HyperCode integration boundary
HyperCode should be treated as an external-but-local substrate.

Responsibilities retained by HyperCode:
- provider routing and failover
- MCP aggregation and server lifecycle
- imported session and memory discovery
- operator control-plane APIs
- mesh and runtime status

Responsibilities owned by the new harness:
- model-facing UX
- exact tool contracts
- agent loop behavior
- session interaction model
- coding-agent specialization

## Key Design Principles

### Principle 1: Behavior compatibility beats source similarity
We are not porting line-by-line. We are reproducing observable capability with a better internal design.

### Principle 2: Exact tool names are first-class data
Tool compatibility should not be implicit in code spread across packages. It should live in an explicit registry.

### Principle 3: Native first, bridge second, never lie
Every imported feature should be labeled as:
- bridged via HyperCode or another substrate,
- specified but not yet implemented,
- native,
- or verified.

### Principle 4: Keep the foundation small
Pi’s best lesson is that the core should stay narrow and extensible.

### Principle 5: Use upstream tools as references, not anchors
We should avoid creating a dependency graph where the new harness cannot stand without submodules.

## Compatibility Strategy

### Phase A: contract-first
Define exact tool and event contracts before rewriting behavior.

### Phase B: bridged parity
Route through HyperCode/Borg or local adapters where native parity is not ready.

### Phase C: native replacement
Replace bridged implementations one capability family at a time.

### Phase D: verified parity
Add contract and snapshot tests for each feature family.

## Current Foundation Artifacts Added in This Phase
- `foundation/pi/foundation.go`
- `foundation/compat/*`
- `foundation/assimilation/*`
- `cmd/foundation.go`

## Baseline Risks Observed
- current repo contains stub implementations advertised as parity-complete but not actually complete,
- current Go test baseline fails in unrelated areas (`mcp`, `orchestrator`, fixture package),
- tool naming/contracts are not centrally defined in the existing Go code,
- current orchestration and TUI layers are not yet backed by a truthful native core.

## Recommended Next Technical Moves
1. Implement a real session store and evented agent runtime under `foundation/pi`.
2. Port Aider-style repo map and edit engines into a dedicated context package.
3. Add HyperCode-backed provider and MCP adapters behind stable interfaces.
4. Replace placeholder tools with native implementations plus contract tests.
5. Add verification, delegation, and background session services.
