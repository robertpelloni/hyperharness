# Hermes Agent Analysis

## Overview
Hermes is a "self-improving" AI agent by Nous Research. It features a wide array of tools and a skill-based architecture.

## Key Features & Tools
- **Code Execution**: `code_execution_tool` for running code in isolated environments.
- **Terminal**: `terminal_tool` for shell access.
- **File Operations**: `file_tools`, `file_operations`, `patch_parser`.
- **Browser**: `browser_tool`, `browser_camofox`, `browser_cdp_tool`.
- **Communication**: `discord_tool`, `send_message_tool`, `feishu_doc_tool`.
- **Skills**: Extensive skill library in `skills/` (smart-home, social-media, software-development, etc.).
- **Others**: `image_generation_tool`, `video_generation_tool`, `kanban_tools`, `todo_tool`, `mcp_tool`.

## Implementation Details in HyperHarness
- **Parity Tools**: Implement Hermes-specific tools in `tools/hermes_parity.go`.
- **Skill System**: Hermes's skill library is a great source for HyperHarness's skill system. We should analyze the structure of Hermes skills and ensure HyperHarness can load or adapt them.
- **Isolated Execution**: Hermes emphasizes isolated environments. HyperHarness should leverage its `internal/codeexec` or `internal/docker_sandbox` for similar functionality.
