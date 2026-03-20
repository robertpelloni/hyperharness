# Handoff – 2026-03-20

## Delta update (latest)

### VS Code task hygiene (Vitest reporter + duplicate task cleanup)
- Fixed `verify: mcp discovery guards` in `.vscode/tasks.json` by removing unsupported Vitest flag `--reporter=basic`.
- Confirmed task execution now passes:
  - `pnpm exec vitest run packages/core/test/mcpDiscoveryFailureHandling.test.ts packages/core/test/HealerReactor.test.ts` ✅
- Removed duplicate task entries for `core: typecheck` / `web: build-webpack` to reduce task picker ambiguity and keep task IDs deterministic.
- Re-validated `core: typecheck-marker` task after cleanup (`CORE_TSC_OK`) ✅
- Follow-up cleanup removed additional duplicate variants:
  - `web: tsc verify current 2`
  - `root: build extensions validate 2`
  - `verify: mcp discovery guards clean`, `clean 2`, `clean 3`
- Re-validated after second pass:
  - `web: tsc verify current` ✅
  - `verify: mcp discovery guards` ✅

### Published catalog → managed install workflow (operator action)
- Added `catalog.installFromRecipe` admin mutation in `packages/core/src/routers/catalogRouter.ts`.
- New install mutation behavior:
  - allows install only for `validated`/`certified` published entries,
  - requires an active recipe,
  - maps recipe template fields into managed `mcp_servers` create input with transport-aware logic,
  - enforces required secrets (`required_secrets`) before create,
  - auto-generates a unique, safe server name if needed.
- Updated `apps/web/src/app/dashboard/registry/page.tsx` with a row-level `Install` action (download icon), enabled only for validated/certified rows.
- Install success now refreshes both published catalog list and managed MCP server list caches.
- Improved `/dashboard/registry` install UX with a guided modal:
  - opens from row-level Install action,
  - loads active recipe detail via `catalog.get`,
  - renders inputs for `required_env` defaults and `required_secrets`,
  - submits merged values into `catalog.installFromRecipe`.

### Verification (post-change)
- `pnpm -C apps/web exec tsc --noEmit --pretty false` ✅
- `pnpm -C apps/web build --webpack` ✅
- Edited files report no diagnostics in workspace Problems for this slice ✅
- Guided install modal wiring in `apps/web/src/app/dashboard/registry/page.tsx` reports no TypeScript diagnostics ✅

### MCP registry/catalog workflow unification (UI)
- Updated `apps/web/src/app/dashboard/mcp/registry/page.tsx` to include a `Published Catalog Intelligence` panel.
- Added live metrics sourced from `trpc.catalog.stats` (total, validated, broken, updated 24h).
- Added an explicit deep-link CTA to `/dashboard/registry` to route operators from quick-install templates to the provenance/validation-first catalog surface.
- Clarified split responsibilities in-page to reduce operator ambiguity between install templates and verified catalog workflows.

### Verification (post-change)
- `pnpm -C apps/web exec tsc --noEmit --pretty false` ✅
- `pnpm -C apps/web build --webpack` ✅
- `pnpm -C packages/core exec tsc --noEmit --pretty false` ✅

### Catalog freshness metric completed
- Implemented a real `recentlyUpdated` metric for the published catalog API:
  - Added `countRecentlyUpdated(hours = 24)` to `packages/core/src/db/repositories/published-catalog.repo.ts`
  - Uses `updated_at >= now - 24h` on `published_mcp_servers` (no placeholder counting)
- Wired `packages/core/src/routers/catalogRouter.ts` `stats` procedure to use this real metric instead of the previous temporary total-count fallback.
- Updated `/dashboard/registry` UI (`apps/web/src/app/dashboard/registry/page.tsx`) to display a new `Updated 24h` stat card.

### Verification (post-change)
- `pnpm -C packages/core exec tsc --noEmit --pretty false` ✅
- `pnpm -C apps/web exec tsc --noEmit --pretty false` ✅
- `pnpm -C apps/web build --webpack` ✅

### Notes for next session
- Build is currently green; no reproducible web build failure remains from earlier task context.
- Next likely high-value step: unify `/dashboard/mcp/registry` and `/dashboard/registry` operator experience (navigation clarity + workflow handoff), then add richer catalog detail drilldown (sources + latest validation run per row).

## What changed in this session

### Runtime warning stabilization
- Patched `packages/ui/src/components/ui/use-toast.tsx` to avoid re-subscribing toast listeners on every state change (`useEffect` deps changed from `[state]` to `[]`).
- This targets the React warning about state updates on unmounted/not-yet-mounted components in the layout/provider tree.

### Verification done
- Type checks and targeted tests passed in prior steps.
- A traced `apps/web` production build completed successfully after earlier intermittent failure, indicating stale/transient artifact behavior rather than deterministic compile break.

### Strategic direction clarified
- Canonical priority is now **MCP Registry Intelligence P0**:
  1. published catalog persistence,
  2. ingestion/provenance/dedupe,
  3. config intelligence,
  4. validation harness,
  5. operator install workflow.
- Updated top authoritative blocks in `ROADMAP.md` and `TODO.md` (both now dated `2026-03-20`).

## External research ingested

High-signal references reviewed:
- MCPProxy docs (`config-file`, `upstream-servers`, `docker-isolation`, `sensitive-data-detection`)
- MetaMCP docs (`quickstart`, concepts)
- MCPHub docs (`quickstart`, endpoint model)
- mcp-use/manufact docs (`agent/client/server overview`, TypeScript agent config)
- ContextForge roadmap page for ecosystem-scale validation and catalog governance patterns

Two repos still unresolved by extractor:
- `https://github.com/robertpelloni/pluggedin-app`
- `https://github.com/robertpelloni/mcp-tool-chainer`

## Recommended next coding step (immediate)

Implement DB primitives in core for published-catalog + validation history:
- `published_mcp_servers`
- `published_mcp_server_sources`
- `published_mcp_validation_runs`
- `published_mcp_config_recipes`

Then wire one ingestion adapter end-to-end (single source) and one UI read path.

## Risks / watchouts
- Avoid mixing published and installed states in one status field.
- Preserve provenance/raw payload references for auditability.
- Ensure secrets are modeled as requirements/placeholders only (never persisted as plaintext values).
