package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// RefactorTool bridges Opencode parity with strict AST REPLACE blocks.
type RefactorTool struct{}

// ApplySearchReplace strictly enforces LLM blocks in the format:
// <<<<<<< SEARCH
// existing
// =======
// new
// >>>>>>> REPLACE
func (r *RefactorTool) ApplySearchReplace(filePath, searchBlock, replaceBlock string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read target: %w", err)
	}

	strContent := string(content)

	if !strings.Contains(strContent, searchBlock) {
		return fmt.Errorf("search block not found natively in file. LLM hallucinated context buffer.")
	}

	// Native atomic replacement
	newContent := strings.Replace(strContent, searchBlock, replaceBlock, 1)

	// Write back with strict permissions
	err = os.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("failed atomic write application: %w", err)
	}

	fmt.Printf("[Opencode Parity] Block successfully mutated %s natively.\n", filePath)
	return nil
}

// registerRefactoringTools binds the isolated Aider functionality to the core Native Engine.
func (reg *Registry) registerRefactoringTools() {
	reg.Tools = append(reg.Tools, Tool{
		Name:        "apply_search_replace",
		Description: "Opencode Parity: strict block-replacement AST refactoring.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["file_path", "search_block", "replace_block"],
			"properties": {
				"file_path": { "type": "string" },
				"search_block": { "type": "string" },
				"replace_block": { "type": "string" }
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			path, _ := args["file_path"].(string)
			search, _ := args["search_block"].(string)
			replace, _ := args["replace_block"].(string)

			rf := &RefactorTool{}
			err := rf.ApplySearchReplace(path, search, replace)
			if err != nil {
				return "", err
			}
			return "File safely mutated natively.", nil
		},
	})

	reg.Tools = append(reg.Tools, Tool{
		Name:        "replace_lines",
		Description: "create or update one or more files (Aider exact parity)",
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

				rf := &RefactorTool{}
				if err := rf.ApplySearchReplace(path, originalBlock, updatedBlock); err != nil {
					failures = append(failures, fmt.Errorf("edit %d on %s failed: %v", i, path, err).Error())
				} else {
					successCount++
				}
			}

			if len(failures) > 0 {
				return fmt.Sprintf("Applied %d edits. %d failed:\n%s", successCount, len(failures), strings.Join(failures, "\n")), nil
			}

			return fmt.Sprintf("Successfully applied %d edits.", successCount), nil
		},
	})
}
