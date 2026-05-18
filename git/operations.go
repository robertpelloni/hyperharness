package git

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// ═══════════════════════════════════════════════════════════════════════
// operations.go — Git operations for HyperHarness
// Unified git operations from claude-code, aider, and goose patterns.
//
// WHAT: High-level git operations: status, diff, log, stash, branch, etc.
// WHY: Every AI coding tool needs git integration for safe code changes.
// HOW: Shell out to git CLI for reliability (matches claude-code approach).
// ═══════════════════════════════════════════════════════════════════════

// GitOps provides high-level git operations.
type GitOps struct {
	WorkingDir string
}

// NewGitOps creates a new git operations helper.
func NewGitOps(workingDir string) *GitOps {
	return &GitOps{WorkingDir: workingDir}
}

func (g *GitOps) run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = g.WorkingDir
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// ── Status ──

// Status returns the output of git status --porcelain.
func (g *GitOps) Status() (string, error) {
	return g.run("status", "--porcelain")
}

// IsClean returns true if there are no uncommitted changes.
func (g *GitOps) IsClean() (bool, error) {
	out, err := g.Status()
	if err != nil {
		return false, err
	}
	return out == "", nil
}

// HasUncommittedChanges checks for uncommitted changes.
func (g *GitOps) HasUncommittedChanges() (bool, error) {
	clean, err := g.IsClean()
	return !clean, err
}

// ChangedFiles returns a list of changed files.
func (g *GitOps) ChangedFiles() ([]string, error) {
	out, err := g.Status()
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}

	var files []string
	for _, line := range strings.Split(out, "\n") {
		if len(line) >= 4 {
			file := strings.TrimSpace(line[3:])
			files = append(files, file)
		}
	}
	return files, nil
}

// StagedFiles returns files that are staged for commit.
func (g *GitOps) StagedFiles() ([]string, error) {
	out, err := g.run("diff", "--cached", "--name-only")
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	return strings.Split(out, "\n"), nil
}

// UnstagedFiles returns files that have unstaged changes.
func (g *GitOps) UnstagedFiles() ([]string, error) {
	out, err := g.run("diff", "--name-only")
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	return strings.Split(out, "\n"), nil
}

// UntrackedFiles returns files not tracked by git.
func (g *GitOps) UntrackedFiles() ([]string, error) {
	out, err := g.run("ls-files", "--others", "--exclude-standard")
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	return strings.Split(out, "\n"), nil
}

// ── Diff ──

// Diff returns the full diff of uncommitted changes.
func (g *GitOps) Diff() (string, error) {
	return g.run("diff")
}

// DiffStaged returns the diff of staged changes.
func (g *GitOps) DiffStaged() (string, error) {
	return g.run("diff", "--cached")
}

// DiffFile returns the diff for a specific file.
func (g *GitOps) DiffFile(file string) (string, error) {
	return g.run("diff", file)
}

// DiffStat returns a summary of changes.
func (g *GitOps) DiffStat() (string, error) {
	return g.run("diff", "--stat")
}

// ── Add & Commit ──

// Add stages files for commit.
func (g *GitOps) Add(files ...string) error {
	args := append([]string{"add"}, files...)
	_, err := g.run(args...)
	return err
}

// AddAll stages all changes.
func (g *GitOps) AddAll() error {
	_, err := g.run("add", "-A")
	return err
}

// Commit creates a commit with the given message.
func (g *GitOps) Commit(message string) (string, error) {
	return g.run("commit", "-m", message)
}

// CommitAll stages all changes and commits.
func (g *GitOps) CommitAll(message string) (string, error) {
	if err := g.AddAll(); err != nil {
		return "", err
	}
	return g.Commit(message)
}

// AmendLast amends the last commit.
func (g *GitOps) AmendLast(message string) (string, error) {
	return g.run("commit", "--amend", "-m", message)
}

// ── Branch ──

// CurrentBranch returns the name of the current branch.
func (g *GitOps) CurrentBranch() (string, error) {
	return g.run("rev-parse", "--abbrev-ref", "HEAD")
}

// Branches returns a list of all branches.
func (g *GitOps) Branches() ([]string, error) {
	out, err := g.run("branch", "--format=%(refname:short)")
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	return strings.Split(out, "\n"), nil
}

// CreateBranch creates a new branch.
func (g *GitOps) CreateBranch(name string) error {
	_, err := g.run("checkout", "-b", name)
	return err
}

// SwitchBranch switches to a different branch.
func (g *GitOps) SwitchBranch(name string) error {
	_, err := g.run("checkout", name)
	return err
}

// DeleteBranch deletes a branch.
func (g *GitOps) DeleteBranch(name string) error {
	_, err := g.run("branch", "-d", name)
	return err
}

// ── Log ──

// LogEntry represents a single git log entry.
type LogEntry struct {
	Hash    string
	Author  string
	Date    time.Time
	Message string
}

// Log returns recent log entries.
func (g *GitOps) Log(count int) ([]LogEntry, error) {
	format := "--pretty=format:%H|%an|%aI|%s"
	out, err := g.run("log", "-n", strconv.Itoa(count), format)
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}

	var entries []LogEntry
	for _, line := range strings.Split(out, "\n") {
		parts := strings.SplitN(line, "|", 4)
		if len(parts) != 4 {
			continue
		}
		date, _ := time.Parse(time.RFC3339, parts[2])
		entries = append(entries, LogEntry{
			Hash:    parts[0],
			Author:  parts[1],
			Date:    date,
			Message: parts[3],
		})
	}
	return entries, nil
}

// LogOneline returns a one-line-per-commit log.
func (g *GitOps) LogOneline(count int) (string, error) {
	return g.run("log", "-n", strconv.Itoa(count), "--oneline")
}

// LastCommit returns the last commit hash.
func (g *GitOps) LastCommit() (string, error) {
	return g.run("rev-parse", "HEAD")
}

// ── Stash ──

// Stash saves working changes to the stash.
func (g *GitOps) Stash(message string) error {
	if message != "" {
		_, err := g.run("stash", "push", "-m", message)
		return err
	}
	_, err := g.run("stash")
	return err
}

// StashPop restores the most recent stashed changes.
func (g *GitOps) StashPop() error {
	_, err := g.run("stash", "pop")
	return err
}

// StashList returns the list of stashes.
func (g *GitOps) StashList() (string, error) {
	return g.run("stash", "list")
}

// ── Remote ──

// Push pushes changes to the remote.
func (g *GitOps) Push() error {
	_, err := g.run("push")
	return err
}

// Pull pulls changes from the remote.
func (g *GitOps) Pull() error {
	_, err := g.run("pull")
	return err
}

// Fetch fetches from the remote without merging.
func (g *GitOps) Fetch() error {
	_, err := g.run("fetch")
	return err
}

// RemoteURL returns the URL of the remote.
func (g *GitOps) RemoteURL() (string, error) {
	return g.run("remote", "get-url", "origin")
}

// ── Reset & Revert ──

// ResetSoft does a soft reset to the given commit.
func (g *GitOps) ResetSoft(commit string) error {
	_, err := g.run("reset", "--soft", commit)
	return err
}

// ResetHard does a hard reset to the given commit.
func (g *GitOps) ResetHard(commit string) error {
	_, err := g.run("reset", "--hard", commit)
	return err
}

// Revert creates a revert commit.
func (g *GitOps) Revert(commit string) error {
	_, err := g.run("revert", "--no-edit", commit)
	return err
}

// ── Utility ──

// IsRepo checks if the directory is a git repo.
func (g *GitOps) IsRepo() bool {
	_, err := g.run("rev-parse", "--git-dir")
	return err == nil
}

// RootDir returns the root directory of the repo.
func (g *GitOps) RootDir() (string, error) {
	return g.run("rev-parse", "--show-toplevel")
}

// FileHistory returns the commit history for a file.
func (g *GitOps) FileHistory(file string, count int) ([]LogEntry, error) {
	format := "--pretty=format:%H|%an|%aI|%s"
	out, err := g.run("log", "-n", strconv.Itoa(count), format, "--", file)
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}

	var entries []LogEntry
	for _, line := range strings.Split(out, "\n") {
		parts := strings.SplitN(line, "|", 4)
		if len(parts) != 4 {
			continue
		}
		date, _ := time.Parse(time.RFC3339, parts[2])
		entries = append(entries, LogEntry{
			Hash:    parts[0],
			Author:  parts[1],
			Date:    date,
			Message: parts[3],
		})
	}
	return entries, nil
}

// ShowFileAtCommit shows a file at a specific commit.
func (g *GitOps) ShowFileAtCommit(file, commit string) (string, error) {
	return g.run("show", fmt.Sprintf("%s:%s", commit, file))
}

// Blame returns blame info for a file.
func (g *GitOps) Blame(file string) (string, error) {
	return g.run("blame", file)
}

// CountCommits returns the number of commits.
func (g *GitOps) CountCommits() (int, error) {
	out, err := g.run("rev-list", "--count", "HEAD")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(out)
}

// ShortHash returns the short hash for a commit.
func (g *GitOps) ShortHash(commit string) (string, error) {
	return g.run("rev-parse", "--short", commit)
}

// FileAtLine returns the file and line number from a path:line spec.
func FileAtLine(spec string) (file string, line int) {
	parts := strings.SplitN(spec, ":", 2)
	file = parts[0]
	if len(parts) > 1 {
		line, _ = strconv.Atoi(parts[1])
	}
	return
}

// RelPath returns the relative path from the git root.
func (g *GitOps) RelPath(absPath string) string {
	root, err := g.RootDir()
	if err != nil {
		return absPath
	}
	rel, err := filepath.Rel(root, absPath)
	if err != nil {
		return absPath
	}
	return rel
}
