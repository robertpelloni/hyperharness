package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/robertpelloni/hyperharness/internal/subagents"
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
			summary, _ := args["summary"].(string)
			basePrompt, _ := args["base_prompt"].(string)

			// We iterate through run configs and spawn tasks in sequence for Warp's orchestration.
			// Since we want to return a single result, we concatenate the outputs.
			resultText := fmt.Sprintf("Warp RunAgents Orchestration: %s\n", summary)

			configs, ok := args["agent_run_configs"].([]interface{})
			if !ok || len(configs) == 0 {
				return "", fmt.Errorf("no valid agent_run_configs provided")
			}

			for i, cfgRaw := range configs {
				cfg, ok := cfgRaw.(map[string]interface{})
				if !ok {
					continue
				}
				name, _ := cfg["name"].(string)
				prompt, _ := cfg["prompt"].(string)

				fullPrompt := fmt.Sprintf("%s\n%s", basePrompt, prompt)

				streamCallback := func(chunk string) {
					fmt.Print(chunk)
				}

				result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, fullPrompt, fullPrompt, "", streamCallback)
				if err != nil {
					resultText += fmt.Sprintf("\n[Agent: %s (Failed)]\nError: %v\n", name, err)
				} else {
					resultText += fmt.Sprintf("\n[Agent: %s (Completed)]\n%s\n", name, result)
				}

				if i < len(configs)-1 {
					resultText += "\n---\n"
				}
			}

			return resultText, nil
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
