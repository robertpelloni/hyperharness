package skills

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// Skill represents a packaged capability that teaches specific patterns and best practices.
// Modeled after pi's skill system and OpenCode's skill tool.
type Skill struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Version     string                 `json:"version"`
	Category    string                 `json:"category"`
	TriggerWords []string              `json:"trigger_words"`
	Content     string                 `json:"content"`
	FilePath    string                 `json:"file_path,omitempty"`
	Tags        []string               `json:"tags"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// SkillResult is the result of executing a skill.
type SkillResult struct {
	Name    string `json:"name"`
	Output  string `json:"output"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// Manager manages skill discovery, loading, and execution.
type Manager struct {
	skills    map[string]*Skill
	loaded    bool
	skillDirs []string
	mu        sync.RWMutex
}

// NewManager creates a new skill manager.
func NewManager() *Manager {
	return &Manager{
		skills: make(map[string]*Skill),
		skillDirs: []string{
			".agents/skills",
			".skills",
			"skills",
		},
	}
}

// Discover scans configured directories for SKILL.md files and loads them.
func (m *Manager) Discover(baseDir string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, skillDir := range m.skillDirs {
		fullPath := filepath.Join(baseDir, skillDir)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			continue
		}

		// Walk the directory looking for SKILL.md files
		err := filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				return nil
			}
			if strings.EqualFold(info.Name(), "SKILL.md") {
				if err := m.loadSkillFile(path); err != nil {
					// Non-fatal, continue scanning
				}
			}
			return nil
		})
		if err != nil {
			continue
		}
	}

	m.loaded = true
	return nil
}

// loadSkillFile loads a single SKILL.md file.
func (m *Manager) loadSkillFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	content := string(data)
	skillName := filepath.Base(filepath.Dir(path))
	if skillName == "." || skillName == "" {
		skillName = filepath.Base(path)
	}

	// Parse frontmatter-like description from first line
	description := ""
	lines := strings.Split(content, "\n")
	if len(lines) > 0 {
		firstLine := lines[0]
		if strings.HasPrefix(firstLine, "# ") {
			description = strings.TrimPrefix(firstLine, "# ")
		}
	}

	skill := &Skill{
		Name:        skillName,
		Description: description,
		Content:     content,
		FilePath:    path,
		Category:    filepath.Base(filepath.Dir(filepath.Dir(path))),
	}

	// Check for skill.json metadata
	metadataPath := filepath.Join(filepath.Dir(path), "skill.json")
	if meta, err := os.ReadFile(metadataPath); err == nil {
		var metaMap map[string]interface{}
		if json.Unmarshal(meta, &metaMap) == nil {
			skill.Metadata = metaMap
			if v, ok := metaMap["name"].(string); ok && v != "" {
				skill.Name = v
			}
			if v, ok := metaMap["description"].(string); ok && v != "" {
				skill.Description = v
			}
			if v, ok := metaMap["version"].(string); ok {
				skill.Version = v
			}
			if v, ok := metaMap["category"].(string); ok {
				skill.Category = v
			}
			if v, ok := metaMap["trigger_words"].([]interface{}); ok {
				for _, tw := range v {
					if s, ok := tw.(string); ok {
						skill.TriggerWords = append(skill.TriggerWords, s)
					}
				}
			}
			if v, ok := metaMap["tags"].([]interface{}); ok {
				for _, t := range v {
					if s, ok := t.(string); ok {
						skill.Tags = append(skill.Tags, s)
					}
				}
			}
		}
	}

	m.skills[skill.Name] = skill
	return nil
}

// Get retrieves a skill by name.
func (m *Manager) Get(name string) (*Skill, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	skill, ok := m.skills[name]
	return skill, ok
}

// List returns all available skills.
func (m *Manager) List() []*Skill {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var skills []*Skill
	for _, skill := range m.skills {
		skills = append(skills, skill)
	}
	sort.Slice(skills, func(i, j int) bool {
		return skills[i].Name < skills[j].Name
	})
	return skills
}

// ListByCategory returns skills grouped by category.
func (m *Manager) ListByCategory() map[string][]*Skill {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string][]*Skill)
	for _, skill := range m.skills {
		result[skill.Category] = append(result[skill.Category], skill)
	}
	return result
}

// Execute runs a skill and returns its result.
func (m *Manager) Execute(name string, params map[string]interface{}) (*SkillResult, error) {
	m.mu.RLock()
	skill, ok := m.skills[name]
	m.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("skill not found: %s", name)
	}

	// Skills are primarily knowledge injection - they return their content
	// with any parameter substitutions applied
	output := skill.Content

	// Simple parameter substitution ({{param_name}} -> value)
	for key, value := range params {
		placeholder := fmt.Sprintf("{{%s}}", key)
		if strVal, ok := value.(string); ok {
			output = strings.ReplaceAll(output, placeholder, strVal)
		}
	}

	return &SkillResult{
		Name:    skill.Name,
		Output:  output,
		Success: true,
	}, nil
}

// MatchByTrigger finds skills that match the given text.
func (m *Manager) MatchByTrigger(text string) []*Skill {
	m.mu.RLock()
	defer m.mu.RUnlock()

	textLower := strings.ToLower(text)
	var matches []*Skill

	for _, skill := range m.skills {
		for _, trigger := range skill.TriggerWords {
			if strings.Contains(textLower, strings.ToLower(trigger)) {
				matches = append(matches, skill)
				break
			}
		}
	}

	return matches
}

// Register adds a skill programmatically.
func (m *Manager) Register(skill *Skill) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.skills[skill.Name] = skill
}

// Builtins returns the built-in skills.
func Builtins() []*Skill {
	return []*Skill{
		{
			Name:         "memory",
			Description:  "Use AI DevKit's memory service for knowledge storage and retrieval.",
			Version:      "1.0.0",
			Category:     "core",
			TriggerWords: []string{"memory", "remember", "recall", "knowledge"},
			Tags:         []string{"memory", "storage", "knowledge"},
			Content: `# Memory Skill

## Usage

### Store Knowledge
Use 'memory_store' to save architectural decisions, patterns, and conventions.

### Search Knowledge
Use 'memory_search' to find relevant knowledge before starting tasks.

### Best Practices
- Search memory before asking questions
- Store decisions with clear titles and tags
- Use appropriate scope (global, project, session)`,
		},
		{
			Name:         "context",
			Description:  "Manage conversation context for optimal AI performance.",
			Version:      "1.0.0",
			Category:     "core",
			TriggerWords: []string{"context", "compact", "inject", "summarize"},
			Tags:         []string{"context", "management"},
			Content: `# Context Management Skill

## Usage
- Use 'context_manager' with action "compact" to reduce context size
- Use action "inject" to add system-level context
- Use action "status" to check context utilization

## Best Practices
- Monitor context utilization percentage
- Compact when utilization exceeds 80%
- Inject relevant project context at session start`,
		},
		{
			Name:         "agent-browser",
			Description:  "Browser automation CLI for AI agents.",
			Version:      "1.0.0",
			Category:     "automation",
			TriggerWords: []string{"browser", "web", "scrape", "screenshot", "navigate"},
			Tags:         []string{"browser", "automation", "web"},
			Content: `# Agent Browser Skill

Automates web browser interactions including:
- Page navigation
- Form filling
- Button clicking
- Screenshots
- Data extraction
- Web app testing`,
		},
		{
			Name:         "simplify",
			Description:  "Simplify and refine recently modified code.",
			Version:      "1.0.0",
			Category:     "code-quality",
			TriggerWords: []string{"simplify", "refactor", "clean", "polish"},
			Tags:         []string{"refactor", "quality"},
			Content: `# Simplify Skill

Use after writing code to improve readability without changing functionality.
Focus on:
- Removing unnecessary complexity
- Improving naming
- Reducing duplication
- Simplifying control flow`,
		},
	}
}
