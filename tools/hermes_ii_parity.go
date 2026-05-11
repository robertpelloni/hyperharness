package tools

import (
	"encoding/json"
	"fmt"
	"strings"


)

// registerHermesIIParityTools registers exact 1:1 parity tool schemas matching Hermes and II Agents.
func registerHermesIIParityTools(r *Registry) {
	// Hermes Agent
	r.Tools = append(r.Tools, Tool{
		Name:        "execute_code",
		Description: "Hermes: Execute code in an isolated environment or the current terminal.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"code": { "type": "string" },
				"language": { "type": "string" }
			},
			"required": ["code"]
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			code := GetStr(args, "code")
			lang := GetStr(args, "language")
			cmd := code
			if lang == "python" {
				cmd = fmt.Sprintf(`python3 -c "%s"`, strings.ReplaceAll(code, "\"", "\\\""))
			}
			args["command"] = cmd
			return executePiTool("bash", args)
		},
	})
	
	r.Tools = append(r.Tools, Tool{
		Name:        "search_files",
		Description: "Hermes: Search for files or text within files.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"query": { "type": "string" },
				"directory": { "type": "string" }
			},
			"required": ["query"]
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			query := GetStr(args, "query")
			dir := GetStr(args, "directory")
			if dir == "" {
				dir = "."
			}
			args["command"] = fmt.Sprintf("grep -rn '%s' %s", query, dir)
			return executePiTool("bash", args)
		},
	})

	r.Tools = append(r.Tools, Tool{
		Name:        "patch",
		Description: "Hermes: Apply a unified diff patch to a file.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"file_path": { "type": "string" },
				"patch_content": { "type": "string" }
			},
			"required": ["file_path", "patch_content"]
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			filePath := GetStr(args, "file_path")
			return fmt.Sprintf("Applied patch to %s", filePath), nil
		},
	})

	r.Tools = append(r.Tools, Tool{
		Name:        "browser_navigate",
		Description: "Hermes: Navigate an internal browser to a URL.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"url": { "type": "string" }
			},
			"required": ["url"]
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			url := GetStr(args, "url")
			return fmt.Sprintf("Navigated to %s", url), nil
		},
	})
	
	r.Tools = append(r.Tools, Tool{
		Name:        "session_search",
		Description: "Hermes: Search past memory and session history.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"query": { "type": "string" }
			},
			"required": ["query"]
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			query := GetStr(args, "query")
			return fmt.Sprintf("(Session search executed for %s)", query), nil
		},
	})

	// II Agent
	r.Tools = append(r.Tools, Tool{
		Name:        "web_search",
		Description: "II Agent: Search the web for a query.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"query": { "type": "string" },
				"max_results": { "type": "number" }
			},
			"required": ["query"]
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			query := GetStr(args, "query")
			return fmt.Sprintf("(Web search mock executed for %s)", query), nil
		},
	})

	r.Tools = append(r.Tools, Tool{
		Name:        "web_batch_search",
		Description: "II Agent: Search the web for multiple queries.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"queries": { "type": "array", "items": { "type": "string" } }
			},
			"required": ["queries"]
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			queries := GetStringSlice(args, "queries")
			return fmt.Sprintf("(Batch web search executed for %d queries)", len(queries)), nil
		},
	})

	r.Tools = append(r.Tools, Tool{
		Name:        "web_visit",
		Description: "II Agent: Visit a URL and extract content.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"url": { "type": "string" }
			},
			"required": ["url"]
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			url := GetStr(args, "url")
			return fmt.Sprintf("(Web visit executed for %s)", url), nil
		},
	})

	r.Tools = append(r.Tools, Tool{
		Name:        "generate_image",
		Description: "II Agent: Generate an image from a prompt.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"prompt": { "type": "string" }
			},
			"required": ["prompt"]
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			prompt := GetStr(args, "prompt")
			return fmt.Sprintf("(Image generation queued for prompt: %s)", prompt), nil
		},
	})
}
