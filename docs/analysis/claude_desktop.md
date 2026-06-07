# Claude Desktop Analysis

Claude Desktop is primarily a wrapper for the Claude AI experience on desktop, but it integrates with **Claude Code** through handoff and deep linking.

## Features

### 1. Session Handoff
- **Deep Linking**: Claude Code can open the current session in the Claude Desktop app using `claude://` URIs.
- **Context Preservation**: Transferring conversation state between the CLI and the Desktop app.

### 2. MCP Client Integration
- Claude Desktop acts as a host for MCP servers.
- **Config**: Reads from `claude_desktop_config.json` (usually in `~/Library/Application Support/Claude/` or `%APPDATA%\Claude\`).

### 3. Tool Parity
Claude Desktop exposes tools through MCP servers it hosts. In HyperHarness, parity means:
- **Auto-Discovery**: Scanning for Claude Desktop's MCP configuration and automatically registering those servers in the Go registry.
- **Handoff Tool**: A tool to generate or trigger the deep link for desktop transition.

## Implementation in Go

- [x] MCP Config Auto-Discovery (`internal/mcp/discovery.go`)
- [x] Deep Link Generator tool (`claude_desktop_handoff`)
