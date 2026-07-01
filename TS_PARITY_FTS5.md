# Memory Isolation and FTS5 TypeScript Parity

The roadmap states:
- [x] Synchronize Memory Isolation and FTS5 to TypeScript client

We need to implement a SQLite FTS5 store in `pi-cli/packages/coding-agent/src/core/memory-store.ts` that provides the same API as the Go implementation in `internal/memory/sqlite_store.go`, specifically:
- `Store(entry MemoryEntry)`
- `Search(query string, limit int)`
- `SearchScoped(query, scope string, limit int)`
