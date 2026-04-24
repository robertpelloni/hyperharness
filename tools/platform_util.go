// Package tools provides platform utility functions for file operations.
// These are used by the download and web fetch tools.
package tools

import (
	"io"
	"os"
	"path/filepath"
)

// createParentDirs creates parent directories for a file path.
func createParentDirs(filePath string) error {
	dir := filepath.Dir(filePath)
	return os.MkdirAll(dir, 0o755)
}

// writeDownloadFile writes the content of a reader to a file.
func writeDownloadFile(filePath string, reader io.Reader) (int64, error) {
	f, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	return io.Copy(f, reader)
}
