package agent

import (
	"errors"
	"strings"
	"testing"
)

func TestOracleQueryWithChatRunsThreeStageFlow(t *testing.T) {
	var prompts []string
	responses := []string{"1. inspect\n2. synthesize", "facts and evidence", "final answer"}
	index := 0

	chat := func(prompt string) (string, error) {
		prompts = append(prompts, prompt)
		if index >= len(responses) {
			return "", errors.New("unexpected extra call")
		}
		resp := responses[index]
		index++
		return resp, nil
	}

	answer, err := oracleQueryWithChat("why is the sky blue?", chat)
	if err != nil {
		t.Fatalf("oracle query failed: %v", err)
	}
	if answer != "final answer" {
		t.Fatalf("expected final synthesized answer, got %q", answer)
	}
	if len(prompts) != 3 {
		t.Fatalf("expected 3 chat prompts, got %d (%#v)", len(prompts), prompts)
	}
	if !strings.Contains(prompts[0], "Create a multi-step research plan") || !strings.Contains(prompts[0], "why is the sky blue?") {
		t.Fatalf("unexpected plan prompt: %q", prompts[0])
	}
	if !strings.Contains(prompts[1], "Plan:\n1. inspect\n2. synthesize") || !strings.Contains(prompts[1], "Original Query: why is the sky blue?") {
		t.Fatalf("unexpected execution prompt: %q", prompts[1])
	}
	if !strings.Contains(prompts[2], "facts and evidence") || !strings.Contains(prompts[2], "Provide a comprehensive, authoritative, and definitive answer.") {
		t.Fatalf("unexpected synthesis prompt: %q", prompts[2])
	}
}

func TestOracleQueryWithChatRejectsEmptyPrompt(t *testing.T) {
	_, err := oracleQueryWithChat("   ", func(string) (string, error) { return "", nil })
	if err == nil || !strings.Contains(err.Error(), "prompt is required") {
		t.Fatalf("expected prompt validation error, got %v", err)
	}
}

func TestOracleQueryWithChatPropagatesStageErrors(t *testing.T) {
	want := errors.New("planning failed")
	_, err := oracleQueryWithChat("query", func(string) (string, error) { return "", want })
	if !errors.Is(err, want) {
		t.Fatalf("expected planning error %v, got %v", want, err)
	}
}

func TestBuildCompletionPromptIncludesPrefixAndSuffix(t *testing.T) {
	prompt := buildCompletionPrompt("func add(a int, b int) {", "}")
	if !strings.Contains(prompt, "Prefix:\nfunc add(a int, b int) {") {
		t.Fatalf("expected prompt to include prefix, got %q", prompt)
	}
	if !strings.Contains(prompt, "Suffix:\n}") {
		t.Fatalf("expected prompt to include suffix, got %q", prompt)
	}
	if !strings.Contains(prompt, "Output ONLY the completion string.") {
		t.Fatalf("expected strict output instruction, got %q", prompt)
	}
}

func TestSuggestCompletionRejectsMissingClient(t *testing.T) {
	a := &Agent{}
	_, err := a.SuggestCompletion("prefix", "suffix")
	if err == nil || !strings.Contains(err.Error(), "openai client is required") {
		t.Fatalf("expected missing client error, got %v", err)
	}
}
