---
name: Owlex AI Council
description: Multi-agent council for parallel CLI deliberation and specialist role critique.
---

# Owlex AI Council

Owlex is a multi-agent deliberation server designed to provide "second opinions" and structured critique between different LLMs (Codex, Gemini, OpenCode, etc.) directly from the CLI.

## Core Concepts

- **Round 1 (Independent):** Questions are sent to agents independently.
- **Round 2 (Deliberation):** Agents review each other's answers and revise their positions.
- **Synthesis:** A primary agent (like Claude) reviews the entire debate to provide a final recommendation.

## Quick Start

### Installation
```bash
uv tool install git+https://github.com/agentic-mcp-tools/owlex.git
```

### Configuration (.mcp.json)
```json
{
  "mcpServers": {
    "owlex": { "command": "owlex-server" }
  }
}
```

## Common Directives

### Council Ask
Query the council for an architecture decision or a tricky bug.
```bash
council_ask prompt="Should I use a monorepo or multiple repos for 5 microservices?"
```

### Specialist Roles
Assign specific perspectives to agents:
- `security`: Focus on vulnerabilities and risks.
- `perf`: Focus on latency and resource usage.
- `skeptic`: Find flaws in the reasoning.
- `architect`: Focus on system design and patterns.
- `maintainer`: Focus on readability and technical debt.

### Team Presets
- `security_audit`
- `code_review`
- `architecture_review`
- `devil_advocate`

## Individual Sessions
Start a continuous session with a specific model:
- `start_codex_session`
- `start_gemini_session`
- `start_opencode_session`
