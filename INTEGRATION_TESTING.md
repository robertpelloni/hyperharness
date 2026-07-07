
## LLM StreamChat & Subagent Loop
The internal subagent manager utilizes the `llm.AutoRouteStream` to dispatch tasks (such as `Agent` spawned by Claude Code Parity or `task` spawned by OpenCode Parity) to live LLM engines.

The e2e lifecycle behaves robustly under load:
1. `tools.Tool` requests a task.
2. `subagents.GlobalManager.Spawn` is called with a `streamCallback` hook directly bound to stdout or the RPC socket.
3. The LLM processes the system prompt mapped via `llm.GetSubagentPrompt`.
4. If an API is offline or unconfigured, the task catches the network exception and successfully cascades into an offline fallback mock (maintaining session continuity during standard offline unit tests).
5. Context is completed correctly, and `.Done` channel cascades to resolve.

This flow completes Phase 4 integration and transitions the Agent parity tools away from mock stubs into real, verifiable task pipelines.

### Validating the Subagent LLM Lifecycle
The end-to-end integration lifecycle test located in `internal/subagents/manager_integration_test.go` exercises the `Spawn` method's connection to the global LLM auto-router. It verifies:
- Dynamic schema injection based on specialized subagent roles (e.g. `TypeResearch`, `TypeCode`).
- Synchronous block waits paired with asynchronous RPC streaming hooks (`callback`).
- Graceful context cancellation and task timeouts when the LLM stream hangs.
- Offline degradation strategies via mock stub mapping ensuring tests succeed when disconnected or run under tight CI limits.
