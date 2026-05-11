package toolregistry

import (
	"encoding/json"
	"fmt"
)

// ValidateParameters checks if the provided arguments match the JSON schema
// defined in the tool's Parameters field. We specifically validate the "required"
// constraints dictated by the JSON schema.
func ValidateParameters(schemaRaw json.RawMessage, args map[string]interface{}) error {
	if len(schemaRaw) == 0 {
		return nil
	}
	
	type JSONSchema struct {
		Required []string `json:"required"`
	}
	
	var schema JSONSchema
	if err := json.Unmarshal(schemaRaw, &schema); err != nil {
		// Ignore malformed schemas instead of crashing the pipeline
		return nil
	}
	
	for _, req := range schema.Required {
		if _, ok := args[req]; !ok {
			return fmt.Errorf("parameter validation failed: missing required parameter '%s'", req)
		}
	}

	return nil
}
