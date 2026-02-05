# Changelog

All notable changes to this project will be documented in this file.

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
