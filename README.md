# Borg

> Local AI operations control plane for MCP routing, provider fallback, session supervision, and a unified dashboard.

Borg is being stabilized toward a focused `1.0` release around four core capabilities:

- **MCP Master Router** — aggregate multiple MCP servers behind one endpoint
- **Model Fallback & Provider Routing** — switch providers/models when quotas or rate limits hit
- **Session Supervisor** — manage long-running external coding sessions
- **Web Dashboard** — one place to see system health, sessions, servers, and providers

## Current status

This repository is in an active cleanup and stabilization phase.

- Root `pnpm install` was verified successfully on Windows in this repo state.
- `docker compose up -d --build` was verified successfully on Windows with Docker Desktop running.
- Dockerized dashboard access was verified at `http://localhost:3001`.
- Root `pnpm run dev` now uses Borg's readiness launcher and may place the dashboard on a fallback web port if `3000` is already occupied; the core bridge remains on `http://127.0.0.1:3001`.
- The web app root redirects from `/` to `/dashboard`.

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

Notes:

- Opening `http://localhost:3001/` redirects to `http://localhost:3001/dashboard`.
- The MCP operator surface lives at `http://localhost:3001/dashboard/mcp`.
- The first image build can take several minutes because it builds the core packages and the Next.js dashboard inside Docker.

### 3) Start local development

```bash
pnpm run dev
```

What the launcher does:

- starts or reuses the Borg core bridge on `http://127.0.0.1:3001`
- waits for the authoritative `startupStatus` contract before reporting success
- builds missing official browser-extension artifacts when needed
- opens the dashboard after the web surface is actually ready

The dashboard usually starts on `http://127.0.0.1:3000/dashboard`, but it may move to a fallback port if `3000` is already taken. To confirm the active URLs after startup, run:

```bash
node scripts/verify_dev_readiness.mjs
```

Useful local probes:

- Core bridge / HTTP probe: `http://127.0.0.1:3001`
- Startup status: `http://127.0.0.1:3001/api/trpc/startupStatus?input=%7B%7D`
- MCP status: `http://127.0.0.1:3001/api/trpc/mcp.getStatus?input=%7B%7D`

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
- supervise external agent and CLI sessions
- manage provider and model selection with fallback
- expose system state through a practical dashboard

## Borg 1.0 focus

The repository still contains legacy and experimental surfaces, but the current product focus is intentionally narrow:

- **MCP Master Router**
- **Model Fallback & Provider Routing**
- **Session Supervisor**
- **Web Dashboard**

If a change does not make those workflows more reliable, more testable, or easier for a new user to run in five minutes, it is probably not a Borg 1.0 priority.

## Repository layout

```text
borg/
├── apps/web/              # Dashboard
├── apps/borg-extension/   # Official browser extension (Chrome/Edge + Firefox builds)
├── apps/extension/        # Legacy/compat browser extension surface
├── packages/core/         # MCP routing, orchestration, services
├── packages/cli/          # CLI entrypoint
├── packages/types/        # Shared types/schemas
├── tasks/                 # Active, backlog, completed task briefs
├── docs/                  # Architecture docs plus archived planning material
└── docker-compose.yml     # Containerized local stack
```

## Canonical docs

These are the root-level documents contributors should trust first:

- `AGENTS.md` — active repo operating directive
- `ARCHITECTURE.md` — high-level system design
- `ROADMAP.md` — current 1.0 / 1.5 / 2.0 milestones
- `CHANGELOG.md` — notable repo changes
- `tasks/active/` — current implementation work

Archive and compatibility material still exists for reference, but it is **not** the source of truth. Treat anything under `docs/archive/` as archive-only.

## Development notes

- Current stabilization work is focused on install, startup, and the core control-plane workflows.
- If you are contributing, prefer `tasks/active/` over older phase-based or parity-based planning artifacts.
- The dashboard landing page is `/dashboard`; use `/dashboard/mcp` when you specifically want the MCP router view.

### Optional startup debug logging

Normal `pnpm run dev` startup is intentionally quiet. If you want verbose Borg core boot/import logging while diagnosing startup issues, enable either of these environment flags before launching dev:

- `BORG_MCP_SERVER_DEBUG=1`
- `DEBUG=borg:mcp-server`

Windows PowerShell examples:

```powershell
$env:BORG_MCP_SERVER_DEBUG='1'; pnpm run dev
```

```powershell
$env:DEBUG='borg:mcp-server'; pnpm run dev
```

## License

MIT
