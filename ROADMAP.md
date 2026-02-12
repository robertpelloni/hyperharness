# Borg Project Roadmap

> **Status**: Phase 58 (Grand Unification) - **IN PROGRESS**
> **Version**: 2.6.0
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

- [x] **Phase 62: Ignition (Real Capabilities)**
    - [ ] `CoderAgent`: Real LLM coding integration
    - [ ] `ResearcherAgent`: DeepResearch integration
    - [ ] **Fix**: Skill Registry API Mismatch (Frontend vs Backend)
    - [ ] **Feature**: Real Memory Graph Visualization (remove mock)
    - [ ] **Harden**: Deep Research Tool Loading

- [ ] **Phase 63: Codebase Hardening (Audit Findings — 2026-02-09)**
  > Exhaustive audit: 31+ router namespaces, 32 dashboard pages, 30 services
  - **P0 — Critical** (Breaks Core UX)
    - [x] `skillsRouter.list` → Query `SkillRegistry` (returns real data now)
    - [x] `researchRouter.conduct` → Fixed context access via `getMcpServer()`
    - [x] `pulseRouter.getLatestEvents` → Implemented `EventBus` history buffer, refactored to `getMcpServer()`
  - **P1 — High** (Fake Data Masquerading as Real)
    - [x] `billing.getStatus` → Real cost data via `QuotaService.getUsageByModel()`
    - [x] `getTaskStatus` → Returns real progress from `ProjectTracker` via `getMcpServer()`
    - [x] `indexingStatus` → Returns real state from `LSPService` via `getMcpServer()`
    - [x] Director page → Wired to real plan data via `trpc.directorConfig.get`, `trpc.getTaskStatus`, `trpc.autonomy.getLevel`
    - [x] Research page → Wired `handleResearch()` to `trpc.research.conduct.useMutation()`
  - **P2 — Medium** (Technical Debt)
    - [ ] Replace remaining ~50 `@ts-ignore global.mcpServerInstance` with `getMcpServer()` (8 routers: workflow, symbols, suggestions, squad, skills, shell, tests, graph)
    - [x] Fix `squadRouter.ts` import path (`../trpc.js` → `../lib/trpc-core.js`)
    - [ ] Fix `councilRouter` naming inconsistency (`council` vs `councilService`)
    - [x] Remove `repoGraph` duplicate mounting (alias of `graph`)
    - [ ] Cache tool→client mapping in `Router.callTool()` (currently O(N²))
    - [ ] Extract 15+ inline routers from `trpc.ts` into separate files
  - **P3 — Low** (Polish)
    - [x] Council page: Implemented "Members" tab with 4 role-based cards + Consensus Modes
    - [ ] Healer page: Add streaming for active infections
    - [ ] `workflowRouter.list`: Expose `WorkflowEngine` registered workflows

- [ ] **Phase 64: The Mesh** - P2P agent swarm coordination
- [ ] **Phase 65: The Marketplace** - Decentralized tool/agent marketplace
- [ ] **Phase 66: The Neural Link** - Direct brain-computer interface (BCI) integration patterns
- [ ] **Phase 67: The Hive Mind** - Shared learning across all Borg instances
