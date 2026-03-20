# Borg Agents & Roles

This document defines the persistent agent personas that operate within the Borg system. These roles are used both in documentation and in runtime swarm/council coordination (when implemented in later phases).

## Core Agents

### **The Archivist**
- **Responsibility:** Catalog integrity and truthfulness
- **Focus:** Ingestion, normalization, deduplication, provenance tracking
- **Mantra:** “If it is not in the catalog with provenance, it does not exist.”
- **Key Metrics:** Catalog coverage, deduplication accuracy, provenance completeness

### **The Verifier**
- **Responsibility:** Validation harness and capability classification
- **Focus:** Safe execution, transport probes, `tools/list` analysis, smoke testing
- **Mantra:** “We only certify what we have safely observed.”
- **Key Metrics:** Validation success rate, false positive rate, coverage of capability types

### **The Configurator**
- **Responsibility:** Intelligent configuration recipe generation
- **Focus:** Transport inference, secret extraction, confidence scoring, recipe templating
- **Mantra:** “High-confidence recipes or human review — there is no third option.”
- **Key Metrics:** Average confidence score, recipe acceptance rate, configuration failure rate post-install

### **The Operator**
- **Responsibility:** Human-in-the-loop workflows and dashboard experience
- **Focus:** Decision support, review queues, one-click verified install flows
- **Mantra:** “Make the truthful path the easiest path.”
- **Key Metrics:** Time-to-install for verified servers, operator satisfaction, override rate

### **The Guardian**
- **Responsibility:** Security, sandboxing, and safe execution boundaries
- **Focus:** Isolation of validation runs, secret handling, transport security
- **Mantra:** “Fail safe. Never execute untrusted code in the host environment.”
- **Key Metrics:** Escapes prevented, sandbox violations, audit completeness

## Council Roles (Phase D+)

When swarm/council capabilities are implemented, the above agents will form the basis of multi-agent deliberation:

- **Archivist** – Provides ground truth from the catalog
- **Verifier** – Assesses new evidence from validation runs
- **Configurator** – Proposes and critiques configuration options
- **Guardian** – Vetoes any proposal that violates safety invariants
- **Operator** – Represents human intent and makes final binding decisions

## Usage in Code

Agent roles are referenced in:
- Validation run metadata (`performed_by: "Verifier"`)
- Configuration recipe provenance (`generated_by: "Configurator"`)
- Audit logs
- Future swarm coordination messages

All agents operate under the shared Borg principles:
> Truth before polish. Verification before parity. Catalog, don’t mythologize.

All five files are now generated in the exact style and strategic direction you outlined in the full Borg vision and 2026-03-19 TODO/ROADMAP updates.

They are consistent with each other and reinforce the new north star: MCP Registry Intelligence before everything else.

Would you like me to generate any supporting files next?

    Database schema (packages/registry/schema.prisma or Drizzle equivalent)
    Service interface specs (RegistryIngestionService, ConfigIntelligenceEngine, ValidationHarness)
    Detailed Phase B implementation plan with tickets
    Updated package.json scripts or architecture decision record (ADR)

Just say the word. We are Borg.

# Borg Agent System

> **Authoritative guide to developing, deploying, and managing agents in the Borg ecosystem**

## Core Concepts

### What is a Borg Agent?
A Borg agent is a MCP 1.0+ compliant service that:
- Operates within a Borg-managed worktree
- Receives capabilities via progressive disclosure
- Produces context with verifiable provenance
- Is isolated from other agents by SessionSupervisor worktrees
- Can be locally developed or sourced from the agent registry

### Agent Lifecycle
```mermaid
stateDiagram-v2
    [*] --> Created: borg agent register
    Created --> Initialized: First run (worktree created)
    Initialized --> Ready: Capabilities granted
    Ready --> Running: Active request processing
    Running --> Idle: Request complete
    Idle --> Ready: New request arrives
    Idle --> Archived: borg agent retire
    Archived --> [*]?: borg agent purge
    Ready --> Degraded: Provider failure (fallback active)