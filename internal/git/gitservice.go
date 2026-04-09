// Package git provides native Go git operations.
// Ported from hypercode/go/internal/gitservice/gitservice.go
//
// WHAT: Full git CLI wrapper for log, status, diff, blame, stash, branch, remote operations
// WHY: Git integration is core to any coding harness - agents need git ops as first-class tools
// HOW: exec.Command wrapper with structured output parsing, thread-safe
package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// CommitLog represents a single git commit entry.
type CommitLog struct {
	Hash    string `json:"hash"`
	Author  string `json:"author"`
	Date    string `json:"date"`
	Message string `json:"message"`
}

// GitStatus represents the working tree status.
type GitStatus struct {
	Branch   string   `json:"branch"`
	Clean    bool     `json:"clean"`
	Modified []string `json:"modified"`
	Staged   []string `json:"staged"`
	Untracked []string `json:"untracked"`
	Ahead    int      `json:"ahead"`
	Behind   int      `json:"behind"`
}

// DiffEntry represents a single file diff.
type DiffEntry struct {
	File     string `json:"file"`
	Added    int    `json:"added"`
	Removed  int    `json:"removed"`
	IsBinary bool   `json:"isBinary"`
}

// BlameLine represents a single line from git blame.
type BlameLine struct {
	Hash    string `json:"hash"`
	Author  string `json:"author"`
	Date    string `json:"date"`
	LineNo  int    `json:"lineNo"`
	Content string `json:"content"`
}

// GitService provides native git operations.
type GitService struct {
	cwd string
}

// NewGitService creates a new GitService for the given working directory.
func NewGitService(cwd string) *GitService {
	return &GitService{cwd: cwd}
}

func (gs *GitService) run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = gs.cwd
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %s: %w", strings.Join(args, " "), strings.TrimSpace(stderr.String()), err)
	}
	return stdout.String(), nil
}

// GetLog returns the last `limit` commits.
func (gs *GitService) GetLog(limit int) ([]CommitLog, error) {
	if limit <= 0 {
		limit = 20
	}
	out, err := gs.run("log", "-n", strconv.Itoa(limit), "--pretty=format:%H|%an|%aI|%s")
	if err != nil {
		return nil, err
	}
	var logs []CommitLog
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 4)
		if len(parts) < 4 {
			continue
		}
		logs = append(logs, CommitLog{
			Hash: parts[0], Author: parts[1], Date: parts[2], Message: parts[3],
		})
	}
	return logs, nil
}

// GetStatus returns the current working tree status.
func (gs *GitService) GetStatus() (*GitStatus, error) {
	branch, err := gs.run("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return &GitStatus{Branch: "unknown"}, nil
	}
	branch = strings.TrimSpace(branch)

	statusOut, err := gs.run("status", "--porcelain")
	if err != nil {
		return &GitStatus{Branch: branch}, nil
	}

	var modified, staged, untracked []string
	for _, line := range strings.Split(statusOut, "\n") {
		if len(line) < 4 {
			continue
		}
		code := line[:2]
		file := strings.TrimSpace(line[3:])
		if code[1] == '?' || (len(code) > 0 && code[0] == '?') {
			untracked = append(untracked, file)
		}
		if strings.Contains(code, "M") {
			modified = append(modified, file)
		}
		if strings.Contains(code, "A") || (code[0] == 'M') {
			staged = append(staged, file)
		}
	}

	status := &GitStatus{
		Branch:    branch,
		Clean:     strings.TrimSpace(statusOut) == "",
		Modified:  modified,
		Staged:    staged,
		Untracked: untracked,
	}

	aheadBehind, err := gs.run("rev-list", "--left-right", "--count", "HEAD...@{upstream}")
	if err == nil {
		parts := strings.Fields(aheadBehind)
		if len(parts) == 2 {
			status.Ahead, _ = strconv.Atoi(parts[0])
			status.Behind, _ = strconv.Atoi(parts[1])
		}
	}
	return status, nil
}

// Revert creates a new commit that undoes a previous commit.
func (gs *GitService) Revert(hash string) (string, error) {
	out, err := gs.run("revert", "--no-edit", hash)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// ResetTo resets HEAD to the given commit.
func (gs *GitService) ResetTo(hash string, hard bool) (string, error) {
	mode := "--soft"
	if hard {
		mode = "--hard"
	}
	out, err := gs.run("reset", mode, hash)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// Diff returns the diff for the working tree or staged changes.
func (gs *GitService) Diff(staged bool) ([]DiffEntry, error) {
	args := []string{"diff", "--numstat"}
	if staged {
		args = append(args, "--cached")
	}
	out, err := gs.run(args...)
	if err != nil {
		return nil, err
	}
	var entries []DiffEntry
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}
		e := DiffEntry{File: parts[2]}
		if parts[0] == "-" || parts[1] == "-" {
			e.IsBinary = true
		} else {
			e.Added, _ = strconv.Atoi(parts[0])
			e.Removed, _ = strconv.Atoi(parts[1])
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// DiffStat returns a short diffstat summary.
func (gs *GitService) DiffStat(staged bool) (string, error) {
	args := []string{"diff", "--stat"}
	if staged {
		args = append(args, "--cached")
	}
	return gs.run(args...)
}

// Blame returns line-by-line blame for a file.
func (gs *GitService) Blame(file string) ([]BlameLine, error) {
	out, err := gs.run("blame", "--porcelain", file)
	if err != nil {
		return nil, err
	}

	var lines []BlameLine
	var current *BlameLine
	lineNo := 0

	for _, line := range strings.Split(out, "\n") {
		if line == "" {
			continue
		}
		if line[0] >= '0' && line[0] <= '9' && line[0] != '\t' {
			parts := strings.Fields(line)
			if current != nil {
				lines = append(lines, *current)
			}
			lineNo++
			current = &BlameLine{Hash: parts[0], LineNo: lineNo}
		} else if strings.HasPrefix(line, "author ") {
			if current != nil {
				current.Author = strings.TrimPrefix(line, "author ")
			}
		} else if strings.HasPrefix(line, "committer-time ") {
			if current != nil {
				current.Date = strings.TrimPrefix(line, "committer-time ")
			}
		} else if strings.HasPrefix(line, "\t") {
			if current != nil {
				current.Content = strings.TrimPrefix(line, "\t")
			}
		}
	}
	if current != nil {
		lines = append(lines, *current)
	}
	return lines, nil
}

// Stash saves local modifications to the stash.
func (gs *GitService) Stash(message string) (string, error) {
	args := []string{"stash", "push"}
	if message != "" {
		args = append(args, "-m", message)
	}
	out, err := gs.run(args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// StashPop applies and removes the most recent stash entry.
func (gs *GitService) StashPop() (string, error) {
	out, err := gs.run("stash", "pop")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// ListBranches returns all local branches.
func (gs *GitService) ListBranches() ([]string, error) {
	out, err := gs.run("branch", "--list")
	if err != nil {
		return nil, err
	}
	var branches []string
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "* ")
		line = strings.TrimSpace(line)
		if line != "" {
			branches = append(branches, line)
		}
	}
	return branches, nil
}

// GetCurrentBranch returns the name of the current branch.
func (gs *GitService) GetCurrentBranch() (string, error) {
	out, err := gs.run("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// IsGitRepo returns true if the cwd is inside a git repository.
func (gs *GitService) IsGitRepo() bool {
	_, err := gs.run("rev-parse", "--git-dir")
	return err == nil
}

// GetRemoteURL returns the URL for the given remote name.
func (gs *GitService) GetRemoteURL(name string) (string, error) {
	out, err := gs.run("remote", "get-url", name)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// Add stages files for the next commit.
func (gs *GitService) Add(files ...string) error {
	args := append([]string{"add"}, files...)
	_, err := gs.run(args...)
	return err
}

// Commit creates a new commit with the given message.
func (gs *GitService) Commit(message string) (string, error) {
	out, err := gs.run("commit", "-m", message)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// Fetch fetches from all remotes.
func (gs *GitService) Fetch(all bool) error {
	args := []string{"fetch"}
	if all {
		args = append(args, "--all")
	}
	_, err := gs.run(args...)
	return err
}

// Pull pulls from the current branch's upstream.
func (gs *GitService) Pull() (string, error) {
	out, err := gs.run("pull")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// Push pushes to the given remote and branch.
func (gs *GitService) Push(remote, branch string) (string, error) {
	out, err := gs.run("push", remote, branch)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// Show returns the content of a commit or file at a given ref.
func (gs *GitService) Show(ref string) (string, error) {
	out, err := gs.run("show", ref)
	if err != nil {
		return "", err
	}
	return out, nil
}

// Init initializes a new git repository.
func (gs *GitService) Init() error {
	_, err := gs.run("init")
	return err
}

// Clone clones a repository into the given directory.
func (gs *GitService) Clone(url, dir string) error {
	_, err := gs.run("clone", url, dir)
	return err
}
