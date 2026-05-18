// Package providers implements LLM provider registry and client abstraction.
// Supports: Anthropic, OpenAI, Google Gemini, Azure OpenAI, Groq, Cerebras,
// xAI, Mistral, OpenRouter, Bedrock, Vertex, and arbitrary custom providers
// via URL configuration. Unified interface matching Pi's provider model
// with enhanced routing, failover, and quota awareness borrowed from hypercode.
package providers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// ProviderType identifies supported LLM providers.
type ProviderType string

const (
	ProviderAnthropic   ProviderType = "anthropic"
	ProviderOpenAI      ProviderType = "openai"
	ProviderGoogle      ProviderType = "google"
	ProviderAzure       ProviderType = "azure"
	ProviderGroq        ProviderType = "groq"
	ProviderCerebras    ProviderType = "cerebras"
	ProviderXAI         ProviderType = "xai"
	ProviderMistral     ProviderType = "mistral"
	ProviderOpenRouter  ProviderType = "openrouter"
	ProviderBedrock     ProviderType = "bedrock"
	ProviderVertex      ProviderType = "vertex"
	ProviderMiniMax     ProviderType = "minimax"
	ProviderHuggingFace ProviderType = "huggingface"
	ProviderOpenCode    ProviderType = "opencode"
	ProviderCustom      ProviderType = "custom"
)

// Transport defines the API transport type.
type Transport string

const (
	TransportSSE       Transport = "sse"
	TransportWebsocket Transport = "websocket"
	TransportAuto      Transport = "auto"
)

// ThinkingLevel maps to provider-specific thinking/thinking budgets.
type ThinkingLevel string

const (
	ThinkingOff     ThinkingLevel = "off"
	ThinkingMinimal ThinkingLevel = "minimal"
	ThinkingLow     ThinkingLevel = "low"
	ThinkingMedium  ThinkingLevel = "medium"
	ThinkingHigh    ThinkingLevel = "high"
	ThinkingXHigh   ThinkingLevel = "xhigh"
)

// ThinkingBudget maps levels to token counts.
var DefaultThinkingBudgets = map[ThinkingLevel]int{
	ThinkingMinimal: 1024,
	ThinkingLow:     4096,
	ThinkingMedium:  10240,
	ThinkingHigh:    32768,
	ThinkingXHigh:   65536,
}

// ProviderConfig holds everything needed to call a provider API.
type ProviderConfig struct {
	Type       ProviderType  `json:"type"`
	APIKey     string        `json:"apiKey,omitempty"`
	BaseURL    string        `json:"baseURL,omitempty"`
	APIVersion string        `json:"apiVersion,omitempty"`
	Region     string        `json:"region,omitempty"`
	ProjectID  string        `json:"projectID,omitempty"`
	Timeout    time.Duration `json:"timeout"`
	MaxRetries int           `json:"maxRetries"`
	RetryDelay time.Duration `json:"retryDelay"`
	Transport  Transport     `json:"transport"`
}

// ModelInfo describes a model available on a provider.
type ModelInfo struct {
	ID                     string       `json:"id"`
	Name                   string       `json:"name"`
	Provider               ProviderType `json:"provider"`
	ContextWindow          int          `json:"contextWindow"`
	MaxOutput              int          `json:"maxOutput"`
	SupportsVision         bool         `json:"supportsVision"`
	SupportsStreaming      bool         `json:"supportsStreaming"`
	SupportsThinking       bool         `json:"supportsThinking"`
	CacheSupport           bool         `json:"cacheSupport"`
	CostPerInputMT         float64      `json:"costPerInputMT"` // per million tokens
	CostPerOutputMT        float64      `json:"costPerOutputMT"`
	CostPerCacheReadMT     float64      `json:"costPerCacheReadMT"`
	CostPerCacheCreationMT float64      `json:"costPerCacheCreationMT"`
	Tags                   []string     `json:"tags,omitempty"`
}

// Provider is the interface all provider implementations must satisfy.
type Provider interface {
	Type() ProviderType
	Name() string
	Configure(config ProviderConfig) error
	GetModels() ([]ModelInfo, error)
	MatchesModel(pattern string) bool
	// StreamCompletion handles streaming chat completions.
	StreamCompletion(ctx context.Context, req CompletionRequest) (chan StreamChunk, error)
	// CreateCompletion handles non-streaming completions.
	CreateCompletion(ctx context.Context, req CompletionRequest) (*CompletionResult, error)
	// EstimateTokens estimates tokens for the given text.
	EstimateTokens(text string) int
}

// StreamingToolUse is a tool call within a streaming chunk.
type StreamingToolUse struct {
	ToolID   string `json:"toolId"`
	ToolType string `json:"toolType"`
	Input    []byte `json:"input"`
}

// StreamChunk is a single chunk from a streaming response.
type StreamChunk struct {
	Type       string             `json:"type"`
	Content    string             `json:"content,omitempty"`
	StopReason string             `json:"stopReason,omitempty"`
	Usage      *TokenUsage        `json:"usage,omitempty"`
	ToolUses   []StreamingToolUse `json:"toolUses,omitempty"`
	Thinking   string             `json:"thinking,omitempty"`
	Err        error              `json:"-"`
}

// TokenUsage tracks token consumption.
type TokenUsage struct {
	InputTokens         int `json:"inputTokens"`
	OutputTokens        int `json:"outputTokens"`
	CacheReadTokens     int `json:"cacheReadTokens"`
	CacheCreationTokens int `json:"cacheCreationTokens"`
}

// CostReport tracks costs.
type CostReport struct {
	InputCost       float64 `json:"inputCost"`
	OutputCost      float64 `json:"outputCost"`
	CacheReadCost   float64 `json:"cacheReadCost"`
	CacheCreateCost float64 `json:"cacheCreateCost"`
	TotalCost       float64 `json:"totalCost"`
	Currency        string  `json:"currency"`
}

// ComputeCost calculates cost based on model pricing and usage.
func (u *TokenUsage) ComputeCost(model ModelInfo) *CostReport {
	report := &CostReport{Currency: "USD"}
	mtIn := float64(u.InputTokens-u.CacheReadTokens) / 1e6
	mtOut := float64(u.OutputTokens) / 1e6
	mtCacheRead := float64(u.CacheReadTokens) / 1e6
	mtCacheCreate := float64(u.CacheCreationTokens) / 1e6

	report.InputCost = mtIn * model.CostPerInputMT
	report.OutputCost = mtOut * model.CostPerOutputMT
	report.CacheReadCost = mtCacheRead * model.CostPerCacheReadMT
	report.CacheCreateCost = mtCacheCreate * model.CostPerCacheCreationMT
	report.TotalCost = report.InputCost + report.OutputCost + report.CacheReadCost + report.CacheCreateCost
	return report
}

// Message represents a chat message (compatible with OpenAI/Anthropic formats).
type Message struct {
	Role    string         `json:"role"` // "user", "assistant", "system", "tool"
	Content []ContentBlock `json:"content"`
}

// ContentBlock is a single block of content in a message.
type ContentBlock struct {
	Type     string                 `json:"type"` // "text", "image", "tool_use", "tool_result"
	Text     string                 `json:"text,omitempty"`
	Image    *ImageContent          `json:"image,omitempty"`
	ToolID   string                 `json:"toolId,omitempty"`
	ToolType string                 `json:"toolType,omitempty"`
	Input    map[string]interface{} `json:"input,omitempty"`
	Name     string                 `json:"name,omitempty"`
}

// ImageContent is an image in a message.
type ImageContent struct {
	MediaType string `json:"mediaType"` // e.g. "image/jpeg"
	Data      []byte `json:"data"`      // base64 encoded
}

// ToolDefinition describes a tool the model can call.
type ToolDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"` // JSON Schema
}

// CompletionRequest is a unified request across all providers.
type CompletionRequest struct {
	ModelID          string            `json:"modelId"`
	System           string            `json:"system,omitempty"`
	Messages         []Message         `json:"messages"`
	MaxTokens        int               `json:"maxTokens,omitempty"`
	Temperature      float64           `json:"temperature,omitempty"`
	TopP             float64           `json:"topP,omitempty"`
	TopK             int               `json:"topK,omitempty"`
	FrequencyPenalty float64           `json:"frequencyPenalty,omitempty"`
	StopSequences    []string          `json:"stop,omitempty"`
	Tools            []ToolDefinition  `json:"tools,omitempty"`
	ToolChoice       string            `json:"toolChoice,omitempty"` // "auto", "any", "tool", "none"
	Stream           bool              `json:"stream"`
	ThinkingBudget   int               `json:"thinkingBudget,omitempty"`
	ThinkingType     ThinkingLevel     `json:"thinkingType,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}

// CompletionResult is a unified completion response.
type CompletionResult struct {
	ID             string         `json:"id"`
	Model          string         `json:"model"`
	Content        []ContentBlock `json:"content"`
	StopReason     string         `json:"stopReason"`
	Usage          *TokenUsage    `json:"usage"`
	Cost           *CostReport    `json:"cost,omitempty"`
	Thinking       string         `json:"thinking,omitempty"`
	CacheID        string         `json:"cacheId,omitempty"`
	ResponseTokens int            `json:"responseTokens"`
}

// Registry manages all known providers and models.
type Registry struct {
	providers map[ProviderType]Provider
	models    []ModelInfo
	configs   map[ProviderType]*ProviderConfig
	mu        sync.RWMutex
}

// NewRegistry creates a provider registry.
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[ProviderType]Provider),
		models:    make([]ModelInfo, 0),
		configs:   make(map[ProviderType]*ProviderConfig),
	}
}

// RegisterProvider registers a provider implementation.
func (r *Registry) RegisterProvider(p Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[p.Type()] = p

	// Fetch and register models
	models, err := p.GetModels()
	if err == nil {
		// Remove existing models from same provider first
		var filtered []ModelInfo
		for _, m := range r.models {
			if m.Provider != p.Type() {
				filtered = append(filtered, m)
			}
		}
		r.models = append(filtered, models...)
	}
}

// ConfigureProvider configures authentication for a provider.
func (r *Registry) ConfigureProvider(pType ProviderType, config ProviderConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	config.Type = pType
	r.configs[pType] = &config

	if provider, ok := r.providers[pType]; ok {
		return provider.Configure(config)
	}

	// Auto-create standard provider implementation if not registered
	provider, err := CreateProvider(pType, config)
	if err != nil {
		return err
	}

	r.providers[pType] = provider
	models, err := provider.GetModels()
	if err == nil {
		r.models = append(r.models, models...)
	}
	return err
}

// GetProvider returns the provider by type.
func (r *Registry) GetProvider(pType ProviderType) (Provider, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	provider, ok := r.providers[pType]
	return provider, ok
}

// GetModels returns all known models, optionally filtered by pattern.
func (r *Registry) GetModels(pattern string) []ModelInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if pattern == "" {
		return r.models
	}
	var filtered []ModelInfo
	for _, m := range r.models {
		if matchesModelPattern(m, pattern) {
			filtered = append(filtered, m)
		}
	}
	return filtered
}

// ResolveModel parses a model spec like "anthropic/claude-sonnet-4-20250514" or "sonnet:high".
func ResolveModel(spec string) (ProviderType, string, ThinkingLevel) {
	var provider ProviderType
	var modelID string
	var thinking ThinkingLevel

	// Check for provider/model format
	if parts := strings.SplitN(spec, "/", 2); len(parts) == 2 {
		provider = ProviderType(parts[0])
		modelID = parts[1]
	} else {
		modelID = spec
		// Try to infer provider from model name
		provider = InferProviderFromModel(modelID)
	}

	// Check for thinking level shorthand like "sonnet:high"
	if parts := strings.SplitN(modelID, ":", 2); len(parts) == 2 {
		modelID = parts[0]
		thinking = ThinkingLevel(parts[1])
	}

	return provider, modelID, thinking
}

// InferProviderFromModel guesses provider from model ID/name.
func InferProviderFromModel(modelID string) ProviderType {
	lower := strings.ToLower(modelID)
	switch {
	case strings.HasPrefix(lower, "claude"):
		return ProviderAnthropic
	case strings.HasPrefix(lower, "gpt-4"), strings.HasPrefix(lower, "o1"), strings.HasPrefix(lower, "o3"), strings.HasPrefix(lower, "o4"):
		return ProviderOpenAI
	case strings.HasPrefix(lower, "gemini"):
		return ProviderGoogle
	case strings.HasPrefix(lower, "grok"):
		return ProviderXAI
	case strings.HasPrefix(lower, "mistral"), strings.HasPrefix(lower, "mixtral"):
		return ProviderMistral
	case strings.HasPrefix(lower, "llama"):
		return ProviderGroq // Groq commonly serves Llama
	case strings.HasPrefix(lower, "codestral"):
		return ProviderMistral
	default:
		return ProviderOpenAI
	}
}

// matchesModelPattern checks if a model matches a glob pattern.
func matchesModelPattern(m ModelInfo, pattern string) bool {
	if pattern == "" {
		return true
	}
	// Simple glob matching
	if strings.Contains(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		if strings.HasPrefix(m.ID, prefix) {
			return true
		}
	}
	return strings.EqualFold(m.ID, pattern) || strings.Contains(strings.ToLower(m.ID), strings.ToLower(pattern))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// CreateProvider creates a concrete provider implementation.
func CreateProvider(pType ProviderType, config ProviderConfig) (Provider, error) {
	switch pType {
	case ProviderAnthropic:
		return newAnthropicProvider(config)
	case ProviderOpenAI:
		return newOpenAIProvider(config)
	case ProviderGoogle:
		return newGoogleProvider(config)
	case ProviderMistral:
		return newMistralProvider(config)
	default:
		return newGenericProvider(config)
	}
}

// ---- Provider Implementations ----

// baseProvider provides common HTTP client and retry logic.
type baseProvider struct {
	config ProviderConfig
	client *http.Client
}

func newBaseProvider(config ProviderConfig) baseProvider {
	if config.Timeout == 0 {
		config.Timeout = 5 * time.Minute
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.RetryDelay == 0 {
		config.RetryDelay = 2 * time.Second
	}
	return baseProvider{
		config: config,
		client: &http.Client{Timeout: config.Timeout},
	}
}

// Anthropic provider implementation.
type anthropicProvider struct {
	baseProvider
	models []ModelInfo
	mu     sync.RWMutex
}

func newAnthropicProvider(config ProviderConfig) (Provider, error) {
	bp := newBaseProvider(config)
	if config.BaseURL == "" {
		config.BaseURL = "https://api.anthropic.com/v1"
	}
	p := &anthropicProvider{
		baseProvider: bp,
		models: []ModelInfo{
			{ID: "claude-sonnet-4-20250514", Name: "Claude Sonnet 4", Provider: ProviderAnthropic, ContextWindow: 200000, MaxOutput: 64000, SupportsVision: true, SupportsStreaming: true, SupportsThinking: true, CacheSupport: true, CostPerInputMT: 3.00, CostPerOutputMT: 15.00, CostPerCacheReadMT: 0.30, CostPerCacheCreationMT: 3.75, Tags: []string{"coding", "analysis", "fast"}},
			{ID: "claude-opus-4-20250514", Name: "Claude Opus 4", Provider: ProviderAnthropic, ContextWindow: 200000, MaxOutput: 64000, SupportsVision: true, SupportsStreaming: true, SupportsThinking: true, CacheSupport: true, CostPerInputMT: 15.00, CostPerOutputMT: 75.00, CostPerCacheReadMT: 1.50, CostPerCacheCreationMT: 18.75, Tags: []string{"coding", "analysis", "smart"}},
			{ID: "claude-3-5-sonnet-20241022", Name: "Claude 3.5 Sonnet", Provider: ProviderAnthropic, ContextWindow: 200000, MaxOutput: 8192, SupportsVision: true, SupportsStreaming: true, SupportsThinking: true, CacheSupport: true, CostPerInputMT: 3.00, CostPerOutputMT: 15.00, CostPerCacheReadMT: 0.30, CostPerCacheCreationMT: 3.75, Tags: []string{"coding", "analysis"}},
			{ID: "claude-3-7-sonnet-20250219", Name: "Claude 3.7 Sonnet", Provider: ProviderAnthropic, ContextWindow: 200000, MaxOutput: 64000, SupportsVision: true, SupportsStreaming: true, SupportsThinking: true, CacheSupport: true, CostPerInputMT: 3.00, CostPerOutputMT: 15.00, CostPerCacheReadMT: 0.30, CostPerCacheCreationMT: 3.75, Tags: []string{"coding", "extended-thinking"}},
			{ID: "claude-3-5-haiku-20241022", Name: "Claude 3.5 Haiku", Provider: ProviderAnthropic, ContextWindow: 200000, MaxOutput: 8192, SupportsVision: true, SupportsStreaming: true, SupportsThinking: false, CacheSupport: true, CostPerInputMT: 0.80, CostPerOutputMT: 4.00, CostPerCacheReadMT: 0.08, CostPerCacheCreationMT: 1.00, Tags: []string{"fast", "cheap"}},
		},
	}
	return p, nil
}

func (p *anthropicProvider) Type() ProviderType { return ProviderAnthropic }
func (p *anthropicProvider) Name() string       { return "Anthropic" }
func (p *anthropicProvider) Configure(config ProviderConfig) error {
	p.baseProvider.config = config
	return nil
}
func (p *anthropicProvider) GetModels() ([]ModelInfo, error) { return p.models, nil }
func (p *anthropicProvider) MatchesModel(pattern string) bool {
	return strings.Contains(strings.ToLower(pattern), "claude")
}
func (p *anthropicProvider) EstimateTokens(text string) int {
	// ~3.6 chars per token for Claude
	return int(float64(len(text)) / 3.6)
}
func (p *anthropicProvider) StreamCompletion(ctx context.Context, req CompletionRequest) (chan StreamChunk, error) {
	return nil, fmt.Errorf("anthropic streaming not yet implemented")
}
func (p *anthropicProvider) CreateCompletion(ctx context.Context, req CompletionRequest) (*CompletionResult, error) {
	return nil, fmt.Errorf("anthropic completion not yet implemented")
}

// OpenAI provider implementation.
type openaiProvider struct {
	baseProvider
	models []ModelInfo
	mu     sync.RWMutex
}

func newOpenAIProvider(config ProviderConfig) (Provider, error) {
	bp := newBaseProvider(config)
	if config.BaseURL == "" {
		config.BaseURL = "https://api.openai.com/v1"
	}
	p := &openaiProvider{
		baseProvider: bp,
		models: []ModelInfo{
			{ID: "gpt-4.1", Name: "GPT-4.1", Provider: ProviderOpenAI, ContextWindow: 1047576, MaxOutput: 16384, SupportsVision: true, SupportsStreaming: true, SupportsThinking: false, CacheSupport: true, CostPerInputMT: 2.00, CostPerOutputMT: 8.00, CostPerCacheReadMT: 0.25, CostPerCacheCreationMT: 2.50, Tags: []string{"general", "fast"}},
			{ID: "gpt-4o", Name: "GPT-4o", Provider: ProviderOpenAI, ContextWindow: 128000, MaxOutput: 16384, SupportsVision: true, SupportsStreaming: true, SupportsThinking: false, CacheSupport: true, CostPerInputMT: 2.50, CostPerOutputMT: 10.00, CostPerCacheReadMT: 1.25, CostPerCacheCreationMT: 2.50, Tags: []string{"general"}},
			{ID: "o3", Name: "O3", Provider: ProviderOpenAI, ContextWindow: 200000, MaxOutput: 100000, SupportsVision: true, SupportsStreaming: true, SupportsThinking: true, CacheSupport: true, CostPerInputMT: 2.00, CostPerOutputMT: 8.00, CostPerCacheReadMT: 0.50, CostPerCacheCreationMT: 2.00, Tags: []string{"reasoning", "coding"}},
			{ID: "o4-mini", Name: "O4 Mini", Provider: ProviderOpenAI, ContextWindow: 200000, MaxOutput: 100000, SupportsVision: true, SupportsStreaming: true, SupportsThinking: true, CacheSupport: true, CostPerInputMT: 1.10, CostPerOutputMT: 4.40, CostPerCacheReadMT: 0.275, CostPerCacheCreationMT: 1.10, Tags: []string{"reasoning"}},
			{ID: "codex-mini", Name: "Codex Mini", Provider: ProviderOpenAI, ContextWindow: 200000, MaxOutput: 100000, SupportsVision: false, SupportsStreaming: true, SupportsThinking: true, CacheSupport: true, CostPerInputMT: 1.50, CostPerOutputMT: 6.00, CostPerCacheReadMT: 0.15, CostPerCacheCreationMT: 1.50, Tags: []string{"coding"}},
		},
	}
	return p, nil
}

func (p *openaiProvider) Type() ProviderType { return ProviderOpenAI }
func (p *openaiProvider) Name() string       { return "OpenAI" }
func (p *openaiProvider) Configure(config ProviderConfig) error {
	p.baseProvider.config = config
	return nil
}
func (p *openaiProvider) GetModels() ([]ModelInfo, error) { return p.models, nil }
func (p *openaiProvider) MatchesModel(pattern string) bool {
	return strings.Contains(strings.ToLower(pattern), "gpt") || strings.Contains(strings.ToLower(pattern), "o3") || strings.Contains(strings.ToLower(pattern), "o4")
}
func (p *openaiProvider) EstimateTokens(text string) int {
	// ~3.5 chars per token for OpenAI
	return int(float64(len(text)) / 3.5)
}
func (p *openaiProvider) StreamCompletion(ctx context.Context, req CompletionRequest) (chan StreamChunk, error) {
	return nil, fmt.Errorf("openai streaming not yet implemented")
}
func (p *openaiProvider) CreateCompletion(ctx context.Context, req CompletionRequest) (*CompletionResult, error) {
	return nil, fmt.Errorf("openai completion not yet implemented")
}

// Google provider implementation.
type googleProvider struct {
	baseProvider
	models []ModelInfo
}

func newGoogleProvider(config ProviderConfig) (Provider, error) {
	bp := newBaseProvider(config)
	if config.BaseURL == "" {
		config.BaseURL = "https://generativelanguage.googleapis.com/v1beta"
	}
	p := &googleProvider{
		baseProvider: bp,
		models: []ModelInfo{
			{ID: "gemini-2.5-pro", Name: "Gemini 2.5 Pro", Provider: ProviderGoogle, ContextWindow: 1000000, MaxOutput: 64000, SupportsVision: true, SupportsStreaming: true, SupportsThinking: true, CacheSupport: true, CostPerInputMT: 2.50, CostPerOutputMT: 15.00, CostPerCacheReadMT: 0.55, CostPerCacheCreationMT: 3.50, Tags: []string{"reasoning", "coding"}},
			{ID: "gemini-2.5-flash", Name: "Gemini 2.5 Flash", Provider: ProviderGoogle, ContextWindow: 1000000, MaxOutput: 64000, SupportsVision: true, SupportsStreaming: true, SupportsThinking: true, CacheSupport: true, CostPerInputMT: 0.15, CostPerOutputMT: 3.50, CostPerCacheReadMT: 0.25, CostPerCacheCreationMT: 1.75, Tags: []string{"fast", "coding"}},
			{ID: "gemini-2.0-flash", Name: "Gemini 2.0 Flash", Provider: ProviderGoogle, ContextWindow: 1000000, MaxOutput: 8192, SupportsVision: true, SupportsStreaming: true, SupportsThinking: false, CacheSupport: true, CostPerInputMT: 10.00, CostPerOutputMT: 2.50, CostPerCacheReadMT: 0.20, CostPerCacheCreationMT: 1.25, Tags: []string{"fast", "free"}},
		},
	}
	return p, nil
}

func (p *googleProvider) Type() ProviderType { return ProviderGoogle }
func (p *googleProvider) Name() string       { return "Google" }
func (p *googleProvider) Configure(config ProviderConfig) error {
	p.baseProvider.config = config
	return nil
}
func (p *googleProvider) GetModels() ([]ModelInfo, error) { return p.models, nil }
func (p *googleProvider) MatchesModel(pattern string) bool {
	return strings.Contains(strings.ToLower(pattern), "gemini")
}
func (p *googleProvider) EstimateTokens(text string) int {
	return len(strings.Fields(text)) * 4 / 3 // ~1.33 tokens per word
}
func (p *googleProvider) StreamCompletion(ctx context.Context, req CompletionRequest) (chan StreamChunk, error) {
	return nil, fmt.Errorf("google streaming not yet implemented")
}
func (p *googleProvider) CreateCompletion(ctx context.Context, req CompletionRequest) (*CompletionResult, error) {
	return nil, fmt.Errorf("google completion not yet implemented")
}

// Mistral provider implementation.
type mistralProvider struct {
	baseProvider
	models []ModelInfo
}

func newMistralProvider(config ProviderConfig) (Provider, error) {
	bp := newBaseProvider(config)
	if config.BaseURL == "" {
		config.BaseURL = "https://api.mistral.ai/v1"
	}
	p := &mistralProvider{
		baseProvider: bp,
		models: []ModelInfo{
			{ID: "mistral-large-latest", Name: "Mistral Large", Provider: ProviderMistral, ContextWindow: 131072, MaxOutput: 128000, SupportsVision: true, SupportsStreaming: true, SupportsThinking: false, CacheSupport: false, CostPerInputMT: 2.00, CostPerOutputMT: 6.00, Tags: []string{"general", "coding"}},
		},
	}
	return p, nil
}

func (p *mistralProvider) Type() ProviderType { return ProviderMistral }
func (p *mistralProvider) Name() string       { return "Mistral" }
func (p *mistralProvider) Configure(config ProviderConfig) error {
	p.baseProvider.config = config
	return nil
}
func (p *mistralProvider) GetModels() ([]ModelInfo, error) { return p.models, nil }
func (p *mistralProvider) MatchesModel(pattern string) bool {
	l := strings.ToLower(pattern)
	return strings.Contains(l, "mistral") || strings.Contains(l, "mixtral")
}
func (p *mistralProvider) EstimateTokens(text string) int {
	return int(float64(len(text)) / 4.0)
}
func (p *mistralProvider) StreamCompletion(ctx context.Context, req CompletionRequest) (chan StreamChunk, error) {
	return nil, fmt.Errorf("mistral streaming not yet implemented")
}
func (p *mistralProvider) CreateCompletion(ctx context.Context, req CompletionRequest) (*CompletionResult, error) {
	return nil, fmt.Errorf("mistral completion not yet implemented")
}

// Generic provider for compatible OpenAI API endpoints.
type genericProvider struct {
	baseProvider
	models []ModelInfo
}

func newGenericProvider(config ProviderConfig) (Provider, error) {
	bp := newBaseProvider(config)
	return &genericProvider{
		baseProvider: bp,
		models: []ModelInfo{
			{ID: "custom", Name: "Custom Model", Provider: ProviderCustom, ContextWindow: 128000, MaxOutput: 4096, SupportsStreaming: true, SupportsVision: false, SupportsThinking: false, CacheSupport: false},
		},
	}, nil
}

func (p *genericProvider) Type() ProviderType { return ProviderCustom }
func (p *genericProvider) Name() string       { return "Custom" }
func (p *genericProvider) Configure(config ProviderConfig) error {
	p.baseProvider.config = config
	return nil
}
func (p *genericProvider) GetModels() ([]ModelInfo, error)  { return p.models, nil }
func (p *genericProvider) MatchesModel(pattern string) bool { return true }
func (p *genericProvider) EstimateTokens(text string) int {
	return int(float64(len(text)) / 4.0)
}
func (p *genericProvider) StreamCompletion(ctx context.Context, req CompletionRequest) (chan StreamChunk, error) {
	return nil, fmt.Errorf("generic streaming not yet implemented")
}
func (p *genericProvider) CreateCompletion(ctx context.Context, req CompletionRequest) (*CompletionResult, error) {
	return nil, fmt.Errorf("generic completion not yet implemented")
}
