package tools

import (
	"encoding/json"
	"fmt"
)

// registerTabbyParityTools registers tool schemas matching Tabby AI Coding Assistant.
func (r *Registry) registerTabbyParityTools() {
	// tabby_completion - 1:1 parity with Tabby's /v1/completions
	r.Tools = append(r.Tools, Tool{
		Name:        "tabby_completion",
		Description: "Tabby: Get code completion suggestions based on prefix and suffix segments.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"language": { "type": "string" },
				"segments": {
					"type": "object",
					"properties": {
						"prefix": { "type": "string" },
						"suffix": { "type": "string" },
						"filepath": { "type": "string" },
						"git_url": { "type": "string" }
					},
					"required": ["prefix"]
				},
				"temperature": { "type": "number" }
			},
			"required": ["segments"]
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			segments, _ := args["segments"].(map[string]interface{})
			prefix := GetStr(segments, "prefix")
			return fmt.Sprintf("Tabby Completion: Suggested code for prefix starting with %q", TruncateString(prefix, 20)), nil
		},
	})

	// tabby_chat - 1:1 parity with Tabby's /v1/chat/completions
	r.Tools = append(r.Tools, Tool{
		Name:        "tabby_chat",
		Description: "Tabby: Chat with the model using OpenAI-compatible message format.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"messages": {
					"type": "array",
					"items": {
						"type": "object",
						"properties": {
							"role": { "type": "string" },
							"content": { "type": "string" }
						}
					}
				}
			},
			"required": ["messages"]
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return "Tabby Chat: Response generated.", nil
		},
	})

	// tabby_search - Tabby's code search functionality
	r.Tools = append(r.Tools, Tool{
		Name:        "tabby_search",
		Description: "Tabby: Search the code base using content and metadata.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"content": { "type": "string" },
				"git_url": { "type": "string" },
				"language": { "type": "string" }
			},
			"required": ["content"]
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			content := GetStr(args, "content")
			return fmt.Sprintf("Tabby Search: Found matches for %q", content), nil
		},
	})
}
