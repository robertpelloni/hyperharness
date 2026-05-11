package tools

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
)

// registerPowerShellTools registers PowerShell execution tools explicitly for Windows parity.
func registerPowerShellTools(r *Registry) {
	r.Tools = append(r.Tools, Tool{
		Name:        "PowerShell",
		Description: "Execute a PowerShell command. This tool provides deep integration for Windows systems running PowerShell scripts or inline commands.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"command": {
					"type": "string",
					"description": "The PowerShell command or script block to execute."
				}
			},
			"required": ["command"]
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			command := GetStr(args, "command")
			if command == "" {
				return "", fmt.Errorf("command is required")
			}

			if runtime.GOOS != "windows" {
				return "", fmt.Errorf("PowerShell tool is only supported on Windows systems, but running on %s", runtime.GOOS)
			}

			cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", command)
			
			output, err := cmd.CombinedOutput()
			if err != nil {
				return string(output), fmt.Errorf("powershell execution failed: %v\nOutput: %s", err, string(output))
			}
			
			return string(output), nil
		},
	})
}
