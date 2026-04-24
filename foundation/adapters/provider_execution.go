package adapters

import (
	"fmt"
	"strings"
)

type ProviderExecutionRequest struct {
	Prompt         string `json:"prompt,omitempty"`
	TaskType       string `json:"taskType,omitempty"`
	CostPreference string `json:"costPreference,omitempty"`
	RequireLocal   bool   `json:"requireLocal,omitempty"`
}

type ProviderExecutionResult struct {
	TaskType       string         `json:"taskType"`
	Route          ProviderRoute  `json:"route"`
	Status         ProviderStatus `json:"status"`
	PromptPreview  string         `json:"promptPreview,omitempty"`
	ExecutionHint  string         `json:"executionHint"`
	SelectionNotes []string       `json:"selectionNotes,omitempty"`
}

func PrepareProviderExecution(req ProviderExecutionRequest) ProviderExecutionResult {
	taskType := normalizeTaskType(req.TaskType, req.Prompt)
	status := BuildProviderStatus()
	route := SelectProviderRoute(ProviderRouteRequest{
		TaskType:       taskType,
		CostPreference: req.CostPreference,
		RequireLocal:   req.RequireLocal,
	})
	preview := promptPreview(req.Prompt)
	hint := buildExecutionHint(route, taskType)
	return ProviderExecutionResult{
		TaskType:       taskType,
		Route:          route,
		Status:         status,
		PromptPreview:  preview,
		ExecutionHint:  hint,
		SelectionNotes: append([]string(nil), route.Reasons...),
	}
}

func normalizeTaskType(taskType, prompt string) string {
	taskType = strings.TrimSpace(taskType)
	if taskType != "" {
		return taskType
	}
	return inferTaskType(prompt)
}

func promptPreview(prompt string) string {
	preview := strings.TrimSpace(prompt)
	if len(preview) > 140 {
		preview = preview[:140] + "..."
	}
	return preview
}

func buildExecutionHint(route ProviderRoute, taskType string) string {
	return fmt.Sprintf("Route provider execution to %s/%s for %s work.", route.Provider, route.Model, taskType)
}

func inferTaskType(prompt string) string {
	lower := strings.ToLower(prompt)
	switch {
	case strings.Contains(lower, "refactor"), strings.Contains(lower, "fix"), strings.Contains(lower, "write code"), strings.Contains(lower, "implement"):
		return "coding"
	case strings.Contains(lower, "search"), strings.Contains(lower, "analyze"), strings.Contains(lower, "investigate"), strings.Contains(lower, "research"):
		return "analysis"
	case strings.Contains(lower, "local model"), strings.Contains(lower, "ollama"), strings.Contains(lower, "offline"):
		return "local"
	default:
		return "general"
	}
}
