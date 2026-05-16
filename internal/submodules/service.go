package submodules

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type SubmoduleStatus struct {
	Path        string    `json:"path"`
	URL         string    `json:"url"`
	Branch      string    `json:"branch"`
	CurrentSHA  string    `json:"currentSha"`
	LatestSHA   string    `json:"latestSha"`
	IsDirty     bool      `json:"isDirty"`
	NeedsUpdate bool      `json:"needsUpdate"`
	LastChecked time.Time `json:"lastChecked"`
}

type SubmoduleService struct {
	rootPath string
}

func NewSubmoduleService(rootPath string) *SubmoduleService {
	return &SubmoduleService{rootPath: rootPath}
}

func (ss *SubmoduleService) List() ([]SubmoduleStatus, error) {
	output, err := ss.git("submodule", "status", "--recursive")
	if err != nil {
		return nil, err
	}
	var subs []SubmoduleStatus
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		prefix := ""
		if len(line) > 0 && (line[0] == '+' || line[0] == '-' || line[0] == ' ') {
			prefix = string(line[0])
			line = line[1:]
		}
		parts := strings.SplitN(line, " ", 3)
		if len(parts) < 2 {
			continue
		}
		subs = append(subs, SubmoduleStatus{
			Path:        parts[1],
			CurrentSHA:  parts[0],
			IsDirty:     prefix == "+",
			LastChecked: time.Now().UTC(),
		})
	}
	return subs, nil
}

func (ss *SubmoduleService) UpdateAll() ([]SubmoduleStatus, error) {
	if _, err := ss.git("submodule", "update", "--remote", "--merge"); err != nil {
		return nil, err
	}
	return ss.List()
}

func (ss *SubmoduleService) Update(path string) error {
	_, err := ss.git("submodule", "update", "--remote", "--merge", path)
	return err
}

func (ss *SubmoduleService) Sync() error {
	_, err := ss.git("submodule", "sync", "--recursive")
	return err
}

func (ss *SubmoduleService) Fetch() error {
	_, err := ss.git("submodule", "foreach", "--recursive", "git fetch")
	return err
}

func (ss *SubmoduleService) CheckUpdates() ([]SubmoduleStatus, error) {
	subs, err := ss.List()
	if err != nil {
		return nil, err
	}
	for i := range subs {
		branch := subs[i].Branch
		if branch == "" {
			branch = "main"
		}
		subPath := filepath.Join(ss.rootPath, subs[i].Path)
		if latest, err := exec.Command("git", "-C", subPath, "rev-parse", "origin/"+branch).Output(); err == nil {
			subs[i].LatestSHA = strings.TrimSpace(string(latest))
			subs[i].NeedsUpdate = subs[i].LatestSHA != subs[i].CurrentSHA
		}
	}
	return subs, nil
}

func (ss *SubmoduleService) Diff(path string) (string, error) {
	subPath := filepath.Join(ss.rootPath, path)
	output, err := exec.Command("git", "-C", subPath, "log", "--oneline", "HEAD..origin/main").Output()
	if err != nil {
		output, err = exec.Command("git", "-C", subPath, "log", "--oneline", "HEAD..origin/master").Output()
		if err != nil {
			return "", err
		}
	}
	return string(output), nil
}

func (ss *SubmoduleService) git(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = ss.rootPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %s: %w", strings.Join(args, " "), string(output), err)
	}
	return string(output), nil
}
