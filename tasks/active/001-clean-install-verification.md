# Task 001: Clean Install Verification

## Context
Borg 1.0 is blocked if a stranger cannot clone, install, and boot the dashboard quickly. This task makes the quick-start path truthful before more feature work piles on.

## Scope
- Files: `package.json`, `pnpm-workspace.yaml`, `docker-compose.yml`, `README.md`, setup scripts, and any build config directly blocking install/start
- Tests: installation/build verification commands and any minimal regression coverage needed for the touched setup code

## Requirements
1. `pnpm install` completes from a fresh clone without manual intervention.
2. `docker-compose up --build` brings up the dashboard with the documented URL.
3. `README.md` quick-start steps reflect reality and match the validated commands.
4. Failures are fixed at the source instead of documented away.

## Acceptance Criteria
- [ ] Fresh install path is validated end to end on the current repo state
- [ ] Documented startup path matches observed behavior
- [ ] Any blocking setup defects are fixed in code/config, not just noted
- [ ] Test file exists and passes, if setup logic is changed in a testable module
- [ ] No `@ts-ignore` added
- [ ] `CHANGELOG.md` updated

## Out of Scope
- Shipping new product features unrelated to install/start
- Refactoring healthy infrastructure just because it looks lonely
- Do not create new task files
- STOP when criteria are met
