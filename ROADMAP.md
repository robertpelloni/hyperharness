# Borg Project Roadmap

> **Status**: Phase 63 (Codebase Hardening) - **IN PROGRESS**
> **Version**: 2.6.1 (canonical from `VERSION.md`)
> **Codename**: AIOS (AI Operating System)

---

## Phase 1-19: Foundations (Completed)

- [x] **Phase 1: Zero-Shot Architecting** - Initial concept and MVP
- [x] **Phase 10: The Local LLM Bridge** - Oufob/Ollama integration
- [x] **Phase 15: Deep Analysis** - Core system introspection
- [x] **Phase 19: The Enterprise Vision** - Defined commercial goals

## Phase 20-29: The Core Architecture (Completed)

- [x] **Phase 20: The Monorepo Transition** - Migrated to pnpm + Turborepo structure types/core/ui/apps
- [x] **Phase 21: The MCP Foundation** - Implemented Model Context Protocol SDK integration
- [x] **Phase 22: The tRPC Backbone** - Replaced REST with tRPC for type-safe frontend-backend comms
- [x] **Phase 23: The Director Agent** - Created the first autonomous loop (Director.ts)
- [x] **Phase 24: The Tool Registry** - Built `packages/tools` with standard file/terminal tools
- [x] **Phase 25: The Shell Service** - Robust PTY management for terminal persistence
- [x] **Phase 26: The Git Service** - Git integration for autonomous commits
- [x] **Phase 27: The System Dashboard** - Launched Next.js dashboard (v1)
- [x] **Phase 28: The Vector Memory** - Integrated LanceDB/ChromaDB for semantic search
- [x] **Phase 29: The Search Service** - Ripgrep + AST implementation

## Phase 30-39: The Intelligence Layer (Completed)

- [x] **Phase 30: Multi-Model Selector** - Support for Claude, Gemini, GPT via `LLMService`
- [x] **Phase 31: The Metrics Service** - Prometheus-style metrics for agent performance
- [x] **Phase 32: The Policy Engine** - Allowed/Blocked commands and file paths
- [x] **Phase 33: The Healer (Self-Correction)** - Automatic error analysis and fix generation
- [x] **Phase 34: The Darwin Protocol** - Prompt mutation and evolution experiments
- [x] **Phase 35: The Council** - Multi-persona voting mechanism (Architect/Critic/Builder)
- [x] **Phase 36: Deep Research V1** - Recursive web scraping and summarization
- [x] **Phase 37: The Knowledge Graph** - Cognee integration for semantic relationships
- [x] **Phase 38: The Skill Assimilator** - Reading docs and generating new tools
- [x] **Phase 39: The Visual Cortex** - Image analysis and screenshot capabilities

## Phase 40-49: The Expansion (Completed)

- [x] **Phase 40: The Browser Extension** - Chrome extension for browser-MCP bridge
- [x] **Phase 41: The IDE Bridge** - VS Code extension for internal observation
- [x] **Phase 42: The Submodule Surge** - Integrated 100+ reference submodules
- [x] **Phase 43: The Nervous System** - Event bus and global state management
- [x] **Phase 44: The CLI Harness** - Commander.js CLI with 11 functional groups
- [x] **Phase 45: Knowledge Assimilation** - Ingesting PDFs, docs, and codebase context
- [x] **Phase 46: Dashboard Rebuild (v2)** - Modern UI with Tailwind, glassmorphism, real-time stats
- [x] **Phase 47: The Supervisor** - High-level task decomposition agent
- [x] **Phase 48: The Worktree Manager** - Git worktrees for parallel agent squads
- [x] **Phase 49: The Sandbox** - Secure code execution environment (VM/Docker)

## Phase 50-59: The Singularity (Current)

- [x] **Phase 50: Grand Unification V1** - Connecting all disparate systems into one OS
- [x] **Phase 51: Infrastructure Hardening** - Robust error handling and process supervision
- [x] **Phase 52: The Universal Interface** - Standardized UI components across all 31 pages
- [x] **Phase 53: The Documentation Overhaul** - Rewrote all LLM instructions (v2.6.0)
- [x] **Phase 54: The Version Sync** - Unified versioning across all packages
- [x] **Phase 55: The Session Manager** - State persistence and crash recovery
- [x] **Phase 56: The Session Router** - Dashboard control of session state
- [x] **Phase 57: Resilient Intelligence** - Automatic fallback, Director hardening, heartbeats
- [x] **Phase 58: Grand Unification V2** (COMPLETED)
  - [x] `test_grand_unification_v2.ts` creation (PASSED)
  - [x] Submodule Dashboard
  - [x] Documentation Consistency
  - [x] Autonomous Loop Execution

## Phase 60+: Future Horizons

- [ ] **Phase 62: Ignition (Real Capabilities)**
    - [ ] `CoderAgent`: Real LLM coding integration
    - [ ] `ResearcherAgent`: DeepResearch integration
    - [ ] **Fix**: Skill Registry API Mismatch (Frontend vs Backend)
    - [ ] **Feature**: Real Memory Graph Visualization (remove mock)
    - [ ] **Harden**: Deep Research Tool Loading

- [ ] **Phase 63: Codebase Hardening + Feature Coverage Reconciliation (IN PROGRESS)**
  > Date: 2026-02-13
  > Scope: `apps/web`, `packages/core`, `packages/ui`, extension surfaces, and planning docs.
  > Goal: convert “build-stable but partially wired” state into fully end-to-end functionality.

  ### 63.A Coverage Matrix (Reality Check)
  - **MCP Router/Aggregator**: Partial
    - [x] `mcpRouter` + `/dashboard/mcp` UI present
    - [ ] MetaMCP-grade aggregation lifecycle (auto-discovery/reconnect/health/namespace) not fully wired
  - **Memory System**: Mostly functional, still partial
    - [x] Memory/agent-memory routers + dashboard pages
    - [ ] `MemoryManager` TODOs remain (`CodeSplitter` integration, metadata filtering, provider abstraction)
  - **Code Intelligence**: Partial
    - [x] Graph/LSP/symbol/search foundations and router coverage
    - [ ] Placeholder/simulated paths still present in orchestration + workflow examples
  - **Agent Orchestration**: Partial
    - [x] Director/Council/Supervisor routers and UI surfaces
    - [ ] `SubAgents.ts` still simulates execution (no real model-selector + tool chain execution)
  - **Provider & Billing**: Mostly complete
    - [x] `billing.getStatus` real data path
    - [ ] Full quota reset/plan transparency across all providers needs final UX pass
  - **Dashboard Enhancements (User Request)**:
    - [ ] **Submodule Dashboard V2**: Comprehensive list of all submodules/versions/dates/builds with directory structure explanation. (Priority: High)
    - [ ] **Jules Autopilot Integration**: Dashboard for Google Jules cloud dev env.
    - [ ] **AI Tools Dashboard**: Billing/Usage/OAuth status for all AI tools (OpenAI, Anthropic, Gemini, etc).
  - **Documentation Overhaul (User Request)**:
    - [ ] "Universal LLM Instructions" consolidation (CLAUDE.md, GEMINI.md, etc).
    - [ ] `VISION.md` update with "Ultimate Goal" text.
    - [ ] `DEPLOY.md` creation/update.
    - [ ] `MEMORY.md` creation/update.

  - **Browser Extension**: Early/stub state
    - [ ] `apps/extension` has minimal scaffold; full MCP bridge/memory sync not fully wired
  - **Session & Cloud Management**: Partial
    - [x] Session manager + session router
    - [ ] Cloud session parity / transfer / broadcast workflows incomplete
  - **Interfaces (WebUI/CLI)**: Partial
    - [x] Broad CLI + dashboard page presence
    - [ ] Several dashboard surfaces intentionally replaced by static placeholder/no-op behavior
  - **Integration Protocols**: Partial
    - [ ] Policy/Mesh/Skill Assimilation capabilities are not comprehensively exposed via router + UI
  - **Advanced Features**: Partial
    - [ ] NotebookLM/computer-use/web-search/spec tooling mostly referenced but not fully end-to-end represented

  ### 63.B P0 — Critical UX and Trust Gaps

  - [x] **Auth completion**: implemented submit flows + backend API routes for:
    - `apps/web/src/components/auth/LoginForm.tsx`
    - `apps/web/src/app/signup/page.tsx`
    - `apps/web/src/app/forgot-password/page.tsx`
  - [x] **Replace placeholder/no-op critical dashboard wiring**:
    - [x] `GlobalSearch.tsx` wired to `lsp.searchSymbols` + VS Code file-open flow
    - [x] `ConfigEditor.tsx` wired to `settings.get` / `settings.update`
    - [x] `SystemStatus.tsx` (live metrics)
    - [x] `TestStatusWidget.tsx` (real test router state)
    - [x] `RemoteAccessCard.tsx` (real tunnel controls)
    - [x] `GraphWidget.tsx` (re-enabled open-file action)
  - [x] **Fix knowledge dashboard integrity**:
    - [x] Removed duplicate coder state block in `apps/web/src/app/dashboard/knowledge/page.tsx`
    - [x] Removed `@ts-ignore` around `trpc.expert.*` (typed shim in place pending TRPC regen)

  ### 63.C P1 — Backend Realism Gaps

  - [x] `billing.getStatus` → real data via `QuotaService.getUsageByModel()`
  - [x] `getTaskStatus` → real progress from `ProjectTracker`
  - [x] `indexingStatus` → real state from `LSPService`
  - [x] Replace simulated sub-agent execution in `packages/core/src/agents/SubAgents.ts`
  - [ ] Replace placeholder assumptions in:
    - `packages/core/src/services/DeepResearchService.ts`
    - `packages/core/src/services/KnowledgeService.ts`
  - [ ] Persist kanban status changes through backend mutation:
    - `packages/ui/src/components/kanban-board.tsx`

  ### 63.D P2 — Technical Debt with Delivery Impact

  - [x] Remove hardcoded TRPC endpoint in `apps/web/src/utils/TRPCProvider.tsx`
  - [ ] Type `getMcpServer()` return type (eliminate broad `(mcp as any)` router pattern)
    - [x] Incremental hardening: added typed workflow-engine accessors in `trpc-core` and removed explicit `(getMcpServer() as any)` usage in `workflowRouter`
    - [x] Incremental hardening: added typed auto-test/healer/git accessors in `trpc-core` and migrated `testsRouter`, `healerRouter`, and `gitRouter` to shared helpers
    - [x] Incremental hardening: added typed permission/audit/director accessors in `trpc-core` and migrated `autonomyRouter` + `auditRouter`, replacing director `as any` daemon hook calls with optional typed invocations
    - [x] Incremental hardening: expanded typed helpers for director/darwin/memory/command/suggestion services and migrated `directorRouter`, `directorConfigRouter`, and `darwinRouter` to shared helper usage
    - [x] Incremental hardening: added typed session/llm/knowledge helpers and migrated `sessionRouter`, `billingRouter`, and `knowledgeRouter` to shared helper usage
    - [x] Incremental hardening: added typed MCP/submodule/skills helpers and migrated `mcpRouter`, `submoduleRouter`, and `skillsRouter` to shared helper usage
    - [x] Incremental hardening: added typed symbol/config/research/tracker/LSP helpers and migrated `symbolsRouter`, `settingsRouter`, `researchRouter`, and `systemProcedures` to shared helper usage
    - [x] Incremental hardening: added typed squad helpers and extended suggestion helper contracts, migrating `squadRouter` and `suggestionsRouter` away from `mcpHelper`
    - [x] Incremental hardening: added typed context/auto-dev/shell/event/memory/agent-memory helpers and migrated `contextRouter`, `commandsRouter`, `autoDevRouter`, `shellRouter`, `pulseRouter`, `memoryRouter`, and `graphRouter` off `mcpHelper`
    - [x] Incremental hardening: removed the remaining router-level broad `as any` cast in `graphRouter` (typed/reflection-safe init check), leaving no `as any` matches under `packages/core/src/routers`
    - [x] Incremental hardening: reduced non-router service `as any` debt in `HealerService`, `DarwinService`, `KnowledgeService`, `MemoryManager`, and `AgentMemoryService` using typed runtime helpers/guards (with explicit reason/what/why comments)
    - [x] Incremental hardening: removed remaining real service-level `as any` usage in `CodeModeService` by switching permissive global exposure to typed `globalThis` key lookup helper
    - [x] Incremental hardening: reduced non-service cast debt in `SuggestionService`, `SandboxService`, and `SpawnerService` via typed response/thenable/fail-invocation guards
    - [x] Incremental hardening: reduced non-service cast debt in `commands/lib/SystemCommands.ts` and `sensors/TerminalSensor.ts` via reflection-safe status extraction and typed stderr hook wrapping
    - [x] Incremental hardening: reduced remaining low-risk cast debt in `MCPServer.ts` native reader fallback paths and `orchestrator/SquadService.ts` council/memory/indexer runtime access
    - [x] Cast-scan milestone: `packages/core/src` now has zero `as any` matches (verified by ripgrep; previously remaining hits were comment text only)
    - [x] Follow-on hardening: started `packages/core/test` cleanup in `test/chaining.test.ts` with typed constructor narrowing and tool-arg typing
    - [x] Follow-on hardening: additional `packages/core/test` cast cleanup landed in `endpoint_namespace_propagation.test.ts`, `UnifiedMemorySystem.test.ts`, `agents/ArchitectMode.test.ts`, `hub_endpoint_auth.test.ts`, `services/RbacService.test.ts`, `tool_annotations_persistence.test.ts`, `Phase24_Browser.test.ts`, `Phase25_Aggregator.test.ts`, `McpmInstaller.test.ts`, and `routes/julesKeeperRoutesHono.test.ts`
    - [x] Follow-on hardening: expanded test cast cleanup in `endpoint_policy_context.test.ts`, `endpoint_policy_regression.test.ts`, `run_code.test.ts`, `session_tool_filtering.test.ts`, `proxy_middleware.test.ts`, `scripts.test.ts`, `nested_trace_logging.test.ts`, `proxy_logging_middleware.test.ts`, and `run_agent.test.ts`
    - [x] Follow-on hardening: completed remaining `packages/core/test` cast cleanup in `namespace_tool_overrides.test.ts`, `namespace_override_annotations.test.ts`, `tool_annotations_call_block.test.ts`, `tool_annotations_middleware.test.ts`, `toolset_load_set.test.ts`, `toolset_meta_tools.test.ts`, `script_dynamic_tools.test.ts`, `run_task.test.ts`, `PolicyVerification.test.ts`, `services/RemoteSupervisor.test.ts`, and `services/AgentMemoryService.test.ts`
    - [x] Cast-scan milestone: `packages/core/test` now has zero `as any` matches (verified by ripgrep)
    - [x] Validation follow-through: `npx tsc -p packages/core/tsconfig.json --noEmit` and root `pnpm run check:placeholders` pass after service-layer hardening batch
    - [x] Validation follow-through: compile + placeholder checks re-verified after final test-layer cast cleanup batch
    - [x] Incremental typing hardening: reduced broad `any` boundaries in `trpc-core` council runtime contracts and added unknown-safe error message narrowing in RBAC middleware
    - [x] Incremental typing hardening: tightened `SuggestionService` metadata payload from `any` to `unknown` while preserving runtime behavior
    - [x] Validation follow-through: compile + placeholder checks re-verified after latest typing micro-batch
    - [x] Incremental typing hardening: tightened `AuditService` payload/result signatures in both `services` and `security` layers from `any` to `unknown`
    - [x] Incremental typing hardening: tightened `PolicyService` rule/check contracts and updated `PolicyEngine` to unknown-safe argument handling via typed target-path extraction helper
    - [x] Validation follow-through: compile + placeholder checks re-verified after service/security typing batch
    - [x] Incremental typing hardening: tightened `ContextPruner` public token/prune method signatures from broad `any[]` to typed `Message[]` with safe shift-guard handling
    - [x] Incremental typing hardening: tightened `GitWorktreeManager.listWorktrees()` from `Promise<any[]>` to `Promise<WorktreeInfo[]>` with explicit parser structure
    - [x] Validation follow-through: compile + placeholder checks re-verified after utility typing batch
    - [x] Incremental typing hardening: tightened `ConfigManager` signatures (`loadConfig`, `saveConfig`) to structured object contracts with runtime JSON object guarding
    - [x] Incremental typing hardening: aligned `MemoryManager.pruneContext/getContextSize` signatures to `Message[]` (matching `ContextPruner` typed contracts)
    - [x] Validation follow-through: compile + placeholder checks re-verified after config/memory typing batch
    - [x] Incremental typing hardening: removed router-local broad `any` usage in `settingsRouter` catch handling, `graphRouter` symbol graph projections, and `billingRouter` model/fallback mapping
    - [x] Incremental typing hardening: exported `graphRouter` symbol graph interfaces to satisfy app router portable type emission (TS4023-safe)
    - [x] Validation follow-through: compile + placeholder checks re-verified after router typing batch
    - [x] Incremental typing hardening: removed remaining router-local `any` task-filter usage in `supervisorRouter` with typed task runtime narrowing
    - [x] Incremental typing hardening: migrated `lspRouter` to typed `getLspService()` access with explicit service-availability guard and expanded LSP runtime contract coverage in `trpc-core`
    - [x] Validation follow-through: compile + placeholder checks re-verified after supervisor/LSP router typing batch
    - [x] Incremental typing hardening: tightened `PermissionManager` guardrail signatures by replacing broad `any` tool-arg inputs with `unknown`
    - [x] Incremental typing hardening: tightened `SystemWorkflows` tool-runner contract to `Record<string, unknown> -> Promise<unknown>` and replaced broad catch typing with unknown-safe error extraction
    - [x] Validation follow-through: compile + placeholder checks re-verified after security/orchestration typing batch
  - [x] Normalize `expertRouter.ts` singleton/import style with current router conventions
  - [x] Replace permissive `executeTool` input schema (`z.any`) with `z.record(z.unknown())`
  - [x] Harden `ProjectTracker` parsing with schema-based task extraction
  - [x] Implement remaining `MemoryManager` TODO paths
  - [x] Replace `AutoConfig.ts` stubbed/ported behavior with production-safe implementation
  - [x] Split oversized inline portions of `packages/core/src/trpc.ts`

  ### 63.E P3 — Representation, Robustness, and Governance

  - [ ] Expose backend-only services through intentional contracts (or explicitly archive):
    - `MeshService`, `SkillAssimilationService`, `PolicyService`, `BrowserService`
  - [x] Add CI check for placeholder/no-op regression patterns before release
  - [ ] Keep canonical doc synchronization strict (`ROADMAP.md` + `STATUS.md` + `docs/DETAILED_BACKLOG.md` + `HANDOFF_ANTIGRAVITY.md`)

### Phase 63 Exit Criteria

- [ ] All auth and critical dashboard actions execute real mutations/queries (no TODO submit handlers).
- [ ] No production page depends on disabled router endpoints without explicit capability fallback.
- [ ] No `@ts-ignore` remains on `trpc.expert` call sites in dashboard knowledge flows.
- [ ] Simulated orchestration paths are replaced with real model/tool execution.
- [ ] Handoff and status docs reflect the same version and priority queue.

- [ ] **Phase 64: The Mesh** - P2P agent swarm coordination
- [ ] **Phase 65: The Marketplace** - Decentralized tool/agent marketplace
- [ ] **Phase 66: The Neural Link** - Direct brain-computer interface (BCI) integration patterns
- [ ] **Phase 67: The Hive Mind** - Shared learning across all Borg instances
