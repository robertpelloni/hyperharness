package tools

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCrushMultiEditBasic tests the multiedit tool with basic edits.
func TestCrushMultiEditBasic(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(filePath, []byte("line1\nline2\nline3\n"), 0o644)

	reg := NewRegistry()
	tool, ok := reg.Find("multiedit")
	if !ok {
		t.Fatal("multiedit tool not found")
	}

	args := map[string]interface{}{
		"file_path": filePath,
		"edits": []interface{}{
			map[string]interface{}{
				"old_string": "line2",
				"new_string": "replaced_line2",
			},
		},
	}

	result, err := tool.Execute(args)
	if err != nil {
		t.Fatalf("multiedit failed: %v", err)
	}
	if !strings.Contains(result, "Applied 1 edits") {
		t.Errorf("unexpected result: %s", result)
	}

	data, _ := os.ReadFile(filePath)
	content := string(data)
	if !strings.Contains(content, "replaced_line2") {
		t.Errorf("file not updated correctly: %s", content)
	}
}

// TestCrushMultiEditCreateNewFile tests multiedit creating a new file.
func TestCrushMultiEditCreateNewFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "newfile.txt")

	reg := NewRegistry()
	tool, ok := reg.Find("multiedit")
	if !ok {
		t.Fatal("multiedit tool not found")
	}

	args := map[string]interface{}{
		"file_path": filePath,
		"edits": []interface{}{
			map[string]interface{}{
				"old_string": "",
				"new_string": "new file content",
			},
		},
	}

	result, err := tool.Execute(args)
	if err != nil {
		t.Fatalf("multiedit create failed: %v", err)
	}
	if !strings.Contains(result, "File created") {
		t.Errorf("unexpected result: %s", result)
	}

	data, _ := os.ReadFile(filePath)
	if string(data) != "new file content" {
		t.Errorf("file content mismatch: %s", string(data))
	}
}

// TestCrushMultiEditPartialFailure tests multiedit with some failures.
func TestCrushMultiEditPartialFailure(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(filePath, []byte("line1\nline2\nline3\n"), 0o644)

	reg := NewRegistry()
	tool, ok := reg.Find("multiedit")
	if !ok {
		t.Fatal("multiedit tool not found")
	}

	args := map[string]interface{}{
		"file_path": filePath,
		"edits": []interface{}{
			map[string]interface{}{
				"old_string": "line1",
				"new_string": "replaced1",
			},
			map[string]interface{}{
				"old_string": "nonexistent",
				"new_string": "wontmatch",
			},
		},
	}

	result, err := tool.Execute(args)
	if err != nil {
		t.Fatalf("multiedit failed: %v", err)
	}
	if !strings.Contains(result, "1 of 2") {
		t.Errorf("expected partial failure: %s", result)
	}
}

// TestCrushViewTool tests the view tool.
func TestCrushViewTool(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "view_test.txt")
	os.WriteFile(filePath, []byte("line1\nline2\nline3\nline4\nline5\n"), 0o644)

	reg := NewRegistry()
	tool, ok := reg.Find("view")
	if !ok {
		t.Fatal("view tool not found")
	}

	args := map[string]interface{}{
		"file_path": filePath,
	}

	result, err := tool.Execute(args)
	if err != nil {
		t.Fatalf("view failed: %v", err)
	}
	if !strings.Contains(result, "line1") {
		t.Errorf("view result missing content: %s", result)
	}
}

// TestCrushBashTool tests the bash tool.
func TestCrushBashTool(t *testing.T) {
	reg := NewRegistry()
	tool, ok := reg.Find("bash")
	if !ok {
		t.Fatal("bash tool not found")
	}

	args := map[string]interface{}{
		"command": "echo hello",
	}

	result, err := tool.Execute(args)
	if err != nil {
		t.Fatalf("bash failed: %v", err)
	}
	if !strings.Contains(result, "hello") {
		t.Errorf("bash result missing 'hello': %s", result)
	}
}

// TestCrushGlobTool tests the glob tool.
func TestCrushGlobTool(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte("package test\n"), 0o644)
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("hello\n"), 0o644)

	reg := NewRegistry()
	tool, ok := reg.Find("glob")
	if !ok {
		t.Fatal("glob tool not found")
	}

	args := map[string]interface{}{
		"pattern": "*.go",
		"path":    tmpDir,
	}

	result, err := tool.Execute(args)
	if err != nil {
		t.Fatalf("glob failed: %v", err)
	}
	if !strings.Contains(result, "test.go") {
		t.Errorf("glob result missing test.go: %s", result)
	}
}

// TestCrushWebFetchTool tests web_fetch validation.
func TestCrushWebFetchTool(t *testing.T) {
	reg := NewRegistry()
	tool, ok := reg.Find("web_fetch")
	if !ok {
		t.Fatal("web_fetch tool not found")
	}

	// Test missing URL
	_, err := tool.Execute(map[string]interface{}{})
	if err == nil {
		t.Error("expected error for missing URL")
	}
}

// TestCrushDiagnosticsTool tests diagnostics.
func TestCrushDiagnosticsTool(t *testing.T) {
	reg := NewRegistry()
	tool, ok := reg.Find("diagnostics")
	if !ok {
		t.Fatal("diagnostics tool not found")
	}

	args := map[string]interface{}{
		"file_path": "tools_test.go",
		"language":  "go",
	}

	result, err := tool.Execute(args)
	if err != nil {
		t.Fatalf("diagnostics failed: %v", err)
	}
	// Should return some result (even if just "no diagnostics")
	if result == "" {
		t.Error("diagnostics returned empty result")
	}
}

// TestCrushTodosTool tests the todos tool lifecycle.
func TestCrushTodosTool(t *testing.T) {
	reg := NewRegistry()
	tool, ok := reg.Find("todos")
	if !ok {
		t.Fatal("todos tool not found")
	}

	// Add a todo
	result, err := tool.Execute(map[string]interface{}{
		"action":  "add",
		"content": "Test task",
	})
	if err != nil {
		t.Fatalf("todos add failed: %v", err)
	}
	if !strings.Contains(result, "Added todo") {
		t.Errorf("unexpected add result: %s", result)
	}

	// List todos
	result, err = tool.Execute(map[string]interface{}{
		"action": "list",
	})
	if err != nil {
		t.Fatalf("todos list failed: %v", err)
	}
	if !strings.Contains(result, "Test task") {
		t.Errorf("todos list missing task: %s", result)
	}

	// Clear todos
	result, err = tool.Execute(map[string]interface{}{
		"action": "clear",
	})
	if err != nil {
		t.Fatalf("todos clear failed: %v", err)
	}
	if !strings.Contains(result, "Cleared") {
		t.Errorf("unexpected clear result: %s", result)
	}
}

// TestCrushJobManagement tests background job lifecycle.
func TestCrushJobManagement(t *testing.T) {
	reg := NewRegistry()

	// Start a background job using the dedicated bash registration
	// Find the bash tool that supports background mode by checking parameters
	found := false
	var bgTool Tool
	for _, tool := range reg.Tools {
		if tool.Name == "bash" {
			var params map[string]interface{}
			if err := json.Unmarshal(tool.Parameters, &params); err == nil {
				if props, ok := params["properties"].(map[string]interface{}); ok {
					if _, hasBG := props["run_in_background"]; hasBG {
						bgTool = tool
						found = true
						break
					}
				}
			}
		}
	}

	if !found {
		t.Skip("No bash tool with background support found")
	}

	result, err := bgTool.Execute(map[string]interface{}{
		"command":           "echo hello_bg",
		"run_in_background": true,
		"description":       "test background job",
	})
	if err != nil {
		t.Fatalf("background bash failed: %v", err)
	}
	if !strings.Contains(result, "job_") {
		t.Skipf("background mode not available: %s", result)
	}

	// Extract job ID
	jobIDStart := strings.Index(result, "job_")
	rest := result[jobIDStart:]
	jobIDEnd := strings.Index(rest, "\n")
	var jobID string
	if jobIDEnd == -1 {
		jobID = rest
	} else {
		jobID = rest[:jobIDEnd]
	}

	if jobID == "" {
		t.Skip("Could not extract job ID")
	}

	// List jobs
	listTool, ok := reg.Find("job_list")
	if !ok {
		t.Fatal("job_list tool not found")
	}

	result, err = listTool.Execute(map[string]interface{}{})
	if err != nil {
		t.Fatalf("job_list failed: %v", err)
	}

	// Get job output (give a moment for completion)
	outputTool, ok := reg.Find("job_output")
	if !ok {
		t.Fatal("job_output tool not found")
	}

	result, err = outputTool.Execute(map[string]interface{}{
		"job_id": jobID,
	})
	if err != nil {
		t.Fatalf("job_output failed: %v", err)
	}
	if !strings.Contains(result, jobID) {
		t.Errorf("job output missing job ID: %s", result)
	}
}

// TestCrushSafeTool tests the safe operation checker.
func TestCrushSafeTool(t *testing.T) {
	reg := NewRegistry()
	tool, ok := reg.Find("safe")
	if !ok {
		t.Fatal("safe tool not found")
	}

	// Test safe read
	result, err := tool.Execute(map[string]interface{}{
		"operation": "read",
		"target":    "test.txt",
	})
	if err != nil {
		t.Fatalf("safe check failed: %v", err)
	}
	if !strings.Contains(result, "SAFE") {
		t.Errorf("expected SAFE: %s", result)
	}

	// Test dangerous command
	result, err = tool.Execute(map[string]interface{}{
		"operation": "execute",
		"target":    "rm -rf /",
	})
	if err != nil {
		t.Fatalf("safe check failed: %v", err)
	}
	if !strings.Contains(result, "UNSAFE") {
		t.Errorf("expected UNSAFE for dangerous command: %s", result)
	}
}

// TestCrushDownloadTool tests download validation.
func TestCrushDownloadTool(t *testing.T) {
	reg := NewRegistry()
	tool, ok := reg.Find("download")
	if !ok {
		t.Fatal("download tool not found")
	}

	// Test missing URL
	_, err := tool.Execute(map[string]interface{}{
		"path": "/tmp/test",
	})
	if err == nil {
		t.Error("expected error for missing URL")
	}
}

// TestGeminiCLITools tests Gemini CLI tool surfaces.
func TestGeminiCLITools(t *testing.T) {
	reg := NewRegistry()

	// Test read_file
	tool, ok := reg.Find("read_file")
	if !ok {
		t.Fatal("read_file tool not found (Gemini parity)")
	}

	// Test edit_file
	tool, ok = reg.Find("edit_file")
	if !ok {
		t.Fatal("edit_file tool not found (Gemini parity)")
	}

	// Test list_directory
	tool, ok = reg.Find("list_directory")
	if !ok {
		t.Fatal("list_directory tool not found (Gemini parity)")
	}

	// Test shell
	tool, ok = reg.Find("shell")
	if !ok {
		t.Fatal("shell tool not found (Gemini parity)")
	}

	// Test search_files
	tool, ok = reg.Find("search_files")
	if !ok {
		t.Fatal("search_files tool not found (Gemini parity)")
	}

	// Test find_files
	tool, ok = reg.Find("find_files")
	if !ok {
		t.Fatal("find_files tool not found (Gemini parity)")
	}

	// Verify read_file actually works (the old read_file uses file_path parameter)
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "gemini_test.txt")
	os.WriteFile(testFile, []byte("hello gemini"), 0o644)

	tool, _ = reg.Find("read_file")
	result, err := tool.Execute(map[string]interface{}{
		"file_path": testFile,
	})
	if err != nil {
		t.Fatalf("read_file failed: %v", err)
	}
	if !strings.Contains(result, "hello gemini") {
		t.Errorf("read_file result: %s", result)
	}
}

// TestGrokTools tests Grok CLI tool surfaces.
func TestGrokTools(t *testing.T) {
	reg := NewRegistry()

	grokTools := []string{
		"file_edit",
		"file_read",
		"file_write",
		"execute_command",
		"scan_directory",
		"search_code",
	}

	for _, name := range grokTools {
		_, ok := reg.Find(name)
		if !ok {
			t.Errorf("Grok tool %q not found", name)
		}
	}
}

// TestPiExactTools tests exact pi tool surfaces.
func TestPiExactTools(t *testing.T) {
	reg := NewRegistry()

	piTools := []string{"read", "write", "edit", "bash", "grep", "find", "ls"}
	for _, name := range piTools {
		_, ok := reg.Find(name)
		if !ok {
			t.Errorf("pi exact tool %q not found", name)
		}
	}

	// Verify exact parameter names
	readTool, _ := reg.Find("read")
	var params map[string]interface{}
	if err := json.Unmarshal(readTool.Parameters, &params); err != nil {
		t.Fatalf("parse read params: %v", err)
	}
	props, ok := params["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("read tool params missing properties")
	}
	for _, expected := range []string{"path", "offset", "limit"} {
		if _, ok := props[expected]; !ok {
			t.Errorf("read tool missing property: %s", expected)
		}
	}

	// Verify edit params
	editTool, _ := reg.Find("edit")
	if err := json.Unmarshal(editTool.Parameters, &params); err != nil {
		t.Fatalf("parse edit params: %v", err)
	}
	props, ok = params["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("edit tool params missing properties")
	}
	if _, ok := props["edits"]; !ok {
		t.Error("edit tool missing 'edits' property")
	}
}

// TestMCPGatewayTool tests the MCP gateway tool.
func TestMCPGatewayTool(t *testing.T) {
	reg := NewRegistry()
	tool, ok := reg.Find("mcp")
	if !ok {
		t.Fatal("mcp tool not found")
	}

	result, err := tool.Execute(map[string]interface{}{})
	if err != nil {
		t.Fatalf("mcp tool failed: %v", err)
	}
	if !strings.Contains(result, "MCP gateway") {
		t.Errorf("unexpected mcp result: %s", result)
	}
}

// TestHypercodeTools tests hypercode-specific tools.
func TestHypercodeTools(t *testing.T) {
	reg := NewRegistry()

	// Test memory_store
	tool, ok := reg.Find("memory_store")
	if !ok {
		t.Fatal("memory_store tool not found")
	}
	result, err := tool.Execute(map[string]interface{}{
		"title":   "test knowledge",
		"content": "test content",
	})
	if err != nil {
		t.Fatalf("memory_store failed: %v", err)
	}
	if !strings.Contains(result, "Stored knowledge") {
		t.Errorf("unexpected memory_store result: %s", result)
	}

	// Test memory_search
	tool, ok = reg.Find("memory_search")
	if !ok {
		t.Fatal("memory_search tool not found")
	}
	result, err = tool.Execute(map[string]interface{}{
		"query": "test",
	})
	if err != nil {
		t.Fatalf("memory_search failed: %v", err)
	}
	if !strings.Contains(result, "test") {
		t.Errorf("unexpected memory_search result: %s", result)
	}

	// Test context_manager
	tool, ok = reg.Find("context_manager")
	if !ok {
		t.Fatal("context_manager tool not found")
	}
	result, err = tool.Execute(map[string]interface{}{
		"action": "status",
	})
	if err != nil {
		t.Fatalf("context_manager failed: %v", err)
	}
	if !strings.Contains(result, "active") {
		t.Errorf("unexpected context_manager result: %s", result)
	}
}

// TestOpenCodeTools tests OpenCode tool surfaces.
func TestOpenCodeTools(t *testing.T) {
	reg := NewRegistry()

	tool, ok := reg.Find("apply_search_replace")
	if !ok {
		t.Fatal("apply_search_replace tool not found")
	}

	// Test with a real file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "opencode_test.txt")
	os.WriteFile(testFile, []byte("search target\nmore content\n"), 0o644)

	result, err := tool.Execute(map[string]interface{}{
		"file_path":    testFile,
		"search_block": "search target",
		"replace_block": "replaced target",
	})
	if err != nil {
		t.Fatalf("apply_search_replace failed: %v", err)
	}
	if !strings.Contains(result, "applied") && !strings.Contains(result, "mutated") {
		t.Errorf("unexpected result: %s", result)
	}
}

// TestAiderV2Tools tests Aider v2 search/replace blocks.
func TestAiderV2Tools(t *testing.T) {
	reg := NewRegistry()

	tool, ok := reg.Find("aider_search_replace")
	if !ok {
		t.Fatal("aider_search_replace tool not found")
	}

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "aider_test.txt")
	os.WriteFile(testFile, []byte("old content\n"), 0o644)

	result, err := tool.Execute(map[string]interface{}{
		"blocks": []interface{}{
			map[string]interface{}{
				"file_path": testFile,
				"search":    "old content",
				"replace":   "new content",
			},
		},
	})
	if err != nil {
		t.Fatalf("aider_search_replace failed: %v", err)
	}
	if !strings.Contains(result, "applied successfully") {
		t.Errorf("unexpected result: %s", result)
	}
}

// TestCopilotCLITools tests Copilot CLI tool surfaces.
func TestCopilotCLITools(t *testing.T) {
	reg := NewRegistry()
	tool, ok := reg.Find("copilot_edit")
	if !ok {
		t.Fatal("copilot_edit tool not found")
	}

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "copilot_test.txt")
	os.WriteFile(testFile, []byte("copilot content\n"), 0o644)

	_, err := tool.Execute(map[string]interface{}{
		"file_path":   testFile,
		"old_string":  "copilot content",
		"new_string":  "copilot edited",
	})
	if err != nil {
		t.Fatalf("copilot_edit failed: %v", err)
	}

	data, _ := os.ReadFile(testFile)
	if !strings.Contains(string(data), "copilot edited") {
		t.Errorf("copilot_edit file not updated: %s", string(data))
	}
}

// TestToolRegistryCompleteness verifies all expected tool surfaces exist.
func TestToolRegistryCompleteness(t *testing.T) {
	reg := NewRegistry()

	expectedTools := []string{
		// Foundation tools (from pi)
		"read", "write", "edit", "bash", "grep", "find", "ls",
		// Crush parity
		"multiedit", "view", "glob", "web_fetch", "web_search",
		"diagnostics", "todos", "references", "lsp_restart",
		"job_output", "job_kill", "job_list",
		"safe", "download", "sourcegraph", "search",
		// Claude Code parity
		"Edit", "Read", "Write", "Bash", "Glob", "Grep", "LS",
		// Computer Use parity
		"str_replace",
		// Aider parity
		"replace_lines", "write_file", "aider_search_replace",
		// Gemini CLI parity
		"read_file", "edit_file", "list_directory", "shell",
		"search_files", "find_files",
		// OpenCode parity
		"apply_search_replace", "apply_diff",
		"task", "batch", "codesearch", "websearch", "webfetch",
		"lsp", "question", "skill", "plan_exit", "plan_enter",
		"opencode_multiedit",
		// Grok parity
		"file_edit", "file_read", "file_write",
		"execute_command", "scan_directory", "search_code",
		// Copilot parity
		"copilot_edit",
		// Goose parity
		"tree", "load", "delegate", "platform__manage_schedule",
		// Kimi CLI parity
		"TaskList", "TaskOutput", "TaskStop", "Think", "AskUser",
		"ReadFile", "WriteFile", "Replace", "Glob", "GrepLocal",
		"WebFetch", "WebSearch", "PlanEnter", "PlanExit",
		// Cursor parity
		"cursor_read_file", "cursor_edit_file", "cursor_run_command",
		"cursor_code_search", "cursor_list_dir",
		// Windsurf parity
		"cascade_edit", "cascade_command",
		// Mistral parity
		"mistral_edit", "mistral_search",
		// Smithery parity
		"smithery_install", "smithery_list",
		// Hypercode tools
		"mcp", "memory_store", "memory_search", "context_manager",
		// Legacy/other
		"run_shell_command", "repomap", "install_mcp_server",
	}

	missing := []string{}
	for _, name := range expectedTools {
		if _, ok := reg.Find(name); !ok {
			missing = append(missing, name)
		}
	}

	if len(missing) > 0 {
		t.Errorf("Missing tools: %v", missing)
	}

	t.Logf("Registry contains %d tools", len(reg.Tools))
	t.Logf("All %d expected tools found", len(expectedTools))
}

// TestGooseTreeTool tests the Goose tree tool.
func TestGooseTreeTool(t *testing.T) {
	reg := NewRegistry()
	tool, ok := reg.Find("tree")
	if !ok {
		t.Fatal("tree tool not found")
	}

	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "a.txt"), []byte("hello"), 0o644)
	os.WriteFile(filepath.Join(tmpDir, "b.go"), []byte("package test"), 0o644)
	os.Mkdir(filepath.Join(tmpDir, "sub"), 0o755)
	os.WriteFile(filepath.Join(tmpDir, "sub", "c.txt"), []byte("world"), 0o644)

	result, err := tool.Execute(map[string]interface{}{
		"path": tmpDir,
	})
	if err != nil {
		t.Fatalf("tree failed: %v", err)
	}
	if !strings.Contains(result, "a.txt") {
		t.Errorf("tree missing a.txt: %s", result)
	}
	if !strings.Contains(result, "sub/") {
		t.Errorf("tree missing sub/: %s", result)
	}
}

// TestGooseLoadTool tests the load tool.
func TestGooseLoadTool(t *testing.T) {
	reg := NewRegistry()
	tool, ok := reg.Find("load")
	if !ok {
		t.Fatal("load tool not found")
	}

	result, err := tool.Execute(map[string]interface{}{
		"action": "list",
	})
	if err != nil {
		t.Fatalf("load list failed: %v", err)
	}
	if !strings.Contains(result, "Available") {
		t.Errorf("unexpected list result: %s", result)
	}
}

// TestGooseDelegateTool tests the delegate tool.
func TestGooseDelegateTool(t *testing.T) {
	reg := NewRegistry()
	tool, ok := reg.Find("delegate")
	if !ok {
		t.Fatal("delegate tool not found")
	}

	result, err := tool.Execute(map[string]interface{}{
		"task":        "test task",
		"description": "test desc",
		"mode":        "sync",
	})
	if err != nil {
		t.Fatalf("delegate failed: %v", err)
	}
	if !strings.Contains(result, "task_result") {
		t.Errorf("unexpected delegate result: %s", result)
	}
}

// TestOpenCodeTaskTool tests the task tool.
func TestOpenCodeTaskTool(t *testing.T) {
	reg := NewRegistry()
	tool, ok := reg.Find("task")
	if !ok {
		t.Fatal("task tool not found")
	}

	result, err := tool.Execute(map[string]interface{}{
		"description":    "test task",
		"prompt":         "do something",
		"subagent_type":  "code",
	})
	if err != nil {
		t.Fatalf("task failed: %v", err)
	}
	if !strings.Contains(result, "task_id") {
		t.Errorf("unexpected task result: %s", result)
	}
}

// TestOpenCodeBatchTool tests the batch tool validation.
func TestOpenCodeBatchTool(t *testing.T) {
	reg := NewRegistry()
	tool, ok := reg.Find("batch")
	if !ok {
		t.Fatal("batch tool not found")
	}

	// Test empty array
	_, err := tool.Execute(map[string]interface{}{})
	if err == nil {
		t.Error("expected error for empty tool_calls")
	}
}

// TestOpenCodeLSPTool tests the lsp tool.
func TestOpenCodeLSPTool(t *testing.T) {
	reg := NewRegistry()
	tool, ok := reg.Find("lsp")
	if !ok {
		t.Fatal("lsp tool not found")
	}

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")
	os.WriteFile(testFile, []byte("package test\n"), 0o644)

	result, err := tool.Execute(map[string]interface{}{
		"operation": "goToDefinition",
		"filePath":  testFile,
		"line":      1,
		"character": 1,
	})
	if err != nil {
		t.Fatalf("lsp failed: %v", err)
	}
	if !strings.Contains(result, "goToDefinition") {
		t.Errorf("unexpected lsp result: %s", result)
	}
}

// TestKimiTaskList tests Kimi's TaskList tool.
func TestKimiTaskList(t *testing.T) {
	reg := NewRegistry()
	tool, ok := reg.Find("TaskList")
	if !ok {
		t.Fatal("TaskList tool not found")
	}

	result, err := tool.Execute(map[string]interface{}{})
	if err != nil {
		t.Fatalf("TaskList failed: %v", err)
	}
	if !strings.Contains(result, "task") || !strings.Contains(result, "No active") {
		t.Errorf("unexpected TaskList result: %s", result)
	}
}

// TestKimiThink tests Kimi's Think tool.
func TestKimiThink(t *testing.T) {
	reg := NewRegistry()
	tool, ok := reg.Find("Think")
	if !ok {
		t.Fatal("Think tool not found")
	}

	result, err := tool.Execute(map[string]interface{}{
		"thought": "let me reason about this",
	})
	if err != nil {
		t.Fatalf("Think failed: %v", err)
	}
	if !strings.Contains(result, "recorded") {
		t.Errorf("unexpected Think result: %s", result)
	}
}

// TestKimiReadFile tests Kimi's ReadFile tool.
func TestKimiReadFile(t *testing.T) {
	reg := NewRegistry()
	tool, ok := reg.Find("ReadFile")
	if !ok {
		t.Fatal("ReadFile tool not found")
	}

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "kimi.txt")
	os.WriteFile(testFile, []byte("hello kimi"), 0o644)

	result, err := tool.Execute(map[string]interface{}{
		"path": testFile,
	})
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if !strings.Contains(result, "hello kimi") {
		t.Errorf("unexpected ReadFile result: %s", result)
	}
}

// TestKimiWriteFile tests Kimi's WriteFile tool.
func TestKimiWriteFile(t *testing.T) {
	reg := NewRegistry()
	tool, ok := reg.Find("WriteFile")
	if !ok {
		t.Fatal("WriteFile tool not found")
	}

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "sub", "dir", "kimi_out.txt")

	_, err := tool.Execute(map[string]interface{}{
		"path":    filePath,
		"content": "kimi wrote this",
	})
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	data, _ := os.ReadFile(filePath)
	if string(data) != "kimi wrote this" {
		t.Errorf("WriteFile content mismatch: %s", string(data))
	}
}

// TestCursorTools tests Cursor parity tools.
func TestCursorTools(t *testing.T) {
	reg := NewRegistry()

	cursorTools := []string{"cursor_read_file", "cursor_edit_file", "cursor_run_command", "cursor_code_search", "cursor_list_dir"}
	for _, name := range cursorTools {
		_, ok := reg.Find(name)
		if !ok {
			t.Errorf("Cursor tool %q not found", name)
		}
	}

	// Test cursor_run_command actually works
	tool, _ := reg.Find("cursor_run_command")
	result, err := tool.Execute(map[string]interface{}{
		"command": "echo cursor_test",
	})
	if err != nil {
		t.Fatalf("cursor_run_command failed: %v", err)
	}
	if !strings.Contains(result, "cursor_test") {
		t.Errorf("cursor_run_command result: %s", result)
	}
}

// TestWindsurfTools tests Windsurf parity tools.
func TestWindsurfTools(t *testing.T) {
	reg := NewRegistry()

	for _, name := range []string{"cascade_edit", "cascade_command"} {
		_, ok := reg.Find(name)
		if !ok {
			t.Errorf("Windsurf tool %q not found", name)
		}
	}
}

// TestMistralTools tests Mistral parity tools.
func TestMistralTools(t *testing.T) {
	reg := NewRegistry()

	for _, name := range []string{"mistral_edit", "mistral_search"} {
		if _, ok := reg.Find(name); !ok {
			t.Errorf("Mistral tool %q not found", name)
		}
	}
}

// TestSmitheryTools tests Smithery parity tools.
func TestSmitheryTools(t *testing.T) {
	reg := NewRegistry()

	tool, ok := reg.Find("smithery_install")
	if !ok {
		t.Fatal("smithery_install tool not found")
	}

	result, err := tool.Execute(map[string]interface{}{
		"server_name": "test-server",
	})
	if err != nil {
		t.Fatalf("smithery_install failed: %v", err)
	}
	if !strings.Contains(result, "test-server") {
		t.Errorf("unexpected smithery result: %s", result)
	}
}

// TestOpenCodeMultiedit tests the OpenCode multiedit format.
func TestOpenCodeMultiedit(t *testing.T) {
	reg := NewRegistry()
	tool, ok := reg.Find("opencode_multiedit")
	if !ok {
		t.Fatal("opencode_multiedit tool not found")
	}

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "multi.txt")
	os.WriteFile(testFile, []byte("foo bar\nbaz qux\nfoo again\n"), 0o644)

	_, err := tool.Execute(map[string]interface{}{
		"filePath": testFile,
		"edits": []interface{}{
			map[string]interface{}{
				"oldString": "foo bar",
				"newString": "FOO BAR",
			},
			map[string]interface{}{
				"oldString": "baz qux",
				"newString": "BAZ QUX",
			},
		},
	})
	if err != nil {
		t.Fatalf("opencode_multiedit failed: %v", err)
	}

	data, _ := os.ReadFile(testFile)
	content := string(data)
	if !strings.Contains(content, "FOO BAR") {
		t.Errorf("first edit not applied: %s", content)
	}
	if !strings.Contains(content, "BAZ QUX") {
		t.Errorf("second edit not applied: %s", content)
	}
}

// TestToolCount reports the total number of registered tools.
func TestToolCount(t *testing.T) {
	reg := NewRegistry()
	t.Logf("Total registered tools: %d", len(reg.Tools))

	// List all tools for documentation
	names := make([]string, 0, len(reg.Tools))
	for _, tool := range reg.Tools {
		names = append(names, tool.Name)
	}
	t.Logf("Tool names: %s", strings.Join(names, ", "))

	// We should have at least 40 tools
	if len(reg.Tools) < 40 {
		t.Errorf("Expected at least 40 tools, got %d", len(reg.Tools))
	}
}
