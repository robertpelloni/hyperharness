package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	foundationpi "github.com/robertpelloni/hyperharness/foundation/pi"
)

// registerAiderTools exposes exact parity with Aider's native tools.
func (r *Registry) registerAiderTools() {
	cwd := "."

	// replace_lines (from editblock_func_coder.py)
	r.Tools = append(r.Tools, Tool{
		Name:        "replace_lines",
		Description: "create or update one or more files",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["explanation", "edits"],
			"properties": {
				"explanation": {
					"type": "string",
					"description": "Step by step plan for the changes to be made to the code (future tense, markdown format)"
				},
				"edits": {
					"type": "array",
					"items": {
						"type": "object",
						"required": ["path", "original_lines", "updated_lines"],
						"properties": {
							"path": {
								"type": "string",
								"description": "Path of file to edit"
							},
							"original_lines": {
								"type": "array",
								"items": { "type": "string" },
								"description": "A unique stretch of lines from the original file, including all whitespace, without skipping any lines"
							},
							"updated_lines": {
								"type": "array",
								"items": { "type": "string" },
								"description": "New content to replace the original_lines with"
							}
						}
					}
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			editsRaw, ok := args["edits"].([]interface{})
			if !ok {
				return "", fmt.Errorf("edits must be an array")
			}

			successCount := 0
			var failures []string

			runtime := foundationpi.NewRuntime(cwd, nil)

			for i, eRaw := range editsRaw {
				e, ok := eRaw.(map[string]interface{})
				if !ok {
					failures = append(failures, fmt.Errorf("edit %d is invalid", i).Error())
					continue
				}

				path, _ := e["path"].(string)

				var originalLines []string
				if origRaw, ok := e["original_lines"].([]interface{}); ok {
					for _, o := range origRaw {
						if s, ok := o.(string); ok {
							originalLines = append(originalLines, s)
						}
					}
				}

				var updatedLines []string
				if updRaw, ok := e["updated_lines"].([]interface{}); ok {
					for _, u := range updRaw {
						if s, ok := u.(string); ok {
							updatedLines = append(updatedLines, s)
						}
					}
				}

				originalBlock := strings.Join(originalLines, "\n")
				updatedBlock := strings.Join(updatedLines, "\n")

				piArgs := map[string]interface{}{
					"path": path,
					"edits": []map[string]interface{}{
						{"oldText": originalBlock, "newText": updatedBlock},
					},
				}
				raw, _ := json.Marshal(piArgs)
				result, err := runtime.ExecuteTool(context.Background(), "", "edit", raw, nil)

				if err != nil {
					failures = append(failures, fmt.Errorf("edit %d on %s failed: %v", i, path, err).Error())
				} else {
					// Check if result has ErrorOutput
					if result.IsError {
						failures = append(failures, fmt.Errorf("edit %d on %s failed: %v", i, path, formatFoundationToolResult(result)).Error())
					} else {
						successCount++
					}
				}
			}

			if len(failures) > 0 {
				return fmt.Sprintf("Applied %d edits. %d failed:\n%s", successCount, len(failures), strings.Join(failures, "\n")), nil
			}

			return fmt.Sprintf("Successfully applied %d edits.", successCount), nil
		},
	})

	// write_file (from single_wholefile_func_coder.py)
	r.Tools = append(r.Tools, Tool{
		Name:        "write_file",
		Description: "write new content into the file",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"explanation": {
					"type": "string",
					"description": "Step by step plan for the changes to be made to the code (future tense, markdown format)"
				},
				"content": {
					"type": "string",
					"description": "Content to write to the file"
				},
				"path": {
					"type": "string",
					"description": "Path to write to (added by hyperharness for multi-file support)"
				}
			},
			"required": ["explanation", "content", "path"],
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			path, ok := args["path"].(string)
			if !ok {
				return "", fmt.Errorf("path is required")
			}
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
}
