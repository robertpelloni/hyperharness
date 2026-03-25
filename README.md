<div align="center">

# 🤖 BORG

## 🧪 ALL CAPS MAD SCIENCE

### The Neural Operating System for AI Agents

> *Connect everything. Orchestrate anything. Remember forever.*

<br />

[![CI](https://github.com/robertpelloni/borg/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/robertpelloni/borg/actions/workflows/ci.yml)
[![CodeQL](https://github.com/robertpelloni/borg/actions/workflows/codeql.yml/badge.svg)](https://github.com/robertpelloni/borg/actions/workflows/codeql.yml)
[![Release](https://github.com/robertpelloni/borg/actions/workflows/release.yml/badge.svg)](https://github.com/robertpelloni/borg/actions/workflows/release.yml)
[![License: MIT](https://img.shields.io/badge/license-MIT-22c55e?style=flat)](./LICENSE)

<br />

![TypeScript](https://img.shields.io/badge/TypeScript-5.x-3178c6?style=flat&logo=typescript&logoColor=white)
![Node.js](https://img.shields.io/badge/Node.js-22+-339933?style=flat&logo=nodedotjs&logoColor=white)
![Next.js](https://img.shields.io/badge/Next.js-15-000000?style=flat&logo=nextdotjs&logoColor=white)
![tRPC](https://img.shields.io/badge/tRPC-11-398ccb?style=flat&logo=trpc&logoColor=white)
![pnpm](https://img.shields.io/badge/pnpm-10-f69220?style=flat&logo=pnpm&logoColor=white)

<br />

![GitHub stars](https://img.shields.io/github/stars/robertpelloni/borg?style=flat&color=f59e0b&logo=github)
![GitHub forks](https://img.shields.io/github/forks/robertpelloni/borg?style=flat&color=0ea5e9&logo=github)
![Commits/month](https://img.shields.io/github/commit-activity/m/robertpelloni/borg?style=flat&label=commits%2Fmonth&color=8b5cf6)
![Last commit](https://img.shields.io/github/last-commit/robertpelloni/borg?style=flat&color=06b6d4)
![Repo size](https://img.shields.io/github/repo-size/robertpelloni/borg?style=flat&color=f97316)

</div>

---

## ⚡ What is Borg?

**Borg is the mecha-suit your AI stack has been missing.** It's a local-first control plane that sits between your AI agents and their infrastructure, turning scattered MCP tools, fragmented memories, and isolated CLI harnesses into one cohesive, autonomous operating system.

```
┌─────────────────────────────────────────────────────────────┐
│                    🧠 BORG CONTROL PLANE                     │
│                                                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │  Memory   │  │   MCP    │  │ Council  │  │ Director │    │
│  │  Engine   │  │  Router  │  │ (Debate) │  │(Autopilot│    │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘    │
│       │              │              │              │          │
│  ┌────┴──────────────┴──────────────┴──────────────┴────┐    │
│  │              60+ tRPC Services                        │    │
│  │    Sessions · Tools · Memory · RAG · Billing · MCP    │    │
│  └───────────────────────────────────────────────────────┘    │
│                                                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │  Claude   │  │   GPT    │  │  Gemini  │  │  Local   │    │
│  │  Code     │  │  5.4     │  │  3.1     │  │  Ollama  │    │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘    │
└─────────────────────────────────────────────────────────────┘
        ▲              ▲              ▲              ▲
   ┌────┴────┐   ┌────┴────┐   ┌────┴────┐   ┌────┴────┐
   │ VS Code │   │ Browser │   │   CLI   │   │  Expo   │
   │Extension│   │Extension│   │Harnesses│   │ Mobile  │
   └─────────┘   └─────────┘   └─────────┘   └─────────┘
```

### 🎯 The Problem

You have 50+ MCP tools, 4 different AI CLI harnesses, browser tabs everywhere, and zero memory between sessions. Your agents are powerful but **fragmented**.

### 💡 The Solution

Borg unifies everything behind **one endpoint**, **one memory system**, and **one orchestration layer**:

| Pillar | What It Does |
|--------|-------------|
| **🔌 MCP Master Router** | Aggregate hundreds of MCP servers behind one endpoint. Auto-conflict resolution, progressive tool disclosure, semantic search. |
| **🧠 Omniscient Memory** | Multi-tier memory (vector + graph + context harvesting). Agents never forget. Sessions resume instantly. |
| **⚡ One-Shot Execution** | Describe an objective → Borg finds the tool, maps the args, executes. One turn, not five. |
| **🏛️ Council of Agents** | Multi-model debate (Claude × GPT × Gemini × Local). Consensus-driven decisions with quota-aware routing. |
| **🤖 Director Autopilot** | Background daemon that monitors, heals, and autonomously drives your development workflow. |
| **📊 Mission Control** | 59-page WebUI dashboard with real-time observability, billing, swarm monitoring, and tool management. |

---

## 🚀 Quick Start

### Prerequisites
- **Node.js 22+** and **pnpm 10+** (required — enforced by `packageManager` lock)

### Option A: Docker
```bash
docker compose up --build
```

### Option B: Local Development
```bash
pnpm install
pnpm run dev
# Dashboard → http://localhost:3001/dashboard
```

### Connect Your AI Agent
Add Borg as an MCP server in any compatible client:

```json
{
  "mcpServers": {
    "borg": {
      "command": "node",
      "args": ["/path/to/borg/packages/core/dist/index.js"]
    }
  }
}
```

---

## 🧠 Key Capabilities

### Memory That Never Forgets
- **Automatic Context Harvesting** — Extracts facts, entities, and relationships from every interaction
- **LanceDB Vector Storage** — High-speed local vector search for semantic retrieval
- **Graph Knowledge** — Tracks relationships between files, tasks, and concepts
- **Session Export/Import** — Auto-detects 8 formats (Claude Code, Cursor, Aider, OpenCode, etc.)
- **Google Workspace Integration** — Syncs with Docs, Gmail, and Drive

### One-Shot Tool Execution
```
Agent → auto_call_tool({
  objective: "Search my emails for the invoice from last Tuesday",
  context: "Tuesday was March 10th"
})

Borg → Finds gmail tool → Maps arguments → Executes → Returns result
```

### The Council: Multi-Model Debate
Configure your routing hierarchy via the dashboard. When one provider exhausts its quota, Borg cascades automatically:

```
1. Claude Opus 4      ← Primary (highest intelligence)
2. GPT-5.4 Turbo      ← Fallback (fast, reliable)
3. Gemini 3.1 Pro      ← Secondary (cost-efficient)
4. Local Gemma3        ← Emergency (free, on-device)
```

### Director Autopilot
- **Background Daemon** — Monitors chat logs and terminal outputs
- **Self-Healing** — Analyzes crashes, generates fixes, proposes patches
- **Session Handoff** — Automatically creates bootstrap prompts for seamless agent resumption

---

## 📂 Repository Layout

```
borg/
├── apps/
│   ├── web/              → Next.js 15 Mission Control Dashboard (59 pages)
│   ├── borg-extension/   → Chrome/Firefox browser bridge
│   ├── vscode/           → VS Code extension
│   └── maestro/          → Electron desktop shell
├── packages/
│   ├── core/             → Engine: 76 tRPC routers, 135+ services
│   ├── ai/               → Model routing, quota tracking, LLM abstraction
│   ├── cli/              → Command-line interface
│   ├── ui/               → Shared React component library
│   └── types/            → Shared TypeScript types
├── docs/                 → Architecture docs, deployment guides
├── scripts/              → Build, version, and deployment automation
└── external/             → MCP submodules and integrations
```

---

## 🔐 Security & Trust Model

Borg is a **local-first** control plane. Your data never leaves your machine unless you connect external API providers.

- **Local process boundary** — Borg runs with your user permissions
- **No hosted backend required** — Fully operational offline
- **Secrets management** — API keys injected via env/config, never committed
- **Audit trail** — Full logging of tool executions and agent actions

---

## 📊 Stats & Activity

<div align="center">

<a href="https://github.com/robertpelloni/borg">
  <img alt="Repo stats" height="165" src="https://github-readme-stats.vercel.app/api/pin/?username=robertpelloni&repo=borg&theme=tokyonight&hide_border=true&show_icons=true" />
</a>

[![Activity Graph](https://github-readme-activity-graph.vercel.app/graph?username=robertpelloni&repo=borg&theme=tokyo-night&hide_border=true&area=true&custom_title=Borg+Commit+Activity)](https://github.com/robertpelloni/borg/commits/main)

[![GitHub Streak](https://streak-stats.demolab.com/?user=robertpelloni&theme=tokyonight&hide_border=true)](https://github.com/robertpelloni/borg)

</div>

---

## 🗺️ Roadmap

| Phase | Status | Description |
|-------|--------|-------------|
| **I** Omniscient Memory & RAG | ✅ Complete | Multi-tier memory, context harvesting, Google Workspace integration |
| **J** Universal Integrations | ✅ Complete | Browser extensions, VS Code plugin, IDE hooks |
| **K** Intelligent Routing | ✅ Complete | Quota-aware cascading, provider fallback, billing |
| **L** Ultimate Coding Harness | ✅ Complete | 11 CLI harness types, auto-start/restart, cloud dev |
| **M** MCP Aggregation | ✅ Complete | Universal directory, traffic inspection, Code Mode |
| **N** Marketplace & Mesh | 🔨 Next | P2P agent federation, community tool marketplace |

---

## 🤝 Contributing

1. Fork the repo
2. Create a feature branch (`git checkout -b feature/amazing`)
3. Run `pnpm install` then `pnpm run build` in `apps/web` to verify
4. Commit and push
5. Open a Pull Request

> **Important:** Use `pnpm v10` (enforced). Import UI components from `@borg/ui`, never `@/components/ui/*`.

---

## ⚖️ License

MIT. See [`LICENSE`](./LICENSE) for details.

---

<div align="center">

**Built with 🧠 by [Robert Pelloni](https://github.com/robertpelloni)**

*Resistance is futile. You will be assimilated.*

</div>