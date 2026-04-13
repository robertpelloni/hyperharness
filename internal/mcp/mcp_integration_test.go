package mcp

import (
	"strings"
	"testing"
)

func TestMCPServerLifecycle(t *testing.T) {
	registry := NewRegistry()

	// Create a mock stdio MCP server pointing to a shell command that echoes an init response
	mockServer := &MCPServer{
		Name:      "test-server",
		Transport: "stdio",
		// We use `echo` to return an error JSON indicating it's an invalid request
		// to test parsing and the execution pipeline. In a real system, we'd mock a Go executable.
		Command: "echo",
		Args:    []string{`{"jsonrpc":"2.0","id":1,"error":{"code":-32600,"message":"Invalid Request"}}`},
	}

	err := registry.RegisterServer(mockServer)
	if err != nil {
		t.Fatalf("RegisterServer failed: %v", err)
	}

	// Connect expects specific initialization responses, so the simple echo command will fail
	// during the handshake. The test verifies it tries to start the process and communicate.
	err = registry.Connect("test-server")

	// Since echo exits immediately, reading the response or subsequent writes will fail,
	// confirming the pipeline ran stdio transport logic.
	if err == nil {
		t.Fatalf("Expected connection to fail due to mock server exit, but it succeeded")
	}

	if !strings.Contains(err.Error(), "decode init response") && !strings.Contains(err.Error(), "read init response") && !strings.Contains(err.Error(), "Invalid Request") && !strings.Contains(err.Error(), "EOF") {
		t.Errorf("Expected an IO or parsing error from mock server, got: %v", err)
	}
}
