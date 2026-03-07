# Borg

> Local AI operations control plane for MCP routing, provider fallback, session supervision, and a unified dashboard.

Borg is being stabilized toward a focused `1.0` release around four core capabilities:

- **MCP Master Router** — aggregate multiple MCP servers behind one endpoint
- **Model Fallback & Provider Routing** — switch providers/models when quotas or rate limits hit
- **Session Supervisor** — manage long-running external coding sessions
- **Web Dashboard** — one place to see system health, sessions, servers, and providers

## Current status

This repository is in an active cleanup/stabilization phase.

- Root `pnpm install` was verified successfully on Windows in this repo state.
- The Docker Compose file parses cleanly.
- Full container boot requires a running Docker engine/Desktop.

## Quick start

### Prerequisites

- Node.js 20+
- `pnpm` 10+
- Docker Desktop or another working Docker engine if you want the containerized stack

### 1) Clone and install

```bash
git clone https://github.com/robertpelloni/borg.git
cd borg
pnpm install
```

### 2) Start with Docker Compose

```bash
docker compose up --build
```

Expected URLs once the stack is up:

- Dashboard: `http://localhost:3001`
- Core API: `http://localhost:3000`

### Windows note

If Docker reports an error like:

- `open //./pipe/dockerDesktopLinuxEngine: The system cannot find the file specified`

that usually means Docker Desktop is not running yet. Start Docker Desktop, wait for the engine to come up, then rerun:

```bash
docker compose up --build
```

## What Borg is

Borg is an **orchestrator**, not a clone of every AI tool.

It is intended to:

- route MCP tools cleanly and safely
- supervise external agent/CLI sessions
- manage provider/model selection and fallback
- expose system state through a practical dashboard

## Repository layout

```text
borg/
├── apps/web/              # Dashboard
├── apps/extension/        # Browser extension
├── packages/core/         # MCP routing, orchestration, services
├── packages/cli/          # CLI entrypoint
├── packages/types/        # Shared types/schemas
├── tasks/                 # Active, backlog, completed task briefs
├── docs/                  # Architecture and archived planning docs
└── docker-compose.yml     # Containerized local stack
```

## Important docs

- `AGENTS.md` — active repo operating directive
- `ARCHITECTURE.md` — high-level system design
- `ROADMAP.md` — current 1.0 / 1.5 / 2.0 milestones
- `tasks/active/` — current implementation work
- `CHANGELOG.md` — notable repo changes

## Development notes

- The repo still contains legacy/reference material outside the Borg 1.0 core path.
- Current stabilization work is focused on making install, startup, and the core control-plane workflows reliable.
- If you are contributing, prefer the active tasks in `tasks/active/` over legacy phase-based documents.

## License

MIT
