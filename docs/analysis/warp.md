# Warp Terminal AI Features Analysis

## Overview
Warp is a modern, Rust-based terminal with agentic AI integration. Its core AI component is "Oz", an agent orchestration platform.

## Agent Features ("Oz")
- **Orchestration**: Manages child agents, execution modes (local/remote), and skills.
- **Actions**:
  - `RequestCommandOutput`: Run a shell command and capture output.
  - `ReadFiles`: Read local files.
  - `SearchCodebase`: Search the project using query and partial paths.
  - `RequestFileEdits`: Apply diffs to files.
  - `Grep`: Search for text across files.
  - `FileGlob`: Find files matching glob patterns.
  - `CallMCPTool`: Invoke Model Context Protocol tools.
  - `ReadMCPResource`: Read MCP resources.
  - `UseComputer`: High-level computer automation actions (keyboard, mouse, screen).
  - `StartAgent` / `RunAgents`: Spawn subagents.
  - `AskUserQuestion`: Interactive multi-choice questions.

## Implementation Details in HyperHarness
- **Parity Tools**: Implement Warp's specific tool signatures in `tools/warp_parity.go`. Many of these (ReadFiles, SearchCodebase, Grep) can delegate to existing `foundation/pi` tools but must maintain Warp's exact parameter names (`path` vs `patterns` vs `queries`).
- **Orchestration**: Warp's `RunAgents` and `StartAgent` logic should map to HyperHarness's `internal/subagents` and `internal/workflow`.
- **Terminal Context**: Warp agents are highly aware of terminal block context. This should be implemented by extending `internal/workspaces/tracker.go` to capture and store command outputs in a structured way (similar to Warp's `BlockId`).
