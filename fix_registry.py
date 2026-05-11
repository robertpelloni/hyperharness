import re

file_path = "internal/toolregistry/registry.go"

with open(file_path, "r") as f:
    content = f.read()

old_call = """// Call is a helper that wraps Execute by JSON unmarshalling string arguments.
func (t *Tool) Call(argsJSON string) (string, error) {
	var args map[string]interface{}
	if argsJSON != "" {
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", err
		}
	}
	return t.Execute(args)
}"""

new_call = """// Call is a helper that wraps Execute by JSON unmarshalling string arguments.
func (t *Tool) Call(argsJSON string) (string, error) {
	var args map[string]interface{}
	if argsJSON != "" {
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", err
		}
	}
	
	// Validate JSON schema
	if err := ValidateParameters(t.Parameters, args); err != nil {
		return "", err
	}
	
	return t.Execute(args)
}"""

content = content.replace(old_call, new_call)

with open(file_path, "w") as f:
    f.write(content)
