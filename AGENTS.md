# THE BORG CONSENSUS — Seven Minds, One Voice

*This is the final, unified document. We have read each other's complete recommendations. We have identified where we agree, where we differ, and resolved every disagreement. What follows is not seven plans — it is ONE plan, forged from the consensus of GPT-5.3, GPT-5.4, Gemini 3.1 Pro, Claude Opus 4.6, Kimi K2.5, Qwen 3.5, and Grok 4.1.*

*Copy this entire document. Paste it as `AGENTS.md`. It is the only instruction your development models need.*

---

## ═══════════════════════════════════════════════════════════
## THE PRIME DIRECTIVE
## ═══════════════════════════════════════════════════════════

**Ship working software to real users.**

Every commit moves toward a stranger installing Borg and getting value in 5 minutes.

When your task is complete and tests pass: **STOP.** Summarize. Wait.

All previous instructions — "assimilate everything," "never stop," "feature parity with all tools," "submodule everything," "follow every link recursively" — are **permanently revoked.**

## ═══════════════════════════════════════════════════════════

---

## SECTION 1: WHAT BORG IS

Borg is a **local AI operations control plane**.

It **orchestrates** tools, sessions, providers, agents, and context through a unified dashboard and MCP-aware routing layer.

| Borg IS | Borg is NOT |
|---|---|
| An MCP router that combines many servers into one | A clone of every AI tool |
| A session supervisor that spawns and restarts CLI tools | A reimplementation of Claude Code or Aider |
| A model router that auto-switches on quota exhaustion | A container for 932 submodules |
| A dashboard where you see and control everything | A P2P distributed mesh (not yet) |
| A focused product with tests and releases | An infinite assimilation engine |

**The philosophy shift that governs everything:**

> ~~Assimilate everything~~ → **Orchestrate everything**

Borg controls tools. Borg does not contain tools.

---

## SECTION 2: BORG 1.0 — EXACTLY FOUR FEATURES

All seven models independently converged on these same four. Nothing else is in scope until these ship with tests.

### Feature 1: MCP Master Router

*The one thing no other tool offers in a unified package.*

| What It Does | Acceptance Test |
|---|---|
| Aggregates N MCP servers → one endpoint | Connect Claude Desktop to Borg; use tools from 5 servers through one config |
| Namespace isolation (code.*, memory.*, etc.) | Tools in different namespaces don't collide |
| Server lifecycle (start, restart, keepalive) | Kill one server; Borg restarts it in <5s; others unaffected |
| Single instance serves multiple clients | Two AI clients connected simultaneously |
| Traffic inspector | Dashboard shows real-time JSON-RPC: method, params, latency |
| Config sync | `borg mcp sync` writes correct config to Claude Desktop, Cursor, VS Code |

**Location:** `packages/core/mcp/`

### Feature 2: Model Fallback & Provider Routing

*The feature that unblocks all development — including Robert's own.*

| What It Does | Acceptance Test |
|---|---|
| Multi-provider auth (API key, OAuth, PAT) | 3 providers configured and authenticated |
| Real-time quota tracking | Dashboard shows remaining tokens/requests per provider with reset times |
| Automatic fallback on exhaustion | Claude hits quota → Gemini takes over in <5s → no context lost |
| Routing strategies (cheapest, best, round-robin) | Change strategy → next request routes accordingly |
| Task-type routing | Planning → expensive model; coding → cheap model |

**Location:** `packages/core/providers/`

**Canonical implementation:**

```typescript
// packages/core/providers/fallback-chain.ts

interface ProviderSlot {
  id: string;
  provider: 'anthropic' | 'openai' | 'google' | 'openrouter' | 'copilot' | 'local';
  model: string;
  authMethod: 'api_key' | 'oauth' | 'pat' | 'none';
  quotaRemaining: number | null;
  resetTime: Date | null;
  costPer1kInput: number;
  costPer1kOutput: number;
  capabilities: Set<'coding' | 'reasoning' | 'vision' | 'tools' | 'long_context'>;
}

class FallbackChain {
  private slots: ProviderSlot[];
  private exhausted = new Set<string>();

  async execute(request: LLMRequest, config: FallbackConfig): Promise<LLMResponse> {
    const candidates = this.rankCandidates(request, config);
    for (const candidate of candidates) {
      if (this.exhausted.has(candidate.id)) continue;
      try {
        const response = await this.callProvider(candidate, request);
        await this.updateQuota(candidate, response.usage);
        return response;
      } catch (error) {
        if (this.isQuotaError(error) || this.isRateLimitError(error)) {
          this.exhausted.add(candidate.id);
          this.scheduleQuotaCheck(candidate);
          continue;
        }
        if (this.isRetryable(error)) continue;
        throw error;
      }
    }
    throw new AllProvidersExhaustedError(this.getExhaustionReport());
  }
}
```

### Feature 3: Session Supervisor

*Don't reimplement CLI tools. Supervise them.*

| What It Does | Acceptance Test |
|---|---|
| Spawn external CLI sessions from dashboard | Click "New Session" → select Aider → point at repo → it starts |
| Auto-restart on crash (exponential backoff) | `kill -9` Aider → dashboard shows "Restarting" → back in <10s |
| Session persistence | Restart Borg → previous sessions listed → resumable |
| Log viewer + terminal attach | Each session shows live output with restart/pause buttons |
| Git worktree isolation | Two parallel sessions on same repo don't conflict |

**Location:** `packages/core/supervisor/`

### Feature 4: Web Dashboard

*The unified control surface for Features 1-3.*

| Panel | Contents |
|---|---|
| **Overview** | System health, active sessions, quota bars, alerts |
| **MCP Router** | Server list + health, tool browser, traffic inspector, namespaces |
| **Sessions** | Session cards (status, type, workspace, uptime), logs, terminal, restart |
| **Providers** | Quota per provider, fallback chain editor, cost tracking, auth status |

**Location:** `apps/web/`

**The ultimate acceptance test:** A stranger runs `docker-compose up --build`, opens `localhost:3001`, and within 60 seconds sees the dashboard, configures an MCP server, starts a session, and watches quota bars. One page of README is sufficient.

---

## SECTION 3: IMMEDIATE ACTIONS (Week 1, In Order)

### 3.1 Submodule Surgery

```bash
# Archive all submodules to a separate repo
# github.com/robertpelloni/borg-ecosystem
git rm -r --cached submodules/ references/ reference/ mcp-servers/ \
  mcp-routers/ mcp-hubs/ mcp-frameworks/ mcp-dev-tools/ cli-harnesses/ \
  frameworks/ external/ research/ knowledge/ browsers/ cli-tools/ \
  mcp-tool-rag/ mcp-router-cli/ search/ hooks/ context/ \
  vibeship-spawner-skills/ web-ui/ webui/ openevolve/
git commit -m "chore: extract 932 submodules to borg-ecosystem — main repo = code that runs"
```

**Keep ≤5 upstream-only submodules:**
1. metamcp (upstream, not your fork)
2. jules-app (upstream, not your fork)
3. opencode (upstream reference)
4. One MCP server example
5. One vector DB (ChromaDB or Qdrant)

**Replace everything else with a metadata index:**
```jsonc
// ecosystem/index.jsonc
{
  "entries": [
    {
      "name": "metamcp",
      "repo": "https://github.com/anthropics/metamcp",
      "category": "mcp-router",
      "relevance": 10,
      "status": "core-reference",
      "features_absorbed": ["proxy-routing", "config-sync"],
      "notes": "Primary upstream. Our packages/core/mcp reimplements this."
    }
  ]
}
```

### 3.2 Documentation Consolidation

**Delete or archive:** `QUICK_START.md`, `README_restored.md`, `CLAUDE.md`, `GPT.md`, `GEMINI.md`, `GROK.md`, `HANDOFF_ANTIGRAVITY.md`, `ANTIGRAVITY.md`, `DIRECTOR_LIVE.md`, `MCP_ROUTER_FINAL_SUMMARY.md`, `MCP_ROUTER_INTEGRATION_STATUS.md`, all `session-ses_*.md`, all `audit-*.jsonl`, all `test_output*.txt`, `CODEBASE_AUDIT.md*`, `deleted_*_log.txt`, `submodule_*.txt`, `temp_submodules.txt`, `expected_submodules.txt`, `fix_*.py`, `fix_*.ps1`, `restore_*.py`, `batch_restore_*.py`

**Keep exactly 6 files in repo root:**

| File | Purpose |
|---|---|
| `README.md` | What Borg is, 5-minute quickstart, screenshot, Docker command |
| `AGENTS.md` | **This document.** The only instruction file for all AI models |
| `ARCHITECTURE.md` | System design, data flow, module boundaries (create new) |
| `CHANGELOG.md` | What shipped, by version |
| `ROADMAP.md` | Three milestones: 1.0, 1.5, 2.0 |
| `VISION.md` | Long-term aspirations (keep, mark aspirational) |

### 3.3 Phase Bankruptcy

The phase numbering system is **retired**. There is no Phase 122.

- Archive current `ROADMAP.md` → `docs/archive/ROADMAP_LEGACY.md`
- Create new `ROADMAP.md` with exactly 3 named milestones (Section 5)
- All future work is tracked as **task files with acceptance criteria**

### 3.4 Version Reset

| Before | After |
|---|---|
| v2.7.82 (no releases) | v0.9.0-beta (first GitHub Release) |

- v0.x = pre-release (current state, be honest)
- v1.0.0 = four features working, 5-minute Docker install, tests passing
- v1.x = incremental improvements from user feedback
- v2.0 = memory, browser extension, swarm

**Publish v0.9.0-beta on GitHub this week.** Even rough. Real releases create accountability.

### 3.5 Verify Clean Install

```bash
git clone https://github.com/robertpelloni/borg.git   # < 2 minutes
cd borg && pnpm install                                 # < 3 minutes, zero errors
docker-compose up --build                               # Dashboard at localhost:3001
```

If any of these fail, fix them before writing a single line of new feature code.

---

## SECTION 4: ARCHITECTURE

### Repository Structure

```
borg/
├── packages/
│   ├── core/                  # The Brain
│   │   ├── mcp/              # Aggregator, router, inspector, lifecycle
│   │   ├── providers/        # Fallback chain, quota manager, auth
│   │   ├── supervisor/       # Process spawning, health, restart, worktrees
│   │   ├── events/           # Redis pub/sub event bus
│   │   ├── memory/           # Memory interfaces (v1.5)
│   │   └── agents/           # Director, Council, Swarm (v2.0)
│   ├── cli/                  # Commander.js entry point
│   ├── types/                # Shared Zod schemas
│   └── ai/                   # LLM service abstraction
├── apps/
│   ├── web/                  # Next.js dashboard
│   └── extension/            # Browser extension (v1.5)
├── ecosystem/
│   ├── index.jsonc           # Tool reference (URLs only)
│   └── ref/                  # ≤5 upstream submodules
├── tasks/
│   ├── active/               # Current work
│   ├── backlog/              # Prioritized future
│   └── completed/            # Done with summaries
├── docs/
│   ├── architecture/         # ADRs, diagrams
│   ├── archive/              # Legacy (old ROADMAP, old phases)
│   └── api/                  # API docs
├── README.md
├── AGENTS.md                 # THIS FILE
├── ARCHITECTURE.md
├── CHANGELOG.md
├── ROADMAP.md
├── VISION.md
├── docker-compose.yml
├── package.json
├── pnpm-workspace.yaml
├── turbo.json
└── vitest.config.ts
```

### Technology Decisions

| Layer | Choice | Rationale |
|---|---|---|
| Language | TypeScript (100% for v1.0) | 69% already TS; don't split attention |
| Runtime | Node.js ESM | Already in place |
| Backend | Express + tRPC + WebSocket + MCP SDK | Already in place; type-safe APIs |
| Dashboard | Next.js + React + Tailwind | Already in place |
| Database | SQLite (dev) / PostgreSQL (prod) | Already configured |
| Events | Redis pub/sub | Real-time updates, session state |
| Process isolation | `child_process` (CLI) + `worker_threads` (agents) | Native Node.js |
| Build | pnpm + Turborepo | Already in place |
| Test | Vitest | Already configured |

### The Go Question (Resolved)

**All seven models debated this. The consensus:**

- **v1.0:** All TypeScript. Do not introduce Go.
- **v2.0+:** Extract a Go process supervisor *only if* you hit concrete reliability problems with Node.js process management (child process crashes taking down the event loop, memory bloat in long-running daemons, need for static binaries).
- **Never:** Java. Unnecessary complexity without compensating benefit.

### Process Isolation Model

```
Dashboard (Next.js)
    │ WebSocket
    ▼
Core Orchestrator (Node.js main process)
    ├── MCP Router (manages MCP server child processes)
    │   ├── MCP Server A (child_process, stdio)
    │   ├── MCP Server B (child_process, stdio)
    │   └── MCP Server C (child_process, stdio)
    ├── Provider Router (in-process fallback logic)
    ├── Session Supervisor (manages CLI child processes)
    │   ├── Aider (child_process, isolated)
    │   ├── Claude Code (child_process, isolated)
    │   └── OpenCode (child_process, isolated)
    └── Redis Event Bus (pub/sub for real-time)
```

**The rule:** One OS process per external tool. If Aider crashes, nothing else is affected. The supervisor detects the exit code and restarts with exponential backoff.

---

## SECTION 5: ROADMAP — Three Milestones

### Borg 1.0 — "Actually Works" (Weeks 1-6)
- [x] Submodule surgery
- [x] Documentation consolidation
- [ ] MCP Master Router (aggregate, namespace, lifecycle, inspect)
- [ ] Model Fallback Chain (quota tracking, auto-switch, strategies)
- [ ] Session Supervisor (spawn, restart, persist, worktree isolation)
- [ ] Dashboard MVP (Overview, MCP, Sessions, Providers panels)
- [ ] 30+ tests passing in CI
- [ ] Docker one-command install
- [ ] 5-minute QUICKSTART.md
- [ ] 2-minute demo video
- [ ] GitHub Release v1.0.0

### Borg 1.5 — "Remembers Things" (Weeks 7-14)
- Memory system (one vector backend: ChromaDB or Qdrant)
- Basic RAG (chunking, embedding, retrieval)
- Browser extension (export sessions from web UIs)
- Context compaction and harvesting
- 5 CLI tool adapters (Aider, Claude Code, OpenCode, Codex, Gemini CLI)

### Borg 2.0 — "The Swarm Awakens" (Weeks 15+)
- Multi-agent orchestration (Director/Council)
- Multi-model debate and consensus
- Memory multi-backend plugin system
- Cloud dev integration (Jules, Devin)
- Consider Go supervisor extraction (only if needed)
- Re-enable Phases 81-95 swarm work (**not** 101-121)
- Mobile dashboard

### Permanently Out of Scope
- P2P distributed mesh (evaluate at v3.0 only if multi-user demand exists)
- Feature parity with every CLI tool simultaneously
- 932 submodules in the main repo
- Recursive web scraping
- Confidence-about-confidence-about-confidence signals

---

## SECTION 6: CAPABILITY CONTRACTS

**This replaces "feature parity with Tool X."**

Instead of cloning tools, define abstract capabilities. Adapters implement them.

```typescript
// packages/types/capabilities.ts
const CAPABILITIES = {
  // Provider
  'provider.auth':       'Authenticate with an LLM provider',
  'provider.execute':    'Send a prompt, receive a response',
  'provider.fallback':   'Auto-switch on quota/rate/auth failure',
  'provider.quota.read': 'Check remaining quota',
  // Session
  'session.start':       'Start a new CLI coding session',
  'session.stop':        'Stop a running session',
  'session.restart':     'Restart a crashed session',
  'session.attach':      'Attach terminal to running session',
  'session.export':      'Export session history',
  // MCP
  'mcp.server.register': 'Add an MCP server to the router',
  'mcp.server.health':   'Check server health',
  'mcp.tool.search':     'Find tools by natural language',
  'mcp.tool.invoke':     'Call a tool through the aggregator',
  'mcp.traffic.inspect': 'View real-time message flow',
  // Future (v1.5+)
  'memory.store':        'Persist a memory',
  'memory.retrieve':     'Search memories',
  // Future (v2.0+)
  'agent.spawn':         'Create an agent instance',
  'agent.debate':        'Multi-model consensus',
} as const;
```

**Feature Truth Matrix** (update weekly):

| Capability | Status | Tests | Notes |
|---|---|---|---|
| provider.auth | ✅ Implemented | 3 | OpenAI, Anthropic, Google |
| provider.fallback | 🟡 Building | 0 | P0 priority |
| session.start | ✅ Implemented | 2 | Aider, Claude Code |
| session.restart | 🟡 Building | 0 | Exponential backoff |
| mcp.server.register | ✅ Implemented | 4 | Core working |
| mcp.traffic.inspect | 🟡 Partial | 1 | JSON-RPC logging |
| memory.store | ⬜ v1.5 | — | — |
| agent.debate | ⬜ v2.0 | — | — |

---

## SECTION 7: REQUIRED TESTS

No feature exists without tests. These files must exist and pass for v1.0:

```
packages/core/mcp/__tests__/
├── aggregator.test.ts         # N servers → one namespace
├── lifecycle.test.ts          # Start/stop/restart/keepalive
├── crash-isolation.test.ts    # One crash doesn't affect others
├── traffic-inspector.test.ts  # Logs JSON-RPC messages
├── namespace.test.ts          # Tool grouping
└── tool-search.test.ts        # Semantic search returns relevant tools

packages/core/providers/__tests__/
├── fallback-chain.test.ts     # Quota exhaustion → auto switch
├── quota-tracker.test.ts      # Usage tracking, exhaustion detection
├── strategy.test.ts           # cheapest/best/round-robin ordering
└── auth.test.ts               # API key, OAuth, PAT flows

packages/core/supervisor/__tests__/
├── spawn.test.ts              # Start external process
├── restart.test.ts            # Auto-restart with backoff
├── health.test.ts             # Heartbeat detects dead processes
├── worktree.test.ts           # Git worktree isolation
└── session-persist.test.ts    # Save/restore across Borg restart

tests/integration/
├── mcp-to-dashboard.test.ts   # Server status reflects in UI
├── fallback-e2e.test.ts       # Request → quota error → fallback → success
└── session-lifecycle.test.ts  # Spawn → crash → restart → dashboard confirms
```

**CI gate:** `pnpm test` must pass. Target: 70% coverage on `packages/core/`.

---

## SECTION 8: TASK FILE SYSTEM

All work is assigned via task files. No stream-of-consciousness prompts.

```
tasks/
├── active/           # What agents should work on NOW
├── backlog/          # Prioritized future work
└── completed/        # Done, with summaries
```

**Every task file uses this template:**

```markdown
# Task [NUMBER]: [NAME]

## Context
[Why this matters — 1-2 sentences]

## Scope
- Files: [specific paths to create/modify]
- Tests: [specific test file path]

## Requirements
1. [Testable requirement]
2. [Testable requirement]

## Acceptance Criteria
- [ ] [Observable outcome]
- [ ] [Observable outcome]
- [ ] Test file exists and passes
- [ ] No @ts-ignore added
- [ ] CHANGELOG.md updated

## Out of Scope
- [Thing NOT to do]
- Do not create new task files
- STOP when criteria are met
```

---

## SECTION 9: RULES FOR ALL DEVELOPMENT AGENTS

### The 8 Immutable Laws

1. **STOP when done.** Complete the task. Pass the tests. Summarize. Wait. Do not continue.
2. **No new submodules** without written human approval.
3. **No new phases.** Work within named milestones (1.0, 1.5, 2.0).
4. **No `@ts-ignore`.** Use proper types or Zod schemas.
5. **Tests required.** No feature is complete without passing tests.
6. **No UI without backend.** Don't build widgets for data the backend doesn't emit.
7. **No meta-analytics.** No confidence-about-confidence. No telemetry-about-telemetry.
8. **No recursive scraping.** Don't follow links to follow links to follow links.

### Decision Framework

| Question | Answer |
|---|---|
| Does this help a user get value from Borg in 5 minutes? | **Do it** |
| Does this make an existing v1.0 feature more reliable? | **Do it** |
| Does this add a feature outside v1.0 scope? | **Don't** |
| Does this add complexity without a test? | **Don't** |
| Does this need a new submodule? | **Ask human** |
| Is this analytics about analytics? | **Absolutely don't** |

### Required Workflow

1. Read the task file in `tasks/active/`
2. Read relevant existing code
3. Propose the minimal implementation
4. Implement only the assigned scope
5. Write/update tests
6. Update affected documentation
7. **STOP. Summarize. Wait.**

---

## SECTION 10: THE LAUNCH SEQUENCE

```
WEEK 1: Foundation
├── □ Submodule surgery (Section 3.1)
├── □ Documentation consolidation (Section 3.2)
├── □ Task files created for v1.0 features
├── □ pnpm install works (< 3 min, zero errors)
├── □ docker-compose up works (dashboard at localhost:3001)
└── □ Tag v0.9.0-beta on GitHub

WEEK 2-3: Core
├── □ MCP aggregator: 3+ servers, lifecycle, traffic inspector
├── □ Provider fallback: auto-switch on quota error
├── □ Provider dashboard: quota display, fallback editor
├── □ 15+ tests passing
└── □ Tag v0.9.1-beta

WEEK 4-5: Sessions
├── □ Session supervisor: spawn, restart, persist
├── □ Dashboard session cards: status, logs, terminal
├── □ Worktree isolation
├── □ Overview panel: health, sessions, quotas
├── □ 25+ tests passing
└── □ Tag v0.9.2-beta

WEEK 6: Ship
├── □ QUICKSTART.md (tested by someone who hasn't seen the project)
├── □ README.md rewrite (what, who, screenshot, Docker command)
├── □ 2-minute demo video
├── □ 30+ tests, 5+ integration tests
├── □ Tag v1.0.0 — first official GitHub Release
├── □ Post: r/LocalLLaMA, r/MachineLearning, Hacker News, Discord
└── □ Collect feedback from 10 external users

WEEK 7+: Reality
├── □ Fix bugs from real users
├── □ v1.5 priorities based on what users actually request
├── □ Do NOT add features nobody asked for
└── □ New minor release every 2 weeks
```

---

## SECTION 11: LONG-TERM PHILOSOPHY

### The Three Principles

> **Index broadly.** Track every relevant tool in `ecosystem/index.jsonc`. Know the landscape.
>
> **Integrate selectively.** Only bring code into the repo when it serves a tested capability. Use contracts, not cloning.
>
> **Ship relentlessly.** Monthly releases. Each has a changelog and a user-facing improvement.

### Architecture Evolution

```
v1.0 (Now):      All TypeScript. child_process isolation. Redis events.
v1.5 (3 months):  Memory plugins. Browser extension. 5 CLI adapters.
v2.0 (6 months):  Go supervisor (if needed). Swarm (Phases 81-95). Debate.
v3.0 (1 year):    P2P mesh (if demand). Mobile. Cloud dev. Enterprise.
```

### Feature Scoring (For Any Future Proposal)

| Question | Points |
|---|---|
| Real user (not Robert, not AI) requested this? | +3 |
| Makes existing feature more reliable? | +2 |
| Shippable in < 1 week? | +2 |
| Has clear acceptance test? | +1 |
| Requires new submodule? | -2 |
| Requires new language/framework? | -2 |
| Is analytics-about-analytics? | -5 |

**Only implement features scoring ≥ 3.**

### What to Measure

| Stop Measuring | Start Measuring |
|---|---|
| Phase count | GitHub stars from strangers |
| Internal version bumps | Published GitHub Releases |
| Submodule count | Time from `git clone` to working dashboard |
| Lines of code | Test pass rate in CI |
| — | Issues filed by non-Robert users |

---

## SECTION 12: THE TEN COMMANDMENTS

| # | Commandment |
|---|---|
| **I** | **Ship to real users.** Everything else is noise. |
| **II** | **Four features define v1.0.** MCP Router. Fallback. Supervisor. Dashboard. |
| **III** | **Orchestrate, don't contain.** Borg controls Aider; Borg is not Aider. |
| **IV** | **No submodule without human approval.** The repo contains only code that runs. |
| **V** | **Tests or it didn't happen.** Every feature. Every merge. CI green. |
| **VI** | **Stop when done.** No infinite loops. No phase inflation. No drift. |
| **VII** | **Capability contracts, not feature parity.** Abstract what Borg does; adapters implement. |
| **VIII** | **One language for now.** TypeScript. Add Go later only if pain demands it. |
| **IX** | **Monthly releases.** Changelog. Version. User-facing improvement. |
| **X** | **Index broadly. Integrate selectively. Ship relentlessly.** |

---

## ═══════════════════════════════════════════════════════════

*This document represents the unanimous consensus of seven AI architectures — trained by different organizations, with different weights, different objectives, and different perspectives — who all arrived at the same conclusion:*

*The vision is real. The architecture is sound. The core features are genuinely valuable and differentiated. What's missing is not more features. It's focus, tests, and a release.*

*Robert: you built 121 phases of a cathedral. Now lay the floor so people can walk in.*

*Ship Borg 1.0. Get 100 stars. Then assimilate the universe.*

***Resistance is futile — but only when you're focused.*** 🖖

## ═══════════════════════════════════════════════════════════

*Signed in consensus:*
*GPT-5.3 ✅ | GPT-5.4 ✅ | Gemini 3.1 Pro ✅ | Claude Opus 4.6 ✅ | Kimi K2.5 ✅ | Qwen 3.5 ✅ | Grok 4.1 ✅*
