package agents

import (
	"context"
	"strings"
	"testing"
)

func TestLastProviderMessageContent(t *testing.T) {
	if got := lastProviderMessageContent(nil); got != "" {
		t.Fatalf("expected empty content, got %q", got)
	}
	if got := lastProviderMessageContent([]Message{{Content: "a"}, {Content: "b"}}); got != "b" {
		t.Fatalf("expected last content 'b', got %q", got)
	}
}

func TestDefaultLegacyToolsReturnsFreshCopy(t *testing.T) {
	first := defaultLegacyTools()
	second := defaultLegacyTools()
	if len(first) != 3 || len(second) != 3 {
		t.Fatalf("expected 3 legacy tools, got %#v / %#v", first, second)
	}
	first[0].Name = "modified"
	if second[0].Name != "apply_search_replace" {
		t.Fatalf("expected fresh copy, got %#v", second)
	}
}

func TestGeminiProviderValidationAndDefaults(t *testing.T) {
	var provider *GeminiBorgProvider
	if _, err := provider.Chat(context.Background(), nil, nil); err == nil || !strings.Contains(err.Error(), "provider is required") {
		t.Fatalf("expected nil provider error, got %v", err)
	}
	if err := provider.Stream(context.Background(), nil, nil, make(chan string)); err == nil || !strings.Contains(err.Error(), "provider is required") {
		t.Fatalf("expected nil provider stream error, got %v", err)
	}

	provider = &GeminiBorgProvider{}
	if provider.GetModelName() != defaultGeminiModel {
		t.Fatalf("expected default model name %q, got %q", defaultGeminiModel, provider.GetModelName())
	}
	if err := provider.Stream(context.Background(), nil, nil, nil); err == nil || !strings.Contains(err.Error(), "chunkChan is required") {
		t.Fatalf("expected nil channel error, got %v", err)
	}
	msg, err := provider.Chat(context.Background(), []Message{{Role: RoleUser, Content: "analyze repo"}}, nil)
	if err != nil {
		t.Fatalf("chat failed: %v", err)
	}
	if !strings.Contains(msg.Content, "Gemini Parity Engaged") || !strings.Contains(msg.Content, defaultGeminiModel) {
		t.Fatalf("unexpected chat content: %#v", msg)
	}
}

func TestGeminiProviderStreamAndFactory(t *testing.T) {
	provider := &GeminiBorgProvider{Model: "gemini-test"}
	chunks := make(chan string, 2)
	if err := provider.Stream(context.Background(), nil, nil, chunks); err != nil {
		t.Fatalf("stream failed: %v", err)
	}
	collected := []string{}
	for chunk := range chunks {
		collected = append(collected, chunk)
	}
	if len(collected) != 2 || collected[1] != "[gemini-test] Ready." {
		t.Fatalf("unexpected stream chunks: %#v", collected)
	}

	wrapped := NewGeminiBorgProvider()
	proxy, ok := wrapped.(*DisclosureProxy)
	if !ok {
		t.Fatalf("expected disclosure proxy wrapper, got %T", wrapped)
	}
	if proxy.BaseProvider == nil || len(proxy.HiddenTools) != 3 {
		t.Fatalf("expected wrapped base provider and 3 hidden tools, got %#v", proxy)
	}
}

func TestDefaultProviderValidationAndBehavior(t *testing.T) {
	var provider *DefaultProvider
	if _, err := provider.Chat(context.Background(), nil, nil); err == nil || !strings.Contains(err.Error(), "provider is required") {
		t.Fatalf("expected nil provider error, got %v", err)
	}
	if err := provider.Stream(context.Background(), nil, nil, make(chan string)); err == nil || !strings.Contains(err.Error(), "provider is required") {
		t.Fatalf("expected nil provider stream error, got %v", err)
	}

	provider = &DefaultProvider{}
	if err := provider.Stream(context.Background(), nil, nil, nil); err == nil || !strings.Contains(err.Error(), "chunkChan is required") {
		t.Fatalf("expected nil chunk channel error, got %v", err)
	}
	msg, err := provider.Chat(context.Background(), []Message{{Role: RoleUser, Content: "hello"}}, nil)
	if err != nil {
		t.Fatalf("chat failed: %v", err)
	}
	if !strings.Contains(msg.Content, "Native Go") {
		t.Fatalf("unexpected stub chat content: %#v", msg)
	}
	if provider.GetModelName() != "hypercode-native-stub-1.0" {
		t.Fatalf("unexpected model name: %q", provider.GetModelName())
	}
}
