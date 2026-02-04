# borg Roadmap: The Ultimate SuperAI

## Phase 13: Local AI Tool Governance (COMPLETED)
- [x] **Tool Spreadsheet:** Dashboard inventory of local tools.
- [x] **Status Detection:** Auto-discovery of installed CLI tools.
- [x] **Process Monitor:** Watchdog for background services.
- [x] **Hardware Stats:** Real-time resource monitoring.
- [x] **Ecosystem Dashboard:** Auto-generated `SUBMODULES.md` tracking 200+ integrated tools.
- [x] **Massive Submodule Integration:** Import and categorize 100+ external tools.
- [x] **Resource Indexing:** `FEATURES_CATALOG.md` created.
- [x] **Unified Dashboarding:** Consolidated "Mission Control" UI established.

## Phase 14: System Polish & Core UX (CURRENT FOCUS)
**Goal:** Streamline existing systems. Fix "interrupting" bugs. Enhance Dashboard utility.
- [x] **Governance:** Establish `RULES.md` and `DESIGN_STANDARDS.md`.
- [x] **Supervisor Repair:** Fix `InputManager` typing interruptions (Added window target logic).
- [x] **Dashboard Config:** Added `DirectorStatusWidget` and consolidated UI.
- [x] **Documentation:** Sync all docs with actual codebase state.
- [x] **Roadmap Sync:** Verify `ROADMAP.md` is the source of truth.
- [x] **Backlog Processing:** Automate triage of `INBOX_LINKS.md` (300+ items).
- [x] **Dashboard Polish:** Unify styling and layout for "Mission Control" (Completed as Phase 22).

## Phase 22: Director & Dashboard 2.0 (COMPLETED)
- [x] **Director Control Panel:** Full configuration UI.
- [x] **Autonomy & Safety:** Kill switch, loop prevention, auto-submit.
- [x] **Draggable Dashboard:** `dnd-kit` grid layout with persistence.
- [x] **Configurable Persona:** Custom prefixes and behavior settings.

## Phase 15: Deep Research & Absorption (Completed)
**Goal:** Achieve feature parity with the "Definitive Borg Resource Index".
- [x] **Ingest:** Processed `INBOX_LINKS.md` (Batches 1-3).
- [x] **Analyze:** Created `web-research` skill and autonomous script.
- [x] **Absorb:** Injected `treesitter-expert`, `frontend-design`, `mcp-builder` skills.
- [x] **Submodules:** Added key repos (`anthropic-skills`, `awesome-mcp-servers`).

## Phase 15b: Deep Analysis & Stabilization (CURRENT FOCUS)
**Goal:** Comprehensive Audit, Documentation, and Test Stabilization.
- [x] **Project Audit:** Created `docs/PROJECT_STATUS.md` analyzing Frontend/Backend stats.
- [x] **Instruction Consolidation:** Unified `UNIVERSAL_LLM_INSTRUCTIONS.md` (v1.7.1).
- [x] **Legacy Cleanup:** Removed `CORE_INSTRUCTIONS.md`.
- [x] **Extension Wiring:** Connect `apps/extension` to `packages/memory`.

## Phase 14: Deep Code Intelligence (The Foundation)
**Goal:** Establish the robust indexing and execution layer required for high-fidelity coding.

### Robust Indexing
- [x] **Vector Integration:** Enhance `RepoMapService` to use vector embeddings for semantic search (RAG) using `fs-lexical` + `txtai`.
- [x] **Tree-sitter Deep Dive:** Upgrade symbol extraction to support full call graphs and type hierarchies.
- [x] **Graph Service:** Finalize `RepoGraphService` to visualize import/export dependencies (Integrated `Cognee` for Graph Memory).

### Secure Execution
- [x] **Sandboxing:** Finalize `SandboxService` for secure Docker/WASM code execution (Integrating `code-sandbox/` tools).
- [x] **Auto-Test Runner:** Automatically detect and run relevant tests for modified files.

## Phase 15: The SuperAI Harness (Feature Parity)
**Goal:** Match and exceed the capabilities of Amp, Auggie, Claude Code, Codebuff, and OpenCode.

### Core Engine
- [x] **Unified CLI Runner:** Create `borg cli` wrapper to transparently invoke any integrated CLI (`superai-cli/`).
- [x] **Unified TUI/WebUI:** Ensure 100% feature parity between terminal and web interfaces.
- [x] **Mobile Remote Control:** Responsive mobile UI for monitoring and intervening in agent sessions.
- [x] **Shell Integration:** Deep shell history integration and context awareness (Warp-style).

### Advanced Coding Features
- [x] **Multi-File Context:** "Add to Context" logic similar to Aider/Claude Code.
- [x] **In-Chat Commands:** Slash commands for git operations, diff reviews, and undo steps.
- [x] **Symbol Pinning:** UI to manually prioritize specific code symbols.
- [x] **Auto-Dev Loops:** "Fix until Pass" modes for tests and linters.

## Phase 16: The SuperAI Browser Extension (COMPLETED)
**Goal:** Bridge local context into web-based AI models (ChatGPT, Gemini, Claude.ai).

### Functionality Injection
- [x] **MCP Injection:** Expose local MCP tools (FS, Git, Terminal) to web chats via browser extension.
- [x] **Context Export:** One-click export of web chats into Borg long-term memory (`director.memorize`).
- [x] **Memory Recording:** Background recording of browsing research.

### Browser Capabilities (via MCP)
- [x] **Page Scraping:** Turn current page content into markdown context (Integrated `browser-use` bridge).
- [x] **Console Reader:** Stream browser console logs to the Borg debugger.
- [x] **History:** Secure access to history via `browser_get_history` (WebSocket Bridge).
- [x] **Tab Mirroring:** Real-time visual oversight of the agent's browser tab in the Dashboard.
- [x] **Active Debugger Proxy:** Port forwarding and remote debugging bridge.

## Phase 17: Universal MCP & Orchestration
- [x] **Traffic Inspector:** Real-time JSON-RPC inspection (Completed in Core).
- [x] **Dynamic Disclosure:** Hide tools until needed to save context.
- [x] **Semantic Reranking:** Optimize tool descriptions for model consumption.
- [x] **Proxy System:** Bridge remote/local servers.
- [x] **Workflow Engine:** Graph-based state machine for durable multi-agent execution (LangGraph-inspired).

## Phase 18: Multi-Agent Squads
- [x] **Consensus Protocol:** Multi-model debate engine.
- [x] **Git Worktree Squads:** Parallel coding agents in isolated branches (Maestro Pattern via `WorktreeManager`).
- [ ] **Local-Remote Bridge:** Sync projects between local machine and cloud instances.

## Phase 19: Enterprise & Ecosystem (Restored)
### Advanced Governance & Security
- [x] **RBAC Foundation:** Role-based access control service and middleware (admin, developer, operator, viewer)
- [ ] **SSO integration:** OIDC/SAML foundation
- [x] **Enhanced Audit Logging:** Comprehensive audit trails for agent, council, and architect actions
- [ ] **Policy Engine:** Declarative policies for agent behavior and resource usage

### Distributed Orchestration
- [ ] **Multi-Node Council:** Distributed supervisor councils across different regions/nodes
- [ ] **Agent-to-Agent (A2A) Mesh:** Secure communication protocol for autonomous agent cooperation
- [ ] **Edge Deployment:** Lightweight Borg runtime for edge devices

### Developer Ecosystem
- [ ] **Borg Marketplace:** Registry for sharing supervisors, skills, and agent templates
- [ ] **Visual Designer:** Low-code interface for designing agent workflows and council structures
- [ ] **OpenAPI / SDK:** Formalized external API and multi-language SDKs (Python, Go, Rust)

### Performance & Scalability
- [ ] **GPU Acceleration:** Native support for local model acceleration (llama.cpp integration)
- [ ] **Tiered Memory:** Advanced memory management with cold storage and semantic caching
- [ ] **Batch Processing:** High-throughput task execution for large-scale operations
