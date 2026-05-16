// helpers.go provides utility functions for pi tools.
package pi

import (
	"os"
	"path/filepath"
	"strings"
)

const defaultGrepLimit = 1000

// relPathOrBase returns the relative path of filePath from basePath,
// or just the base filename if they share no common path.
func relPathOrBase(filePath, basePath string) string {
	rel, err := filepath.Rel(basePath, filePath)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return filepath.Base(filePath)
	}
	return rel
}

// readFileLines reads a file and returns its lines.
func readFileLines(filePath string) ([]string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n"), nil
}

// truncateLineStr truncates a line to a reasonable display width.
func truncateLineStr(line string) (string, error) {
	max := 200
	if len(line) <= max {
		return line, nil
	}
	return line[:max] + "...", nil
}