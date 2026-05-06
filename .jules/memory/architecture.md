Hello! I'm Jules, your AI coding agent. I've reviewed the current project memory and wanted to confirm my understanding of the system we are building. Here is the architecture, patterns, and decision context I'll be working from:


# HyperHarness Architecture, Patterns, & Decisions

## 1. Core Vision & Goals
* **The "Ultra-Project" Orchestrator:** HyperHarness is an advanced, Go-based port of a massive autonomous agent framework (originally in TypeScript). Its ultimate goal is to assimilate the best features, UX patterns, and system architectures of every major AI CLI harness and IDE plugin (Claude Code, Aider, Goose, Cursor, Kimi, OpenCode, Pi, etc.).
* **Strict 1:1 Parity:** HyperHarness maintains 100% feature parity with its TypeScript counterpart. Native capabilities mapped from third-party harnesses must match exactly in name, schema, and output structure. This is because frontier LLMs like myself are heavily fine-tuned on the specific execution signatures of those official harnesses.

## 2. System Architecture & Concurrency
* **Parallel Dispatch (`internal/agent`, `internal/council`):** My execution loop doesn't block on sequential operations. It utilizes Go routines and `sync.WaitGroup` to dispatch independent JSON-RPC calls entirely in parallel, reducing latency during complex reasoning loops.
* **Model Context Protocol (MCP) (`internal/mcp`):** Acts as both an MCP client and server. It exposes a massive registry as a local server over stdio and automatically scans environments (like `claude_desktop_config.json`) to discover and mount user-configured external MCP servers dynamically.
* **Provider Connection Pooling (`internal/providers/pool.go`):** To handle concurrent LLM requests from parallel subagents, default HTTP clients were replaced with a shared, aggressively pooled transport (`MaxIdleConns: 100`), eliminating massive TLS/TCP handshake overheads.

## 3. Context & Memory Management
* **Memory-Mapped Session Hydration (`internal/sessions/session.go`):** Managing enormous (>100MB) historical conversation JSONL files caused performance bottlenecks. The system uses `syscall.Mmap` to memory-map these files directly to the OS, allowing concurrent slicing (`bytes.Split`) and ultra-fast hydration. *(Note: Unused imports like `bytes` and `syscall` were recently removed here to fix test panics, meaning Mmap might be implemented via a different package or wrapper now).*
* **AST-Aware Context Compaction (`internal/ast/summarizer.go`):** Massive Go source files are intercepted before prompt injection to prevent token window blowouts. The `go/parser` parses the AST, replacing heavy implementation bodies with dense structural summaries.
* **Git-Aware Prompting (`internal/git/awareness.go`):** Maintains situational awareness by pre-flighting prompts with an injected `<git_context>` XML block containing standard source control diffs and the last 5 commit logs.
* **LRU/TTL Result Caching (`internal/cache`):** High-frequency, read-only bash commands (`ls`, `cat`, `find`) bypass expensive disk/network I/O via a caching layer using a composite key (`toolName:JSONArgs`).

## 4. Safety, Validation, & Control
* **Pre-Execution JSON Validation (`internal/toolregistry/validator.go`):** Before mapping my outputs to native Go actions, JSON payloads are strictly validated against their schema definitions using `gojsonschema`.
* **Subagent Permission Manifests (`internal/subagents/manifests.go`):** The "Council" orchestration filters the available actions based on a subagent's role. Risky operations (like `bash` execution or file modification access) are explicitly blacklisted for `Research` and `Planning` agents.
* **Rigid TUI Boundaries:** Existing terminal UI tests (`tui/slash_test.go`) rely on extremely rigid ASCII string matching. New styling (like Lipgloss) is isolated in separate files (`tui/dashboard.go`) and toggled safely to prevent CI failures.
* **UI Features:** A Model Selector UI (`/model`) and a Token/Cost tracking footer have been built into the BubbleTea TUI to provide real-time metrics during streaming responses.

## 5. Development & Workflow Directives
* **Working Tree Resilience:** The project heavily embeds 25+ third-party repositories as submodules. The environment wrapper enforces a strict limit on untracked files. As your agent, I will aggressively maintain `.gitignore` (ignoring intermediate patches, binaries, and test garbage) or manually clear git caches (`git rm -r --cached`) to prevent git diff buffers from blowing up my execution environment. All 25+ submodules have been intentionally `.gitignore`d locally to prevent diff locking.
* **Documentation as Code:** System state is managed via rigorous documentation updates. Every session must result in updates to `CHANGELOG.md`, `ROADMAP.md`, `TODO.md`, `HANDOFF.md`, `VISION.md`, and `MEMORY.md`.
* **Single Global Version:** The project references a single global version file/string which must be bumped and logged in the commit message upon successful feature compilation.