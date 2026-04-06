package agents

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	codeModeTypeScriptFile = "hypercode_codemode.ts"
	codeModeShellFile      = "hypercode_codemode.sh"
)

// CodeMode executes an LLM-generated string of code directly using an external compiler/runtime natively
// This enforces the "Escape Hatch - Code Mode - 94% context reduction" specification from AGENTS.md
func ExecuteCodeMode(scriptBundle string, language string) (string, error) {
	resolved, err := resolveCodeMode(language)
	if err != nil {
		return "", err
	}

	tmpFile := filepath.Join(os.TempDir(), resolved.fileName)
	if err := os.WriteFile(tmpFile, []byte(scriptBundle), 0o644); err != nil {
		return "", err
	}
	defer os.Remove(tmpFile)

	cmd := exec.Command(resolved.command, append(resolved.args, tmpFile)...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("%s: %w\n%s", resolved.errorPrefix, err, output)
	}
	return string(output), nil
}

type codeModeExecution struct {
	fileName    string
	command     string
	args        []string
	errorPrefix string
}

func resolveCodeMode(language string) (codeModeExecution, error) {
	switch strings.ToLower(strings.TrimSpace(language)) {
	case "typescript", "ts":
		return codeModeExecution{fileName: codeModeTypeScriptFile, command: "bun", args: []string{"run"}, errorPrefix: "CodeMode exception"}, nil
	case "bash", "sh", "pwsh":
		return codeModeExecution{fileName: codeModeShellFile, command: "pwsh", args: []string{"-File"}, errorPrefix: "CodeMode native exception"}, nil
	default:
		return codeModeExecution{}, fmt.Errorf("missing language execution binding for %s", language)
	}
}
