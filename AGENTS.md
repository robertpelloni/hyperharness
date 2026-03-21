# BORG Agents Directory

> **CRITICAL MANDATE: READ `docs/UNIVERSAL_LLM_INSTRUCTIONS.md` FIRST.**
> All agents must adhere to the universal guidelines.

## 🤖 The Council of Agents

Borg relies on a diverse set of specialized AI models to accomplish tasks. Each model has its own specific instruction set.

### 1. Gemini (`GEMINI.md`)
*   **Role**: The Architect & Analyst.
*   **Best for**: Massive context processing, full-repo scans, cross-file architectural analysis, resolving complex git conflicts.
*   **Characteristics**: Fast, handles massive inputs, excellent at recognizing high-level patterns.

### 2. Claude (`CLAUDE.md`)
*   **Role**: The Senior Engineer & Code Reviewer.
*   **Best for**: Nuanced TypeScript refactoring, strict type checking, detailed markdown documentation, complex logic debugging.
*   **Characteristics**: Methodical, extremely detail-oriented, prone to writing robust step-by-step plans in `TODO.md` before acting.

### 3. GPT (`GPT.md`)
*   **Role**: The Rapid Implementer.
*   **Best for**: Quick scaffolding, regex generation, simple shell scripts, boilerplate generation.
*   **Characteristics**: Fast execution, highly proficient with standard CLI tools.

### 4. Copilot (`copilot-instructions.md`)
*   **Role**: The IDE Companion.
*   **Best for**: Inline code completion, writing unit tests for active files, micro-edits.
*   **Characteristics**: Operates within the IDE context, constrained to the immediate files open in the editor.

## 🔄 Agentic Workflows

Borg agents follow the 7-Step Workflow outlined in `UNIVERSAL_LLM_INSTRUCTIONS.md`. They are expected to operate autonomously, utilizing subagents where necessary, and continually synchronizing with the central memory and context files (`VISION.md`, `ROADMAP.md`, `MEMORY.md`, etc.).

## 🎓 Skills

Agents can extend their capabilities by activating tools from the `skills/` directory. If an agent encounters a problem outside its immediate training, it should search the skills library and assimilate the necessary knowledge.
