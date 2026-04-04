package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

type localAutoDevLoopConfig struct {
	MaxAttempts int    `json:"maxAttempts"`
	Type        string `json:"type"`
	Target      string `json:"target,omitempty"`
	Command     string `json:"command,omitempty"`
}

type localAutoDevLoop struct {
	ID             string                 `json:"id"`
	Config         localAutoDevLoopConfig `json:"config"`
	Status         string                 `json:"status"`
	CurrentAttempt int                    `json:"currentAttempt"`
	StartTime      int64                  `json:"startTime"`
	LastOutput     string                 `json:"lastOutput"`
}

type persistedAutoDevState struct {
	NextID int64              `json:"nextId"`
	Loops  []localAutoDevLoop `json:"loops"`
}

type localAutoDevManager struct {
	mu          sync.RWMutex
	persistPath string
	workspace   string
	nextID      int64
	loops       map[string]*localAutoDevLoop
}

func newLocalAutoDevManager(workspaceRoot string, persistPath string) *localAutoDevManager {
	m := &localAutoDevManager{
		persistPath: persistPath,
		workspace:   workspaceRoot,
		nextID:      1,
		loops:       map[string]*localAutoDevLoop{},
	}
	m.load()
	for _, loop := range m.loops {
		if loop != nil && loop.Status == "running" {
			loop.Status = "failed"
			loop.LastOutput = strings.TrimSpace(loop.LastOutput + "\n[Native Go AutoDev fallback note] Loop resumed after process restart as failed because in-flight execution state is not yet recoverable.")
		}
	}
	m.saveLocked()
	return m
}

func (m *localAutoDevManager) startLoop(config localAutoDevLoopConfig) localAutoDevLoop {
	m.mu.Lock()
	defer m.mu.Unlock()
	id := fmt.Sprintf("loop-%d", m.nextID)
	m.nextID++
	if config.MaxAttempts <= 0 {
		config.MaxAttempts = 3
	}
	loop := &localAutoDevLoop{
		ID:             id,
		Config:         config,
		Status:         "running",
		CurrentAttempt: 0,
		StartTime:      time.Now().UTC().UnixMilli(),
		LastOutput:     "",
	}
	m.loops[id] = loop
	m.saveLocked()
	go m.runLoop(id)
	return *loop
}

func (m *localAutoDevManager) cancelLoop(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	loop := m.loops[id]
	if loop == nil || loop.Status != "running" {
		return false
	}
	loop.Status = "cancelled"
	if strings.TrimSpace(loop.LastOutput) == "" {
		loop.LastOutput = "Cancellation requested. Native Go AutoDev fallback does not kill an in-flight process; it only prevents further retry attempts."
	} else {
		loop.LastOutput += "\nCancellation requested. Native Go AutoDev fallback does not kill an in-flight process; it only prevents further retry attempts."
	}
	m.saveLocked()
	return true
}

func (m *localAutoDevManager) getLoops() []localAutoDevLoop {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]localAutoDevLoop, 0, len(m.loops))
	for _, loop := range m.loops {
		if loop == nil {
			continue
		}
		result = append(result, *loop)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].StartTime > result[j].StartTime
	})
	return result
}

func (m *localAutoDevManager) getLoop(id string) (*localAutoDevLoop, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	loop := m.loops[id]
	if loop == nil {
		return nil, false
	}
	copyLoop := *loop
	return &copyLoop, true
}

func (m *localAutoDevManager) clearCompleted() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	count := 0
	for id, loop := range m.loops {
		if loop == nil {
			delete(m.loops, id)
			continue
		}
		if loop.Status != "running" {
			delete(m.loops, id)
			count++
		}
	}
	m.saveLocked()
	return count
}

func (m *localAutoDevManager) runLoop(id string) {
	for {
		m.mu.Lock()
		loop := m.loops[id]
		if loop == nil {
			m.mu.Unlock()
			return
		}
		if loop.Status != "running" {
			m.saveLocked()
			m.mu.Unlock()
			return
		}
		if loop.CurrentAttempt >= loop.Config.MaxAttempts {
			loop.Status = "failed"
			if strings.TrimSpace(loop.LastOutput) == "" {
				loop.LastOutput = "Native Go AutoDev fallback exhausted all retry attempts."
			}
			m.saveLocked()
			m.mu.Unlock()
			return
		}
		loop.CurrentAttempt++
		attempt := loop.CurrentAttempt
		config := loop.Config
		m.saveLocked()
		m.mu.Unlock()

		output, success := m.executeAttempt(config)

		m.mu.Lock()
		loop = m.loops[id]
		if loop == nil {
			m.mu.Unlock()
			return
		}
		loop.LastOutput = output
		if loop.Status != "running" {
			m.saveLocked()
			m.mu.Unlock()
			return
		}
		if success {
			loop.Status = "success"
			m.saveLocked()
			m.mu.Unlock()
			return
		}
		if attempt >= config.MaxAttempts {
			loop.Status = "failed"
			m.saveLocked()
			m.mu.Unlock()
			return
		}
		m.saveLocked()
		m.mu.Unlock()
		time.Sleep(minDuration(250*time.Millisecond*time.Duration(1<<max(0, attempt-1)), 2*time.Second))
	}
}

func (m *localAutoDevManager) executeAttempt(config localAutoDevLoopConfig) (string, bool) {
	command := strings.TrimSpace(config.Command)
	if command == "" {
		command = m.defaultCommand(config)
	}
	if strings.TrimSpace(command) == "" {
		return "Native Go AutoDev fallback could not determine a command to run.", false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd := shellCommand(ctx, command)
	cmd.Dir = m.workspace
	output, err := cmd.CombinedOutput()
	text := strings.TrimSpace(string(output))
	if text == "" {
		text = command
	}
	if err == nil {
		return text, true
	}
	if text == command {
		text = strings.TrimSpace(err.Error())
	}
	return text, false
}

func (m *localAutoDevManager) defaultCommand(config localAutoDevLoopConfig) string {
	switch config.Type {
	case "test":
		if strings.TrimSpace(config.Target) != "" {
			return "npx vitest run " + config.Target
		}
		return "npm test"
	case "lint":
		if strings.TrimSpace(config.Target) != "" {
			return "npx eslint --fix " + config.Target
		}
		return "npm run lint -- --fix"
	case "build":
		return "npm run build"
	default:
		return ""
	}
}

func shellCommand(ctx context.Context, command string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.CommandContext(ctx, "cmd", "/C", command)
	}
	return exec.CommandContext(ctx, "sh", "-c", command)
}

func (m *localAutoDevManager) load() {
	if strings.TrimSpace(m.persistPath) == "" {
		return
	}
	data, err := os.ReadFile(m.persistPath)
	if err != nil {
		return
	}
	var state persistedAutoDevState
	if err := json.Unmarshal(data, &state); err != nil {
		return
	}
	if state.NextID > 0 {
		m.nextID = state.NextID
	}
	for _, loop := range state.Loops {
		copyLoop := loop
		m.loops[loop.ID] = &copyLoop
	}
}

func (m *localAutoDevManager) saveLocked() {
	if strings.TrimSpace(m.persistPath) == "" {
		return
	}
	state := persistedAutoDevState{NextID: m.nextID, Loops: make([]localAutoDevLoop, 0, len(m.loops))}
	for _, loop := range m.loops {
		if loop == nil {
			continue
		}
		state.Loops = append(state.Loops, *loop)
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return
	}
	_ = os.MkdirAll(filepath.Dir(m.persistPath), 0o755)
	_ = os.WriteFile(m.persistPath, data, 0o644)
}

func minDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
