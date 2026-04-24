package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	foundationpi "github.com/robertpelloni/hyperharness/foundation/pi"
)

// gemini_parity.go provides exact Gemini CLI tool surfaces.
// Models trained on Gemini CLI expect these exact tool names and parameters.

// registerGeminiCLITools adds Gemini CLI compatible tool surfaces.
func (r *Registry) registerGeminiCLITools() {
	// read_file - Gemini exact naming
	r.Tools = append(r.Tools, Tool{
		Name:        "read_file",
		Description: "Read the contents of a file.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["path"],
			"properties": {
				"path": {
					"type": "string",
					"description": "The path to the file to read"
				},
				"offset": {
					"type": "integer",
					"description": "Line number to start reading from (1-indexed)"
				},
				"limit": {
					"type": "integer",
					"description": "Maximum number of lines to read"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			path, _ := args["path"].(string)
			if path == "" {
				return "", fmt.Errorf("path is required")
			}
			piArgs := map[string]interface{}{
				"path": path,
			}
			if offset, ok := args["offset"]; ok {
				piArgs["offset"] = offset
			}
			if limit, ok := args["limit"]; ok {
				piArgs["limit"] = limit
			}
			raw, _ := json.Marshal(piArgs)
			runtime := foundationpi.NewRuntime(".", nil)
			result, err := runtime.ExecuteTool(context.Background(), "", "read", raw, nil)
			if err != nil {
				return "", err
			}
			return formatFoundationToolResult(result), nil
		},
	})

	// write_file - Gemini exact naming
	r.Tools = append(r.Tools, Tool{
		Name:        "write_file",
		Description: "Write content to a file. Creates the file if it doesn't exist, overwrites if it does.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["path", "content"],
			"properties": {
				"path": {
					"type": "string",
					"description": "The path to the file to write"
				},
				"content": {
					"type": "string",
					"description": "The content to write"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			path, _ := args["path"].(string)
			content, _ := args["content"].(string)
			if path == "" {
				return "", fmt.Errorf("path is required")
			}
			piArgs := map[string]interface{}{"path": path, "content": content}
			raw, _ := json.Marshal(piArgs)
			runtime := foundationpi.NewRuntime(".", nil)
			result, err := runtime.ExecuteTool(context.Background(), "", "write", raw, nil)
			if err != nil {
				return "", err
			}
			return formatFoundationToolResult(result), nil
		},
	})

	// edit_file - Gemini exact naming for search/replace editing
	r.Tools = append(r.Tools, Tool{
		Name:        "edit_file",
		Description: "Edit a file by providing the exact string to find and its replacement. The old_string must match exactly.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["path", "old_string", "new_string"],
			"properties": {
				"path": {
					"type": "string",
					"description": "The path to the file to edit"
				},
				"old_string": {
					"type": "string",
					"description": "The exact text to find in the file"
				},
				"new_string": {
					"type": "string",
					"description": "The text to replace old_string with"
				},
				"replace_all": {
					"type": "boolean",
					"description": "Replace all occurrences of old_string (default false)"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			path, _ := args["path"].(string)
			oldStr, _ := args["old_string"].(string)
			newStr, _ := args["new_string"].(string)

			if path == "" {
				return "", fmt.Errorf("path is required")
			}

			data, err := os.ReadFile(path)
			if err != nil {
				return "", fmt.Errorf("failed to read file: %w", err)
			}

			content := strings.ReplaceAll(string(data), "\r\n", "\n")
			if !strings.Contains(content, oldStr) {
				return "", fmt.Errorf("old_string not found in file. Make sure it matches exactly, including whitespace and line breaks")
			}

			replaceAll, _ := args["replace_all"].(bool)
			if replaceAll {
				content = strings.ReplaceAll(content, oldStr, newStr)
			} else {
				count := strings.Count(content, oldStr)
				if count > 1 {
					return "", fmt.Errorf("old_string appears %d times in the file. Please provide more context to ensure a unique match, or set replace_all to true", count)
				}
				content = strings.Replace(content, oldStr, newStr, 1)
			}

			if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
				return "", fmt.Errorf("failed to write file: %w", err)
			}

			return fmt.Sprintf("File edited successfully: %s", path), nil
		},
	})

	// list_directory - Gemini exact naming for ls
	r.Tools = append(r.Tools, Tool{
		Name:        "list_directory",
		Description: "List the contents of a directory.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"path": {
					"type": "string",
					"description": "The directory path to list"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			path, _ := args["path"].(string)
			if path == "" {
				path = "."
			}
			piArgs := map[string]interface{}{"path": path}
			raw, _ := json.Marshal(piArgs)
			runtime := foundationpi.NewRuntime(".", nil)
			result, err := runtime.ExecuteTool(context.Background(), "", "ls", raw, nil)
			if err != nil {
				return "", err
			}
			return formatFoundationToolResult(result), nil
		},
	})

	// shell - Gemini exact naming for bash
	r.Tools = append(r.Tools, Tool{
		Name:        "shell",
		Description: "Execute a shell command and return its output.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["command"],
			"properties": {
				"command": {
					"type": "string",
					"description": "The shell command to execute"
				},
				"working_directory": {
					"type": "string",
					"description": "Working directory for the command"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			command, _ := args["command"].(string)
			if command == "" {
				return "", fmt.Errorf("command is required")
			}
			piArgs := map[string]interface{}{"command": command}
			if wd, ok := args["working_directory"].(string); ok && wd != "" {
				piArgs["working_dir"] = wd
			}
			raw, _ := json.Marshal(piArgs)
			runtime := foundationpi.NewRuntime(".", nil)
			result, err := runtime.ExecuteTool(context.Background(), "", "bash", raw, nil)
			if err != nil {
				return "", err
			}
			return formatFoundationToolResult(result), nil
		},
	})

	// search_files - Gemini exact naming for grep
	r.Tools = append(r.Tools, Tool{
		Name:        "search_files",
		Description: "Search for a pattern in files. Returns matching lines with file paths and line numbers.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["pattern"],
			"properties": {
				"pattern": {
					"type": "string",
					"description": "The search pattern (supports regex)"
				},
				"path": {
					"type": "string",
					"description": "The directory or file to search in"
				},
				"ignore_case": {
					"type": "boolean",
					"description": "Case-insensitive search"
				},
				"file_pattern": {
					"type": "string",
					"description": "Filter files by glob pattern (e.g. '*.go')"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			pattern, _ := args["pattern"].(string)
			if pattern == "" {
				return "", fmt.Errorf("pattern is required")
			}
			piArgs := map[string]interface{}{
				"pattern": pattern,
			}
			if path, ok := args["path"].(string); ok && path != "" {
				piArgs["path"] = path
			}
			if ignoreCase, ok := args["ignore_case"].(bool); ok {
				piArgs["ignore_case"] = ignoreCase
			}
			if filePattern, ok := args["file_pattern"].(string); ok && filePattern != "" {
				piArgs["glob"] = filePattern
			}
			raw, _ := json.Marshal(piArgs)
			runtime := foundationpi.NewRuntime(".", nil)
			result, err := runtime.ExecuteTool(context.Background(), "", "grep", raw, nil)
			if err != nil {
				return "", err
			}
			return formatFoundationToolResult(result), nil
		},
	})

	// find_files - Gemini exact naming for find
	r.Tools = append(r.Tools, Tool{
		Name:        "find_files",
		Description: "Find files matching a pattern. Returns matching file paths.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["pattern"],
			"properties": {
				"pattern": {
					"type": "string",
					"description": "The glob pattern to match against filenames"
				},
				"path": {
					"type": "string",
					"description": "The directory to search in"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			pattern, _ := args["pattern"].(string)
			if pattern == "" {
				return "", fmt.Errorf("pattern is required")
			}
			piArgs := map[string]interface{}{"pattern": pattern}
			if path, ok := args["path"].(string); ok && path != "" {
				piArgs["path"] = path
			}
			raw, _ := json.Marshal(piArgs)
			runtime := foundationpi.NewRuntime(".", nil)
			result, err := runtime.ExecuteTool(context.Background(), "", "find", raw, nil)
			if err != nil {
				return "", err
			}
			return formatFoundationToolResult(result), nil
		},
	})
}

// registerOpenCodeTools adds OpenCode-compatible tool surfaces.
// OpenCode uses a search/replace block format for file editing.
func (r *Registry) registerOpenCodeTools() {
	// apply_search_replace - OpenCode's primary edit mechanism
	r.Tools = append(r.Tools, Tool{
		Name:        "apply_search_replace",
		Description: "Apply a search/replace block to a file. Uses the SEARCH/REPLACE diff format for precise file editing.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["file_path", "search_block", "replace_block"],
			"properties": {
				"file_path": {
					"type": "string",
					"description": "Path to the file to edit"
				},
				"search_block": {
					"type": "string",
					"description": "The exact text to find in the file"
				},
				"replace_block": {
					"type": "string",
					"description": "The text to replace the search block with"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			filePath, _ := args["file_path"].(string)
			searchBlock, _ := args["search_block"].(string)
			replaceBlock, _ := args["replace_block"].(string)

			if filePath == "" {
				return "", fmt.Errorf("file_path is required")
			}

			data, err := os.ReadFile(filePath)
			if err != nil {
				return "", fmt.Errorf("failed to read file: %w", err)
			}

			content := string(data)
			if !strings.Contains(content, searchBlock) {
				return "", fmt.Errorf("search block not found in file. Make sure it matches exactly, including whitespace and line breaks")
			}

			newContent := strings.Replace(content, searchBlock, replaceBlock, 1)

			if err := os.WriteFile(filePath, []byte(newContent), 0o644); err != nil {
				return "", fmt.Errorf("failed to write file: %w", err)
			}

			return fmt.Sprintf("Search/replace applied to %s", filePath), nil
		},
	})

	// apply_diff - OpenCode's diff application tool
	r.Tools = append(r.Tools, Tool{
		Name:        "apply_diff",
		Description: "Apply a unified diff to a file.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["file_path", "diff"],
			"properties": {
				"file_path": {
					"type": "string",
					"description": "Path to the file to apply the diff to"
				},
				"diff": {
					"type": "string",
					"description": "The unified diff to apply"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			filePath, _ := args["file_path"].(string)
			diffContent, _ := args["diff"].(string)

			if filePath == "" {
				return "", fmt.Errorf("file_path is required")
			}

			// Use the agent diff package for applying diffs
			piArgs := map[string]interface{}{
				"path":  filePath,
				"patch": diffContent,
			}
			raw, _ := json.Marshal(piArgs)
			runtime := foundationpi.NewRuntime(".", nil)
			result, err := runtime.ExecuteTool(context.Background(), "", "edit", raw, nil)
			if err != nil {
				return "", err
			}
			return formatFoundationToolResult(result), nil
		},
	})
}

// registerGrokTools adds Grok CLI compatible tool surfaces.
func (r *Registry) registerGrokTools() {
	// file_edit - Grok's primary editing tool
	r.Tools = append(r.Tools, Tool{
		Name:        "file_edit",
		Description: "Edit a file by specifying the exact content to find and replace.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["path", "find", "replace"],
			"properties": {
				"path": {
					"type": "string",
					"description": "Path to the file"
				},
				"find": {
					"type": "string",
					"description": "Exact text to find"
				},
				"replace": {
					"type": "string",
					"description": "Text to replace with"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			path, _ := args["path"].(string)
			find, _ := args["find"].(string)
			replace, _ := args["replace"].(string)

			piArgs := map[string]interface{}{
				"path": path,
				"edits": []map[string]interface{}{
					{"oldText": find, "newText": replace},
				},
			}
			raw, _ := json.Marshal(piArgs)
			runtime := foundationpi.NewRuntime(".", nil)
			result, err := runtime.ExecuteTool(context.Background(), "", "edit", raw, nil)
			if err != nil {
				return "", err
			}
			return formatFoundationToolResult(result), nil
		},
	})

	// file_read - Grok's reading tool
	r.Tools = append(r.Tools, Tool{
		Name:        "file_read",
		Description: "Read the contents of a file.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["path"],
			"properties": {
				"path": {
					"type": "string",
					"description": "Path to the file to read"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			path, _ := args["path"].(string)
			piArgs := map[string]interface{}{"path": path}
			raw, _ := json.Marshal(piArgs)
			runtime := foundationpi.NewRuntime(".", nil)
			result, err := runtime.ExecuteTool(context.Background(), "", "read", raw, nil)
			if err != nil {
				return "", err
			}
			return formatFoundationToolResult(result), nil
		},
	})

	// file_write - Grok's writing tool
	r.Tools = append(r.Tools, Tool{
		Name:        "file_write",
		Description: "Write content to a file.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["path", "content"],
			"properties": {
				"path": {
					"type": "string",
					"description": "Path to the file to write"
				},
				"content": {
					"type": "string",
					"description": "Content to write"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			path, _ := args["path"].(string)
			content, _ := args["content"].(string)
			piArgs := map[string]interface{}{"path": path, "content": content}
			raw, _ := json.Marshal(piArgs)
			runtime := foundationpi.NewRuntime(".", nil)
			result, err := runtime.ExecuteTool(context.Background(), "", "write", raw, nil)
			if err != nil {
				return "", err
			}
			return formatFoundationToolResult(result), nil
		},
	})

	// execute_command - Grok's shell tool
	r.Tools = append(r.Tools, Tool{
		Name:        "execute_command",
		Description: "Execute a shell command and return its output.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["command"],
			"properties": {
				"command": {
					"type": "string",
					"description": "The command to execute"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			command, _ := args["command"].(string)
			piArgs := map[string]interface{}{"command": command}
			raw, _ := json.Marshal(piArgs)
			runtime := foundationpi.NewRuntime(".", nil)
			result, err := runtime.ExecuteTool(context.Background(), "", "bash", raw, nil)
			if err != nil {
				return "", err
			}
			return formatFoundationToolResult(result), nil
		},
	})

	// scan_directory - Grok's directory listing tool
	r.Tools = append(r.Tools, Tool{
		Name:        "scan_directory",
		Description: "Scan a directory and list its contents with file types.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"path": {
					"type": "string",
					"description": "Path to the directory to scan"
				},
				"recursive": {
					"type": "boolean",
					"description": "Whether to scan recursively (default false)"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			path, _ := args["path"].(string)
			recursive, _ := args["recursive"].(bool)

			if path == "" {
				path = "."
			}

			if !recursive {
				piArgs := map[string]interface{}{"path": path}
				raw, _ := json.Marshal(piArgs)
				runtime := foundationpi.NewRuntime(".", nil)
				result, err := runtime.ExecuteTool(context.Background(), "", "ls", raw, nil)
				if err != nil {
					return "", err
				}
				return formatFoundationToolResult(result), nil
			}

			// Recursive listing using find
			piArgs := map[string]interface{}{"pattern": "*", "path": path}
			raw, _ := json.Marshal(piArgs)
			runtime := foundationpi.NewRuntime(".", nil)
			result, err := runtime.ExecuteTool(context.Background(), "", "find", raw, nil)
			if err != nil {
				return "", err
			}
			return formatFoundationToolResult(result), nil
		},
	})

	// search_code - Grok's code search tool
	r.Tools = append(r.Tools, Tool{
		Name:        "search_code",
		Description: "Search for a pattern across the codebase. Returns matching files and lines.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["pattern"],
			"properties": {
				"pattern": {
					"type": "string",
					"description": "The pattern to search for"
				},
				"path": {
					"type": "string",
					"description": "Directory to search in"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			pattern, _ := args["pattern"].(string)
			piArgs := map[string]interface{}{"pattern": pattern}
			if path, ok := args["path"].(string); ok && path != "" {
				piArgs["path"] = path
			}
			raw, _ := json.Marshal(piArgs)
			runtime := foundationpi.NewRuntime(".", nil)
			result, err := runtime.ExecuteTool(context.Background(), "", "grep", raw, nil)
			if err != nil {
				return "", err
			}
			return formatFoundationToolResult(result), nil
		},
	})
}

// registerCopilotCLITools adds GitHub Copilot CLI compatible tool surfaces.
func (r *Registry) registerCopilotCLITools() {
	// copilot_edit - Copilot's primary edit interface
	r.Tools = append(r.Tools, Tool{
		Name:        "copilot_edit",
		Description: "Edit a file using GitHub Copilot CLI format. Replaces old_string with new_string.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["file_path", "old_string", "new_string"],
			"properties": {
				"file_path": {
					"type": "string",
					"description": "Path to the file to edit"
				},
				"old_string": {
					"type": "string",
					"description": "Text to find"
				},
				"new_string": {
					"type": "string",
					"description": "Text to replace with"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			filePath, _ := args["file_path"].(string)
			oldStr, _ := args["old_string"].(string)
			newStr, _ := args["new_string"].(string)

			piArgs := map[string]interface{}{
				"path": filePath,
				"edits": []map[string]interface{}{
					{"oldText": oldStr, "newText": newStr},
				},
			}
			raw, _ := json.Marshal(piArgs)
			runtime := foundationpi.NewRuntime(".", nil)
			result, err := runtime.ExecuteTool(context.Background(), "", "edit", raw, nil)
			if err != nil {
				return "", err
			}
			return formatFoundationToolResult(result), nil
		},
	})
}

// registerAiderV2Tools adds Aider v2 search/replace block format tools.
func (r *Registry) registerAiderV2Tools() {
	// aider_search_replace - Aider's SEARCH/REPLACE block editing
	r.Tools = append(r.Tools, Tool{
		Name:        "aider_search_replace",
		Description: "Apply Aider-style SEARCH/REPLACE blocks to files. Each block finds the exact text and replaces it.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["blocks"],
			"properties": {
				"blocks": {
					"type": "array",
					"items": {
						"type": "object",
						"required": ["file_path", "search", "replace"],
						"properties": {
							"file_path": {
								"type": "string",
								"description": "Path to the file"
							},
							"search": {
								"type": "string",
								"description": "The exact text to search for"
							},
							"replace": {
								"type": "string",
								"description": "The replacement text"
							}
						}
					},
					"description": "Array of SEARCH/REPLACE blocks"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			blocksRaw, ok := args["blocks"].([]interface{})
			if !ok || len(blocksRaw) == 0 {
				return "", fmt.Errorf("blocks must be a non-empty array")
			}

			var results []string
			for i, blockRaw := range blocksRaw {
				block, ok := blockRaw.(map[string]interface{})
				if !ok {
					results = append(results, fmt.Sprintf("Block %d: invalid format", i))
					continue
				}

				filePath, _ := block["file_path"].(string)
				search, _ := block["search"].(string)
				replace, _ := block["replace"].(string)

				if filePath == "" {
					results = append(results, fmt.Sprintf("Block %d: file_path is required", i))
					continue
				}

				data, err := os.ReadFile(filePath)
				if err != nil {
					results = append(results, fmt.Sprintf("Block %d on %s: read error: %v", i, filePath, err))
					continue
				}

				content := string(data)
				if !strings.Contains(content, search) {
					results = append(results, fmt.Sprintf("Block %d on %s: search text not found", i, filePath))
					continue
				}

				newContent := strings.Replace(content, search, replace, 1)
				if err := os.WriteFile(filePath, []byte(newContent), 0o644); err != nil {
					results = append(results, fmt.Sprintf("Block %d on %s: write error: %v", i, filePath, err))
					continue
				}

				results = append(results, fmt.Sprintf("Block %d on %s: applied successfully", i, filepath.Base(filePath)))
			}

			return strings.Join(results, "\n"), nil
		},
	})
}
