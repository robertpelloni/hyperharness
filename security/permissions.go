package security

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// ═══════════════════════════════════════════════════════════════════════
// security.go — Comprehensive security sandbox and permission system
// Ported from goose SandboxService, claude-code permissions, and pi-cli
// security patterns.
//
// WHAT: Multi-layer security with autonomy levels, dangerous command
//       detection, path sandboxing, and tool-specific permission rules.
// WHY: Prevent AI agents from executing destructive or unauthorized
//       operations on the user's system.
// HOW: PermissionManager with configurable rules, Sandbox for path/file
//       restrictions, and CommandChecker for dangerous pattern detection.
// ═══════════════════════════════════════════════════════════════════════

// AutonomyLevel defines the freedom of action for Agent operations.
type AutonomyLevel int

const (
	AutonomyLevelNone    AutonomyLevel = iota // Explicit approval for all commands
	AutonomyLevelRead                         // Can read files without asking
	AutonomyLevelWrite                        // Can read/write files without asking
	AutonomyLevelExecute                      // Can run safe builds without asking
	AutonomyLevelGod                          // Total bypass (dangerous)
)

func (l AutonomyLevel) String() string {
	switch l {
	case AutonomyLevelNone:
		return "none"
	case AutonomyLevelRead:
		return "read"
	case AutonomyLevelWrite:
		return "write"
	case AutonomyLevelExecute:
		return "execute"
	case AutonomyLevelGod:
		return "god"
	default:
		return "unknown"
	}
}

// ActionType categorizes what kind of operation is being requested.
type ActionType string

const (
	ActionFileRead     ActionType = "file_read"
	ActionFileWrite    ActionType = "file_write"
	ActionFileDelete   ActionType = "file_delete"
	ActionFileMove     ActionType = "file_move"
	ActionCommandExec  ActionType = "command_exec"
	ActionWebFetch     ActionType = "web_fetch"
	ActionWebSearch    ActionType = "web_search"
	ActionMCPToolCall  ActionType = "mcp_tool_call"
	ActionCodeExec     ActionType = "code_exec"
	ActionGitOperation ActionType = "git_operation"
	ActionEnvRead      ActionType = "env_read"
	ActionNetwork      ActionType = "network"
)

// PermissionDecision represents the outcome of a permission check.
type PermissionDecision string

const (
	DecisionAllow      PermissionDecision = "allow"
	DecisionDeny       PermissionDecision = "deny"
	DecisionAsk        PermissionDecision = "ask"      // Requires human approval
	DecisionAllowOnce  PermissionDecision = "allow_once"
	DecisionDenyOnce   PermissionDecision = "deny_once"
)

// PermissionRule is a single rule in the permission system.
type PermissionRule struct {
	Action     ActionType         `json:"action"`
	Resource   string             `json:"resource,omitempty"`   // Glob pattern for file paths
	Decision   PermissionDecision `json:"decision"`
	Reason     string             `json:"reason,omitempty"`
	ToolName   string             `json:"toolName,omitempty"`   // Specific tool restriction
}

// PermissionManager manages the security permission system.
type PermissionManager struct {
	mu       sync.RWMutex
	Level    AutonomyLevel
	Rules    []PermissionRule
	WorkDir  string
	AllowDir string // Sandboxed directory (all ops restricted to this)
	DeniedCount int
	AllowedCount int
}

// NewPermissionManager creates a new permission manager.
func NewPermissionManager(level AutonomyLevel) *PermissionManager {
	pm := &PermissionManager{
		Level: level,
		Rules: DefaultPermissionRules(),
	}
	return pm
}

// DefaultPermissionRules returns the default set of permission rules.
func DefaultPermissionRules() []PermissionRule {
	return []PermissionRule{
		// File operations
		{Action: ActionFileRead, Decision: DecisionAllow, Reason: "reading files is safe"},
		{Action: ActionFileWrite, Decision: DecisionAsk, Reason: "writing files requires approval"},
		{Action: ActionFileDelete, Decision: DecisionAsk, Reason: "deleting files requires approval"},
		{Action: ActionFileMove, Decision: DecisionAsk, Reason: "moving files requires approval"},

		// Command execution
		{Action: ActionCommandExec, Decision: DecisionAsk, Reason: "executing commands requires approval"},
		{Action: ActionCodeExec, Decision: DecisionAsk, Reason: "executing code requires approval"},
		{Action: ActionGitOperation, Decision: DecisionAsk, Reason: "git operations require approval"},

		// Network
		{Action: ActionWebFetch, Decision: DecisionAllow, Reason: "web fetching is generally safe"},
		{Action: ActionWebSearch, Decision: DecisionAllow, Reason: "web search is safe"},
		{Action: ActionNetwork, Decision: DecisionAsk, Reason: "network access requires approval"},

		// Environment
		{Action: ActionEnvRead, Decision: DecisionAllow, Reason: "reading env vars is safe"},

		// MCP
		{Action: ActionMCPToolCall, Decision: DecisionAsk, Reason: "MCP tool calls require approval"},
	}
}

// Check evaluates a permission request and returns the decision.
func (pm *PermissionManager) Check(action ActionType, resource string) PermissionDecision {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	// God mode bypasses everything
	if pm.Level == AutonomyLevelGod {
		pm.AllowedCount++
		return DecisionAllow
	}

	// Check autonomy level first
	switch pm.Level {
	case AutonomyLevelExecute:
		if action == ActionFileRead || action == ActionFileWrite ||
			action == ActionCommandExec || action == ActionCodeExec ||
			action == ActionGitOperation {
			pm.AllowedCount++
			return DecisionAllow
		}
	case AutonomyLevelWrite:
		if action == ActionFileRead || action == ActionFileWrite {
			pm.AllowedCount++
			return DecisionAllow
		}
	case AutonomyLevelRead:
		if action == ActionFileRead {
			pm.AllowedCount++
			return DecisionAllow
		}
	case AutonomyLevelNone:
		// All require approval — fall through to rule check
		// but default to Ask instead of Allow
	}

	// Check path sandboxing
	if pm.AllowDir != "" && isPathAction(action) {
		if !isPathWithin(resource, pm.AllowDir) {
			pm.DeniedCount++
			return DecisionDeny
		}
	}

	// Check explicit rules (last match wins)
	decision := DecisionAsk // Default: ask
	for _, rule := range pm.Rules {
		if rule.Action != action {
			continue
		}
		if rule.Resource != "" && !matchGlob(resource, rule.Resource) {
			continue
		}
		decision = rule.Decision
	}

	if decision == DecisionAllow || decision == DecisionAllowOnce {
		pm.AllowedCount++
	} else {
		pm.DeniedCount++
	}

	return decision
}

// RequiresApproval determines if the current request must block for human input.
func (pm *PermissionManager) RequiresApproval(actionType string, resource string) bool {
	decision := pm.Check(ActionType(actionType), resource)
	return decision == DecisionAsk || decision == DecisionDeny || decision == DecisionDenyOnce
}

// AddRule adds a permission rule.
func (pm *PermissionManager) AddRule(rule PermissionRule) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.Rules = append(pm.Rules, rule)
}

// RemoveRule removes a permission rule by index.
func (pm *PermissionManager) RemoveRule(idx int) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	if idx >= 0 && idx < len(pm.Rules) {
		pm.Rules = append(pm.Rules[:idx], pm.Rules[idx+1:]...)
	}
}

// GetRules returns the current rules.
func (pm *PermissionManager) GetRules() []PermissionRule {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	result := make([]PermissionRule, len(pm.Rules))
	copy(result, pm.Rules)
	return result
}

// Stats returns permission check statistics.
func (pm *PermissionManager) Stats() (allowed, denied int) {
	return pm.AllowedCount, pm.DeniedCount
}

// ── Dangerous Command Detection (ported from goose/cline) ──

var dangerousPatterns = []struct {
	Pattern *regexp.Regexp
	Reason  string
}{
	{regexp.MustCompile(`(?i)\brm\s+(-rf|-fr)\s+/(etc|usr|var|boot|root)\b`), "recursive removal of system directories"},
	{regexp.MustCompile(`(?i)\brm\s+(-rf|-fr)\s+~`), "recursive removal of home directory"},
	{regexp.MustCompile(`(?i)\brm\s+(-rf|-fr)\s+/`), "recursive removal from root"},
	{regexp.MustCompile(`(?i)\bdd\s+if=`), "dd command can overwrite disks"},
	{regexp.MustCompile(`(?i)\bmkfs\b`), "filesystem formatting"},
	{regexp.MustCompile(`(?i)\bformat\s+[a-zA-Z]:`), "disk formatting (Windows)"},
	{regexp.MustCompile(`(?i)\b:(?:` + "`" + `|()%.*%>)\s*rm\b`), "fork bomb variant"},
	{regexp.MustCompile(`(?i)\bchmod\s+(-R\s+)?777\s+/`), "recursively setting world-writable on root"},
	{regexp.MustCompile(`(?i)\bchown\s+(-R\s+)?\w+\s+/`), "recursively changing ownership of root"},
	{regexp.MustCompile(`(?i)\bshutdown\b`), "system shutdown"},
	{regexp.MustCompile(`(?i)\breboot\b`), "system reboot"},
	{regexp.MustCompile(`(?i)\bhalt\b`), "system halt"},
	{regexp.MustCompile(`(?i)\bpoweroff\b`), "system power off"},
	{regexp.MustCompile(`(?i)\biptables\s+-F`), "flushing firewall rules"},
	{regexp.MustCompile(`(?i)\bcurl\b.*\|\s*(ba)?sh`), "piping remote content to shell"},
	{regexp.MustCompile(`(?i)\bwget\b.*\|\s*(ba)?sh`), "piping remote content to shell"},
	{regexp.MustCompile(`(?i)\beval\s+["']`), "eval with quoted string (potential injection)"},
	{regexp.MustCompile(`(?i)\b(iex|Invoke-Expression)\b`), "PowerShell command injection"},
	{regexp.MustCompile(`(?i)\bsudo\s+rm\b`), "sudo rm is dangerous"},
	{regexp.MustCompile(`(?i)\bgit\s+push\s+.*--force`), "force pushing to remote"},
	{regexp.MustCompile(`(?i)\bgit\s+reset\s+--hard`), "hard reset can lose uncommitted work"},
	{regexp.MustCompile(`(?i)\bdocker\s+rm\s+(-f|--force)`), "force removing docker containers"},
	{regexp.MustCompile(`(?i)\bdocker\s+(system|volume)\s+prune`), "docker prune can remove data"},
	{regexp.MustCompile(`(?i)\bkubectl\s+delete\s+--all`), "deleting all k8s resources"},
	{regexp.MustCompile(`(?i)\baws\s+.*\bdelete\b`), "AWS delete operations"},
	{regexp.MustCompile(`(?i)\bgcloud\s+.*\bdelete\b`), "GCloud delete operations"},
}

// CommandChecker evaluates shell commands for dangerous patterns.
type CommandChecker struct {
	extraPatterns []*regexp.Regexp
}

// NewCommandChecker creates a command safety checker.
func NewCommandChecker() *CommandChecker {
	return &CommandChecker{}
}

// AddPattern adds a custom dangerous pattern.
func (cc *CommandChecker) AddPattern(pattern string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid pattern: %w", err)
	}
	cc.extraPatterns = append(cc.extraPatterns, re)
	return nil
}

// Check evaluates a command for dangerous patterns.
func (cc *CommandChecker) Check(command string) (safe bool, reason string) {
	for _, dp := range dangerousPatterns {
		if dp.Pattern.MatchString(command) {
			return false, dp.Reason
		}
	}
	for _, ep := range cc.extraPatterns {
		if ep.MatchString(command) {
			return false, "matched custom dangerous pattern"
		}
	}
	return true, ""
}

// InterceptDangerousAction checks if a command is dangerous and should be blocked.
func (pm *PermissionManager) InterceptDangerousAction(command string) bool {
	checker := NewCommandChecker()
	safe, reason := checker.Check(command)
	if !safe {
		_ = reason // Log in production
		return true
	}
	return false
}

// ── Path Sandboxing (ported from claude-code) ──

// Sandbox restricts file operations to allowed directories.
type Sandbox struct {
	AllowedDirs []string
	DeniedDirs  []string
	WorkingDir  string
}

// NewSandbox creates a sandbox with the given allowed directories.
func NewSandbox(workingDir string) *Sandbox {
	return &Sandbox{
		WorkingDir:  workingDir,
		AllowedDirs: []string{workingDir},
		DeniedDirs:  defaultDeniedDirs(),
	}
}

func defaultDeniedDirs() []string {
	return []string{
		"/etc",
		"/usr",
		"/var",
		"/boot",
		"/root",
		"/System",
		"/Library",
		filepath.Join(os.Getenv("HOME"), ".ssh"),
		filepath.Join(os.Getenv("HOME"), ".gnupg"),
		filepath.Join(os.Getenv("HOME"), ".password-store"),
	}
}

// AllowDir adds a directory to the allowed list.
func (s *Sandbox) AllowDir(path string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return
	}
	s.AllowedDirs = append(s.AllowedDirs, absPath)
}

// DenyDir adds a directory to the denied list.
func (s *Sandbox) DenyDir(path string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return
	}
	s.DeniedDirs = append(s.DeniedDirs, absPath)
}

// IsPathAllowed checks if a path is within an allowed directory.
func (s *Sandbox) IsPathAllowed(path string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	// Check denied first
	for _, denied := range s.DeniedDirs {
		if strings.HasPrefix(absPath, denied) {
			return false
		}
	}

	// Check allowed
	for _, allowed := range s.AllowedDirs {
		if strings.HasPrefix(absPath, allowed) {
			return true
		}
	}

	return false
}

// ValidatePath validates and returns a safe path within the sandbox.
func (s *Sandbox) ValidatePath(path string) (string, error) {
	// Resolve relative paths against working dir
	var absPath string
	if filepath.IsAbs(path) {
		absPath = path
	} else {
		absPath = filepath.Join(s.WorkingDir, path)
	}

	// Clean the path
	absPath = filepath.Clean(absPath)

	if !s.IsPathAllowed(absPath) {
		return "", fmt.Errorf("path %s is outside allowed directories", path)
	}

	return absPath, nil
}

// ── Helper functions ──

func isPathAction(action ActionType) bool {
	return action == ActionFileRead || action == ActionFileWrite ||
		action == ActionFileDelete || action == ActionFileMove
}

func isPathWithin(path, base string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	absBase, err := filepath.Abs(base)
	if err != nil {
		return false
	}
	return strings.HasPrefix(absPath, absBase)
}

func matchGlob(s, pattern string) bool {
	matched, err := filepath.Match(pattern, s)
	if err != nil {
		return false
	}
	return matched
}
