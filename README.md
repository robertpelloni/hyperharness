# Borg (v2.7.333)

> The Unified AI Operations Control Plane.

Borg is a high-performance local control plane designed to sit between your AI agents and their underlying infrastructure. It provides a single point of truth for MCP routing, provider fallback, session supervision, and real-time observability.

## 🚀 Key Capabilities

- **One-Shot Discovery & Execution**: Instead of multi-turn tool discovery, Borg's `auto_call_tool` allows agents to search for, map, and execute any tool in a single turning point.
- **Universal MCP Master Router**: Aggregate hundreds of MCP servers behind one endpoint. Borg manages connections, tool conflicts, and namespacing automatically.
- **Dynamic Tool Advertising**: To stay within LLM function declaration limits (like Gemini's 512-tool cap), Borg only advertises core Meta-Tools by default. Standard and skill tools are summoned on-demand via semantic search.
- **Provider Fallback & Quota Management**: Chain multiple LLM providers (OpenAI, Anthropic, Gemini, Ollama) with automatic failover when rate limits or quotas are reached.
- **Session Supervision**: Attach Borg to external agent sessions (Claude Code, OpenCode, Jules) to monitor logs, broadcast messages, and enforce safety policies.
- **Mission Control Dashboard**: A comprehensive web UI for managing system health, MCP catalogs, active sessions, and provider billing.

---

## 🏁 Quick Start

### Prerequisites
- Node.js 20+
- `pnpm` 10+
- Docker Desktop (optional for containerized stack)

### Option A: Docker (Recommended)
```bash
docker compose up --build
```
- **Dashboard**: `http://localhost:3001/dashboard`
- **Core API**: `http://localhost:3000`

### Option B: Local Development
```bash
pnpm install
pnpm run dev
```
Borg will boot up, verify subsystem readiness, and provide a local dashboard URL (usually `http://127.0.0.1:3000/dashboard`).

---

## 🛠️ Configuration

1. **Environment Variables**:
   ```bash
   cp apps/web/.env.example apps/web/.env.local
   cp packages/core/.env.example packages/core/.env
   ```
2. **Provider Keys**: Add your keys to `packages/core/.env`. Borg degrades gracefully; missing providers are skipped in the fallback chain.
3. **MCP Config**: Edit `mcp.json` in the root to add or remove servers. This follows the standard MCP client configuration format.

---

## 🏗️ Core Architecture (The Meta-Tool Pattern)

To ensure maximum efficiency and compatibility with modern LLMs, Borg utilizes a "Meta-Tool" first approach:

1. **Discovery**: `search_tools(query)` ranks tools across all servers using semantic embedding.
2. **Loading**: `load_tool(name)` hydrates a specific tool's schema into the active session.
3. **Execution**: `auto_call_tool(objective, context)` combines discovery and execution into one turn.

This architecture allows Borg to scale to thousands of available tools without exceeding context windows or API limits.

---

## 📂 Repository Layout

Borg is structured as a lean, demo-ready monorepo:

- `apps/web`: Next.js Mission Control dashboard.
- `apps/borg-extension`: Browser side-panel integration.
- `packages/core`: The core engine (MCP routing, LLM orchestration, proxy).
- `packages/cli`: Command-line interface for direct interaction.
- `packages/ui`: Shared Tailwind + Shadcn/ui component library.
- `packages/memory`: Vector memory (LanceDB) integration.
- `archive/`: Compressed history, old documentation, and unused modules.

---

## 📄 Documentation

- [CHANGELOG.md](CHANGELOG.md) — Release history and migration notes.
- [VERSION](VERSION) — The current semantic version.
- `archive/docs/` — Comprehensive architectural archives and legacy planning.

---

## ⚖️ License

MIT. See `LICENSE` for details.
