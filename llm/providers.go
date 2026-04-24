// Package llm provides multi-provider LLM routing and generation.
// Ported from hypercode/go/internal/ai/llm.go with enhancements.
//
// WHAT: Unified LLM interface supporting OpenAI, Anthropic, Gemini, DeepSeek, OpenRouter, LMStudio, Ollama
// WHY: All major AI providers have compatible APIs; auto-routing selects best available
// HOW: Provider interface pattern with environment-based selection and priority ordering
package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Message represents a chat message for LLM APIs.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// TokenUsage tracks token consumption.
type TokenUsage struct {
	InputTokens  int `json:"inputTokens"`
	OutputTokens int `json:"outputTokens"`
}

// LLMResponse is the unified response from any provider.
type LLMResponse struct {
	Content  string    `json:"content"`
	Usage    TokenUsage `json:"usage"`
	Provider string    `json:"provider"`
	Model    string    `json:"model"`
}

// Provider is the interface that all LLM providers implement.
type Provider interface {
	// Name returns the provider identifier (e.g. "openai", "anthropic").
	Name() string
	// GenerateText sends messages to the LLM and returns a response.
	GenerateText(ctx context.Context, model string, messages []Message) (*LLMResponse, error)
	// Models returns the list of known model IDs for this provider.
	Models() []string
}

// ---------------------------------------------------------------------------
// OpenAI Provider (also base for DeepSeek, OpenRouter, LMStudio)
// ---------------------------------------------------------------------------

// OpenAIProvider implements the OpenAI Chat Completions API.
// Also used as base for any OpenAI-compatible endpoint.
type OpenAIProvider struct {
	APIKey  string
	BaseURL string
}

func (p *OpenAIProvider) Name() string { return "openai" }

func (p *OpenAIProvider) Models() []string {
	return []string{"gpt-4o", "gpt-4o-mini", "gpt-4-turbo", "gpt-3.5-turbo", "o1", "o1-mini", "o3-mini"}
}

func (p *OpenAIProvider) GenerateText(ctx context.Context, model string, messages []Message) (*LLMResponse, error) {
	if p.BaseURL == "" {
		p.BaseURL = "https://api.openai.com/v1/chat/completions"
	}

	reqBody, _ := json.Marshal(map[string]interface{}{
		"model":    model,
		"messages": messages,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", p.BaseURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("openai: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.APIKey)

	return doOpenAICompatibleRequest(req, "openai", model)
}

// ---------------------------------------------------------------------------
// Anthropic Provider
// ---------------------------------------------------------------------------

// AnthropicProvider implements the Anthropic Messages API.
type AnthropicProvider struct {
	APIKey  string
	BaseURL string
}

func (p *AnthropicProvider) Name() string { return "anthropic" }

func (p *AnthropicProvider) Models() []string {
	return []string{"claude-sonnet-4-20250514", "claude-3-7-sonnet-20250219", "claude-3-5-sonnet-20241022", "claude-3-haiku-20240307"}
}

func (p *AnthropicProvider) GenerateText(ctx context.Context, model string, messages []Message) (*LLMResponse, error) {
	if p.BaseURL == "" {
		p.BaseURL = "https://api.anthropic.com/v1/messages"
	}

	// Extract system message (Anthropic handles it separately)
	var systemContent string
	var filtered []Message
	for _, m := range messages {
		if m.Role == "system" {
			systemContent = m.Content
		} else {
			filtered = append(filtered, m)
		}
	}

	body := map[string]interface{}{
		"model":      model,
		"max_tokens": 16384,
		"messages":   filtered,
	}
	if systemContent != "" {
		body["system"] = systemContent
	}

	reqBody, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, "POST", p.BaseURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("anthropic: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("anthropic: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("anthropic: API error %s: %s", resp.Status, body)
	}

	var payload struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("anthropic: decode: %w", err)
	}
	if len(payload.Content) == 0 {
		return nil, fmt.Errorf("anthropic: no content returned")
	}

	return &LLMResponse{
		Content:  payload.Content[0].Text,
		Provider: "anthropic",
		Model:    model,
		Usage: TokenUsage{
			InputTokens:  payload.Usage.InputTokens,
			OutputTokens: payload.Usage.OutputTokens,
		},
	}, nil
}

// ---------------------------------------------------------------------------
// Google Gemini Provider
// ---------------------------------------------------------------------------

// GeminiProvider implements the Google Generative AI API.
type GeminiProvider struct {
	APIKey  string
	BaseURL string
}

func (p *GeminiProvider) Name() string { return "google" }

func (p *GeminiProvider) Models() []string {
	return []string{"gemini-2.5-flash", "gemini-2.5-pro", "gemini-2.0-flash", "gemini-1.5-pro", "gemini-1.5-flash"}
}

func (p *GeminiProvider) GenerateText(ctx context.Context, model string, messages []Message) (*LLMResponse, error) {
	if p.BaseURL == "" {
		p.BaseURL = "https://generativelanguage.googleapis.com/v1beta"
	}

	// Convert messages to Gemini content format
	contents := make([]map[string]interface{}, 0, len(messages))
	var systemInstruction string
	for _, msg := range messages {
		if msg.Role == "system" {
			systemInstruction = msg.Content
			continue
		}
		role := msg.Role
		if role == "assistant" {
			role = "model"
		}
		contents = append(contents, map[string]interface{}{
			"role":  role,
			"parts": []map[string]string{{"text": msg.Content}},
		})
	}

	body := map[string]interface{}{"contents": contents}
	if systemInstruction != "" {
		body["systemInstruction"] = map[string]interface{}{
			"parts": []map[string]string{{"text": systemInstruction}},
		}
	}

	reqBody, _ := json.Marshal(body)
	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", p.BaseURL, model, p.APIKey)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("gemini: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gemini: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gemini: API error %s: %s", resp.Status, respBody)
	}

	var payload struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
		UsageMetadata struct {
			PromptTokenCount     int `json:"promptTokenCount"`
			CandidatesTokenCount int `json:"candidatesTokenCount"`
		} `json:"usageMetadata"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("gemini: decode: %w", err)
	}
	if len(payload.Candidates) == 0 || len(payload.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("gemini: no content returned")
	}

	var sb strings.Builder
	for _, part := range payload.Candidates[0].Content.Parts {
		sb.WriteString(part.Text)
	}

	return &LLMResponse{
		Content:  sb.String(),
		Provider: "google",
		Model:    model,
		Usage: TokenUsage{
			InputTokens:  payload.UsageMetadata.PromptTokenCount,
			OutputTokens: payload.UsageMetadata.CandidatesTokenCount,
		},
	}, nil
}

// ---------------------------------------------------------------------------
// DeepSeek Provider (OpenAI-compatible)
// ---------------------------------------------------------------------------

// DeepSeekProvider uses the OpenAI-compatible API at api.deepseek.com.
type DeepSeekProvider struct{ APIKey string }

func (p *DeepSeekProvider) Name() string { return "deepseek" }

func (p *DeepSeekProvider) Models() []string {
	return []string{"deepseek-chat", "deepseek-reasoner"}
}

func (p *DeepSeekProvider) GenerateText(ctx context.Context, model string, messages []Message) (*LLMResponse, error) {
	resp, err := (&OpenAIProvider{APIKey: p.APIKey, BaseURL: "https://api.deepseek.com/v1/chat/completions"}).GenerateText(ctx, model, messages)
	if err != nil {
		return nil, err
	}
	resp.Provider = "deepseek"
	return resp, nil
}

// ---------------------------------------------------------------------------
// OpenRouter Provider (OpenAI-compatible)
// ---------------------------------------------------------------------------

// OpenRouterProvider uses the OpenAI-compatible API at openrouter.ai.
type OpenRouterProvider struct{ APIKey string }

func (p *OpenRouterProvider) Name() string { return "openrouter" }

func (p *OpenRouterProvider) Models() []string {
	return []string{"openrouter/free", "anthropic/claude-3.5-sonnet", "google/gemini-2.0-flash-exp"}
}

func (p *OpenRouterProvider) GenerateText(ctx context.Context, model string, messages []Message) (*LLMResponse, error) {
	resp, err := (&OpenAIProvider{
		APIKey:  p.APIKey,
		BaseURL: "https://openrouter.ai/api/v1/chat/completions",
	}).GenerateText(ctx, model, messages)
	if err != nil {
		return nil, err
	}
	resp.Provider = "openrouter"
	return resp, nil
}

// ---------------------------------------------------------------------------
// LM Studio Provider (local OpenAI-compatible)
// ---------------------------------------------------------------------------

// LMStudioProvider uses a local OpenAI-compatible server (default :1234).
type LMStudioProvider struct{ BaseURL string }

func (p *LMStudioProvider) Name() string { return "lmstudio" }

func (p *LMStudioProvider) Models() []string { return []string{"local-model"} }

func (p *LMStudioProvider) GenerateText(ctx context.Context, model string, messages []Message) (*LLMResponse, error) {
	baseURL := p.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:1234/v1/chat/completions"
	}
	resp, err := (&OpenAIProvider{APIKey: "lm-studio", BaseURL: baseURL}).GenerateText(ctx, model, messages)
	if err != nil {
		return nil, err
	}
	resp.Provider = "lmstudio"
	return resp, nil
}

// ---------------------------------------------------------------------------
// Ollama Provider
// ---------------------------------------------------------------------------

// OllamaProvider uses the local Ollama API at localhost:11434.
type OllamaProvider struct{ BaseURL string }

func (p *OllamaProvider) Name() string { return "ollama" }

func (p *OllamaProvider) Models() []string {
	return []string{"gemma:2b", "llama3:8b", "codellama:7b", "mistral:7b", "phi3:mini"}
}

func (p *OllamaProvider) GenerateText(ctx context.Context, model string, messages []Message) (*LLMResponse, error) {
	baseURL := p.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:11434/api/chat"
	}

	reqBody, _ := json.Marshal(map[string]interface{}{
		"model":    model,
		"messages": messages,
		"stream":   false,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", baseURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("ollama: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ollama: request: %w", err)
	}
	defer resp.Body.Close()

	var payload struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
		PromptEvalCount int `json:"prompt_eval_count"`
		EvalCount       int `json:"eval_count"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("ollama: decode: %w", err)
	}

	return &LLMResponse{
		Content:  payload.Message.Content,
		Provider: "ollama",
		Model:    model,
		Usage: TokenUsage{
			InputTokens:  payload.PromptEvalCount,
			OutputTokens: payload.EvalCount,
		},
	}, nil
}

// ---------------------------------------------------------------------------
// Groq Provider (OpenAI-compatible, fast inference)
// ---------------------------------------------------------------------------

// GroqProvider uses the OpenAI-compatible API at api.groq.com.
type GroqProvider struct{ APIKey string }

func (p *GroqProvider) Name() string { return "groq" }

func (p *GroqProvider) Models() []string {
	return []string{"llama-3.3-70b-versatile", "mixtral-8x7b-32768", "gemma2-9b-it"}
}

func (p *GroqProvider) GenerateText(ctx context.Context, model string, messages []Message) (*LLMResponse, error) {
	resp, err := (&OpenAIProvider{
		APIKey:  p.APIKey,
		BaseURL: "https://api.groq.com/openai/v1/chat/completions",
	}).GenerateText(ctx, model, messages)
	if err != nil {
		return nil, err
	}
	resp.Provider = "groq"
	return resp, nil
}

// ---------------------------------------------------------------------------
// Auto-Router
// ---------------------------------------------------------------------------

// ProviderEntry defines a provider for priority-based auto-routing.
type ProviderEntry struct {
	EnvVar       string
	ProviderName string
	DefaultModel string
	Factory      func(apiKey string) Provider
}

// ProviderPriority defines the selection order for auto-routing.
// First provider with a valid API key wins.
var ProviderPriority = []ProviderEntry{
	{"ANTHROPIC_API_KEY", "anthropic", "claude-sonnet-4-20250514", func(k string) Provider { return &AnthropicProvider{APIKey: k} }},
	{"GOOGLE_API_KEY", "google", "gemini-2.5-flash", func(k string) Provider { return &GeminiProvider{APIKey: k} }},
	{"GEMINI_API_KEY", "google", "gemini-2.5-flash", func(k string) Provider { return &GeminiProvider{APIKey: k} }},
	{"OPENAI_API_KEY", "openai", "gpt-4o", func(k string) Provider { return &OpenAIProvider{APIKey: k} }},
	{"DEEPSEEK_API_KEY", "deepseek", "deepseek-chat", func(k string) Provider { return &DeepSeekProvider{APIKey: k} }},
	{"OPENROUTER_API_KEY", "openrouter", "openrouter/free", func(k string) Provider { return &OpenRouterProvider{APIKey: k} }},
	{"GROQ_API_KEY", "groq", "llama-3.3-70b-versatile", func(k string) Provider { return &GroqProvider{APIKey: k} }},
	{"", "lmstudio", "local-model", func(_ string) Provider { return &LMStudioProvider{} }},
	{"", "ollama", "gemma:2b", func(_ string) Provider { return &OllamaProvider{} }},
}

type providerSelection struct {
	ProviderName string
	DefaultModel string
	Factory      func(apiKey string) Provider
	APIKey       string
}

// resolveProvider finds the best available provider based on environment.
func resolveProvider() (providerSelection, bool) {
	for _, entry := range ProviderPriority {
		key := os.Getenv(entry.EnvVar)
		if entry.EnvVar == "" || key != "" {
			return providerSelection{
				ProviderName: entry.ProviderName,
				DefaultModel: entry.DefaultModel,
				Factory:      entry.Factory,
				APIKey:       key,
			}, true
		}
	}
	return providerSelection{}, false
}

// AutoRoute selects the best available provider and generates text.
// Priority: Anthropic > Gemini > OpenAI > DeepSeek > OpenRouter > Groq > LMStudio > Ollama.
func AutoRoute(ctx context.Context, messages []Message) (*LLMResponse, error) {
	sel, ok := resolveProvider()
	if !ok {
		return nil, fmt.Errorf("no LLM provider configured (set ANTHROPIC_API_KEY, GOOGLE_API_KEY, OPENAI_API_KEY, etc.)")
	}
	return sel.Factory(sel.APIKey).GenerateText(ctx, sel.DefaultModel, messages)
}

// AutoRouteWithModel auto-routes but allows model override.
func AutoRouteWithModel(ctx context.Context, model string, messages []Message) (*LLMResponse, error) {
	sel, ok := resolveProvider()
	if !ok {
		return nil, fmt.Errorf("no LLM provider configured")
	}
	if model == "" {
		model = sel.DefaultModel
	}
	return sel.Factory(sel.APIKey).GenerateText(ctx, model, messages)
}

// ListConfiguredProviders returns which providers have API keys set.
func ListConfiguredProviders() []string {
	var configured []string
	seen := map[string]struct{}{}
	for _, entry := range ProviderPriority {
		if entry.EnvVar != "" && os.Getenv(entry.EnvVar) == "" {
			continue
		}
		if _, ok := seen[entry.ProviderName]; ok {
			continue
		}
		seen[entry.ProviderName] = struct{}{}
		configured = append(configured, entry.ProviderName)
	}
	return configured
}

// GetProvider creates a specific provider by name.
func GetProvider(name string) (Provider, bool) {
	switch name {
	case "openai":
		if key := os.Getenv("OPENAI_API_KEY"); key != "" {
			return &OpenAIProvider{APIKey: key}, true
		}
	case "anthropic":
		if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
			return &AnthropicProvider{APIKey: key}, true
		}
	case "google", "gemini":
		for _, env := range []string{"GOOGLE_API_KEY", "GEMINI_API_KEY"} {
			if key := os.Getenv(env); key != "" {
				return &GeminiProvider{APIKey: key}, true
			}
		}
	case "deepseek":
		if key := os.Getenv("DEEPSEEK_API_KEY"); key != "" {
			return &DeepSeekProvider{APIKey: key}, true
		}
	case "openrouter":
		if key := os.Getenv("OPENROUTER_API_KEY"); key != "" {
			return &OpenRouterProvider{APIKey: key}, true
		}
	case "groq":
		if key := os.Getenv("GROQ_API_KEY"); key != "" {
			return &GroqProvider{APIKey: key}, true
		}
	case "lmstudio":
		return &LMStudioProvider{}, true
	case "ollama":
		return &OllamaProvider{}, true
	}
	return nil, false
}

// ---------------------------------------------------------------------------
// Shared helpers
// ---------------------------------------------------------------------------

// doOpenAICompatibleRequest executes a request against any OpenAI-compatible endpoint.
func doOpenAICompatibleRequest(req *http.Request, provider, model string) (*LLMResponse, error) {
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: request: %w", provider, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%s: API error %s: %s", provider, resp.Status, body)
	}

	var payload struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
		} `json:"usage"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("%s: decode: %w", provider, err)
	}
	if len(payload.Choices) == 0 {
		return nil, fmt.Errorf("%s: no choices returned", provider)
	}

	return &LLMResponse{
		Content:  payload.Choices[0].Message.Content,
		Provider: provider,
		Model:    model,
		Usage: TokenUsage{
			InputTokens:  payload.Usage.PromptTokens,
			OutputTokens: payload.Usage.CompletionTokens,
		},
	}, nil
}
