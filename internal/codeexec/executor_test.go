package codeexec

import (
	"context"
	"testing"
	"time"
)

func TestNewCodeExecutor(t *testing.T) {
	ce := NewCodeExecutor()
	if ce == nil {
		t.Fatal("should create executor")
	}
}

func TestExecuteShell(t *testing.T) {
	ce := NewCodeExecutor()
	result, err := ce.Execute(context.Background(), ExecutionConfig{
		Language: Shell,
		Code:     "echo hello world",
		Timeout:  10 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.ExitCode != 0 {
		t.Errorf("exit code: %d, stderr: %s", result.ExitCode, result.Stderr)
	}
	if result.Language != "shell" {
		t.Errorf("language: %s", result.Language)
	}
	if result.TimedOut {
		t.Error("should not timeout")
	}
}

func TestExecuteShellWithOutput(t *testing.T) {
	ce := NewCodeExecutor()
	result, err := ce.Execute(context.Background(), ExecutionConfig{
		Language: Shell,
		Code:     "echo test-output",
		Timeout:  5 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.ExitCode != 0 {
		t.Errorf("exit: %d", result.ExitCode)
	}
	t.Logf("stdout: %q stderr: %q", result.Stdout, result.Stderr)
}

func TestExecuteShellError(t *testing.T) {
	ce := NewCodeExecutor()
	result, err := ce.Execute(context.Background(), ExecutionConfig{
		Language: Shell,
		Code:     "exit 42",
		Timeout:  5 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.ExitCode != 42 {
		t.Errorf("exit code: %d", result.ExitCode)
	}
}

func TestExecuteTimeout(t *testing.T) {
	ce := NewCodeExecutor()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	result, err := ce.Execute(ctx, ExecutionConfig{
		Language: Shell,
		Code:     "sleep 10",
		Timeout:  50 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.TimedOut {
		t.Error("should have timed out")
	}
}

func TestExecuteWithEnv(t *testing.T) {
	ce := NewCodeExecutor()
	result, err := ce.Execute(context.Background(), ExecutionConfig{
		Language: Shell,
		Code:     "echo $TEST_VAR",
		Timeout:  5 * time.Second,
		Env:      map[string]string{"TEST_VAR": "hello-env"},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("output: %q", result.Stdout)
}

func TestExecuteGo(t *testing.T) {
	ce := NewCodeExecutor()
	result, err := ce.Execute(context.Background(), ExecutionConfig{
		Language: Go,
		Code: `package main
import "fmt"
func main() { fmt.Println("go-exec") }`,
		Timeout: 30 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.ExitCode != 0 {
		t.Errorf("exit: %d, stderr: %s", result.ExitCode, result.Stderr)
	}
	t.Logf("stdout: %q", result.Stdout)
}

func TestIsLanguageAvailable(t *testing.T) {
	if !IsLanguageAvailable(Shell) {
		t.Error("shell should always be available")
	}
	if !IsLanguageAvailable(Go) {
		t.Error("Go should be available (we're running Go tests)")
	}
}

func TestListAvailableLanguages(t *testing.T) {
	langs := ListAvailableLanguages()
	if len(langs) == 0 {
		t.Error("should have at least shell available")
	}
	t.Logf("available: %v", langs)
}

func TestUnsupportedLanguage(t *testing.T) {
	ce := NewCodeExecutor()
	_, err := ce.Execute(context.Background(), ExecutionConfig{
		Language: "cobol",
		Code:     "DISPLAY 'HELLO'",
	})
	if err == nil {
		t.Error("should error for unsupported language")
	}
}
