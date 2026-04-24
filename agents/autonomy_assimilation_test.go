package agents

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestCouncilHelpers(t *testing.T) {
	if !councilApprovedFromReason("  VOTE: APPROVED because safe") {
		t.Fatal("expected approved vote detection")
	}
	if councilApprovedFromReason("VOTE: DENIED because risky") {
		t.Fatal("expected denied vote detection")
	}
	if !councilFinalVerdict([]CouncilVote{{Approved: true}, {Approved: true}, {Approved: false}}) {
		t.Fatal("expected 2/3 approval to pass")
	}
	if councilFinalVerdict([]CouncilVote{{Approved: true}, {Approved: false}, {Approved: false}}) {
		t.Fatal("expected 1/3 approval to fail")
	}
	personas := defaultCouncilPersonas()
	if len(personas) != 3 {
		t.Fatalf("expected 3 council personas, got %d", len(personas))
	}
}

func TestAutoDriveStartValidatesInputs(t *testing.T) {
	if err := autoDriveStart(context.Background(), nil, "goal", "sandbox", nil); err == nil || !strings.Contains(err.Error(), "autodrive is required") {
		t.Fatalf("expected nil autodrive error, got %v", err)
	}
	a := &AutoDrive{}
	if err := autoDriveStart(context.Background(), a, "goal", "sandbox", nil); err == nil || !strings.Contains(err.Error(), "director is required") {
		t.Fatalf("expected director error, got %v", err)
	}
	a.Director = &Director{}
	if err := autoDriveStart(context.Background(), a, "goal", "sandbox", nil); err == nil || !strings.Contains(err.Error(), "director provider is required") {
		t.Fatalf("expected provider error, got %v", err)
	}
	a.Director.Provider = &recordingProvider{}
	a.MaxIterations = 1
	if err := autoDriveStart(context.Background(), a, " ", "sandbox", nil); err == nil || !strings.Contains(err.Error(), "objective is required") {
		t.Fatalf("expected objective error, got %v", err)
	}
	if err := autoDriveStart(context.Background(), a, "goal", " ", nil); err == nil || !strings.Contains(err.Error(), "sandboxDir is required") {
		t.Fatalf("expected sandbox error, got %v", err)
	}
	a.MaxIterations = 0
	if err := autoDriveStart(context.Background(), a, "goal", "sandbox", nil); err == nil || !strings.Contains(err.Error(), "maxIterations must be positive") {
		t.Fatalf("expected iteration error, got %v", err)
	}
}

func TestAutoDriveStartAppendsPromptAndToolResults(t *testing.T) {
	provider := &recordingProvider{
		chatFn: func(ctx context.Context, messages []Message, tools []Tool) (Message, error) {
			return Message{Role: RoleAssistant, Content: "working", ToolCalls: []ToolCall{{ID: "1", Name: "bash", Args: "echo hi"}}}, nil
		},
	}
	director := &Director{Provider: provider, History: []Message{{Role: RoleSystem, Content: "sys"}}}
	a := &AutoDrive{Director: director, MaxIterations: 1}

	err := autoDriveStart(context.Background(), a, "ship feature", "/tmp/sandbox", func() {})
	if err != nil {
		t.Fatalf("autodrive failed: %v", err)
	}
	if a.IsRunning {
		t.Fatal("expected autodrive to stop after completion")
	}
	if len(director.History) < 4 {
		t.Fatalf("expected prompt, assistant, and tool history entries, got %#v", director.History)
	}
	if !strings.Contains(director.History[1].Content, "CRITICAL: All commands MUST be executed exclusively within '/tmp/sandbox'.") {
		t.Fatalf("unexpected prompt entry: %#v", director.History[1])
	}
	if director.History[len(director.History)-1].Role != RoleTool || director.History[len(director.History)-1].Name != "bash" {
		t.Fatalf("expected tool result message appended, got %#v", director.History[len(director.History)-1])
	}
	if !strings.Contains(director.History[len(director.History)-1].Content, "bash(echo hi)") {
		t.Fatalf("unexpected tool result content: %#v", director.History[len(director.History)-1])
	}
}

func TestAutoDriveStartPropagatesAbortAndChatErrors(t *testing.T) {
	abortProvider := &recordingProvider{chatFn: func(ctx context.Context, messages []Message, tools []Tool) (Message, error) {
		return Message{Role: RoleAssistant, Content: "working", ToolCalls: []ToolCall{{ID: "1", Name: "bash", Args: "echo hi"}}}, nil
	}}
	abortDrive := &AutoDrive{Director: &Director{Provider: abortProvider}, MaxIterations: 2}
	waitCalls := 0
	err := autoDriveStart(context.Background(), abortDrive, "goal", "sandbox", func() {
		waitCalls++
		if waitCalls == 2 {
			abortDrive.Abort()
		}
	})
	if err == nil || !strings.Contains(err.Error(), "autodrive aborted early") {
		t.Fatalf("expected abort error, got %v", err)
	}

	want := errors.New("chat failed")
	errProvider := &recordingProvider{chatFn: func(ctx context.Context, messages []Message, tools []Tool) (Message, error) { return Message{}, want }}
	errDrive := &AutoDrive{Director: &Director{Provider: errProvider}, MaxIterations: 1}
	if err := autoDriveStart(context.Background(), errDrive, "goal", "sandbox", func() {}); !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

func TestShellTranslatorValidationAndMessages(t *testing.T) {
	var nilTranslator *ShellTranslator
	if _, err := nilTranslator.Translate(context.Background(), "ls"); err == nil || !strings.Contains(err.Error(), "provider is required") {
		t.Fatalf("expected provider error for nil translator, got %v", err)
	}
	translator := &ShellTranslator{}
	if _, err := translator.Translate(context.Background(), "ls"); err == nil || !strings.Contains(err.Error(), "provider is required") {
		t.Fatalf("expected provider error, got %v", err)
	}
	translator.Provider = &recordingProvider{}
	if _, err := translator.Translate(context.Background(), "   "); err == nil || !strings.Contains(err.Error(), "intent is required") {
		t.Fatalf("expected intent error, got %v", err)
	}

	msgs := buildShellTranslationMessages("show git status")
	if len(msgs) != 2 || msgs[0].Role != RoleSystem || msgs[1].Role != RoleUser {
		t.Fatalf("unexpected translation messages: %#v", msgs)
	}
	if msgs[1].Content != "show git status" {
		t.Fatalf("unexpected user intent: %#v", msgs[1])
	}
}
