---
name: Roundtable AI Delegation
description: Zero-configuration orchestration for delegating tasks to multiple AI assistants (Gemini, Claude, Codex, Cursor) in parallel.
---

# Roundtable AI Delegation

Roundtable is a local Model Context Protocol (MCP) server that coordinates specialized AI sub-agents to solve complex engineering problems. It allows a primary AI (like Claude) to delegate specific sub-tasks to other agents like Gemini, Codex, or Cursor.

## Core Concepts

- **Parallel Execution:** Delegate to multiple models simultaneously.
- **Model Specialization:** Assign the right AI for each task (e.g., Gemini for high-context analysis, Claude for reasoning).
- **Context Continuity:** Shared project context across all sub-agents.

## Quick Start

### Installation
```bash
pip install roundtable-ai
```

### Usage (Claude Code)
```bash
claude mcp add roundtable-ai -- roundtable-ai --agents gemini,claude,codex,cursor
```

## Available Tools

### Availability Checks
Ensure your local CLI tools are working before starting a task.
- `check_codex_availability`
- `check_claude_availability`
- `check_cursor_availability`
- `check_gemini_availability`

### Task Execution
Delegate specific work to a sub-agent.
- `execute_codex_task`
- `execute_claude_task`
- `execute_cursor_task`
- `execute_gemini_task`

## Example Multi-Agent Prompt
"Use Gemini SubAgent to analyze frontend performance issues. Use Codex SubAgent to examine the backend API for N+1 queries. Use Claude SubAgent to review infrastructure logs."
