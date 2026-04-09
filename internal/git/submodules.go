// Package git provides submodule management operations.
// Ported from hypercode/go/internal/git/submodules.go
//
// WHAT: Git submodule listing, updating, and status tracking
// WHY: HyperHarness uses 30+ submodules from superai that need management
// HOW: Concurrent git submodule update with structured error reporting
package git

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// SubmoduleResult tracks the result of a submodule operation.
type SubmoduleResult struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Success bool   `json:"success"`
	Output  string `json:"output,omitempty"`
	Error   string `json:"error,omitempty"`
}

// UpdateReport summarizes a batch submodule update.
type UpdateReport struct {
	Total      int               `json:"total"`
	Successful int               `json:"successful"`
	Failed     int               `json:"failed"`
	Details    []SubmoduleResult `json:"details"`
}

// ListSubmodules parses `git submodule status` to get submodule paths.
func ListSubmodules(ctx context.Context, repoRoot string) ([]string, error) {
	cmd := exec.CommandContext(ctx, "git", "submodule", "status")
	cmd.Dir = repoRoot
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list submodules: %w", err)
	}
	return parseSubmoduleStatus(string(output)), nil
}

func parseSubmoduleStatus(output string) []string {
	var paths []string
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		path := strings.TrimSpace(parts[1])
		if path != "" {
			paths = append(paths, path)
		}
	}
	sort.Strings(paths)
	return paths
}

// UpdateAll concurrently updates all submodules to their remote tracked branches.
func UpdateAll(ctx context.Context, repoRoot string) (*UpdateReport, error) {
	paths, err := ListSubmodules(ctx, repoRoot)
	if err != nil {
		return nil, err
	}

	report := &UpdateReport{
		Total:   len(paths),
		Details: make([]SubmoduleResult, 0, len(paths)),
	}
	if len(paths) == 0 {
		return report, nil
	}

	// Init recursively first
	initCmd := exec.CommandContext(ctx, "git", "submodule", "update", "--init", "--recursive")
	initCmd.Dir = repoRoot
	_ = initCmd.Run() // Non-fatal

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, path := range paths {
		wg.Add(1)
		go func(subPath string) {
			defer wg.Done()
			name := filepath.Base(strings.ReplaceAll(subPath, "\\", "/"))

			cmd := exec.CommandContext(ctx, "git", "submodule", "update", "--remote", subPath)
			cmd.Dir = repoRoot
			output, err := cmd.CombinedOutput()

			mu.Lock()
			defer mu.Unlock()

			res := SubmoduleResult{
				Name:    name,
				Path:    subPath,
				Success: err == nil,
				Output:  string(output),
			}
			if err != nil {
				res.Error = err.Error()
				report.Failed++
			} else {
				report.Successful++
			}
			report.Details = append(report.Details, res)
		}(path)
	}

	wg.Wait()
	sort.Slice(report.Details, func(i, j int) bool {
		return report.Details[i].Path < report.Details[j].Path
	})
	return report, nil
}
