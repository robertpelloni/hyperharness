# HyperHarness Vision

## The Universal Go-Native AI Coding Agent Harness

HyperHarness is the **definitive** Go-based AI coding agent that assimilates features from every major CLI coding tool on earth with 100% parameter parity and superior Go-native architecture.

### Core Philosophy

**Assimilation, not Integration.** We don't wrap external tools — we absorb their best features into a single, cohesive, Go-native system. Every tool surface from every harness is implemented natively with the exact same API contract, so AI models trained on any specific harness work identically with HyperHarness.

### Why This Matters

Major AI models (Claude, GPT-4, Gemini, etc.) are trained on specific tool interfaces. When they see `Read`, `Write`, `Edit`, `Bash`, `Glob`, `Grep` — they know exactly how to use them. By providing these exact tool surfaces natively in Go, we get:

1. **Better Model Performance** — Models produce better code when tools match their training
2. **Single Binary Deployment** — No Node.js, Python, or external runtime dependencies
3. **Unified Architecture** — One codebase, one design philosophy, one API surface
4. **Cross-Harness Compatibility** — Works with prompts designed for any supported harness

### Architecture Principles

1. **Tool Delegation Pattern** — All harness-specific tools delegate to foundation `pi` tools (read, write, edit, bash, grep, find, ls) for actual execution. Single source of truth.

2. **Registry-First Design** — Every tool is a `Tool` struct with `Name`, `Description`, `Parameters` (JSON Schema), and `Execute` function. Clean discovery and invocation.

3. **Borg Adapter Pattern** — All subsystems (memory, sessions, MCP, extensions, tools) register as Borg adapters with unified lifecycle.

4. **Case-Sensitive Naming** — Tools maintain exact casing from source harness (Claude Code PascalCase, Pi lowercase, Gemini snake_case).

5. **Singleton Infrastructure** — Thread-safe lazy initialization with `sync.Once` for knowledge base, job manager, and core systems.

### Supported Harness Sources

| Harness | Tools | Key Features Absorbed |
|---|---|---|
| Pi | 7 | Foundation tools (read, write, edit, bash, grep, find, ls) |
| Claude Code | 28 | Agent spawning, TodoWrite, LSP, plan mode, web search, skills, worktrees |
| Gemini CLI | 7 | File operations with schema validation |
| OpenCode | 14 | Advanced file editing, multi-edit, plan mode |
| Kimi CLI | 14 | Task management, search, file operations |
| Crush | 18 | Background jobs, batch execution, diagnostics, web tools |
| Goose | 4 | Developer tool, browser automation, computer control |
| Cursor | 5 | Code search, edit, run |
| Windsurf | 2 | File operations |
| Copilot CLI | 1 | Terminal execution |
| Aider | 1 | Search-replace editing |
| Mistral | 2 | Shell and editor |
| Grok | 6 | File ops, search, listing |
| Smithery | 2 | MCP server discovery |
| **TOTAL** | **145+** | |

### Integration with Borg/HyperCode

HyperHarness serves as:
- **Built-in Tool Provider** — All 145+ tools available to any agent
- **MCP Client/Server Aggregator** — Route MCP calls through unified registry
- **Memory System** — Persistent knowledge base with FTS, tags, scopes
- **Context Manager** — Automatic compaction, injection, and budget management
- **Subagent Orchestrator** — 10 specialized agent types for delegation

### Target Users

1. **AI Agent Developers** — Need a reliable, fast tool provider for their agents
2. **Enterprise Teams** — Need single-binary deployment with no external dependencies
3. **Researchers** — Need consistent tool surfaces for model evaluation
4. **Open Source Contributors** — Want to extend the harness with new tools

### Long-Term Goals

1. **Zero External Dependencies for Core** — Everything native Go
2. **Full MCP Protocol Support** — Act as both MCP client and server
3. **Multi-Agent Orchestration** — Council debates, director mode, autonomous loops
4. **Plugin System** — Go plugins for custom tool providers
5. **Performance Benchmarks** — Sub-millisecond tool dispatch
6. **Cloud Deployment** — gRPC/HTTP API for remote tool execution
