package adapters

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHyperCodeRepoCandidatesDeduplicateAndOrder(t *testing.T) {
	candidates := hyperCodeRepoCandidates(filepath.Join("/tmp", "workspace", "hyperharness"), "/home/demo")
	if len(candidates) != 2 {
		t.Fatalf("expected 2 unique candidates, got %#v", candidates)
	}
	if candidates[0] != filepath.Clean(filepath.Join("/tmp", "workspace", "hyperharness", "..", "hypercode")) {
		t.Fatalf("unexpected first candidate: %#v", candidates)
	}
	if candidates[1] != filepath.Clean(filepath.Join("/home/demo", "workspace", "hypercode")) {
		t.Fatalf("unexpected second candidate: %#v", candidates)
	}
}

func TestRenderHyperCodeSystemContextIncludesImportantSections(t *testing.T) {
	status := HyperCodeStatus{
		Assimilated:       true,
		BorgCoreURL:       "http://127.0.0.1:7331",
		MemoryContext:     "MEMCTX",
		Provider:          ProviderStatus{Available: []string{"openai"}},
		MCPServerNames:    []string{"alpha", "beta"},
		HyperCodeRepoPath: "/repo/hypercode",
		Warnings:          []string{"warn-a", "warn-b"},
	}
	text := renderHyperCodeSystemContext(status)
	for _, needle := range []string{"[HyperCode Adapter]", "Assimilated: true", "Borg Core URL: http://127.0.0.1:7331", "MEMCTX", "Configured MCP servers: alpha, beta", "HyperCode repo: /repo/hypercode", "Warnings: warn-a; warn-b"} {
		if !strings.Contains(text, needle) {
			t.Fatalf("expected context to contain %q, got %q", needle, text)
		}
	}
}

func TestHyperCodeAdapterNilAndRepoDiscoveryBehavior(t *testing.T) {
	var nilAdapter *HyperCodeAdapter
	status := nilAdapter.Status()
	if len(status.Warnings) == 0 || !strings.Contains(strings.Join(status.Warnings, "; "), "hypercode adapter is nil") {
		t.Fatalf("expected nil adapter warning, got %#v", status)
	}
	if ctx := nilAdapter.BuildSystemContext(); !strings.Contains(ctx, "[HyperCode Adapter]") {
		t.Fatalf("expected system context for nil adapter, got %q", ctx)
	}

	dir := t.TempDir()
	hypercodeDir := filepath.Join(dir, "..", "hypercode")
	if err := os.MkdirAll(hypercodeDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(hypercodeDir, "README.md"), []byte("# HyperCode"), 0o644); err != nil {
		t.Fatal(err)
	}
	adapter := NewHyperCodeAdapter(dir)
	adapter.homeDir = t.TempDir()
	if repo, ok := adapter.findHyperCodeRepo(); !ok || repo == "" {
		t.Fatalf("expected repo discovery, got %q ok=%v", repo, ok)
	}
}
