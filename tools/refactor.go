package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (r *Registry) registerRefactoringTools() {
	r.Tools = append(r.Tools, Tool{
		Name:        "bulk_refactor",
		Description: "Performs an automated mass find-and-replace across the codebase. Arguments: target_dir (string), old_pattern (string), new_pattern (string)",
		Execute: func(args map[string]interface{}) (string, error) {
			targetDir, ok := args["target_dir"].(string)
			if !ok {
				return "", fmt.Errorf("target_dir must be a string")
			}
			oldPattern, ok := args["old_pattern"].(string)
			if !ok {
				return "", fmt.Errorf("old_pattern must be a string")
			}
			newPattern, ok := args["new_pattern"].(string)
			if !ok {
				return "", fmt.Errorf("new_pattern must be a string")
			}

			modifiedFiles := 0

			err := filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				if info.IsDir() && (info.Name() == ".git" || info.Name() == "node_modules" || info.Name() == "vendor") {
					return filepath.SkipDir
				}
				if !info.IsDir() {
					content, err := os.ReadFile(path)
					if err != nil {
						return nil
					}
					
					strContent := string(content)
					if strings.Contains(strContent, oldPattern) {
						newContent := strings.ReplaceAll(strContent, oldPattern, newPattern)
						os.WriteFile(path, []byte(newContent), info.Mode())
						modifiedFiles++
					}
				}
				return nil
			})

			if err != nil {
				return "", err
			}

			return fmt.Sprintf("Bulk refactoring complete. Modified %d files.", modifiedFiles), nil
		},
	})
}
