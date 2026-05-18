package providers

import (
	"encoding/json"
	"os"
	"testing"
)

func TestProviderTypes(t *testing.T) {
	types := map[ProviderType]string{
		ProviderAnthropic:   "anthropic",
		ProviderOpenAI:      "openai",
		ProviderGoogle:      "google",
		ProviderAzure:       "azure",
		ProviderGroq:        "groq",
		ProviderCerebras:    "cerebras",
		ProviderXAI:         "xai",
		ProviderMistral:     "mistral",
		ProviderOpenRouter:  "openrouter",
		ProviderBedrock:     "bedrock",
		ProviderVertex:      "vertex",
		ProviderMiniMax:     "minimax",
		ProviderHuggingFace: "huggingface",
		ProviderOpenCode:    "opencode",
		ProviderCustom:      "custom",
	}
	for pt, expected := range types {
		if string(pt) != expected {
			t.Errorf("provider type %q != %q", pt, expected)
		}
	}
}

func TestThinkingLevels(t *testing.T) {
	levels := map[ThinkingLevel]string{
		ThinkingOff:     "off",
		ThinkingMinimal: "minimal",
		ThinkingLow:     "low",
		ThinkingMedium:  "medium",
		ThinkingHigh:    "high",
		ThinkingXHigh:   "xhigh",
	}
	for tl, expected := range levels {
		if string(tl) != expected {
			t.Errorf("thinking level %q != %q", tl, expected)
		}
	}
}

func TestDefaultThinkingBudgets(t *testing.T) {
	if DefaultThinkingBudgets[ThinkingOff] != 0 {
		t.Errorf("off budget should be 0")
	}
	if DefaultThinkingBudgets[ThinkingHigh] <= DefaultThinkingBudgets[ThinkingMedium] {
		t.Errorf("high budget should exceed medium")
	}
}

func TestProviderConfigJSONRoundtrip(t *testing.T) {
	config := ProviderConfig{
		Type:       ProviderOpenAI,
		APIKey:     "sk-test123",
		BaseURL:    "https://api.openai.com/v1",
		MaxRetries: 3,
		Transport:  TransportSSE,
	}

	data, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}

	var decoded ProviderConfig
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}

	if decoded.Type != ProviderOpenAI {
		t.Errorf("type: %s", decoded.Type)
	}
	if decoded.MaxRetries != 3 {
		t.Errorf("max retries: %d", decoded.MaxRetries)
	}
}

func TestRegistryNew(t *testing.T) {
	registry := NewRegistry()
	if registry == nil {
		t.Fatal("registry should not be nil")
	}
}

func TestInferProviderFromModel(t *testing.T) {
	tests := []struct {
		model    string
		expected ProviderType
	}{
		{"gpt-4o", ProviderOpenAI},
		{"gpt-3.5-turbo", ProviderOpenAI},
		{"claude-3-opus", ProviderAnthropic},
		{"claude-3.5-sonnet", ProviderAnthropic},
		{"gemini-2.0-flash", ProviderGoogle},
		{"gemini-1.5-pro", ProviderGoogle},
		{"llama-3.1-8b", ProviderGroq},
		{"deepseek-r1", ProviderOpenAI},
	}

	for _, tc := range tests {
		result := InferProviderFromModel(tc.model)
		if result != tc.expected {
			t.Errorf("InferProviderFromModel(%q) = %q, want %q", tc.model, result, tc.expected)
		}
	}
}

func TestResolveModel(t *testing.T) {
	pt, _, thinking := ResolveModel("anthropic/claude-3-opus:high")
	if pt != ProviderAnthropic {
		t.Errorf("provider: %s", pt)
	}
	if thinking != ThinkingHigh {
		t.Errorf("thinking: %s", thinking)
	}
}

func TestModelInfo(t *testing.T) {
	info := ModelInfo{
		ID:            "gpt-4o",
		Name:          "GPT-4o",
		Provider:      ProviderOpenAI,
		ContextWindow: 128000,
		MaxOutput:     16384,
	}

	if info.ID != "gpt-4o" {
		t.Errorf("ID: %s", info.ID)
	}
	if info.ContextWindow != 128000 {
		t.Errorf("context window: %d", info.ContextWindow)
	}
}

func TestTransportTypes(t *testing.T) {
	transports := map[Transport]string{
		TransportSSE:       "sse",
		TransportWebsocket: "websocket",
		TransportAuto:      "auto",
	}
	for tr, expected := range transports {
		if string(tr) != expected {
			t.Errorf("transport %q != %q", tr, expected)
		}
	}
}

func TestMessageTypes(t *testing.T) {
	msg := Message{
		Role: "user",
		Content: []ContentBlock{{Type: "text", Text: "Hello"}},
	}

	data, _ := json.Marshal(msg)
	var decoded Message
	json.Unmarshal(data, &decoded)

	if decoded.Role != "user" {
		t.Errorf("role: %s", decoded.Role)
	}
}

func TestTokenUsage(t *testing.T) {
	usage := TokenUsage{
		InputTokens:     100,
		OutputTokens:    50,
		CacheReadTokens: 20,
	}

	if usage.InputTokens != 100 {
		t.Errorf("input tokens: %d", usage.InputTokens)
	}
}

func TestContentBlock(t *testing.T) {
	block := ContentBlock{
		Type: "text",
		Text: "Hello world",
	}

	data, _ := json.Marshal(block)
	var decoded ContentBlock
	json.Unmarshal(data, &decoded)

	if decoded.Type != "text" {
		t.Errorf("type: %s", decoded.Type)
	}
	if decoded.Text != "Hello world" {
		t.Errorf("text: %s", decoded.Text)
	}
}

func TestCatalogEntries(t *testing.T) {
	if len(Catalog) == 0 {
		t.Error("catalog should have entries")
	}
	for _, entry := range Catalog {
		if entry.Provider == "" {
			t.Error("provider should have a name")
		}
		if entry.DefaultModel == "" {
			t.Errorf("%s should have a default model", entry.Provider)
		}
	}
}

func TestFindProviderForModel(t *testing.T) {
	p := FindProviderForModel("gpt-4o")
	if p == nil {
		t.Error("should find provider for gpt-4o")
	}
	if p.Provider != "openai" {
		t.Errorf("expected openai, got %s", p.Provider)
	}

	p = FindProviderForModel("nonexistent-model")
	if p != nil {
		t.Error("should not find nonexistent model")
	}
}

func TestFindCatalogModelInfo(t *testing.T) {
	info := FindCatalogModelInfo("gpt-4o")
	if info == nil {
		t.Error("should find model info for gpt-4o")
	}
	if info.Speed != "medium" {
		t.Errorf("gpt-4o speed should be medium, got %s", info.Speed)
	}

	info = FindCatalogModelInfo("claude-sonnet-4-20250514")
	if info == nil {
		t.Error("should find Claude Sonnet 4")
	}
}

func TestStatusChecker(t *testing.T) {
	sc := NewStatusChecker()
	statuses := sc.CheckAll()
	if len(statuses) == 0 {
		t.Error("should have statuses")
	}

	for _, s := range statuses {
		if s.Provider == "" {
			t.Error("status should have provider name")
		}
	}
}

func TestStatusCheckerLocalProviders(t *testing.T) {
	sc := NewStatusChecker()
	statuses := sc.CheckAll()

	for _, s := range statuses {
		if s.Provider == "ollama" || s.Provider == "lmstudio" {
			if !s.Configured {
				t.Errorf("%s should be configured (local)", s.Provider)
			}
			if !s.Authenticated {
				t.Errorf("%s should be authenticated (local)", s.Provider)
			}
		}
	}
}

func TestSelectModelForTask(t *testing.T) {
	// Clear env vars for deterministic test
	for _, key := range []string{"ANTHROPIC_API_KEY", "GOOGLE_API_KEY", "OPENAI_API_KEY"} {
		os.Unsetenv(key)
	}

	// With no keys, should find local providers
	entry := SelectModelForTask("local-inference")
	if entry == nil {
		t.Error("should find local provider for local-inference")
	}
}

func TestFormatCatalogSummary(t *testing.T) {
	sc := NewStatusChecker()
	statuses := sc.CheckAll()
	summary := FormatCatalogSummary(statuses)
	if summary == "" {
		t.Error("summary should not be empty")
	}
	if len(summary) < 50 {
		t.Error("summary should be substantial")
	}
}
