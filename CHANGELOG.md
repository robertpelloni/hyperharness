# Changelog

All notable changes to this project will be documented in this file.

## [2.4.0] - 2026-02-07

### Added

- **Phase 45: Knowledge Assimilation & Dashboard Expansion**:
  - **Knowledge Dashboard**: `/dashboard/knowledge` now acts as a central hub for all external intelligence (MCPs, Skills, Docs).
  - **Ingestion Pipeline**: Automated script (`scripts/ingest_resources.ts`) to assimilate vast resource lists into structured knowledge.
  - **Submodule Service**: New `SubmoduleService` to health-check and sync the massive plugin ecosystem.
- **Phase 44: The Immune System (Self-Healing)**:
  - **HealerReactor**: Autonomous error detection loop (Terminal -> Sensor -> Healer).
  - **Auto-Fix**: System can now self-diagnose and patch simple code errors without human intervention.
  - **Immune Dashboard**: `/dashboard/healer` visualizes active infections and immune responses.

## [2.3.0] - 2026-02-06

### Added

- **Phase 36 (Release & Observability)** (In Progress):
  - Enhanced Submodules Dashboard with project structure visualization and detailed versioning.
- **Phase 35 (Standardization & Polish)**:
  - **Code Hygiene**: Resolved 150+ lint errors in frontend.
  - **Library UI**: Launched `/dashboard/library` showing Prompts and Skills.
- **Phase 34 (Evolution - The Darwin Protocol)**:
  - **DarwinService**: Mutation engine for A/B testing prompts.
  - **Evolution UI**: `/dashboard/evolution` for managing agent experiments.
- **Phase 33 (Self-Correction - The Healer)**:
  - **HealerService**: Automated error analysis and fix suggestion loop.
- **Phase 32 (Security & Governance - The Guardian)**:
  - **PolicyService**: RBAC and Scope enforcement for tools.
  - **Security Dashboard**: `/dashboard/security` for audit logs and policy management.
- **Phase 31 (Deep Research - The Scholar)**:
  - **DeepResearchService**: Multi-step recursive research agent.

## [2.2.1] - 2026-02-05

### Fixed

- **Development Experience**: Resolved persistent console clearing issues during `pnpm run dev`.
  - Added `--preserveWatchOutput` to all `tsc` watch scripts.
  - Added `clearScreen: false` to Vite config.
  - Monkey-patched `console.clear` in Next.js config.
  - Updated CLI to use `tsx watch --clear-screen=false`.
- **Extension Bridge**: Verified active WebSocket server on port 3001.

## [2.2.0] - 2026-02-04

### Added

- **Phase 23 (Deep Data Search)**:
  - Created `Indexer` and `CodeSplitter` in `@borg/memory`.
  - Added AST-based symbol extraction for TypeScript.
  - Added semantic chunking logic.
  - Exposed `memory_index_codebase` and `memory_search` tools in `MCPServer`.
  - Implemented `VectorStore.addDocuments` for batch processing.

## [2.1.0] - 2026-02-04

### Added

- **Master MCP Server Architecture (Phase 21)**: Borg now aggregates downstream MCP servers (`git`, `filesystem`, etc.) via `MCPAggregator`.
- **Stdio Client**: Native integration for spawning and controlling local MCP tools.
- **Unified Documentation**: Centralized all agent instructions into `docs/LLM_INSTRUCTIONS.md`.
- **Vision Document**: Published `VISION.md` outlining the "Neural Operating System" goal.
- **Centralized Versioning**: Project version is now master-controlled by the `VERSION` file.

### Changed

- **Config**: `borg.config.json` is now the primary configuration point for adding tools.
- **Routing**: Tool calls are now prefixed (e.g., `git_commit`) to allow namespace isolation between multiple servers.

## [1.7.0] - 2026-02-03

### Added

- **Core Infrastructure**: Registered LSP, Code Mode, and Plan Mode tools in MCPServer.
- **Verification**: Added `CoreInfra.test.ts`.
