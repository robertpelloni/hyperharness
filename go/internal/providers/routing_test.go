package providers

import "testing"

func TestBuildRoutingSummaryReflectsConfiguredProviders(t *testing.T) {
	summary := BuildRoutingSummary([]Status{
		{Provider: "google", Configured: true, Authenticated: true},
		{Provider: "anthropic", Configured: true, Authenticated: true},
		{Provider: "openai", Configured: false, Authenticated: false},
		{Provider: "deepseek", Configured: false, Authenticated: false},
		{Provider: "openrouter", Configured: false, Authenticated: false},
		{Provider: "github-copilot", Configured: true, Authenticated: true},
	})

	if summary.DefaultStrategy != "best" {
		t.Fatalf("expected default strategy best, got %q", summary.DefaultStrategy)
	}
	if len(summary.Tasks) != 6 {
		t.Fatalf("expected 6 task summaries, got %d", len(summary.Tasks))
	}
	if summary.Tasks[0].TaskType != "coding" {
		t.Fatalf("expected coding task first, got %q", summary.Tasks[0].TaskType)
	}
	if !summary.Tasks[0].Candidates[0].Configured {
		t.Fatalf("expected first coding candidate to be configured")
	}
}
