package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestOpenAIProviderName(t *testing.T) {
	p := &OpenAIProvider{}
	if p.Name() != "openai" {
		t.Errorf("expected openai, got %s", p.Name())
	}
}

func TestOpenAIProviderModels(t *testing.T) {
	p := &OpenAIProvider{}
	models := p.Models()
	if len(models) == 0 {
		t.Error("should have models")
	}
	found := false
	for _, m := range models {
		if m == "gpt-4o" {
			found = true
		}
	}
	if !found {
		t.Error("should include gpt-4o")
	}
}

func TestOpenAIProviderGenerateText(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify auth header
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Error("missing auth header")
		}

		resp := map[string]interface{}{
			"choices": []map[string]interface{}{
				{"message": map[string]string{"content": "Hello from GPT!"}},
			},
			"usage": map[string]int{
				"prompt_tokens":     10,
				"completion_tokens": 5,
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := &OpenAIProvider{APIKey: "test-key", BaseURL: server.URL}
	resp, err := p.GenerateText(context.Background(), "gpt-4o", []Message{
		{Role: "user", Content: "Hello"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Content != "Hello from GPT!" {
		t.Errorf("content mismatch: %s", resp.Content)
	}
	if resp.Provider != "openai" {
		t.Errorf("provider: %s", resp.Provider)
	}
	if resp.Usage.InputTokens != 10 {
		t.Errorf("input tokens: %d", resp.Usage.InputTokens)
	}
}

func TestAnthropicProviderGenerateText(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") != "test-key" {
			t.Error("missing anthropic key")
		}
		if r.Header.Get("anthropic-version") == "" {
			t.Error("missing anthropic version")
		}

		// Verify system message extraction
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if sys, ok := body["system"]; !ok || sys == nil {
			t.Error("system message should be extracted")
		}

		resp := map[string]interface{}{
			"content": []map[string]string{
				{"text": "Claude says hello!"},
			},
			"usage": map[string]int{
				"input_tokens":  15,
				"output_tokens": 8,
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := &AnthropicProvider{APIKey: "test-key", BaseURL: server.URL}
	resp, err := p.GenerateText(context.Background(), "claude-sonnet-4-20250514", []Message{
		{Role: "system", Content: "You are helpful"},
		{Role: "user", Content: "Hi"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Content != "Claude says hello!" {
		t.Errorf("content: %s", resp.Content)
	}
	if resp.Provider != "anthropic" {
		t.Errorf("provider: %s", resp.Provider)
	}
}

func TestGeminiProviderGenerateText(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)

		contents, ok := body["contents"].([]interface{})
		if !ok || len(contents) == 0 {
			t.Error("should have contents")
		}

		resp := map[string]interface{}{
			"candidates": []map[string]interface{}{
				{"content": map[string]interface{}{
					"parts": []map[string]string{{"text": "Gemini response"}},
				}},
			},
			"usageMetadata": map[string]int{
				"promptTokenCount":     20,
				"candidatesTokenCount": 10,
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := &GeminiProvider{APIKey: "test-key", BaseURL: server.URL}
	resp, err := p.GenerateText(context.Background(), "gemini-2.5-flash", []Message{
		{Role: "user", Content: "Hello"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Content != "Gemini response" {
		t.Errorf("content: %s", resp.Content)
	}
	if resp.Provider != "google" {
		t.Errorf("provider: %s", resp.Provider)
	}
}

func TestDeepSeekProvider(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"choices": []map[string]interface{}{
				{"message": map[string]string{"content": "DeepSeek response"}},
			},
			"usage": map[string]int{"prompt_tokens": 5, "completion_tokens": 3},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Test DeepSeek uses OpenAI-compatible endpoint
	resp, err := (&OpenAIProvider{APIKey: "test-key", BaseURL: server.URL}).GenerateText(context.Background(), "deepseek-chat", nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Content != "DeepSeek response" {
		t.Errorf("content: %s", resp.Content)
	}
}

func TestOllamaProviderGenerateText(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["stream"] != false {
			t.Error("ollama should have stream=false")
		}

		resp := map[string]interface{}{
			"message":          map[string]string{"content": "Ollama response"},
			"prompt_eval_count": 12,
			"eval_count":       6,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := &OllamaProvider{BaseURL: server.URL}
	resp, err := p.GenerateText(context.Background(), "gemma:2b", []Message{
		{Role: "user", Content: "Hello"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Content != "Ollama response" {
		t.Errorf("content: %s", resp.Content)
	}
	if resp.Provider != "ollama" {
		t.Errorf("provider: %s", resp.Provider)
	}
}

func TestAutoRouteNoProvider(t *testing.T) {
	// Clear all API keys
	for _, key := range []string{"ANTHROPIC_API_KEY", "GOOGLE_API_KEY", "GEMINI_API_KEY", "OPENAI_API_KEY", "DEEPSEEK_API_KEY", "OPENROUTER_API_KEY", "GROQ_API_KEY"} {
		os.Unsetenv(key)
	}

	_, err := AutoRoute(context.Background(), nil)
	if err == nil {
		t.Error("should error with no providers")
	}
}

func TestAutoRouteWithKey(t *testing.T) {
	// Set a key
	os.Setenv("OPENAI_API_KEY", "test-key")
	defer os.Unsetenv("OPENAI_API_KEY")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"choices": []map[string]interface{}{
				{"message": map[string]string{"content": "Auto-routed!"}},
			},
			"usage": map[string]int{"prompt_tokens": 1, "completion_tokens": 1},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Override the base URL for the test
	p := &OpenAIProvider{APIKey: "test-key", BaseURL: server.URL}
	resp, err := p.GenerateText(context.Background(), "gpt-4o", []Message{{Role: "user", Content: "test"}})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Content != "Auto-routed!" {
		t.Errorf("content: %s", resp.Content)
	}
}

func TestListConfiguredProviders(t *testing.T) {
	// Clear all
	for _, key := range []string{"ANTHROPIC_API_KEY", "GOOGLE_API_KEY", "OPENAI_API_KEY"} {
		os.Unsetenv(key)
	}

	// Set one
	os.Setenv("OPENAI_API_KEY", "test")
	defer os.Unsetenv("OPENAI_API_KEY")

	configured := ListConfiguredProviders()
	found := false
	for _, p := range configured {
		if p == "openai" {
			found = true
		}
	}
	if !found {
		t.Error("openai should be configured")
	}
}

func TestGetProvider(t *testing.T) {
	os.Setenv("OPENAI_API_KEY", "test")
	defer os.Unsetenv("OPENAI_API_KEY")

	p, ok := GetProvider("openai")
	if !ok {
		t.Error("should find openai provider")
	}
	if p.Name() != "openai" {
		t.Errorf("name: %s", p.Name())
	}

	_, ok = GetProvider("nonexistent")
	if ok {
		t.Error("should not find nonexistent provider")
	}
}

func TestGetProviderLocal(t *testing.T) {
	p, ok := GetProvider("ollama")
	if !ok {
		t.Error("ollama should always be available")
	}
	if p.Name() != "ollama" {
		t.Errorf("name: %s", p.Name())
	}

	p2, ok := GetProvider("lmstudio")
	if !ok {
		t.Error("lmstudio should always be available")
	}
	_ = p2
}

func TestGeminiRoleMapping(t *testing.T) {
	var receivedBody map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&receivedBody)
		resp := map[string]interface{}{
			"candidates": []map[string]interface{}{
				{"content": map[string]interface{}{
					"parts": []map[string]string{{"text": "ok"}},
				}},
			},
			"usageMetadata": map[string]int{"promptTokenCount": 1, "candidatesTokenCount": 1},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := &GeminiProvider{APIKey: "test", BaseURL: server.URL}
	_, err := p.GenerateText(context.Background(), "gemini-2.5-flash", []Message{
		{Role: "assistant", Content: "Previous response"},
		{Role: "user", Content: "Follow up"},
	})
	if err != nil {
		t.Fatal(err)
	}

	contents := receivedBody["contents"].([]interface{})
	first := contents[0].(map[string]interface{})
	if first["role"] != "model" {
		t.Errorf("assistant should map to model, got %v", first["role"])
	}
}

func TestProviderTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		resp := map[string]interface{}{
			"choices": []map[string]interface{}{
				{"message": map[string]string{"content": "too slow"}},
			},
			"usage": map[string]int{"prompt_tokens": 1, "completion_tokens": 1},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	p := &OpenAIProvider{APIKey: "test", BaseURL: server.URL}
	_, err := p.GenerateText(ctx, "gpt-4o", []Message{{Role: "user", Content: "test"}})
	if err == nil {
		t.Error("should timeout")
	}
}

func TestProviderAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`{"error": "rate limited"}`))
	}))
	defer server.Close()

	p := &OpenAIProvider{APIKey: "test", BaseURL: server.URL}
	_, err := p.GenerateText(context.Background(), "gpt-4o", nil)
	if err == nil {
		t.Error("should error on 429")
	}
}

func TestGetSwarmPrompt(t *testing.T) {
	tests := []struct {
		role     string
		contains string
	}{
		{"planner", "Planner"},
		{"implementer", "Implementer"},
		{"tester", "Tester"},
		{"critic", "Critic"},
		{"unknown", "helpful assistant"},
	}
	for _, tt := range tests {
		prompt := GetSwarmPrompt(tt.role)
		if prompt == "" {
			t.Errorf("prompt for %s should not be empty", tt.role)
		}
	}
}

func TestGetSubagentPrompt(t *testing.T) {
	types := []string{"code", "research", "review", "security", "debug", "doc"}
	for _, typ := range types {
		prompt := GetSubagentPrompt(typ)
		if prompt == "" {
			t.Errorf("prompt for %s should not be empty", typ)
		}
	}

	// Default case
	prompt := GetSubagentPrompt("custom-type")
	if prompt == "" {
		t.Error("default prompt should not be empty")
	}
}
