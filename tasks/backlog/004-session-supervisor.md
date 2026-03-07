# Task 004: Session Supervisor

## Context
Borg should supervise coding tools, not impersonate them. The session supervisor makes external CLI sessions reliable, resumable, and isolated.

## Scope
- Files: `packages/core/supervisor/**`, related persistence/state wiring, worktree helpers, and dashboard contracts required for session lifecycle visibility
- Tests: `packages/core/supervisor/__tests__/spawn.test.ts`, `restart.test.ts`, `health.test.ts`, `worktree.test.ts`, `session-persist.test.ts`

## Requirements
1. Spawn external CLI sessions from Borg with typed/session metadata.
2. Detect crashes and restart with exponential backoff.
3. Persist session state across Borg restarts.
4. Support terminal attach/log viewing contracts.
5. Isolate parallel coding sessions through git worktrees when they target the same repo.

## Acceptance Criteria
- [ ] Session start/stop/restart flows are implemented through the supervisor layer
- [ ] Crash recovery is automatic and test-covered
- [ ] Persisted sessions can be restored after Borg restart
- [ ] Worktree isolation is enforced for parallel coding sessions
- [ ] All listed test files exist and pass
- [ ] No `@ts-ignore` added
- [ ] `CHANGELOG.md` updated

## Out of Scope
- Reimplementing the full UX of Aider, Claude Code, or OpenCode
- Building dashboard chrome before backend session events exist
- Do not create new task files
- STOP when criteria are met
