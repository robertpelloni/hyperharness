package tools

import (
	"context"
	"encoding/json"

	foundationpi "github.com/robertpelloni/hyperharness/foundation/pi"
)

// registerComputerUseTools exposes exact parity with Anthropic's Computer Use / Open-Interpreter standard.
func (r *Registry) registerComputerUseTools() {
	cwd := "."

	// Computer / bash
	r.Tools = append(r.Tools, Tool{
		Name:        "bash",
		Description: "Run bash commands. Anthropic computer-use parity.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["command"],
			"properties": {
				"command": { "type": "string", "description": "The command to execute" },
				"timeout": { "type": "number", "description": "Optional timeout" }
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			command, _ := args["command"].(string)
			piArgs := map[string]interface{}{"command": command}
			raw, _ := json.Marshal(piArgs)
			runtime := foundationpi.NewRuntime(cwd, nil)
			result, err := runtime.ExecuteTool(context.Background(), "", "bash", raw, nil)
			if err != nil {
				return "", err
			}
			return formatFoundationToolResult(result), nil
		},
	})

	// str_replace
	r.Tools = append(r.Tools, Tool{
		Name:        "str_replace",
		Description: "Replace string in file. Anthropic computer-use parity.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["path", "old_str", "new_str"],
			"properties": {
				"path": { "type": "string" },
				"old_str": { "type": "string" },
				"new_str": { "type": "string" }
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			path, _ := args["path"].(string)
			oldStr, _ := args["old_str"].(string)
			newStr, _ := args["new_str"].(string)

			piArgs := map[string]interface{}{
				"path": path,
				"edits": []map[string]interface{}{
					{"oldText": oldStr, "newText": newStr},
				},
			}
			raw, _ := json.Marshal(piArgs)
			runtime := foundationpi.NewRuntime(cwd, nil)
			result, err := runtime.ExecuteTool(context.Background(), "", "edit", raw, nil)
			if err != nil {
				return "", err
			}
			return formatFoundationToolResult(result), nil
		},
	})
}
