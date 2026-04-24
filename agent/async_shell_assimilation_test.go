package agent

import (
	"strings"
	"testing"
)

func TestNewAsyncWorkerInitializesChannels(t *testing.T) {
	worker := NewAsyncWorker()
	defer worker.Close()

	if worker == nil {
		t.Fatal("expected async worker")
	}
	if worker.tasks == nil {
		t.Fatal("expected task channel to be initialized")
	}
	if worker.closed == nil {
		t.Fatal("expected closed channel to be initialized")
	}
	if cap(worker.tasks) != asyncWorkerBufferSize {
		t.Fatalf("expected task buffer size %d, got %d", asyncWorkerBufferSize, cap(worker.tasks))
	}
}

func TestAsyncWorkerEnqueueNilReceiverIsNoOp(t *testing.T) {
	var worker *AsyncWorker
	worker.Enqueue("task")
}

func TestAsyncWorkerCloseIsIdempotent(t *testing.T) {
	worker := NewAsyncWorker()
	worker.Close()
	worker.Close()
	worker.Enqueue("after-close")
}

func TestAsyncTaskMessages(t *testing.T) {
	enqueued := asyncTaskEnqueuedMessage("lint repo")
	if enqueued != "[Async] Task enqueued: lint repo" {
		t.Fatalf("unexpected enqueued message: %q", enqueued)
	}
	completed := asyncTaskCompletedMessage("lint repo")
	if !strings.Contains(completed, "Background task completed: lint repo") {
		t.Fatalf("unexpected completed message: %q", completed)
	}
}

func TestBuildShellPromptIncludesRequestAndInstruction(t *testing.T) {
	prompt := buildShellPrompt("list files recursively")
	if !strings.Contains(prompt, "single, valid shell command for a Windows environment") {
		t.Fatalf("expected environment instruction, got %q", prompt)
	}
	if !strings.Contains(prompt, "Output ONLY the command, nothing else.") {
		t.Fatalf("expected strict output instruction, got %q", prompt)
	}
	if !strings.Contains(prompt, "Request: list files recursively") {
		t.Fatalf("expected prompt to include request, got %q", prompt)
	}
}

func TestSuggestShellCommandRejectsMissingClient(t *testing.T) {
	a := &Agent{}
	_, err := a.SuggestShellCommand("show git status")
	if err == nil || !strings.Contains(err.Error(), "openai client is required") {
		t.Fatalf("expected missing client error, got %v", err)
	}
}
