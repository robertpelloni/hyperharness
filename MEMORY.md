# HyperHarness Memory System

## Architecture
- **internal/memory/knowledge.go**: In-memory KnowledgeBase with CRUD, tagging, scoping
- **internal/memory/sqlite_store.go**: SQLite FTS5 full-text search over memories
- **tools/pi_exact_parity.go**: `memory_store`/`memory_search` tools wired to KnowledgeBase
- **internal/ctxharvester/harvester.go**: Semantic chunking with time decay and token budgets

## Key Insights
- `sync.Once` for singleton KnowledgeBase (thread-safe lazy init)
- FTS5 requires `glebarez/sqlite` (not `modernc`) due to orchestrator/database.go conflict
- Build tag `!nosqlite` on sqlite_store.go and test for future flexibility

## Provider Routing
- Priority: Anthropic > Gemini > OpenAI > DeepSeek > OpenRouter > Groq > LMStudio > Ollama
- `llm.AutoRoute()` selects first with valid API key
- `llm.AutoRouteWithModel()` allows model override
- `providers.GetProvider("name")` for explicit provider selection

## Tool Surface Coverage
- 136+ tools across: Claude Code (23), Crush (18), Pi (7+7), Gemini (7), OpenCode (14), Grok (6),
  Goose (4), Kimi (14), Cursor (5), Windsurf (2), Mistral (2), Smithery (2), Advanced (14)
- Foundation tools delegate to `foundation/pi/tools_native.go` for actual execution
