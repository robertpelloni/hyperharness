# HyperHarness Submodule Dashboard

## Project Directory Structure

```
hyperharness/
├── agent/                    # Agent loop: tool execution, streaming, compaction
├── agents/                   # Multi-agent: Director, Council, RAG, Autonomy
│   └── *_test.go             # Assimilation tests for agent coordination
├── borg/                     # Legacy Borg adapter compatibility
├── cmd/                      # CLI commands (serve, foundation, HTTP)
│   └── hyperharness/         # Main binary entry point
├── config/                   # Global configuration
├── foundation/
│   ├── adapters/             # HyperCode/MCP adapters
│   ├── assimilation/         # Submodule assimilation framework
│   ├── compat/               # Backward compatibility catalog
│   ├── orchestration/        # Planner and plan executor
│   ├── pi/                   # Core tool runtime (the foundation)
│   └── repomap/              # Repository map generator
├── git/                      # Git worktree management
├── hypercode/                # Unified Harness integration layer
├── internal/
│   ├── agent/                # Deep agent runtime (26 functions)
│   ├── borg/                 # Borg core engine with adapters/events
│   ├── config/               # Settings and configuration management
│   ├── context/              # Conversation context management
│   ├── extensions/           # MCP extension manager
│   ├── fs/                   # File system utilities
│   ├── mcp/                  # MCP protocol implementation
│   ├── memory/               # Knowledge base system
│   ├── providers/            # LLM provider abstraction (15+ providers)
│   ├── sessions/             # Session tree management
│   ├── skills/               # Skill discovery and execution
│   ├── subagents/            # Subagent lifecycle (10 types)
│   ├── tools/                # Internal tool registry
│   └── ui/                   # TUI components (23 functions)
├── llm/                      # LLM abstraction layer
├── mcp/                      # MCP client/server (RemoteMCP, deferred loading)
├── orchestrator/             # Job queue, database, RAG store, vectors
│   └── services/             # Orchestrator services
├── repl/                     # REPL interface
├── rpc/                      # JSON-RPC 2.0 (TCP/Unix)
├── security/                 # Permission manager with autonomy levels
├── tools/                    # Tool registry (145+ tools, 17 registration functions)
└── tui/                      # Terminal UI (slash commands, chat, bridge)
```

## Submodule Reference

These submodules from `superai/` were analyzed for tool surfaces:

### Fully Ported (Tool Surfaces Implemented)

| Submodule | Version | Location | Tools Ported | Notes |
|---|---|---|---|---|
| **pi-cli** | v0.0.2-3195 | `superai/pi-cli` | 7 | Foundation tools: read, write, edit, bash, grep, find, ls |
| **claude-code** | main | `superai/claude-code` | 28 | PascalCase: Read, Write, Edit, Bash, Glob, Grep, Agent, TodoWrite, LSP, WebSearch, etc. |
| **crush** | v0.55.0-2 | `superai/crush` | 18 | multiedit, view, glob, webfetch, websearch, diagnostics, todos, tree, batch, delegate |
| **gemini-cli** | MK_TAG_TEST-4562 | `superai/gemini-cli` | 7 | snake_case: read_file, write_file, edit_file, shell, grep, list_directory |
| **opencode** | latest-4261 | `superai/opencode` | 14 | file_read, file_write, file_edit, bash, grep, list_files, multiedit, plan mode |
| **goose** | ls-1360 | `superai/goose` | 4 | developer, text_editor, browser, computer |
| **kimi-cli** | 1.30.0-3 | `superai/kimi-cli` | 14 | read_file, write_file, edit_file, shell, search_files, task management |
| **cursor** | (via code-cli) | `superai/code-cli` | 5 | cursor_read_file, cursor_edit_file, cursor_code_search, cursor_run_command |
| **copilot-cli** | v1.0.18-1 | `superai/copilot-cli` | 1 | copilot_edit (run_in_terminal) |
| **grok-cli** | grok-dev@1.1.4-rc5 | `superai/grok-cli` | 6 | read_file, write_file, edit_file, execute_command, search, list_directory |
| **aider** | v0.86.3.dev | `superai/aider` | 1 | aider_search_replace |
| **mistral-vibe** | v2.7.3 | `superai/mistral-vibe` | 2 | mistral_search, mistral_edit |
| **windsurf** | (via code-cli) | `superai/code-cli` | 2 | read_file, write_file |
| **smithery-cli** | v3.1.6-171 | `superai/smithery-cli` | 2 | smithery_list, smithery_install |

### No Unique Tool Surfaces (Not Ported)

| Submodule | Version | Location | Reason |
|---|---|---|---|
| **adrenaline** | master | `superai/adrenaline` | No CLI tools (concept/documentation only) |
| **auggie** | v0.22.0-pre | `superai/auggie` | Different architecture (Electron app) |
| **azure-ai-cli** | 1.0.0-preview | `superai/azure-ai-cli` | Azure-specific, no unique tools |
| **bito-cli** | v5.3.0-2 | `superai/bito-cli` | No unique tool surfaces |
| **byterover-cli** | 2.6.0-32 | `superai/byterover-cli` | No unique tool surfaces |
| **claude-code-templates** | v1.26.0-488 | `superai/claude-code-templates` | Templates only, not tools |
| **factory-cli** | main | `superai/factory-cli` | No unique tool surfaces |
| **kilocode** | v7.1.20-31 | `superai/kilocode` | VS Code extension, no CLI tools |
| **open-interpreter** | v0.4.2-47 | `superai/open-interpreter` | Python-based, tools covered by other harnesses |
| **qwen-code-cli** | clean-main | `superai/qwen-code-cli` | No unique tool surfaces |
| **rowboat** | v0.1.85 | `superai/rowboat` | No unique tool surfaces |

### Infrastructure (Not Tool Surfaces)

| Submodule | Version | Location | Purpose |
|---|---|---|---|
| **dolt** | v0.22.8-22180 | `superai/dolt` | SQL database with Git semantics |
| **ollama** | v0.13.4-rc2-366 | `superai/ollama` | Local model server |
| **litellm** | v1.83.1-nightly | `superai/litellm` | LLM proxy/gateway |
| **llamafile** | 0.10.0-5 | `superai/llamafile` | Single-file LLM runtime |
| **llm-cli** | 0.30-1 | `superai/llm-cli` | LLM CLI wrapper |
| **jules-extension** | v0.1.0-3 | `superai/jules-extension` | IDE extension |
