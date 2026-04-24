import sys

with open("tools/claude_code_parity.go", "r") as f:
    content = f.read()

pattern = """// ---- Claude Code: TodoWrite ----
// Manages a session-level task checklist with pending/in_progress/completed states.
// Models use this proactively to track multi-step tasks.

type ClaudeTodoItem struct {
	Content    string `json:"content"`              // Imperative: "Fix auth bug"
	ActiveForm string `json:"activeForm"`            // Present continuous: "Fixing auth bug"
	Status     string `json:"status"`                // "pending" | "in_progress" | "completed"
}

var (
	globalTodos []ClaudeTodoItem
	todoMu      sync.Mutex
)

func init() {
	globalTodos = make([]ClaudeTodoItem, 0)
}

func registerClaudeCodeTools(r *Registry) {
	// TodoWrite - Claude Code's task tracking tool
	r.Tools = append(r.Tools, Tool{
		Name:        "TodoWrite",
		Description: "Update the todo list for the current session. To be used proactively and often to track progress and pending tasks. Make sure that at least one task is in_progress at all times. Always provide both content (imperative) and activeForm (present continuous) for each task.",
		Parameters:  json.RawMessage(`{"type":"object","properties":{"todos":{"type":"array","items":{"type":"object","properties":{"content":{"type":"string","description":"The imperative form describing what needs to be done"},"activeForm":{"type":"string","description":"The present continuous form shown during execution"},"status":{"type":"string","enum":["pending","in_progress","completed"],"description":"Task status"}},"required":["content","activeForm","status"]}}},"required":["todos"]}`),
		Execute: func(args map[string]interface{}) (string, error) {
			todoMu.Lock()
			defer todoMu.Unlock()

			todosRaw, ok := args["todos"]
			if !ok {
				return "", fmt.Errorf("todos parameter required")
			}

			todosArr, ok := todosRaw.([]interface{})
			if !ok {
				return "", fmt.Errorf("todos must be an array")
			}

			oldTodos := make([]ClaudeTodoItem, len(globalTodos))
			copy(oldTodos, globalTodos)

			newTodos := make([]ClaudeTodoItem, 0, len(todosArr))
			for i, t := range todosArr {
				tMap, ok := t.(map[string]interface{})
				if !ok {
					continue
				}
				item := ClaudeTodoItem{
					Content:    getStr(tMap, "content"),
					ActiveForm: getStr(tMap, "activeForm"),
					Status:     getStr(tMap, "status"),
				}
				if item.Status == "" {
					item.Status = "pending"
				}
				_ = i
				newTodos = append(newTodos, item)
			}

			globalTodos = newTodos

			result := map[string]interface{}{
				"oldTodos": oldTodos,
				"newTodos": newTodos,
			}
			data, _ := json.Marshal(result)
			return string(data), nil
		},
	})"""

replacement = """// ---- Claude Code: TodoWrite ----
// Manages a session-level task checklist with pending/in_progress/completed states.
// Models use this proactively to track multi-step tasks.

var (
	// Shared global store for TodoWrite
	todoStore = NewSessionTodoStore(".hypercode/todos.json")
)

func registerClaudeCodeTools(r *Registry) {
	// TodoWrite - Claude Code's task tracking tool
	r.Tools = append(r.Tools, Tool{
		Name:        "TodoWrite",
		Description: "Update the todo list for the current session. To be used proactively and often to track progress and pending tasks. Make sure that at least one task is in_progress at all times. Always provide both content (imperative) and activeForm (present continuous) for each task.",
		Parameters:  json.RawMessage(`{"type":"object","properties":{"todos":{"type":"array","items":{"type":"object","properties":{"content":{"type":"string","description":"The imperative form describing what needs to be done"},"activeForm":{"type":"string","description":"The present continuous form shown during execution"},"status":{"type":"string","enum":["pending","in_progress","completed"],"description":"Task status"}},"required":["content","activeForm","status"]}}},"required":["todos"]}`),
		Execute: func(args map[string]interface{}) (string, error) {
			todosRaw, ok := args["todos"]
			if !ok {
				return "", fmt.Errorf("todos parameter required")
			}

			todosArr, ok := todosRaw.([]interface{})
			if !ok {
				return "", fmt.Errorf("todos must be an array")
			}

			newTodos := make([]SessionTodo, 0, len(todosArr))
			for _, t := range todosArr {
				tMap, ok := t.(map[string]interface{})
				if !ok {
					continue
				}
				item := SessionTodo{
					Content:    getStr(tMap, "content"),
					ActiveForm: getStr(tMap, "activeForm"),
					Status:     getStr(tMap, "status"),
				}
				if item.Status == "" {
					item.Status = "pending"
				}
				newTodos = append(newTodos, item)
			}

			oldTodos, currentTodos := todoStore.Update(newTodos)

			result := map[string]interface{}{
				"oldTodos": oldTodos,
				"newTodos": currentTodos,
			}
			data, _ := json.Marshal(result)
			return string(data), nil
		},
	})"""

if pattern in content:
    new_content = content.replace(pattern, replacement)
    with open("tools/claude_code_parity.go", "w") as f:
        f.write(new_content)
    print("Replaced successfully")
else:
    print("Pattern not found")
