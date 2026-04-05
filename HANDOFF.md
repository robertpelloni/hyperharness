# Handoff

## What was done
- Added a Go-native foundation bootstrap for a Pi-derived harness:
  - `foundation/pi`
  - `foundation/compat`
  - `foundation/assimilation`
  - `cmd/foundation.go`
- Extended the bootstrap into a **truthful native runtime baseline**:
  - native `read`, `write`, `edit`, `bash` handlers
  - evented runtime execution
  - JSONL-backed session persistence with create/list/load/fork
  - CLI execution and session management commands under `foundation`
- Added the first native **Aider-inspired repo map baseline**:
  - ranked source scanning
  - lightweight symbol extraction
  - deterministic `<repo_map>` output
  - `hypercode foundation repomap`
- Strengthened repo-map quality with first **graph-ranking groundwork**:
  - files that reference symbols defined in other files now push ranking weight toward those defining files
  - mention-based ranking still works, but cross-file code relationships now matter too
- Routed top-level tool registration closer to the new foundation:
  - `tools.Registry` now exposes exact-name native Pi-compatible tools from `foundation/pi`
  - `agent.Agent` now registers per-tool schemas instead of one placeholder schema
  - legacy `tools/repomap.go` now delegates to `foundation/repomap`
- Added the first HyperCode/Borg adapter seam:
  - `foundation/adapters/hypercode.go` exposes assimilation status, memory context, provider status, MCP config visibility, and adjacent HyperCode repo discovery
  - `foundation/adapters/providers.go` exposes current provider/model visibility, detected providers, Ollama model discovery, and provider-route selection groundwork
  - `foundation/adapters/provider_routing.go` provides shared route-selection logic for CLI and HTTP surfaces
  - `foundation/adapters/provider_execution.go` provides shared provider execution-preparation behavior for CLI, HTTP, and provider stubs
  - `foundation/adapters/mcp_config.go` and `foundation/adapters/mcp.go` expose MCP config parsing, server discovery, tool hints, route hints, mediated tool-call preparation, and configured-server startup seams
  - `hypercode foundation adapters` now inspects that seam directly
  - `agent.Agent` now incorporates adapter-derived system context into its system prompt
- Migrated more top-level surfaces onto the foundation/adapter layer:
  - `mcp/client.go`, `mcp/manager.go`, and `mcp/config.go` now consume adapter-backed behavior
  - `cmd/foundation_http.go` provides reusable foundation-backed HTTP helpers
  - `cmd/serve.go` now exposes `/api/v1/foundation/*` endpoints and routes `/fs/read` through the foundation `read` tool
  - foundation-backed MCP HTTP helper and route surfaces now expose MCP tool listing and mediated call preparation
  - foundation-backed provider helper and route surfaces now expose provider visibility, route selection, and execution preparation
  - foundation-backed orchestration helper and route surfaces now expose plan generation
  - `hypercode foundation providers status/select/prepare` now exposes provider routing groundwork from the CLI
  - `hypercode foundation plan` now exposes orchestration planning from the CLI
  - `agents/provider_stub.go` and `agents/provider.go` now consume provider execution-preparation hints
  - `agents/director.go` and `agent/orchestrator.go` now consume `foundation/orchestration` planning primitives
  - `orchestrator/webhooks.go`, daemon sweep planning, and autodrive objective generation now consume foundation orchestration helpers
  - `tui/slash.go` now exposes foundation-backed `/plan`, `/repomap`, `/providers`, `/adapters`, and `/mcp`
  - `tui/foundation_bridge.go` now routes normal prompt and shell proposal flows through foundation-aware helpers
  - `tui/slash.go` now exposes foundation-backed `/plan` and `/repomap` slash commands
- Added deeper verification coverage:
  - snapshot-style tests for baseline tool results
  - top-level agent schema registration test
  - HyperCode/Borg adapter seam test
  - provider adapter seam test
  - provider-route selection test
  - provider execution-preparation test
  - MCP adapter seam tests and top-level MCP package tests
  - foundation-backed HTTP helper tests
  - MCP mediation helper tests
  - orchestration planner, daemon planner, webhook planner, and migrated director/orchestrator tests
  - foundation-backed TUI slash-command tests
  - foundation-backed TUI prompt/shell helper tests
  - provider CLI smoke checks
  - foundation plan CLI smoke checks
  - TUI provider/adapter introspection smoke coverage via tests
- Added comprehensive `docs/ai/` documentation for requirements, design, planning, implementation, and testing.
- Added `docs/ai/design/upstream-toolchain-analysis.md` summarizing the imported upstream systems and assimilation strategy.
- Fixed the duplicate SQLite driver registration issue in `orchestrator/queue.go` and `orchestrator/vectors.go` by removing redundant `modernc.org/sqlite` imports.

## Validation completed
- `gofmt -w cmd/foundation.go cmd/foundation_http.go cmd/foundation_http_test.go cmd/serve.go foundation/adapters/hypercode.go foundation/adapters/hypercode_test.go foundation/adapters/providers.go foundation/adapters/providers_test.go foundation/adapters/mcp_config.go foundation/adapters/mcp.go foundation/adapters/mcp_test.go foundation/compat/types.go foundation/compat/catalog.go foundation/compat/default_catalog.go foundation/compat/catalog_test.go foundation/pi/foundation.go foundation/pi/foundation_test.go foundation/pi/runtime_types.go foundation/pi/runtime.go foundation/pi/runtime_test.go foundation/pi/session.go foundation/pi/session_test.go foundation/pi/tool_parity_test.go foundation/pi/tool_snapshot_test.go foundation/pi/tools_native.go foundation/assimilation/inventory.go foundation/assimilation/summary.go foundation/assimilation/inventory_test.go foundation/repomap/repomap.go foundation/repomap/repomap_test.go tools/registry.go tools/repomap.go tools/registry_test.go agent/agent.go agent/agent_test.go mcp/config.go mcp/client.go mcp/manager.go mcp/client_test.go mcp/manager_test.go mcp/mcphost.go mcp/mcphost_test.go orchestrator/vectors.go orchestrator/queue.go`
- `go test ./foundation/... ./cmd ./orchestrator ./tools ./agent ./mcp`
- `go run . foundation adapters`
- `go run . foundation tools`
- `go run . foundation inventory`
- `go run . foundation repomap --dir foundation --max-files 5`
- `go run . foundation session create --name smoke`
- `go run . foundation exec --tool write --input '{"path":"smoke.txt","content":"hello"}'`

## Important repo state notes
- The broader repository still has unrelated baseline failures under full `go test ./...`:
  - `aider/tests/fixtures/languages/go/test.go` has an unused import.
  - `mcp/mcphost_test.go` is out of sync with the host API.
- Many existing Go packages still contain placeholder or aspirational parity logic and should be migrated to the new `foundation/*` packages over time.

## Recommended next steps
1. Continue routing remaining top-level placeholder orchestration/tool surfaces onto the new `foundation/pi` runtime.
2. Expand verified result-shape and snapshot tests for `read`, `write`, `edit`, and `bash` plus CLI/HTTP smoke coverage.
3. Deepen `foundation/repomap` from graph-ranking groundwork toward fuller Aider-style graph ranking and richer edit engines.
4. Expand `foundation/adapters` from visibility, route-selection, and execution-preparation seams into richer provider routing, memory, and richer MCP runtime adapters backed by HyperCode/Borg.
5. Migrate TUI and orchestration code to the new truthful foundation instead of placeholder parity claims, with special attention to adapter-backed execution paths.

## Additional work completed on 2026-04-04
- Reconciled the repository module path back to the actual local import graph:
  - `go.mod` now uses `module github.com/robertpelloni/hypercode`
- Ran dependency reconciliation successfully via `go mod tidy`.
- Verified the **canonical native foundation track** instead of the duplicate experimental internal track:
  - `foundation/pi`
  - `foundation/compat`
  - `foundation/assimilation`
  - `foundation/repomap`
  - `foundation/adapters`
  - `foundation/orchestration`
- Expanded the verified native Pi-compatible tool surface from **4** tools to **7** tools:
  - `read`
  - `write`
  - `edit`
  - `bash`
  - `grep`
  - `find`
  - `ls`
- Added native input types for those tools in `foundation/pi/runtime_types.go`.
- Expanded `foundation/pi/BuiltinTools()` JSON schema descriptors for those tools.
- Expanded `foundation/pi/DefaultToolHandlers()` and added native handlers:
  - `executeGrepTool`
  - `executeFindTool`
  - `executeLSTool`
- Expanded `foundation/compat/DefaultCatalog()` so exact-name Pi contracts now include the same 7 tools.
- Updated and extended tests to match the broader native parity surface:
  - `foundation/pi/foundation_test.go`
  - `foundation/pi/runtime_test.go`
  - `foundation/pi/tool_snapshot_test.go`
  - `foundation/compat/catalog_test.go`
- Added comprehensive architectural reconciliation doc:
  - `docs/analysis/FOUNDATION_RECONCILIATION_2026-04-04.md`

## Latest validation completed
- `go test ./foundation/pi/...`
- `go test ./foundation/...`
- Verified snapshot coverage specifically for the expanded native Pi tool surface.

## Updated architectural conclusion
- The repo contains both a newer `internal/*` Go track and an older but more integrated `foundation/*` Go track.
- The correct near-term implementation truth should be the **existing `foundation/*` track**, because it is already:
  - test-backed,
  - wired into the repo’s existing CLI foundation surfaces,
  - aligned with `foundation/compat`, and
  - better suited to honest parity expansion.
- Recommendation remains to **merge good ideas from `internal/*` into `foundation/*`**, not the other way around.

## Updated next steps
1. Keep `go test ./foundation/...` green at all times.
2. Continue expanding exact Pi parity in `foundation/pi` beyond the current 7 native tools.
3. Promote `foundation/compat` as the single source of truth for model-facing tool contracts.
4. Deepen Aider-style repo map and edit-engine work in `foundation/repomap` and `foundation/pi`.
5. Continue bridging HyperCode via `foundation/adapters` rather than duplicating control-plane responsibilities inside the harness.

## Additional work completed later on 2026-04-04
- Expanded native `foundation/pi` session semantics significantly beyond the earlier basic persistence layer.
- `foundation/pi/session.go` now supports richer Pi-style session metadata and entry shapes:
  - versioned session metadata (`Version`, `ParentSession`)
  - model changes
  - thinking-level changes
  - compaction entries
  - branch summary entries
  - custom entries
  - custom message entries
  - session-info entries
  - label entries
- Added native session append helpers:
  - `AppendThinkingLevelChange`
  - `AppendModelChange`
  - `AppendCompaction`
  - `AppendBranchSummary`
  - `AppendCustomEntry`
  - `AppendCustomMessage`
  - `AppendSessionInfo`
  - `AppendLabelChange`
- Added session inspection and tree helpers:
  - `GetEntry`
  - `GetChildren`
  - `GetBranch`
  - `GetLabel`
  - `GetSessionName`
- Added native session context reconstruction via:
  - `BuildSessionContext(sessionID, leafID)`
  - `SessionContext`
- Added runtime convenience wrappers in `foundation/pi/runtime.go` for the richer session APIs.
- Added and expanded tests in `foundation/pi/session_test.go` to verify:
  - versioned session creation
  - parent-session preservation on fork
  - labels
  - session-info naming
  - model/thinking changes
  - custom entries and custom messages
  - compaction entries
  - branch reconstruction
  - session context reconstruction
- Added detailed analysis doc:
  - `docs/analysis/PI_SESSION_PARITY_TRANCHE_2026-04-04.md`

## Latest validation completed after session work
- `gofmt -w foundation/pi/session.go foundation/pi/runtime.go foundation/pi/session_test.go`
- `go test ./foundation/pi/...`
- `go test ./foundation/...`

## Updated recommendation after session tranche
1. Keep `foundation/*` green and truthful.
2. Continue Pi-native parity next with explicit branch/leaf management and branch-summary semantics.
3. After that, deepen CLI/TUI/session-mode behavior on top of the now-richer native session model.

## Additional work completed later on 2026-04-04 (branch/leaf tranche)
- Added explicit active-leaf tracking to `foundation/pi` sessions via `SessionMetadata.LeafID`.
- Updated `AppendEntry()` to parent new entries to the active leaf when present instead of always using the latest entry in file order.
- Updated `AppendEntry()` to advance the active leaf to the newly appended entry.
- Updated `Fork()` to default to the active leaf when choosing a fork point.
- Added native branch/leaf helpers on `SessionStore`:
  - `GetLeafID`
  - `Branch`
  - `ResetLeaf`
- Updated `GetBranch()` to default to the active leaf.
- Updated `BuildSessionContext()` to default to the active leaf instead of only latest-entry order.
- Added runtime wrappers in `foundation/pi/runtime.go`:
  - `GetLeafID`
  - `BranchSession`
  - `ResetSessionLeaf`
- Updated runtime tool-run persistence so tool executions attach to and advance the active leaf correctly.
- Added explicit branch/leaf verification to `foundation/pi/session_test.go`.
- Added detailed analysis doc:
  - `docs/analysis/PI_BRANCH_LEAF_PARITY_TRANCHE_2026-04-04.md`

## Latest validation after branch/leaf tranche
- `gofmt -w foundation/pi/session.go foundation/pi/runtime.go foundation/pi/session_test.go`
- `go test ./foundation/pi/...`
- `go test ./foundation/...`

## Updated recommendation after branch/leaf tranche
1. Keep `foundation/*` green and canonical.
2. Implement native branch-summary/common-ancestor helpers next on top of explicit leaf semantics.
3. Then layer `/tree`-style CLI/TUI workflows over the truthful foundation behavior.

## Additional work completed later on 2026-04-04 (branch-summary groundwork tranche)
- Added native branch-summary preparation substrate to `foundation/pi`.
- Added `BranchSummaryPreparation` with:
  - `TargetID`
  - `OldLeafID`
  - `CommonAncestorID`
  - `EntriesToSummarize`
- Added native common-ancestor computation:
  - `GetCommonAncestor(sessionID, firstLeafID, secondLeafID)`
- Added native branch-summary preparation:
  - `PrepareBranchSummary(sessionID, targetID)`
- Added native branch-and-append-summary helper:
  - `BranchWithSummary(sessionID, targetID, summary, details)`
- Added runtime wrappers in `foundation/pi/runtime.go`:
  - `GetCommonAncestor`
  - `PrepareBranchSummary`
  - `BranchWithSummary`
- Added explicit verification in `foundation/pi/session_test.go` for:
  - common ancestor correctness
  - abandoned branch entry collection
  - branch-summary attachment on the destination branch
  - `fromId` preservation
- Added detailed analysis doc:
  - `docs/analysis/PI_BRANCH_SUMMARY_GROUNDWORK_TRANCHE_2026-04-04.md`

## Latest validation after branch-summary groundwork
- `gofmt -w foundation/pi/session.go foundation/pi/runtime.go foundation/pi/session_test.go`
- `go test ./foundation/pi/...`
- `go test ./foundation/...`

## Updated recommendation after branch-summary groundwork
1. Keep `foundation/*` green and canonical.
2. Implement budget-aware branch-summary preparation and structured summary-generation hooks next.
3. Then build truthful `/tree`-style CLI/TUI behavior on top of the now-real branch substrate.

## Additional work completed later on 2026-04-04 (budget-aware branch-summary tranche)
- Expanded `BranchSummaryPreparation` to include:
  - `SerializedConversation`
  - `FileOps`
  - `EstimatedTokens`
  - `MaxTokens`
- Added `BranchSummaryFileOps` with cumulative:
  - `ReadFiles`
  - `ModifiedFiles`
- Added budget-aware preparation API:
  - `PrepareBranchSummaryWithBudget(sessionID, targetID, maxTokens)`
- Added serialization helpers in `foundation/pi/session.go`:
  - `serializeConversation`
  - `serializeSingleEntry`
  - `flattenToolResultText`
  - `truncateSummaryText`
  - `estimateSerializedTokens`
- Added budget trimming helper:
  - `trimEntriesToBudget`
- Added cumulative file-op extraction helpers:
  - `collectBranchFileOps`
  - `collectFileOpsFromEntry`
- Added structured summary hook:
  - `DefaultStructuredSummaryTemplate(prep)`
- Added runtime wrapper:
  - `PrepareBranchSummaryWithBudget`
- Important design choice: file-op tracking is cumulative across the full abandoned branch path even when the serialized conversation window is budget-trimmed.
- Added detailed analysis doc:
  - `docs/analysis/PI_BUDGET_AWARE_BRANCH_SUMMARY_TRANCHE_2026-04-04.md`

## Latest validation after budget-aware branch-summary tranche
- `gofmt -w foundation/pi/session.go foundation/pi/runtime.go foundation/pi/session_test.go`
- `go test ./foundation/pi/...`
- `go test ./foundation/...`

## Updated recommendation after budget-aware branch-summary tranche
1. Keep `foundation/*` green and canonical.
2. Add native summary-generation hooks and/or provider-backed summary generation next.
3. Then begin surfacing truthful `/tree`-style CLI/TUI branch workflows on top of the new summarization substrate.

## Additional work completed later on 2026-04-04 (summary-generation hooks tranche)
- Added pluggable `SummaryGenerator` interface in `foundation/pi/summary.go`.
- Added `SummaryGeneratorFunc` adapter for lightweight generator injection.
- Added default local `DeterministicSummaryGenerator` that produces structured branch summaries without any network/provider dependency.
- Added deterministic structured-summary inference helpers:
  - `inferGoal`
  - `inferCompletedItems`
  - `inferNextSteps`
  - `inferCriticalContext`
- Added session-store generation APIs:
  - `GenerateBranchSummary(ctx, prep, generator)`
  - `BranchWithGeneratedSummary(ctx, sessionID, targetID, maxTokens, generator, details)`
- Added runtime wrappers in `foundation/pi/runtime.go`:
  - `GenerateBranchSummary`
  - `BranchWithGeneratedSummary`
- Added verification in new file:
  - `foundation/pi/summary_test.go`
- Added detailed analysis doc:
  - `docs/analysis/PI_SUMMARY_GENERATION_HOOKS_TRANCHE_2026-04-04.md`

## Latest validation after summary-generation hooks tranche
- `gofmt -w foundation/pi/summary.go foundation/pi/runtime.go foundation/pi/summary_test.go`
- `go test ./foundation/pi/...`
- `go test ./foundation/...`

## Updated recommendation after summary-generation hooks tranche
1. Keep `foundation/*` green and canonical.
2. Mirror the same hook/generator architecture for compaction preparation next.
3. Then start surfacing the now-real summary-generation workflow in CLI/TUI branch navigation (`/tree`-style behavior).

## Additional work completed later on 2026-04-04 (compaction generation hooks tranche)
- Added native `CompactionPreparation` carrying:
  - `LeafID`
  - `EntriesToSummarize`
  - `SerializedConversation`
  - `FileOps`
  - `EstimatedTokens`
  - `TokensBefore`
  - `FirstKeptEntryID`
  - `KeepRecentTokens`
- Added session-store preparation APIs:
  - `PrepareCompaction(sessionID)`
  - `PrepareCompactionWithBudget(sessionID, keepRecentTokens)`
- Added pluggable `CompactionSummaryGenerator` interface.
- Extended `DeterministicSummaryGenerator` to implement deterministic local compaction summary generation.
- Added session-store generation APIs:
  - `GenerateCompactionSummary(ctx, prep, generator)`
  - `CompactWithGeneratedSummary(ctx, sessionID, keepRecentTokens, generator, details)`
- Added runtime wrappers in `foundation/pi/runtime.go`:
  - `PrepareCompaction`
  - `PrepareCompactionWithBudget`
  - `GenerateCompactionSummary`
  - `CompactWithGeneratedSummary`
- Expanded `foundation/pi/summary_test.go` to verify compaction generation and append behavior.
- Added detailed analysis doc:
  - `docs/analysis/PI_COMPACTION_GENERATION_HOOKS_TRANCHE_2026-04-04.md`

## Latest validation after compaction generation hooks tranche
- `gofmt -w foundation/pi/summary.go foundation/pi/runtime.go foundation/pi/summary_test.go foundation/pi/session.go`
- `go test ./foundation/pi/...`
- `go test ./foundation/...`

## Updated recommendation after compaction generation hooks tranche
1. Keep `foundation/*` green and canonical.
2. Start surfacing the now-truthful summary-generation flows in CLI/TUI behavior next (`/tree` and compaction command flows).
3. After that, consider provider-backed summary generation as an optional enhancement on top of the deterministic baseline.

## Additional work completed later on 2026-04-04 (foundation summary CLI tranche)
- Added new `foundation summary` command group in `cmd/foundation.go`.
- Added `foundation summary branch` command with:
  - `--session`
  - `--target`
  - `--max-tokens`
- Added `foundation summary compact` command with:
  - `--session`
  - `--keep-recent-tokens`
- Added reusable helper functions in `cmd/foundation_http.go`:
  - `prepareFoundationBranchSummary`
  - `generateFoundationBranchSummary`
  - `prepareFoundationCompaction`
  - `generateFoundationCompaction`
- Expanded `cmd/foundation_http_test.go` to verify the new summary helper flows.
- Added detailed analysis doc:
  - `docs/analysis/FOUNDATION_SUMMARY_CLI_TRANCHE_2026-04-04.md`

## Latest validation after foundation summary CLI tranche
- `gofmt -w cmd/foundation.go cmd/foundation_http.go cmd/foundation_http_test.go`
- `go test ./cmd ./foundation/...`

## Updated recommendation after foundation summary CLI tranche
1. Keep `foundation/*` green and canonical.
2. Surface the same truthful summary-generation flows in the TUI next.
3. Then consider HTTP exposure and provider-backed generation as layered enhancements, not replacements for the deterministic baseline.

## Additional work completed later on 2026-04-04 (TUI summary slash tranche)
- Added minimal canonical foundation session bridge into the TUI model via `foundationSessionID`.
- Added lazy foundation session creation helper:
  - `ensureFoundationSession`
- Added TUI→foundation session append helpers:
  - `appendFoundationUserText`
  - `appendFoundationAssistantText`
- Added TUI summary bridge helpers:
  - `buildFoundationCompactionDisplay`
  - `buildFoundationBranchSummaryDisplay`
  - `parseSummaryArgs`
- Added new TUI slash commands:
  - `/fsession`
  - `/summary-compact [keepRecentTokens]`
  - `/summary-branch <targetEntryId> [maxTokens]`
- Updated `/help` output in the TUI to document those new slash commands.
- Expanded `tui/slash_test.go` to verify the new TUI summary surfaces.
- Added detailed analysis doc:
  - `docs/analysis/TUI_SUMMARY_SLASH_TRANCHE_2026-04-04.md`

## Latest validation after TUI summary slash tranche
- `gofmt -w tui/chat.go tui/foundation_bridge.go tui/slash.go tui/slash_test.go`
- `go test ./tui ./cmd ./foundation/...`

## Updated recommendation after TUI summary slash tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Begin implementing truthful `/tree`-style TUI navigation next on top of the now-real branch/session/summary substrate.
3. Keep provider-backed generation optional and layered above the deterministic baseline.

## Additional work completed later on 2026-04-04 (TUI tree navigation tranche)
- Added foundation-backed TUI tree display helper:
  - `buildFoundationTreeDisplay`
- Added foundation-backed TUI tree switch helper:
  - `switchFoundationTreeDisplay`
- Added TUI slash command surfaces:
  - `/tree`
  - `/tree <targetEntryId> [maxTokens]`
- Updated TUI `/help` output to document the new `/tree` behavior.
- Expanded `tui/slash_test.go` to verify:
  - tree display output
  - tree switch output
  - structured branch-summary content during switch
  - actual leaf movement after switch
- Added detailed analysis doc:
  - `docs/analysis/TUI_TREE_NAVIGATION_TRANCHE_2026-04-04.md`

## Latest validation after TUI tree navigation tranche
- `gofmt -w tui/foundation_bridge.go tui/slash.go tui/slash_test.go`
- `go test ./tui ./cmd ./foundation/...`

## Updated recommendation after TUI tree navigation tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Expand `/tree` from minimal command-driven navigation into a richer TUI branch explorer next.
3. Continue resisting duplicated semantics in the UI layer; richer UX should sit on top of the same verified foundation behavior.

## Additional work completed later on 2026-04-04 (TUI tree explorer tranche)
- Added runtime wrappers in `foundation/pi/runtime.go` for richer tree inspection:
  - `GetChildren`
  - `GetLabel`
  - `GetSessionName`
- Expanded TUI tree display to include:
  - session name
  - active leaf
  - per-entry child counts
  - per-entry labels
- Added TUI foundation bridge helpers:
  - `buildFoundationChildrenDisplay`
  - `setFoundationLabel`
- Added new TUI slash commands:
  - `/tree-children <entryId>`
  - `/label <entryId> <label>`
- Updated TUI `/help` output to document the richer tree-explorer surfaces.
- Expanded `tui/slash_test.go` to verify:
  - label-setting output
  - labeled tree output
  - tree-children output
  - child-branch previews
- Added detailed analysis doc:
  - `docs/analysis/TUI_TREE_EXPLORER_TRANCHE_2026-04-04.md`

## Latest validation after TUI tree explorer tranche
- `gofmt -w foundation/pi/runtime.go tui/foundation_bridge.go tui/slash.go tui/slash_test.go`
- `go test ./tui ./cmd ./foundation/...`

## Updated recommendation after TUI tree explorer tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Begin building an interactive TUI tree/branch selector next on top of the now-richer runtime surfaces.
3. Keep all richer UX as a thin layer over canonical runtime truth, not a competing state model.

## Additional work completed later on 2026-04-04 (TUI tree selector tranche)
- Added selector state to the TUI model:
  - `foundationTreeSelection []string`
- Added foundation bridge helpers:
  - `buildFoundationTreeSelectionDisplay`
  - `switchFoundationTreeSelection`
- Added new TUI slash commands:
  - `/tree-select`
  - `/tree-go <index> [maxTokens]`
- Updated TUI `/help` output to document those selector-style tree commands.
- Expanded `tui/slash_test.go` to verify:
  - selector output
  - selector state persistence in the TUI model
  - selector-driven branch switching
  - actual leaf movement after selection-based switch
- Added detailed analysis doc:
  - `docs/analysis/TUI_TREE_SELECTOR_TRANCHE_2026-04-04.md`

## Latest validation after TUI tree selector tranche
- `gofmt -w tui/chat.go tui/foundation_bridge.go tui/slash.go tui/slash_test.go foundation/pi/runtime.go`
- `go test ./tui ./cmd ./foundation/...`

## Updated recommendation after TUI tree selector tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Consider an actual interactive in-TUI tree/selector view next, now that the selector-state groundwork exists.
3. Continue layering UX over canonical session/summary truth rather than introducing parallel TUI-specific state machines.

## Additional work completed later on 2026-04-04 (interactive tree selector tranche)
- Added lightweight selector state to the TUI model:
  - `foundationTreeSelection []string`
- Added foundation bridge helpers:
  - `buildFoundationTreeSelectionDisplay`
  - `switchFoundationTreeSelection`
- Added new TUI slash commands:
  - `/tree-select`
  - `/tree-go <index> [maxTokens]`
- Updated TUI `/help` output to document selector-based tree workflows.
- Expanded `tui/slash_test.go` to verify:
  - selector output rendering
  - selector ID caching in the TUI model
  - selector-driven leaf movement
  - structured branch-summary output after selection-based switching
- Added detailed analysis doc:
  - `docs/analysis/TUI_INTERACTIVE_TREE_SELECTOR_TRANCHE_2026-04-04.md`

## Latest validation after interactive tree selector tranche
- `gofmt -w tui/chat.go tui/foundation_bridge.go tui/slash.go tui/slash_test.go foundation/pi/runtime.go`
- `go test ./tui ./cmd ./foundation/...`

## Updated recommendation after interactive tree selector tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Begin a truly cursor-driven in-TUI tree browser next on top of the selector groundwork.
3. Preserve the rule that the UI can cache selector state but must not become the owner of session/branch truth.

## Additional work completed later on 2026-04-04 (cursor tree browser tranche)
- Added browser state to the TUI model:
  - `browserActive`
  - `browserItems`
  - `browserIndex`
- Added `TreeBrowserItem` in `tui/foundation_bridge.go`.
- Added browser helpers:
  - `buildFoundationTreeBrowser`
  - `renderTreeBrowser`
  - `openSelectedTreeBrowser`
- Added new TUI slash command:
  - `/tree-browser`
- Updated TUI keyboard handling so that when browser mode is active:
  - `Up` / `Down` move selection
  - `Enter` performs a real branch switch through the canonical runtime
  - `Esc` closes the browser without quitting the app
- Updated TUI view rendering to show browser content when browser mode is active.
- Expanded `tui/slash_test.go` to verify browser activation, cursor movement, selection, branch switching, and browser close behavior.
- Added detailed analysis doc:
  - `docs/analysis/TUI_CURSOR_TREE_BROWSER_TRANCHE_2026-04-04.md`

## Latest validation after cursor tree browser tranche
- `gofmt -w tui/chat.go tui/foundation_bridge.go tui/slash.go tui/slash_test.go`
- `go test ./tui ./cmd ./foundation/...`

## Updated recommendation after cursor tree browser tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Add richer browser ergonomics next (preview-before-switch, filtering/folding, better tree layout).
3. Keep browser state transient and UI-local while all session/branch truth remains in the canonical foundation runtime.

## Additional work completed later on 2026-04-04 (browser preview/filter tranche)
- Added browser filter state to the TUI model:
  - `browserFilter string`
- Added browser filter helper:
  - `filterTreeBrowserItems`
- Expanded browser rendering to show:
  - filter string
  - filtered match count
  - selected-entry preview pane
- Expanded browser keyboard handling so that in browser mode:
  - typed runes and spaces update the filter
  - backspace/delete shrinks the filter
  - selection index clamps to the filtered result set
- Expanded `tui/slash_test.go` to verify preview rendering, filter updates, and successful branch switching after filter interactions.
- Added detailed analysis doc:
  - `docs/analysis/TUI_BROWSER_PREVIEW_FILTER_TRANCHE_2026-04-04.md`

## Latest validation after browser preview/filter tranche
- `gofmt -w tui/chat.go tui/foundation_bridge.go tui/slash_test.go`
- `go test ./tui ./cmd ./foundation/...`

## Updated recommendation after browser preview/filter tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Add richer browser ergonomics next (folding/grouping or pre-switch summary-preparation preview).
3. Preserve the architecture rule that filter/preview remain UI-only enhancements and do not become competing branch/session truth.

## Additional work completed later on 2026-04-04 (pre-switch summary preview tranche)
- Expanded `TreeBrowserItem` to carry branch-summary preparation preview fields:
  - `SummaryEntries`
  - `CommonAncestorID`
  - `ReadFilesCount`
  - `ModifiedFilesCount`
- Updated `buildFoundationTreeBrowser` so non-leaf entries query the canonical runtime via `PrepareBranchSummaryWithBudget(..., 128)` and cache preview facts in the browser items.
- Expanded browser preview rendering to display:
  - branch summary entry count estimate
  - common ancestor id
  - read file count
  - modified file count
  - special-case note when already on active leaf
- Updated `tui/slash_test.go` to verify that browser preview includes branch-summary preparation details.
- Added detailed analysis doc:
  - `docs/analysis/TUI_PRE_SWITCH_SUMMARY_PREVIEW_TRANCHE_2026-04-04.md`

## Latest validation after pre-switch summary preview tranche
- `gofmt -w tui/foundation_bridge.go tui/slash_test.go`
- `go test ./tui ./cmd ./foundation/...`

## Updated recommendation after pre-switch summary preview tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Add another browser ergonomics layer next (folding/grouping or explicit confirmation-before-switch).
3. Keep all pre-switch guidance derived from canonical runtime preparation rather than UI-invented heuristics.

## Additional work completed later on 2026-04-04 (confirm-before-switch tranche)
- Added transient confirm state to the TUI browser model:
  - `browserConfirmPending bool`
- Updated browser keyboard behavior so that:
  - first `Enter` arms the selected switch
  - second `Enter` confirms it
  - `Y` confirms while pending
  - `N`, `Esc`, or `Backspace` cancel while pending
  - navigation is paused while confirmation is pending
- Updated browser rendering to show:
  - normal browser instructions in non-confirm state
  - confirm-specific instructions in pending state
  - `[Confirm]` block in the preview pane
- Expanded `tui/slash_test.go` to verify:
  - first Enter enters confirm-pending mode
  - confirm prompt appears in the browser view
  - second Enter performs the real canonical switch
- Added detailed analysis doc:
  - `docs/analysis/TUI_CONFIRM_BEFORE_SWITCH_TRANCHE_2026-04-04.md`

## Latest validation after confirm-before-switch tranche
- `gofmt -w tui/chat.go tui/slash_test.go`
- `go test ./tui ./cmd ./foundation/...`

## Updated recommendation after confirm-before-switch tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Add folding/grouping or richer tree layout next, since navigation, preview, filtering, and confirmation are now all in place.
3. Preserve the rule that confirm/preview are UI workflow layers only and do not alter the canonical switch semantics underneath.

## Additional documentation completed after cursor tree browser tranche
- Added appendix doc:
  - `docs/analysis/TUI_CURSOR_TREE_BROWSER_TRANCHE_2026-04-04_APPENDIX.md`
- Re-validated the cursor-browser tranche remains green:
  - `go test ./tui ./cmd ./foundation/...`

## Additional work completed later on 2026-04-04 (tree layout cues tranche)
- Expanded `TreeBrowserItem` with canonical-structure-derived fields:
  - `ParentID`
  - `Depth`
  - `ChildCount`
- Updated `buildFoundationTreeBrowser` to compute depth and child counts from canonical session relationships.
- Updated browser rendering to show:
  - indentation by depth
  - simple branch glyphs
  - child-count cues
- Expanded preview pane to show:
  - depth
  - child count
- Expanded `tui/slash_test.go` to assert richer tree-layout cues are present in browser output.
- Added detailed analysis doc:
  - `docs/analysis/TUI_TREE_LAYOUT_CUES_TRANCHE_2026-04-04.md`

## Latest validation after tree layout cues tranche
- `gofmt -w tui/foundation_bridge.go tui/slash_test.go`
- `go test ./tui ./cmd ./foundation/...`

## Updated recommendation after tree layout cues tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Add folding/grouping next so larger session trees become easier to navigate.
3. Keep layout cues derived from canonical runtime structure rather than creating a separate UI-owned tree model.

## Additional work completed later on 2026-04-04 (graph-style rendering tranche)
- Expanded `TreeBrowserItem` with additional structure-derived fields:
  - `ParentID`
  - `Prefix`
  - `IsLastChild`
- Added canonical prefix helper:
  - `buildTreePrefix`
- Updated browser rendering to distinguish sibling/last-child structure with graph-style connectors:
  - `├─`
  - `└─`
- Kept rendering fully derived from canonical session relationships rather than a separate UI tree model.
- Expanded `tui/slash_test.go` to assert graph-style connector glyphs appear in browser output.
- Added detailed analysis doc:
  - `docs/analysis/TUI_GRAPH_STYLE_RENDERING_TRANCHE_2026-04-04.md`

## Latest validation after graph-style rendering tranche
- `gofmt -w tui/foundation_bridge.go tui/slash_test.go`
- `go test -v ./tui`
- `go test -v ./cmd ./foundation/...`

## Updated recommendation after graph-style rendering tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Consider richer grouping modes or persistent browser panes next.
3. Continue deriving every visual tree cue from canonical runtime structure rather than UI-invented semantics.

## Additional work completed later on 2026-04-04 (tree folding tranche)
- Added transient browser collapse state:
  - `browserCollapsed map[string]bool`
- Added fold-aware visibility helper:
  - `visibleTreeBrowserItems`
- Updated browser keyboard handling so that in browser mode:
  - `Left` collapses the selected subtree when it has children
  - `Right` expands the selected subtree
- Updated browser rendering to show fold state cues:
  - `[+]` collapsed
  - `[-]` expanded
- Ensured selection index clamps correctly when collapse state changes the visible result set.
- Expanded `tui/slash_test.go` to verify collapse/expand cues through keyboard interaction.
- Added detailed analysis doc:
  - `docs/analysis/TUI_TREE_FOLDING_TRANCHE_2026-04-04.md`

## Latest validation after tree folding tranche
- `gofmt -w tui/chat.go tui/foundation_bridge.go tui/slash_test.go`
- `go test ./tui ./cmd ./foundation/...`

## Updated recommendation after tree folding tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Add richer grouping/graph-style rendering next so the browser becomes even clearer at scale.
3. Keep folding as a UI-only projection over canonical runtime structure, never a competing state model.

## Additional work completed later on 2026-04-04 (graph-style rendering tranche)
- Expanded `TreeBrowserItem` with additional structure-derived fields:
  - `ParentID`
  - `Prefix`
  - `IsLastChild`
- Added canonical prefix helper:
  - `buildTreePrefix`
- Updated browser rendering to distinguish sibling/last-child structure with graph-style connectors:
  - `├─`
  - `└─`
- Kept rendering fully derived from canonical session relationships rather than a separate UI tree model.
- Focused tests passed against the modified TUI/browser surfaces:
  - `go test -run TestTreeBrowserModeNavigation -v ./tui`
  - `go test -run TestProcessSlashCommandTreeExplorerLabelAndChildren -v ./tui`
  - `go test -run TestFoundationSummaryHelpers -v ./cmd`
- Added detailed analysis doc:
  - `docs/analysis/TUI_GRAPH_STYLE_RENDERING_TRANCHE_2026-04-04.md`

## Updated recommendation after graph-style rendering tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Consider richer grouping modes or persistent browser panes next.
3. Continue deriving every visual tree cue from canonical runtime structure rather than UI-invented semantics.

## Additional work completed later on 2026-04-04 (grouped browser tranche)
- Added grouped-mode browser state:
  - `browserGrouped bool`
- Expanded `TreeBrowserItem` with:
  - `GroupKey`
- Added canonical lineage-based group helper:
  - `buildGroupKey`
- Updated browser rendering to insert `[Group] ...` headings when grouped mode is enabled.
- Added browser keyboard toggle:
  - `Tab` toggles grouped mode on/off
- Expanded browser-mode verification to assert grouped rendering appears after toggling.
- Added detailed analysis doc:
  - `docs/analysis/TUI_GROUPED_BROWSER_TRANCHE_2026-04-04.md`

## Latest validation after grouped browser tranche
- `gofmt -w tui/chat.go tui/foundation_bridge.go tui/slash_test.go`
- `go test -run TestTreeBrowserModeNavigation -v ./tui`
- `go test -run TestProcessSlashCommandTreeExplorerLabelAndChildren -v ./tui`
- `go test -run TestFoundationSummaryHelpers -v ./cmd`

## Updated recommendation after grouped browser tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Consider a persistent browser pane/view or stronger branch-head/divergence grouping next.
3. Keep grouped mode as a presentation-only layer over canonical runtime structure.

## Additional work completed later on 2026-04-04 (persistent tree pane tranche)
- Added persistent pane state to the TUI model:
  - `browserPinned bool`
- Added foundation bridge helpers:
  - `pinFoundationTreeBrowser`
  - `unpinFoundationTreeBrowser`
- Added new TUI slash command:
  - `/tree-pane`
- Updated `model.View()` so a pinned browser pane renders alongside the normal prompt when browser mode is not active.
- Expanded `tui/slash_test.go` to verify pin/unpin behavior and pane visibility in the TUI view.
- Added detailed analysis doc:
  - `docs/analysis/TUI_PERSISTENT_TREE_PANE_TRANCHE_2026-04-04.md`

## Latest validation after persistent tree pane tranche
- `gofmt -w tui/chat.go tui/foundation_bridge.go tui/slash.go tui/slash_test.go`
- `go test ./tui ./cmd ./foundation/...`

## Updated recommendation after persistent tree pane tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Consider refresh and stronger pane ergonomics next, or evolve toward a richer split-view browser layout.
3. Keep the pane as a view over canonical runtime truth rather than a parallel session/branch model.

## Additional work completed later on 2026-04-04 (pinned pane refresh tranche)
- Added pinned-pane refresh helper:
  - `refreshPinnedFoundationTreeBrowser`
- Wired refresh into normal TUI message flow so the pinned pane updates after:
  - user text append
  - assistant text append
- Added focused regression coverage for pinned-pane refresh behavior.
- Added detailed analysis doc:
  - `docs/analysis/TUI_PINNED_PANE_REFRESH_TRANCHE_2026-04-04.md`

## Latest validation after pinned-pane refresh tranche
- `gofmt -w tui/chat.go tui/foundation_bridge.go tui/slash_test.go`
- `go test ./tui ./cmd ./foundation/...`

## Updated recommendation after pinned-pane refresh tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Move toward stronger pane ergonomics or split-view browser behavior next, since the pane is now both persistent and live.
3. Keep all pane refreshes driven by canonical runtime reads rather than local UI-only session mutation logic.

## Additional work completed later on 2026-04-04 (focusable pane tranche)
- Added focus state to the TUI model:
  - `browserPinnedFocus bool`
- Added new TUI slash command:
  - `/tree-pane-focus`
- Updated key handling so that when the pane is pinned and focused it supports:
  - Up/Down navigation
  - Left/Right collapse/expand
  - Tab grouping toggle
  - type-to-filter
  - Enter/Y/N confirm-before-switch flow
- Updated the view to show:
  - `[Tree Pane Focused]`
- Ensured pane hide clears focus and pending-confirm state.
- Expanded `tui/slash_test.go` with a focused pane navigation regression.
- Added detailed analysis doc:
  - `docs/analysis/TUI_FOCUSABLE_PANE_TRANCHE_2026-04-04.md`

## Latest validation after focusable pane tranche
- `gofmt -w tui/chat.go tui/slash.go tui/slash_test.go`
- `go test ./tui ./cmd ./foundation/...`

## Updated recommendation after focusable pane tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Improve split-view ergonomics or pane viewport behavior next.
3. Keep the focusable pane as another controller over canonical runtime truth rather than a separate navigation model.

## Additional work completed later on 2026-04-04 (pane viewport tranche)
- Extended the shared browser renderer to accept viewport parameters (`maxVisible`, `title`).
- Added selection-centered viewport windowing logic for the browser renderer.
- Updated the pinned tree pane to render through a bounded viewport instead of always rendering the full visible item list.
- Added pane viewport metadata (`showing=X-Y of Z`) when the visible item set exceeds the pane window.
- Added detailed analysis doc:
  - `docs/analysis/TUI_PANE_VIEWPORT_TRANCHE_2026-04-04.md`

## Latest validation after pane viewport tranche
- `gofmt -w tui/foundation_bridge.go tui/chat.go tui/slash_test.go`
- `go test ./tui ./cmd ./foundation/...`

## Updated recommendation after pane viewport tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Continue improving pane/browser coexistence or move toward stronger split-view layout behavior next.
3. Keep viewport behavior strictly as a rendering concern over canonical runtime-derived browser state.

## Additional work completed later on 2026-04-04 (pane size tranche)
- Added persistent pane height state to the TUI model:
  - `browserPaneHeight int`
- Initialized pane height with a default of `8`.
- Added new TUI slash command:
  - `/tree-pane-size <n>`
- Updated the pinned pane render path to use the adjustable pane height with a safe fallback.
- Added focused regression coverage:
  - `TestProcessSlashCommandTreePaneSize`
- Added detailed analysis doc:
  - `docs/analysis/TUI_PANE_SIZE_TRANCHE_2026-04-04.md`

## Latest validation after pane size tranche
- `gofmt -w tui/chat.go tui/slash.go tui/slash_test.go`
- `go test ./tui`
- `go test ./cmd ./foundation/...`

## Updated recommendation after pane size tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Continue strengthening pane/browser coexistence ergonomics next.
3. Keep pane sizing strictly as a rendering concern over canonical runtime-backed browser state.

## Additional documentation completed after pane size tranche
- Added notes doc:
  - `docs/analysis/TUI_SPLIT_VIEW_ERGONOMICS_NOTES_2026-04-04.md`
- Reconfirmed current pane-size tranche validation:
  - `go test ./tui`
  - `go test ./cmd ./foundation/...`

## Additional work completed later on 2026-04-04 (pane position tranche)
- Added pane position state to the TUI model:
  - `browserPanePosition string`
- Default pane position is now `top`.
- Added new TUI slash command:
  - `/tree-pane-position <top|bottom>`
- Updated `model.View()` so the pinned pane can render above or below the main prompt/history flow.
- Added focused regression coverage:
  - `TestProcessSlashCommandTreePanePosition`
- Added detailed analysis doc:
  - `docs/analysis/TUI_PANE_POSITION_TRANCHE_2026-04-04.md`

## Latest validation after pane position tranche
- `gofmt -w tui/chat.go tui/slash.go tui/slash_test.go`
- `go test ./tui ./cmd ./foundation/...`

## Updated recommendation after pane position tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Continue improving split-view ergonomics next, especially clearer pane/history separation or stronger pane viewport behavior.
3. Keep pane position strictly as a layout concern over canonical runtime-backed browser state.

## Additional work completed later on 2026-04-04 (pane separation tranche)
- Updated the TUI view to render a stronger divider between the pinned tree pane and the main prompt/history flow.
- Applied the same separation behavior for both top and bottom pane positions.
- Expanded pinned-pane verification to assert the presence of the divider.
- Added detailed analysis doc:
  - `docs/analysis/TUI_PANE_SEPARATION_TRANCHE_2026-04-04.md`

## Latest validation after pane separation tranche
- `gofmt -w tui/chat.go tui/slash_test.go`
- `go test ./tui ./cmd ./foundation/...`

## Updated recommendation after pane separation tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Continue improving pane/browser coexistence through stronger viewport controls or more explicit pane-focused interaction affordances.
3. Keep split-view clarity improvements as rendering-only changes layered over canonical runtime-backed state.

## Additional work completed later on 2026-04-04 (viewport controls tranche)
- Added Home/End/PgUp/PgDn navigation for the modal tree browser.
- Added Home/End/PgUp/PgDn navigation for the focused pinned tree pane.
- Tied focused pane PgUp/PgDn movement to the current pane height.
- Added `min` helper for safe index clamping.
- Updated browser guidance text to mention viewport controls.
- Added focused regression coverage:
  - `TestTreePaneViewportControls`
- Added detailed analysis doc:
  - `docs/analysis/TUI_VIEWPORT_CONTROLS_TRANCHE_2026-04-04.md`

## Latest validation after viewport controls tranche
- `gofmt -w tui/chat.go tui/slash_test.go`
- `go test ./tui ./cmd ./foundation/...`

## Updated recommendation after viewport controls tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Continue improving split-view ergonomics, especially stronger active-pane affordances or layout controls.
3. Keep viewport controls as navigation over canonical runtime-derived visible items, not a new tree semantics layer.

## Additional work completed later on 2026-04-04 (pane preview toggle tranche)
- Added persistent pane preview state to the TUI model:
  - `browserPanePreview bool`
- Extended the shared tree renderer to accept an explicit preview-visibility flag.
- Added new TUI slash command:
  - `/tree-pane-preview <on|off>`
- Updated the pinned pane render path to honor the preview toggle while browser mode remains preview-rich.
- Added focused regression coverage for pane-preview toggling.
- Added detailed analysis doc:
  - `docs/analysis/TUI_PANE_PREVIEW_TOGGLE_TRANCHE_2026-04-04.md`

## Latest validation after pane preview toggle tranche
- `gofmt -w tui/chat.go tui/foundation_bridge.go tui/slash.go tui/slash_test.go`
- `go test ./tui ./cmd ./foundation/...`

## Updated recommendation after pane preview toggle tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Continue refining pane/browser coexistence, likely with stronger pane-focused interaction affordances or alternate layout presets.
3. Keep pane preview visibility as a pure rendering concern over canonical runtime-backed browser state.

## Additional work completed later on 2026-04-04 (pane preset tranche)
- Added new TUI slash command:
  - `/tree-pane-preset <compact|detailed>`
- Added `compact` preset:
  - height `6`
  - preview `off`
  - position `bottom`
- Added `detailed` preset:
  - height `12`
  - preview `on`
  - position `top`
- Updated TUI help text to document the preset command.
- Added focused regression coverage for preset application.
- Added detailed analysis doc:
  - `docs/analysis/TUI_PANE_PRESET_TRANCHE_2026-04-04.md`

## Latest validation after pane preset tranche
- `gofmt -w tui/slash.go tui/slash_test.go`
- `go test ./tui ./cmd ./foundation/...`

## Updated recommendation after pane preset tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Continue refining pane/browser coexistence with higher-level operator presets or split-view controls next.
3. Keep presets as bundles of rendering/layout state only, never as alternate session/branch semantics.

## Additional work completed later on 2026-04-04 (pane preset extension tranche)
- Extended `/tree-pane-preset` to support:
  - `navigation`
  - `review`
- Added `navigation` preset:
  - height `10`
  - preview `off`
  - position `bottom`
  - grouped `true`
- Added `review` preset:
  - height `14`
  - preview `on`
  - position `top`
  - grouped `true`
- Expanded pane preset regression coverage to verify all four presets.
- Added detailed analysis doc:
  - `docs/analysis/TUI_PANE_PRESET_EXTENSION_TRANCHE_2026-04-04.md`

## Latest validation after pane preset extension tranche
- `gofmt -w tui/slash.go tui/slash_test.go`
- `go test ./tui ./cmd ./foundation/...`

## Updated recommendation after pane preset extension tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Continue refining split-view ergonomics with stronger pane/browser mode distinctions or additional workflow presets next.
3. Keep presets as bundles of rendering/layout state only, never as alternate session/branch semantics.

## Additional work completed later on 2026-04-04 (mode distinction tranche)
- Updated modal browser title to:
  - `[Foundation Tree Browser :: Modal]`
- Updated passive pinned pane title to:
  - `[Foundation Tree Pane :: Passive]`
- Updated focused pinned pane title to:
  - `[Foundation Tree Pane :: Focused]`
- Expanded TUI regression assertions so the mode-specific titles are explicitly verified.
- Added detailed analysis doc:
  - `docs/analysis/TUI_MODE_DISTINCTION_TRANCHE_2026-04-04.md`

## Latest validation after mode distinction tranche
- `gofmt -w tui/chat.go tui/slash_test.go`
- `go test ./tui ./cmd ./foundation/...`

## Updated recommendation after mode distinction tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Continue refining split-view ergonomics with clearer layout and interaction affordances.
3. Keep mode distinctions explicit in rendering so operators do not have to infer control state from behavior alone.

## Additional work completed later on 2026-04-04 (pane status tranche)
- Added new TUI slash command:
  - `/tree-pane-status`
- Added pane status output covering:
  - pinned
  - focus
  - height
  - position
  - preview
  - grouped
  - filter
- Added focused regression coverage for pane-status output.
- Added detailed analysis doc:
  - `docs/analysis/TUI_PANE_STATUS_TRANCHE_2026-04-04.md`

## Latest validation after pane status tranche
- `gofmt -w tui/slash.go tui/slash_test.go`
- `go test ./tui ./cmd ./foundation/...`

## Updated recommendation after pane status tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Continue improving operator ergonomics around the pane/browser subsystem with complementary reset or quick-toggle controls where useful.
3. Keep pane status as introspection over UI/layout state, not as a competing source of session/branch truth.

## Additional work completed later on 2026-04-04 (reset controls tranche)
- Added new TUI slash commands:
  - `/tree-browser-clear`
  - `/tree-pane-reset`
- `/tree-browser-clear` now clears transient browser state:
  - filter
  - pending confirmation
  - collapse state
  - selection index
- `/tree-pane-reset` now restores pane defaults:
  - height `8`
  - position `top`
  - preview `true`
  - grouped `false`
  - focus `false`
  - pending confirmation `false`
  - filter cleared
  - collapsed state cleared
- Added focused regression coverage for both reset paths.
- Added detailed analysis doc:
  - `docs/analysis/TUI_RESET_CONTROLS_TRANCHE_2026-04-04.md`

## Latest validation after reset controls tranche
- `gofmt -w tui/slash.go tui/slash_test.go`
- `go test ./tui ./cmd ./foundation/...`

## Updated recommendation after reset controls tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Continue improving operator control with higher-level presets or additional quick toggles where they reduce friction.
3. Keep reset/clear behavior confined to UI/layout/control state, never canonical session/branch truth.

## Additional work completed later on 2026-04-04 (pane cycle tranche)
- Added new TUI slash command:
  - `/tree-pane-cycle`
- Implemented state-based cycling across pane presets:
  - `compact`
  - `navigation`
  - `detailed`
  - `review`
- Added focused regression coverage verifying the full cycle order and state transitions.
- Added detailed analysis doc:
  - `docs/analysis/TUI_PANE_CYCLE_TRANCHE_2026-04-04.md`

## Latest validation after pane cycle tranche
- `gofmt -w tui/slash.go tui/slash_test.go`
- `go test ./tui ./cmd ./foundation/...`

## Updated recommendation after pane cycle tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Continue adding operator polish around pane/browser control and recovery where it reduces friction.
3. Keep convenience commands like cycles as lightweight layers over the same underlying pane state bundles.

## Additional work completed later on 2026-04-04 (grouped toggle tranche)
- Added new TUI slash command:
  - `/tree-pane-grouped <on|off|toggle>`
- Added explicit command-level control for grouped rendering outside of keyboard-only browser interaction.
- Added focused regression coverage verifying grouped mode on/toggle behavior.
- Added detailed analysis doc:
  - `docs/analysis/TUI_GROUPED_TOGGLE_TRANCHE_2026-04-04.md`

## Latest validation after grouped toggle tranche
- `gofmt -w tui/slash.go tui/slash_test.go`
- `go test ./tui ./cmd ./foundation/...`

## Updated recommendation after grouped toggle tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Continue adding small operator-control affordances where they reduce friction, especially around pane/browser layout and state recovery.
3. Keep grouped-mode controls as rendering-only state over the same canonical runtime-backed browser model.

## Additional work completed later on 2026-04-04 (pane preview quick-toggle tranche)
- Added new TUI slash command:
  - `/tree-pane-preview-toggle`
- Extended the preview handler to support toggle semantics in addition to explicit on/off.
- Added focused regression coverage for repeated preview toggling.
- Added detailed analysis doc:
  - `docs/analysis/TUI_PANE_PREVIEW_QUICK_TOGGLE_TRANCHE_2026-04-04.md`

## Latest validation after pane preview quick-toggle tranche
- `gofmt -w tui/slash.go tui/slash_test.go`
- `go test ./tui ./cmd ./foundation/...`

## Updated recommendation after pane preview quick-toggle tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Continue adding small operator-polish controls where they reduce repeated pane/browser friction.
3. Keep quick toggles as convenience layers over the same pane layout state rather than separate semantics.

## Additional work completed later on 2026-04-04 (pane position quick-toggle tranche)
- Added new TUI slash command:
  - `/tree-pane-position-toggle`
- Extended the pane-position handler to support toggle semantics in addition to explicit top/bottom values.
- Added focused regression coverage verifying top↔bottom quick toggling.
- Added detailed analysis doc:
  - `docs/analysis/TUI_PANE_POSITION_QUICK_TOGGLE_TRANCHE_2026-04-04.md`

## Latest validation after pane position quick-toggle tranche
- `gofmt -w tui/slash.go tui/slash_test.go`
- `go test ./tui ./cmd ./foundation/...`

## Updated recommendation after pane position quick-toggle tranche
1. Keep `foundation/*`, `cmd`, and `tui` green and aligned to the same canonical runtime.
2. Continue adding small convenience controls where they reduce repeated pane/browser friction.
3. Keep quick position toggles as layout-only controls over the same pane/browser state model.
