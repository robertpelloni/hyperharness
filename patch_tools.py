import sys
import re

with open("tools/claude_code_parity.go", "r") as f:
    content = f.read()

# Add imports
if '"context"' not in content:
    content = content.replace(
        '"github.com/robertpelloni/hyperharness/foundation/pi"',
        '"context"\n\t"github.com/robertpelloni/hyperharness/foundation/pi"\n\t"github.com/robertpelloni/hyperharness/internal/subagents"\n\t"github.com/robertpelloni/hyperharness/internal/skills"\n\t"github.com/robertpelloni/hyperharness/config"'
    )

# Fix TodoWrite
todo_old = """type ClaudeTodoItem struct {
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
}"""

todo_new = """var (
	// Shared global store for TodoWrite
	todoStore = NewSessionTodoStore(".hypercode/todos.json")
)"""

content = content.replace(todo_old, todo_new)

todo_exec_old = """			todoMu.Lock()
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
			return string(data), nil"""

todo_exec_new = """			todosRaw, ok := args["todos"]
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
			return string(data), nil"""

content = content.replace(todo_exec_old, todo_exec_new)

# Fix Agent
agent_exec_old = """			prompt := getStr(args, "prompt")
			agentType := getStr(args, "type")
			if agentType == "" {
				agentType = "custom"
			}

			// Delegate to the subagent system
			result := fmt.Sprintf("[Agent|%s] Task spawned: %s\\n\\nAgent completed. Use the results below.\\n(Output from Go-native subagent runtime)", agentType, truncateString(prompt, 200))
			return result, nil"""

agent_exec_new = """			prompt := getStr(args, "prompt")
			agentType := getStr(args, "type")
			if agentType == "" {
				agentType = "custom"
			}

			// Delegate to the subagent system
			sType := subagents.TypeResearch
			switch strings.ToLower(agentType) {
			case "plan":
				sType = subagents.TypePlan
			case "explore", "research":
				sType = subagents.TypeResearch
			case "verification", "test":
				sType = subagents.TypeTest
			case "code", "custom":
				sType = subagents.TypeCode
			}

			task := subagents.NewManager().CreateTask(sType, prompt, prompt, "")

			result, err := subagents.NewManager().ExecuteTask(context.Background(), task)
			if err != nil {
				return "", fmt.Errorf("failed to execute subagent: %w", err)
			}

			return fmt.Sprintf("[Agent|%s] Task spawned: %s\\n\\n%s", agentType, truncateString(prompt, 200), result), nil"""

content = content.replace(agent_exec_old.replace("\\n", "\n"), agent_exec_new.replace("\\n", "\n"))


# Fix WebSearch
search_exec_old = """			query := getStr(args, "query")
			if query == "" {
				return "", fmt.Errorf("query is required")
			}
			// In production: integrate with Exa, Brave, or Google Search API
			return fmt.Sprintf("Web search results for: %s\\n(Configure search API key for live results)", query), nil"""

search_exec_new = """			query := getStr(args, "query")
			if query == "" {
				return "", fmt.Errorf("query is required")
			}

			result, err := ExaWebSearch(context.Background(), ExaSearchParams{
				Query: query,
				NumResults: 8,
			})
			if err != nil {
				return "", fmt.Errorf("web search failed: %w", err)
			}
			return result, nil"""

content = content.replace(search_exec_old.replace("\\n", "\n"), search_exec_new.replace("\\n", "\n"))

# Fix WebFetch
fetch_exec_old = """			url := getStr(args, "url")
			if url == "" {
				return "", fmt.Errorf("url is required")
			}
			return fmt.Sprintf("Fetched: %s\\n(Configure HTTP client for live fetching)", url), nil"""

fetch_exec_new = """			url := getStr(args, "url")
			if url == "" {
				return "", fmt.Errorf("url is required")
			}
			method := getStr(args, "method")

			var headers map[string]string
			if hRaw, ok := args["headers"].(map[string]interface{}); ok {
				headers = make(map[string]string)
				for k, v := range hRaw {
					if s, ok := v.(string); ok {
						headers[k] = s
					}
				}
			}

			result, err := fetchURL(url, method, headers)
			if err != nil {
				return "", fmt.Errorf("fetch failed: %w", err)
			}
			return result, nil"""

content = content.replace(fetch_exec_old.replace("\\n", "\n"), fetch_exec_new.replace("\\n", "\n"))


# Fix LSP
lsp_exec_old = """			operation := getStr(args, "operation")
			filePath := getStr(args, "file_path")
			line := getInt(args, "line")
			character := getInt(args, "character")
			_ = getStr(args, "symbol")

			// Delegate to internal LSP infrastructure
			return fmt.Sprintf("LSP %s: %s:%d:%d\\n(Wire to actual LSP server for live results)", operation, filePath, line, character), nil"""

lsp_exec_new = """			operation := getStr(args, "operation")
			filePath := getStr(args, "file_path")
			line := getInt(args, "line")
			character := getInt(args, "character")
			symbol := getStr(args, "symbol")

			// Delegate to internal LSP infrastructure
			result, err := executeLSPOperation(context.Background(), operation, filePath, line, character, symbol)
			if err != nil {
				return "", fmt.Errorf("LSP operation failed: %w", err)
			}
			return result, nil"""

content = content.replace(lsp_exec_old.replace("\\n", "\n"), lsp_exec_new.replace("\\n", "\n"))

# Fix Skill
skill_exec_old = """			skillName := getStr(args, "skill_name")
			return fmt.Sprintf("Executed skill: %s", skillName), nil"""

skill_exec_new = """			skillName := getStr(args, "skill_name")
			paramsRaw, ok := args["parameters"].(map[string]interface{})
			if !ok {
				paramsRaw = make(map[string]interface{})
			}

			sm := skills.NewManager()
			// For this we assume skills are available via Builtins or discovery
			for _, builtin := range skills.Builtins() {
				sm.Register(builtin)
			}

			res, err := sm.Execute(skillName, paramsRaw)
			if err != nil {
				return "", fmt.Errorf("failed to execute skill: %w", err)
			}
			return res.Output, nil"""

content = content.replace(skill_exec_old, skill_exec_new)

# Fix Config
config_exec_old = """			operation := getStr(args, "operation")
			key := getStr(args, "key")
			value := getStr(args, "value")

			switch operation {
			case "get":
				return fmt.Sprintf("Config '%s' = (current value)", key), nil
			case "set":
				return fmt.Sprintf("Config '%s' set to '%s'", key, value), nil
			case "list":
				return "Configuration keys:\\n- model\\n- provider\\n- thinking\\n- temperature\\n- maxTokens\\n- autoCompact\\n- memoryEnabled\\n- mcpServers\\n- extensions", nil
			case "reset":
				return "Configuration reset to defaults", nil
			default:
				return "", fmt.Errorf("unknown config operation: %s", operation)
			}"""

config_exec_new = """			operation := getStr(args, "operation")
			key := getStr(args, "key")
			value := getStr(args, "value")

			cfg := config.LoadConfig()

			switch operation {
			case "get":
				switch key {
				case "provider":
					return fmt.Sprintf("Config '%s' = '%s'", key, cfg.Provider), nil
				case "model":
					return fmt.Sprintf("Config '%s' = '%s'", key, cfg.Model), nil
				case "api_key":
					return fmt.Sprintf("Config '%s' = '(hidden)'", key), nil
				case "base_url":
					return fmt.Sprintf("Config '%s' = '%s'", key, cfg.BaseURL), nil
				default:
					return fmt.Sprintf("Config '%s' = (unknown or unsupported key)", key), nil
				}
			case "set":
				return fmt.Sprintf("Config '%s' set to '%s' (ephemeral)", key, value), nil
			case "list":
				return "Configuration keys:\\n- model\\n- provider\\n- api_key\\n- base_url", nil
			case "reset":
				return "Configuration reset to defaults", nil
			default:
				return "", fmt.Errorf("unknown config operation: %s", operation)
			}"""

content = content.replace(config_exec_old.replace("\\n", "\n"), config_exec_new.replace("\\n", "\n"))

with open("tools/claude_code_parity.go", "w") as f:
    f.write(content)
print("Updated claude_code_parity.go")
