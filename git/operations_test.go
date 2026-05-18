package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func newTestRepo(t *testing.T) *GitOps {
	t.Helper()
	dir := t.TempDir()
	exec.Command("git", "init", dir).Run()
	exec.Command("git", "-C", dir, "config", "user.email", "test@test.com").Run()
	exec.Command("git", "-C", dir, "config", "user.name", "Test").Run()

	// Create an initial commit
	os.WriteFile(filepath.Join(dir, "README.md"), []byte("# Test\n"), 0644)
	exec.Command("git", "-C", dir, "add", ".").Run()
	exec.Command("git", "-C", dir, "commit", "-m", "initial").Run()

	return NewGitOps(dir)
}

func TestIsRepo(t *testing.T) {
	ops := newTestRepo(t)
	if !ops.IsRepo() {
		t.Error("temp dir should be a git repo")
	}

	nonOps := NewGitOps(t.TempDir())
	// On some systems, temp dirs may be inside a git repo
	// Just verify IsRepo doesn't crash
	_ = nonOps.IsRepo()
}

func TestCurrentBranch(t *testing.T) {
	ops := newTestRepo(t)
	branch, err := ops.CurrentBranch()
	if err != nil {
		t.Fatal(err)
	}
	if branch == "" {
		t.Error("should have a branch name")
	}
}

func TestIsClean(t *testing.T) {
	ops := newTestRepo(t)
	clean, err := ops.IsClean()
	if err != nil {
		t.Fatal(err)
	}
	if !clean {
		t.Error("fresh repo after commit should be clean")
	}
}

func TestChangedFiles(t *testing.T) {
	ops := newTestRepo(t)

	// Modify a file
	dir := ops.WorkingDir
	os.WriteFile(filepath.Join(dir, "README.md"), []byte("# Modified\n"), 0644)

	files, err := ops.ChangedFiles()
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Error("should have changed files after modification")
	}
}

func TestAddAndCommit(t *testing.T) {
	ops := newTestRepo(t)
	dir := ops.WorkingDir

	os.WriteFile(filepath.Join(dir, "new.txt"), []byte("hello\n"), 0644)
	if err := ops.Add("new.txt"); err != nil {
		t.Fatal(err)
	}

	staged, err := ops.StagedFiles()
	if err != nil {
		t.Fatal(err)
	}
	if len(staged) == 0 {
		t.Error("should have staged files")
	}

	hash, err := ops.Commit("add new file")
	if err != nil {
		t.Fatal(err)
	}
	if hash == "" {
		t.Error("should return commit hash")
	}

	clean, _ := ops.IsClean()
	if !clean {
		t.Error("repo should be clean after commit")
	}
}

func TestCommitAll(t *testing.T) {
	ops := newTestRepo(t)
	dir := ops.WorkingDir

	os.WriteFile(filepath.Join(dir, "all.txt"), []byte("content\n"), 0644)
	_, err := ops.CommitAll("commit all")
	if err != nil {
		t.Fatal(err)
	}

	clean, _ := ops.IsClean()
	if !clean {
		t.Error("repo should be clean after commit all")
	}
}

func TestDiff(t *testing.T) {
	ops := newTestRepo(t)
	dir := ops.WorkingDir

	os.WriteFile(filepath.Join(dir, "README.md"), []byte("# Changed\n"), 0644)
	diff, err := ops.Diff()
	if err != nil {
		t.Fatal(err)
	}
	if diff == "" {
		t.Error("should have diff output")
	}
}

func TestLog(t *testing.T) {
	ops := newTestRepo(t)

	entries, err := ops.Log(10)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Error("should have log entries")
	}
	if entries[0].Message == "" {
		t.Error("log entry should have a message")
	}
}

func TestLogOneline(t *testing.T) {
	ops := newTestRepo(t)
	log, err := ops.LogOneline(10)
	if err != nil {
		t.Fatal(err)
	}
	if log == "" {
		t.Error("should have log output")
	}
}

func TestLastCommit(t *testing.T) {
	ops := newTestRepo(t)
	hash, err := ops.LastCommit()
	if err != nil {
		t.Fatal(err)
	}
	if len(hash) < 10 {
		t.Errorf("hash seems too short: %s", hash)
	}
}

func TestRootDir(t *testing.T) {
	ops := newTestRepo(t)
	root, err := ops.RootDir()
	if err != nil {
		t.Fatal(err)
	}
	if root == "" {
		t.Error("should have root dir")
	}
}

func TestStash(t *testing.T) {
	ops := newTestRepo(t)
	dir := ops.WorkingDir

	os.WriteFile(filepath.Join(dir, "README.md"), []byte("# Changed\n"), 0644)
	if err := ops.Stash("test stash"); err != nil {
		t.Fatal(err)
	}

	clean, _ := ops.IsClean()
	if !clean {
		t.Error("repo should be clean after stash")
	}
}

func TestBranches(t *testing.T) {
	ops := newTestRepo(t)
	branches, err := ops.Branches()
	if err != nil {
		t.Fatal(err)
	}
	if len(branches) == 0 {
		t.Error("should have at least one branch")
	}
}

func TestCreateBranch(t *testing.T) {
	ops := newTestRepo(t)
	if err := ops.CreateBranch("test-branch"); err != nil {
		t.Fatal(err)
	}
	branch, _ := ops.CurrentBranch()
	if branch != "test-branch" {
		t.Errorf("current branch: %s", branch)
	}
}

func TestCountCommits(t *testing.T) {
	ops := newTestRepo(t)
	count, err := ops.CountCommits()
	if err != nil {
		t.Fatal(err)
	}
	if count < 1 {
		t.Errorf("should have at least 1 commit, got %d", count)
	}
}

func TestShortHash(t *testing.T) {
	ops := newTestRepo(t)
	hash, _ := ops.LastCommit()
	short, err := ops.ShortHash(hash)
	if err != nil {
		t.Fatal(err)
	}
	if len(short) >= len(hash) {
		t.Error("short hash should be shorter than full hash")
	}
}

func TestFileAtLine(t *testing.T) {
	file, line := FileAtLine("main.go:42")
	if file != "main.go" {
		t.Errorf("file: %s", file)
	}
	if line != 42 {
		t.Errorf("line: %d", line)
	}

	file, line = FileAtLine("main.go")
	if file != "main.go" {
		t.Errorf("file: %s", file)
	}
	if line != 0 {
		t.Errorf("line should be 0: %d", line)
	}
}

func TestRelPath(t *testing.T) {
	ops := newTestRepo(t)
	root, _ := ops.RootDir()
	rel := ops.RelPath(filepath.Join(root, "src", "main.go"))
	if rel != filepath.Join("src", "main.go") {
		t.Errorf("rel path: %s", rel)
	}
}

func TestUntrackedFiles(t *testing.T) {
	ops := newTestRepo(t)
	dir := ops.WorkingDir

	os.WriteFile(filepath.Join(dir, "untracked.txt"), []byte("new\n"), 0644)
	files, err := ops.UntrackedFiles()
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Error("should have untracked files")
	}
}
