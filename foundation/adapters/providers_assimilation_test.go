package adapters

import (
	"testing"

	"github.com/robertpelloni/hyperharness/config"
)

func TestDetectAvailableProvidersAndSortStrings(t *testing.T) {
	cfg := &config.Config{Provider: "ollama"}
	providers := detectAvailableProviders(cfg)
	if len(providers) == 0 {
		t.Fatalf("expected at least one provider, got %#v", providers)
	}
	foundOllama := false
	for i, provider := range providers {
		if provider == "ollama" {
			foundOllama = true
		}
		if i > 0 && providers[i-1] > provider {
			t.Fatalf("expected sorted providers, got %#v", providers)
		}
	}
	if !foundOllama {
		t.Fatalf("expected ollama detection, got %#v", providers)
	}

	values := []string{"zeta", "alpha", "beta"}
	sortStrings(values)
	if values[0] != "alpha" || values[1] != "beta" || values[2] != "zeta" {
		t.Fatalf("unexpected sorted values: %#v", values)
	}
}
