package lockfile

import (
	"path/filepath"
	"testing"
)

func TestWriteAndRead(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.lock")
	rec := Record{Host: "localhost", Port: 8080, Version: "1.0.0", StartedAt: "2025-01-01"}
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
}

func TestReadNonexistent(t *testing.T) {
	_, err := Read("/nonexistent/lock")
	if err == nil {
		t.Error("should fail")
	}
}
