# Borg TODO List

_Short-term tasks, bug fixes, and feature implementations. For long-term goals, see ROADMAP.md._

## Critical — Phase O Convergence & v1.0.0
- [ ] **Zero-Error Console Policy**: Fix remaining "Access to storage" and "tRPC SSE connection" errors in browser/VSCode extensions. (Partially done: safeStorage fallback implemented).
- [ ] **Audit all package.json lockfiles**: Ensure exact version pinning for v1.0.0 release to prevent dependency drift.
- [ ] **Finalize Dashboard Data-Binding**: Verify all 59+ pages in `apps/web` are fully wired to real tRPC services (brain, cloud-dev, chronicle, experts, evolution, reader).
- [ ] **Submodule Dashboard**: Implement the UI for monitoring and updating submodules directly from the dashboard.
- [ ] **Mobile App Hardening**: Link `@borg/mobile` React Native wireframes to real WebSocket/tRPC endpoints for remote monitoring.

## Critical — Build & Type Safety
- [x] Fix `council/index.ts` — expose `members`/`updateMembers` at top level for `providers/routing/page.tsx` build fix.
- [ ] Ensure `council.json` config file exists at `packages/core/config/council.json` with default members for fresh installs.
- [ ] Clean up orphaned `routers/councilRouter.ts` (top-level) — its logic is now inlined into `council/index.ts`.

## UI & Dashboard (apps/web)
- [x] Implement Dashboard page/panel listing submodules, versions, dates, and exact repository locations.
- [x] Polish the Roundtable/Council UI enough to render live sessions, history, and Smart Pilot compatibility.
- [x] Wire a non-destructive "Add Borg as MCP server" action into the Integration Hub.
- [x] Implement "Code Mode" escape hatch interface in the dashboard.
- [x] Create detailed usage/billing subpanels tracking credit balances per provider.
- [x] Build Unified Directory merging installed servers + backlog links.
- [x] **Marketplace page** — currently calls `MarketplaceService` which has TODO stubs for MeshService peer queries. Wire up actual community entries or catalog data.
- [ ] **Config page** (`/dashboard/config`) — minimal wrapper around DirectorConfig component. Expand with the full system settings surface (themes, notifications, data retention, etc.).
- [x] **Workshop page** — delegates to `@borg/ui` `WorkshopPage`. Verify this component is fully implemented in the UI package.

## Orchestration & Models
- [x] Implement robust model fallback logic with quota-aware cascading.
- [x] Ensure auto-start/restart logic handles all 11 CLI harness types.
- [x] Implement OAuth logic for subscribing to premium models.
- [x] Fully wire up the Council debate to SmartPilot.
- [x] **MeshRouter** (`meshRouter`) — currently commented out in `trpc.ts`. The P2P mesh networking layer for distributed Borg instances is scaffolded but disabled. Re-enable when MeshService is stable.

## Documentation & Instructions
- [x] Restore "ALL CAPS MAD SCIENCE" to README.md.
- [x] Create/Update MEMORY.md, DEPLOY.md, CHANGELOG.md.
- [ ] **Universal Instruction Unification**: Refine CLAUDE.md, GEMINI.md, GPT.md to all reference `docs/UNIVERSAL_LLM_INSTRUCTIONS.md` as the primary source of truth.
- [x] Ensure `docs/SUBMODULES.md` is current with all 7 active submodules.
- [x] Update `docs/VERSION_LOCATIONS.md` after version bump to 0.99.5.

## Mad Science & Ambitious Pivots (IDEAS.md)
- [ ] **Rust Micro-Kernel**: Prototype the core event loop in Rust for sub-millisecond response times.
- [ ] **P2P Hive Mind**: Activate fact-gossiping across the Borg mesh.
- [ ] **Bobcoin Integration**: Implement micro-payments for tool/skill marketplace.
- [ ] **3D Brain Dashboard**: Visualization of cognitive activity using Three.js.
