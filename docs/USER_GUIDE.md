# HyperHarness User Guide

HyperHarness is a powerful, Go-native control plane for AI-driven development. It allows you to orchestrate multiple LLMs, MCP servers, and specialized agents from a single terminal interface.

## Core Commands

### 1. Interactive TUI
Launch the advanced BubbleTea-based orchestrator:
```bash
./hypercode tui
```

### 2. Autonomous Autopilot
Give the agent a goal and let it plan and execute sub-tasks using the Council architecture:
```bash
# Inside TUI
/autopilot "Refactor the authentication logic to use JWT"
```

### 3. Data Ingestion
Index local files or directories into the persistent knowledge base for RAG:
```bash
./hypercode ingest ./path/to/project
```

### 4. Background Sync
Sync bookmarks and crawl links to ground your agents in external knowledge:
```bash
# Native tools available to agents
sync_bobby_bookmarks
crawl_links
```

## Configuration
HyperHarness reads model provider API keys from environment variables:
- `ANTHROPIC_API_KEY`
- `OPENAI_API_KEY`
- `GOOGLE_API_KEY`
- `OPENROUTER_API_KEY`

## Harness Parity
HyperHarness supports 100% exact tool signatures for:
- Claude Code
- Codex
- Warp / Wave
- Tabby
- Aider
- pi-coding-agent
