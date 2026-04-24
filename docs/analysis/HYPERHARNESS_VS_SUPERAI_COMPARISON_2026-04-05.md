# Hyperharness vs SuperAI — Architecture & Functionality Comparison

Date: 2026-04-05

## Executive Summary

Both repositories share a common ancestor module (`github.com/robertpelloni/hypercode`) and nearly
identical top-level package layout (`agent/`, `agents/`, `hypercode/`, `cmd/`, `config/`, `foundation/`,
`git/`, `llm/`, `mcp/`, `orchestrator/`, `repl/`, `security/`, `tools/`, `tui/`). They are **sibling
projects** that have diverged in scope:

| Dimension | hyperharness | superai |
|---|---|---|
| Go module | `github.com/robertpelloni/hyperharness` | `github.com/robertpelloni/hypercode` |
| Go source files | **138** | **2 703** |
| Lines of Go | **~19 K** | **~123 K** |
| Test files (*_test.go) | **32** | **733** |
| Own Go packages | **9** (+ `internal/*`) | **8** |
| Build | `go build .` + `go test ./tui ./cmd ./foundation/...` green | Same command, same result |
| TUI slash commands | 45+ (tree-pane, browser, presets, help, summary) | 22 (foundation session, tree, summary) |
| Analysis docs | 48+ tranches in `docs/analysis/` | 1 empty `README.md` |

**Verdict:** superai is the **broader** codebase by raw volume (6.5× more Go, 23× more tests).
hyperharness is the **narrower but more polished** port of the `foundation/pi` track with a richer
TUI surface and a documented tranche-by-tranche development trail.

---

## Deep-Dive: Architecture

### Shared Architecture (both repos)

```
┌─────────────────────────────────────────────────────────┐
│                        CLI (cobra)                      │
├───────────────┬───────────────┬─────────────────────────┤
│    TUI        │   foundation  │   agents/agent          │
│  bubbletea    │  ├─ pi (rt)   │   orchestrator          │
│  chat+browser │  ├─ compat    │   mcp                   │
│  slash cmds   │  ├─ adapters  │   security              │
│               │  ├─ repomap   │   git                   │
│               │  └─ assimilation                      │
└───────────────┴───────────────┴─────────────────────────┘
```

Both rely on:
- **charmbracelet/bubbletea** for TUI
- **spf13/cobra** for CLI commands
- **mark3labs/mcp-go** for MCP server support
- **go-git/go-git/v5** for repository introspection
- **gofiber/fiber** for optional HTTP endpoints

### Structural Differences

| Aspect | hyperharness | superai |
|---|---|---|
| Extra top-level packages | `internal/*` (8 sub-packages), `cmd/hyperharness` | none |
| Extra agent files | `auton.go`, `codemode.go`, `council.go`, `disclosure.go`, `interfaces.go`, `provider_hypercode.go`, `rag.go`, `shell_assistant.go` | `director.go` (basic), `provider.go`, `provider_stub.go` |
| Extra agent logic | async, autocompleter, autopilot, diff, oracle, shell (all in superai) | similar set but fewer files |
| foundation/pi extra | `tools_extra.go`, `tools_extra_test.go` | none |
| Analysis trail | 48 markdown docs in `docs/analysis/` | none |
| HANDOFF.md | updated every tranche | present but minimal |
| Module path correctness | `github.com/robertpelloni/hyperharness` (**canonical**) | `github.com/robertpelloni/hypercode` (legacy name) |
| Go 1.26 feature flag | yes | yes |

### Package-by-Package Comparison

#### tui/

| Feature | hyperharness | superai |
|---|---|---|
| `chat.go` — Bubble Tea model | ✅ | ✅ (same) |
| `slash.go` — slash commands | 45+ commands | 22 commands |
| `slash_test.go` | 32 regression test func | fewer |
| `foundation_bridge.go` | ✅ | ✅ (same) |
| `shell.go` | ✅ | ✅ (same) |
| Tree-pane subsystem | extensive (help, summary, presets, sizing, focus, preview, grouping, position, filter, viewport, refresh, reset) | minimal |

**Winner: hyperharness** — the TUI control surface is richer and more fully tested.

#### foundation/pi/

| Feature | hyperharness | superai |
|---|---|---|
| `runtime.go` | ✅ | ✅ |
| `session.go` | ✅ | ✅ |
| `summary.go` | ✅ | ✅ |
| `tools_native.go` | ✅ | ✅ |
| `tools_extra.go` | **NO** (removed during port) | **YES** |
| `foundation.go` | ✅ | ✅ |
| Deterministic generators | ✅ | ✅ |

**Winner: superai** — includes `tools_extra.go` which hyperharness port doesn't (yet).

#### agent/

| Feature | hyperharness | superai |
|---|---|---|
| `agent.go` | ✅ | ✅ |
| `orchestrator.go` | ✅ | ✅ |
| `pipe.go` | ✅ | ✅ |
| `async.go` | NO | ✅ |
| `autocomplete.go` | NO | ✅ |
| `autopilot.go` | NO | ✅ |
| `compare.go` | NO | ✅ |
| `context.go` | NO | ✅ |
| `diff.go` | NO | ✅ |
| `oracle.go` | NO | ✅ |
| `shell.go` | NO | ✅ |
| `agent_test.go` | NO | ✅ |
| `orchestrator_test.go` | NO | ✅ |

**Winner: superai** — significantly more agent-side functionality and tests. The hyperharness agent
subset is a minimal viable port; superai has deep agent behavior (async orchestration, diff
generation, context management, autonomous shell execution, etc.).

#### agents/ (director sub-pkg)

| Feature | hyperharness | superai |
|---|---|---|
| `director.go` | ✅ (basic) | ✅ (basic) |
| `director_test.go` | NO | ✅ |
| `provider.go` | ✅ (stub) | ✅ (stub) |
| `provider_stub.go` | ✅ | ✅ |
| `provider_hypercode.go` | NO | ✅ |
| `auton.go`, `codemode.go`, `council.go`, etc. | NO | ✅ |
| `rag.go`, `shell_assistant.go` | NO | ✅ |

**Winner: superai** — the agents package in superai has extensive functionality including
autonomous agent orchestration, multi-agent council, RAG pipeline, and shell assistant.
hyperharness has only a minimal stub director + provider.

#### internal/* (hyperharness only)

| Package | Purpose | Present in superai? |
|---|---|---|
| `internal/agent` | Agent runtime wrapper | NO (separate `agent/` pkg) |
| `internal/config` | Configuration management | NO (uses `config/` only) |
| `internal/mcp` | MCP client layer | NO (uses `mcp/` only) |
| `internal/memory` | Knowledge memory service | NO |
| `internal/providers` | Provider registry | NO |
| `internal/sessions` | Session management | NO |
| `internal/tools` | Tool registry | NO (uses `tools/` only) |
| `internal/ui` | UI components | NO |

**Winner: hyperharness** — the `internal/*` packages provide architectural separation between
internal plumbing and public-facing packages, which is a Go best-practice pattern superai
does not use.

---

## Deep-Dive: Functionality

### TUI Capability Matrix

| Capability | hyperharness | superai |
|---|---|---|
| Chat with AI | ✅ | ✅ |
| Slash commands | 45+ | 22 |
| Tree browser (modal) | ✅ | ✅ |
| Persistent pane | ✅ | ✅ |
| Pane focus | ✅ (on/off/toggle) | minimal |
| Pane sizing | ✅ (set, cycle, presets) | minimal |
| Pane position | ✅ (top/bottom/toggle) | minimal |
| Pane preview | ✅ (on/off/toggle) | minimal |
| Pane grouping | ✅ (on/off/toggle) | minimal |
| Pane presets | ✅ (compact/detailed/nav/review) | minimal |
| Pane refresh | ✅ (auto + manual) | manual only |
| Pane summary | ✅ (compact one-line) | NO |
| Pane help | ✅ (subsystem-local) | NO |
| Pane reset | ✅ | ✅ |
| Pane status | ✅ (detailed + summary) | minimal |
| Fold/collapse | ✅ | ✅ |
| Graph connectors | ✅ (`├─`, `└─`) | ✅ |
| Filter | ✅ | ✅ |
| Confirm-before-switch | ✅ | NO |
| Pre-switch summary preview | ✅ | NO |
| Position divider | ✅ | NO |

### Agent Capability Matrix

| Capability | hyperharness | superai |
|---|---|---|
| Basic agent orchestration | ✅ | ✅ |
| Async task dispatch | ❌ | ✅ |
| Autocomplete suggestions | ❌ | ✅ |
| Autopilot mode | ❌ | ✅ |
| Context management | ❌ | ✅ |
| Diff generation & application | ❌ | ✅ |
| Code comparison | ❌ | ✅ |
| Oracle (prediction/analysis) | ❌ | ✅ |
| Autonomous shell | ❌ | ✅ |
| Agent tests | ❌ | ✅ (2 files) |

### Foundation Core Capability Matrix

| Capability | hyperharness | superai |
|---|---|---|
| Runtime (session, budget, budget-aware branching) | ✅ | ✅ |
| Session store (JSONL, tree, branch) | ✅ | ✅ |
| Summary generation hooks | ✅ | ✅ |
| Compaction generation hooks | ✅ | ✅ |
| Native tools (read/write/edit/bash/grep/find/ls) | ✅ | ✅ |
| Extra tools | ❌ | ✅ (tools_extra.go) |
| Tool parity tests | ✅ | ✅ |

### Testing Discipline

| Metric | hyperharness | superai |
|---|---|---|
| Test files | 32 | 733 |
| Test pass (`./tui ./cmd ./foundation/...`) | ✅ green | ✅ green |
| TUI slash command tests | 27+ test funcs | fewer focused |
| Foundation tests | comprehensive per-package | comprehensive per-package |
| Agent tests | **none** (package not yet ported) | 2 files (agent, orchestrator) |
| Agent subpackage tests (agents/) | **none** | director_test.go |

### Documentation

| Metric | hyperharness | superai |
|---|---|---|
| README | comprehensive | minimal (just `# superai`) |
| Analysis docs | 48 tranches in `docs/analysis/` | 0 |
| HANDOFF.md | continuously updated | present, minimal |
| Architecture doc | `docs/architecture/HYPERHARNESS_ARCHITECTURE.md` | none |

---

## Overall Ratings

| Criterion | hyperharness | superai |
|---|---|---|
| **TUI surface** | **A+** (45+ commands, full pane/browser) | **B** (basic tree browser, minimal pane) |
| **Agent depth** | **C** (minimal port) | **A+** (async, autopilot, diff, oracle, shell) |
| **Foundation core** | **A** (deterministic, well-tested) | **A+** (plus extra tools) |
| **Test coverage breadth** | **B** (32 files, focused) | **A+** (733 files, extensive) |
| **Test pass reliability** | **A** (all green, fast) | **A** (all green, fast) |
| **Documentation** | **A+** (48 analysis docs + handoff) | **D** (empty README) |
| **Module identity** | **A+** (canonical `hyperharness`) | **B** (legacy `hypercode` name) |
| **Go idiomatic structure** | **A** (uses `internal/` pattern) | **B+** (flat package structure) |
| **Overall architectural polish** | **A** | **A-** |
| **Overall feature breadth** | **B+** | **A+** |

## Synthesis

**hyperharness** is the more **polished, documentable, and TUI-capable** port with a cleaner Go
module identity, better idiomatic structure (`internal/` packages), and an exhaustive tranche-by-tranche
development trail. It excels at:
- TUI pane/browser ergonomics
- Documentation density
- Module path correctness
- Test-focused confidence

**superai** is the more **functionally complete** sibling with significantly more agent-level
autonomy, extensive test coverage, and extra tool capabilities. It excels at:
- Agent-side depth (async, autopilot, diff, oracle, shell)
- Test breadth (733 test files)
- Foundation extras (`tools_extra.go`)
- Agents subpackage (auton, council, rag, etc.)

**The highest-value next step** is a bidirectional synthesis:
1. **Port superai's agent depth into hyperharness** (async, autocomplete, autopilot, diff, oracle, shell)
2. **Port superai's `tools_extra.go` into hyperharness** foundation/pi
3. **Normalize superai's module path** to `hyperharness`
4. **Apply hyperharness's TUI richness** onto superai if desired
