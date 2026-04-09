---
description: Your role, own self-concept, personality traits, and behavioral guidelines
limit: 2000
---
# HyperHarness Agent Context

## Current: v0.3.0 - 509 tests, 33 packages
- Multi-provider LLM routing (8 providers: OpenAI, Anthropic, Gemini, DeepSeek, OpenRouter, Groq, LMStudio, Ollama)
- 136+ tool surfaces across 15+ source harnesses
- Key subsystems: cache, eventbus, healer, ctxharvester, git service, tool detector, council, A2A broker
- Build: `go build -buildvcs=false ./...` with Go 1.26.2
- SQLite: glebarez/sqlite (NOT modernc - conflicts with orchestrator)
- Version in VERSION file, read by internal/buildinfo

## Critical Type Names (avoid clashes)
- CouncilDelegate (not CouncilMember for role constant)
- PriorityQueue (not TaskQueue - SQLite version exists)
- SessionTodoStore (not TodoStore - crush version exists)

## Architecture
- Tools delegate to foundation/pi for execution
- Provider priority: Anthropic > Gemini > OpenAI > DeepSeek > OpenRouter > Groq > LMStudio > Ollama
- EventBus: exact match + wildcard (tool:*) + global listeners
