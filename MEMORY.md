# HyperHarness Memory - Session Observations

## Codebase Design Patterns

### Tool Registration Pattern
All tools use `(r *Registry) registerXxxTools()` methods with `r.Tools = append(r.Tools, Tool{...})`. 
The dispatcher `registerAllParityTools(r)` in `tools/pi_exact_parity.go` calls all registration functions.
New tool files must: 1) define a registration function, 2) add it to the dispatcher.

### Type Naming Conflicts
Multiple parity files define similar helper types. Solution: use harness-specific prefixes 
(e.g., `ClaudeTodoItem` instead of `TodoItem` which exists in crush_parity.go).

### Platform Considerations
- Windows: `filepath.IsAbs("/path")` returns false. Use `filepath.ToSlash()` for comparison in tests.
- DiffPrettyText: Produces ANSI escape codes. Strip with regex `\x1b\[[0-9;]*m` for test comparison.
- Go binary path: On this system, `PATH` must include `/c/Program Files/Go/bin` explicitly.

### Test Infrastructure
- 315 tests across 25 packages
- Tests run with `go test -buildvcs=false ./... -count=1 -timeout 180s`
- The `-buildvcs=false` flag is needed due to git submodule status issues

### Go Module Dependencies
- `github.com/sashabaranov/go-openai` — OpenAI client
- `github.com/google/uuid` — UUID generation  
- `github.com/spf13/cobra` — CLI framework
- `github.com/sergi/go-diff/diffmatchpatch` — Diff algorithm
- `github.com/mattn/go-sqlite3` — SQLite driver (for memory/knowledge)

## Submodule Analysis

### Has Unique Tool Surfaces (Ported)
crush, gemini-cli, opencode, goose, kimi-cli, cursor, windsurf, copilot-cli, mistral-vibe, grok-cli, claude-code, smithery-cli, aider, pi-cli

### No Unique Tool Surfaces (Skip)
adrenaline, auggie, azure-ai-cli, bito-cli, byterover-cli, factory-cli, kilocode, qwen-code-cli, open-interpreter, rowboat, claude-code-templates

### Infrastructure (Not Tools)
dolt (database), ollama (model server), litellm (proxy), llamafile (runtime), llm-cli (CLI wrapper)

## Build Observations
- Total: ~28K lines Go production code, ~8K lines tests
- Build time: ~10s for full project
- Test time: ~80s (dominated by tools package at 60s)
