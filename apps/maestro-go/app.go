package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"maestro-go/internal/agents"
)

// App struct
type App struct {
	ctx        context.Context
	processes  map[string]*exec.Cmd
	processMu  sync.Mutex
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		processes: make(map[string]*exec.Cmd),
	}
}

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// GetAgents returns the list of detected agents
func (a *App) ListAgents() []agents.Agent {
	return agents.Detect()
}

// ExecuteCommand runs a command and streams output to frontend
func (a *App) ExecuteCommand(id string, command string, args []string, cwd string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = cwd
	
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	a.processMu.Lock()
	a.processes[id] = cmd
	a.processMu.Unlock()

	go a.streamOutput(id, stdout, "stdout")
	go a.streamOutput(id, stderr, "stderr")

	go func() {
		err := cmd.Wait()
		runtime.EventsEmit(a.ctx, "process:exit", map[string]interface{}{
			"id":       id,
			"exitCode": cmd.ProcessState.ExitCode(),
			"error":    err != nil,
		})
		a.processMu.Lock()
		delete(a.processes, id)
		a.processMu.Unlock()
	}()

	return nil
}

func (a *App) streamOutput(id string, rc io.ReadCloser, source string) {
	reader := bufio.NewReader(rc)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		runtime.EventsEmit(a.ctx, "process:data", map[string]interface{}{
			"id":     id,
			"data":   line,
			"source": source,
		})
	}
}

// ResizeProcess changes terminal dimensions
func (a *App) ResizeProcess(id string, cols, rows int) error {
	// PTY resize not implemented in standard os/exec
	// requires a PTY library
	return nil
}

// GetConfig returns the current settings
func (a *App) GetConfig() (map[string]interface{}, error) {
	return map[string]interface{}{
		"activeThemeId": "dracula",
		"llmProvider":   "openrouter",
		"modelSlug":     "anthropic/claude-3.5-sonnet",
	}, nil
}

// ReadDir returns the contents of a directory
func (a *App) ReadDir(path string) ([]map[string]interface{}, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0, len(entries))
	for _, entry := range entries {
		info, _ := entry.Info()
		result = append(result, map[string]interface{}{
			"name":  entry.Name(),
			"isDir": entry.IsDir(),
			"size":  info.Size(),
			"mtime": info.ModTime().UnixMilli(),
		})
	}
	return result, nil
}

// ReadFile returns the content of a file
func (a *App) ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
