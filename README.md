# HyperHarness (HyperCode)

HyperHarness is a Go control plane for AI development harnesses: the goal is a single local-first substrate that exposes the tool surface of the popular AI coding harnesses (Claude Code, Codex, Aider, and 15+ others) through one binary — with a TUI dashboard, session continuity, MCP tooling, and knowledge-base ingest. The binary is called `hypercode`.

Version **0.5.0-alpha.5** (see `VERSION`) — alpha; interfaces are still moving.

## Status

**Active, alpha.** Last significant activity: **2026-11-12** — merge of a Jules feature branch (Protocol #281). Nov 8–9 work was infrastructure sanitization (web manifests, build scripts, docs) and a deployment-guide update for v0.5.0-alpha.5 and dual binaries.

## A note on what's in this repo

GitHub reports this repo's language as TypeScript, which is misleading: the project's own code is **Go** (module `github.com/robertpelloni/hyperharness`, `go.mod` at the root, `cmd/`, `cli/`, `agent/`, `foundation/`). The bulk of the tree is **34 git submodules / vendored reference codebases** — the harness implementations under study for tool-parity work (directories named after the harnesses: `claude-code`, `aider`, `crush`, `copilot-cli`, `gemini-cli`, …) — plus bot-workflow state files (`.borg_queue.db`, session markers) from the development process.

Compiled binaries (`hypercode`, `hyperharness`, `hypercode.exe~`, ~30 MB each) are checked into the repo root — build artifacts that should not be tracked.

A near-duplicate of this repo exists under the candlestixxx account (`candlestixxx/hyperharness`, a fork); this is the primary copy.

## Build & run

Prerequisite: Go (the module targets Go 1.26; go.mod is authoritative). `build.sh` / `build.bat` at the root wrap the build.

```bash
git clone https://github.com/robertpelloni/hyperharness.git
cd hyperharness
go build ./...
./hypercode --help
./hypercode tui              # TUI dashboard
./hypercode ingest ./src     # ingest a tree into the knowledge base
```

There is also a `Dockerfile` / `Dockerfile.prod` and `docker-compose.yml`.

## Docs

`VISION.md`, `DEPLOY.md`, `docs/PROJECT_STRUCTURE.md`, `docs/analysis/` (harness-by-harness analysis), plus the working notes (`ANALYSIS.md`, `ANALYSIS_GO_PORT.md`, `PORTING_MAP.md`, `AUDIT.md`, `TODO.md`, `HANDOFF.md`) — development artifacts of varying age; where they conflict with the code, believe the code.

## License

MIT (per the project's stated license).
