package subagents

// GetPermissionsForType returns the tool access manifest for a given SubagentType.
// By default, agents are granted access to basic read and search tools.
// Risky tools (write, edit, bash) are restricted to specific subagents.
func GetPermissionsForType(t SubagentType) []string {
	// Base safe tools available to almost all agents
	baseReadTools := []string{
		"read", "read_file", "file_read", "ls", "list_directory", "tree",
		"grep", "find", "search", "search_code", "search_files", "codesearch",
		"websearch", "webfetch", "web_search", "web_fetch", "memory_search",
	}

	// Edit and execution tools
	writeTools := []string{
		"write", "write_file", "file_write", "edit", "edit_file", "file_edit",
		"replace", "replace_lines", "apply_search_replace", "apply_diff", "patch",
		"multiedit", "cascade_edit",
	}
	
	execTools := []string{
		"bash", "execute_command", "shell", "run_shell_command", "PowerShell",
	}

	switch t {
	case TypeCode:
		// Full access for coding
		return append(append(baseReadTools, writeTools...), execTools...)
	
	case TypeResearch, TypeReview, TypePlan, TypeDoc:
		// Read-only agents
		return baseReadTools
	
	case TypeBuild, TypeTest, TypeDebug, TypeSecurity, TypeDevOps:
		// Specialized agents needing execution but minimal editing
		return append(baseReadTools, execTools...)
	
	default:
		return baseReadTools
	}
}

// IsToolAllowed checks if a specific tool is within the permission manifest.
func IsToolAllowed(t SubagentType, toolName string) bool {
	allowed := GetPermissionsForType(t)
	for _, a := range allowed {
		if a == toolName {
			return true
		}
	}
	// Fallback allow foundation meta tools like 'plan_enter', 'ask_user'
	if toolName == "plan_enter" || toolName == "plan_exit" || toolName == "AskUser" {
		return true
	}
	return false
}
