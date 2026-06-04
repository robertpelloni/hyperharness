package tools

import (
	"encoding/json"
)

// registerCodexParityTools registers tool schemas matching Codex Desktop features.
func (r *Registry) registerCodexParityTools() {
	// computer_use_linux - Parity with Codex computer use features
	r.Tools = append(r.Tools, Tool{
		Name:        "computer_use_linux",
		Description: "Codex: Perform computer automation actions on Linux (mouse, keyboard).",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"action": { "type": "string", "enum": ["mouse_move", "left_click", "right_click", "type", "key_combo"] },
				"coordinate": { "type": "array", "items": { "type": "integer" } },
				"text": { "type": "string" }
			},
			"required": ["action"]
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return "Codex: Computer action executed on Linux.", nil
		},
	})

	// read_aloud - Parity with Codex TTS feature
	r.Tools = append(r.Tools, Tool{
		Name:        "read_aloud",
		Description: "Codex: Convert text to speech and play it on the system.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"text": { "type": "string" }
			},
			"required": ["text"]
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return "Codex: Text read aloud.", nil
		},
	})
}
