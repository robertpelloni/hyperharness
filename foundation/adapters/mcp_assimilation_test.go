package adapters

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestSortedMCPServerNamesAndStatusHelpers(t *testing.T) {
	conf := &MCPConfig{MCPServers: map[string]MCPServerConfig{"zeta": {}, "alpha": {}, "beta": {}}}
	names := sortedMCPServerNames(conf)
	joined := strings.Join(names, ",")
	if joined != "alpha,beta,zeta" {
		t.Fatalf("unexpected sorted names: %q", joined)
	}

	status := buildMCPServerStatus(nil, "demo", MCPServerConfig{Command: "cmd", Args: []string{"/c", "echo"}, Env: map[string]string{"FOO": "BAR"}})
	if status.Name != "demo" || status.Command != "cmd" || !status.HasEnv {
		t.Fatalf("unexpected server status: %#v", status)
	}
	if len(status.ToolHints) != 2 || status.RouteHint != "demo" {
		t.Fatalf("unexpected tool hints/route hint: %#v", status)
	}
}

func TestMCPAdapterNilAndLookupHelpers(t *testing.T) {
	var nilAdapter *MCPAdapter
	status := nilAdapter.Status()
	if len(status.Warnings) == 0 || !strings.Contains(strings.Join(status.Warnings, "; "), "mcp adapter is nil") {
		t.Fatalf("expected nil adapter warning, got %#v", status)
	}
	if routeHintForAdapter(nil, "demo") != "demo" {
		t.Fatal("expected nil adapter route hint to fall back to name")
	}
	if mcpAdapterWorkingDir(nil) != "" {
		t.Fatal("expected nil adapter working dir to be empty")
	}
	if _, ok := nilAdapter.LookupServer("demo"); ok {
		t.Fatal("expected nil adapter lookup to fail")
	}
}

func TestStartMCPServerValidation(t *testing.T) {
	if _, err := startMCPServer(MCPServerConfig{}, "demo", ""); err == nil || !strings.Contains(err.Error(), "MCP server demo has no command") {
		t.Fatalf("expected missing command error, got %v", err)
	}
}

func TestDefaultToolHintsAndCommandResolvable(t *testing.T) {
	hints := defaultToolHintsForServer(" Demo ", MCPServerConfig{})
	if len(hints) != 2 || hints[0] != "mcp:demo:list-tools" || hints[1] != "mcp:demo:call-tool" {
		t.Fatalf("unexpected tool hints: %#v", hints)
	}
	blankHints := defaultToolHintsForServer("   ", MCPServerConfig{})
	if blankHints[0] != "mcp:mcp:list-tools" {
		t.Fatalf("expected fallback base, got %#v", blankHints)
	}
	if commandResolvable("") {
		t.Fatal("expected blank command to be non-resolvable")
	}
	path, err := exec.LookPath("go")
	if err == nil && path != "" && !commandResolvable("go") {
		t.Fatal("expected go command to be resolvable when on PATH")
	}
}

func TestStartMCPServerUsesWorkingDir(t *testing.T) {
	server := MCPServerConfig{Command: "go", Args: []string{"version"}}
	cmd, err := startMCPServer(server, "demo", filepath.Clean("."))
	if err != nil {
		// acceptable on systems without go on PATH; the validation path is already covered
		return
	}
	defer func() {
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
	}()
	if cmd.Dir == "" {
		t.Fatal("expected working dir to be set")
	}
}
