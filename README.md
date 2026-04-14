# 🔥 HyperCode 

**The Local-First Control Plane for AI Operations.**  
*(Now rewritten in lightning-fast Go)*

> **Status:** `1.0.0-alpha.32` — **Native Go Convergence**  
> **Focus:** Speed, Unfailing Reliability, Actionable Context, & Total Tool Autonomy.

HyperCode is the ultimate locally-hosted control system orchestrating your entire AI tool stack. Designed for power-users leveraging multiple MCP (Model Context Protocol) servers, diverse model providers, and sprawling context across AI coding sessions, **HyperCode unifies them all into one blindingly fast, operator-owned command center.**

---

## 🚀 Why HyperCode Exists

Modern AI-assisted development is deeply fragmented:
- **Too many disparate MCP servers** with overlapping capabilities.
- **Provider quota exhaustion** interrupting deep workflows.
- **Zero continuity** between coding sessions across IDEs and CLIs.
- **Blind LLMs** that constantly have to "guess and check" what tools exist.

**HyperCode fixes this.** It acts as an autonomous nervous system for your models, providing them with proactive context, high-value tool suggestions, shared memory, and automatic failovers without requiring a gated cloud backend. 

---

## ⚡ What Makes HyperCode Special

### 1. 🧠 The "Decision System" (Smart MCP Router)
HyperCode is not just another dumb MCP proxy. It is an **autonomous decision system**:
- Instead of overwhelming models with 50+ raw tools, HyperCode exposes **6 permanent meta-tools**.
- Features **Ranked Discovery & Auto-Loading**: When a model queries a capability, HyperCode scores matches (via BM25 + semantic groupings) and silently auto-loads high-confidence tools in the background.
- Employs **LRU Tool Eviction** to aggressively manage the model’s context window so it never bloats.
- **Built-in Parity Aliases:** Models trained to use specific internal tools (e.g., *Claude Code's* `edit_file`, *Codex's* `shell`, or *Cursor's* `codebase_search`) automatically get transparently aliased to unified robust implementations natively.

### 2. 📡 Preemptive Tool Advertiser
Why force the model to search for tools? HyperCode continuously parses your conversation transcripts in real-time, extracts intent via bigrams and token distribution, and **proactively injects short tool advertisements** for capabilities it knows will solve the current roadblock. 

### 3. 💾 Unified Tiered MemoryManager
Say goodbye to AI amnesia. HyperCode features a native `MemoryManager` that categorizes data into **Short-term (current session)**, **Medium-term (recent summaries)**, and **Long-term (permanent facts/skills)** storage. It autonomously handles promotion, demotion, and deduplication of context using smart access decay.

### 4. 🕵️ Universal Session Autodiscovery
HyperCode acts as a black box flight recorder. It recursively scans and ingests session logs across all your tools—**Claude Code, Cursor, Windsurf, OpenCode, Gemini CLI, Aider, Copilot, etc.**—parsing arbitrary JSON/MD histories into a continuous unified timeline context.

### 5. 🛡️ Robust Go-Native Sandboxed Execution
We threw out the sluggish Node wrappers. HyperCode now features a **Go-native CodeExecutor** that sandboxes scripts (`bash`, `javascript`, `python`, `rust`, etc.) on the fly with strict timeouts and output truncations, letting you safely run complex `code-mode` scripts to chain MCP tools together autonomously. 

---

## 🛠️ Quick Start

Starting with `1.0.0-alpha.32`, HyperCode's core has been completely extracted out of TypeScript into a pure compiled Go binary, granting massive performance gains and zero-dependency startup times. 

### Requirements
- **Go 1.22+**

### Launching the Daemon
Simply run the included startup script. It will transparently handle building the binary.
```bash
./start.sh
# or on Windows:
start.bat
```

Interact with the daemon via its 40+ extensive native HTTP API endpoints or run the `hypercode` CLI binary directly for terminal operation.

---

## 🏗️ Repository Architecture

With the aggressive transition to native Go, the monolithic legacy Typescript codebase has been successfully archived. 

```text
/bin                 -> Compiled HyperCode binaries
/go
  /cmd/hypercode     -> CLI & Daemon Entrypoint
  /internal
    /mcp             -> The Decision System, Aggregator & Preemptive Advertiser
    /memory          -> Multi-tiered context engine
    /codeexec        -> Sandboxed arbitrary multi-language runner
    /ctxharvester    -> Semantic chunker & context pruner
    /sessionimport   -> Universal session schema normalizer
    /httpapi         -> 40+ Control Plane Endpoints
    /gitservice      -> Native git bridging
 
/archive             -> Legacy Typescript packages, UIs, and historical research
```

---

## 🔥 The Vision

HyperCode is aiming for absolute feature parity with the leading proprietary AI tools (*Cursor, Windsurf, Claude Code*), but keeping the orchestration layer entirely open, local, and operator-owned. 

In the immediate roadmap, HyperCode will introduce:
- **Vector embedding overlays** for deep semantic memory recall.
- **Progressive Skill Disclosure** mirroring the seamless tool decision system.
- **Autonomous Multi-Model Chatrooms** where Claude, GPT, and Gemini natively hand off implementation, testing, and planning roles to each other over a shared context. 

*HyperCode: Keep your data local. Keep your models honest. Don't stop the party.*
