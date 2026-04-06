package adapters

import "testing"

func TestProviderRoutingHelpers(t *testing.T) {
	if got := normalizedPreference(" low "); got != "budget" {
		t.Fatalf("expected budget normalization, got %q", got)
	}
	if got := normalizedPreference("best"); got != "quality" {
		t.Fatalf("expected quality normalization, got %q", got)
	}
	if got := normalizedPreference("other"); got != "" {
		t.Fatalf("expected empty normalization, got %q", got)
	}

	if got := normalizedTaskType(" refactor "); got != "coding" {
		t.Fatalf("expected coding normalization, got %q", got)
	}
	if got := normalizedTaskType("research"); got != "analysis" {
		t.Fatalf("expected analysis normalization, got %q", got)
	}
	if got := normalizedTaskType("unknown"); got != "" {
		t.Fatalf("expected empty task type normalization, got %q", got)
	}

	if got := firstOr([]string{"model-a"}, "fallback"); got != "model-a" {
		t.Fatalf("expected first value, got %q", got)
	}
	if got := firstOr(nil, "fallback"); got != "fallback" {
		t.Fatalf("expected fallback, got %q", got)
	}

	if got := defaultModelForProvider("ollama", []string{"llama3.2"}); got != "llama3.2" {
		t.Fatalf("unexpected ollama default: %q", got)
	}
	if got := defaultModelForProvider("ollama", nil); got != defaultOllamaModel {
		t.Fatalf("unexpected ollama fallback: %q", got)
	}
	if got := defaultModelForProvider("google", nil); got != defaultGoogleModel {
		t.Fatalf("unexpected google default: %q", got)
	}
	if got := defaultModelForProvider("openai", nil); got != defaultOpenAIModel {
		t.Fatalf("unexpected openai default: %q", got)
	}
}
