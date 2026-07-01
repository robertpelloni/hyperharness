# TS Parity Tools Sync (Goose, OpenCode, Kimi, Cursor, Windsurf)

To complete the syncing of the remaining parity tool sets:
1. Ensure `scripts/extract_schemas.py` target directories include `goose`, `opencode`, `kimi-cli`, etc.
2. Re-run `python3 scripts/extract_schemas.py` and `python3 scripts/generate_tools.py` to regenerate the TS definitions inside `pi-cli/packages/coding-agent/src/core/tools/parity/generated/index.ts`.
3. Create end-to-end parity validation tests in `pi-cli/packages/coding-agent/test/` to confirm consistent behavior between Go core and TS client.
4. Implement daemon hardening, crash recovery, auto-restart logic, and observability hooks.
