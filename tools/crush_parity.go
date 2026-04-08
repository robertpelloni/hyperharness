package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	foundationpi "github.com/robertpelloni/hyperharness/foundation/pi"
)

// crush_parity.go provides comprehensive tool surfaces ported from the Crush CLI agent
// (github.com/charmbracelet/crush). These tools match the exact names, parameters,
// and result formats that models trained on Crush expect.
//
// Key features ported:
// - multiedit: Multi-operation file editing with partial failure handling
// - view: File viewing with offset/limit, line numbers, and base64 image support
// - diagnostics: Code diagnostics gathering
// - todos: Task tracking
// - web_fetch: HTTP content fetching with truncation
// - web_search: Web search (placeholder routing to search providers)
// - glob: File globbing with path filtering
// - job management: Background command execution, output retrieval, job killing
// - references: Code reference finding
// - lsp_restart: Language server protocol restart capability
// - safe operations: Safety-checked file operations

// ==========================
// Background Job Management
// ==========================

// Job represents a running background command.
type Job struct {
	ID          string
	Command     string
	Description string
	WorkingDir  string
	StartTime   time.Time
	Cmd         *exec.Cmd
	Output      *strings.Builder
	Done        bool
	ExitCode    int
	mu          sync.Mutex
	cancel      context.CancelFunc
}

// JobManager manages background jobs across the tool registry.
type JobManager struct {
	jobs map[string]*Job
	mu   sync.RWMutex
}

// GlobalJobManager is the singleton job manager for background processes.
var GlobalJobManager = NewJobManager()

func NewJobManager() *JobManager {
	return &JobManager{jobs: make(map[string]*Job)}
}

func (jm *JobManager) AddJob(job *Job) {
	jm.mu.Lock()
	defer jm.mu.Unlock()
	jm.jobs[job.ID] = job
}

func (jm *JobManager) GetJob(id string) (*Job, bool) {
	jm.mu.RLock()
	defer jm.mu.RUnlock()
	job, ok := jm.jobs[id]
	return job, ok
}

func (jm *JobManager) RemoveJob(id string) {
	jm.mu.Lock()
	defer jm.mu.Unlock()
	delete(jm.jobs, id)
}

func (jm *JobManager) ListJobs() []*Job {
	jm.mu.RLock()
	defer jm.mu.RUnlock()
	var jobs []*Job
	for _, job := range jm.jobs {
		jobs = append(jobs, job)
	}
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].StartTime.Before(jobs[j].StartTime)
	})
	return jobs
}

// ==========================
// Crush-Parity Tool Registration
// ==========================

// registerCrushParityTools adds all Crush-compatible tool surfaces.
func (r *Registry) registerCrushParityTools() {
	r.registerMultiEditTool()
	r.registerViewTool()
	r.registerWriteTool()
	r.registerGlobTool()
	r.registerBashTool()
	r.registerWebFetchTool()
	r.registerWebSearchTool()
	r.registerDiagnosticsTool()
	r.registerTodosTool()
	r.registerReferencesTool()
	r.registerLSPRestartTool()
	r.registerJobManagementTools()
	r.registerSafeTool()
	r.registerDownloadTool()
	r.registerSourcegraphTool()
	r.registerSearchTool()
}

// ==========================
// multiedit tool
// ==========================

func (r *Registry) registerMultiEditTool() {
	r.Tools = append(r.Tools, Tool{
		Name:        "multiedit",
		Description: "Apply multiple edit operations to a single file sequentially. Each edit replaces old_string with new_string. If the first edit has empty old_string, creates a new file. Partial failures are reported individually.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["file_path", "edits"],
			"properties": {
				"file_path": {
					"type": "string",
					"description": "The absolute path to the file to modify"
				},
				"edits": {
					"type": "array",
					"items": {
						"type": "object",
						"required": ["old_string", "new_string"],
						"properties": {
							"old_string": {
								"type": "string",
								"description": "The text to replace. Empty for first edit means create new file."
							},
							"new_string": {
								"type": "string",
								"description": "The text to replace it with"
							},
							"replace_all": {
								"type": "boolean",
								"description": "Replace all occurrences of old_string (default false)"
							}
						}
					},
					"description": "Array of edit operations to apply sequentially"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			filePath, _ := args["file_path"].(string)
			if filePath == "" {
				return "", fmt.Errorf("file_path is required")
			}

			editsRaw, ok := args["edits"].([]interface{})
			if !ok || len(editsRaw) == 0 {
				return "", fmt.Errorf("edits must be a non-empty array")
			}

			// Parse all edits
			type editOp struct {
				OldString  string
				NewString  string
				ReplaceAll bool
			}
			var edits []editOp
			for i, eRaw := range editsRaw {
				e, ok := eRaw.(map[string]interface{})
				if !ok {
					return "", fmt.Errorf("edit %d is invalid", i)
				}
				oldStr, _ := e["old_string"].(string)
				newStr, _ := e["new_string"].(string)
				replaceAll := false
				if ra, ok := e["replace_all"].(bool); ok {
					replaceAll = ra
				}
				edits = append(edits, editOp{OldString: oldStr, NewString: newStr, ReplaceAll: replaceAll})
			}

			// Handle file creation case
			if len(edits) > 0 && edits[0].OldString == "" {
				content := edits[0].NewString
				// Apply remaining edits to the content
				var failures []string
				applied := 1
				for i := 1; i < len(edits); i++ {
					newContent, err := applyEditToContent(content, edits[i].OldString, edits[i].NewString, edits[i].ReplaceAll)
					if err != nil {
						failures = append(failures, fmt.Sprintf("edit %d: %s", i+1, err.Error()))
						continue
					}
					content = newContent
					applied++
				}

				dir := filepath.Dir(filePath)
				if err := os.MkdirAll(dir, 0o755); err != nil {
					return "", fmt.Errorf("failed to create parent directories: %w", err)
				}

				if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
					return "", fmt.Errorf("failed to write file: %w", err)
				}

				if len(failures) > 0 {
					return fmt.Sprintf("File created with %d of %d edits applied to %s.\n%d failed:\n%s",
						applied, len(edits), filePath, len(failures), strings.Join(failures, "\n")), nil
				}
				return fmt.Sprintf("File created with %d edits applied: %s", len(edits), filePath), nil
			}

			// Read existing file
			data, err := os.ReadFile(filePath)
			if err != nil {
				return "", fmt.Errorf("failed to read file: %w", err)
			}
			content := normalizeLineEndings(string(data))
			oldContent := content

			// Apply edits sequentially
			var failures []string
			applied := 0
			for i, edit := range edits {
				newContent, err := applyEditToContent(content, edit.OldString, edit.NewString, edit.ReplaceAll)
				if err != nil {
					failures = append(failures, fmt.Sprintf("edit %d: %s", i+1, err.Error()))
					continue
				}
				content = newContent
				applied++
			}

			if content == oldContent {
				if len(failures) > 0 {
					return fmt.Sprintf("No changes made - all %d edit(s) failed:\n%s", len(failures), strings.Join(failures, "\n")), nil
				}
				return "No changes made - all edits resulted in identical content", nil
			}

			if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
				return "", fmt.Errorf("failed to write file: %w", err)
			}

			if len(failures) > 0 {
				return fmt.Sprintf("Applied %d of %d edits to file: %s\n%d failed:\n%s",
					applied, len(edits), filePath, len(failures), strings.Join(failures, "\n")), nil
			}
			return fmt.Sprintf("Applied %d edits to file: %s", len(edits), filePath), nil
		},
	})
}

func applyEditToContent(content, oldString, newString string, replaceAll bool) (string, error) {
	if oldString == "" {
		return "", fmt.Errorf("old_string cannot be empty for content replacement")
	}

	if replaceAll {
		if !strings.Contains(content, oldString) {
			return "", fmt.Errorf("old_string not found in content")
		}
		return strings.ReplaceAll(content, oldString, newString), nil
	}

	index := strings.Index(content, oldString)
	if index == -1 {
		return "", fmt.Errorf("old_string not found in content. Make sure it matches exactly, including whitespace and line breaks")
	}

	lastIndex := strings.LastIndex(content, oldString)
	if index != lastIndex {
		return "", fmt.Errorf("old_string appears multiple times in the content. Please provide more context to ensure a unique match, or set replace_all to true")
	}

	return content[:index] + newString + content[index+len(oldString):], nil
}

// ==========================
// view tool (with offset/limit)
// ==========================

func (r *Registry) registerViewTool() {
	r.Tools = append(r.Tools, Tool{
		Name:        "view",
		Description: "Read the contents of a file. Supports text files and images (jpg, png, gif, webp). Images are sent as attachments. For text files, output is truncated to 2000 lines or 50KB (whichever is hit first). Use offset/limit for large files.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["file_path"],
			"properties": {
				"file_path": {
					"type": "string",
					"description": "Path to the file to read (relative or absolute)"
				},
				"offset": {
					"type": "number",
					"description": "Line number to start reading from (1-indexed)"
				},
				"limit": {
					"type": "number",
					"description": "Maximum number of lines to read"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			path, _ := args["file_path"].(string)
			if path == "" {
				return "", fmt.Errorf("file_path is required")
			}

			offset := toInt(args["offset"], 1)
			limit := toInt(args["limit"], 2000)

			if offset < 1 {
				offset = 1
			}
			if limit < 1 {
				limit = 2000
			}

			piArgs := map[string]interface{}{
				"path":   path,
				"offset": offset,
				"limit":  limit,
			}
			raw, _ := json.Marshal(piArgs)
			runtime := foundationpi.NewRuntime(".", nil)
			result, err := runtime.ExecuteTool(context.Background(), "", "read", raw, nil)
			if err != nil {
				return "", err
			}
			return formatFoundationToolResult(result), nil
		},
	})
}

// ==========================
// write tool (enhanced with diff)
// ==========================

func (r *Registry) registerWriteTool() {
	r.Tools = append(r.Tools, Tool{
		Name:        "write",
		Description: "Write content to a file. Creates the file if it doesn't exist, overwrites if it does. Automatically creates parent directories.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["file_path", "content"],
			"properties": {
				"file_path": {
					"type": "string",
					"description": "The path to the file to write"
				},
				"content": {
					"type": "string",
					"description": "The content to write to the file"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			path, _ := args["file_path"].(string)
			content, _ := args["content"].(string)
			if path == "" {
				return "", fmt.Errorf("file_path is required")
			}
			if content == "" {
				return "", fmt.Errorf("content is required")
			}

			dir := filepath.Dir(path)
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return "", fmt.Errorf("failed to create parent directories: %w", err)
			}

			if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
				return "", fmt.Errorf("failed to write file: %w", err)
			}

			return fmt.Sprintf("File written successfully: %s (%d bytes)", path, len(content)), nil
		},
	})
}

// ==========================
// glob tool
// ==========================

func (r *Registry) registerGlobTool() {
	r.Tools = append(r.Tools, Tool{
		Name:        "glob",
		Description: "Fast file pattern matching tool that supports glob patterns like \"**/*.js\" or \"src/**/*.ts\". Returns matching file paths sorted by modification time. Respects .gitignore by default.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["pattern"],
			"properties": {
				"pattern": {
					"type": "string",
					"description": "The glob pattern to match files against (e.g. \"**/*.go\", \"src/**/*.ts\")"
				},
				"path": {
					"type": "string",
					"description": "The directory to search in. Defaults to current working directory."
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			pattern, _ := args["pattern"].(string)
			searchPath, _ := args["path"].(string)
			if pattern == "" {
				return "", fmt.Errorf("pattern is required")
			}
			if searchPath == "" {
				searchPath = "."
			}

			piArgs := map[string]interface{}{"pattern": pattern, "path": searchPath}
			raw, _ := json.Marshal(piArgs)
			runtime := foundationpi.NewRuntime(".", nil)
			result, err := runtime.ExecuteTool(context.Background(), "", "find", raw, nil)
			if err != nil {
				return "", err
			}
			return formatFoundationToolResult(result), nil
		},
	})
}

// ==========================
// bash tool (enhanced with background jobs)
// ==========================

func (r *Registry) registerBashTool() {
	r.Tools = append(r.Tools, Tool{
		Name:        "bash",
		Description: "Execute a bash command and return its output. Supports background execution for long-running commands.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["command"],
			"properties": {
				"command": {
					"type": "string",
					"description": "The command to execute"
				},
				"description": {
					"type": "string",
					"description": "A brief description of what the command does"
				},
				"working_dir": {
					"type": "string",
					"description": "The working directory to execute the command in"
				},
				"run_in_background": {
					"type": "boolean",
					"description": "Set to true to run this command in the background. Use job_output to read the output later."
				},
				"timeout": {
					"type": "number",
					"description": "Optional timeout in seconds"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			command, _ := args["command"].(string)
			if command == "" {
				return "", fmt.Errorf("command is required")
			}

			workingDir, _ := args["working_dir"].(string)
			runInBackground, _ := args["run_in_background"].(bool)
			description, _ := args["description"].(string)
			timeoutSeconds := toInt(args["timeout"], 0)

			if workingDir == "" {
				workingDir = "."
			}

			if runInBackground {
				return executeBackgroundCommand(command, description, workingDir)
			}

			return executeForegroundCommand(command, workingDir, timeoutSeconds)
		},
	})
}

func executeForegroundCommand(command, workingDir string, timeoutSeconds int) (string, error) {
	var cmd *exec.Cmd
	ctx := context.Background()
	if timeoutSeconds > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
		defer cancel()
	}

	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "cmd", "/C", command)
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-lc", command)
	}
	cmd.Dir = workingDir

	output, err := cmd.CombinedOutput()
	result := string(output)
	if result == "" && err != nil {
		return "", err
	}
	if len(result) > 30000 {
		result = result[:30000] + "\n... (output truncated at 30000 bytes)"
	}
	if err != nil {
		return result, nil // Return output even on error
	}
	return result, nil
}

func executeBackgroundCommand(command, description, workingDir string) (string, error) {
	jobID := fmt.Sprintf("job_%d", time.Now().UnixNano())

	var cmd *exec.Cmd
	ctx, cancel := context.WithCancel(context.Background())

	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "cmd", "/C", command)
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-lc", command)
	}
	cmd.Dir = workingDir

	job := &Job{
		ID:          jobID,
		Command:     command,
		Description: description,
		WorkingDir:  workingDir,
		StartTime:   time.Now(),
		Cmd:         cmd,
		Output:      &strings.Builder{},
		cancel:      cancel,
	}

	GlobalJobManager.AddJob(job)

	go func() {
		output, err := cmd.CombinedOutput()
		job.mu.Lock()
		defer job.mu.Unlock()
		job.Output.Write(output)
		job.Done = true
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				job.ExitCode = exitErr.ExitCode()
			} else {
				job.ExitCode = -1
			}
		}
	}()

	return fmt.Sprintf("Background job started: %s\nCommand: %s\nDescription: %s",
		jobID, command, description), nil
}

// ==========================
// web_fetch tool
// ==========================

func (r *Registry) registerWebFetchTool() {
	r.Tools = append(r.Tools, Tool{
		Name:        "web_fetch",
		Description: "Fetch content from a URL. Returns the response body as text, truncated to a reasonable size. Useful for reading web pages, API responses, or documentation.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["url"],
			"properties": {
				"url": {
					"type": "string",
					"description": "The URL to fetch"
				},
				"max_length": {
					"type": "number",
					"description": "Maximum response length in bytes (default 50000)"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			url, _ := args["url"].(string)
			if url == "" {
				return "", fmt.Errorf("url is required")
			}

			maxLength := toInt(args["max_length"], 50000)

			client := &http.Client{Timeout: 30 * time.Second}
			resp, err := client.Get(url)
			if err != nil {
				return "", fmt.Errorf("failed to fetch URL: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
			}

			limitedReader := io.LimitReader(resp.Body, int64(maxLength+1))
			data, err := io.ReadAll(limitedReader)
			if err != nil {
				return "", fmt.Errorf("failed to read response: %w", err)
			}

			content := string(data)
			if len(data) > maxLength {
				content = content[:maxLength] + fmt.Sprintf("\n\n... (response truncated at %d bytes)", maxLength)
			}

			return content, nil
		},
	})
}

// ==========================
// web_search tool
// ==========================

func (r *Registry) registerWebSearchTool() {
	r.Tools = append(r.Tools, Tool{
		Name:        "web_search",
		Description: "Search the web for information. Returns search results with titles, URLs, and snippets.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["query"],
			"properties": {
				"query": {
					"type": "string",
					"description": "The search query"
				},
				"num_results": {
					"type": "number",
					"description": "Number of results to return (default 10)"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			query, _ := args["query"].(string)
			if query == "" {
				return "", fmt.Errorf("query is required")
			}

			// Web search requires external provider integration.
			// Return a structured response that the agent loop can route to
			// the configured search provider (e.g., via MCP).
			return fmt.Sprintf("Web search query: %q\n(Note: Configure a search provider via MCP or environment to enable live web search results)", query), nil
		},
	})
}

// ==========================
// diagnostics tool
// ==========================

func (r *Registry) registerDiagnosticsTool() {
	r.Tools = append(r.Tools, Tool{
		Name:        "diagnostics",
		Description: "Get code diagnostics (errors, warnings) for a file or project. Returns compiler/linter output.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"file_path": {
					"type": "string",
					"description": "The file to get diagnostics for. If empty, checks the entire project."
				},
				"language": {
					"type": "string",
					"description": "The language of the file (go, typescript, python, etc.)"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			filePath, _ := args["file_path"].(string)
			language, _ := args["language"].(string)

			if language == "" {
				ext := strings.ToLower(filepath.Ext(filePath))
				switch ext {
				case ".go":
					language = "go"
				case ".ts", ".tsx":
					language = "typescript"
				case ".js", ".jsx":
					language = "javascript"
				case ".py":
					language = "python"
				case ".rs":
					language = "rust"
				default:
					language = "unknown"
				}
			}

			switch language {
			case "go":
				return runGoDiagnostics(filePath)
			default:
				return fmt.Sprintf("Diagnostics for %s (language: %s): No built-in diagnostics available for this language. Consider using an LSP server via MCP.", filePath, language), nil
			}
		},
	})
}

func runGoDiagnostics(filePath string) (string, error) {
	cmd := exec.Command("go", "vet", "./...")
	if filePath != "" {
		dir := filepath.Dir(filePath)
		cmd.Dir = dir
	} else {
		cmd.Dir = "."
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), nil
	}
	if len(output) == 0 {
		return "No diagnostics issues found.", nil
	}
	return string(output), nil
}

// ==========================
// todos tool
// ==========================

// TodoItem represents a tracked task.
type TodoItem struct {
	ID          string `json:"id"`
	Content     string `json:"content"`
	Status      string `json:"status"` // "pending", "in_progress", "completed"
	Priority    string `json:"priority,omitempty"`
	File        string `json:"file,omitempty"`
	Line        int    `json:"line,omitempty"`
}

// TodoStore manages todos in memory.
var TodoStore = struct {
	sync.RWMutex
	Items []TodoItem
}{}

func (r *Registry) registerTodosTool() {
	r.Tools = append(r.Tools, Tool{
		Name:        "todos",
		Description: "Manage a todo list for tracking tasks during development. Actions: list, add, update, remove, clear.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["action"],
			"properties": {
				"action": {
					"type": "string",
					"description": "The action to perform: list, add, update, remove, clear",
					"enum": ["list", "add", "update", "remove", "clear"]
				},
				"id": {
					"type": "string",
					"description": "Todo item ID (for update/remove)"
				},
				"content": {
					"type": "string",
					"description": "Todo item content (for add/update)"
				},
				"status": {
					"type": "string",
					"description": "Todo item status: pending, in_progress, completed",
					"enum": ["pending", "in_progress", "completed"]
				},
				"priority": {
					"type": "string",
					"description": "Priority level: low, medium, high"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			action, _ := args["action"].(string)
			switch action {
			case "list":
				return todoList()
			case "add":
				return todoAdd(args)
			case "update":
				return todoUpdate(args)
			case "remove":
				return todoRemove(args)
			case "clear":
				return todoClear()
			default:
				return "", fmt.Errorf("unknown action: %s (supported: list, add, update, remove, clear)", action)
			}
		},
	})
}

func todoList() (string, error) {
	TodoStore.RLock()
	defer TodoStore.RUnlock()
	if len(TodoStore.Items) == 0 {
		return "No todos.", nil
	}
	var lines []string
	for _, item := range TodoStore.Items {
		status := item.Status
		if status == "" {
			status = "pending"
		}
		line := fmt.Sprintf("[%s] %s: %s", status, item.ID, item.Content)
		if item.Priority != "" {
			line += fmt.Sprintf(" (priority: %s)", item.Priority)
		}
		if item.File != "" {
			line += fmt.Sprintf(" (%s:%d)", item.File, item.Line)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n"), nil
}

func todoAdd(args map[string]interface{}) (string, error) {
	content, _ := args["content"].(string)
	if content == "" {
		return "", fmt.Errorf("content is required for add")
	}
	id := fmt.Sprintf("todo_%d", time.Now().UnixNano())
	status, _ := args["status"].(string)
	if status == "" {
		status = "pending"
	}
	priority, _ := args["priority"].(string)

	TodoStore.Lock()
	TodoStore.Items = append(TodoStore.Items, TodoItem{
		ID:       id,
		Content:  content,
		Status:   status,
		Priority: priority,
	})
	TodoStore.Unlock()

	return fmt.Sprintf("Added todo %s: %s [%s]", id, content, status), nil
}

func todoUpdate(args map[string]interface{}) (string, error) {
	id, _ := args["id"].(string)
	if id == "" {
		return "", fmt.Errorf("id is required for update")
	}

	TodoStore.Lock()
	defer TodoStore.Unlock()

	for i, item := range TodoStore.Items {
		if item.ID == id {
			if content, ok := args["content"].(string); ok && content != "" {
				TodoStore.Items[i].Content = content
			}
			if status, ok := args["status"].(string); ok && status != "" {
				TodoStore.Items[i].Status = status
			}
			if priority, ok := args["priority"].(string); ok {
				TodoStore.Items[i].Priority = priority
			}
			return fmt.Sprintf("Updated todo %s", id), nil
		}
	}
	return "", fmt.Errorf("todo %s not found", id)
}

func todoRemove(args map[string]interface{}) (string, error) {
	id, _ := args["id"].(string)
	if id == "" {
		return "", fmt.Errorf("id is required for remove")
	}

	TodoStore.Lock()
	defer TodoStore.Unlock()

	for i, item := range TodoStore.Items {
		if item.ID == id {
			TodoStore.Items = append(TodoStore.Items[:i], TodoStore.Items[i+1:]...)
			return fmt.Sprintf("Removed todo %s", id), nil
		}
	}
	return "", fmt.Errorf("todo %s not found", id)
}

func todoClear() (string, error) {
	TodoStore.Lock()
	count := len(TodoStore.Items)
	TodoStore.Items = nil
	TodoStore.Unlock()
	return fmt.Sprintf("Cleared %d todos", count), nil
}

// ==========================
// references tool
// ==========================

func (r *Registry) registerReferencesTool() {
	r.Tools = append(r.Tools, Tool{
		Name:        "references",
		Description: "Find references to a symbol in the codebase. Uses grep as a fallback when LSP is not available.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["symbol"],
			"properties": {
				"symbol": {
					"type": "string",
					"description": "The symbol to find references for"
				},
				"file_path": {
					"type": "string",
					"description": "The file containing the symbol (helps narrow results)"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			symbol, _ := args["symbol"].(string)
			filePath, _ := args["file_path"].(string)
			if symbol == "" {
				return "", fmt.Errorf("symbol is required")
			}

			// Use word-boundary grep for symbol references
			pattern := fmt.Sprintf("\\b%s\\b", regexp.QuoteMeta(symbol))
			searchPath := "."
			if filePath != "" {
				searchPath = filepath.Dir(filePath)
			}

			piArgs := map[string]interface{}{
				"pattern": pattern,
				"path":   searchPath,
			}
			raw, _ := json.Marshal(piArgs)
			runtime := foundationpi.NewRuntime(".", nil)
			result, err := runtime.ExecuteTool(context.Background(), "", "grep", raw, nil)
			if err != nil {
				return "", err
			}
			return formatFoundationToolResult(result), nil
		},
	})
}

// ==========================
// lsp_restart tool
// ==========================

func (r *Registry) registerLSPRestartTool() {
	r.Tools = append(r.Tools, Tool{
		Name:        "lsp_restart",
		Description: "Restart the Language Server Protocol server for a given language. Useful when the LSP server becomes unresponsive or out of sync.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"language": {
					"type": "string",
					"description": "The language server to restart (go, typescript, python, etc.)"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			language, _ := args["language"].(string)
			// LSP management would be handled by MCP server integration
			return fmt.Sprintf("LSP restart requested for language: %s. Configure an LSP MCP server for full restart capability.", language), nil
		},
	})
}

// ==========================
// job management tools
// ==========================

func (r *Registry) registerJobManagementTools() {
	// job_output - read output from a background job
	r.Tools = append(r.Tools, Tool{
		Name:        "job_output",
		Description: "Get the output from a background job. Returns the current stdout/stderr output.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["job_id"],
			"properties": {
				"job_id": {
					"type": "string",
					"description": "The ID of the background job"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			jobID, _ := args["job_id"].(string)
			if jobID == "" {
				return "", fmt.Errorf("job_id is required")
			}

			job, ok := GlobalJobManager.GetJob(jobID)
			if !ok {
				return "", fmt.Errorf("job %s not found", jobID)
			}

			job.mu.Lock()
			defer job.mu.Unlock()

			output := job.Output.String()
			status := "running"
			if job.Done {
				status = fmt.Sprintf("exited with code %d", job.ExitCode)
			}

			if len(output) > 30000 {
				output = output[:30000] + "\n... (output truncated at 30000 bytes)"
			}

			return fmt.Sprintf("Job %s [%s]:\n%s", jobID, status, output), nil
		},
	})

	// job_kill - kill a background job
	r.Tools = append(r.Tools, Tool{
		Name:        "job_kill",
		Description: "Kill a running background job.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["job_id"],
			"properties": {
				"job_id": {
					"type": "string",
					"description": "The ID of the background job to kill"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			jobID, _ := args["job_id"].(string)
			if jobID == "" {
				return "", fmt.Errorf("job_id is required")
			}

			job, ok := GlobalJobManager.GetJob(jobID)
			if !ok {
				return "", fmt.Errorf("job %s not found", jobID)
			}

			job.mu.Lock()
			defer job.mu.Unlock()

			if job.Done {
				return fmt.Sprintf("Job %s already completed.", jobID), nil
			}

			job.cancel()
			return fmt.Sprintf("Job %s killed.", jobID), nil
		},
	})

	// job_list - list all background jobs
	r.Tools = append(r.Tools, Tool{
		Name:        "job_list",
		Description: "List all background jobs and their status.",
		Parameters:  json.RawMessage(`{"type":"object","properties":{},"additionalProperties":false}`),
		Execute: func(args map[string]interface{}) (string, error) {
			jobs := GlobalJobManager.ListJobs()
			if len(jobs) == 0 {
				return "No background jobs.", nil
			}

			var lines []string
			for _, job := range jobs {
				status := "running"
				if job.Done {
					status = fmt.Sprintf("exited(%d)", job.ExitCode)
				}
				desc := job.Description
				if desc == "" {
					desc = job.Command
				}
				lines = append(lines, fmt.Sprintf("%s [%s] %s (started %s)", job.ID, status, desc, job.StartTime.Format(time.RFC3339)))
			}
			return strings.Join(lines, "\n"), nil
		},
	})
}

// ==========================
// safe tool
// ==========================

func (r *Registry) registerSafeTool() {
	r.Tools = append(r.Tools, Tool{
		Name:        "safe",
		Description: "Check if an operation is safe to perform. Validates file paths, commands, and operations against security policies.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["operation"],
			"properties": {
				"operation": {
					"type": "string",
					"description": "The operation to check: read, write, delete, execute"
				},
				"target": {
					"type": "string",
					"description": "The target file path or command to check"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			operation, _ := args["operation"].(string)
			target, _ := args["target"].(string)

			if operation == "" {
				return "", fmt.Errorf("operation is required")
			}

			switch operation {
			case "read", "write", "delete":
				if target == "" {
					return "", fmt.Errorf("target path is required")
				}
				absPath, err := filepath.Abs(target)
				if err != nil {
					return fmt.Sprintf("UNSAFE: Cannot resolve path: %v", err), nil
				}
				// Basic safety checks
				if strings.Contains(absPath, "..") {
					return fmt.Sprintf("UNSAFE: Path traversal detected: %s", target), nil
				}
				return fmt.Sprintf("SAFE: %s operation on %s (resolved: %s)", operation, target, absPath), nil

			case "execute":
				if target == "" {
					return "", fmt.Errorf("target command is required")
				}
				// Check for dangerous commands
				dangerous := []string{"rm -rf /", "format ", "del /f /s /q C:", "shutdown", "mkfs.", "dd if="}
				lowerTarget := strings.ToLower(target)
				for _, d := range dangerous {
					if strings.Contains(lowerTarget, strings.ToLower(d)) {
						return fmt.Sprintf("UNSAFE: Potentially dangerous command: %s", target), nil
					}
				}
				return fmt.Sprintf("SAFE: Execute command: %s", target), nil

			default:
				return fmt.Sprintf("UNKNOWN: Unsupported operation: %s", operation), nil
			}
		},
	})
}

// ==========================
// download tool
// ==========================

func (r *Registry) registerDownloadTool() {
	r.Tools = append(r.Tools, Tool{
		Name:        "download",
		Description: "Download a file from a URL to a local path.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["url", "path"],
			"properties": {
				"url": {
					"type": "string",
					"description": "The URL to download from"
				},
				"path": {
					"type": "string",
					"description": "The local file path to save to"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			url, _ := args["url"].(string)
			path, _ := args["path"].(string)
			if url == "" {
				return "", fmt.Errorf("url is required")
			}
			if path == "" {
				return "", fmt.Errorf("path is required")
			}

			dir := filepath.Dir(path)
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return "", fmt.Errorf("failed to create directory: %w", err)
			}

			client := &http.Client{Timeout: 5 * time.Minute}
			resp, err := client.Get(url)
			if err != nil {
				return "", fmt.Errorf("failed to download: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
			}

			f, err := os.Create(path)
			if err != nil {
				return "", fmt.Errorf("failed to create file: %w", err)
			}
			defer f.Close()

			written, err := io.Copy(f, resp.Body)
			if err != nil {
				return "", fmt.Errorf("failed to write file: %w", err)
			}

			return fmt.Sprintf("Downloaded %s to %s (%d bytes)", url, path, written), nil
		},
	})
}

// ==========================
// sourcegraph tool
// ==========================

func (r *Registry) registerSourcegraphTool() {
	r.Tools = append(r.Tools, Tool{
		Name:        "sourcegraph",
		Description: "Search code on Sourcegraph. Returns matching code snippets from public repositories.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["query"],
			"properties": {
				"query": {
					"type": "string",
					"description": "The search query (supports Sourcegraph search syntax)"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			query, _ := args["query"].(string)
			if query == "" {
				return "", fmt.Errorf("query is required")
			}

			// Sourcegraph search via their public API
			url := fmt.Sprintf("https://sourcegraph.com/.api/search/stream?q=%s", query)
			client := &http.Client{Timeout: 15 * time.Second}
			resp, err := client.Get(url)
			if err != nil {
				return "", fmt.Errorf("failed to search Sourcegraph: %w", err)
			}
			defer resp.Body.Close()

			data, err := io.ReadAll(io.LimitReader(resp.Body, 50000))
			if err != nil {
				return "", fmt.Errorf("failed to read response: %w", err)
			}

			content := string(data)
			if len(content) > 20000 {
				content = content[:20000] + "\n... (truncated)"
			}
			return content, nil
		},
	})
}

// ==========================
// search tool (enhanced grep)
// ==========================

func (r *Registry) registerSearchTool() {
	r.Tools = append(r.Tools, Tool{
		Name:        "search",
		Description: "Search for a pattern in the codebase. Returns file paths and matching lines with context. Uses ripgrep if available, falls back to native Go search.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["pattern"],
			"properties": {
				"pattern": {
					"type": "string",
					"description": "The pattern to search for (supports regex)"
				},
				"path": {
					"type": "string",
					"description": "The directory or file to search in"
				},
				"ignore_case": {
					"type": "boolean",
					"description": "Case-insensitive search"
				},
				"literal": {
					"type": "boolean",
					"description": "Treat pattern as literal string"
				},
				"glob": {
					"type": "string",
					"description": "Filter files by glob pattern (e.g. '*.go')"
				},
				"context": {
					"type": "number",
					"description": "Number of context lines to show around matches"
				},
				"limit": {
					"type": "number",
					"description": "Maximum number of matches to return"
				}
			},
			"additionalProperties": false
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			piArgs := make(map[string]interface{})
			for k, v := range args {
				piArgs[k] = v
			}
			raw, _ := json.Marshal(piArgs)
			runtime := foundationpi.NewRuntime(".", nil)
			result, err := runtime.ExecuteTool(context.Background(), "", "grep", raw, nil)
			if err != nil {
				return "", err
			}
			return formatFoundationToolResult(result), nil
		},
	})
}

// normalizeLineEndings converts CRLF to LF.
func normalizeLineEndings(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}
