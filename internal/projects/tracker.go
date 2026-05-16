package projects

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

type Project struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Path         string            `json:"path"`
	Language     string            `json:"language,omitempty"`
	Framework    string            `json:"framework,omitempty"`
	Description  string            `json:"description,omitempty"`
	Tags         []string          `json:"tags,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	LastActiveAt time.Time         `json:"lastActiveAt"`
	SessionCount int               `json:"sessionCount"`
	FileCount    int               `json:"fileCount"`
	GitBranch    string            `json:"gitBranch,omitempty"`
	GitRemote    string            `json:"gitRemote,omitempty"`
	DiscoveredAt time.Time         `json:"discoveredAt"`
}

type ProjectTracker struct {
	mu        sync.RWMutex
	projects  map[string]*Project
	storePath string
}

func NewProjectTracker(storePath string) *ProjectTracker {
	pt := &ProjectTracker{projects: make(map[string]*Project), storePath: storePath}
	if storePath != "" {
		_ = pt.loadFromDisk()
	}
	return pt
}

func (pt *ProjectTracker) Register(path string) (*Project, error) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}
	if existing, ok := pt.projects[absPath]; ok {
		existing.LastActiveAt = time.Now().UTC()
		existing.SessionCount++
		_ = pt.saveToDisk()
		return existing, nil
	}
	project := &Project{
		ID:           generateProjectID(absPath),
		Name:         filepath.Base(absPath),
		Path:         absPath,
		LastActiveAt: time.Now().UTC(),
		DiscoveredAt: time.Now().UTC(),
		Metadata:     make(map[string]string),
		SessionCount: 1,
	}
	pt.detectProjectType(absPath, project)
	pt.projects[absPath] = project
	_ = pt.saveToDisk()
	return project, nil
}

func (pt *ProjectTracker) Get(path string) (*Project, bool) {
	pt.mu.RLock()
	defer pt.mu.RUnlock()
	absPath, _ := filepath.Abs(path)
	p, ok := pt.projects[absPath]
	return p, ok
}

func (pt *ProjectTracker) List() []*Project {
	pt.mu.RLock()
	defer pt.mu.RUnlock()
	var projects []*Project
	for _, p := range pt.projects {
		projects = append(projects, p)
	}
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].LastActiveAt.After(projects[j].LastActiveAt)
	})
	return projects
}

func (pt *ProjectTracker) Recent(n int) []*Project {
	all := pt.List()
	if len(all) > n {
		return all[:n]
	}
	return all
}

func (pt *ProjectTracker) Remove(path string) bool {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	absPath, _ := filepath.Abs(path)
	if _, ok := pt.projects[absPath]; ok {
		delete(pt.projects, absPath)
		_ = pt.saveToDisk()
		return true
	}
	return false
}

func (pt *ProjectTracker) Count() int {
	pt.mu.RLock()
	defer pt.mu.RUnlock()
	return len(pt.projects)
}

func (pt *ProjectTracker) detectProjectType(path string, project *Project) {
	checks := []struct{ file, language, framework string }{
		{"go.mod", "go", ""},
		{"package.json", "javascript", "node"},
		{"tsconfig.json", "typescript", ""},
		{"Cargo.toml", "rust", ""},
		{"pyproject.toml", "python", ""},
		{"requirements.txt", "python", ""},
		{"Gemfile", "ruby", ""},
		{"pom.xml", "java", "maven"},
		{"build.gradle", "java", "gradle"},
	}
	for _, c := range checks {
		if _, err := os.Stat(filepath.Join(path, c.file)); err == nil {
			if project.Language == "" {
				project.Language = c.language
			}
			if project.Framework == "" {
				project.Framework = c.framework
			}
		}
	}
	if project.Language == "javascript" {
		if _, err := os.Stat(filepath.Join(path, "tsconfig.json")); err == nil {
			project.Language = "typescript"
		}
	}
}

func (pt *ProjectTracker) saveToDisk() error {
	if pt.storePath == "" {
		return nil
	}
	var projects []*Project
	for _, p := range pt.projects {
		projects = append(projects, p)
	}
	data, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(pt.storePath), 0755); err != nil {
		return err
	}
	return os.WriteFile(pt.storePath, data, 0644)
}

func (pt *ProjectTracker) loadFromDisk() error {
	data, err := os.ReadFile(pt.storePath)
	if err != nil {
		return err
	}
	var projects []*Project
	if err := json.Unmarshal(data, &projects); err != nil {
		return err
	}
	for _, p := range projects {
		pt.projects[p.Path] = p
	}
	return nil
}

func generateProjectID(path string) string {
	h := 0
	for _, c := range path {
		h = h*31 + int(c)
	}
	return fmt.Sprintf("proj_%08x", uint32(h))
}

var _ = strings.TrimSpace
