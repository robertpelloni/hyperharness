---
name: VS Code Workspace Bridge
description: MCP server exposing VS Code's filesystem, editing, and symbol search capabilities to AI models.
---

# VS Code Workspace Bridge

The VS Code MCP Server bridges your IDE workspace directly to AI models, allowing them to exploration the codebase, search for symbols, and make precise edits using VS Code's internal APIs.

## Core Features

- **Symbol Search:** Find functions, classes, and variables across the whole workspace.
- **Diagnostics:** Get real-time errors and warnings from the VS Code Problems panel.
- **Precise Edits:** Use `replace_lines_code` for safe, validated code modifications.
- **Context Efficiency:** Use `get_document_symbols_code` to get file outlines instead of reading full files.

## Supported Tools

### Exploration
- `list_files_code`: List files/directories.
- `search_symbols_code`: Search for names across project.
- `get_document_symbols_code`: Get hierarchical file outline.

### Analysis
- `read_file_code`: Read content with size limits.
- `get_symbol_definition_code`: Get type info and documentation.
- `get_diagnostics_code`: Check for workspace errors.

### Modification
- `create_file_code`: Create or overwrite files.
- `replace_lines_code`: Replace exact line ranges.
- `move_file_code` / `rename_file_code`: Refactoring-aware file moves.

## Best Practices
1. **Outline First:** Always run `get_document_symbols_code` before reading a large file.
2. **Verify Changes:** Run `get_diagnostics_code` after every edit to ensure no new errors were introduced.
3. **Plan Context:** Read your `CLAUDE.md` or `reference/` docs for project-specific patterns.
