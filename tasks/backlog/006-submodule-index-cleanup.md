# Task 006: Submodule Index Cleanup

## Context
Borg's live submodule registry now lists only 5 tracked submodules in `.gitmodules`, but the Git index still contains 849 gitlinks. That leaves 844 orphaned gitlinks that break `git submodule status`, distort inventory tooling, and keep legacy submodule sprawl alive in the repository state.

## Scope
- Files: `.gitmodules`, `docs/SUBMODULES.md`, `docs/SUBMODULE_DASHBOARD.md`, supporting cleanup scripts under `scripts/**`, and any narrowly required inventory/report docs
- Tests: validation script coverage where practical for submodule inventory or cleanup tooling; command-level verification that `git submodule status` succeeds after cleanup

## Requirements
1. Identify and remove or intentionally reconcile gitlinks that are still tracked in the index but no longer registered in `.gitmodules`.
2. Preserve the currently approved live registry entries:
   - `external/MetaMCP`
   - `packages/MCP-SuperAssistant`
   - `packages/opencode-autopilot`
   - `submodules/mcpproxy`
   - `submodules/litellm`
3. Ensure submodule tooling operates from the same source of truth so `git submodule status` no longer fails on orphaned paths.
4. Keep the cleanup scoped to registry/index consistency; do not re-expand the old submodule corpus.

## Acceptance Criteria
- [ ] `git submodule status` succeeds from the repository root
- [ ] All orphaned gitlinks not present in `.gitmodules` are either removed from the index or intentionally re-registered with written justification
- [ ] `docs/SUBMODULES.md` and `docs/SUBMODULE_DASHBOARD.md` reflect the post-cleanup registry state
- [ ] Cleanup tooling or verification steps are documented in-repo
- [ ] No `@ts-ignore` added
- [ ] `CHANGELOG.md` updated

## Out of Scope
- Reintroducing archived submodule trees or restoring the old 800+ reference corpus
- Adding new submodules beyond the currently approved five live entries without human approval
- Broad ecosystem reclassification beyond what is required for registry/index consistency
- Do not create new task files
- STOP when criteria are met
