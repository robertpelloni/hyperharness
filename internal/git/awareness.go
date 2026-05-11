package git

import (
	
	"fmt"
	"os/exec"
	"strings"
)

// GetAwarenessContext gathers git diff and recent commits to inject into the system prompt.
// If the directory is not a git repository, it returns empty string and no error.
func GetAwarenessContext(workingDir string) (string, error) {
	// Check if git repo
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = workingDir
	if err := cmd.Run(); err != nil {
		return "", nil // Not a git repo
	}

	var sb strings.Builder
	sb.WriteString("<git_context>\n")

	// Get current branch
	branchCmd := exec.Command("git", "branch", "--show-current")
	branchCmd.Dir = workingDir
	branchOut, _ := branchCmd.Output()
	branch := strings.TrimSpace(string(branchOut))
	if branch != "" {
		sb.WriteString(fmt.Sprintf("Current Branch: %s\n\n", branch))
	}

	// Get uncommitted changes (staged and unstaged)
	diffCmd := exec.Command("git", "diff", "HEAD")
	diffCmd.Dir = workingDir
	diffOut, _ := diffCmd.Output()
	diff := strings.TrimSpace(string(diffOut))
	
	if diff != "" {
		// Truncate diff if it's too huge to avoid blowing up context
		if len(diff) > 4000 {
			diff = diff[:4000] + "\n... (diff truncated due to length)"
		}
		sb.WriteString("Uncommitted Changes (Diff against HEAD):\n```diff\n")
		sb.WriteString(diff)
		sb.WriteString("\n```\n\n")
	} else {
		sb.WriteString("Working tree is clean. No uncommitted changes.\n\n")
	}

	// Get last 5 commits
	logCmd := exec.Command("git", "log", "-n", "5", "--oneline")
	logCmd.Dir = workingDir
	logOut, _ := logCmd.Output()
	logText := strings.TrimSpace(string(logOut))
	
	if logText != "" {
		sb.WriteString("Recent Commits:\n```\n")
		sb.WriteString(logText)
		sb.WriteString("\n```\n")
	}

	sb.WriteString("</git_context>\n")

	return sb.String(), nil
}
