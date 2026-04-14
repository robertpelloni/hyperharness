import sys

with open("tools/claude_code_parity.go", "r") as f:
    content = f.read()

# Fix TodoWrite
todo_pattern = """// ---- Claude Code: TodoWrite ----
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

todo_replacement = """// ---- Claude Code: TodoWrite ----
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

content = content.replace(todo_pattern, todo_replacement)

agent_pattern = """	// Agent - Claude Code's subagent spawning tool (renamed from "Task")
	r.Tools = append(r.Tools, Tool{
		Name:        "Agent",
		Description: "Launch a new agent that has access to the following tools: Read, Write, Edit, Bash, Glob, Grep, WebSearch, WebFetch, TodoWrite, Skill, ToolSearch, LSP. The agent will work on the given task and return a result. Use this for multi-step research, exploration, or implementation tasks.",
		Parameters:  json.RawMessage(`{"type":"object","properties":{"prompt":{"type":"string","description":"The task for the agent to perform"},"type":{"type":"string","description":"Agent type: Explore, Plan, or custom","enum":["Explore","Plan","verification","custom"]}},"required":["prompt"]}`),
		Execute: func(args map[string]interface{}) (string, error) {
			prompt := getStr(args, "prompt")
			agentType := getStr(args, "type")
			if agentType == "" {
				agentType = "custom"
			}

			// Delegate to the subagent system
			result := fmt.Sprintf("[Agent|%s] Task spawned: %s\\n\\nAgent completed. Use the results below.\\n(Output from Go-native subagent runtime)", agentType, truncateString(prompt, 200))
			return result, nil
		},
	})"""

agent_replacement = """	// Agent - Claude Code's subagent spawning tool (renamed from "Task")
	r.Tools = append(r.Tools, Tool{
		Name:        "Agent",
		Description: "Launch a new agent that has access to the following tools: Read, Write, Edit, Bash, Glob, Grep, WebSearch, WebFetch, TodoWrite, Skill, ToolSearch, LSP. The agent will work on the given task and return a result. Use this for multi-step research, exploration, or implementation tasks.",
		Parameters:  json.RawMessage(`{"type":"object","properties":{"prompt":{"type":"string","description":"The task for the agent to perform"},"type":{"type":"string","description":"Agent type: Explore, Plan, or custom","enum":["Explore","Plan","verification","custom"]}},"required":["prompt"]}`),
		Execute: func(args map[string]interface{}) (string, error) {
			prompt := getStr(args, "prompt")
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

			return fmt.Sprintf("[Agent|%s] Task spawned: %s\\n\\n%s", agentType, truncateString(prompt, 200), result), nil
		},
	})"""

content = content.replace(agent_pattern, agent_replacement)

search_pattern = """	// WebSearch - Claude Code's web search tool
	r.Tools = append(r.Tools, Tool{
		Name:        "WebSearch",
		Description: "Search the web for information. Returns search results with titles, URLs, and snippets.",
		Parameters:  json.RawMessage(`{"type":"object","properties":{"query":{"type":"string","description":"The search query"},"allowed_domains":{"type":"array","items":{"type":"string"},"description":"Only include results from these domains"},"blocked_domains":{"type":"array","items":{"type":"string"},"description":"Exclude results from these domains"}},"required":["query"]}`),
		Execute: func(args map[string]interface{}) (string, error) {
			query := getStr(args, "query")
			if query == "" {
				return "", fmt.Errorf("query is required")
			}
			// In production: integrate with Exa, Brave, or Google Search API
			return fmt.Sprintf("Web search results for: %s\\n(Configure search API key for live results)", query), nil
		},
	})"""

search_replacement = """	// WebSearch - Claude Code's web search tool
	r.Tools = append(r.Tools, Tool{
		Name:        "WebSearch",
		Description: "Search the web for information. Returns search results with titles, URLs, and snippets.",
		Parameters:  json.RawMessage(`{"type":"object","properties":{"query":{"type":"string","description":"The search query"},"allowed_domains":{"type":"array","items":{"type":"string"},"description":"Only include results from these domains"},"blocked_domains":{"type":"array","items":{"type":"string"},"description":"Exclude results from these domains"}},"required":["query"]}`),
		Execute: func(args map[string]interface{}) (string, error) {
			query := getStr(args, "query")
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
			return result, nil
		},
	})"""

content = content.replace(search_pattern, search_replacement)

fetch_pattern = """	// WebFetch - Claude Code's URL fetcher
	r.Tools = append(r.Tools, Tool{
		Name:        "WebFetch",
		Description: "Fetch content from a URL. Returns the page content as text.",
		Parameters:  json.RawMessage(`{"type":"object","properties":{"url":{"type":"string","description":"The URL to fetch"},"method":{"type":"string","description":"HTTP method","enum":["GET","POST"],"default":"GET"},"headers":{"type":"object","description":"HTTP headers"}},"required":["url"]}`),
		Execute: func(args map[string]interface{}) (string, error) {
			url := getStr(args, "url")
			if url == "" {
				return "", fmt.Errorf("url is required")
			}
			return fmt.Sprintf("Fetched: %s\\n(Configure HTTP client for live fetching)", url), nil
		},
	})"""

fetch_replacement = """	// WebFetch - Claude Code's URL fetcher
	r.Tools = append(r.Tools, Tool{
		Name:        "WebFetch",
		Description: "Fetch content from a URL. Returns the page content as text.",
		Parameters:  json.RawMessage(`{"type":"object","properties":{"url":{"type":"string","description":"The URL to fetch"},"method":{"type":"string","description":"HTTP method","enum":["GET","POST"],"default":"GET"},"headers":{"type":"object","description":"HTTP headers"}},"required":["url"]}`),
		Execute: func(args map[string]interface{}) (string, error) {
			url := getStr(args, "url")
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
			return result, nil
		},
	})"""

content = content.replace(fetch_pattern, fetch_replacement)

lsp_pattern = """	// LSP - Language Server Protocol tool
	r.Tools = append(r.Tools, Tool{
		Name:        "LSP",
		Description: "Interact with Language Server Protocol (LSP) servers to get code intelligence features. Supported operations: goToDefinition, findReferences, hover, documentSymbol, workspaceSymbol, goToImplementation, diagnostics.",
		Parameters:  json.RawMessage(`{"type":"object","properties":{"operation":{"type":"string","enum":["goToDefinition","findReferences","hover","documentSymbol","workspaceSymbol","goToImplementation","diagnostics"],"description":"The LSP operation to perform"},"file_path":{"type":"string","description":"Path to the file"},"line":{"type":"integer","description":"Line number (0-indexed)"},"character":{"type":"integer","description":"Character offset (0-indexed)"},"symbol":{"type":"string","description":"Symbol name for workspaceSymbol search"}},"required":["operation","file_path"]}`),
		Execute: func(args map[string]interface{}) (string, error) {
			operation := getStr(args, "operation")
			filePath := getStr(args, "file_path")
			line := getInt(args, "line")
			character := getInt(args, "character")
			_ = getStr(args, "symbol")

			// Delegate to internal LSP infrastructure
			return fmt.Sprintf("LSP %s: %s:%d:%d\\n(Wire to actual LSP server for live results)", operation, filePath, line, character), nil
		},
	})"""

lsp_replacement = """	// LSP - Language Server Protocol tool
	r.Tools = append(r.Tools, Tool{
		Name:        "LSP",
		Description: "Interact with Language Server Protocol (LSP) servers to get code intelligence features. Supported operations: goToDefinition, findReferences, hover, documentSymbol, workspaceSymbol, goToImplementation, diagnostics.",
		Parameters:  json.RawMessage(`{"type":"object","properties":{"operation":{"type":"string","enum":["goToDefinition","findReferences","hover","documentSymbol","workspaceSymbol","goToImplementation","diagnostics"],"description":"The LSP operation to perform"},"file_path":{"type":"string","description":"Path to the file"},"line":{"type":"integer","description":"Line number (0-indexed)"},"character":{"type":"integer","description":"Character offset (0-indexed)"},"symbol":{"type":"string","description":"Symbol name for workspaceSymbol search"}},"required":["operation","file_path"]}`),
		Execute: func(args map[string]interface{}) (string, error) {
			operation := getStr(args, "operation")
			filePath := getStr(args, "file_path")
			line := getInt(args, "line")
			character := getInt(args, "character")
			symbol := getStr(args, "symbol")

			// Delegate to internal LSP infrastructure
			result, err := executeLSPOperation(context.Background(), operation, filePath, line, character, symbol)
			if err != nil {
				return "", fmt.Errorf("LSP operation failed: %w", err)
			}
			return result, nil
		},
	})"""

content = content.replace(lsp_pattern, lsp_replacement)

skill_pattern = """	// Skill - Claude Code's skill tool
	r.Tools = append(r.Tools, Tool{
		Name:        "Skill",
		Description: "Invoke a specific skill or capability. Skills are predefined sets of instructions or patterns.",
		Parameters:  json.RawMessage(`{"type":"object","properties":{"skill_name":{"type":"string","description":"The name of the skill to invoke"},"parameters":{"type":"object","description":"Parameters for the skill"}},"required":["skill_name"]}`),
		Execute: func(args map[string]interface{}) (string, error) {
			skillName := getStr(args, "skill_name")
			return fmt.Sprintf("Executed skill: %s", skillName), nil
		},
	})"""

skill_replacement = """	// Skill - Claude Code's skill tool
	r.Tools = append(r.Tools, Tool{
		Name:        "Skill",
		Description: "Invoke a specific skill or capability. Skills are predefined sets of instructions or patterns.",
		Parameters:  json.RawMessage(`{"type":"object","properties":{"skill_name":{"type":"string","description":"The name of the skill to invoke"},"parameters":{"type":"object","description":"Parameters for the skill"}},"required":["skill_name"]}`),
		Execute: func(args map[string]interface{}) (string, error) {
			skillName := getStr(args, "skill_name")
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
			return res.Output, nil
		},
	})"""

content = content.replace(skill_pattern, skill_replacement)

config_pattern = """	// Config - Claude Code's configuration tool
	r.Tools = append(r.Tools, Tool{
		Name:        "Config",
		Description: "Read or modify the current configuration. Supports getting and setting configuration values.",
		Parameters:  json.RawMessage(`{"type":"object","properties":{"operation":{"type":"string","enum":["get","set","list","reset"],"description":"The config operation to perform"},"key":{"type":"string","description":"Configuration key"},"value":{"type":"string","description":"Configuration value (for set operation)"}},"required":["operation"]}`),
		Execute: func(args map[string]interface{}) (string, error) {
			operation := getStr(args, "operation")
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
			}
		},
	})"""

config_replacement = """	// Config - Claude Code's configuration tool
	r.Tools = append(r.Tools, Tool{
		Name:        "Config",
		Description: "Read or modify the current configuration. Supports getting and setting configuration values.",
		Parameters:  json.RawMessage(`{"type":"object","properties":{"operation":{"type":"string","enum":["get","set","list","reset"],"description":"The config operation to perform"},"key":{"type":"string","description":"Configuration key"},"value":{"type":"string","description":"Configuration value (for set operation)"}},"required":["operation"]}`),
		Execute: func(args map[string]interface{}) (string, error) {
			operation := getStr(args, "operation")
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
			}
		},
	})"""

content = content.replace(config_pattern, config_replacement)


if '"github.com/robertpelloni/hyperharness/internal/subagents"' not in content:
    content = content.replace('"github.com/robertpelloni/hyperharness/foundation/pi"',
                              '"github.com/robertpelloni/hyperharness/foundation/pi"\n\t"github.com/robertpelloni/hyperharness/internal/subagents"\n\t"github.com/robertpelloni/hyperharness/internal/skills"\n\t"github.com/robertpelloni/hyperharness/config"\n\t"context"')

with open("tools/claude_code_parity.go", "w") as f:
    f.write(content)
print("Updated claude_code_parity.go")
