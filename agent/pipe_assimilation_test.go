package agent

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func TestBuildPipePromptAndFormatResponse(t *testing.T) {
	prompt := buildPipePrompt("summarize", "hint text", "input body")
	for _, needle := range []string{"summarize", "Execution Hint:\nhint text", "Input Data:\ninput body"} {
		if !strings.Contains(prompt, needle) {
			t.Fatalf("expected prompt to contain %q, got %q", needle, prompt)
		}
	}
	formatted := formatPipeResponse("hint text", "done")
	if !strings.Contains(formatted, "[Pipe Execution]") || !strings.Contains(formatted, "hint text") || !strings.Contains(formatted, "done") {
		t.Fatalf("unexpected formatted response: %q", formatted)
	}
	if plain := formatPipeResponse("   ", "done"); plain != "done" {
		t.Fatalf("expected plain response when no hint, got %q", plain)
	}
}

func TestProcessPipeWithReaderValidation(t *testing.T) {
	if _, err := processPipeWithReader("prompt", os.ModeCharDevice, strings.NewReader("x"), func(string) (string, error) { return "", nil }); err == nil || !strings.Contains(err.Error(), "no data piped to stdin") {
		t.Fatalf("expected no-pipe error, got %v", err)
	}
	if _, err := processPipeWithReader("prompt", 0, nil, func(string) (string, error) { return "", nil }); err == nil || !strings.Contains(err.Error(), "stdin reader is required") {
		t.Fatalf("expected nil reader error, got %v", err)
	}
	if _, err := processPipeWithReader("prompt", 0, strings.NewReader("x"), nil); err == nil || !strings.Contains(err.Error(), "chat function is required") {
		t.Fatalf("expected nil chat error, got %v", err)
	}
}

func TestProcessPipeWithReaderSuccessAndErrorPaths(t *testing.T) {
	var seenPrompt string
	result, err := processPipeWithReader("summarize this", 0, strings.NewReader("alpha\nbeta"), func(prompt string) (string, error) {
		seenPrompt = prompt
		return "analysis complete", nil
	})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if !strings.Contains(seenPrompt, "summarize this") || !strings.Contains(seenPrompt, "Input Data:\nalpha\nbeta") {
		t.Fatalf("unexpected combined prompt: %q", seenPrompt)
	}
	if !strings.Contains(result, "[Pipe Execution]") || !strings.Contains(result, "analysis complete") {
		t.Fatalf("unexpected result: %q", result)
	}

	want := errors.New("chat failed")
	if _, err := processPipeWithReader("prompt", 0, strings.NewReader("x"), func(string) (string, error) { return "", want }); !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

func TestAgentProcessPipeRejectsNilReceiver(t *testing.T) {
	var a *Agent
	if _, err := a.ProcessPipe("hello"); err == nil || !strings.Contains(err.Error(), "agent is required") {
		t.Fatalf("expected nil agent error, got %v", err)
	}
}
