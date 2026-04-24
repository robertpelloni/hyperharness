# SuperAI Agent Core Assimilation Tranche

Date: 2026-04-05

## Summary

This tranche continues the `../superai` → `hyperharness` assimilation effort at the central agent runtime seam:
- `agent/agent.go`

This is an especially important seam because it sits between:
- the OpenAI-compatible chat/tool-call loop
- the Go-native tool registry
- the Pi-compatible foundation tool layer
- the broader Borg/HyperCode adapter context

## Key finding

`agent/agent.go` was already present in `hyperharness` and effectively parity-aligned with `../superai`.

The highest-value work here was therefore not redundant copying. It was to:
- preserve the public behavior surface
- extract implicit assumptions into explicit helpers
- add validation for missing core dependencies
- add regression coverage around tool-call and response-choice handling

## What changed

### 1. System prompt extraction
Added:
- `buildAgentSystemPrompt(...)`

This makes the core system-prompt contract explicit and testable.

### 2. Response-choice validation
Added:
- `firstChoiceMessage(...)`

This isolates the previously implicit assumption that a completion response always contains at least one choice.

Behavioral improvement:
- both `Chat(...)` and `handleToolCalls(...)` now fail clearly if no completion choices are returned

### 3. Tool-call execution extraction
Added:
- `executeToolCall(...)`

This isolates tool-call resolution against the Go-native registry and makes the behavior directly testable.

Behavior preserved:
- unknown tools still report `Unknown tool: <name>`
- if a tool returns both output and an error, the output still wins
- otherwise the formatted execution error string is returned

### 4. Core validation guards
Added explicit validation for:
- nil agent receiver
- missing OpenAI client
- missing tool registry

This removes obvious crash paths in the central agent loop.

## Regression coverage added

Added:
- `agent/core_assimilation_test.go`

Coverage includes:
- system prompt construction
- empty-choice validation
- tool-call unknown/error handling
- nil-agent validation
- missing client validation
- missing registry validation

## Validation

Verified successfully:

```bash
go test ./agent ./agents ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because `agent/agent.go` is the core bridge between:
- model-facing tool calling
- exact-name tool registry exposure
- Go-native runtime behavior

By making the assumptions here explicit and tested, `hyperharness` becomes a stronger canonical base for further exact-contract assimilation.

## Design insight

A useful lesson from this tranche is:

> the highest-leverage hardening often happens at the central seams where multiple already-good subsystems meet.

This file is one of those seams.

## Recommended next move

The strongest next assimilation step is now:

1. continue through any remaining central `agent/` files not yet hardened, especially:
   - `agent/pipe.go`
2. then widen into adjacent integration seams only where they add real parity/confidence
3. keep validating with:
   - `go test ./agent ./agents ./tui ./cmd ./foundation/...`

## Bottom line

This tranche confirms that `hyperharness` already contained the `superai` core agent loop, and upgrades it into a clearer, safer, regression-tested base for the broader Go-native assimilation effort.
