package ast

import (
	"strings"
	"testing"
)

func TestSummarizeGoFile(t *testing.T) {
	src := []byte(`package main

import (
	"fmt"
	"strings"
)

const DefaultLimit = 100

type Config struct {
	Timeout int
}

func (*Config) Get() int {
	return c.Timeout
}

func ExecuteTask(taskName string) error {
	fmt.Println("Executing", taskName)
	return nil
}
`)

	summary, err := SummarizeGoFile("test.go", src)
	if err != nil {
		t.Fatalf("SummarizeGoFile failed: %v", err)
	}

	if !strings.Contains(summary, "package main") {
		t.Error("Missing package main")
	}
	if !strings.Contains(summary, `"fmt"`) {
		t.Error("Missing fmt import")
	}
	if !strings.Contains(summary, "const DefaultLimit") {
		t.Error("Missing constant DefaultLimit")
	}
	if !strings.Contains(summary, "type Config struct{ ... }") {
		t.Error("Missing struct type Config")
	}
	if !strings.Contains(summary, "func (*Config) Get() int") {
		t.Errorf("Missing Get() method. Got:\n%s", summary)
	}
	if !strings.Contains(summary, "func ExecuteTask(taskName string) error") {
		t.Error("Missing ExecuteTask function")
	}
}
