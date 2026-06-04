# Wave (Waveterm) AI Features Analysis

## Overview
Wave is an open-source, AI-integrated terminal. It features "Wave AI", a context-aware assistant that has access to the terminal's workspace, blocks, and scrollback.

## AI Features & Tools
- **Context-Awareness**: The assistant sees the "Current Tab State", including all open widgets (terminals, previews, web browsers, etc.).
- **Tools**:
  - `CaptureScreenshot`: Take a screenshot of a specific tab.
  - `ReadTextFile` / `WriteTextFile` / `EditTextFile` / `DeleteTextFile`: Basic file operations.
  - `ReadDir`: List directory contents.
  - `TermGetScrollback`: Retrieve the scrollback buffer of a terminal block.
  - `WebNavigate`: Navigate the web browser widget.
  - `TsunamiGetData` / `TsunamiGetConfig`: Interact with Tsunami blocks (likely a specific Wave widget type).

## Implementation Details in HyperHarness
- **Parity Tools**: Implement Wave's tools in `tools/wave_parity.go`.
- **Workspace Integration**: Wave's concept of "Current Tab State" matches HyperHarness's `internal/workspaces`. We should implement a `GetTabState` function in `internal/workspaces/tracker.go` that aggregates information about active sessions and their metadata.
- **Scrollback Capture**: Implement `TermGetScrollback` by reading from the session logs in `internal/sessions`.
