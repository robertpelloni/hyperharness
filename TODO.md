# Borg TODO List

_Short-term tasks, bug fixes, and feature implementations. For long-term goals, see ROADMAP.md._

## UI & Dashboard (apps/web)
- [ ] Implement Dashboard page/panel listing all submodules, versions, build numbers, dates, and their exact locations in the directory structure.
- [ ] Add explicit explanations of the project directory structure layout directly in the UI.
- [ ] Polish the Roundtable/Council UI: Hook up all "Placeholder" sections (Session Grid, History Table) to actual live state.
- [ ] Ensure all features are thoroughly documented within the UI (tooltips, labels, descriptions).
- [ ] Implement "Code Mode" escape hatch interface in the dashboard.
- [ ] Create detailed usage/billing subpanels tracking credit balances per provider.

## Orchestration & Models
- [ ] Implement robust model fallback logic: when quota limit errors occur, automatically switch to the next appropriate model/provider (e.g., Gemini 3 Pro -> Codex 5.3 -> Opus 4.6).
- [ ] Ensure auto-start/restart logic can handle opencode, superai-cli, codebuff, codemachine, factory droid, codex, claude code, and gemini CLI instances.
- [ ] Implement OAuth logic for subscribing to premium models (Copilot Premium Plus, ChatGPT Plus, Claude Pro, Google AI Plus).
- [ ] Fully wire up the Council debate to the `SmartPilot` trigger so sessions can autonomously self-correct based on consensus.

## MCP Substrate & Proxies
- [ ] Improve MCP router startup: load last-known-good configuration to report immediately, regardless of whether the core is fully ready.
- [ ] Implement dynamic progressive tool disclosure (show only 5-6 permanent meta tools initially, auto-load others on high confidence).
- [ ] Add explicit tool semantic search / tool RAG.
- [ ] Implement TOON format parsing and MCP traffic inspection panels.
- [ ] Build environment variable and secrets management tool inside the dashboard.

## Memory & RAG
- [ ] Connect the memory subsystem to Google Docs, Gmail, and Google Drive for seamless RAG.
- [ ] Implement memory browser extension endpoints (saving web memories, parsing DOM).
- [ ] Integrate NotebookLM-style features for source-grounded, citation-backed answers.

## Automation & Scripts
- [ ] Create automated script to check versions of all submodules, update them, and generate the UI data feed for the submodule dashboard.
- [ ] Standardize and document the 7-step merge protocol for handling upstream forks.

## Documentation
- [ ] Create/Update `MEMORY.md` with ongoing observations.
- [ ] Create/Update `DEPLOY.md` with explicit deployment/startup instructions.
- [ ] Update `CHANGELOG.md` for this sprint.
- [ ] Refine `AGENTS.md`, `GEMINI.md`, `CLAUDE.md`, `GPT.md` to reference the Universal LLM Instructions and the new rigorous documentation protocols.
