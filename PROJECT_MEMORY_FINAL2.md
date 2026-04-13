[PROJECT_MEMORY]

# HyperHarness Project Architecture & Findings

## Overview
HyperHarness is a highly ambitious project designed to assimilate, port, and integrate the functionalities of multiple AI CLI harnesses (such as Claude Code, Aider, Copilot CLI, Gemini CLI, etc.) into a single, unified, Go-native orchestrator. The primary goal is achieving 100% exact tool parity, matching parameter schemas and naming conventions, since major LLMs have been specifically trained on these exact tool signatures.

## Core Architecture
The system is built with a modular, highly decoupled architecture heavily utilizing the Adapter pattern.
Key architectural components include:

- **Foundation (`foundation/`)**: Provides the base tools derived from Pi (e.g., read, write, edit, bash, grep) which other higher-level or third-party parity tools delegate to.
- **Tool Parity Layer (`tools/`)**: Houses exact replicas of tool interfaces from multiple CLI tools (e.g., `claude_code_parity.go`, `goose_opencode_kimi_parity.go`, `crush_parity.go`). These tools execute by mapping their inputs into the internal subsystems.
- **Borg Core Engine (`internal/borg/`)**: A central orchestrator that provides an adapter pattern to connect various models and agents, managing their lifecycles, hooks, and events.
- **Context Manager (`internal/context/`)**: Handles token tracking, context compaction, history reduction, and prompt injection to efficiently manage LLM context windows.
- **Memory System (`internal/memory/`)**: Implements a Knowledge Base for long-term storage and retrieval. It currently utilizes SQLite FTS5 (`glebarez/sqlite`) to prevent conflicts with other CGO dependencies.
- **Session Manager (`internal/sessions/`)**: Manages conversational state, storing session trees as JSONL files. It supports complex behaviors like branching, forking, and compaction.
- **Subagent Manager (`internal/subagents/`)**: Spawns specialized agents (Code, Research, Test, Plan, etc.) with specific configurations, allowed tools, and token limits.
- **MCP (Model Context Protocol) (`internal/mcp/`)**: A transport layer that acts both as a client to external MCP servers (via `stdio` and `sse`) and can serve internal tools outwardly.
- **Providers (`llm/`, `internal/providers/`)**: Abstracts integration with various LLM backends (Anthropic, OpenAI, Gemini, Ollama, etc.) with auto-routing and failover capabilities.
- **Skills System (`internal/skills/`)**: Discovers and executes predefined instructional patterns stored in `SKILL.md` files.
- **TUI (`tui/`)**: A Bubbletea-based interactive REPL.

## Key Design Decisions & Patterns
- **Strict Parity over "Better Naming"**: Tools from source harnesses are named EXACTLY as they are in the original source (e.g., `TodoWrite`, `platform__manage_schedule`) to leverage the training data of the LLMs.
- **Delegation over Duplication**: Rather than implementing "grep" 5 times for 5 different CLI tools, parity tools validate their specific schema parameters and delegate to a core underlying implementation (like `foundation/pi.Grep`).
- **Concurrent & Safe Execution**: The codebase heavily relies on `sync.RWMutex` to ensure thread-safe operations across global states, caches, registries, and session stores.
- **Interface & Dependency Decoupling**: Components are designed to interact via well-defined structs and interfaces (e.g., the `Tool` struct, `JSONRPCMessage` for MCP), allowing easy mocking and testing.
- **CGO-Free Preference**: The project prefers pure Go implementations (like `glebarez/sqlite` over `modernc.org/sqlite`) to simplify cross-platform compilation and deployment.
- **Unified Helpers over Duplicate Logic**: A centralized `tools/helpers.go` exists to abstract shared parameter type coercion (like converting JSON Numbers gracefully to `int` or extracting `map[string]interface{}` configurations), minimizing scattered utility instances.

## Current State & Recent Improvements
- **Deep Tool Wiring**: Completed connecting stub parity tools into real backend modules:
  - `TodoWrite` interacts with `SessionTodoStore` tracking session-level to-dos.
  - `Agent` connects to `subagents.Manager` to create tasks dynamically.
  - `WebSearch` invokes live web search (like Exa).
  - `WebFetch` parses remote content correctly via HTTP.
  - `LSP` coordinates IDE-like capabilities securely.
  - `Config` dynamically uses runtime configurations (`config.LoadConfig()`).
  - **OpenCode** `plan_enter` / `plan_exit` and **Crush** `delegate` are now securely routed towards the orchestration and subagent state layers mapping specific task requirements.
- **MCP Transport Advancements**:
  - Full client support for `stdio` MCP server spawning via `exec.Command` and JSON-RPC lifecycle connections.
  - Initial `sse` (Server-Sent Events) support over HTTP clients matching the latest protocol definitions.
- **Memory System Progression**:
  - Introduced schema capabilities for SQLite stores parsing vector embeddings (`[]float32`) alongside metadata.
  - Enhanced search indexing with Memory Decay Weighting – computing relevancy via FTS5 `rank` and penalizing entries based on age (`julianday`).

## Future Roadmap Focus (TODO.md)
The next primary objectives involve:
1. **Tool Expansion & Parallelism**: Finalizing parallel job coordination features from older shells like Crush (`batch` execution mappings, deeper background process manipulation for `job_*`).
2. **Testing Infrastructure**: Formulating comprehensive integration testing for execution paths spanning agent loops, `Registry` tool discovery, and verifying correct MCP process lifecycles.
3. **Extensibility Additions**: Maturing Smithery integrations and allowing flexible dynamic tool loading.
