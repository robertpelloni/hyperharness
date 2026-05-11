# HyperHarness IDEAS & Future Enhancements

Based on an extremely deep analysis of the current HyperHarness architecture, code, conversation history, and assimilated submodules, here is a comprehensive list of ideas, pivots, and improvements to elevate the project to an "insanely great" state.

## 1. Advanced UI/TUI & Frontend Pivots
- **Bubbletea-driven Terminal Dashboard:** Currently, the TUI is slated for the medium term. We should build an advanced dashboard using `charmbracelet/bubbletea` and `charmbracelet/lipgloss` that features:
  - Split panes: A left pane for the conversation thread, a right pane for active tool executions (e.g., tailing bash outputs or MCP server logs).
  - A bottom status bar tracking token usage, active subagents, and memory scope.
- **Native GUI Frontends:** Using Wails (Go) or Fyne, we can wrap the Go core into a native desktop app for macOS, Windows, and Linux. This bypasses terminal limitations and allows rich markdown rendering, diff visualization (like VSCode's diff viewer), and interactive architectural graphs.

## 2. Multi-Agent "Council" Architecture
- **Director-Worker Paradigm:** Implement a strict hierarchy where a "Director" agent breaks down a user prompt, delegates to "Researcher", "Architect", and "Coder" subagents, and then reviews their work before presenting it to the user.
- **Debate Mode:** If a code edit fails tests, spawn a temporary "Debate" session between two agents with different system prompts to argue over the best approach, summarizing the consensus back to the main thread.

## 3. Advanced Memory & Context Management
- **Git-Aware Context:** The context manager should automatically pull in the `git diff` of the current working tree and the commit messages of the last 3 commits to ground the LLM in the current state of work.
- **AST-Aware Chunking:** Instead of simple line truncations for files, integrate an AST parser (like tree-sitter for Go, TS, Python) to summarize files by their function signatures and struct definitions when injecting context, saving massive amounts of tokens.

## 4. MCP Ecosystem Expansion
- **Auto-Discovery of MCP Servers:** Scan the local system (e.g., standard VSCode extension directories, global NPM prefix) for known MCP servers and automatically offer to spin them up.
- **Smithery Global Search:** Build a TUI command `/install <query>` that searches the Smithery registry live and hot-loads the MCP server into the current running session without a restart.

## 5. Security & Sandboxing
- **Containerized Bash:** Currently, `Bash` executes directly on the host. Add an optional Docker-backend adapter for the `Bash` tool, executing shell commands in an ephemeral container to prevent catastrophic LLM mistakes (e.g., `rm -rf /`).
- **Granular Tool Permissions:** Create a permission manifest per session or subagent. For instance, a "Researcher" agent is allowed to use `WebSearch` and `read`, but strictly denied `write` and `bash`.

## 6. Language & Runtime Pivots
- **Typescript/Node Parity via WebAssembly:** Since the prompt requests exact parity with the TS version, we can compile parts of the Go core to WebAssembly (Wasm) and run it inside the Node.js context, ensuring 100% identical business logic for file parsing and memory decay across both ecosystems.

## 7. Refactoring & Code Quality
- **Code Generation for Parity Tools:** Maintaining 145+ parity tools manually is prone to drift. We should create a Go `//go:generate` script that reads the JSON schemas of upstream tools (extracted via git submodules) and automatically generates the `Tool{}` structs and parameter parsing boilerplate in Go.
- **Unified Event Bus:** Replace direct callbacks with a robust Pub/Sub event bus (e.g., using Go channels) for `RunEventType`. This allows multiple listeners (TUI, logger, debug server) to cleanly subscribe to agent states without coupling.
