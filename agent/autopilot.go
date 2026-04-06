package agent

import (
	"fmt"
	"strings"
)

const autopilotMaxIterations = 5

// AutopilotMode mimics Opencode-Autopilot and CLI-Orchestrator.
// It enters an autonomous loop until a goal is achieved.
func (a *Agent) AutopilotMode(goal string) (string, error) {
	return autopilotModeWithChat(strings.TrimSpace(goal), autopilotMaxIterations, a.Chat)
}

func autopilotModeWithChat(goal string, maxIterations int, chat func(string) (string, error)) (string, error) {
	fmt.Printf("[Autopilot] Goal set: %s\n", goal)
	if strings.TrimSpace(goal) == "" {
		return "", fmt.Errorf("goal is required")
	}
	if maxIterations <= 0 {
		return "", fmt.Errorf("maxIterations must be positive")
	}
	if chat == nil {
		return "", fmt.Errorf("chat function is required")
	}

	for i := 0; i < maxIterations; i++ {
		fmt.Printf("[Autopilot] Iteration %d/%d...\n", i+1, maxIterations)

		prompt := buildAutopilotPrompt(goal, i+1)
		response, err := chat(prompt)
		if err != nil {
			return "", err
		}

		if strings.TrimSpace(response) == "GOAL_ACHIEVED" {
			return "Autopilot completed: Goal achieved.", nil
		}

		fmt.Printf("[Autopilot] Progress: %s\n", response)
	}

	return "Autopilot halted: Max iterations reached.", nil
}

func buildAutopilotPrompt(goal string, iteration int) string {
	return fmt.Sprintf("Goal: %s\nIteration: %d\nYou are in Autopilot Mode. If the goal is not yet achieved, use your tools to make progress. If achieved, respond with 'GOAL_ACHIEVED'.", goal, iteration)
}
