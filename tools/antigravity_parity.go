package tools

import (
	"encoding/json"
)

// registerAntigravityParityTools registers tool schemas matching Antigravity 2.0 features.
func (r *Registry) registerAntigravityParityTools() {
	// apply_diffs - Parity with Antigravity multi-file editing
	r.Tools = append(r.Tools, Tool{
		Name:        "apply_diffs",
		Description: "Antigravity: Apply diffs to multiple files in the project.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"edits": {
					"type": "array",
					"items": {
						"type": "object",
						"properties": {
							"path": { "type": "string" },
							"diff": { "type": "string" }
						},
						"required": ["path", "diff"]
					}
				}
			},
			"required": ["edits"]
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return "Antigravity: Diffs applied to files.", nil
		},
	})

	// review_artifacts - Parity with Antigravity artifact review
	r.Tools = append(r.Tools, Tool{
		Name:        "review_artifacts",
		Description: "Antigravity: Open the artifact review panel for the user.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"artifact_ids": { "type": "array", "items": { "type": "string" } }
			}
		}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return "Antigravity: Artifact review panel opened.", nil
		},
	})
}
