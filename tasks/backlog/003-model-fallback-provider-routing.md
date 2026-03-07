# Task 003: Model Fallback & Provider Routing

## Context
Provider fallback keeps Borg usable when quotas or rate limits bite. It also unlocks Borg's own internal development loops by making model switching automatic.

## Scope
- Files: `packages/core/providers/**`, provider wiring, shared provider schemas, and dashboard state needed for quota visibility
- Tests: `packages/core/providers/__tests__/fallback-chain.test.ts`, `quota-tracker.test.ts`, `strategy.test.ts`, `auth.test.ts`

## Requirements
1. Support multiple provider auth modes needed by the current codebase.
2. Track provider quota/availability state in a normalized way.
3. Fail over automatically on quota exhaustion or rate limiting without losing request flow.
4. Support routing strategies: cheapest, best, and round-robin.
5. Support task-type-aware routing rules.

## Acceptance Criteria
- [ ] At least three providers can be configured/authenticated through the normalized provider layer
- [ ] Automatic fallback succeeds on quota/rate-limit errors
- [ ] Routing strategy changes affect candidate ordering deterministically
- [ ] Quota data is exposed for dashboard consumption
- [ ] All listed test files exist and pass
- [ ] No `@ts-ignore` added
- [ ] `CHANGELOG.md` updated

## Out of Scope
- Billing analytics beyond what is required for quota and routing decisions
- Reinventing provider SDKs instead of adapting them
- Do not create new task files
- STOP when criteria are met
