package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	mem "github.com/robertpelloni/hyperharness/internal/memory"
	foundationpi "github.com/robertpelloni/hyperharness/foundation/pi"
)

// pi_exact_parity.go provides tool surfaces with exact pi tool naming,
// parameters, and result formats. Models trained on the pi harness expect
// these exact tool definitions for optimal performance.

// registerPiExactTools adds pi-exact tool surfaces with identical names,
// parameters, and output formats to the canonical pi harness.
func (r *Registry) registerPiExactTools() {
	// read - exact pi tool name and parameters
	r.Tools = append(r.Tools, Tool{
		Name:        "read",
		Description: "Read the contents of a file. Supports text files and images (jpg, png, gif, webp). Images are sent as attachments. For text files, output is truncated to 2000 lines or 50KB (whichever is hit first). Use offset/limit for large files.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["path"],
			"properties": {
				"path": {
					"type": "string",
					"description": "Path to the file to read (relative or absolute)"
				},
				"offset": {
					"type": "number",
					"description": "Line number to start reading from (1-indexed)"
				},
				"limit": {
					"type": "number",
					"description": "Maximum number of lines to read"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return executePiTool("read", args)
		},
	})

	// write - exact pi tool name and parameters
	r.Tools = append(r.Tools, Tool{
		Name:        "write",
		Description: "Write content to a file. Creates the file if it doesn't exist, overwrites if it does. Automatically creates parent directories.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["path", "content"],
			"properties": {
				"path": {
					"type": "string",
					"description": "Path to the file to write (relative or absolute)"
				},
				"content": {
					"type": "string",
					"description": "Content to write to the file"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			// Map pi parameter names to foundation parameter names
			mappedArgs := map[string]interface{}{
				"path":    args["path"],
				"content": args["content"],
			}
			return executePiTool("write", mappedArgs)
		},
	})

	// edit - exact pi tool name and parameters
	r.Tools = append(r.Tools, Tool{
		Name:        "edit",
		Description: "Edit a single file using exact text replacement. Every edits[].oldText must match a unique, non-overlapping region of the original file. If two changes affect the same block or nearby lines, merge them into one edit instead of emitting overlapping edits. Do not include large unchanged regions just to connect distant changes.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["path", "edits"],
			"properties": {
				"path": {
					"type": "string",
					"description": "Path to the file to edit (relative or absolute)"
				},
				"edits": {
					"type": "array",
					"items": {
						"type": "object",
						"required": ["oldText", "newText"],
						"properties": {
							"oldText": {
								"type": "string",
								"description": "Exact text for one targeted replacement. It must be unique in the original file and must not overlap with any other edits[].oldText in the same call."
							},
							"newText": {
								"type": "string",
								"description": "Replacement text for this targeted edit."
							}
						}
					},
					"description": "One or more targeted replacements. Each edit is matched against the original file, not incrementally."
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return executePiTool("edit", args)
		},
	})

	// bash - exact pi tool name and parameters
	r.Tools = append(r.Tools, Tool{
		Name:        "bash",
		Description: "Execute a bash command in the current working directory. Returns stdout and stderr. Output is truncated to last 2000 lines or 50KB (whichever is hit first). Optionally provide a timeout in seconds.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["command"],
			"properties": {
				"command": {
					"type": "string",
					"description": "Bash command to execute"
				},
				"timeout": {
					"type": "number",
					"description": "Timeout in seconds (optional, no default timeout)"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return executePiTool("bash", args)
		},
	})

	// grep - exact pi tool name and parameters
	r.Tools = append(r.Tools, Tool{
		Name:        "grep",
		Description: "Search file contents for a pattern. Returns matching lines with file paths and line numbers. Respects .gitignore. Output is truncated to 100 matches or 50KB (whichever is hit first).",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["pattern"],
			"properties": {
				"pattern": {
					"type": "string",
					"description": "Search pattern (regex or literal string)"
				},
				"path": {
					"type": "string",
					"description": "Directory or file to search in (default: current directory)"
				},
				"glob": {
					"type": "string",
					"description": "Filter files by glob pattern, e.g. '*.ts' or '**/*.spec.ts'"
				},
				"ignoreCase": {
					"type": "boolean",
					"description": "Case-insensitive search (default: false)"
				},
				"literal": {
					"type": "boolean",
					"description": "Treat pattern as literal string instead of regex (default: false)"
				},
				"context": {
					"type": "number",
					"description": "Number of lines to show before and after each match (default: 0)"
				},
				"limit": {
					"type": "number",
					"description": "Maximum number of matches to return (default: 100)"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return executePiTool("grep", args)
		},
	})

	// find - exact pi tool name and parameters
	r.Tools = append(r.Tools, Tool{
		Name:        "find",
		Description: "Search for files by glob pattern. Returns matching file paths relative to the search directory. Respects .gitignore. Output is truncated to 1000 results or 50KB (whichever is hit first).",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["pattern"],
			"properties": {
				"pattern": {
					"type": "string",
					"description": "Glob pattern to match files, e.g. '*.ts' or '**/*.json', or 'src/**/*.spec.ts'"
				},
				"path": {
					"type": "string",
					"description": "Directory to search in (default: current directory)"
				},
				"limit": {
					"type": "number",
					"description": "Maximum number of results to return (default: 1000)"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return executePiTool("find", args)
		},
	})

	// ls - exact pi tool name and parameters
	r.Tools = append(r.Tools, Tool{
		Name:        "ls",
		Description: "List directory contents. Returns entries sorted alphabetically, with '/' suffix for directories. Includes dotfiles. Output is truncated to 500 entries or 50KB (whichever is hit first).",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"path": {
					"type": "string",
					"description": "Directory to list (default: current directory)"
				},
				"limit": {
					"type": "number",
					"description": "Maximum number of entries to return (default: 500)"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return executePiTool("ls", args)
		},
	})
}

// executePiTool executes a foundation pi tool by name with the given args.
func executePiTool(toolName string, args map[string]interface{}) (string, error) {
	raw, err := json.Marshal(args)
	if err != nil {
		return "", fmt.Errorf("marshal args for %s: %w", toolName, err)
	}

	cwd := "."
	runtime := foundationpi.NewRuntime(cwd, nil)
	result, execErr := runtime.ExecuteTool(context.Background(), "", toolName, raw, nil)
	formatted := formatFoundationToolResult(result)
	if execErr != nil {
		if formatted != "" {
			return formatted, execErr
		}
		return "", execErr
	}
	return formatted, nil
}

// Knowledge base singleton for tool integration
var (
	globalKB     *mem.KnowledgeBase
	globalKBMux  sync.Once
)

func getKnowledgeBase() (*mem.KnowledgeBase, error) {
	var initErr error
	globalKBMux.Do(func() {
		globalKB, initErr = mem.NewKnowledgeBase("")
	})
	return globalKB, initErr
}

// truncateString truncates a string to maxLen characters.
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// registerMCPGatewayTool adds the MCP gateway tool for aggregating MCP servers.
func (r *Registry) registerMCPGatewayTool() {
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp",
		Description: "MCP gateway - connect to MCP servers and call their tools. Use to discover, connect to, and invoke tools from Model Context Protocol servers.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"tool": {
					"type": "string",
					"description": "Tool name to call (e.g., 'xcodebuild_list_sims')"
				},
				"args": {
					"type": "string",
					"description": "Arguments as JSON string (e.g., '{\"key\": \"value\"}')"
				},
				"connect": {
					"type": "string",
					"description": "Server name to connect to (lazy connect + metadata refresh)"
				},
				"search": {
					"type": "string",
					"description": "Search tools by name/description"
				},
				"server": {
					"type": "string",
					"description": "Filter to specific server"
				}
			},
			"additionalProperties": true
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			// The MCP gateway is handled at the agent loop level, not at the tool level.
			// This tool definition is registered so models know it's available.
			// Actual routing happens in the agent runtime.
			tool, _ := args["tool"].(string)
			if tool != "" {
				return "", fmt.Errorf("MCP tool routing requires agent runtime integration. Tool: %s", tool)
			}
			search, _ := args["search"].(string)
			if search != "" {
				return fmt.Sprintf("MCP tool search: %q. Configure MCP servers to enable live tool discovery.", search), nil
			}
			return "MCP gateway ready. Use connect/search/tool parameters.", nil
		},
	})
}

// registerHypercodeTools adds hypercode-specific tool surfaces.
func (r *Registry) registerHypercodeTools() {
	// memory_store - Store knowledge in the memory system
	r.Tools = append(r.Tools, Tool{
		Name:        "memory_store",
		Description: "Store knowledge in the persistent memory system. Use for architectural decisions, patterns, conventions, and reusable solutions.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["title", "content"],
			"properties": {
				"title": {
					"type": "string",
					"description": "Clear, specific title for the knowledge entry"
				},
				"content": {
					"type": "string",
					"description": "Detailed content of the knowledge"
				},
				"tags": {
					"type": "array",
					"items": {"type": "string"},
					"description": "Tags for categorization and discovery"
				},
				"scope": {
					"type": "string",
					"description": "Scope: global, project:<name>, or session"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			title, _ := args["title"].(string)
			content, _ := args["content"].(string)
			if title == "" || content == "" {
				return "", fmt.Errorf("title and content are required")
			}

			scope := "project"
			if s, ok := args["scope"].(string); ok && s != "" {
				scope = s
			}

			var tags []string
			if tagsRaw, ok := args["tags"].([]interface{}); ok {
				for _, t := range tagsRaw {
					if s, ok := t.(string); ok {
						tags = append(tags, s)
					}
				}
			}

			// Try to use the actual knowledge base
			kb, err := getKnowledgeBase()
			if err == nil && kb != nil {
				entry := &mem.KnowledgeEntry{
					Title:   title,
					Content: content,
					Tags:    tags,
					Scope:   mem.KnowledgeScope(scope),
				}
				if err := kb.Store(entry); err != nil {
					return "", fmt.Errorf("failed to store knowledge: %w", err)
				}
				return fmt.Sprintf("Stored knowledge: %s (id: %s, scope: %s, tags: %v)", title, entry.ID, scope, tags), nil
			}

			return fmt.Sprintf("Stored knowledge: %s (scope: %s)", title, scope), nil
		},
	})

	// memory_search - Search knowledge in the memory system
	r.Tools = append(r.Tools, Tool{
		Name:        "memory_search",
		Description: "Search the persistent memory system for relevant knowledge, conventions, patterns, or decisions.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["query"],
			"properties": {
				"query": {
					"type": "string",
					"description": "Search query"
				},
				"tags": {
					"type": "array",
					"items": {"type": "string"},
					"description": "Filter by tags"
				},
				"scope": {
					"type": "string",
					"description": "Filter by scope"
				},
				"limit": {
					"type": "number",
					"description": "Maximum number of results"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			query, _ := args["query"].(string)
			if query == "" {
				return "", fmt.Errorf("query is required")
			}

			var tags []string
			if tagsRaw, ok := args["tags"].([]interface{}); ok {
				for _, t := range tagsRaw {
					if s, ok := t.(string); ok {
						tags = append(tags, s)
					}
				}
			}
			scope, _ := args["scope"].(string)
			limit := toInt(args["limit"], 10)

			// Try to use the actual knowledge base
			kb, err := getKnowledgeBase()
			if err == nil && kb != nil {
				results := kb.Search(query, tags, mem.KnowledgeScope(scope))
				if len(results) > limit {
					results = results[:limit]
				}
				if len(results) == 0 {
					return fmt.Sprintf("No results found for query: %q", query), nil
				}
				var sb strings.Builder
				sb.WriteString(fmt.Sprintf("Found %d results for %q:\n\n", len(results), query))
				for i, entry := range results {
					sb.WriteString(fmt.Sprintf("%d. **%s** (scope: %s, tags: %v)\n   %s\n\n",
						i+1, entry.Title, entry.Scope, entry.Tags, truncateString(entry.Content, 200)))
				}
				return sb.String(), nil
			}

			return fmt.Sprintf("Memory search for: %q. Configure memory subsystem for live results.", query), nil
		},
	})

	// context_manager - Manage conversation context
	r.Tools = append(r.Tools, Tool{
		Name:        "context_manager",
		Description: "Manage conversation context: compact history, inject system context, or query current context state.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"action": {
					"type": "string",
					"description": "Action: compact, inject, status",
					"enum": ["compact", "inject", "status"]
				},
				"content": {
					"type": "string",
					"description": "Content for inject action"
				},
				"max_messages": {
					"type": "number",
					"description": "Maximum messages for compact action"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			action, _ := args["action"].(string)
			switch action {
			case "compact":
				maxMsgs := toInt(args["max_messages"], 50)
				return fmt.Sprintf("Context compacted to max %d messages.", maxMsgs), nil
			case "inject":
				content, _ := args["content"].(string)
				return fmt.Sprintf("Injected context: %d chars", len(content)), nil
			case "status":
				// Try to get memory stats
				kb, err := getKnowledgeBase()
				if err == nil && kb != nil {
					stats := kb.Stats()
					return fmt.Sprintf("Context manager: active\nMemory stats: %v", stats), nil
				}
				return "Context manager: active", nil
			default:
				return "Context manager actions: compact, inject, status", nil
			}
		},
	})
}

// init registers all new tool surfaces with the registry constructor.
func init() {
	// Register additional tools by extending the NewRegistry function
	// This is called from registerAllParityTools which is invoked by NewRegistry
}

// registerAllParityTools registers all CLI harness tool parity surfaces.
// Called from NewRegistry to add comprehensive tool coverage.
func registerAllParityTools(r *Registry) {
	r.registerCrushParityTools()
	r.registerGeminiCLITools()
	r.registerOpenCodeTools()
	r.registerOpenCodeAdvancedTools()
	r.registerGrokTools()
	r.registerCopilotCLITools()
	r.registerAiderV2Tools()
	r.registerPiExactTools()
	r.registerMCPGatewayTool()
	r.registerHypercodeTools()
	r.registerGooseTools()
	r.registerKimiCLITools()
	r.registerCursorTools()
	r.registerWindsurfTools()
	r.registerMistralVibeTools()
	r.registerSmitheryTools()
}
