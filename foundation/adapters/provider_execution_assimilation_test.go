package adapters

import (
	"strings"
	"testing"
)

func TestExecutionHelpers(t *testing.T) {
	if got := normalizeTaskType(" coding ", "ignored"); got != "coding" {
		t.Fatalf("expected explicit task type, got %q", got)
	}
	if got := normalizeTaskType("", "please analyze this repo"); got != "analysis" {
		t.Fatalf("expected inferred analysis task, got %q", got)
	}
	short := promptPreview("  hello world  ")
	if short != "hello world" {
		t.Fatalf("unexpected short preview: %q", short)
	}
	long := promptPreview(strings.Repeat("a", 150))
	if len(long) != 143 || !strings.HasSuffix(long, "...") {
		t.Fatalf("unexpected long preview: %q", long)
	}
	hint := buildExecutionHint(ProviderRoute{Provider: "openai", Model: "gpt-4o"}, "analysis")
	if !strings.Contains(hint, "openai/gpt-4o") || !strings.Contains(hint, "analysis") {
		t.Fatalf("unexpected execution hint: %q", hint)
	}
}

func TestInferTaskTypeCoverage(t *testing.T) {
	cases := map[string]string{
		"please refactor this":          "coding",
		"search and analyze logs":       "analysis",
		"run on local model via ollama": "local",
		"say hello":                     "general",
	}
	for prompt, want := range cases {
		if got := inferTaskType(prompt); got != want {
			t.Fatalf("inferTaskType(%q) = %q, want %q", prompt, got, want)
		}
	}
}
