package agent

import (
	"errors"
	"strings"
	"testing"
)

func TestAutopilotModeWithChatStopsWhenGoalAchieved(t *testing.T) {
	var prompts []string
	responses := []string{"progress step 1", " GOAL_ACHIEVED \n"}
	index := 0

	chat := func(prompt string) (string, error) {
		prompts = append(prompts, prompt)
		resp := responses[index]
		index++
		return resp, nil
	}

	result, err := autopilotModeWithChat("finish task", 5, chat)
	if err != nil {
		t.Fatalf("autopilot failed: %v", err)
	}
	if result != "Autopilot completed: Goal achieved." {
		t.Fatalf("unexpected result: %q", result)
	}
	if len(prompts) != 2 {
		t.Fatalf("expected 2 prompts before completion, got %d", len(prompts))
	}
	if !strings.Contains(prompts[0], "Goal: finish task") || !strings.Contains(prompts[0], "Iteration: 1") {
		t.Fatalf("unexpected first prompt: %q", prompts[0])
	}
	if !strings.Contains(prompts[1], "Iteration: 2") {
		t.Fatalf("unexpected second prompt: %q", prompts[1])
	}
}

func TestAutopilotModeWithChatHaltsAtMaxIterations(t *testing.T) {
	calls := 0
	chat := func(prompt string) (string, error) {
		calls++
		return "still working", nil
	}

	result, err := autopilotModeWithChat("finish task", 3, chat)
	if err != nil {
		t.Fatalf("autopilot failed: %v", err)
	}
	if result != "Autopilot halted: Max iterations reached." {
		t.Fatalf("unexpected result: %q", result)
	}
	if calls != 3 {
		t.Fatalf("expected 3 iterations, got %d", calls)
	}
}

func TestAutopilotModeWithChatValidatesInputs(t *testing.T) {
	if _, err := autopilotModeWithChat("   ", 3, func(string) (string, error) { return "", nil }); err == nil || !strings.Contains(err.Error(), "goal is required") {
		t.Fatalf("expected goal validation error, got %v", err)
	}
	if _, err := autopilotModeWithChat("goal", 0, func(string) (string, error) { return "", nil }); err == nil || !strings.Contains(err.Error(), "maxIterations must be positive") {
		t.Fatalf("expected iteration validation error, got %v", err)
	}
	if _, err := autopilotModeWithChat("goal", 1, nil); err == nil || !strings.Contains(err.Error(), "chat function is required") {
		t.Fatalf("expected chat validation error, got %v", err)
	}
}

func TestAutopilotModeWithChatPropagatesErrors(t *testing.T) {
	want := errors.New("chat failed")
	_, err := autopilotModeWithChat("goal", 2, func(string) (string, error) { return "", want })
	if !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

func TestBuildAutopilotPromptIncludesGoalAndIteration(t *testing.T) {
	prompt := buildAutopilotPrompt("ship feature", 4)
	if !strings.Contains(prompt, "Goal: ship feature") {
		t.Fatalf("expected goal in prompt, got %q", prompt)
	}
	if !strings.Contains(prompt, "Iteration: 4") {
		t.Fatalf("expected iteration in prompt, got %q", prompt)
	}
	if !strings.Contains(prompt, "respond with 'GOAL_ACHIEVED'") {
		t.Fatalf("expected completion contract in prompt, got %q", prompt)
	}
}
