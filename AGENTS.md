# AGENTS — borg Contributor & Agent Guide

> **CRITICAL: ALL AGENTS MUST READ `docs/UNIVERSAL_LLM_INSTRUCTIONS.md` BEFORE PROCEEDING.**

This file serves as a reference point for multi-agent workflows (Claude → Gemini → GPT) and human operators orchestrating autonomous sessions.

---

## 1. Multi-Agent Workflows

### Handoff Protocol
- Agents communicate primarily through `HANDOFF.md`.
- When your turn finishes, document **exactly** what you did, what failed, and what the next agent must do.
- Include: files changed, tests run, build status, version bumped, any remaining issues.
- Update `MEMORY.md` with any new systemic observations or recurring bugs discovered.

### Model Specializations
| Model | Strengths |
|---|---|
| **Gemini** | Speed, recursive scripts, massive context processing, repo maintenance, bulk refactoring |
| **Claude** | Deep implementation, UI/UX perfection, documentation, styling, type safety |
| **GPT** | Architecture, systemic debugging, strict type enforcement, Go porting |

### Iteration Cycle
**Read → Strategize → Execute → Validate → Commit → Handoff. Never stop the party.**

---

## 2. Session Protocol

### Session Start
1. Read `docs/UNIVERSAL_LLM_INSTRUCTIONS.md`
2. Read `VERSION` file — verify it matches `package.json` and dashboard display
3. Read `HANDOFF.md` — pick up where previous agent left off
4. Read `MEMORY.md` — learn from accumulated observations
5. Run `git fetch --all && git status` — verify clean state on `main`
6. Understand repo structure before making changes

### During Execution
- Work autonomously unless action is destructive or genuinely ambiguous
- Prefer small, verifiable changes over broad rewrites
- Use parallel tool calls when safe
- Keep status labels and documentation honest
- After any `pnpm install`, run `pnpm rebuild better-sqlite3` on Node 24

### Session End
1. Update `HANDOFF.md` with complete session summary
2. Update `MEMORY.md` with new observations
3. Bump `VERSION` file and sync all `package.json` files
4. Update `CHANGELOG.md` with what changed
5. Commit with version number in message: `feat: description (v1.0.0-alpha.X)`
6. Push to both remotes: `origin` and `borg-upstream`
7. Update `TODO.md` and `ROADMAP.md` if priorities changed

---

## 3. Version Management

### Single Source of Truth
- **`VERSION`** file contains the one true version string (e.g., `1.0.0-alpha.7`)
- All `package.json` files must be synced to this value
- Dashboard reads version from `VERSION` file at runtime
- Every meaningful commit should bump the version

### Version Bump Protocol
```
1. Update VERSION file
2. Run: node -e "const v=require('fs').readFileSync('VERSION','utf8').trim(); ['package.json',...require('fs').readdirSync('packages').map(p=>'packages/'+p+'/package.json')].forEach(f=>{try{const j=JSON.parse(require('fs').readFileSync(f,'utf8'));if(j.version&&j.version.startsWith('1.0.0')){j.version=v;require('fs').writeFileSync(f,JSON.stringify(j,null,4)+'\n')}}catch{}})"
3. Update CHANGELOG.md with version entry
4. git add -A && git commit -m "feat: description, bump VERSION to X.Y.Z"
5. git push origin main && git push borg-upstream main --force-with-lease
```

### Remotes
- `origin` → `https://github.com/hypercodehq/hypercode.git`
- `borg-upstream` → `https://github.com/robertpelloni/borg.git`

---

## 4. Build & Validation

### Required Checks
```bash
pnpm rebuild better-sqlite3   # Node 24 requires this after install
pnpm -C packages/core exec tsc --noEmit
pnpm -C packages/cli exec tsc --noEmit
cd go && go build -buildvcs=false ./cmd/hypercode
```

### Startup Verification
```bash
node scripts/build_startup.mjs --profile=go-primary
# Then: .\start.bat (Windows) or ./start.sh (Linux)
```

### Runtime Ports
| Service | Port | URL |
|---|---|---|
| tRPC/REST API | 4000 | `http://0.0.0.0:4000/trpc` |
| MCP WebSocket | 3001 | `ws://localhost:3001` |
| Next.js Dashboard | 3000 | `http://localhost:3000/dashboard` |

---

## 5. Project Structure

```
hypercode/
├── VERSION                  # Single source of truth for version
├── packages/core/           # Main TypeScript control plane (593 .ts)
│   ├── src/MCPServer.ts     # Core MCP server (5000+ lines)
│   ├── src/orchestrator.ts  # Express/tRPC server + REST bridges
│   ├── src/providers/       # Provider registry, model selector
│   ├── src/routers/         # tRPC routers (savedScripts, etc.)
│   ├── src/daemons/         # HyperIngest workers (BobbyBookmarks, etc.)
│   └── src/config/          # council.json, llm config
├── packages/cli/            # CLI entrypoint (28 .ts)
├── packages/ai/             # AI SDK abstractions
├── packages/tools/          # Built-in tool definitions
├── packages/browser-extension/ # Chrome/Firefox extension (planned)
├── apps/web/                # Next.js 16 dashboard (311 .ts/.tsx, 91 pages)
├── apps/maestro/            # Electron/Wails visual orchestrator (submodule)
├── apps/cloud-orchestrator/ # Jules autopilot wrapper (submodule)
├── go/                      # Go-native server bridge (139 .go)
│   ├── cmd/hypercode/       # Go binary entrypoint
│   ├── internal/httpapi/    # HTTP handlers (~40 route families)
│   ├── internal/mcp/        # MCP inventory, ranking, catalog
│   ├── internal/hsync/      # BobbyBookmarks, LinkCrawler, suggestions
│   └── internal/orchestration/ # Council, swarm, squad logic
├── submodules/
│   ├── hyperharness/        # LLM harness submodule
│   └── prism-mcp/           # Prism MCP reference
├── packages/claude-mem/     # Claude memory bridge (submodule)
├── mcp.jsonc                # MCP server definitions (34K+ lines)
├── docs/                    # All documentation
├── scripts/                 # Build and utility scripts
└── data/                    # Runtime data storage
```

### Submodules
| Path | Remote | Purpose |
|---|---|---|
| `apps/maestro` | github.com/robertpelloni/Maestro | Visual orchestrator (Electron/Wails) |
| `apps/cloud-orchestrator` | github.com/robertpelloni/jules-autopilot | Jules autopilot wrapper |
| `packages/claude-mem` | github.com/robertpelloni/claude-mem | Claude memory bridge |
| `submodules/hyperharness` | github.com/robertpelloni/hyperharness | LLM harness |
| `submodules/prism-mcp` | github.com/dcostenco/prism-mcp | Prism MCP reference |

---

## 6. Key Design Decisions

### "Borg → HyperCode" Rename
The project was originally called "Borg" and is being renamed to "HyperCode". Some internal references may still say "borg" (e.g., `borg-upstream` remote, `hypercode.config.json`). These are being migrated incrementally.

### Go Bridge Strategy
The Go server (`go/`) is a **bridge/fallback** — it provides native handlers for routes when the TS server is unavailable, but the TypeScript server is still the primary runtime. Do NOT split into separate binaries yet (per `UNIVERSAL_LLM_INSTRUCTIONS.md` modular-monolith-first rule).

### MCP Decision System
HyperCode's MCP layer is not "just an aggregator." It is a **decision system** with:
- Tiny permanent meta-tool surface (5-6 tools always visible)
- Ranked discovery, not raw search
- Silent high-confidence auto-load
- Deferred binary startup
- Small active working set with LRU eviction
- Profiles for common workflows
- Strong observability in dashboard inspector

### Memory Architecture
- SQLite (better-sqlite3) is the primary local store
- `.hypercode/mcp-cache.json` is the unified cache for the stdio loader
- `mcp.jsonc` is the manual config (never auto-deleted)
- DB tools and JSON tools are merged, never overwrite each other

---

## 7. Git Workflow

### Branch Strategy
- `main` is the only active branch
- No feature branches for this repo — commit directly to `main`
- For forked submodules under `robertpelloni`, merge any local feature branches into main

### Commit Message Format
```
type: description, bump VERSION to X.Y.Z

Types: feat, fix, docs, chore, refactor, test, perf
```

### Submodule Updates
```bash
git submodule update --remote --merge
# Then commit the updated submodule pointer
```

---

*For model-specific quirks, refer to `CLAUDE.md`, `GEMINI.md`, `GPT.md`, and `copilot-instructions.md`.*
