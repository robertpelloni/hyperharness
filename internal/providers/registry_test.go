package providers

import (
	"encoding/json"
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
