# Extension Surfaces Parity Matrix

> **Updated**: 2026-03-06 (Phase 146)
> **Version**: 2.7.108

This document tracks capability parity across the three Borg client surfaces: **Web Dashboard**, **Browser Extension**, and **VS Code Extension**.

---

## Capability Matrix

| Capability | Dashboard | Browser Ext | VS Code Ext | Status |
|---|:---:|:---:|:---:|---|
| **Connection & Health** | | | | |
| Core health check | ✅ tRPC | ✅ HTTP `/health` | ✅ WS auto-connect | 🟢 Shipped |
| Status bar indicator | N/A | N/A | ✅ StatusBarItem | 🟢 Shipped |
| Auto-reconnect | ✅ TRPCProvider | ✅ 5s retry | ✅ 5s retry | 🟢 Shipped |
| **Memory & Knowledge** | | | | |
| Page/context ingestion | ✅ Knowledge dashboard + browser activity feed | ✅ `ABSORB_PAGE` (Readability→Markdown) | ✅ `AIOS: Remember Selection` | 🟢 Shipped |
| URL ingestion | ✅ `knowledge.ingest` | ✅ Popup URL input → Core compatibility endpoint | ✅ Command + mini-dashboard action → Core compatibility endpoint | 🟢 Shipped |
| RAG file ingestion | ✅ `rag.ingestFile` | ❌ — | ❌ — | 🟡 Partial |
| RAG text/page ingestion | ✅ `rag.ingestText` + browser activity feed | ✅ Popup action → Core compatibility endpoint | ✅ `AIOS: Ingest Selection to RAG` | 🟢 Shipped |
| Memory save from context | ✅ Live core ingestion pipeline | ✅ `SAVE_CONTEXT` to Core | ✅ `aios.rememberSelection` | 🟢 Shipped |
| **Agent Interaction** | | | | |
| Research agent dispatch | ✅ `expert.research` | ❌ — | ✅ Sidebar + `AIOS: Run Agent` | 🟡 Partial |
| Coder agent dispatch | ✅ `expert.code` | ❌ — | ✅ Sidebar + `AIOS: Run Agent` | 🟡 Partial |
| Chat paste/submit | N/A | ✅ `PASTE_INTO_CHAT` | ✅ `PASTE_INTO_CHAT` (anti-hijack) | 🟢 Shipped |
| Chat history scraping | N/A | ✅ `SCRAPE_CHAT` (DOM) | ✅ Interaction-backed history + best-effort editor snapshot | 🟡 Partial |
| Tool execution proxy | N/A | ✅ `EXECUTE_TOOL` | ✅ `VSCODE_COMMAND` | 🟢 Shipped |
| **Monitoring** | | | | |
| Console log capture | ✅ Traffic Inspector live stream + filters | ✅ Interceptor + buffer | ✅ Output channel + live inspector mirror | 🟢 Shipped |
| CDP debug event stream | ✅ Traffic Inspector live stream + filters | ✅ `chrome.debugger.onEvent` relay | N/A | 🟢 Shipped |
| User activity tracking | N/A | ✅ Focus/click/keydown/scroll activity bridge | ✅ Selection + Document changes | 🟢 Shipped |
| Tab mirroring (screenshots) | ✅ Dedicated Browser dashboard mirror panel | ✅ JPEG mirror stream | ❌ — | 🟢 Shipped |
| **Browser-Specific** | | | | |
| Page scraping (Readability) | ✅ Dedicated Browser dashboard scrape panel | ✅ Full pipeline | N/A | 🟢 Shipped |
| Screenshot capture | ✅ Dedicated Browser dashboard screenshot panel | ✅ `captureVisibleTab` | N/A | 🟢 Shipped |
| Browser history search | ✅ Dedicated Browser dashboard search panel | ✅ `chrome.history` API | N/A | 🟢 Shipped |
| CDP debug proxy | ✅ Dedicated Browser dashboard CDP panel + inspector events | ✅ Attach/detach/command | N/A | 🟢 Shipped |
| Proxy fetch | ✅ Dedicated Browser dashboard request panel | ✅ `browser_proxy_fetch` | N/A | 🟢 Shipped |
| **VS Code–Specific** | | | | |
| Active editor info | ❌ — | N/A | ✅ `GET_STATUS` | 🟡 VSCode-only |
| Text selection extraction | ❌ — | N/A | ✅ `GET_SELECTION` | 🟡 VSCode-only |
| Terminal name reporting | ❌ — | N/A | ✅ `GET_TERMINAL` | 🟢 Shipped |
| Terminal content reading | ❌ — | N/A | ✅ Rolling buffer via terminal write event | 🟢 Shipped |
| Anti-hijack paste guard | N/A | N/A | ✅ 2s activity check | 🟢 Shipped |
| **Infrastructure** | | | | |
| Swarm orchestration | ✅ Full UI | ❌ — | ❌ — | 🟡 Dashboard-only |
| Director config | ✅ Full UI | ❌ — | ❌ — | 🟡 Dashboard-only |
| Submodule management | ✅ Full UI | ❌ — | ❌ — | 🟡 Dashboard-only |
| MCP server management | ✅ Full UI | ❌ — | ❌ — | 🟡 Dashboard-only |

---

## Parity Milestones

### Milestone 1: Extension Hardening (P2)
- [x] Replace hardcoded `localhost:3001` in browser extension `background.ts` with configurable URL
- [x] Add configurable WS URL to browser extension (matching VS Code's `borg.coreUrl` config)
- [x] Add error boundary / offline UI to browser extension popup

### Milestone 2: Cross-Surface Intelligence (P3)
- [x] Forward RAG ingestion capability to browser extension (ingest current page via Core compatibility endpoint)
- [x] Add research/code agent dispatch UI to VS Code sidebar
- [x] Surface console log buffer from browser extension in dashboard Traffic Inspector
- [x] Add terminal content reading to VS Code extension (beyond name-only)

### Milestone 3: Full Parity (P4)
- [x] Unified WebSocket protocol specification document
- [x] Browser extension options page for URL/port configuration
- [x] VS Code webview panel for mini-dashboard (status, recent tasks, quick actions)

---

## Architecture Notes

| Surface | Transport | Entry Point |
|---|---|---|
| Dashboard | tRPC over HTTP (batch) | `apps/web` → `@borg/ui` TRPCProvider |
| Browser Extension | WS + Chrome messaging | `apps/extension/src/background.ts` |
| VS Code Extension | WS (node `ws` lib) | `packages/vscode/src/extension.ts` |

All three surfaces connect to the same Borg Core instance. The browser and VS Code extensions use WebSocket for real-time bidirectional communication, while the dashboard uses tRPC HTTP batch for request/response patterns.
