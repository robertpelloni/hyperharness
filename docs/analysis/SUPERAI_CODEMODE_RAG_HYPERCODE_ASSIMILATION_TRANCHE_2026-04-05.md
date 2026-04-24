# SuperAI CodeMode, RAG, and HyperCode Provider Assimilation Tranche

Date: 2026-04-05

## Summary

This tranche continues the `../superai` → `hyperharness` assimilation effort with the next richer `agents/` surfaces:
- `agents/codemode.go`
- `agents/rag.go`
- `agents/provider_hypercode.go`

## Key finding

These files were already present in `hyperharness` and effectively parity-aligned with `../superai`.

The highest-value work here was therefore:
- preserve the intended public surface
- extract testable helpers
- harden contract edges
- add regression coverage
- keep the entire validated package target green

## What changed

### 1. `agents/codemode.go` hardening
Added helper extraction and normalization:
- `resolveCodeMode(...)`
- `codeModeExecution` helper struct
- explicit constants for generated filenames:
  - `codeModeTypeScriptFile`
  - `codeModeShellFile`

Behavioral improvements:
- language selection is now normalized with trim + lowercase
- alias handling is explicit and testable
- cleanup is handled via `defer`
- unsupported language behavior remains explicit and stable

This keeps the same intended execution bindings while making the contract more explicit.

### 2. `agents/rag.go` hardening
Added explicit constants and validation:
- `defaultRAGVectorStore`
- `defaultEmbeddingVector`

Behavioral improvements:
- `DocumentIntakeService.Ingest(...)` now rejects nil receiver and empty filepath
- `EmbeddingService.Compute(...)` now rejects empty text
- embedding results are returned as a copied slice rather than aliasing a shared backing array

This upgrades the RAG stubs from silent placeholders into explicit, safer contracts.

### 3. `agents/provider_hypercode.go` hardening
Added helper extraction:
- `buildHyperCodeChatRequest(...)`
- `parseHyperCodeChatResponse(...)`
- `normalizeHyperCodeBaseURL(...)`
- explicit request/response helper types

Behavioral improvements:
- nil provider receiver validation in `Chat`
- nil chunk channel validation in `Stream`
- default base URL fallback when unset
- trailing slash normalization for base URL
- response parsing separated from transport logic
- HTTP error, parse error, and rejection paths are now directly testable

This is particularly important because this provider is a bridge into the broader HyperCode control-plane world you want the harness to interoperate with.

## Regression coverage added

Added:
- `agents/codemode_rag_hypercode_assimilation_test.go`

Coverage includes:
- CodeMode alias resolution and unsupported-language behavior
- RAG vector-store default and embedding/input validation
- HyperCode request building and response parsing helpers
- HyperCode provider successful chat/stream behavior via test server
- HyperCode provider validation and failure paths:
  - nil provider
  - nil chunk channel
  - HTTP error
  - invalid JSON
  - explicit API rejection
  - unreachable transport

## Validation

Verified successfully:

```bash
go test ./agent ./agents ./tui ./cmd ./foundation/...
```

Everything remained green.

## Why this matters

This tranche matters because it strengthens three important bridge surfaces at once:
- code execution escape hatch behavior
- local RAG/document intake contracts
- HyperCode control-plane interoperability

These are all likely to matter in the broader unified harness vision.

## Design insight

A useful lesson from this tranche is:

> bridge layers become dramatically more maintainable once request-building, response-parsing, and normalization logic are extracted into isolated helpers.

That is especially true when the long-term target is exact-name, exact-contract tool and provider compatibility across multiple runtimes.

## Recommended next move

The strongest next assimilation step is now:

1. continue through the remaining richer `agents/` or provider-adjacent surfaces
2. inspect next files such as:
   - `agents/provider.go`
   - `agents/provider_stub.go`
   - `agents/codemode.go`'s surrounding runtime call sites if deeper parity hooks are needed
3. keep validating with:
   - `go test ./agent ./agents ./tui ./cmd ./foundation/...`

## Bottom line

This tranche confirms that `hyperharness` already contained the `superai` CodeMode, RAG, and HyperCode-provider surface, and upgrades it into a clearer, safer, regression-tested base for continued Go-native assimilation.
