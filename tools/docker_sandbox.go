package tools

import (
	"fmt"
	"os/exec"
	"strings"
)

// DockerSandbox intercepts bash tools and wraps them in an ephemeral container.
type DockerSandbox struct {
	Enabled bool
}

func NewDockerSandbox() *DockerSandbox {
	return &DockerSandbox{Enabled: false} // Optional plugin
}

// WrapCommand modifies a raw command string to execute securely inside an Ubuntu container,
// mapping the provided `cwd` into `/workspace`.
func (d *DockerSandbox) WrapCommand(command, cwd string) string {
	if !d.Enabled {
		return command
	}

	fmt.Printf("[DockerSandboxPlugin] Intercepting command for sandboxing: %s\n", command)
	
	// Escape quotes properly for bash -c
	escapedCmd := strings.ReplaceAll(command, "\"", "\\\"")
	
	return fmt.Sprintf("docker run --rm -v \"%s:/workspace\" -w /workspace ubuntu:latest bash -c \"%s\"", cwd, escapedCmd)
}

// CheckDocker checks if Docker daemon is responsive.
func (d *DockerSandbox) CheckDocker() error {
	cmd := exec.Command("docker", "info")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker not available: %w", err)
	}
	return nil
}
