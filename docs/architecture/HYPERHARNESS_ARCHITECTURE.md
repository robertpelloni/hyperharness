# Hyperharness Architecture

## Vision

Hyperharness is the unified Go-based AI coding agent harness that assimilates EVERY feature from EVERY AI coding tool on earth, implemented with 100% parameter parity and superior architecture.

## Core Design Principles

1. **Go-First**: Built in Go for maximum performance, minimal dependencies, cross-platform binaries
2. **100% Tool Parity**: Every tool has exact name, parameter signature, and output format as models expect
3. **Modular Architecture**: Clean seams between config, provider, agent, session, tools, MCP, memory, and UI
4. **Extensible**: Extension system for custom tools, skills, prompt templates, themes, and packages
5. **MCP Native**: Built-in MCP client + server, tool aggregation, dynamic tool registration
6. **Memory-Integrated**: Knowledge base with full-text search, semantic indexing, scoped memories
7. **Session-First**: JSONL session files with tree structure, branching, forking, compaction
8. **Multi-Provider**: Unified provider interface with routing, failover, quota awareness
9. **Terminal-Optimized**: Rich TUI with themes, keybindings, autocomplete, streaming output
10. **Compatible**: Drop-in replacement for Pi, compatible with existing skills, extensions, AGENTS.md

## Architectural Layers

```
┌─────────────────────────────────────────────────────────┐
│                     CLI / TUI Layer                     │
│  (bubbletea/repl, commands, model selector, tree view)  │
├─────────────────────────────────────────────────────────┤
│                   Agent Runtime Layer                    │
│  (agent loop, context building, compaction, tool exec)  │
├────────────┬──────────────┬──────────────┬──────────────┤
│  Provider   │    Tool       │   Memory     │     MCP       │
│  Interface  │   Registry    │ Knowledge    │  Registry     │
│            │               │   Base       │              │
├────────────┴──────────────┴──────────────┴──────────────┤
│                    Config / Settings                     │
├─────────────────────────────────────────────────────────┤
│                    Session / Journal                     │
└─────────────────────────────────────────────────────────┘
```

## Package Structure

```
hyperharness/
├── cmd/hyperharness/          # Main CLI entry point
├── internal/
│   ├── config/                # Settings, config management
│   ├── providers/             # LLM provider registry (Anthropic, OpenAI, Google, etc.)
│   ├── sessions/              # Session management (JSONL tree, branching, forking)
│   ├── agent/                 # Agent runtime (loop, context, compaction)
│   ├── tools/                 # Built-in tools (read, write, edit, bash, grep, find, ls)
│   ├── mcp/                   # MCP client/server layer
│   ├── memory/                # Knowledge base (search, store, scope)
│   ├── ui/                    # System prompt builder, theming, formatting
│   ├── tui/                   # Terminal UI (bubbletea) - FUTURE
│   ├── extensions/            # Extension system - FUTURE
│   ├── skills/                # Skill loading - FUTURE
│   ├── compaction/            # Context compaction strategies - FUTURE
│   ├── hypercode/                  # Tool inventory management for HyperCode
│   └── rpc/                   # JSON-RPC protocol (stdin/stdout) - FUTURE
└── docs/                      # Architecture and design documentation
```

## Feature Assimilation Map

### From Pi (Core Foundation)
- [x] Configuration system (global + project, env overrides)
- [x] Provider registry with model catalog
- [x] Session management (JSONL tree, branching, forking, compaction)
- [x] Agent runtime with tool execution loop
- [x] Built-in tools (read, write, edit, bash, grep, find, ls) - exact parameter parity
- [x] System prompt builder with AGENTS.md, skills, memory context
- [x] Message queue (steering/follow-up)
- [x] Token/cost tracking
- [x] Thinking level support (minimal through xhigh)
- [x] Retry with exponential backoff
- [x] Tool output truncation (line-based, byte-based, with temp file fallback)
- [x] Package management (extensions, skills, prompts, themes)
- [ ] TUI with bubbletea/charmbracelet
- [ ] Extension system (TypeScript equivalent in Go)
- [ ] Session tree navigation UI (/tree command)
- [ ] Model cycling (Ctrl+P)
- [ ] HTML export/session sharing

### From Claude Code
- [x] SEARCH/REPLACE patch application
- [x] Compaction with LLM summarization
- [ ] Todo/task tracking system
- [ ] Permission model with configurable gates
- [ ] Sub-agent delegation

### From Aider
- [ ] Multi-file editing with git awareness
- [ ] Auto-commit after changes
- [ ] Read-only file mode
- [ ] Voice input support

### From OpenAI Codex CLI
- [ ] Code execution sandbox
- [ ] Safety checks and guardrails
- [ ] File change preview before applying

### From Goose
- [ ] Web search integration
- [ ] Multi-step planning
- [ ] Autonomous workflow execution

### From Open Interpreter
- [ ] Code execution in multiple languages
- [ ] Browser automation
- [ ] System control

### From Crush (Charm)
- [ ] Permission gates on every tool execution
- [ ] TUI with charmbracelet primitives
- [ ] Shell selection (bash, zsh, powershell)

### From Factory AI
- [ ] Sub-agent orchestration with isolated contexts
- [ ] Concurrent agent execution
- [ ] Result aggregation

### From HyperCode
- [ ] MCP server aggregation and management
- [ ] Provider routing with failover
- [ ] Memory retrieval across sessions
- [ ] Tool inventory exposure to models
- [ ] Mesh capability discovery

### From Gemini CLI
- [ ] Multi-turn conversation with grounding
- [ ] File-level code navigation
- [ ] Google Workspace integration

### From Ollama
- [ ] Local model support
- [ ] Model pulling and management
- [ ] Offloading to GPU

### All Other Tools
- [ ] Web search (from Goose/Open Interpreter)
- [ ] Browser automation (from Agent-Browser/Firecrawl)
- [ ] Image analysis (from Pi image support)
- [ ] MCP tool calling (from Claude Code/Copilot)
- [ ] Git integration (from Aider/Codex)
- [ ] Database connectivity (from various CLIs)
- [ ] SSH/remote execution (from various CLIs)

## Go Architecture Decisions

1. **No dynamic code loading**: Unlike Pi's TypeScript extensions, Go extensions will be compiled as plugins or use RPC
2. **Concurrency-first**: Multiple tools can execute simultaneously using goroutines
3. **Type-safe tools**: All tool parameters validated at compile time via generated code
4. **Zero-copy where possible**: Memory-mapped files for large session reading
5. **SQLite for memory/KB**: Embedded sqlite3 for persistent knowledge storage
6. **Standard library first**: Minimize external dependencies
7. **Cross-platform**: Windows, macOS, Linux support
8. **Modular tool loading**: Tools registered via constructor functions, enabling conditional compilation

## Integration with HyperCode

Hyperharness serves as the harness layer in HyperCode's architecture:
- Tool provider for HyperCode's MCP layer
- Session manager for HyperCode's control plane
- Memory substrate for HyperCode's context system
- Provider routing backend for HyperCode's AI layer
- Submodule of HyperCode at `submodules/hyperharness`

## Implementation Priority

### Phase 1: Foundation (Current)
- [x] Go module and package structure
- [x] Configuration system
- [x] Provider registry
- [x] Session management
- [x] Agent runtime
- [x] Core tools with parameter parity
- [x] MCP registry and client
- [x] Memory knowledge base
- [x] System prompt builder

### Phase 2: TUI & UX
- [ ] Bubbletea TUI shell
- [ ] Model selector
- [ ] Session tree navigation
- [ ] File autocomplete
- [ ] Theme system
- [ ] Keybindings
- [ ] Image handling

### Phase 3: Advanced Features
- [ ] Extension system
- [ ] Package management
- [ ] Sub-agent delegation
- [ ] Compaction implementation (full LLM summarization)
- [ ] Web search
- [ ] Browser automation

### Phase 4: HyperCode Integration
- [ ] Full tool inventory API
- [ ] Provider routing integration
- [ ] Memory sync
- [ ] Session import/export

### Phase 5: Parity Complete
- [ ] Every tool from every agent
- [ ] 100% functional parity
- [ ] Comprehensive test suite
- [ ] Cross-platform binaries
- [ ] Documentation
