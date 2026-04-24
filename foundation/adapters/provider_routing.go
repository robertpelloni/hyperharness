package adapters

import "strings"

const (
	defaultOpenAIModel = "gpt-4o"
	defaultGoogleModel = "gemini-2.5-flash"
	defaultOllamaModel = "llama3"
)

type ProviderRouteRequest struct {
	TaskType       string `json:"taskType,omitempty"`
	CostPreference string `json:"costPreference,omitempty"`
	RequireLocal   bool   `json:"requireLocal,omitempty"`
}

type ProviderRoute struct {
	Provider string   `json:"provider"`
	Model    string   `json:"model"`
	Reasons  []string `json:"reasons,omitempty"`
}

func SelectProviderRoute(req ProviderRouteRequest) ProviderRoute {
	status := BuildProviderStatus()
	provider := status.CurrentProvider
	model := status.CurrentModel
	reasons := make([]string, 0, 4)

	if req.RequireLocal && containsString(status.Available, "ollama") {
		provider = "ollama"
		model = firstOr(status.OllamaModels, model)
		reasons = append(reasons, "local execution requested and ollama is available")
	}

	switch normalizedPreference(req.CostPreference) {
	case "budget":
		if containsString(status.Available, "google") {
			provider = "google"
			if model == "" || provider != status.CurrentProvider {
				model = defaultGoogleModel
			}
			reasons = append(reasons, "budget preference favored a lower-cost provider profile")
		} else {
			reasons = append(reasons, "budget preference requested but no alternate provider detected")
		}
	case "quality":
		if containsString(status.Available, "openai") {
			provider = "openai"
			if model == "" || provider != status.CurrentProvider {
				model = defaultOpenAIModel
			}
			reasons = append(reasons, "quality preference favored the strongest available default route")
		}
	}

	switch normalizedTaskType(req.TaskType) {
	case "coding":
		reasons = append(reasons, "coding workload detected")
	case "analysis":
		reasons = append(reasons, "analysis workload detected")
	case "local":
		if containsString(status.Available, "ollama") {
			provider = "ollama"
			model = firstOr(status.OllamaModels, model)
			reasons = append(reasons, "local-only task type requested")
		}
	}

	if provider == "" {
		provider = "openai"
		reasons = append(reasons, "fallback provider defaulted to openai")
	}
	if model == "" {
		model = defaultModelForProvider(provider, status.OllamaModels)
	}
	if len(reasons) == 0 {
		reasons = append(reasons, "current configured provider and model remained suitable")
	}
	return ProviderRoute{Provider: provider, Model: model, Reasons: reasons}
}

func normalizedPreference(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "cheap", "low", "budget":
		return "budget"
	case "high", "quality", "best":
		return "quality"
	default:
		return ""
	}
}

func normalizedTaskType(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "code", "coding", "edit", "refactor":
		return "coding"
	case "search", "analysis", "research":
		return "analysis"
	case "local":
		return "local"
	default:
		return ""
	}
}

func firstOr(values []string, fallback string) string {
	if len(values) > 0 && strings.TrimSpace(values[0]) != "" {
		return values[0]
	}
	return fallback
}

func defaultModelForProvider(provider string, ollamaModels []string) string {
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "ollama":
		return firstOr(ollamaModels, defaultOllamaModel)
	case "google":
		return defaultGoogleModel
	default:
		return defaultOpenAIModel
	}
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if strings.EqualFold(strings.TrimSpace(value), target) {
			return true
		}
	}
	return false
}
