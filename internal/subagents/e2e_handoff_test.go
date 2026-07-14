package subagents_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/robertpelloni/hyperharness/internal/subagents"
)

// TestE2ESessionHandoff validates complex multi-agent handoffs.
func TestE2ESessionHandoff(t *testing.T) {
	// e.g., Code -> Research -> Plan loop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Coder spawns task
	coderResult, err := subagents.GlobalManager.Spawn(ctx, subagents.TypeCode, "Generate Initial Snippet", "", "", nil)
	if err != nil {
		t.Fatalf("Coder failed to execute: %v", err)
	}

	// 2. Pass context to Researcher
	researcherResult, err := subagents.GlobalManager.Spawn(ctx, subagents.TypeResearch, "Verify syntax for snippet", coderResult, "", nil)
	if err != nil {
		t.Fatalf("Researcher failed to execute: %v", err)
	}

	// 3. Pass context to Planner
	plannerResult, err := subagents.GlobalManager.Spawn(ctx, subagents.TypePlan, "Design refactor based on research", researcherResult, "", nil)
	if err != nil {
		t.Fatalf("Planner failed to execute: %v", err)
	}

	// Validate context compaction/propagation mock fallbacks exist properly
	if plannerResult == "" {
		t.Fatalf("Handoff sequence broke: Empty final plan")
	}

	// Since we use the fallback LLM outputs in testing when no API key is provided, we check for those signals
	if !strings.Contains(plannerResult, "Analyze codebase") {
		t.Errorf("Expected planner to emit baseline mock output during test handoff chain. Got: %s", plannerResult)
	}
}

// TestAllSubagentTypes validates that all 10 defined subagent types can successfully spawn and return.
func TestAllSubagentTypes(t *testing.T) {
    types := []subagents.SubagentType{
        subagents.TypeCode, subagents.TypeResearch, subagents.TypeReview, subagents.TypePlan,
        subagents.TypeBuild, subagents.TypeTest, subagents.TypeDebug, subagents.TypeDoc,
        subagents.TypeSecurity, subagents.TypeDevOps,
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    for _, st := range types {
        t.Run(string(st), func(t *testing.T) {
            output, err := subagents.GlobalManager.Spawn(ctx, st, "ping", "test context", "", nil)
            if err != nil && !strings.Contains(err.Error(), "timeout") {
                t.Fatalf("Subagent %s failed to spawn: %v", st, err)
            }
            if output == "" {
                t.Errorf("Subagent %s returned empty output", st)
            }
        })
    }
}
