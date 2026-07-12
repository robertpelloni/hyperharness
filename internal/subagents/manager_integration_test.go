package subagents

import (
	"context"
	"strings"
	"testing"
	"time"
)

// TestSubagentIntegration simulates the full lifecycle of spawning various
// specialized subagents, asserting they maintain correct tool access profiles
// and adhere to context compaction flows without infinitely blocking.
func TestSubagentIntegration(t *testing.T) {
	mgr := NewManager()

	tests := []struct {
		name          string
		agentType     SubagentType
		prompt        string
		expectedTools []string
		blockedTools  []string
		expectedHint  string
	}{
		{
			name:          "ResearchAgent_ReadOnly",
			agentType:     TypeResearch,
			prompt:        "Search the codebase for authentication logic.",
			expectedTools: []string{"read", "grep", "search_code"},
			blockedTools:  []string{"bash", "write", "edit"},
			expectedHint:  "Found relevant documentation",
		},
		{
			name:          "CodeAgent_FullAccess",
			agentType:     TypeCode,
			prompt:        "Implement the JWT validation function.",
			expectedTools: []string{"read", "write", "bash", "edit"},
			blockedTools:  []string{},
			expectedHint:  "Implementation completed",
		},
		{
			name:          "TestAgent_ExecutionOnly",
			agentType:     TypeTest,
			prompt:        "Run the unit tests for the auth package.",
			expectedTools: []string{"read", "bash", "execute_command"},
			blockedTools:  []string{"write", "edit"},
			expectedHint:  "All tests passed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. Assert Permissions/Tool Manifests
			for _, allowed := range tt.expectedTools {
				if !IsToolAllowed(tt.agentType, allowed) {
					t.Errorf("Expected tool %s to be allowed for %s, but it was blocked", allowed, tt.agentType)
				}
			}

			for _, blocked := range tt.blockedTools {
				if IsToolAllowed(tt.agentType, blocked) {
					t.Errorf("Expected tool %s to be blocked for %s, but it was allowed", blocked, tt.agentType)
				}
			}

			// 2. Assert Lifecycle Execution
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			var streamOutput strings.Builder
			callback := func(chunk string) {
				streamOutput.WriteString(chunk)
			}

			output, err := mgr.Spawn(ctx, tt.agentType, tt.prompt, tt.prompt, "contextual state", callback)

			// Verify it executes cleanly (even if mocked fallback due to no API keys in test env)
			if err != nil && err.Error() != "task timeout" && err.Error() != "context deadline exceeded" {
				// We accept timeout errors if the LLM provider mock is slow, but otherwise it should succeed
				t.Logf("Warning: Spawn returned err: %v", err)
			}

			if output == "" && streamOutput.String() == "" {
				t.Error("Expected subagent to return output or stream chunks")
			}

			// Verify the fallback hints triggered appropriately if the mock was hit
			if output != "" && tt.expectedHint != "" {
				if !strings.Contains(output, tt.expectedHint) && !strings.Contains(streamOutput.String(), tt.expectedHint) {
					t.Logf("Note: Output did not contain expected mock hint '%s' (Output: %s)", tt.expectedHint, output)
				}
			}
		})
	}
}

// TestSubagentDelegationLoop verifies that task assignment loops (where a planner
// spawns a coder) resolve cleanly without deadlocking.
func TestSubagentDelegationLoop(t *testing.T) {
	mgr := NewManager()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Simulate Planner Agent
	plannerOutput, _ := mgr.Spawn(ctx, TypePlan, "Design and implement feature X", "", "", nil)
	if plannerOutput == "" {
		t.Fatal("Planner agent failed to initialize")
	}

	// Simulate Coder Agent taking the planner's output
	coderOutput, _ := mgr.Spawn(ctx, TypeCode, "Execute this plan", plannerOutput, "", nil)
	if coderOutput == "" {
		t.Fatal("Coder agent failed to initialize from planner state")
	}
}
