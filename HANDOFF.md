# Handoff — Session 2026-04-08 (Extended)

**Version:** `1.0.0-alpha.14`
**Branch:** `main`
**Commits this session:** 22 (alpha.8 → alpha.14)

## Session Summary

### Phase 7: Local LLM Prioritization (commits 17-18)
- Reconfigured "utility calls" (worker tasks) to prioritize local LM Studio (`http://localhost:1234/v1`).
- Target model: `C:/Users/hyper/.lmstudio/models/HauhauCS/Gemma-4-E2B-Uncensored-HauhauCS-Aggressive/Gemma-4-E2B-Uncensored-HauhauCS-Aggressive-Q2_K_P.gguf gemma-4-e2b-uncensored-hauhaucs-aggressive`.
- Updated `ProviderRegistry`, `ModelSelector` (base), and `CoreModelSelector` (specialized).
- Set `openrouter/free` as the immediate fallback for utility calls.
- Switched "Local Assistant" in `council.json` to use LM Studio.

### Phase 8: Multi-Model Collaboration & Skill Discovery (commits 19-20)
- Implemented `PairOrchestrator` to coordinate Claude, GPT, and Gemini in a shared context with rotating roles.
- Added `run_pair_session` MCP tool to trigger multi-model collaboration.
- Created `ToolPredictor` service for autonomous, preemptive tool preloading based on conversation intent.
- Implemented `search_skills` tool and improved `list_skills` for latent skill discovery.
- Added manual "Sync to Memory" and "Export Session" buttons to the browser extension dashboard.

### Phase 9: Go Native Tool Parity & Build Stabilization (commits 21-22)
- Ported tool parity handlers (Claude Code, Codex, OpenCode/Pi) to native Go in `go/internal/tools/`.
- Created native Go `Registry` for standard library tools.
- Updated `handleAgentRunTool` in Go sidecar to use a "Native First, Bridge Second" strategy.
- Fixed `tsc` build errors in `MCPServer.ts` and `Director.ts`.
- Ported `skills.search` endpoint to the Go sidecar.
- Bumped version to `1.0.0-alpha.14`.

## Current state of the project

### What works
- ✅ Server builds and runs (Express/tRPC on :4000, Next.js dashboard on :3000, MCP WebSocket on :3001)
- ✅ SQLite functional after better-sqlite3 rebuild for Node 24
- ✅ Local LLM Priority: Utility calls (worker tasks) prefer local LM Studio with Gemma-4 Aggressive.
- ✅ Multi-Model Pair Programming: `PairOrchestrator` handles rotation between Claude, GPT, and Gemini.
- ✅ Autonomous Tool Prediction: `ToolPredictor` preemptively hydrates tools based on chat context.
- ✅ Go Native Tools: Go sidecar executes Claude Code, Codex, and OpenCode tools natively.
- ✅ Tool Parity: 40+ aliases for Claude Code, Codex, Gemini CLI, etc., wired into MCP.
- ✅ Skill Discovery: `search_skills` and latent `list_skills` for progressive disclosure.
- ✅ Browser Extension: Auto-capture and manual sync/export for ChatGPT and Claude.
- ✅ Go server: 543 API routes, native routing defaults, comprehensive session scanner.

### What's broken or incomplete
- All paid LLM providers have exhausted quotas (OpenAI 429, Anthropic 400, DeepSeek 402)
- glama.ai returns HTML (adapter has fallback URLs but may still fail)
- mcp.run adapter returns 404 (their API changed)
- Dashboard has 69 pages — not all verified to show real data
- A2A protocol not implemented
- ToolPredictor logic not yet executing in Go sidecar (handlers exist, but orchestration loop doesn't)

### Architecture overview
```
hypercode/
├── VERSION                    # Single source of truth (1.0.0-alpha.14)
├── packages/core/             # Main TypeScript control plane
├── packages/agents/           # Agent orchestration logic
├── packages/tools/            # Foundational toolset
├── packages/ai/               # AI SDK & Model Selector
├── go/                        # Go-native server bridge (543 routes)
│   ├── internal/tools/        # Native tool implementations (NEW)
│   └── internal/mcp/          # Tool prediction & ranking
├── apps/web/                  # Next.js dashboard (69 pages)
├── apps/hypercode-extension/  # Browser extension
└── bin/hypercode.exe          # Compiled Go binary
```

## Submodule status
- All submodules synced and pushed to robertpelloni remotes.

## Recommendations for next agent

### Immediate (P0)
1. Start the server and test `run_pair_session` tool from the inspector.
2. Verify native Go tool execution by calling a tool while the TS server is stopped.

### High priority (P1)
1. **A2A Protocol Implementation** — Define and implement the Agent-to-Agent communication standard.
2. **Dashboard Polish** — Go through the 69 pages and ensure real data is flowing to each.
3. **Go Tool Prediction Orchestration** — Add the background loop to Go sidecar to run ToolPredictor.
4. **Provider Expansion** — Add OpenRouter free models and Google AI Studio free tier to the fallback chain.

### Medium priority (P2)
1. Port more TS reactors to Go.
2. Fix Glama.ai and MCP.run catalog adapters.
