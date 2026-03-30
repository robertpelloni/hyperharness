package tools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Registry holds all available tools for the agent
type Registry struct {
	Tools []Tool
}

type Tool struct {
	Name        string
	Description string
	Execute     func(args map[string]interface{}) (string, error)
}

func NewRegistry() *Registry {
	r := &Registry{}
	r.registerCoreTools()
	r.registerFileTools()
	r.registerInterpreterTools()
	r.registerSearchTools()
	r.registerAdvancedTools()
	r.registerRefactoringTools()
	r.registerCloudTools()
	r.registerIntegrationsTools()
	r.registerBookmarkTools()
	r.registerGUITools()
	r.registerCloudOrchestratorTools()
	r.registerSystemTools()
	r.registerLlamafileTools()
	return r
}
func (r *Registry) registerCoreTools() {
	r.Tools = append(r.Tools, Tool{
		Name:        "run_shell_command",
		Description: "Executes a shell command.",
		Execute: func(args map[string]interface{}) (string, error) {
			cmdStr := args["command"].(string)
			cmd := exec.Command("cmd", "/c", cmdStr)
			out, err := cmd.CombinedOutput()
			return string(out), err
		},
	})
}

func (r *Registry) registerSearchTools() {
	r.Tools = append(r.Tools, Tool{
		Name:        "search",
		Description: "Recursively searches for a pattern in the codebase. Arguments: pattern (string)",
		Execute: func(args map[string]interface{}) (string, error) {
			pattern, _ := args["pattern"].(string)
			// Simple fallback if ripgrep isn't there, matching Aider/Adrenaline
			var results string
			filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				if info.IsDir() && (info.Name() == ".git" || info.Name() == "node_modules") {
					return filepath.SkipDir
				}
				if !info.IsDir() {
					content, _ := os.ReadFile(path)
					if strings.Contains(string(content), pattern) {
						rel, _ := filepath.Rel(".", path)
						results += fmt.Sprintf("Match in %s\n", rel)
					}
				}
				return nil
			})
			return results + "Search functionality complete", nil
		},
	})
}
