import sys
import re

with open("HANDOFF.md", "r") as f:
    handoff = f.read()

# Update HANDOFF.md next steps
handoff_old = """## Next Steps
1. Wire tool detectors into actual tool dispatch (use detected tools as backends)
2. MCP deep integration: actual stdio transport, SSE transport
3. Memory SQLite backend with vector embeddings
4. Integration tests with actual AI model tool calling
5. Port remaining submodule agent loop patterns (streaming, context budgets)
6. Performance benchmarks
7. Wire all tools to real backends: TodoWrite→SessionTodoStore, Agent→subagents, LSP→lsp_client"""

handoff_new = """## Next Steps
1. Wire tool detectors into actual tool dispatch (use detected tools as backends)
2. Memory SQLite backend with vector embeddings
3. Integration tests with actual AI model tool calling
4. Port remaining submodule agent loop patterns (streaming, context budgets)
5. Performance benchmarks
6. MCP deep integration: Bidirectional routing (expose internal tools via MCP server)"""

handoff = handoff.replace(handoff_old, handoff_new)

# Update state in HANDOFF.md
handoff_state_old = "## Current State: v0.3.0"
handoff_state_new = "## Current State: v0.3.0\n- **Deep Tool Wiring Complete**: TodoWrite, Agent, LSP, WebSearch, WebFetch, Config, Skill are now wired to real backends.\n- **MCP Transports Complete**: Both stdio and SSE client transports are fully implemented."

handoff = handoff.replace(handoff_state_old, handoff_state_new)

with open("HANDOFF.md", "w") as f:
    f.write(handoff)

with open("MEMORY.md", "r") as f:
    memory = f.read()

memory_new = """[PROJECT_MEMORY]

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
- **Memory System (`internal/memory/`)**: Implements a Knowledge Base for long-term storage and retrieval. It currently aims to utilize SQLite FTS5 (`glebarez/sqlite`) to prevent conflicts with other CGO dependencies.
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

## Current State & Recent Wiring
- **Deep Tool Wiring**: Recently, dummy stub implementations in the parity layer were successfully wired to their actual backend Go implementations.
  - `TodoWrite` now updates a persistent `SessionTodoStore`.
  - `Agent` connects to the `subagents.Manager` to spawn tasks.
  - Web tools (`WebSearch`, `WebFetch`) map to real HTTP/Exa clients.
  - `LSP` maps to actual Language Server Protocol clients.
  - `Config` and `Skill` map to their respective internal managers.
- **MCP Transport**: Both `stdio` and `sse` transports have been fully implemented, handling JSON-RPC handshakes (`initialize`, `tools/list`, `tools/call`).

## Future Roadmap Focus (TODO.md)
The next primary objectives involve:
1. **Memory Enhancements**: Wiring up SQLite FTS5 for knowledge bases and vector embeddings for semantic search.
2. **Deep MCP Integration**: Finalizing the bidirectional routing of internal tools exposed over the MCP server interface.
3. **Comprehensive Testing**: Adding robust integration tests spanning the entire agent loop, tool execution pipelines, and MCP lifecycle.
"""

with open("MEMORY.md", "w") as f:
    f.write(memory_new)
print("Updated docs")
