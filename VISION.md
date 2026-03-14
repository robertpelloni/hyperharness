# VISION.md - Borg Long-Term Direction

> **Current release track**: `0.9.0-beta`
> **Near-term mission**: Ship a trustworthy Borg `1.0`
> **Long-term ambition**: Become the best local AI operations control plane for builders who want one place to run, supervise, and route their AI tooling.

---

## 1. The long view

Borg exists to make AI-assisted development feel operable instead of chaotic.

The long-term vision is not "own every tool" or "clone every interface." It is to give a single operator one dependable place to:

- route tools through a unified MCP-aware control plane
- switch between providers and models without losing momentum
- supervise long-running coding sessions and recover from crashes
- observe what the system is doing in real time
- gradually add memory, agent orchestration, and browser/cloud integrations without breaking the core

In short: **orchestrate everything, reimplement almost nothing unless it must exist inside Borg to make the whole system reliable.**

---

## 2. What success looks like

In the near term, success is practical:

- a new user can clone the repo, install it, start it, and reach the dashboard quickly
- the dashboard tells the truth about MCP servers, providers, sessions, and system health
- Borg can keep working when one provider hits quota or one supervised process crashes
- external AI clients can connect once and use many tools through Borg cleanly

In the longer term, success becomes leverage:

- teams can plug Borg into their preferred editors and AI clients
- provider choice becomes a routing policy, not a workflow disruption
- session history and memory become durable operational context
- agent workflows become supervised, inspectable, and recoverable rather than magical black boxes

---

## 3. Product philosophy

### Orchestrate, don't absorb

Borg should coordinate other tools, not turn into a museum of embedded clones.

That means:

- adapters over parity checklists
- capability contracts over brand-specific duplication
- clear boundaries over sprawling internal copies
- upstream references where useful, runtime ownership where necessary

### Make the boring parts solid

The most valuable features are often the least glamorous:

- restart behavior
- config sync
- health visibility
- quota awareness
- deterministic routing
- session persistence
- truthful logs

If Borg does those better than everyone else, it becomes the place operators trust.

### Prefer focus over mythology

Big ideas are welcome. Infinite scope is not.

The vision only matters if each release leaves the product easier to install, easier to operate, and easier to explain.

---

## 4. End-state capability picture

Over time, Borg should grow into a cohesive control plane with five durable layers.

### 4.1 MCP operations layer

Borg should be the easiest way to:

- aggregate many MCP servers behind one endpoint
- keep tool routing collision-safe while discovery happens through semantic search and grouping
- inspect live JSON-RPC traffic
- manage health, lifecycle, and config sync

### 4.2 Provider routing layer

Borg should make model usage resilient by handling:

- normalized provider authentication
- quota and rate-limit awareness
- routing strategies by cost, quality, and task type
- fallback without losing request continuity

### 4.3 Session supervision layer

Borg should supervise external coding tools and agent loops with:

- isolated child processes
- restart policies and crash recovery
- worktree isolation
- persistent session state, logs, and operator controls

### 4.4 Operator dashboard layer

The dashboard should remain an operations console, not a decorative feature dump.

It should help operators answer:

- what is running?
- what is failing?
- what is costing money?
- what was restarted?
- what tool traffic just happened?

### 4.5 Expansion layer

Once the control plane is trustworthy, Borg can responsibly expand into:

- memory and retrieval systems
- browser and editor companion surfaces
- multi-agent orchestration
- cloud-dev bridges
- richer automation workflows

Those are important, but they should compound on a stable base instead of becoming excuses to avoid shipping one.

---

## 5. Milestone horizon

### Borg 1.0

Ship the four core features:

1. MCP Master Router
2. Model Fallback & Provider Routing
3. Session Supervisor
4. Web Dashboard

### Borg 1.5

Add memory and context systems that support real workflows without destabilizing the core product.

### Borg 2.0

Add multi-agent orchestration and broader integrations only after the control plane has earned trust with external users.

---

## 6. Non-goals

The long-term vision explicitly does **not** require:

- cloning every AI coding tool feature-for-feature
- shipping every experimental idea at once
- hiding instability behind clever branding
- measuring meta-analytics that do not improve operator outcomes
- keeping massive reference collections in the main runtime path

If Borg ever has to choose between being impressive on paper and dependable in practice, dependable wins.

---

## 7. The standard for future decisions

Future proposals should make Borg more of the following:

- **runnable**
- **observable**
- **recoverable**
- **composable**
- **useful to a stranger in under five minutes**

That is the real vision.

Everything else is garnish.



## 2. VISION.md

```markdown
# VISION.md — Borg Long-Term Direction

> **Current version:** `1.0.0-alpha.1`
> **Near-term mission:** Ship a trustworthy Borg 1.0
> **Long-term ambition:** The best local AI operations control plane for builders.

---

## 1. What Borg Is

Borg exists to make AI-assisted development **operable instead of chaotic**.

It gives a single operator one dependable place to:

- Route tools through a unified MCP-aware control plane
- Switch between providers and models without losing momentum
- Supervise long-running coding sessions and recover from crashes
- Observe what the system is doing in real time
- Gradually add memory, agent orchestration, and integrations without
  breaking the core

**In short:** orchestrate everything, reimplement almost nothing.

## 2. The Primitive: Session

The fundamental unit of Borg is the **Session** — a supervised, checkpointable,
resumable AI workflow. Sessions are to Borg what containers are to Kubernetes.

A session encapsulates:
Session {
process // Child process or remote connection
state // Serializable, resumable conversation + tool state
routing_policy // Model preferences, fallback chain, cost budget
loaded_tools // Currently active MCP tools
memory_refs // Pointers to relevant long-term context
supervisor // Council config, restart policy
logs // Structured event stream
metrics // Cost, latency, token usage
checkpoints // Git commits + serialized state snapshots
}



**If you can serialize it, persist it, restore it, and hand it to a different
model — that's the Borg primitive.**

## 3. The Five Engines

### 3.1 MCP Master Router

Aggregates N MCP servers behind one endpoint. Models see a **tiny permanent
meta-tool surface** (6 tools), not hundreds:

- `search_tools` — Ranked semantic discovery
- `load_tool` — Activate a capability
- `get_tool_schema` — Inspect before calling
- `list_loaded_tools` — What's active
- `unload_tool` — Release
- `run_code` — Code-mode escape hatch

**Progressive disclosure tiers:**

| Tier | Strategy | Example |
|------|----------|---------|
| 0: Meta | Always visible | `search_tools`, `load_tool` |
| 1: Core | Always visible | `read_file`, `write_file`, `execute_command` |
| 2: Lazy | On-demand | Database, browser, Slack tools |
| 3: Deferred | Binary not spawned until first call | Heavy tools |

**Auto-load thresholds:**

- Confidence > 0.85 → silent auto-load
- Confidence 0.60–0.85 → present ranked alternatives
- Confidence < 0.60 → require explicit search

**Key insight from our MCP aggregator research:**

> "Most aggregators solve plumbing. Borg solves selection friction.
> The model should almost never face more than a handful of visible choices."

### 3.2 Provider Router

Normalized API access to all model providers with automatic quota-aware fallback.

**Routing equation:**

$$m^* = \argmax_{m \in M} \left( S(m) \cdot R(m) \cdot \left[ w_1 Q(m,T) - w_2 C(m,T) - w_3 L(m,T) \right] \right)$$

Where:
- $S(m) \in \{0, 1\}$ — health status (circuit breaker)
- $R(m) \in [0, 1]$ — remaining quota ratio
- $w_1, w_2, w_3$ — user-configurable weights

**In practice:** ordered waterfall fallback chains per task type. Exhaust all
providers of the best model before cascading to the next tier.

### 3.3 Session Supervisor

Manages long-running coding sessions as supervised child processes with
**crash-safe checkpoint/restore**.

**Checkpoint = git commit + serialized conversation state + environment snapshot.**

Recovery: checkout the commit, rehydrate the conversation, inject into a fresh
process with the next available model. The session never dies.

### 3.4 Council (Autonomous Supervisor)

When a session is idle and waiting for human input, the Council injects
intelligent continuation. Three modes:

1. **Cheerleader** (~free) — Rotates through "Keep going" prompts
2. **Informed** (~$0.001) — Flash model reads TODO.md and suggests next step
3. **Expert Panel** (~$0.05) — Multi-model debate on direction

### 3.5 Operator Dashboard

A web-based operations console that **tells stories, not just state**.

Instead of: "MCP Server X: unhealthy"
Show: "File search became unavailable 4 min ago. Borg fell back to grep.
       Here's the log."

## 4. Product Philosophy

### Orchestrate, don't absorb

Borg coordinates tools. It does not clone them. Adapters over parity checklists.
Capability contracts over brand-specific duplication.

### Make the boring parts solid

Restart behavior. Config sync. Health visibility. Quota awareness. Deterministic
routing. Session persistence. Truthful logs. If Borg does those better than
everyone else, it becomes the place operators trust.

### Prefer focus over mythology

Big ideas are welcome. Infinite scope is not. The vision only matters if each
release leaves the product easier to install, easier to operate, and easier
to explain.

## 5. Milestone Horizon

### Borg 1.0 — "It Runs and Recovers"
- MCP Master Router with progressive disclosure
- Provider Router with waterfall fallback
- Session Supervisor with checkpoint/restore
- Web Dashboard
- Council (cheerleader mode)

### Borg 1.5 — "It Remembers"
- Memory plugin system (file, vector DB, graph)
- Semantic tool search / tool RAG
- Context harvesting and compaction
- Browser extension (memory capture)

### Borg 2.0 — "It Orchestrates"
- Multi-model council (expert panel)
- Architect-implementer-reviewer workflows
- Cloud dev session management
- IDE extensions
- Skill and prompt libraries

## 6. Non-Goals

- Cloning every AI coding tool feature-for-feature
- Shipping every experimental idea at once
- Hiding instability behind clever branding
- Measuring meta-analytics that don't improve operator outcomes

**If Borg ever has to choose between being impressive on paper and
dependable in practice, dependable wins.**

## 7. The Standard

Future proposals should make Borg more:

- **Runnable**
- **Observable**
- **Recoverable**
- **Composable**
- **Useful to a stranger in under five minutes**

That is the real vision. Everything else is garnish.