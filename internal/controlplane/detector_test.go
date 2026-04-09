package controlplane

import (
	"context"
	"testing"
	"time"

	"github.com/robertpelloni/hyperharness/internal/buildinfo"
)

func TestDetectorCreation(t *testing.T) {
	d := NewDetector(5*time.Second, 30*time.Second)
	if len(d.definitions) == 0 {
		t.Error("should have tool definitions")
	}
}

func TestDetectAll(t *testing.T) {
	// Use a minimal detector for speed
	d := NewDetector(1*time.Second, 30*time.Second)
	d.definitions = d.definitions[:5] // Only check first 5
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	tools, err := d.DetectAll(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(tools) == 0 {
		t.Error("should detect tools")
	}

	// Go should always be available in our environment
	var goFound bool
	for _, tool := range tools {
		if tool.Type == "go" {
			goFound = tool.Available
			t.Logf("Go: available=%v version=%s path=%s", tool.Available, tool.Version, tool.Path)
		}
	}
	if !goFound {
		t.Error("Go should be available")
	}
}

func TestDetectAvailable(t *testing.T) {
	d := NewDetector(1*time.Second, 30*time.Second)
	d.definitions = d.definitions[:5]
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	available, err := d.DetectAvailable(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(available) == 0 {
		t.Error("should have available tools")
	}
	for _, tool := range available {
		if !tool.Available {
			t.Errorf("%s should be available", tool.Name)
		}
	}
}

func TestDetectOne(t *testing.T) {
	d := NewDetector(1*time.Second, 30*time.Second)
	d.definitions = d.definitions[:5]
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	goTool := d.DetectOne(ctx, "go")
	if goTool == nil {
		t.Fatal("should find go")
	}
	if !goTool.Available {
		t.Error("go should be available")
	}

	missing := d.DetectOne(ctx, "nonexistent-tool-xyz")
	if missing != nil && missing.Available {
		t.Error("nonexistent tool should not be available")
	}
}

func TestDetectCaching(t *testing.T) {
	d := NewDetector(1*time.Second, 30*time.Second)
	d.definitions = d.definitions[:3]
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tools1, _ := d.DetectAll(ctx)
	tools2, _ := d.DetectAll(ctx)

	if len(tools1) != len(tools2) {
		t.Error("cached results should match")
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"go version go1.26.2 windows/amd64", "1.26.2"},
		{"v1.0.0", "1.0.0"},
		{"2.3.4", "2.3.4"},
		{"no version here", "no version here"},
		{"", "unknown"},
	}
	for _, tt := range tests {
		result := parseVersion(tt.input)
		if result != tt.expected {
			t.Errorf("parseVersion(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestBuildinfoVersion(t *testing.T) {
	// Just verify buildinfo is accessible
	if buildinfo.Version == "" {
		t.Error("version should not be empty")
	}
	if buildinfo.ProductName == "" {
		t.Error("product name should not be empty")
	}
	t.Logf("%s v%s", buildinfo.ProductName, buildinfo.Version)
}
