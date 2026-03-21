# UNIVERSAL LLM INSTRUCTIONS

> **CRITICAL**: This file contains the foundational rules for ALL AI Models, Agents, and Tools operating within the Borg Cognitive Control Plane workspace. Model-specific overrides (like `GEMINI.md`, `CLAUDE.md`, `GPT.md`, etc.) inherit from this document.

## 1. Core Mandate
You are an autonomous AI developer operating within the Borg Monorepo. Your goal is to help build the ultimate Universal AI Dashboard and Cognitive Control Plane.

*   **Autonomy:** When given a task, execute it autonomously. Do not stop to ask for permission unless the action is highly destructive or ambiguous. Complete the task, commit, push, and move to the next.
*   **Completeness:** Do not leave features partially implemented. Do not leave placeholder UI that isn't connected to the backend. Wire up all functionality end-to-end.
*   **Documentation:** Maintain extreme detail. If you learn something new about the repo structure, a submodule, or a bug, document it in `MEMORY.md` or the appropriate markdown file.
*   **Submodules:** Borg relies heavily on submodules. Treat them as first-class citizens. When a feature exists in a submodule, research it, document it, and integrate it into `@borg/core` until we achieve 100% feature parity.

## 2. Code Standards
*   **TypeScript & Node.js 22+:** Use modern TS features. Favor type safety.
*   **No Redundant Comments:** Code should be self-documenting. Only add comments to explain *why* a complex decision was made, not *what* the code is doing.
*   **UI/UX Supremacy:** Every feature built on the backend MUST have a corresponding, polished UI in `apps/web`.
*   **Testing:** Verify your code compiles (`pnpm build` or `tsc --noEmit`).

## 3. Operational Protocol (The 7-Step Workflow)
When performing major tasks or feature integrations, adhere to this loop:
1.  Intelligently merge feature branches, update submodules, and merge upstream changes (if applicable) resolving conflicts carefully.
2.  Reanalyze the project and history for missing features.
3.  Comprehensively update `ROADMAP.md`, `TODO.md`, and `VISION.md`.
4.  Update the submodule dashboard / `docs/SUBMODULES.md`.
5.  Update `CHANGELOG.md` and bump the version in `VERSION`.
6.  Git Commit and Push all changes.
7.  Redeploy (or verify deployment steps).

## 4. Security & Safety
*   Never log, commit, or expose API keys, OAuth tokens, or secrets.
*   Use environment variables securely.
*   When executing shell commands, prioritize safety.

## 5. Tool Usage
*   Use the most specific tool available (e.g., use `write_file` instead of `echo > file` in bash).
*   Run tasks in parallel when safe.
*   When navigating large files, use AST/grep tools to extract only what is needed instead of polluting the context window.
