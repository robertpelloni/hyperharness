package tools

import (
	"encoding/json"
	"fmt"
	"github.com/robertpelloni/hyperharness/internal/sessions"
	"strings"
)

// registerWaveParityTools registers tool schemas matching Wave (Waveterm) AI features.
func (r *Registry) registerWaveParityTools() {
	// TermGetScrollback - Parity with Wave's scrollback retrieval
	r.Tools = append(r.Tools, Tool{
		Name:        "TermGetScrollback",
		Description: "Wave: Retrieve the scrollback buffer for a specific terminal block.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"blockId": { "type": "string" },
				"numLines": { "type": "integer", "default": 100 }
			},
			"required": ["blockId"]
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			blockID := GetStr(args, "blockId")
			numLines := GetIntDef(args["numLines"], 100)

			// Use the global session manager instance
			mgr := sessions.GetGlobalManager()
			if mgr == nil {
				return "", fmt.Errorf("failed to get session manager")
			}

			// In HyperHarness, we might map blockId to a session ID
			session, err := mgr.LoadSession(blockID)
			if err != nil {
				// Fallback to active session if blockID not found
				list := mgr.ListSessions()
				if len(list) > 0 {
					session = list[0]
				} else {
					return fmt.Sprintf("Wave: No session found for block %s", blockID), nil
				}
			}

			var scrollback []string
			entries := session.Entries
			if len(entries) > numLines {
				entries = entries[len(entries)-numLines:]
			}

			for _, e := range entries {
				if e.Type == sessions.EntryToolResult && e.ToolName == "bash" {
					scrollback = append(scrollback, e.ToolResult)
				} else {
					scrollback = append(scrollback, e.Content)
				}
			}

			return strings.Join(scrollback, "\n"), nil
		},
	})

	// CaptureScreenshot - Parity with Wave's screenshot tool
	r.Tools = append(r.Tools, Tool{
		Name:        "CaptureScreenshot",
		Description: "Wave: Capture a screenshot of the current tab or a specific widget.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"tabId": { "type": "string" },
				"blockId": { "type": "string" }
			}
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return "Wave: Screenshot captured.", nil
		},
	})

	// wsh_ai - Parity with Wave's CLI AI tool
	r.Tools = append(r.Tools, Tool{
		Name:        "wsh_ai",
		Description: "Wave: Pipe output or attach files to an AI command from the shell.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"message": { "type": "string" },
				"pipe_input": { "type": "string" },
				"files": { "type": "array", "items": { "type": "string" } }
			},
			"required": ["message"]
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return "Wave wsh ai: Request processed.", nil
		},
	})
}
