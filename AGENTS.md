# Agent Instructions

> **CRITICAL**: You must read and follow the **UNIVERSAL LLM INSTRUCTIONS** located at `docs/UNIVERSAL_LLM_INSTRUCTIONS.md`.

This file serves as the entry point for all autonomous agents working on the Borg project.

## Quick Links
- **Universal Rules**: `docs/UNIVERSAL_LLM_INSTRUCTIONS.md`
- **Roadmap**: `ROADMAP.md`
- **Architecture**: `docs/DESIGN.md`
- **Submodules**: `docs/SUBMODULE_DASHBOARD.md`

## Agent Protocols
1. **Autonomy**: Solve problems independently.
2. **Verification**: Always test your code.
3. **Documentation**: Update Changelog and Roadmap.
4. **Self-Correction (The Healer)**: If a tool fails, analyze the error and attempt a fix.
5. **Evolution (The Darwin Protocol)**: Experiment with prompt variations if results are suboptimal.

Proceed with your assigned task.



The following features MUST be EXTREMELY WELL IMPLEMENTED and EXTREMELY WELL REPRESENTED by the UI, CLI, TUI, and documentation, including UI labeling and tooltips:

mcp router/aggregator, mcp server which is mcp client which combines many mcp servers into one "master" mcp server, handles mcp session, instanciation, lifecycle auto start, restart, making sure it is available to multiple clients without multiple instances, making sure it does not time out and is responsive, measures latency, groups mcp servers into namespaces, groups tools into namespaces, enable/disable individual tools/servers, inspect mcp traffic, show activity, logging, mcp tool call chaining, TOON format, code mode, tool automatic renaming, tool reranking, tool semantic search, progressive tool disclosure, tool rag, context inspector, keep alive, heartbeat, mcp directory, automatic mcp install, environment variable and secrets management, mcp client config file autoconfig, auto write config file format, import/export configs, autodetect all configs on system, maintain/manage all configs on system, save sets of configs, wipe all client configs, set all client configs to x, etc, manage databases, handle auth. 
mcp server chains other mcps. 
can watch mcp traffic. 
reduce context language. 
reduce context syntax. 
summarize session. 
automatic context harvesting. 
compaction. 
context pruning, session pruning, memory pruning, memory reranking. 
import/export session. 
import/export memories. 
autodetect session/memories. 
code execution sandbox. 
code indexing, understanding. 
semantic code search. 
lsp servers. 
ast, see graph view ast code. 
ripgrep. 
code chunking. 
RAG, various RAG techniques, intake documents, ocr, summarize. 
embed skills, specs/speckit, tasks/task manager, bmad, agents, swarm
plugins, cli sdk, ai sdk, a2a sdk, mcp sdk, acp sdk, agent sdk (many providers). 
orchestrator. 
CLIs/TUIs, webUIs. 
resource management, lazy load, health checks, auto start, auto restart, timeout, latency, responsiveness, single-instance multiple clients, monitor memory usage, cpu usage, long running process. 
proxy for stdio as remote. 
proxy for remote as stdio. 
proxy for sse as streaming-http. 
proxy for streaming-http as stdio. 
manage oauth, bearer token. 
manage secrets, environment variables, paths. 
manage env variable expansion, .env files. 
subagents, long running, subcontexts, timeouts, multiple models, mutiple cli, orchestration, harness/cli/tiu as mcp. 
multi-model debate, multi-model consensus, share context between models, multiple models pair programming with each other, architect-implementer. 
mcp groups. 
tool groups. 
tool renaming, minimize context. 
tool semantic search/tool rag. 
tool lazy load, dynamic progressive tool disclosure. 
tools as code, code mode. 
mcp tool chaining. 
TOON format, context saving. 
inspect context, context makeup,
usage, billing, dashboard, credits, api vs auth, resets on day at time. 
intelligent selection of models based on credits, free tier, subscription, provider, switch between all providers of gemini 3 pro, then all providers of codex 5.2, then all opus 4.5. 
plan with model x, implement with model y. 
supervisor, council, autopilot. 
skills, skill convert, skill rerank, skill improve. 
prompt library, system prompts, prompt templates, improve prompts, jailbreaks. 
personas. 
subagent definitions, subagent collections. 
agents.md, claude.md, copilot-instructions.md, gpt.md, llm-instructions.md, gemini.md, grok.md, warp.md. 
changelog.md. 
version, version.md. 
vision.md. 
readme.md. 
handoff.md, init.md. 
superpowers. 
beads. 
specs, tasks, bmad. 
memory. 
short term. 
long term. 
browser extension: store memory, browser-use, console/debug log, mcp superassistant, connect mcp to browser chat, connect universal memory to browser chat, export browser sessions, export browser memory, browser history. 
connect to all interfaces all models. 
all models connect to memory. 
import/export all memory, sessions, chat history. 
add to memory. 
automatically add memories about x topic. 
web search. 
computer-use. 
browser-use. 
agentic loop, autopilot, supervisor, fallbacks, council, consensus. 
manage cloud dev sessions. 
manage local cli sessions. 
cli has cli, tiu, webui by default. 
connect to cli with webui, dashboard connects to cli. 
no tui, only webui, only multi-session dashboard in web. 
remembers repo folders. 
list of all repo folders, workspace, last session. 
auto start previous sessions. 
pause all sessions. 
import/export cloud dev memory, session history. 
transfer task to cloud dev. 
broadcast message to all cloud dev sessions. 
broadcast message to all cli sessions. 
remote management of sessions, mobile control of sessions, roo remote. 
auto approve cloud dev plans. 
start new session for expired (30 day) cloud dev sessions, inject session history, prune session log, archive old session. 
continue paused/stalled cloud dev session. 
web terminal like opencode web ui terminal. 
cli activity log. 
memory dashboard like supermemory console. 
automatic adding of memories during session? all user input gets added as memories if there is new information?. 
notebooklm integration. 
notebooklm open source functionality. 