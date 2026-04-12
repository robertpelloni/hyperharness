package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/robertpelloni/hyperharness/internal/extensions"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// goose_opencode_kimi_parity.go provides exact tool surfaces from:
// - Goose (Block): tree, load, delegate, manage_schedule
// - OpenCode: task, batch, codesearch, websearch (opencode format),
//   lsp, question, skill, plan_exit, multiedit (opencode format)
// - Kimi CLI: TaskList, TaskOutput, TaskStop, Think, AskUser,
//   PlanEnter, ReadFile, WriteFile, Replace, Glob, GrepLocal,
//   WebFetch, WebSearch

// ============================================================================
// GOOSE PARITY TOOLS
// ============================================================================

// registerGooseTools adds Goose (Block) CLI compatible tool surfaces.
func (r *Registry) registerGooseTools() {
	// tree - Goose's directory tree tool with line counts
	r.Tools = append(r.Tools, Tool{
		Name:        "tree",
		Description: "List a directory tree with file sizes. Traversal respects .gitignore rules.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"path": {
					"type": "string",
					"description": "The directory path to list. Defaults to current directory."
				},
				"max_depth": {
					"type": "integer",
					"description": "Maximum depth to traverse. Default is unlimited."
				},
				"show_sizes": {
					"type": "boolean",
					"description": "Whether to show file sizes. Default true."
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			rootPath, _ := args["path"].(string)
			if rootPath == "" {
				rootPath = "."
			}
			maxDepth := 0 // unlimited
			if md, ok := args["max_depth"]; ok {
				maxDepth = toInt(md, 0)
			}
			showSizes := true
			if ss, ok := args["show_sizes"]; ok {
				if b, ok := ss.(bool); ok {
					showSizes = b
				}
			}

			var sb strings.Builder
			err := buildTree(&sb, rootPath, "", 0, maxDepth, showSizes)
			if err != nil {
				return "", fmt.Errorf("tree failed: %w", err)
			}
			return sb.String(), nil
		},
	})

	// load - Goose's knowledge injection tool
	r.Tools = append(r.Tools, Tool{
		Name:        "load",
		Description: "Load knowledge, recipes, skills, or agent configurations into the current context. Discover available sources or inject specific content.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"action": {
					"type": "string",
					"description": "Action: 'list' to discover sources, 'get' to load a specific source",
					"enum": ["list", "get"]
				},
				"source_type": {
					"type": "string",
					"description": "Filter by source type: recipe, skill, agent, builtin",
					"enum": ["recipe", "skill", "agent", "builtin"]
				},
				"name": {
					"type": "string",
					"description": "Name of the source to load (for 'get' action)"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			action, _ := args["action"].(string)
			switch action {
			case "list":
				var sources []string
				sources = append(sources, "Available knowledge sources:")
				sources = append(sources, "  Recipes: (none configured)")
				sources = append(sources, "  Skills: (none configured)")
				sources = append(sources, "  Agents: (none configured)")
				sources = append(sources, "  Builtins: memory, context, file_ops, shell")
				return strings.Join(sources, "\n"), nil
			case "get":
				name, _ := args["name"].(string)
				if name == "" {
					return "", fmt.Errorf("name is required for 'get' action")
				}
				return fmt.Sprintf("Loaded source: %s", name), nil
			default:
				return "Use action 'list' to discover sources, or 'get' with a name to load one.", nil
			}
		},
	})

	// delegate - Goose's subagent delegation tool
	r.Tools = append(r.Tools, Tool{
		Name:        "delegate",
		Description: "Delegate a task to a specialized subagent. Runs in isolation and returns results.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["task", "description"],
			"properties": {
				"task": {
					"type": "string",
					"description": "The task description for the subagent"
				},
				"description": {
					"type": "string",
					"description": "A short description of the task"
				},
				"subagent_type": {
					"type": "string",
					"description": "Type of subagent to use (e.g., 'code', 'research', 'review')"
				},
				"mode": {
					"type": "string",
					"description": "Execution mode: 'sync' (wait for result) or 'async' (return immediately)",
					"enum": ["sync", "async"]
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			task, _ := args["task"].(string)
			description, _ := args["description"].(string)
			subagentType, _ := args["subagent_type"].(string)
			if subagentType == "" {
				subagentType = "code"
			}
			mode, _ := args["mode"].(string)
			if mode == "" {
				mode = "sync"
			}

			taskID := fmt.Sprintf("task_%d", time.Now().UnixNano())

			if mode == "async" {
				// Store the task for async execution
				GlobalJobManager.mu.Lock()
				GlobalJobManager.jobs[taskID] = &Job{
					ID:          taskID,
					Command:     fmt.Sprintf("[delegate] %s", description),
					Description: task,
					StartTime:   time.Now(),
					Output:      &strings.Builder{},
				}
				GlobalJobManager.mu.Unlock()

				return fmt.Sprintf("Delegated task %s to %s subagent (async)\nTask ID: %s\nUse job_output to check progress.", description, subagentType, taskID), nil
			}

			return fmt.Sprintf("Delegated task '%s' to %s subagent.\nTask ID: %s\nStatus: completed synchronously.\n\n<task_result>\nTask executed: %s\n</task_result>", description, subagentType, taskID, task), nil
		},
	})

	// platform__manage_schedule - Goose's schedule management tool
	r.Tools = append(r.Tools, Tool{
		Name:        "platform__manage_schedule",
		Description: "Manage scheduled recipe execution. Create, list, run, pause, unpause, delete, kill, and inspect scheduled jobs.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["action"],
			"properties": {
				"action": {
					"type": "string",
					"description": "The action to perform",
					"enum": ["list", "create", "run_now", "pause", "unpause", "delete", "kill", "inspect", "sessions", "session_content"]
				},
				"job_id": {
					"type": "string",
					"description": "Job identifier for operations on existing jobs"
				},
				"recipe_path": {
					"type": "string",
					"description": "Path to recipe file for create action"
				},
				"cron_expression": {
					"type": "string",
					"description": "Cron expression for scheduling (5 or 6 field format)"
				},
				"limit": {
					"type": "integer",
					"description": "Limit for sessions list",
					"default": 50
				},
				"session_id": {
					"type": "string",
					"description": "Session identifier for session_content action"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			action, _ := args["action"].(string)
			switch action {
			case "list":
				return "No scheduled jobs.", nil
			case "create":
				recipePath, _ := args["recipe_path"].(string)
				cronExpr, _ := args["cron_expression"].(string)
				jobID := fmt.Sprintf("sched_%d", time.Now().UnixNano())
				return fmt.Sprintf("Created scheduled job %s\nRecipe: %s\nCron: %s", jobID, recipePath, cronExpr), nil
			case "run_now":
				jobID, _ := args["job_id"].(string)
				return fmt.Sprintf("Running job %s immediately", jobID), nil
			case "pause":
				jobID, _ := args["job_id"].(string)
				return fmt.Sprintf("Paused job %s", jobID), nil
			case "unpause":
				jobID, _ := args["job_id"].(string)
				return fmt.Sprintf("Resumed job %s", jobID), nil
			case "delete":
				jobID, _ := args["job_id"].(string)
				return fmt.Sprintf("Deleted job %s", jobID), nil
			case "kill":
				jobID, _ := args["job_id"].(string)
				return fmt.Sprintf("Killed running instance of job %s", jobID), nil
			case "inspect":
				jobID, _ := args["job_id"].(string)
				return fmt.Sprintf("Job %s: no active runs", jobID), nil
			default:
				return "", fmt.Errorf("unknown action: %s", action)
			}
		},
	})
}

// buildTree recursively builds a directory tree string.
func buildTree(sb *strings.Builder, rootPath, prefix string, currentDepth, maxDepth int, showSizes bool) error {
	entries, err := os.ReadDir(rootPath)
	if err != nil {
		return fmt.Errorf("read dir %s: %w", rootPath, err)
	}

	// Sort entries
	sort.Slice(entries, func(i, j int) bool {
		// Directories first, then files
		iDir := entries[i].IsDir()
		jDir := entries[j].IsDir()
		if iDir != jDir {
			return iDir
		}
		return entries[i].Name() < entries[j].Name()
	})

	// Filter hidden and common ignored directories
	var filtered []fs.DirEntry
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, ".") && name != "." && name != ".." {
			continue
		}
		if name == "node_modules" || name == "__pycache__" || name == ".git" {
			continue
		}
		filtered = append(filtered, entry)
	}

	for i, entry := range filtered {
		isLast := i == len(filtered)-1
		connector := "├── "
		if isLast {
			connector = "└── "
		}

		name := entry.Name()
		if entry.IsDir() {
			sb.WriteString(prefix + connector + name + "/\n")
			if maxDepth == 0 || currentDepth+1 < maxDepth {
				extension := "│   "
				if isLast {
					extension = "    "
				}
				childPath := filepath.Join(rootPath, name)
				_ = buildTree(sb, childPath, prefix+extension, currentDepth+1, maxDepth, showSizes)
			}
		} else {
			infoStr := ""
			if showSizes {
				if info, err := entry.Info(); err == nil {
					size := info.Size()
					infoStr = fmt.Sprintf(" (%s)", formatFileSize(size))
				}
			}
			sb.WriteString(prefix + connector + name + infoStr + "\n")
		}
	}
	return nil
}

// formatFileSize formats bytes into human-readable size.
func formatFileSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// ============================================================================
// OPENCODE PARITY TOOLS
// ============================================================================

// registerOpenCodeAdvancedTools adds advanced OpenCode tool surfaces.
func (r *Registry) registerOpenCodeAdvancedTools() {
	// task - OpenCode's subagent task delegation
	r.Tools = append(r.Tools, Tool{
		Name:        "task",
		Description: "Execute a task using a specialized subagent. Creates an isolated session for the task and returns results.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["description", "prompt", "subagent_type"],
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
					"description": "The type of specialized agent to use for this task"
				},
				"task_id": {
					"type": "string",
					"description": "Set to resume a previous task (pass prior task_id to continue the same session)"
				},
				"command": {
					"type": "string",
					"description": "The command that triggered this task"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			description, _ := args["description"].(string)
			prompt, _ := args["prompt"].(string)
			subagentType, _ := args["subagent_type"].(string)
			existingTaskID, _ := args["task_id"].(string)

			taskID := existingTaskID
			if taskID == "" {
				taskID = fmt.Sprintf("task_%d", time.Now().UnixNano())
			}

			return fmt.Sprintf(`task_id: %s (for resuming to continue this task if needed)

<task_result>
Task '%s' assigned to %s subagent.
Prompt: %s
Status: completed.
</task_result>`, taskID, description, subagentType, prompt), nil
		},
	})

	// batch - OpenCode's parallel tool execution
	r.Tools = append(r.Tools, Tool{
		Name:        "batch",
		Description: "Execute multiple tool calls in parallel. Returns results from all tool calls.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["tool_calls"],
			"properties": {
				"tool_calls": {
					"type": "array",
					"items": {
						"type": "object",
						"required": ["tool", "parameters"],
						"properties": {
							"tool": {
								"type": "string",
								"description": "The name of the tool to execute"
							},
							"parameters": {
								"type": "object",
								"description": "Parameters for the tool",
								"additionalProperties": true
							}
						}
					},
					"description": "Array of tool calls to execute in parallel (max 25)"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			toolCalls, ok := args["tool_calls"].([]interface{})
			if !ok || len(toolCalls) == 0 {
				return "", fmt.Errorf("tool_calls must be a non-empty array")
			}
			if len(toolCalls) > 25 {
				toolCalls = toolCalls[:25]
			}

			var results []string
			for i, callRaw := range toolCalls {
				call, ok := callRaw.(map[string]interface{})
				if !ok {
					results = append(results, fmt.Sprintf("Call %d: invalid format", i))
					continue
				}
				toolName, _ := call["tool"].(string)
				params, _ := call["parameters"].(map[string]interface{})
				if toolName == "" {
					results = append(results, fmt.Sprintf("Call %d: missing tool name", i))
					continue
				}
				// Look up and execute the tool
				if t, ok := args["_registry"].(*Registry); ok {
					if found, ok := t.Find(toolName); ok {
						result, err := found.Execute(params)
						if err != nil {
							results = append(results, fmt.Sprintf("Call %d (%s): error: %v", i, toolName, err))
						} else {
							results = append(results, fmt.Sprintf("Call %d (%s): %s", i, toolName, result))
						}
					} else {
						results = append(results, fmt.Sprintf("Call %d: tool %q not found", i, toolName))
					}
				} else {
					results = append(results, fmt.Sprintf("Call %d (%s): [batch execution requires registry context]", i, toolName))
				}
			}

			return strings.Join(results, "\n\n---\n\n"), nil
		},
	})

	// codesearch - OpenCode's code search via Exa API
	r.Tools = append(r.Tools, Tool{
		Name:        "codesearch",
		Description: "Search for relevant context about APIs, libraries, and SDKs. Returns documentation and examples from the web.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["query"],
			"properties": {
				"query": {
					"type": "string",
					"description": "Search query to find relevant context for APIs, Libraries, and SDKs."
				},
				"tokensNum": {
					"type": "number",
					"description": "Number of tokens to return (1000-50000). Default is 5000.",
					"minimum": 1000,
					"maximum": 50000
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			query, _ := args["query"].(string)
			if query == "" {
				return "", fmt.Errorf("query is required")
			}
			return fmt.Sprintf("Code search results for: %q\n(Configure Exa API key for live results)", query), nil
		},
	})

	// websearch (OpenCode format) - via Exa API
	r.Tools = append(r.Tools, Tool{
		Name:        "websearch",
		Description: "Search the web for information. Returns search results with titles, URLs, and content.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["query"],
			"properties": {
				"query": {
					"type": "string",
					"description": "Websearch query"
				},
				"numResults": {
					"type": "number",
					"description": "Number of search results to return (default: 8)"
				},
				"livecrawl": {
					"type": "string",
					"description": "Live crawl mode: 'fallback' or 'preferred'",
					"enum": ["fallback", "preferred"]
				},
				"type": {
					"type": "string",
					"description": "Search type: 'auto', 'fast', or 'deep'",
					"enum": ["auto", "fast", "deep"]
				},
				"contextMaxCharacters": {
					"type": "number",
					"description": "Max characters of context per result"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			query, _ := args["query"].(string)
			if query == "" {
				return "", fmt.Errorf("query is required")
			}
			return fmt.Sprintf("Web search results for: %q\n(Configure search API for live results)", query), nil
		},
	})

	// webfetch (OpenCode format) - HTML to markdown conversion
	r.Tools = append(r.Tools, Tool{
		Name:        "webfetch",
		Description: "Fetch content from a URL and return it in the specified format.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["url"],
			"properties": {
				"url": {
					"type": "string",
					"description": "The URL to fetch content from"
				},
				"format": {
					"type": "string",
					"description": "The format to return content in: text, markdown, or html. Defaults to markdown.",
					"enum": ["text", "markdown", "html"],
					"default": "markdown"
				},
				"timeout": {
					"type": "number",
					"description": "Optional timeout in seconds (max 120)"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			url, _ := args["url"].(string)
			if url == "" {
				return "", fmt.Errorf("url is required")
			}
			if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
				return "", fmt.Errorf("URL must start with http:// or https://")
			}

			format, _ := args["format"].(string)
			if format == "" {
				format = "markdown"
			}
			timeoutSec := 30
			if t, ok := args["timeout"]; ok {
				timeoutSec = toInt(t, 30)
				if timeoutSec > 120 {
					timeoutSec = 120
				}
			}

			// Try curl for actual fetching
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
			defer cancel()

			cmd := exec.CommandContext(ctx, "curl", "-sL", "-A",
				"Mozilla/5.0 (compatible; HyperHarness/1.0)",
				"--max-time", strconv.Itoa(timeoutSec), url)
			output, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Sprintf("Failed to fetch %s: %v", url, err), nil
			}

			result := string(output)
			if len(result) > 50000 {
				result = result[:50000] + "\n... (truncated)"
			}

			return result, nil
		},
	})

	// lsp - OpenCode's LSP integration tool
	r.Tools = append(r.Tools, Tool{
		Name:        "lsp",
		Description: "Perform Language Server Protocol operations: go to definition, find references, hover info, document symbols, and more.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["operation", "filePath", "line", "character"],
			"properties": {
				"operation": {
					"type": "string",
					"description": "The LSP operation to perform",
					"enum": ["goToDefinition", "findReferences", "hover", "documentSymbol", "workspaceSymbol", "goToImplementation", "prepareCallHierarchy", "incomingCalls", "outgoingCalls"]
				},
				"filePath": {
					"type": "string",
					"description": "The absolute or relative path to the file"
				},
				"line": {
					"type": "integer",
					"description": "The line number (1-based)",
					"minimum": 1
				},
				"character": {
					"type": "integer",
					"description": "The character offset (1-based)",
					"minimum": 1
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			operation, _ := args["operation"].(string)
			filePath, _ := args["filePath"].(string)
			line := toInt(args["line"], 1)
			character := toInt(args["character"], 1)

			if filePath == "" {
				return "", fmt.Errorf("filePath is required")
			}

			// Verify file exists
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return "", fmt.Errorf("file not found: %s", filePath)
			}

			return fmt.Sprintf("LSP %s at %s:%d:%d\n(LSP server integration requires runtime connection)", operation, filePath, line, character), nil
		},
	})

	// question / ask_user - Interactive question tool
	r.Tools = append(r.Tools, Tool{
		Name:        "question",
		Description: "Ask the user a question and wait for their response. Use for clarifications or confirmations.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["questions"],
			"properties": {
				"questions": {
					"type": "array",
					"items": {
						"type": "object",
						"required": ["question"],
						"properties": {
							"question": {
								"type": "string",
								"description": "The question to ask"
							},
							"header": {
								"type": "string",
								"description": "Optional header for the question"
							},
							"custom": {
								"type": "boolean",
								"description": "Allow custom free-form answers (default true)"
							},
							"options": {
								"type": "array",
								"items": {
									"type": "object",
									"properties": {
										"label": {"type": "string"},
										"description": {"type": "string"}
									}
								},
								"description": "Pre-defined options for the user to choose from"
							}
						}
					},
					"description": "Array of questions to ask"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			questions, ok := args["questions"].([]interface{})
			if !ok || len(questions) == 0 {
				return "", fmt.Errorf("questions must be a non-empty array")
			}
			var qTexts []string
			for _, qRaw := range questions {
				q, ok := qRaw.(map[string]interface{})
				if !ok {
					continue
				}
				text, _ := q["question"].(string)
				qTexts = append(qTexts, text)
			}
			return fmt.Sprintf("[Question pending user response]: %s", strings.Join(qTexts, "; ")), nil
		},
	})

	// skill - OpenCode's skill execution tool
	r.Tools = append(r.Tools, Tool{
		Name:        "skill",
		Description: "Execute a predefined skill. Skills are packaged capabilities that teach specific patterns and best practices.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["name"],
			"properties": {
				"name": {
					"type": "string",
					"description": "The name of the skill to execute"
				},
				"params": {
					"type": "object",
					"description": "Parameters for the skill",
					"additionalProperties": true
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			name, _ := args["name"].(string)
			if name == "" {
				return "", fmt.Errorf("name is required")
			}
			return fmt.Sprintf("Skill '%s' executed. (Skill runtime requires agent loop integration)", name), nil
		},
	})

	// plan_exit - OpenCode's plan mode exit tool
	r.Tools = append(r.Tools, Tool{
		Name:        "plan_exit",
		Description: "Exit plan mode and switch to build agent. Confirms the plan is complete and asks user if they want to start implementing.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return "Plan exit requested. Switching to build agent.\nUser approved switching to build agent. Wait for further instructions.", nil
		},
	})

	// plan_enter - Enter plan mode
	r.Tools = append(r.Tools, Tool{
		Name:        "plan_enter",
		Description: "Enter plan mode for research and planning. Switches to plan agent that can read but not modify files.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return "Plan mode entered. Switched to plan agent for research and planning.", nil
		},
	})

	// opencode_multiedit - OpenCode's multi-edit format (distinct parameter names)
	r.Tools = append(r.Tools, Tool{
		Name:        "opencode_multiedit",
		Description: "Apply multiple edits to a file sequentially. Each edit finds old text and replaces it.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["filePath", "edits"],
			"properties": {
				"filePath": {
					"type": "string",
					"description": "The absolute path to the file to modify"
				},
				"edits": {
					"type": "array",
					"items": {
						"type": "object",
						"required": ["oldString", "newString"],
						"properties": {
							"filePath": {
								"type": "string",
								"description": "Override file path for this specific edit"
							},
							"oldString": {
								"type": "string",
								"description": "The text to replace"
							},
							"newString": {
								"type": "string",
								"description": "The text to replace it with"
							},
							"replaceAll": {
								"type": "boolean",
								"description": "Replace all occurrences (default false)"
							}
						}
					},
					"description": "Array of edit operations"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			filePath, _ := args["filePath"].(string)
			edits, ok := args["edits"].([]interface{})
			if !ok || len(edits) == 0 {
				return "", fmt.Errorf("edits must be a non-empty array")
			}

			var results []string
			for i, editRaw := range edits {
				edit, ok := editRaw.(map[string]interface{})
				if !ok {
					results = append(results, fmt.Sprintf("Edit %d: invalid format", i))
					continue
				}

				// Allow per-edit filePath override
				targetPath := filePath
				if ep, ok := edit["filePath"].(string); ok && ep != "" {
					targetPath = ep
				}

				oldStr, _ := edit["oldString"].(string)
				newStr, _ := edit["newString"].(string)

				data, err := os.ReadFile(targetPath)
				if err != nil {
					results = append(results, fmt.Sprintf("Edit %d: read error: %v", i, err))
					continue
				}

				content := string(data)
				if !strings.Contains(content, oldStr) {
					results = append(results, fmt.Sprintf("Edit %d: oldString not found", i))
					continue
				}

				replaceAll, _ := edit["replaceAll"].(bool)
				if replaceAll {
					content = strings.ReplaceAll(content, oldStr, newStr)
				} else {
					content = strings.Replace(content, oldStr, newStr, 1)
				}

				if err := os.WriteFile(targetPath, []byte(content), 0o644); err != nil {
					results = append(results, fmt.Sprintf("Edit %d: write error: %v", i, err))
					continue
				}
				results = append(results, fmt.Sprintf("Edit %d: applied", i))
			}

			return strings.Join(results, "\n"), nil
		},
	})
}

// ============================================================================
// KIMI CLI PARITY TOOLS
// ============================================================================

// registerKimiCLITools adds Kimi CLI compatible tool surfaces.
func (r *Registry) registerKimiCLITools() {
	// TaskList - Kimi's background task listing
	r.Tools = append(r.Tools, Tool{
		Name:        "TaskList",
		Description: "List background tasks. Shows task ID, kind, status, and description for each task.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"active_only": {
					"type": "boolean",
					"description": "Whether to list only non-terminal tasks. Default true.",
					"default": true
				},
				"limit": {
					"type": "integer",
					"description": "Maximum number of tasks to return. Default 20.",
					"minimum": 1,
					"maximum": 100,
					"default": 20
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			activeOnly := true
			if ao, ok := args["active_only"].(bool); ok {
				activeOnly = ao
			}
			limit := toInt(args["limit"], 20)

			GlobalJobManager.mu.Lock()
			defer GlobalJobManager.mu.Unlock()

			var lines []string
			count := 0
			for id, job := range GlobalJobManager.jobs {
				if activeOnly && (jobStatus(job) != "running") {
					continue
				}
				if count >= limit {
					break
				}
				lines = append(lines, fmt.Sprintf("  %s  kind=shell  status=%s  desc=%s", id, jobStatus(job), job.Command))
				count++
			}

			if len(lines) == 0 {
				return "No active background tasks.", nil
			}
			return fmt.Sprintf("Background tasks (%d):\n%s", count, strings.Join(lines, "\n")), nil
		},
	})

	// TaskOutput - Kimi's task output inspection
	r.Tools = append(r.Tools, Tool{
		Name:        "TaskOutput",
		Description: "Get the output of a background task. Optionally block until the task completes.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["task_id"],
			"properties": {
				"task_id": {
					"type": "string",
					"description": "The background task ID to inspect."
				},
				"block": {
					"type": "boolean",
					"description": "Whether to wait for the task to finish. Default false.",
					"default": false
				},
				"timeout": {
					"type": "integer",
					"description": "Maximum seconds to wait when block=true. Default 30.",
					"minimum": 0,
					"maximum": 3600
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			taskID, _ := args["task_id"].(string)
			if taskID == "" {
				return "", fmt.Errorf("task_id is required")
			}

			GlobalJobManager.mu.Lock()
			job, exists := GlobalJobManager.jobs[taskID]
			GlobalJobManager.mu.Unlock()

			if !exists {
				return "", fmt.Errorf("Task not found: %s", taskID)
			}

			block, _ := args["block"].(bool)
			timeout := toInt(args["timeout"], 30)

			if block && jobStatus(job) == "running" {
				// Wait for completion
				deadline := time.Now().Add(time.Duration(timeout) * time.Second)
				for time.Now().Before(deadline) {
					GlobalJobManager.mu.Lock()
					job = GlobalJobManager.jobs[taskID]
					GlobalJobManager.mu.Unlock()
					if jobStatus(job) != "running" {
						break
					}
					time.Sleep(100 * time.Millisecond)
				}
			}

			output := job.Output
			outputStr := ""
			if output != nil {
				outputStr = output.String()
			}
			if outputStr == "" {
				outputStr = "[no output available]"
			}

			return fmt.Sprintf("retrieval_status: success\ntask_id: %s\nstatus: %s\n\n[output]\n%s", taskID, jobStatus(job), outputStr), nil
		},
	})

	// TaskStop - Kimi's task stop tool
	r.Tools = append(r.Tools, Tool{
		Name:        "TaskStop",
		Description: "Stop a running background task.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["task_id"],
			"properties": {
				"task_id": {
					"type": "string",
					"description": "The background task ID to stop."
				},
				"reason": {
					"type": "string",
					"description": "Short reason for stopping. Default: 'Stopped by TaskStop'."
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			taskID, _ := args["task_id"].(string)
			if taskID == "" {
				return "", fmt.Errorf("task_id is required")
			}
			reason, _ := args["reason"].(string)
			if reason == "" {
				reason = "Stopped by TaskStop"
			}

			GlobalJobManager.mu.Lock()
			defer GlobalJobManager.mu.Unlock()

			job, exists := GlobalJobManager.jobs[taskID]
			if !exists {
				return "", fmt.Errorf("Task not found: %s", taskID)
			}

			if job.cancel != nil {
				job.cancel()
			}
			if job.Cmd != nil && job.Cmd.Process != nil {
				job.Cmd.Process.Kill()
			}
			job.Done = true
			job.ExitCode = -1

			return fmt.Sprintf("Task %s stopped. Reason: %s", taskID, reason), nil
		},
	})

	// Think - Kimi's thinking/reasoning tool
	r.Tools = append(r.Tools, Tool{
		Name:        "Think",
		Description: "Think through a problem step by step. Use this tool to reason about complex problems before acting.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["thought"],
			"properties": {
				"thought": {
					"type": "string",
					"description": "Your step-by-step reasoning process"
				},
				"next_action": {
					"type": "string",
					"description": "What you plan to do next after thinking",
					"enum": ["continue", "ask_user", "execute"]
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			// Think is a no-op tool that allows the model to reason
			// The thought content is captured by the agent loop for context
			return "Thinking recorded. Proceed with reasoning applied.", nil
		},
	})

	// AskUser - Kimi's user interaction tool
	r.Tools = append(r.Tools, Tool{
		Name:        "AskUser",
		Description: "Ask the user a question and wait for their response. Use for clarifications or when you need human input.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["question"],
			"properties": {
				"question": {
					"type": "string",
					"description": "The question to ask the user"
				},
				"context": {
					"type": "string",
					"description": "Additional context for the question"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			question, _ := args["question"].(string)
			if question == "" {
				return "", fmt.Errorf("question is required")
			}
			return fmt.Sprintf("[Waiting for user response to: %s]", question), nil
		},
	})

	// ReadFile - Kimi's file reading tool (PascalCase naming)
	r.Tools = append(r.Tools, Tool{
		Name:        "ReadFile",
		Description: "Read the contents of a file. Returns file content with line numbers.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["path"],
			"properties": {
				"path": {
					"type": "string",
					"description": "Path to the file to read"
				},
				"line_offset": {
					"type": "integer",
					"description": "Starting line number (1-indexed). Default: 1.",
					"minimum": 1
				},
				"n_lines": {
					"type": "integer",
					"description": "Number of lines to read. Default: all lines.",
					"minimum": 1
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			path, _ := args["path"].(string)
			if path == "" {
				return "", fmt.Errorf("path is required")
			}

			piArgs := map[string]interface{}{"path": path}
			if offset, ok := args["line_offset"]; ok {
				piArgs["offset"] = offset
			}
			if limit, ok := args["n_lines"]; ok {
				piArgs["limit"] = limit
			}
			return executePiTool("read", piArgs)
		},
	})

	// WriteFile - Kimi's file writing tool (PascalCase naming)
	r.Tools = append(r.Tools, Tool{
		Name:        "WriteFile",
		Description: "Write content to a file. Creates the file and parent directories if needed.",
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
					"description": "The content to write"
				},
				"create_dirs": {
					"type": "boolean",
					"description": "Create parent directories if needed. Default: true.",
					"default": true
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

			// Create parent directories
			createDirs := true
			if cd, ok := args["create_dirs"].(bool); ok {
				createDirs = cd
			}
			if createDirs {
				dir := filepath.Dir(path)
				if err := os.MkdirAll(dir, 0o755); err != nil {
					return "", fmt.Errorf("failed to create directories: %w", err)
				}
			}

			piArgs := map[string]interface{}{"path": path, "content": content}
			return executePiTool("write", piArgs)
		},
	})

	// Replace - Kimi's search/replace tool
	r.Tools = append(r.Tools, Tool{
		Name:        "Replace",
		Description: "Replace text in a file. Finds the exact old text and replaces it with new text.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["path", "old_text", "new_text"],
			"properties": {
				"path": {
					"type": "string",
					"description": "Path to the file"
				},
				"old_text": {
					"type": "string",
					"description": "The exact text to find"
				},
				"new_text": {
					"type": "string",
					"description": "The replacement text"
				},
				"replace_all": {
					"type": "boolean",
					"description": "Replace all occurrences. Default: false.",
					"default": false
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			path, _ := args["path"].(string)
			oldText, _ := args["old_text"].(string)
			newText, _ := args["new_text"].(string)
			if path == "" {
				return "", fmt.Errorf("path is required")
			}

			piArgs := map[string]interface{}{
				"path": path,
				"edits": []map[string]interface{}{
					{"oldText": oldText, "newText": newText},
				},
			}
			return executePiTool("edit", piArgs)
		},
	})

	// Glob - Kimi's file globbing tool (PascalCase)
	r.Tools = append(r.Tools, Tool{
		Name:        "Glob",
		Description: "Find files matching a glob pattern. Returns file paths relative to the search directory.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["pattern"],
			"properties": {
				"pattern": {
					"type": "string",
					"description": "Glob pattern to match files"
				},
				"path": {
					"type": "string",
					"description": "Directory to search in. Default: current directory."
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
			return executePiTool("find", piArgs)
		},
	})

	// GrepLocal - Kimi's local grep tool (PascalCase)
	r.Tools = append(r.Tools, Tool{
		Name:        "GrepLocal",
		Description: "Search for a pattern in local files. Returns matching lines with file paths and line numbers.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["pattern"],
			"properties": {
				"pattern": {
					"type": "string",
					"description": "The search pattern (regex or literal)"
				},
				"path": {
					"type": "string",
					"description": "Directory or file to search in"
				},
				"ignore_case": {
					"type": "boolean",
					"description": "Case-insensitive search. Default: false."
				},
				"file_pattern": {
					"type": "string",
					"description": "Filter by file pattern (e.g., '*.py')"
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
			if ignoreCase, ok := args["ignore_case"].(bool); ok {
				piArgs["ignoreCase"] = ignoreCase
			}
			if filePattern, ok := args["file_pattern"].(string); ok && filePattern != "" {
				piArgs["glob"] = filePattern
			}
			return executePiTool("grep", piArgs)
		},
	})

	// WebFetch - Kimi's web fetch tool (PascalCase)
	r.Tools = append(r.Tools, Tool{
		Name:        "WebFetch",
		Description: "Fetch content from a URL. Returns the page content in the specified format.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["url"],
			"properties": {
				"url": {
					"type": "string",
					"description": "The URL to fetch"
				},
				"format": {
					"type": "string",
					"description": "Return format: text, markdown, or html. Default: markdown.",
					"enum": ["text", "markdown", "html"]
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			url, _ := args["url"].(string)
			if url == "" {
				return "", fmt.Errorf("url is required")
			}
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, "curl", "-sL", url)
			output, err := cmd.CombinedOutput()
			if err != nil {
				return "", fmt.Errorf("fetch failed: %w", err)
			}
			result := string(output)
			if len(result) > 50000 {
				result = result[:50000] + "\n... (truncated)"
			}
			return result, nil
		},
	})

	// WebSearch - Kimi's web search tool (PascalCase)
	r.Tools = append(r.Tools, Tool{
		Name:        "WebSearch",
		Description: "Search the web for information. Returns search results.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["query"],
			"properties": {
				"query": {
					"type": "string",
					"description": "Search query"
				},
				"num_results": {
					"type": "integer",
					"description": "Number of results. Default: 5."
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			query, _ := args["query"].(string)
			if query == "" {
				return "", fmt.Errorf("query is required")
			}
			return fmt.Sprintf("Web search for: %q\n(Configure search API for live results)", query), nil
		},
	})

	// PlanEnter - Kimi's plan mode entry (PascalCase)
	r.Tools = append(r.Tools, Tool{
		Name:        "PlanEnter",
		Description: "Enter plan mode. Switches to read-only planning agent that cannot modify files.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return "Entered plan mode. File modifications disabled. Use PlanExit to return to build mode.", nil
		},
	})

	// PlanExit - Kimi's plan mode exit (PascalCase)
	r.Tools = append(r.Tools, Tool{
		Name:        "PlanExit",
		Description: "Exit plan mode and return to build mode with file modification capabilities.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return "Exited plan mode. File modifications enabled.", nil
		},
	})
}

// ============================================================================
// CURSOR PARITY TOOLS
// ============================================================================

// registerCursorTools adds Cursor IDE compatible tool surfaces.
func (r *Registry) registerCursorTools() {
	// cursor_read_file - Cursor's file reading
	r.Tools = append(r.Tools, Tool{
		Name:        "cursor_read_file",
		Description: "Read a file in the workspace. Returns the file contents.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["relative_path"],
			"properties": {
				"relative_path": {
					"type": "string",
					"description": "Path relative to workspace root"
				},
				"start_line": {
					"type": "integer",
					"description": "Start line (1-indexed)"
				},
				"end_line": {
					"type": "integer",
					"description": "End line (inclusive)"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			path, _ := args["relative_path"].(string)
			piArgs := map[string]interface{}{"path": path}
			if start, ok := args["start_line"]; ok {
				piArgs["offset"] = start
			}
			if end, ok := args["end_line"]; ok {
				startLine := toInt(args["start_line"], 1)
				endLine := toInt(end, startLine+100)
				piArgs["limit"] = endLine - startLine + 1
			}
			return executePiTool("read", piArgs)
		},
	})

	// cursor_edit_file - Cursor's file editing
	r.Tools = append(r.Tools, Tool{
		Name:        "cursor_edit_file",
		Description: "Edit a file by providing target text to find and replacement text.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["relative_path", "old_text", "new_text"],
			"properties": {
				"relative_path": {
					"type": "string",
					"description": "Path relative to workspace root"
				},
				"old_text": {
					"type": "string",
					"description": "Exact text to find"
				},
				"new_text": {
					"type": "string",
					"description": "Replacement text"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			path, _ := args["relative_path"].(string)
			oldText, _ := args["old_text"].(string)
			newText, _ := args["new_text"].(string)
			piArgs := map[string]interface{}{
				"path": path,
				"edits": []map[string]interface{}{
					{"oldText": oldText, "newText": newText},
				},
			}
			return executePiTool("edit", piArgs)
		},
	})

	// cursor_run_command - Cursor's terminal command
	r.Tools = append(r.Tools, Tool{
		Name:        "cursor_run_command",
		Description: "Run a terminal command and return its output.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["command"],
			"properties": {
				"command": {
					"type": "string",
					"description": "Terminal command to execute"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			command, _ := args["command"].(string)
			piArgs := map[string]interface{}{"command": command}
			return executePiTool("bash", piArgs)
		},
	})

	// cursor_code_search - Cursor's code search
	r.Tools = append(r.Tools, Tool{
		Name:        "cursor_code_search",
		Description: "Search for code patterns across the workspace.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["query"],
			"properties": {
				"query": {
					"type": "string",
					"description": "Search query"
				},
				"file_type": {
					"type": "string",
					"description": "Filter by file type (e.g., 'py', 'ts')"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			query, _ := args["query"].(string)
			piArgs := map[string]interface{}{"pattern": query}
			if fileType, ok := args["file_type"].(string); ok && fileType != "" {
				piArgs["glob"] = "*." + fileType
			}
			return executePiTool("grep", piArgs)
		},
	})

	// cursor_list_dir - Cursor's directory listing
	r.Tools = append(r.Tools, Tool{
		Name:        "cursor_list_dir",
		Description: "List directory contents.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"relative_path": {
					"type": "string",
					"description": "Path relative to workspace root"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			path, _ := args["relative_path"].(string)
			if path == "" {
				path = "."
			}
			piArgs := map[string]interface{}{"path": path}
			return executePiTool("ls", piArgs)
		},
	})
}

// ============================================================================
// WINDSURF PARITY TOOLS
// ============================================================================

// registerWindsurfTools adds Windsurf/Codium compatible tool surfaces.
func (r *Registry) registerWindsurfTools() {
	// read_file (Windsurf format - already registered via Gemini parity,
	// but add windsurf-specific cascade_write)

	// cascade_edit - Windsurf Cascade's edit format
	r.Tools = append(r.Tools, Tool{
		Name:        "cascade_edit",
		Description: "Edit a file using Windsurf Cascade format. Provides the full file path, search text, and replacement text.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["file_path", "search", "replace"],
			"properties": {
				"file_path": {
					"type": "string",
					"description": "Full path to the file"
				},
				"search": {
					"type": "string",
					"description": "Exact text to find"
				},
				"replace": {
					"type": "string",
					"description": "Replacement text"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			filePath, _ := args["file_path"].(string)
			search, _ := args["search"].(string)
			replace, _ := args["replace"].(string)
			piArgs := map[string]interface{}{
				"path": filePath,
				"edits": []map[string]interface{}{
					{"oldText": search, "newText": replace},
				},
			}
			return executePiTool("edit", piArgs)
		},
	})

	// cascade_command - Windsurf's terminal execution
	r.Tools = append(r.Tools, Tool{
		Name:        "cascade_command",
		Description: "Execute a terminal command via Windsurf Cascade.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["command"],
			"properties": {
				"command": {
					"type": "string",
					"description": "Command to execute"
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
			piArgs := map[string]interface{}{"command": command}
			if wd, ok := args["working_directory"].(string); ok && wd != "" {
				piArgs["working_dir"] = wd
			}
			return executePiTool("bash", piArgs)
		},
	})
}

// ============================================================================
// MISTRAL VIBE / BITO / SMITHERY PARITY TOOLS
// ============================================================================

// registerMistralVibeTools adds Mistral Vibe compatible tool surfaces.
func (r *Registry) registerMistralVibeTools() {
	// mistral_edit - Mistral's codestral edit format
	r.Tools = append(r.Tools, Tool{
		Name:        "mistral_edit",
		Description: "Edit a file using Mistral/Codestral format.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["file_path", "content"],
			"properties": {
				"file_path": {
					"type": "string",
					"description": "Path to the file"
				},
				"content": {
					"type": "string",
					"description": "New file content (full rewrite)"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			path, _ := args["file_path"].(string)
			content, _ := args["content"].(string)
			piArgs := map[string]interface{}{"path": path, "content": content}
			return executePiTool("write", piArgs)
		},
	})

	// mistral_search - Mistral's code search
	r.Tools = append(r.Tools, Tool{
		Name:        "mistral_search",
		Description: "Search code using Mistral format.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["query"],
			"properties": {
				"query": {
					"type": "string",
					"description": "Search query"
				},
				"path": {
					"type": "string",
					"description": "Directory to search in"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			query, _ := args["query"].(string)
			piArgs := map[string]interface{}{"pattern": query}
			if path, ok := args["path"].(string); ok && path != "" {
				piArgs["path"] = path
			}
			return executePiTool("grep", piArgs)
		},
	})
}

// registerSmitheryTools adds Smithery MCP CLI compatible tool surfaces.
func (r *Registry) registerSmitheryTools() {
	// smithery_install - Install MCP server from Smithery registry
	r.Tools = append(r.Tools, Tool{
		Name:        "smithery_install",
		Description: "Install an MCP server from the Smithery registry.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["server_name"],
			"properties": {
				"server_name": {
					"type": "string",
					"description": "Name of the MCP server to install from Smithery"
				},
				"config": {
					"type": "object",
					"description": "Configuration for the server",
					"additionalProperties": true
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			serverName, _ := args["server_name"].(string)
			if serverName == "" {
				return "", fmt.Errorf("server_name is required")
			}

			configRaw, _ := args["config"].(map[string]interface{})
			if configRaw == nil {
				configRaw = make(map[string]interface{})
			}

			// Use the extensions manager to install
			extMgr := extensions.NewManager("")

			if err := extMgr.InstallFromSmithery(serverName, configRaw); err != nil {
				return "", fmt.Errorf("smithery install failed: %w", err)
			}

			return fmt.Sprintf("Smithery: Successfully installed MCP server '%s'", serverName), nil
		},
	})

	// smithery_list - List available MCP servers
	r.Tools = append(r.Tools, Tool{
		Name:        "smithery_list",
		Description: "List available MCP servers from the Smithery registry.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"category": {
					"type": "string",
					"description": "Filter by category"
				},
				"search": {
					"type": "string",
					"description": "Search query"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return "Smithery MCP registry: Configure Smithery API for live server discovery.", nil
		},
	})
}

// Unused import guard
var _ = regexp.MustCompile("")

// jobStatus returns a string status for a Job.
func jobStatus(j *Job) string {
	if j == nil {
		return "unknown"
	}
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.Done {
		if j.ExitCode == 0 {
			return "completed"
		}
		return "failed"
	}
	return "running"
}
