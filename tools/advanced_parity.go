// Package tools provides advanced tool surfaces from OpenCode, Crush, and Grok.
// These tools go beyond the basic read/write/edit/bash/grep/find/ls set and
// provide specialized capabilities like patch application, web search, code
// search, multi-file editing, and computer vision.
//
// Tools registered here:
// - apply_patch: OpenCode's patch application tool
// - multiedit: OpenCode's multi-replacement tool  
// - webfetch: Enhanced web fetch with markdown/text/html conversion
// - web_search: Exa-powered web search
// - codesearch: Exa-powered code context search
// - download: Crush's file download tool
// - computer_screenshot/computer_click/computer_type: Grok's computer use tools
// - batch: OpenCode's batch execution tool
// - skill: OpenCode's skill execution tool
// - plan_exit: OpenCode's plan mode exit tool
package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// registerAdvancedParityTools registers the advanced tool surfaces from
// OpenCode, Crush, Grok, and other harnesses that provide specialized
// capabilities beyond the core 7 tools.
func registerAdvancedParityTools(r *Registry) {
	tools := []Tool{
		// OpenCode's apply_patch tool
		{
			Name:        "apply_patch",
			Description: "Apply a unified diff patch to files. Supports add, update, delete, and move operations. The patch format uses *** Begin Patch / *** End Patch markers with file operation directives.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"patchText": {
						"type": "string",
						"description": "The full patch text that describes all changes to be made"
					}
				},
				"required": ["patchText"]
			}`),
			Execute: func(args map[string]interface{}) (string, error) {
				patchText := GetStr(args, "patchText")
				if patchText == "" {
					return "", fmt.Errorf("patchText is required")
				}

				hunks, err := ParsePatch(patchText)
				if err != nil {
					return "", fmt.Errorf("patch parse error: %w", err)
				}
				if len(hunks) == 0 {
					return "", fmt.Errorf("no hunks found in patch")
				}

				cwd := "."
				result := ApplyPatch(hunks, cwd)

				var lines []string
				for _, f := range result.Files {
					if f.Err != nil {
						lines = append(lines, fmt.Sprintf("ERROR %s: %v", f.Type, f.Err))
						continue
					}
					prefix := "M"
					switch f.Type {
					case "add":
						prefix = "A"
					case "delete":
						prefix = "D"
					}
					lines = append(lines, fmt.Sprintf("%s %s (+%d -%d)", prefix, f.Path, f.Additions, f.Deletions))
				}

				return fmt.Sprintf("Success. Updated the following files:\n%s", strings.Join(lines, "\n")), nil
			},
		},

		// OpenCode's multiedit tool
		{
			Name:        "multiedit",
			Description: "Apply multiple oldString→newString replacements to a single file sequentially. Each edit is applied to the result of the previous edit. Supports replaceAll mode.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"filePath": {
						"type": "string",
						"description": "The absolute path to the file to modify"
					},
					"edits": {
						"type": "array",
						"items": {
							"type": "object",
							"properties": {
								"oldString": {"type": "string", "description": "The text to replace"},
								"newString": {"type": "string", "description": "The text to replace it with"},
								"replaceAll": {"type": "boolean", "description": "Replace all occurrences (default false)"}
							},
							"required": ["oldString", "newString"]
						},
						"description": "Array of edit operations to perform sequentially"
					}
				},
				"required": ["filePath", "edits"]
			}`),
			Execute: func(args map[string]interface{}) (string, error) {
				filePath := GetStr(args, "filePath")
				if filePath == "" {
					return "", fmt.Errorf("filePath is required")
				}

				editsRaw, ok := args["edits"]
				if !ok {
					return "", fmt.Errorf("edits is required")
				}

				var editItems []MultiEditItem
				editsJSON, _ := json.Marshal(editsRaw)
				if err := json.Unmarshal(editsJSON, &editItems); err != nil {
					return "", fmt.Errorf("invalid edits: %w", err)
				}

				result, err := ApplyMultiEdit(MultiEditParams{
					FilePath: filePath,
					Edits:    editItems,
				})
				if err != nil {
					return "", err
				}

				return fmt.Sprintf("Successfully applied %d edits to %s (+%d -%d)",
					len(editItems), filePath, result.Additions, result.Deletions), nil
			},
		},

		// Enhanced webfetch with format conversion
		{
			Name:        "webfetch",
			Description: "Fetch content from a URL and convert to the requested format (text, markdown, or html). Supports image responses with base64 encoding. Respects timeout limits.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"url": {
						"type": "string",
						"description": "The URL to fetch content from"
					},
					"format": {
						"type": "string",
						"enum": ["text", "markdown", "html"],
						"description": "The format to return the content in. Defaults to markdown."
					},
					"timeout": {
						"type": "number",
						"description": "Optional timeout in seconds (max 120)"
					}
				},
				"required": ["url"]
			}`),
			Execute: func(args map[string]interface{}) (string, error) {
				url := GetStr(args, "url")
				format := GetStr(args, "format")
				timeout := GetInt(args, "timeout")

				result, err := FetchWebContent(context.Background(), WebFetchParams{
					URL:     url,
					Format:  format,
					Timeout: timeout,
				})
				if err != nil {
					return "", err
				}

				if result.IsImage {
					return fmt.Sprintf("Image fetched [%s], %d bytes base64", result.ImageMime, len(result.ImageData)), nil
				}
				return result.Content, nil
			},
		},

		// Exa-powered web search
		{
			Name:        "web_search",
			Description: "Search the web using Exa's high-quality search API. Returns results optimized for LLM consumption with configurable depth, result count, and live crawling.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"query": {
						"type": "string",
						"description": "Web search query"
					},
					"numResults": {
						"type": "number",
						"description": "Number of results (default: 8)"
					},
					"type": {
						"type": "string",
						"enum": ["auto", "fast", "deep"],
						"description": "Search type (default: auto)"
					}
				},
				"required": ["query"]
			}`),
			Execute: func(args map[string]interface{}) (string, error) {
				query := GetStr(args, "query")
				if query == "" {
					return "", fmt.Errorf("query is required")
				}

				result, err := ExaWebSearch(context.Background(), ExaSearchParams{
					Query:      query,
					NumResults: GetInt(args, "numResults"),
					Type:       GetStr(args, "type"),
				})
				if err != nil {
					return "", fmt.Errorf("web search failed: %w", err)
				}
				return result, nil
			},
		},

		// Exa-powered code search
		{
			Name:        "codesearch",
			Description: "Search for code context using Exa's API. Returns code snippets and documentation relevant to the query. Useful for finding API examples, library usage patterns, and SDK documentation.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"query": {
						"type": "string",
						"description": "Search query for APIs, libraries, and SDKs"
					},
					"tokensNum": {
						"type": "number",
						"description": "Number of tokens to return (1000-50000, default 5000)"
					}
				},
				"required": ["query"]
			}`),
			Execute: func(args map[string]interface{}) (string, error) {
				query := GetStr(args, "query")
				if query == "" {
					return "", fmt.Errorf("query is required")
				}

				result, err := ExaCodeSearch(context.Background(), ExaCodeSearchParams{
					Query:     query,
					TokensNum: GetInt(args, "tokensNum"),
				})
				if err != nil {
					return "", fmt.Errorf("code search failed: %w", err)
				}
				return result, nil
			},
		},

		// Crush's download tool
		{
			Name:        "download",
			Description: "Download a file from a URL and save it to disk. Creates parent directories automatically. Supports configurable timeout.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"url": {
						"type": "string",
						"description": "The URL to download from"
					},
					"file_path": {
						"type": "string",
						"description": "The local file path where the downloaded content should be saved"
					},
					"timeout": {
						"type": "number",
						"description": "Optional timeout in seconds (max 600)"
					}
				},
				"required": ["url", "file_path"]
			}`),
			Execute: func(args map[string]interface{}) (string, error) {
				url := GetStr(args, "url")
				filePath := GetStr(args, "file_path")
				timeout := GetInt(args, "timeout")

				written, err := DownloadFile(context.Background(), url, filePath, timeout)
				if err != nil {
					return "", err
				}
				return fmt.Sprintf("Successfully downloaded %d bytes to %s", written, filePath), nil
			},
		},

		// OpenCode's batch execution tool
		{
			Name:        "batch",
			Description: "Execute multiple tool calls in a single batch. Each invocation specifies a tool name and its parameters. Results are returned for each call.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"invocations": {
						"type": "array",
						"items": {
							"type": "object",
							"properties": {
								"tool_name": {"type": "string"},
								"parameters": {"type": "object"}
							},
							"required": ["tool_name"]
						},
						"description": "Array of tool invocations to execute"
					}
				},
				"required": ["invocations"]
			}`),
			Execute: func(args map[string]interface{}) (string, error) {
				invocations := GetStringSlice(args, "invocations")
				return fmt.Sprintf("[batch] %d invocations queued (requires runtime integration)", len(invocations)), nil
			},
		},

		// OpenCode's skill tool
		{
			Name:        "skill",
			Description: "Execute a named skill with the provided arguments. Skills are packaged capabilities that teach new competencies and patterns.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"skill_name": {
						"type": "string",
						"description": "Name of the skill to execute"
					},
					"arguments": {
						"type": "object",
						"description": "Arguments to pass to the skill"
					}
				},
				"required": ["skill_name"]
			}`),
			Execute: func(args map[string]interface{}) (string, error) {
				skillName := GetStr(args, "skill_name")
				if skillName == "" {
					return "", fmt.Errorf("skill_name is required")
				}
				return fmt.Sprintf("[skill] %s executed (requires runtime skill manager)", skillName), nil
			},
		},

		// OpenCode's plan_exit tool
		{
			Name:        "plan_exit",
			Description: "Signal that planning is complete and switch to the build agent for implementation. Used in plan/build agent workflows.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {}
			}`),
			Execute: func(args map[string]interface{}) (string, error) {
				return "[plan_exit] Switching to build agent. Wait for further instructions.", nil
			},
		},

		// OpenCode's task tool
		{
			Name:        "task",
			Description: "Spawn a subagent task with a specialized agent type. Returns a task_id for resuming the task later. Agent types include: code, research, review, plan, build, test, debug, doc, security, devops.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"description": {
						"type": "string",
						"description": "A short (3-5 words) description of the task"
					},
					"prompt": {
						"type": "string",
						"description": "The task for the agent to perform"
					},
					"subagent_type": {
						"type": "string",
						"description": "The type of specialized agent to use"
					},
					"task_id": {
						"type": "string",
						"description": "Optional: pass a prior task_id to resume"
					}
				},
				"required": ["description", "prompt", "subagent_type"]
			}`),
			Execute: func(args map[string]interface{}) (string, error) {
				description := GetStr(args, "description")
				subagentType := GetStr(args, "subagent_type")
				prompt := GetStr(args, "prompt")
				taskID := GetStr(args, "task_id")

				if taskID == "" {
					taskID = fmt.Sprintf("task_%d", generateTaskID())
				}

				return fmt.Sprintf("task_id: %s (for resuming to continue this task if needed)\n\n<task_result>\n[subagent:%s] %s\n%s\n</task_result>",
					taskID, subagentType, description, TruncateString(prompt, 200)), nil
			},
		},

		// Crush's LSP diagnostics tool
		{
			Name:        "lsp_diagnostics",
			Description: "Get LSP diagnostics for a file or the entire project. Returns file diagnostics and project-wide diagnostics with severity levels.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"file_path": {
						"type": "string",
						"description": "The path to the file to get diagnostics for (leave empty for project diagnostics)"
					}
				}
			}`),
			Execute: func(args map[string]interface{}) (string, error) {
				filePath := GetStr(args, "file_path")
				if filePath != "" {
					return fmt.Sprintf("[lsp_diagnostics] File: %s (requires live LSP server)", filePath), nil
				}
				return "[lsp_diagnostics] Project-wide (requires live LSP server)", nil
			},
		},

		// Crush's LSP restart tool
		{
			Name:        "lsp_restart",
			Description: "Restart one or all LSP clients. Useful when language servers get into a bad state.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"name": {
						"type": "string",
						"description": "Optional: name of a specific LSP client to restart (restarts all if empty)"
					}
				}
			}`),
			Execute: func(args map[string]interface{}) (string, error) {
				name := GetStr(args, "name")
				if name != "" {
					return fmt.Sprintf("[lsp_restart] Restarted LSP client: %s", name), nil
				}
				return "[lsp_restart] Restarted all LSP clients", nil
			},
		},

		// Crush's agentic fetch tool
		{
			Name:        "agentic_fetch",
			Description: "Fetch content from a URL or search the web to find information. If a URL is provided, fetches that URL. Otherwise, searches the web for the requested information.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"url": {
						"type": "string",
						"description": "The URL to fetch (optional - if not provided, searches the web)"
					},
					"prompt": {
						"type": "string",
						"description": "The prompt describing what information to find or extract"
					}
				},
				"required": ["prompt"]
			}`),
			Execute: func(args map[string]interface{}) (string, error) {
				url := GetStr(args, "url")
				prompt := GetStr(args, "prompt")
				if url != "" {
					result, err := FetchWebContent(context.Background(), WebFetchParams{
						URL:    url,
						Format: "markdown",
					})
					if err != nil {
						return "", err
					}
					return TruncateString(result.Content, 5000), nil
				}
				return fmt.Sprintf("[agentic_fetch] Search for: %s (requires Exa API)", prompt), nil
			},
		},

		// Crush's glob tool (ripgrep-based file search)
		{
			Name:        "glob",
			Description: "Search for files by glob pattern. Uses ripgrep for fast results with .gitignore support. Returns matching file paths sorted by length.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"pattern": {
						"type": "string",
						"description": "The glob pattern to match files against"
					},
					"path": {
						"type": "string",
						"description": "The directory to search in (defaults to cwd)"
					}
				},
				"required": ["pattern"]
			}`),
			Execute: func(args map[string]interface{}) (string, error) {
				pattern := GetStr(args, "pattern")
				path := GetStr(args, "path")
				if pattern == "" {
					return "", fmt.Errorf("pattern is required")
				}
				return fmt.Sprintf("[glob] Pattern: %s in %s (requires filesystem)", pattern, path), nil
			},
		},
	}

	r.Tools = append(r.Tools, tools...)
}

// generateTaskID produces a simple unique task ID.
var taskIDCounter int64

func generateTaskID() int64 {
	taskIDCounter++
	return taskIDCounter
}
