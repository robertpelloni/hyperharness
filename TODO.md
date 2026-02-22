# Borg — Master TODO List (Implementor Edition)

> **Updated**: 2026-02-22 | **Version**: 2.7.3 | **Phase**: 64 (Release Readiness)
> **Execution Order**: P0 (must close before release) → P1 (critical feature parity) → P2 (quality/polish) → P3 (future)

---

## P0 — Backend Reality Closure (Blockers)

### 1) Remove stub/simulated execution paths
- [x] **Replace simulated Saved Script execution with real runner path**
  - File: `packages/core/src/routers/savedScriptsRouter.ts`
  - Status: now executes via `CodeExecutorService` and returns structured execution metadata.

- [x] **Implement baseline OAuth token exchange flow (provider-specific)**
  - File: `packages/core/src/routers/oauthRouter.ts`
  - Status: now validates session/client context, performs provider token POST, parses/validates token payload, and persists tokens.
  - Follow-up hardening: encrypted token-at-rest + stricter state nonce contract.

- [x] **Finish agent chat orchestration wiring**
  - File: `packages/core/src/routers/agentRouter.ts`
  - Status: route now uses live `llmService` generation with graceful degraded fallback.

- [x] **Replace Researcher stub output with model-backed synthesis**
  - File: `packages/core/src/agents/Researcher.ts`
  - Status: now uses model-selected `generateText`, JSON extraction, and fallback tool-shape normalization.

- [x] **Remove remaining MetaMCP proxy stub branches/imports**
  - File: `packages/core/src/services/metamcp-proxy.service.ts`
  - Status: `run_code`, `run_python`, `save_script`, saved-script execution, tool search/sync, and `run_agent` now use live services/repositories or LLM-backed orchestration.
  - Residual note: `toon.serializer.stub` naming remains, but is now a utility concern rather than a runtime execution stub blocker.

### 2) Hybrid Storage Migration (Postgres + JSON)
- [x] **Migrate MCP server configurations to `mcp.json`**
  - Scope: `packages/core/src/db/repositories/mcp-servers.repo.ts`, `McpConfigService.ts`
  - Status: Dual-write implemented, startup sync active. `mcp.json` is now the source of truth for MCP config, with Postgres for auth/memory.
  - Verification: `scripts/verify-sync.ts` created and logic unit-tested.

### 2) Type-safety debt in persistence layer
- [x] **Eliminate broad `@ts-ignore` usage in repository layer**
  - Scope: `packages/core/src/db/repositories/*.ts`
  - Status: Completed in Phase 63 type-hardening batches I-XVI.
  - Validation: `rg -n "@ts-(ignore|nocheck)" packages/core/src` returns no matches.

---

## P1 — Frontend Representation Parity (Critical UX)

### 1) Dashboard route parity with backend capability map
- [x] **Create baseline AI Tools dashboard route with live core data**
  - Implemented: `apps/web/src/app/dashboard/mcp/ai-tools/page.tsx`
  - Live sources: `tools.list`, `mcpServers.list`, `apiKeys.list` + namespace coverage widgets (`expert.getStatus`, `session.getState`, `agentMemory.stats`, `serverHealth.check`, `shell.getSystemHistory`).
  - Remaining follow-up: provider-level billing/usage/OAuth matrix.
  - Target route: `/dashboard/mcp/ai-tools`
  - Remaining gap: nearest broader tools page (`packages/ui/src/app/dashboard/tools/page.tsx`) still enriches with mock status data.
  - Acceptance (advanced): provider auth/install/usage/billing are driven by live API responses.

- [x] **Add Jules dashboard route and wire baseline backend access**
  - Target route: `/dashboard/jules`
  - Implemented in `apps/web`:
    - `apps/web/src/app/dashboard/jules/page.tsx` (embedded Jules autopilot launch surface)
    - `apps/web/src/app/api/jules/route.ts` (Jules API proxy)
    - dashboard home card entry for discoverability
  - Follow-up completed: `packages/ui/src/lib/jules/client.ts` now has baseline `updateSession` (`PATCH`) support with graceful fallback for unsupported API versions.
  - Follow-up completed: `packages/ui/src/components/kanban-board.tsx` now attempts cloud sync on drag/drop status changes while preserving local persistence.

### 2) Router namespace coverage gap closure
- [x] **Add explicit UI usage surfaces for baseline coverage of**:
  - `agentMemory`, `expert`, `serverHealth`, `session`, `shell`
  - Source baseline: router usage audit from `packages/core/src/trpc.ts` vs `trpc.*` UI usage.
  - Implemented in: `apps/web/src/app/dashboard/mcp/ai-tools/page.tsx`
  - Follow-up: deepen controls/actions per namespace beyond baseline observability.

### 3) UI TODO/hardcoded fallback cleanup
- [ ] **Resolve Director endpoint test TODO**
  - File: `apps/web/src/components/DirectorConfig.tsx`
- [ ] **Resolve MCP namespace placeholder behavior**
  - File: `apps/web/src/components/mcp/EditEndpoint.tsx`
- [ ] **Complete prompt variable extraction implementation**
  - File: `apps/web/src/app/api/prompts/route.ts`

---

## P2 — Robustness, Security, and Extension Completion

### 0) First-party marker sweep follow-through (2026-02-16 audit)
- [ ] **Reduce remaining critical-path “for now/mock/placeholder” assumptions in first-party source**
  - Scope anchors from latest scan:
    - `packages/core/src/services/*` (`DeepResearchService`, `KnowledgeService`, `HealerService`, `MetaMCPController`, selected adapters)
    - `apps/web/src/app/dashboard/*` (monitoring/workflows/research/mcp pages with mock/fallback comments)
    - `packages/ui/src/components/*` (dashboard/tools/security/visualization placeholders)
  - Acceptance: production-path assumptions are replaced with explicit capability checks, telemetry, or completed implementations.

### 1) Browser extension completion
- [ ] **Finish fuzzy text matching click action**
  - File: `packages/browser-extension/background.js`
- [ ] **End-to-end MCP bridge validation**
  - Scope: extension ↔ Borg websocket bridge, memory sync hooks, reconnect behavior.

### 2) Shared UI `@ts-ignore` reduction pass
- [ ] **Remove remaining `@ts-ignore` in user-facing components**
  - Scope: `apps/web/src/components/*`, `packages/ui/src/components/*`
  - Acceptance: annotated temporary waivers only where upstream typing cannot yet be corrected.

### 3) Canonical docs consistency hardening
- [ ] **Keep `ROADMAP.md`, `TODO.md`, `STATUS.md`, `HANDOFF.md`, `CHANGELOG.md` synchronized after each merged fix batch**
  - Acceptance: no stale contradiction between code reality and canonical docs.

---

## P3 — Future Phases (Post-63)

- [x] **Phase 64: Release Readiness (v2.7.0)** — Documentation Freeze + E2E builds verified
- [x] **Phase 65: Marketplace & Ecosystem** — MCP registry, plugin engine, skills repos
- [x] **Phase 66: AI Command Center & Dashboards** — Jules, OpenCode, Billing, Tool Detector
- [x] **Phase 67: MetaMCP Submodule Assimilation** — Git submodule + MetaMCPBridgeService
- [ ] **Phase 68: Memory System Multi-Backend** — selectable vector stores, import/export
- [ ] **Phase 69: RAG Pipeline & Context Harvesting** — document intake, chunking, embedding

---

## Release Verification Checklist (Gate)

- [ ] `pnpm run check:dev-readiness` passes in strict mode with all critical services up
- [ ] `apps/web` typecheck + build pass
- [ ] `packages/core` typecheck pass
- [ ] `pnpm run check:placeholders` passes (or equivalent placeholder scan)
- [ ] No unresolved stub/simulated critical execution paths in P0 scope
- [ ] Canonical docs + version markers (`VERSION.md`, `package.json`, `CHANGELOG.md`) are synchronized

---

## Immediate Follow-Ups (Current Session)

- [ ] Add `--json` output mode to `scripts/verify_dev_readiness.mjs` for CI/dashboard ingestion.
- [ ] Add optional retries/backoff in readiness checks for startup race tolerance.
- [ ] Wire readiness command into release checklist automation script(s).

---

*This file is intentionally execution-ordered for implementor-model handoff and should be treated as the active backlog source of truth for Phase 63 closure.*
