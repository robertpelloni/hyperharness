# HyperHarness Performance Benchmarks

This document records the baseline metrics for core subsystems in Phase 5 to validate production readiness and identify optimization targets before final release.

## Context Compaction
Measures the overhead of iterating history messages during the `Compact` operation for different lengths of history sizes.

| Test Case | Time / op | B / op | Allocs / op |
|---|---|---|---|
| Size=10 | ~278 ns/op | 0 | 0 |
| Size=100 | ~311 ns/op | 0 | 0 |
| Size=1000 | ~489 ns/op | 0 | 0 |

*Note: History array modifications execute exceptionally fast due to pass-by-reference design.*

## Memory (SQLite FTS5)
Measures the performance of retrieval from the knowledge base holding ~1,000 document records.

| Test Case | Time / op | B / op | Allocs / op |
|---|---|---|---|
| Search (Global limit 10) | ~5.4 ms/op | 20,916 | 426 |
| SearchScoped (Limit 10) | ~3.7 ms/op | 20,996 | 428 |

*Note: Scoped search adds efficiency due to the filtering at the query plan level.*

## Subagent Spawn Time
Measures the time to spawn subagents and allocate new task constraints in the `internal/subagents.Manager.CreateTask()` structure.

| Test Case | Time / op | B / op | Allocs / op |
|---|---|---|---|
| Spawning SubagentType | ~2.7 µs/op | 406 | 4 |

## Tool Dispatch Overhead
Measures the execution latency overhead comparing directly calling a foundational native tool versus calling a parity delegated wrapper that verifies schema and routes.

| Test Case | Time / op | B / op | Allocs / op |
|---|---|---|---|
| Foundation Tool (\`read\`) | ~7.4 µs/op | 772 | 19 |
| Parity Tool (\`read_file\`) | ~2.0 µs/op | 72 | 2 |

*Note: Parity tools routing appears highly efficient. Foundation metrics inherently measure disk I/O in the test mock.*

## Summary

The core engine is comfortably sub-millisecond, adhering to production goals:
- **Tool Dispatch:** O(µs) per call
- **Subagent Routing:** O(µs) allocation overhead
- **Knowledge Retrieval:** O(ms) SQLite latency
- **Context:** O(ns) slice slicing operations
