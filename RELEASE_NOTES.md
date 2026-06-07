# Release Notes - HyperHarness v0.4.5

This release marks the completion of the Go-native multi-harness porting phase. HyperHarness now achieves 1:1 tool parity with over 15 major AI coding harnesses.

## Key Features
- **Multi-Harness Parity**: Full tool signature and functional parity for Claude (Code/Desktop), Codex (Desktop/CLI), Hermes (Agent/Desktop), Tabby, Warp, Wave, Antigravity, and pi-coding-agent.
- **HyperSync (hsync)**: Integrated background synchronization for BobbyBookmarks and LinkCrawler content extraction.
- **Unified LLM Routing (ai)**: Core module for provider-agnostic model routing and prompt template management.
- **Data Ingestion**: High-quality data normalization and AST-aware Go code summarization for knowledge base seeding.
- **TUI Dashboard**: Advanced BubbleTea-based REPL and real-time system monitoring.

## Ported Harnesses
1. Tabby AI
2. Warp Terminal
3. Wave Terminal
4. Hermes Agent
5. Antigravity 2.0
6. Codex Desktop
7. Claude Code
8. Claude Desktop
9. Codex CLI
10. pi-coding-agent
11. Gemini CLI
12. OpenCode
13. Goose
14. Kimi
15. Aider

## Documentation
- Architectural analyses for each harness can be found in `docs/analysis/`.
- Usage instructions are available in `README.md` and `VISION.md`.
