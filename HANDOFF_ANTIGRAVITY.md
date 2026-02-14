# BORG — Complete Handoff Document
> **Generated**: 2026-02-13T19:13 EST  
> **Analysts**: Gemini 3 (first pass) → Claude Opus 4.6 (second pass)  
> **Purpose**: Self-contained reference for any AI model to continue work on this codebase.

> **Update (Copilot / GPT-5.3-Codex): 2026-02-13 late-pass reconciliation**
> This handoff was revalidated against live code evidence to prepare implementation by multiple models.

### Delta Summary — What Was Revalidated

- Previously flagged P0 auth TODOs are now implemented:
    - `apps/web/src/components/auth/LoginForm.tsx`
    - `apps/web/src/app/signup/page.tsx`
    - `apps/web/src/app/forgot-password/page.tsx`
- Previously flagged critical placeholder/no-op UI surfaces are now wired:
    - `GlobalSearch.tsx`, `ConfigEditor.tsx`, `SystemStatus.tsx`, `TestStatusWidget.tsx`, `RemoteAccessCard.tsx`, `GraphWidget.tsx`
- Confirmed orchestration realism gap remains:
    - `packages/core/src/agents/SubAgents.ts` has now been rewired to execute via researcher/coder pathways instead of staged delay simulation.
- Knowledge dashboard integrity issue is resolved in current slice:
    - duplicate coder state/mutation block removed in `apps/web/src/app/dashboard/knowledge/page.tsx`
- `packages/core/src/lib/trpc-core.ts` now exposes a **transitional typed** `getMcpServer()` runtime shape (including `directorConfig.demo_mode`) instead of raw `any`.

### Documentation Sync Performed in This Pass

- Updated `ROADMAP.md` Phase 63 into explicit coverage-reconciliation mode (domain matrix + P0/P1/P2/P3 backlog).
- Updated `STATUS.md` to align canonical version with `VERSION.md` (2.6.1) and explicitly call out changelog drift.
- Updated `docs/DETAILED_BACKLOG.md` with domain coverage matrix and implementation tracks for extension/parity/service-exposure.

### Immediate Execution Order for Implementor Models

1. **P0 UX integrity**: auth flows + six disabled dashboard surfaces.
2. **P0/P1 knowledge/orchestration integrity**: remove duplicate knowledge block; replace simulated sub-agent execution.
3. **P2 maintainability**: type `getMcpServer()`, remove hardcoded TRPC endpoint, normalize expert router pattern.
4. **Service representation pass**: resolve router/UI strategy for Mesh/SkillAssimilation/Policy/Browser services.

### Critical Governance Note

- `VERSION.md` currently indicates `2.6.1`, while `CHANGELOG.md` includes `2.6.2` entries.
- Treat this as active documentation/version drift until resolved in a dedicated version-sync change.

### Implementation Slice Update (Auth + Search + Config + Widgets + Knowledge)

- ✅ Added auth API routes in `apps/web`:
    - `src/app/api/auth/signup/route.ts`
    - `src/app/api/auth/login/route.ts`
    - `src/app/api/auth/forgot-password/route.ts`
- ✅ Added local auth persistence helper: `src/lib/authStore.ts`.
- ✅ Wired submit handlers in:
    - `src/components/auth/LoginForm.tsx`
    - `src/app/signup/page.tsx`
    - `src/app/forgot-password/page.tsx`
- ✅ Restored real behavior in:
    - `src/components/GlobalSearch.tsx` (`trpc.lsp.searchSymbols` + `vscode://file/...` open flow)
    - `src/components/ConfigEditor.tsx` (`trpc.settings.get/update`)
- ✅ Restored live dashboard widgets:
    - `src/components/SystemStatus.tsx` (`trpc.metrics.systemSnapshot`)
    - `src/components/TestStatusWidget.tsx` (`trpc.tests.status/start/stop`)
    - `src/components/RemoteAccessCard.tsx` (remote access tool controls via `trpc.executeTool`)
    - `src/components/GraphWidget.tsx` (VS Code deep-link open flow)
- ✅ Cleaned `src/app/dashboard/knowledge/page.tsx`:
    - removed duplicate coder state/mutation block
    - removed `@ts-ignore` at `trpc.expert.*` call sites (typed shim)
- ✅ P1 realism update:
    - `packages/core/src/agents/SubAgents.ts` now dispatches through `researcherAgent`/`deepResearchService` and `coderAgent` with structured completion/failure handling.
- ✅ P2 hardening update:
    - `apps/web/src/utils/TRPCProvider.tsx` no longer hardcodes a single localhost endpoint; it now resolves from `NEXT_PUBLIC_TRPC_URL` with runtime-safe local/prod fallbacks.
    - `packages/core/src/trpc.ts` `executeTool.args` now uses `z.record(z.unknown())` (replacing permissive `z.any()`).
    - `packages/core/src/services/ProjectTracker.ts` now performs schema-based task-line extraction (Zod) and indentation-aware subtask traversal.
    - Oversized inline system procedures were extracted from `packages/core/src/trpc.ts` into `packages/core/src/routers/systemProcedures.ts` (API shape preserved).
    - `packages/core/src/services/MemoryManager.ts` now chunks large/code contexts with `CodeSplitter` and uses metadata-aware symbol filtering in `searchSymbols()`.
    - `packages/core/src/installer/AutoConfig.ts` is now production-safe: validated config schema, sanitized env allowlist, robust container detection, and explicit config write helpers.
    - Added placeholder/no-op regression guard (`scripts/check-placeholder-regressions.mjs`, exposed as `npm run check:placeholders`) and validated it against current sources.
    - Security dashboard now performs real autonomy lockdown mutation instead of disabled-router placeholder behavior.
    - Incremental MCP accessor typing hardening landed for workflow surfaces: shared typed workflow helpers in `trpc-core` now back `workflowRouter`, removing explicit `(getMcpServer() as any)` usage there and tightening workflow `initialState` input to `z.record(z.unknown())`.
    - Incremental MCP accessor typing hardening expanded to service routers: shared typed helpers now back `testsRouter`, `healerRouter`, and `gitRouter`, reducing direct untyped service property access from those endpoints.
    - Incremental MCP accessor typing hardening continued into governance/control paths: `autonomyRouter` and `auditRouter` now use shared typed helpers, and director daemon hooks no longer use raw `as any` casts.
    - Incremental MCP accessor typing hardening continued into orchestration paths: `directorRouter`, `directorConfigRouter`, and `darwinRouter` now use shared typed helpers from `trpc-core`, replacing direct untyped service property access.
    - Incremental MCP accessor typing hardening expanded into session/billing/knowledge paths: those routers now use shared typed helpers for session manager, LLM/quota, and knowledge/deep-research service access.
    - Incremental MCP accessor typing hardening expanded into MCP/submodule/skills paths: `mcpRouter`, `submoduleRouter`, and `skillsRouter` now use shared typed helpers for aggregator/submodule/skills services from `trpc-core`.
    - Incremental MCP accessor typing hardening expanded into symbols/settings/research/system paths: `symbolsRouter`, `settingsRouter`, `researchRouter`, and `systemProcedures` now use shared typed helpers for symbol pin/config/research/tracker/LSP access from `trpc-core`.
    - Incremental MCP accessor typing hardening expanded into collaboration paths: `squadRouter` and `suggestionsRouter` now use shared typed helpers for squad/suggestion services from `trpc-core`, removing direct `mcpHelper` imports.
    - Incremental MCP accessor typing hardening expanded into context/automation/runtime paths: `contextRouter`, `commandsRouter`, `autoDevRouter`, `shellRouter`, `pulseRouter`, `memoryRouter`, and `graphRouter` now use shared typed helpers from `trpc-core`, removing the remaining direct `mcpHelper` router imports.
    - Router-level cast cleanup completed: remaining broad router cast in `graphRouter` removed in favor of reflection-safe init-state access; router tree now has no `as any` matches.
    - Non-router cast hardening advanced across services: `HealerService` and `DarwinService` now use typed LLM-response extraction helpers, `KnowledgeService` now guards optional graph snapshot capability via runtime narrowing, `MemoryManager` now uses typed dynamic module/store contracts and provider capability guards, and `AgentMemoryService` now hydrates memory metadata through typed narrowing helpers.
    - Service hardening follow-through: `CodeModeService` permissive sandbox global exposure now uses a typed `globalThis` lookup helper, removing the last real service-level `(global as any)` pattern.
    - Additional non-service cast hardening: `SuggestionService` now normalizes LLM context-analysis responses through typed extraction, `SandboxService` now uses a typed promise-like runtime guard for async VM results, and `SpawnerService` now uses reflection-safe internal fail invocation with typed snapshot output.
    - Follow-up non-service cast hardening: `SystemCommands` now extracts director status via reflection-safe typed narrowing, and `TerminalSensor` now applies/removes a typed stderr hook without broad write casts.
    - Latest low-risk cast hardening: `MCPServer` native reader fallback now uses typed tool-handler guards and promise normalization, and `SquadService` now uses reflection-safe council/memory runtime access plus typed lazy `IndexerJob` contracts.
    - Milestone: `packages/core/src` now reports zero `as any` matches under ripgrep scan; next frontier is primarily `packages/core/test/**` mock-cast cleanup.
    - Test cleanup kickoff: `test/chaining.test.ts` now uses typed mock constructor narrowing and typed internal-tool callback args.
    - Test cleanup expansion: additional cast-hardening landed in `endpoint_namespace_propagation.test.ts`, `UnifiedMemorySystem.test.ts`, `agents/ArchitectMode.test.ts`, `hub_endpoint_auth.test.ts`, `services/RbacService.test.ts`, `tool_annotations_persistence.test.ts`, `Phase24_Browser.test.ts`, `Phase25_Aggregator.test.ts`, `McpmInstaller.test.ts`, and `routes/julesKeeperRoutesHono.test.ts`.
    - Test cleanup expansion (wave 2): additional cast-hardening landed in `endpoint_policy_context.test.ts`, `endpoint_policy_regression.test.ts`, `run_code.test.ts`, `session_tool_filtering.test.ts`, `proxy_middleware.test.ts`, `scripts.test.ts`, `nested_trace_logging.test.ts`, `proxy_logging_middleware.test.ts`, and `run_agent.test.ts`.
    - Validation pass: `npx tsc -p packages/core/tsconfig.json --noEmit` and root `pnpm run check:placeholders` both pass after the service-layer hardening batch.
    - Test cleanup completion: remaining cast-hardening landed in `namespace_tool_overrides.test.ts`, `namespace_override_annotations.test.ts`, `tool_annotations_call_block.test.ts`, `tool_annotations_middleware.test.ts`, `toolset_load_set.test.ts`, `toolset_meta_tools.test.ts`, `script_dynamic_tools.test.ts`, `run_task.test.ts`, `PolicyVerification.test.ts`, `services/RemoteSupervisor.test.ts`, and `services/AgentMemoryService.test.ts`.
    - Test cast-scan milestone: `packages/core/test` now reports zero `as any` matches, and compile + placeholder checks were re-verified after the final test cleanup pass.
    - Typing micro-batch: reduced broad `any` in `trpc-core` council runtime contracts (`Promise<unknown>` / `unknown` result placeholders) and added unknown-safe RBAC error message narrowing helper.
    - Typing micro-batch: tightened `SuggestionService` payload type from `any` to `unknown` while preserving behavior and storage shape.
    - Validation pass: `npx tsc -p packages/core/tsconfig.json --noEmit` and root `pnpm run check:placeholders` re-verified after latest typing hardening.
    - Typing micro-batch: tightened `AuditService` payload/result signatures from `any` to `unknown` in both service-layer and security-layer audit implementations.
    - Typing micro-batch: tightened `PolicyService` rule/check signatures and upgraded `PolicyEngine` filesystem argument handling with unknown-safe target-path extraction.
    - Validation pass: compile + placeholder checks re-verified after the service/security hardening batch.
    - Typing micro-batch: tightened `ContextPruner` token/prune method signatures to `Message[]` and added a defensive guard for shifted pruning entries.
    - Typing micro-batch: tightened `GitWorktreeManager` parser/output typing with explicit `WorktreeInfo` instead of broad `any[]` structures.
    - Validation pass: compile + placeholder checks re-verified after the utility hardening batch.
    - Typing micro-batch: tightened `ConfigManager` load/save contracts with structured config typing and runtime JSON object validation.
    - Typing micro-batch: aligned `MemoryManager.pruneContext/getContextSize` to `Message[]` signatures for consistent typed context-pruning boundaries.
    - Validation pass: compile + placeholder checks re-verified after config/memory hardening.
    - Typing micro-batch: removed router-local broad `any` patterns in `settingsRouter`, `graphRouter`, and `billingRouter` via unknown-safe error handling and typed projection interfaces.
    - Typing micro-batch: exported graph symbol projection interfaces to satisfy portable `appRouter` type emission (TS4023-safe).
    - Validation pass: compile + placeholder checks re-verified after router typing hardening.
    - Typing micro-batch: removed remaining router-local `any` in `supervisorRouter` task filtering using typed task-status narrowing.
    - Typing micro-batch: migrated `lspRouter` to typed `getLspService()` access with explicit service-availability guards and expanded LSP runtime contracts in `trpc-core`.
    - Validation pass: compile + placeholder checks re-verified after supervisor/LSP router hardening.
    - Typing micro-batch: tightened `PermissionManager` guardrail signatures by replacing broad tool-arg `any` usage with `unknown`.
    - Typing micro-batch: tightened `SystemWorkflows` tool-runner contract and catch handling to unknown-safe patterns.
    - Validation pass: compile + placeholder checks re-verified after security/orchestration hardening.

---

## 1. Project Architecture

**Borg** is a Neural Operating System built as a monorepo:

```
borg/
├── packages/
│   ├── core/       ← MCPServer (2809 lines), 34 TRPC routers, 30 services
│   ├── agents/     ← Director, Council, CoderAgent, ResearcherAgent, etc.
│   ├── tools/      ← MCP tool implementations (file, browser, search, etc.)
│   ├── ui/         ← Shared React components (@borg/ui)
│   └── cli/        ← Terminal interface
├── apps/
│   └── web/        ← Next.js dashboard (31 pages)
└── ROADMAP.md      ← Phase tracking (Phases 1-67+)
```

**Key Pattern**: `MCPServer.ts` (the god class, 2809 lines) instantiates ~40 services and exposes them as public properties. TRPC routers in `packages/core/src/routers/*.ts` access these via `getMcpServer()` (currently a transitional runtime-typed accessor, not yet fully strict interface-typed). The Next.js dashboard consumes these routers via `@trpc/react-query`.

---

## 2. MCPServer Public Interface (Complete)

Every public property on `MCPServer` (from `packages/core/src/MCPServer.ts`, lines 153-260):

| Property | Type | Phase | Has Router? | Has Dashboard? |
|----------|------|-------|-------------|----------------|
| `modelSelector` | `ModelSelector` | Core | `billingRouter` | `billing/` |
| `llmService` | `LLMService` | Core | `metricsRouter` | `metrics/` |
| `director` | `Director` | Core | `directorRouter` | `director/` |
| `permissionManager` | `PermissionManager` | Core | — | `security/` |
| `auditService` | `AuditService` | Core | `auditRouter` | — |
| `shellService` | `ShellService` | Core | `shellRouter` | — |
| `memoryManager` | `MemoryManager` | Core | `memoryRouter` | `memory/` |
| `suggestionService` | `SuggestionService` | Core | `suggestionsRouter` | — |
| `configManager` | `ConfigManager` | Core | `settingsRouter` | `settings/` |
| `autoTestService` | `AutoTestService` | Core | `testsRouter` | — |
| `sandboxService` | `SandboxService` | Core | — | — |
| `spawnerService` | `SpawnerService` | Core | — | — |
| `squadService` | `SquadService` | Core | `squadRouter` | `squads/` |
| `gitWorktreeManager` | `GitWorktreeManager` | Core | `gitRouter` | — |
| `commandRegistry` | `CommandRegistry` | Core | `commandsRouter` | `command/` |
| `contextManager` | `ContextManager` | Core | `contextRouter` | — |
| `symbolPinService` | `SymbolPinService` | Core | `symbolsRouter` | — |
| `autoDevService` | `AutoDevService` | Core | `autoDevRouter` | — |
| `researchService` | `ResearchService` | Core | `researchRouter` | `research/` |
| `gitService` | `GitService` | 30 | `gitRouter` | — |
| `metricsService` | `MetricsService` | 31 | `metricsRouter` | `metrics/` |
| `policyService` | `PolicyService` | 32 | — | `security/` |
| `knowledgeService` | `KnowledgeService` | — | `knowledgeRouter` | `knowledge/` |
| `healerService` | `HealerService` | — | `healerRouter` | `healer/` |
| `darwinService` | `DarwinService` | — | `darwinRouter` | `evolution/` |
| `promptRegistry` | `PromptRegistry` | — | — | — |
| `skillRegistry` | `SkillRegistry` | — | `skillsRouter` | `skills/` |
| `skillAssimilationService` | `SkillAssimilationService` | — | **NONE** | **NONE** |
| `mcpAggregator` | `MCPAggregator` | — | `mcpRouter` | `mcp/` |
| `eventBus` | `EventBus` | — | `pulseRouter` | `events/`, `pulse/` |
| `deepResearchService` | `DeepResearchService` | — | `knowledgeRouter` | `knowledge/` |
| `councilService` | `CouncilService` | — | `councilRouter` | `council/` |
| `browserService` | `BrowserService` | — | — | — |
| `sessionManager` | `SessionManager` | 57 | `sessionRouter` | — |
| `projectTracker` | `ProjectTracker` | 59 | inline in `trpc.ts` | `director/` |
| `meshService` | `MeshService \| undefined` | 60 | **NONE** | **NONE** |
| `lspService` | `LSPService` | 51 | `lspRouter` | — |
| `planService` | `PlanService` | 51 | `planRouter` | `plans/` |
| `codeModeService` | `CodeModeService` | 51 | — | `code/` |
| `workflowEngine` | `WorkflowEngine` | 51 | `workflowRouter` | `workflows/` |
| `agentMemoryService` | `AgentMemoryService` | 53 | `agentMemoryRouter` | `memory/` |
| `council` | `Council` | — | `councilRouter` | `council/` |
| `supervisor` | `Supervisor` | — | `supervisorRouter` | `supervisor/` |
| `coderAgent` | `CoderAgent` | — | `expertRouter` | `knowledge/` |
| `researcherAgent` | `ResearcherAgent` | — | `expertRouter` | `knowledge/` |

---

## 3. TRPC Router Mount Map

From `packages/core/src/trpc.ts` — every namespace and its mount key:

```typescript
export const appRouter = t.router({
    graph:          graphRouter,         // RepoGraphService / SymbolGraph
    workflow:       workflowRouter,      // WorkflowEngine (start/pause/resume/approve/reject)
    tests:          testsRouter,         // AutoTestService
    borgContext:    contextRouter,       // ContextManager
    commands:       commandsRouter,      // CommandRegistry
    symbols:        symbolsRouter,       // SymbolPinService
    autoDev:        autoDevRouter,       // AutoDevService
    shell:          shellRouter,         // ShellService
    memory:         memoryRouter,        // MemoryManager + AgentMemoryService
    knowledge:      knowledgeRouter,     // KnowledgeService + DeepResearchService
    research:       researchRouter,      // ResearchService
    pulse:          pulseRouter,         // EventBus
    skills:         skillsRouter,        // SkillRegistry
    squad:          squadRouter,         // SquadService
    suggestions:    suggestionsRouter,   // SuggestionService
    council:        councilRouter,       // Council + CouncilService ⚠️
    supervisor:     supervisorRouter,    // Supervisor
    metrics:        metricsServiceRouter,// MetricsService
    lsp:            lspRouter,           // LSPService
    agentMemory:    agentMemoryRouter,   // AgentMemoryService
    planService:    planServiceRouter,   // PlanService
    settings:       settingsRouter,      // ConfigManager
    session:        sessionRouter,       // SessionManager
    billing:        billingRouter,       // LLMService.getCostStats()
    mcp:            mcpRouter,           // MCPAggregator
    healer:         healerRouter,        // HealerService
    darwin:         darwinRouter,        // DarwinService
    autonomy:       autonomyRouter,      // Director.autonomyLevel
    director:       directorRouter,      // Director
    directorConfig: directorConfigRouter,// directorConfig object
    git:            gitRouter,           // GitService
    audit:          auditRouter,         // AuditService
    submodule:      submoduleRouter,     // SubmoduleService
    expert:         expertRouter,        // CoderAgent + ResearcherAgent ⚠️

    // Inline procedures (not in separate files):
    health:         /* inline */ -> { status: 'running', service: '@borg/core' }
    getTaskStatus:  /* inline */ -> ProjectTracker.getStatus()
    indexingStatus: /* inline */ -> LSPService.getStatus()
    executeTool:    /* inline */ -> MCPServer.executeTool(name, args)
});
```

---

## 4. Dashboard Pages vs Routers (Complete Cross-Reference)

### ✅ Confirmed Wired (11 pages)
| Dashboard Dir | Router Used | Status |
|---|---|---|
| `billing/` | `trpc.billing.getStatus` | ✅ Real data |
| `healer/` | `useHealerStream` (WebSocket) | ✅ Real data |
| `knowledge/` | `trpc.knowledge.*`, `trpc.expert.*`, `trpc.submodule.*` | ⚠️ Has bugs |
| `mcp/` | `trpc.mcp.listServers/listTools/addServer/removeServer` | ✅ Wired (no aggregator yet) |
| `memory/` | `trpc.memory.*`, `trpc.knowledge.getGraph` | ✅ Real data |
| `metrics/` | `trpc.metrics.getStats` | ✅ Real data |
| `skills/` | `trpc.skills.*` | ✅ Real data |
| `squads/` | `trpc.squad.*` | ✅ Real data |
| `submodules/` | `trpc.submodule.*` | ✅ Real data |
| `workflows/` | `trpc.workflow.*` | ✅ Wired (may return empty) |
| `research/` | `trpc.research.*` | ✅ Wired |

### ❓ Unclear / Needs Audit (10 pages)
| Dashboard Dir | Best Guess Router | Notes |
|---|---|---|
| `architecture/` | None | Likely static/hardcoded visualization |
| `brain/` | `trpc.memory`? | May be alternate memory view |
| `chronicle/` | `trpc.pulse` or `trpc.audit`? | Event history? |
| `code/` | `trpc.autoDev`? | Code mode interface |
| `command/` | `trpc.commands`? | Command registry UI |
| `council/` | `trpc.council.*` | ⚠️ Has 3 `@ts-ignore` |
| `director/` | `trpc.director.*`, `trpc.directorConfig.*` | ✅ Likely wired |
| `inspector/` | WebSocket / pulse? | Traffic inspector |
| `library/` | `trpc.skills`? | Skill library browser |
| `reader/` | File system? | Code reader |

### 🔴 No Router Exists (4 pages)
| Dashboard Dir | Missing Router | What's Needed |
|---|---|---|
| `config/` | `configRouter`? | Settings already has `settingsRouter` — may be duplicate |
| `manual/` | None | Static docs page? |
| `security/` | `securityRouter` | `PolicyService` exists but no router exposes it |
| `workshop/` | None | Unknown purpose |

---

## 5. Bugs & Issues (Prioritized)

### 🔴 P0 — Will Crash at Runtime

| # | File | Line(s) | Issue | Fix |
|---|---|---|---|---|
| 1 | `knowledge/page.tsx` | 33-51 | **Duplicate code block**: `coderTask`, `coderMutation`, `handleCode` declared twice. React Hooks violation hidden by `@ts-ignore`. | Delete lines 43-51. |
| 2 | `expertRouter.ts` | 1, 13, 42, 58-59 | Uses raw `global.mcpServerInstance` (bypasses `getMcpServer()` null check). Wrong import path `../trpc.js`. | Replace with `getMcpServer()`, fix import to `../lib/trpc-core.js`. |
| 3 | `mcpRouter.ts` | 71 | `addServer` throws if `mcpAggregator` is null, but `mcpAggregator` may not be initialized. | Add null check or initialize aggregator eagerly. |

### 🟡 P1 — Incorrect Behavior

| # | File | Line(s) | Issue | Fix |
|---|---|---|---|---|
| 4 | `councilRouter.ts` | 8 vs 15 | Accesses `.council.runConsensusSession()` AND `.councilService.listSessions()` — two different objects. | Standardize. Both exist on MCPServer (line 201 + 232). Verify which methods live where. |
| 5 | `graphRouter.ts` | 12-19 | Fallback creates standalone `RepoGraphService(process.cwd())` if `autoTestService.repoGraph` missing. Two divergent instances possible. | Always use `getMcpServer().autoTestService.repoGraph` or register graph service directly. |
| 6 | `trpc-core.ts` | 10 | Duplicate comment: `// For local desktop app, we default to ADMIN unless config says 'demo_mode'` appears twice. | Delete duplicate line. |
| 7 | `trpc.ts` | 100 | `executeTool` input uses `z.any()` for args — no validation. | Define proper schema or use `z.record(z.unknown())`. |
| 8 | `knowledgeService` | MCPServer L191 | Declared `private` but accessed directly by `knowledgeRouter`. TypeScript allows this via `any` cast. | Make `public` or add accessor. |

### 🟠 P2 — Technical Debt

| # | Scope | Count | Issue |
|---|---|---|---|
| 9 | `apps/web/src/**` | **28** | `@ts-ignore` directives (13 components, 3 dashboard pages). Root cause: stale TRPC types. |
| 10 | `packages/core/src/routers/**` | **80+** | `(mcp as any).someService` casts. Root cause: `getMcpServer()` returns `any`. |
| 11 | `expertRouter.ts` | 4 | Raw `global.mcpServerInstance` instead of `getMcpServer()`. |
| 12 | Services index | — | `services/index.ts` exports 17 of 30 services. 13 services not exported. |
| 13 | `MeshService` | 171 lines | Fully implemented P2P swarm (Hyperswarm). Zero integration — no router, no dashboard. |
| 14 | `SkillAssimilationService` | 4KB | Exists on MCPServer. No router, no dashboard. |

---

## 6. The Root Cause: `getMcpServer()` Returns `any`

```typescript
// packages/core/src/lib/trpc-core.ts, line 32
export function getMcpServer(): any {    // ← THIS IS THE PROBLEM
    const server = global.mcpServerInstance;
    if (!server) throw new Error("MCPServer instance not found");
    return server;
}
```

**Impact**: Because this returns `any`, every router must cast: `(mcp as any).someService`. This means:
- No autocomplete
- No compile-time verification that services exist
- No type checking on method calls
- Refactoring `MCPServer` properties will silently break routers

**Fix**: Create and export an `MCPServerInterface`:

```typescript
// packages/core/src/lib/MCPServerInterface.ts
export interface MCPServerInterface {
    modelSelector: ModelSelector;
    llmService: LLMService;
    director: Director;
    // ... all 40+ public properties from MCPServer class
    executeTool(name: string, args: Record<string, unknown>): Promise<any>;
}

// trpc-core.ts
export function getMcpServer(): MCPServerInterface {
    const server = global.mcpServerInstance;
    if (!server) throw new Error("MCPServer instance not found");
    return server as MCPServerInterface;
}
```

This single change eliminates **all 80+ `as any` casts** and **all 4 `global.mcpServerInstance`** references.

---

## 7. `@ts-ignore` Locations (All 28 in Frontend)

| File | Count | Lines | What It Suppresses |
|------|-------|-------|-------------------|
| `SuggestionsPanel.tsx` | 5 | 9, 42, 45, 50, 53 | TRPC type mismatches |
| `SquadWidget.tsx` | 4 | 7, 12, 15, 23 | TRPC type mismatches |
| `knowledge/page.tsx` | 3 | 25, 35, 45 | `trpc.expert` not in generated types |
| `council/page.tsx` | 3 | 54, 58, 63 | TRPC type mismatches |
| `DirectorChat.tsx` | 2 | 41, 51 | TRPC type mismatches |
| `CouncilWidget.tsx` | 2 | 10, 16 | TRPC type mismatches |
| `TrafficInspector.tsx` | 1 | 57 | Unknown |
| `SkillsViewer.tsx` | 1 | 58 | TRPC type mismatches |
| `ShellHistoryWidget.tsx` | 1 | 7 | TRPC type mismatches |
| `SandboxWidget.tsx` | 1 | 12 | TRPC type mismatches |
| `Mermaid.tsx` | 1 | 17 | Library type issue |
| `HealerWidget.tsx` | 1 | 13 | TRPC type mismatches |
| `AutonomyControl.tsx` | 1 | 27 | TRPC type mismatches |
| `AuditLogViewer.tsx` | 1 | 7 | TRPC type mismatches |

**Resolution**: Rebuild `@borg/core` → regenerate TRPC types → most of these go away.

---

## 8. MetaMCP Assimilation Targets

The MetaMCP repo (`c:\Users\hyper\workspace\metamcp`) contains the aggregator logic that should be ported to fill Borg's `mcpAggregator` property.

### Source Files (in `metamcp/apps/backend/src/lib/metamcp/`)

| File | Size | Maps To (Borg) | Purpose |
|------|------|-----------------|---------|
| `metamcp-proxy.ts` | 43KB | `services/MetaMCPService.ts` | Core aggregating proxy — routes tool calls to downstream MCPs |
| `metamcp-server-pool.ts` | 19KB | `services/ConnectionPool.ts` | Manages lifecycle of downstream MCP connections |
| `mcp-server-pool.ts` | 20KB | Merge with above | Server process spawning and monitoring |
| `metamcp-middleware/` | Dir | `middleware/MiddlewareManager.ts` | Request/response interception pipeline |
| `auto-discovery.service.ts` | 11KB | `services/DiscoveryService.ts` | Auto-detect MCP servers on system |
| `server-health.service.ts` | 11KB | Enhance `mcpRouter` | Health check and heartbeat |
| `auto-reconnect.service.ts` | 13KB | Built into ConnectionPool | Resilient reconnection logic |
| `client.ts` | 8KB | Utility | MCP client wrapper |
| `tool-name-parser.ts` | 2KB | Utility | Namespaced tool name parsing |
| `log-store.ts` | 2KB | Utility | Request/response logging |
| `sessions.ts` | 2KB | Utility | Session management |

### What Borg Already Has
- `mcpRouter.ts` (4KB) — TRPC router with `listServers`, `listTools`, `addServer`, `removeServer`, `getStatus`
- `mcp/page.tsx` (12KB) — Dashboard UI with server cards, tool list, add-server form
- `MCPAggregator` import (MCPServer.ts line 197) — type exists but implementation may be incomplete

### What Borg Needs
1. The actual `MCPAggregator` class populated with MetaMCP's proxy + pool logic
2. Middleware pipeline for auth, logging, filtering
3. Auto-discovery for locally installed MCP servers
4. Health monitoring with auto-reconnect

---

## 9. Services NOT Exported from `index.ts`

These 13 services exist in `packages/core/src/services/` but are NOT re-exported from `services/index.ts`:

| Service | Size | Has Router? | Status |
|---------|------|-------------|--------|
| `AuditService.ts` | 2KB | `auditRouter` | Missing from index |
| `BrowserService.ts` | 3KB | None | Missing from index |
| `ContextPruner.ts` | 3KB | None | Internal utility |
| `CouncilService.ts` | 2KB | `councilRouter` | Missing from index |
| `DarwinService.ts` | 5KB | `darwinRouter` | Missing from index |
| `DeepResearchService.ts` | 9KB | `knowledgeRouter` | Missing from index |
| `EventBus.ts` | 2KB | `pulseRouter` | Missing from index |
| `MeshService.ts` | 5KB | **None** | Orphaned |
| `PolicyService.ts` | 1KB | **None** | No router exists |
| `ProjectTracker.ts` | 6KB | Inline in `trpc.ts` | Missing from index |
| `SessionManager.ts` | 2KB | `sessionRouter` | Missing from index |
| `SkillAssimilationService.ts` | 4KB | **None** | Orphaned |
| `SubmoduleService.ts` | 8KB | `submoduleRouter` | Missing from index |

---

## 10. Recommended Fix Order

### Tier 1: Quick Wins (< 30 min total)
1. **Delete duplicate code** in `knowledge/page.tsx` (lines 43-51)
2. **Fix `expertRouter.ts`**: change import to `../lib/trpc-core.js`, replace 4× `global.mcpServerInstance` with `getMcpServer()`
3. **Fix `trpc-core.ts`**: delete duplicate comment on line 10
4. **Fix `trpc.ts` line 100**: replace `z.any()` with `z.record(z.unknown())`

### Tier 2: Architectural Foundation (1-2 hours)
5. **Create `MCPServerInterface`** in `packages/core/src/lib/MCPServerInterface.ts` — type every public property. Change `getMcpServer()` return type. This eliminates 80+ `as any` casts.
6. **Rebuild TRPC types** — run full pipeline (`npm run build` in `packages/core`, then `apps/web`) to regenerate frontend types. This eliminates most of the 28 `@ts-ignore`.
7. **Standardize `councilRouter`** — both `.council` and `.councilService` exist; verify which methods belong where and use consistent access.

### Tier 3: Feature Completion (2+ hours each)
8. **MetaMCP Assimilation** — Port `metamcp-proxy.ts` → `MCPAggregator`. Wire to existing `mcpRouter` and `mcp/page.tsx`.
9. **Wire `MeshService`** — Create `meshRouter.ts`, create dashboard page.
10. **Wire `PolicyService`** — Create `securityRouter.ts`, connect to `security/page.tsx`.
11. **Audit orphan dashboard pages** — Check `architecture/`, `brain/`, `chronicle/`, `inspector/`, `library/`, `manual/`, `reader/`, `workshop/` for functionality.

### Tier 4: Polish
12. Export remaining 13 services from `services/index.ts`
13. Make `knowledgeService` public on MCPServer (currently `private`)
14. Add streaming support for Healer page active infections
15. Expose `WorkflowEngine` registered workflows properly (not via `(engine as any).workflows`)

---

## 11. Build & Dev Commands

```bash
# Install
npm install          # from monorepo root

# Build core (must be first — generates TRPC types)
cd packages/core && npm run build

# Build dashboard
cd apps/web && npm run build

# Dev mode
cd apps/web && npm run dev     # Dashboard at localhost:3000
cd packages/core && npm start  # MCP Server

# Type check
npx tsc --noEmit               # from monorepo root
```

---

## 12. Key Files Quick Reference

| Purpose | Path |
|---------|------|
| God class | `packages/core/src/MCPServer.ts` (2809 lines) |
| TRPC mount | `packages/core/src/trpc.ts` (113 lines) |
| TRPC init | `packages/core/src/lib/trpc-core.ts` (37 lines) |
| Server accessor | `packages/core/src/lib/mcpHelper.ts` |
| All routers | `packages/core/src/routers/*.ts` (34 files) |
| All services | `packages/core/src/services/*.ts` (30 files) |
| Dashboard pages | `apps/web/src/app/dashboard/*/page.tsx` (31 dirs) |
| Shared UI | `packages/ui/src/components/*.tsx` |
| Roadmap | `ROADMAP.md` (root) |
| MetaMCP source | `../metamcp/apps/backend/src/lib/metamcp/` |
