// Package providers provides provider discovery, health checking, and model routing.
// Ported from hypercode/go/internal/providers/catalog.go with enhancements.
//
// WHAT: Provider catalog with status tracking, health checks, and model selection
// WHY: Models need to know which providers are available and healthy for routing
// HOW: Registry of providers with status, authentication, and model preference metadata
package providers

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// ProviderStatus represents the health/authentication state of a provider.
type ProviderStatus struct {
	Provider      string    `json:"provider"`
	Configured    bool      `json:"configured"`
	Authenticated bool      `json:"authenticated"`
	Reachable     bool      `json:"reachable"`
	LastChecked   time.Time `json:"lastChecked"`
	Latency       string    `json:"latency,omitempty"`
	Error         string    `json:"error,omitempty"`
}

// CatalogEntry describes a provider and its capabilities.
type CatalogEntry struct {
	Provider       string   `json:"provider"`
	Name           string   `json:"name"`
	AuthMethod     string   `json:"authMethod"`
	EnvVar         string   `json:"envVar"`
	DefaultModel   string   `json:"defaultModel"`
	PreferredTasks []string `json:"preferredTasks,omitempty"`
	Executable     bool     `json:"executable"`
	Models         []string `json:"models,omitempty"`
	BaseURL        string   `json:"baseUrl,omitempty"`
}

// CatalogModelInfo describes a specific model in the catalog.
type CatalogModelInfo struct {
	ID            string  `json:"id"`
	Provider     string  `json:"provider"`
	Name         string  `json:"name"`
	ContextWindow int    `json:"contextWindow"`
	CostPer1kIn  float64 `json:"costPer1kIn"`
	CostPer1kOut float64 `json:"costPer1kOut"`
	Speed        string  `json:"speed"` // "fast", "medium", "slow"
}

// Catalog contains all known providers.
var Catalog = []CatalogEntry{
	{
		Provider: "anthropic", Name: "Anthropic", AuthMethod: "api_key",
		EnvVar: "ANTHROPIC_API_KEY", DefaultModel: "claude-sonnet-4-20250514",
		PreferredTasks: []string{"planning", "coding", "research"},
		Executable: true,
		Models: []string{"claude-sonnet-4-20250514", "claude-3-7-sonnet-20250219", "claude-3-5-sonnet-20241022", "claude-3-haiku-20240307"},
		BaseURL: "https://api.anthropic.com",
	},
	{
		Provider: "google", Name: "Google Gemini", AuthMethod: "api_key",
		EnvVar: "GOOGLE_API_KEY", DefaultModel: "gemini-2.5-flash",
		PreferredTasks: []string{"coding", "research"},
		Executable: true,
		Models: []string{"gemini-2.5-flash", "gemini-2.5-pro", "gemini-2.0-flash", "gemini-1.5-pro"},
		BaseURL: "https://generativelanguage.googleapis.com",
	},
	{
		Provider: "openai", Name: "OpenAI", AuthMethod: "api_key",
		EnvVar: "OPENAI_API_KEY", DefaultModel: "gpt-4o",
		PreferredTasks: []string{"planning", "coding"},
		Executable: true,
		Models: []string{"gpt-4o", "gpt-4o-mini", "gpt-4-turbo", "o1", "o1-mini", "o3-mini"},
		BaseURL: "https://api.openai.com",
	},
	{
		Provider: "deepseek", Name: "DeepSeek", AuthMethod: "api_key",
		EnvVar: "DEEPSEEK_API_KEY", DefaultModel: "deepseek-chat",
		PreferredTasks: []string{"coding", "reasoning"},
		Executable: true,
		Models: []string{"deepseek-chat", "deepseek-reasoner"},
		BaseURL: "https://api.deepseek.com",
	},
	{
		Provider: "openrouter", Name: "OpenRouter", AuthMethod: "api_key",
		EnvVar: "OPENROUTER_API_KEY", DefaultModel: "openrouter/free",
		Executable: true,
		Models: []string{"anthropic/claude-3.5-sonnet", "google/gemini-2.0-flash-exp"},
		BaseURL: "https://openrouter.ai",
	},
	{
		Provider: "groq", Name: "Groq", AuthMethod: "api_key",
		EnvVar: "GROQ_API_KEY", DefaultModel: "llama-3.3-70b-versatile",
		PreferredTasks: []string{"fast-inference", "coding"},
		Executable: true,
		Models: []string{"llama-3.3-70b-versatile", "mixtral-8x7b-32768"},
		BaseURL: "https://api.groq.com",
	},
	{
		Provider: "lmstudio", Name: "LM Studio", AuthMethod: "local",
		DefaultModel: "local-model",
		PreferredTasks: []string{"local-inference"},
		Executable: true,
		Models: []string{"local-model"},
		BaseURL: "http://localhost:1234",
	},
	{
		Provider: "ollama", Name: "Ollama", AuthMethod: "local",
		DefaultModel: "gemma:2b",
		PreferredTasks: []string{"local-inference", "coding"},
		Executable: true,
		Models: []string{"gemma:2b", "llama3:8b", "codellama:7b", "mistral:7b"},
		BaseURL: "http://localhost:11434",
	},
}

// ModelCatalog contains detailed info about popular models.
var ModelCatalog = []CatalogModelInfo{
	{ID: "gpt-4o", Provider: "openai", Name: "GPT-4o", ContextWindow: 128000, CostPer1kIn: 0.0025, CostPer1kOut: 0.01, Speed: "medium"},
	{ID: "gpt-4o-mini", Provider: "openai", Name: "GPT-4o Mini", ContextWindow: 128000, CostPer1kIn: 0.00015, CostPer1kOut: 0.0006, Speed: "fast"},
	{ID: "claude-sonnet-4-20250514", Provider: "anthropic", Name: "Claude Sonnet 4", ContextWindow: 200000, CostPer1kIn: 0.003, CostPer1kOut: 0.015, Speed: "medium"},
	{ID: "claude-3-5-sonnet-20241022", Provider: "anthropic", Name: "Claude 3.5 Sonnet", ContextWindow: 200000, CostPer1kIn: 0.003, CostPer1kOut: 0.015, Speed: "medium"},
	{ID: "gemini-2.5-flash", Provider: "google", Name: "Gemini 2.5 Flash", ContextWindow: 1000000, CostPer1kIn: 0.0, CostPer1kOut: 0.0, Speed: "fast"},
	{ID: "gemini-2.5-pro", Provider: "google", Name: "Gemini 2.5 Pro", ContextWindow: 1000000, CostPer1kIn: 0.0, CostPer1kOut: 0.0, Speed: "medium"},
	{ID: "deepseek-chat", Provider: "deepseek", Name: "DeepSeek V3", ContextWindow: 64000, CostPer1kIn: 0.00014, CostPer1kOut: 0.00028, Speed: "fast"},
	{ID: "deepseek-reasoner", Provider: "deepseek", Name: "DeepSeek R1", ContextWindow: 64000, CostPer1kIn: 0.00055, CostPer1kOut: 0.00219, Speed: "slow"},
}

// StatusChecker checks provider health and authentication.
type StatusChecker struct {
	mu      sync.RWMutex
	statuses map[string]*ProviderStatus
}

// NewStatusChecker creates a new status checker.
func NewStatusChecker() *StatusChecker {
	return &StatusChecker{
		statuses: make(map[string]*ProviderStatus),
	}
}

// CheckAll checks all providers in the catalog.
func (sc *StatusChecker) CheckAll() []ProviderStatus {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	var statuses []ProviderStatus
	for _, entry := range Catalog {
		status := sc.checkEntry(entry)
		sc.statuses[entry.Provider] = status
		statuses = append(statuses, *status)
	}
	return statuses
}

// checkEntry checks a single catalog entry.
func (sc *StatusChecker) checkEntry(entry CatalogEntry) *ProviderStatus {
	status := &ProviderStatus{
		Provider:    entry.Provider,
		LastChecked: time.Now().UTC(),
	}

	// Local providers are always "configured"
	if entry.AuthMethod == "local" {
		status.Configured = true
		status.Authenticated = true
		status.Reachable = true
		return status
	}

	// Check env var
	apiKey := os.Getenv(entry.EnvVar)
	status.Configured = apiKey != ""
	status.Authenticated = apiKey != ""
	status.Reachable = apiKey != "" // Best guess without actual HTTP call

	return status
}

// GetStatus returns the cached status for a provider.
func (sc *StatusChecker) GetStatus(provider string) (*ProviderStatus, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	status, ok := sc.statuses[provider]
	if !ok {
		return nil, false
	}
	return status, true
}

// ListConfigured returns only providers that are configured.
func (sc *StatusChecker) ListConfigured() []CatalogEntry {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	var configured []CatalogEntry
	for _, entry := range Catalog {
		if status, ok := sc.statuses[entry.Provider]; ok {
			if status.Configured {
				configured = append(configured, entry)
			}
		} else if entry.AuthMethod == "local" {
			configured = append(configured, entry)
		} else if os.Getenv(entry.EnvVar) != "" {
			configured = append(configured, entry)
		}
	}
	return configured
}

// FindProviderForModel finds the provider that has a specific model.
func FindProviderForModel(modelID string) *CatalogEntry {
	for i := range Catalog {
		for _, m := range Catalog[i].Models {
			if m == modelID {
				return &Catalog[i]
			}
		}
	}
	return nil
}

// FindCatalogModelInfo returns detailed info about a model.
func FindCatalogModelInfo(modelID string) *CatalogModelInfo {
	for i := range ModelCatalog {
		if ModelCatalog[i].ID == modelID {
			return &ModelCatalog[i]
		}
	}
	return nil
}

// SelectModelForTask picks the best model for a given task type.
func SelectModelForTask(taskType string) *CatalogEntry {
	for i := range Catalog {
		for _, task := range Catalog[i].PreferredTasks {
			if task == taskType {
				envVar := Catalog[i].EnvVar
				if envVar == "" || os.Getenv(envVar) != "" {
					return &Catalog[i]
				}
			}
		}
	}
	return nil
}

// FormatCatalogSummary returns a formatted string of provider status.
func FormatCatalogSummary(statuses []ProviderStatus) string {
	var summary string
	summary += fmt.Sprintf("%-15s %-12s %-12s %-10s\n", "PROVIDER", "CONFIGURED", "AUTHED", "REACHABLE")
	summary += fmt.Sprintf("%-15s %-12s %-12s %-10s\n", "--------", "----------", "------", "--------")
	for _, s := range statuses {
		summary += fmt.Sprintf("%-15s %-12v %-12v %-10v\n",
			s.Provider, s.Configured, s.Authenticated, s.Reachable)
	}
	return summary
}
