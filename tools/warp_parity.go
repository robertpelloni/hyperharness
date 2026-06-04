package tools

import (
	"encoding/json"
)

// registerWarpParityTools registers tool schemas matching Warp's AI Agent actions.
func (r *Registry) registerWarpParityTools() {
	// RequestCommandOutput - Parity with Warp's command execution
	r.Tools = append(r.Tools, Tool{
		Name:        "RequestCommandOutput",
		Description: "Warp: Execute a command and get its output as context.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"command": { "type": "string" },
				"is_read_only": { "type": "boolean" },
				"wait_until_completion": { "type": "boolean", "default": true }
			},
			"required": ["command"]
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return executePiTool("bash", args)
		},
	})

	// SearchCodebase - Parity with Warp's codebase search
	r.Tools = append(r.Tools, Tool{
		Name:        "SearchCodebase",
		Description: "Warp: Search the codebase for a query, optionally filtered by partial paths.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"query": { "type": "string" },
				"partial_paths": { "type": "array", "items": { "type": "string" } },
				"codebase_path": { "type": "string" }
			},
			"required": ["query"]
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			// Map to grep or internal search
			return executePiTool("grep", map[string]interface{}{
				"pattern": args["query"],
				"path":    args["codebase_path"],
			})
		},
	})

	// RunAgents - Warp's subagent orchestration
	r.Tools = append(r.Tools, Tool{
		Name:        "RunAgents",
		Description: "Warp: Orchestrate one or more child agents to perform tasks.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"summary": { "type": "string" },
				"base_prompt": { "type": "string" },
				"agent_run_configs": {
					"type": "array",
					"items": {
						"type": "object",
						"properties": {
							"name": { "type": "string" },
							"prompt": { "type": "string" }
						}
					}
				}
			},
			"required": ["summary", "agent_run_configs"]
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return "Warp: Orchestrating agents...", nil
		},
	})

	// FileGlob - Parity with Warp's file globbing
	r.Tools = append(r.Tools, Tool{
		Name:        "FileGlob",
		Description: "Warp: Find files matching glob patterns.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"patterns": { "type": "array", "items": { "type": "string" } },
				"path": { "type": "string" }
			},
			"required": ["patterns"]
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			patterns, _ := args["patterns"].([]interface{})
			pattern := "*"
			if len(patterns) > 0 {
				pattern, _ = patterns[0].(string)
			}
			return executePiTool("find", map[string]interface{}{
				"pattern": pattern,
				"path":    args["path"],
			})
		},
	})
}
