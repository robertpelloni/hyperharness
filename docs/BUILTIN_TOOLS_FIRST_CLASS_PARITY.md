# Built-In Tools First-Class Parity (Borg)

> Date: 2026-03-18  
> Scope requested: Antigravity, Kiro, Windsurf, Gemini CLI, Codex CLI, Claude Code, Copilot CLI, VS Code/Copilot IDE, Cursor, OpenCode

## Objective

Borg should expose **native, first-class tools** aligned with the tools and workflows models are already trained to use in leading IDE/CLI agents.

This document is the active, execution-ready parity tracker:

- **Verified**: direct evidence from source/docs already in this repository.
- **Unverified**: placeholders to be validated from official upstream docs/repositories.
- **Target**: Borg-native equivalents that preserve naming/semantics where feasible.

---

## Verification Legend

- ✅ **Verified** — captured from source material available in-repo (`archive/docs/*`).
- 🟡 **Draft** — likely accurate but not yet source-locked in current branch docs.
- ❌ **Unverified** — pending authoritative upstream source confirmation.

---

## Competitor Tool Inventory (Current)

## 1) OpenCode (SST) — ✅ Verified

Source: `archive/docs/RESEARCH_COMPETITORS.md`

Built-in tools observed:

- `ls(path, recursive)`
- `grep(pattern, path, include, literal_text)`
- `read(file_path)`
- `view(file_path, offset, limit)`
- `edit(file_path, ...)`
- `patch(file_path, diff)`
- `diagnostics(file_path)`
- `bash(command, timeout)`
- `fetch(url, format)`
- `agent(prompt)`

Borg parity status:

- Already strong coverage via workspace search/read/edit, diagnostics, shell, web fetch, and sub-agent orchestration.
- Gap: preserve **tool-shape compatibility aliases** (OpenCode-style names/signatures) for drop-in behavior.

---

## 2) Gemini CLI — 🟡 Draft

Source: `archive/docs/RESEARCH_COMPETITORS.md` (summary-level)

Capabilities noted:

- Search grounding
- File read/write
- Shell execution
- Web fetch
- Trusted folders / execution policies

Borg parity status:

- Strong baseline present.
- Gap: codify Gemini-compatible policy modes and context-file conventions as first-class profiles.

---

## 3) GitHub Copilot CLI (`gh copilot`) — ✅ Verified (high-level)

Source: `archive/docs/RESEARCH_COMPETITORS.md`

Commands/tools noted:

- `suggest`
- `explain`

Borg parity status:

- Existing capability can emulate these through wrappers.
- Gap: add explicit command aliases and output formatting parity.

---

## 4) Aider — ✅ Verified (high-level)

Source: `archive/docs/RESEARCH_COMPETITORS.md`, `archive/docs/CLI_FEATURE_PARITY.md`

Notable built-ins:

- Repo map (AST/index overview)
- Git-first edit loop (diff/commit ergonomics)
- Lint/test assisted fixing

Borg parity status:

- Most primitives exist.
- Gap: dedicated `repo_map`-style tool contract with deterministic output shape.

---

## 5) Codex CLI — 🟡 Draft

Source: `archive/docs/CLI_FEATURE_PARITY.md` (feature summary)

Capabilities referenced:

- Cloud task integration
- Parallel task execution
- Reviewer/diff analysis mode
- Multimodal inputs
- Skills/runbooks

Borg parity status:

- Strong overlap (tasks, diffs, sub-agents, skills).
- Gap: Codex-style compatibility profile + command idioms.

---

## 6) Claude Code — ❌ Unverified (tool-level)

Current status:

- Mentioned across archived docs but no authoritative tool schema captured in active docs.

Required next step:

- Capture official command/tool list and permission model from upstream source docs.

---

## 7) Cursor — ❌ Unverified (tool-level)

Current status:

- Mentioned in archived strategy docs; no validated built-in tool contract in active docs.

Required next step:

- Source-lock exact built-ins from official docs/changelog/API references.

---

## 8) VS Code + Copilot IDE Agent — ❌ Unverified (tool-level)

Current status:

- No authoritative, versioned tool manifest captured in active docs.

Required next step:

- Capture tool categories + operation semantics from official VS Code/Copilot documentation.

---

## 9) Windsurf — ❌ Unverified (tool-level)

Current status:

- No source-locked built-in tool manifest captured.

Required next step:

- Document official built-ins, mode controls, and approval model.

---

## 10) Kiro — ❌ Unverified (tool-level)

Current status:

- Mentioned in competitor ecosystem docs; no exact first-party built-in tool schema captured.

Required next step:

- Capture Kiro official tool/mode contract.

---

## 11) Antigravity — ❌ Unverified (tool-level)

Current status:

- Mentioned in archived integration docs; no exact built-in tool inventory captured in active docs.

Required next step:

- Capture official tool list + invocation semantics.

---

## Borg First-Class Tool Contract (Target)

The following Borg-native tool families should be treated as **immutable first-class surfaces** (stable names + schemas):

1. **Workspace I/O**
   - `list_files`, `read_file`, `write_file`, `apply_patch`, `search_text`, `search_semantic`
2. **Code Intelligence**
   - `get_diagnostics`, `find_references`, `go_to_definition`, `rename_symbol`
3. **Execution**
   - `run_command`, `run_task`, `get_terminal_output`, `kill_terminal`
4. **Web/Docs Grounding**
   - `fetch_webpage`, `browser_open`, `browser_read`, `browser_action`
5. **Planning/Autonomy**
   - `plan_update`, `subagent_run`, `checkpoint_create`, `checkpoint_restore`
6. **Memory**
   - `memory_save`, `memory_recall`, `memory_forget`
7. **Change Control**
   - `propose_change`, `approve_change`, `plan_mode`, `build_mode`

Compatibility profiles (aliases/adapters):

- `opencode` profile
- `copilot-cli` profile
- `codex-cli` profile
- `gemini-cli` profile
- `cursor` profile
- `claude-code` profile

---

## Implementation Milestones

### Milestone A — Tool Schema Canonicalization (P0)

- [ ] Freeze Borg canonical tool names and JSON schemas.
- [ ] Add semver for tool contracts.
- [ ] Add backward-compatible alias layer for OpenCode/Copilot CLI/Codex CLI style names.

### Milestone B — Compatibility Profiles (P1)

- [ ] Add profile switch: `--profile opencode|copilot-cli|codex-cli|gemini-cli|cursor|claude-code`.
- [ ] Map profile tool calls to canonical Borg tools.
- [ ] Snapshot tests for request/response shape equivalence.

### Milestone C — Approval & Safety Parity (P1)

- [ ] Standardize permission scopes (`once`, `session`, `always`, `deny`).
- [ ] Add per-tool/per-arg allowlists.
- [ ] Emit audit trail for every privileged call.

### Milestone D — Repo Map + Context Contract (P1)

- [ ] Add deterministic `repo_map` tool.
- [ ] Support universal context files: `AGENTS.md`, `CLAUDE.md`, `OpenCode.md`, `context.md`, `BORG.md`.

### Milestone E — Evidence Locking (P0 ongoing)

- [ ] For each competitor above, add:
  - exact built-in tool names,
  - parameter schemas,
  - return shapes,
  - approval model,
  - source URLs + commit/date pin.

---

## Acceptance Criteria for “First-Class”

Borg is first-class when all are true:

1. **Deterministic contracts**: stable, versioned schemas.
2. **Compatibility layer**: alias profiles pass golden tests.
3. **Operational parity**: file/code/shell/web/planning/memory primitives all native.
4. **Safety parity**: permission gates + auditability equivalent or better.
5. **Documentation parity**: one-page matrix with source-pinned competitor evidence.

---

## Immediate Next Work Queue

1. Build `docs/BUILTIN_TOOLS_EVIDENCE_LOCK.md` with per-platform source pins.
2. Generate golden compatibility fixtures for OpenCode and Copilot CLI first.
3. Implement profile adapter scaffolding in Borg CLI/core.
4. Add CI gate: fail if canonical tool schema changes without version bump.
