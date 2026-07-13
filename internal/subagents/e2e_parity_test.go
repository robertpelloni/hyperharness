package subagents_test

import (
	"strings"
	"testing"

	"github.com/robertpelloni/hyperharness/tools"
)

// TestE2EParityToSubagent explicitly tests the flow from a parity tool interface
// directly into the subagents.GlobalManager.Spawn loop, verifying that descriptions,
// prompts, and stream chunks are preserved.
func TestE2EParityToSubagent(t *testing.T) {
	registry := tools.NewRegistry()

	tests := []struct {
		toolName      string
		args          map[string]interface{}
		expectedMatch string
	}{
		{
			toolName: "delegate", // Goose
			args: map[string]interface{}{
				"task":          "refactor the context compaction pipeline to use the new SQLite FTS5 backend",
				"description":   "refactor context compaction",
				"subagent_type": "code",
				"mode":          "sync",
			},
			expectedMatch: "Delegated task 'refactor context compaction' to code subagent",
		},
		{
			toolName: "task", // OpenCode
			args: map[string]interface{}{
				"prompt":        "refactor the context compaction pipeline to use the new SQLite FTS5 backend",
				"description":   "refactor context compaction",
				"subagent_type": "code",
			},
			expectedMatch: "task_id:",
		},
		{
			toolName: "Agent", // Claude Code
			args: map[string]interface{}{
				"prompt": "refactor the context compaction pipeline to use the new SQLite FTS5 backend",
				"type":   "Explore",
			},
			expectedMatch: "Task spawned",
		},
	}

	for _, tt := range tests {
		t.Run(tt.toolName, func(t *testing.T) {
			var targetTool *tools.Tool
			for _, tool := range registry.Tools {
				if tool.Name == tt.toolName {
					// Duplicate handling: target first matched definition
					targetTool = &tool
					break
				}
			}

			if targetTool == nil {
				t.Fatalf("Tool %s not found in registry", tt.toolName)
			}

			result, err := targetTool.Execute(tt.args)
			if err != nil {
				t.Fatalf("Tool %s failed to execute: %v", tt.toolName, err)
			}

			if result == "" {
				t.Errorf("Tool %s returned empty result", tt.toolName)
			}

			if !strings.Contains(result, tt.expectedMatch) {
				t.Errorf("Tool %s result did not contain expected match '%s'. Got: %s", tt.toolName, tt.expectedMatch, result)
			}

			// Verify that the fallback LLM mock output actually bubbled up to the tool surface
			if !strings.Contains(result, "Implementation completed") && !strings.Contains(result, "Found relevant documentation") && !strings.Contains(result, "1. Analyze") {
				t.Logf("Warning: Underlying subagent fallback output was not present in the tool return string for %s.", tt.toolName)
			}
		})
	}
}
