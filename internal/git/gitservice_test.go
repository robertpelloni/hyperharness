package git

import (
	"os"
	"os/exec"
	"testing"
)

func TestNewGitService(t *testing.T) {
	gs := NewGitService(".")
	if gs == nil {
		t.Fatal("should create git service")
	}
}

func TestIsGitRepo(t *testing.T) {
	// The hyperharness repo should be a git repo
	gs := NewGitService(".")
	if !gs.IsGitRepo() {
		t.Error("should be in a git repo")
	}

	// Temp dir outside the repo should not be a git repo
	tmpDir := t.TempDir()
	gs2 := NewGitService(tmpDir)
	// This may or may not be a git repo depending on if tmpdir is inside the repo
	// Just verify it doesn't crash
	_ = gs2.IsGitRepo()
}

func TestGetCurrentBranch(t *testing.T) {
	gs := NewGitService(".")
	branch, err := gs.GetCurrentBranch()
	if err != nil {
		t.Skipf("can't get branch: %v", err)
	}
	if branch == "" {
		t.Error("branch should not be empty")
	}
	t.Logf("current branch: %s", branch)
}

func TestGetStatus(t *testing.T) {
	gs := NewGitService(".")
	status, err := gs.GetStatus()
	if err != nil {
		t.Skipf("can't get status: %v", err)
	}
	if status.Branch == "" {
		t.Error("branch should be set")
	}
	t.Logf("branch=%s clean=%v modified=%d staged=%d",
		status.Branch, status.Clean, len(status.Modified), len(status.Staged))
}

func TestGetLog(t *testing.T) {
	gs := NewGitService(".")
	logs, err := gs.GetLog(5)
	if err != nil {
		t.Skipf("can't get log: %v", err)
	}
	if len(logs) == 0 {
		t.Error("should have commits")
	}
	t.Logf("latest commit: %s by %s: %s", logs[0].Hash[:8], logs[0].Author, logs[0].Message)
}

func TestListBranches(t *testing.T) {
	gs := NewGitService(".")
	branches, err := gs.ListBranches()
	if err != nil {
		t.Skipf("can't list branches: %v", err)
	}
	if len(branches) == 0 {
		t.Error("should have branches")
	}
	t.Logf("branches: %v", branches)
}

func TestGetRemoteURL(t *testing.T) {
	gs := NewGitService(".")
	url, err := gs.GetRemoteURL("origin")
	if err != nil {
		t.Skipf("no origin remote: %v", err)
	}
	if url == "" {
		t.Error("URL should not be empty")
	}
	t.Logf("origin: %s", url)
}

func TestDiff(t *testing.T) {
	gs := NewGitService(".")
	entries, err := gs.Diff(false)
	if err != nil {
		t.Skipf("diff error: %v", err)
	}
	// May be empty if working tree is clean
	t.Logf("diff entries: %d", len(entries))
}

func TestDiffStat(t *testing.T) {
	gs := NewGitService(".")
	stat, err := gs.DiffStat(false)
	if err != nil {
		t.Skipf("diff stat error: %v", err)
	}
	t.Logf("diff stat length: %d", len(stat))
}

func TestParseSubmoduleStatus(t *testing.T) {
	output := " abc123d4 superai/adrenaline (v1.0.0)\n def456a1 superai/aider (heads/main)\n"
	paths := parseSubmoduleStatus(output)
	if len(paths) != 2 {
		t.Fatalf("expected 2 paths, got %d", len(paths))
	}
	if paths[0] != "superai/adrenaline" {
		t.Errorf("first path: %s", paths[0])
	}
}

func TestParseSubmoduleStatusEmpty(t *testing.T) {
	paths := parseSubmoduleStatus("")
	if len(paths) != 0 {
		t.Error("should be empty")
	}
}

func TestGitServiceInTempDir(t *testing.T) {
	tmpDir := t.TempDir()
	gs := NewGitService(tmpDir)

	// Init a repo
	if err := gs.Init(); err != nil {
		t.Fatal(err)
	}
	if !gs.IsGitRepo() {
		t.Error("should be a git repo after init")
	}

	// Write a file and commit
	if err := os.WriteFile(tmpDir+"/test.txt", []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := gs.Add("test.txt"); err != nil {
		t.Fatal(err)
	}
	// Need to configure user for commit
	execInDir(tmpDir, "git", "config", "user.email", "test@test.com")
	execInDir(tmpDir, "git", "config", "user.name", "Test")

	msg, err := gs.Commit("initial commit")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("commit: %s", msg)

	logs, _ := gs.GetLog(1)
	if len(logs) != 1 {
		t.Fatal("should have 1 commit")
	}
	if logs[0].Message != "initial commit" {
		t.Errorf("message: %s", logs[0].Message)
	}
}

func execInDir(dir string, args ...string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Run()
}
