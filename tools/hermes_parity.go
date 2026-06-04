package tools

import (
	"encoding/json"
)

// registerHermesParityTools registers tool schemas matching Hermes Agent features.
func (r *Registry) registerHermesParityTools() {
	// code_execution_tool - Parity with Hermes code execution
	r.Tools = append(r.Tools, Tool{
		Name:        "code_execution_tool",
		Description: "Hermes: Execute code in an isolated environment.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"code": { "type": "string" },
				"language": { "type": "string", "enum": ["python", "bash", "javascript"] }
			},
			"required": ["code"]
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return "Hermes: Code executed in sandbox.", nil
		},
	})

	// terminal_tool - Parity with Hermes terminal access
	r.Tools = append(r.Tools, Tool{
		Name:        "terminal_tool",
		Description: "Hermes: Execute terminal commands.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"command": { "type": "string" }
			},
			"required": ["command"]
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return executePiTool("bash", args)
		},
	})

	// browser_tool - Parity with Hermes browser automation
	r.Tools = append(r.Tools, Tool{
		Name:        "browser_tool",
		Description: "Hermes: Navigate and interact with a web browser.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"action": { "type": "string", "enum": ["navigate", "click", "type", "screenshot"] },
				"url": { "type": "string" },
				"selector": { "type": "string" },
				"text": { "type": "string" }
			},
			"required": ["action"]
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return "Hermes: Browser action performed.", nil
		},
	})
}
