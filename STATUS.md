# Borg Status Report

**Last Updated:** 2026-02-13  
**Canonical Version (from `VERSION.md`):** 2.6.1  
**Current Focus:** Phase 63 — Codebase Hardening + Feature Coverage Reconciliation

## Executive Snapshot

- **Backend baseline:** Strong service/router breadth in `packages/core`.
- **End-to-end reality:** Several user-facing flows are still placeholder, simulated, or TODO-backed.
- **Primary risk:** “Compiled/builds” has outpaced “fully wired + comprehensively represented in UI”.

## Latest Implementation Delta (2026-02-13)

- ✅ Added concrete auth API contract in `apps/web`:
	- `POST /api/auth/signup`
	- `POST /api/auth/login`
	- `POST /api/auth/forgot-password`
- ✅ Wired auth submit handlers in:
	- `apps/web/src/components/auth/LoginForm.tsx`
	- `apps/web/src/app/signup/page.tsx`
	- `apps/web/src/app/forgot-password/page.tsx`
- ✅ Replaced placeholder behavior in:
	- `apps/web/src/components/GlobalSearch.tsx` (real `trpc.lsp.searchSymbols` path)
	- `apps/web/src/components/ConfigEditor.tsx` (real `trpc.settings.get/update` path)
	- `apps/web/src/components/SystemStatus.tsx` (real `trpc.metrics.systemSnapshot` path)
	- `apps/web/src/components/TestStatusWidget.tsx` (real `trpc.tests.status/start/stop` path)
	- `apps/web/src/components/RemoteAccessCard.tsx` (real remote-access tool orchestration via `trpc.executeTool`)
	- `apps/web/src/components/GraphWidget.tsx` (restored VS Code deep-link open behavior)
- ✅ Fixed knowledge dashboard integrity in:
	- `apps/web/src/app/dashboard/knowledge/page.tsx` (removed duplicate coder state block + removed `@ts-ignore` usage)

## Coverage Reality (High Level)

- **MCP router/aggregator:** Partial (router/UI exists, full aggregation lifecycle still incomplete)
- **Memory:** Mostly working, with TODO debt in chunking/filtering/provider abstraction
- **Code intelligence:** Partial (good base + remaining simulated/placeholder behavior)
- **Orchestration:** Partial (`SubAgents.ts` still simulated)
- **Provider/billing:** Mostly complete (`billing.getStatus` real)
- **Browser extension:** Early/stub state
- **Session/cloud management:** Partial
- **Interfaces (CLI/WebUI):** Broad coverage, mixed end-to-end fidelity
- **Integration protocols / advanced features:** Mixed; many capabilities present in code but not fully exposed in router/UI

## Verified Gaps (Code-Evidenced)

### P0 — Critical UX Wiring

- Remaining P0 follow-up:
	- Validate these newly wired widgets against live backend runtime in an integrated dashboard smoke pass.
- Knowledge dashboard integrity issue:
	- Core duplicate/`@ts-ignore` issues resolved; follow-up is TRPC type regeneration cleanup across app.

### P1 — Backend Realism / Persistence Gaps

- `packages/core/src/agents/SubAgents.ts` now dispatches to real researcher/coder execution paths (remaining realism gaps are in deeper knowledge/research services).
- `packages/core/src/services/DeepResearchService.ts` and `KnowledgeService.ts` include placeholder assumptions.
- `packages/ui/src/components/kanban-board.tsx` drag/drop status is local-only (no persistent mutation path).

### P2 — Delivery-Impact Technical Debt

- `apps/web/src/utils/TRPCProvider.tsx` now resolves endpoint via `NEXT_PUBLIC_TRPC_URL` with local/prod-safe fallback behavior.
- `packages/core/src/lib/trpc-core.ts` returns `any` from `getMcpServer()` (wide cast usage across routers).
- `packages/core/src/installer/AutoConfig.ts` explicitly marked stubbed/ported.
- TODO debt remains in `MemoryManager.ts` and `ProjectTracker.ts`.

## Documentation Drift to Resolve

- `CHANGELOG.md` contains `2.6.2` entries while `VERSION.md` is `2.6.1`.
- Some prior “fixed” wording actually reflects stabilization via placeholders, not full functionality restoration.

## Canonical Doc Flow

- `ROADMAP.md` → phase status and priority ladder
- `docs/DETAILED_BACKLOG.md` → implementation queue with acceptance criteria
- `HANDOFF_ANTIGRAVITY.md` / `handoff.md` → handoff package for implementor models

Any non-canonical or stale documents should be treated as archival unless synchronized with the files above.
