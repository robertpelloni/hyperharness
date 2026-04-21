package lockfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteAndRead(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.lock")

	rec := Record{
		Host:      "localhost",
		Port:      8080,
		Version:   "1.0.0",
		StartedAt: "2025-01-01T00:00:00Z",
		Startup: &StartupProvenance{
			ActiveRuntime: "go",
			ActivePort:    8080,
			PortDecision:  "default",
		},
	}

	if err := Write(path, rec); err != nil {
		t.Fatal(err)
	}

	got, err := Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Host != "localhost" {
		t.Errorf("host: %s", got.Host)
	}
	if got.Port != 8080 {
		t.Errorf("port: %d", got.Port)
	}
	if got.Startup == nil {
		t.Fatal("startup should not be nil")
	}
	if got.Startup.ActiveRuntime != "go" {
		t.Errorf("runtime: %s", got.Startup.ActiveRuntime)
	}
}

func TestReadNonexistent(t *testing.T) {
	_, err := Read("/nonexistent/lock")
	if err == nil {
		t.Error("should fail for nonexistent")
	}
}

func TestWriteCreatesDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "deep", "test.lock")
	err := Write(path, Record{Host: "test"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal("file should exist")
	}
}
