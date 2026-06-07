package agent

import (
	"context"
	"fmt"

	"github.com/robertpelloni/hyperharness/internal/council"
)

// AutopilotMode mimics Opencode-Autopilot and CLI-Orchestrator.
// It enters an autonomous loop using the Council architecture to achieve a goal.
func (a *Agent) AutopilotMode(goal string) (string, error) {
	fmt.Printf("[Autopilot] Goal set: %s\n", goal)

	director := council.NewDirectorAgent()

	fmt.Printf("[Autopilot] Decomposing goal into tasks...\n")
	_, err := director.PlanDelegation(context.Background(), a, goal)
	if err != nil {
		return "", fmt.Errorf("autopilot failed to plan: %w", err)
	}

	fmt.Printf("[Autopilot] Executing autonomous plan...\n")
	result, err := director.ExecutePlan()
	if err != nil {
		return "", fmt.Errorf("autopilot execution failed: %w", err)
	}

	return "Autopilot completed:\n" + result, nil
}
