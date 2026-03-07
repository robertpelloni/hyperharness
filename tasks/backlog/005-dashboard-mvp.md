Proceed # Task 005: Dashboard MVP

## Context
Borg needs one control surface that makes the router, provider state, and supervised sessions visible within a minute of startup. This is the face users will judge first.

## Scope
- Files: `apps/web/**`, supporting shared UI packages, and backend contracts strictly required for overview, MCP, sessions, and providers panels
- Tests: integration coverage for dashboard-backed system status where practical, especially `tests/integration/mcp-to-dashboard.test.ts`, `fallback-e2e.test.ts`, and `session-lifecycle.test.ts`

## Requirements
1. Provide Overview, MCP Router, Sessions, and Providers panels.
2. Surface real status data from the backend instead of placeholder optimism.
3. Show quota bars, session state, server health, and traffic inspector output.
4. Keep the initial dashboard path understandable for a first-time operator.

## Acceptance Criteria
- [ ] The four v1 dashboard panels exist and are wired to real backend data
- [ ] Operators can inspect server health, provider quota state, and session lifecycle from the UI
- [ ] Integration coverage exists for the dashboard-facing system flows listed in scope
- [ ] No placeholder-only widgets are added without live backing data
- [ ] No `@ts-ignore` added
- [ ] `CHANGELOG.md` updated

## Out of Scope
- Nice-to-have pages outside the v1.0 control surface
- Re-skinning the entire app without improving the core workflows
- Do not create new task files
- STOP when criteria are met
