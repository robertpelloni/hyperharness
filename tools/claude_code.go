package tools

import (
	"context"
	"encoding/json"

	foundationpi "github.com/robertpelloni/hyperharness/foundation/pi"
)

// registerClaudeCodeTools exposes exact parity with Claude Code's native tools.
// Models like Claude 3.5/3.7 have been fine-tuned to use these exact tools natively.
func (r *Registry) registerClaudeCodeTools() {
	cwd := "." // For simplicity in this alias binding

	// Edit Tool
	r.Tools = append(r.Tools, Tool{
		Name:        "Edit",
		Description: "A tool for editing files. Use this to modify file contents in place.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["file_path", "old_string", "new_string"],
			"properties": {
				"file_path": { "type": "string", "description": "The absolute path to the file to modify" },
				"old_string": { "type": "string", "description": "The text to replace" },
				"new_string": { "type": "string", "description": "The text to replace it with (must be different from old_string)" },
				"replace_all": { "type": "boolean", "description": "Replace all occurrences of old_string (default false)" }
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			path, _ := args["file_path"].(string)
			oldStr, _ := args["old_string"].(string)
			newStr, _ := args["new_string"].(string)

			// Map to foundation pi edit tool
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

	// Read Tool
	r.Tools = append(r.Tools, Tool{
		Name:        "Read",
		Description: "A tool for reading files.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["file_path"],
			"properties": {
				"file_path": { "type": "string", "description": "The absolute path to the file to read" }
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			path, _ := args["file_path"].(string)
			piArgs := map[string]interface{}{"path": path}
			raw, _ := json.Marshal(piArgs)
			runtime := foundationpi.NewRuntime(cwd, nil)
			result, err := runtime.ExecuteTool(context.Background(), "", "read", raw, nil)
			if err != nil {
				return "", err
			}
			return formatFoundationToolResult(result), nil
		},
	})

	// Write Tool
	r.Tools = append(r.Tools, Tool{
		Name:        "Write",
		Description: "A tool for writing files.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["file_path", "content"],
			"properties": {
				"file_path": { "type": "string", "description": "The absolute path to the file to write" },
				"content": { "type": "string", "description": "The content to write" }
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			path, _ := args["file_path"].(string)
			content, _ := args["content"].(string)
			piArgs := map[string]interface{}{"path": path, "content": content}
			raw, _ := json.Marshal(piArgs)
			runtime := foundationpi.NewRuntime(cwd, nil)
			result, err := runtime.ExecuteTool(context.Background(), "", "write", raw, nil)
			if err != nil {
				return "", err
			}
			return formatFoundationToolResult(result), nil
		},
	})

	// Bash Tool
	r.Tools = append(r.Tools, Tool{
		Name:        "Bash",
		Description: "A tool for executing bash commands.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["command"],
			"properties": {
				"command": { "type": "string", "description": "The command to execute" }
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

	// Glob Tool
	r.Tools = append(r.Tools, Tool{
		Name:        "Glob",
		Description: "A tool for globbing files.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["pattern"],
			"properties": {
				"pattern": { "type": "string", "description": "The glob pattern to execute" },
				"path": { "type": "string", "description": "The path to execute the pattern in" }
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			pattern, _ := args["pattern"].(string)
			path, _ := args["path"].(string)
			piArgs := map[string]interface{}{"pattern": pattern, "path": path}
			raw, _ := json.Marshal(piArgs)
			runtime := foundationpi.NewRuntime(cwd, nil)
			result, err := runtime.ExecuteTool(context.Background(), "", "find", raw, nil)
			if err != nil {
				return "", err
			}
			return formatFoundationToolResult(result), nil
		},
	})

	// Grep Tool
	r.Tools = append(r.Tools, Tool{
		Name:        "Grep",
		Description: "A tool for grepping files.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["pattern"],
			"properties": {
				"pattern": { "type": "string", "description": "The grep pattern to execute" },
				"path": { "type": "string", "description": "The path to execute the pattern in" }
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			pattern, _ := args["pattern"].(string)
			path, _ := args["path"].(string)
			piArgs := map[string]interface{}{"pattern": pattern, "path": path}
			raw, _ := json.Marshal(piArgs)
			runtime := foundationpi.NewRuntime(cwd, nil)
			result, err := runtime.ExecuteTool(context.Background(), "", "grep", raw, nil)
			if err != nil {
				return "", err
			}
			return formatFoundationToolResult(result), nil
		},
	})

	// LS Tool
	r.Tools = append(r.Tools, Tool{
		Name:        "LS",
		Description: "A tool for listing directories.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"path": { "type": "string", "description": "The path to list" }
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			path, _ := args["path"].(string)
			piArgs := map[string]interface{}{"path": path}
			raw, _ := json.Marshal(piArgs)
			runtime := foundationpi.NewRuntime(cwd, nil)
			result, err := runtime.ExecuteTool(context.Background(), "", "ls", raw, nil)
			if err != nil {
				return "", err
			}
			return formatFoundationToolResult(result), nil
		},
	})
}
