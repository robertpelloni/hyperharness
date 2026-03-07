# Borg Roadmap

This roadmap replaces the legacy phase system with three milestone tracks. Historical phase planning now lives in `docs/archive/ROADMAP_LEGACY.md`.

## Borg 1.0 — Actually Works

- [x] Root `AGENTS.md` aligned to the focused Borg 1.0 directive
- [x] Legacy roadmap archived to `docs/archive/ROADMAP_LEGACY.md`
- [x] Task file system created under `tasks/`
- [ ] Verify clean install from fresh clone (`pnpm install`, `docker-compose up --build`)
- [ ] MCP Master Router
  - Aggregate multiple MCP servers behind one Borg endpoint
  - Namespace tools without collisions
  - Restart crashed child servers without disturbing healthy peers
  - Expose traffic inspector data for the dashboard
- [ ] Model Fallback & Provider Routing
  - Track provider auth and quota state
  - Fail over automatically on quota/rate-limit exhaustion
  - Support cheapest, best, and round-robin strategies
  - Support task-type-aware routing
- [ ] Session Supervisor
  - Spawn external CLI sessions
  - Auto-restart with exponential backoff
  - Persist sessions across Borg restarts
  - Isolate parallel coding sessions with worktrees
- [ ] Dashboard MVP
  - Overview panel
  - MCP Router panel
  - Sessions panel
  - Providers panel
- [ ] 30+ tests passing in CI
- [ ] Docker one-command install
- [ ] 5-minute quickstart verified against reality
- [ ] GitHub Release `v1.0.0`

## Borg 1.5 — Remembers Things

- [ ] Memory system with one vector backend
- [ ] Basic RAG ingestion, chunking, embedding, and retrieval
- [ ] Browser extension export and context capture workflows
- [ ] Context compaction and harvesting
- [ ] Five CLI tool adapters: Aider, Claude Code, OpenCode, Codex, Gemini CLI

## Borg 2.0 — The Swarm Awakens

- [ ] Multi-agent orchestration via Director/Council flows
- [ ] Multi-model debate and consensus
- [ ] Memory multi-backend plugin system
- [ ] Cloud dev integration (Jules, Devin, related tooling)
- [ ] Revisit Go supervisor extraction only if Node.js supervision proves insufficient
- [ ] Re-enable focused swarm work from legacy phases 81-95 when 1.x is stable
- [ ] Mobile dashboard

## Out of Scope for Now

- P2P distributed mesh as a primary roadmap track
- Feature parity with every CLI tool at once
- Recursive scraping and assimilation work
- Analytics-about-analytics detours

