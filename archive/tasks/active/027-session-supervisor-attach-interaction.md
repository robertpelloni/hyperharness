# Task: Session Supervisor Attach and Interaction

## Context
Session supervision is implemented and stable for supervision basics (creation, restart, termination), but operators lack a clear, shipped attach/interaction path for connecting to running supervised sessions. This creates ambiguity about whether supervised sessions are extensible for end-to-end workflows or are terminal/isolated.

## Scope
- `packages/core/src/supervisor/SessionSupervisor.ts` attach mechanics
- `packages/core/src/routers/sessionRouter.ts` attach API contract
- `apps/web/src/app/dashboard/session/**` attach/interaction UI
- focused tests for attach workflows and error paths
- dashboard labels/status indicators showing attach readiness

## Requirements
1. Define a clear attach path: what it means to "attach" to a supervised session (stdio/signal handling).
2. Implement attach HTTP/WebSocket endpoint if operators are expected to use it.
3. Provide truthful attach readiness in dashboard session cards (green/ready, yellow/pending, red/unavailable).
4. Document attach non-goal if supervisor is intentionally terminal/isolated.
5. Ensure restart/recovery behavior is explicit and operator-visible during attach.

## Acceptance Criteria
- [ ] Attach path is explicitly defined (either functional feature or documented non-goal)
- [ ] Dashboard session cards show true attach readiness status
- [ ] Supervisor attach tests pass (happy path + error cases)
- [ ] Attach documentation reflects current implementation truthfully
- [ ] Session restart during active attach is safe and visibly handled
- [ ] `CHANGELOG.md` updated with attach story clarity

## Out of Scope
- Multi-user concurrent attach scenarios
- Advanced debugging integration (breakpoints, profiling)
- Supervisor clustering / multi-node attach

## Notes
- Task 026 partially addressed attach reliability at the supervisor level, but dashboard integration and operator-facing attach path still need clarity.
- Attach readiness is distinct from session health; a healthy session might not be attachable if the supervisor is in a transitional state.
