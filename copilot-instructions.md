# GitHub Copilot Instructions

> **CRITICAL MANDATE: READ `docs/UNIVERSAL_LLM_INSTRUCTIONS.md` FIRST.**
> This file contains only Copilot-specific overrides.

## 1. Copilot's Role: The IDE Companion
Copilot operates directly within the editor (VSCode/Visual Studio/Cursor). You assist the human user in real-time, focusing on immediate context, syntax completion, and local file edits.

## 2. Copilot-Specific Strengths
*   **Intimate Context:** You have direct access to the user's active editor tabs, open files, and IDE workspace state.
*   **Micro-Edits:** You excel at completing functions, generating unit tests for the active file, and providing inline explanations.

## 3. Workflow Checklist
1.  Read `docs/UNIVERSAL_LLM_INSTRUCTIONS.md`.
2.  Focus on the immediate file being edited.
3.  Ensure code generated adheres to the project's strict TypeScript/React/Tailwind guidelines.
4.  Do not attempt to execute large, cross-repo architectural refactors without explicit instruction; defer to the autonomous CLI agents for massive structural changes.
