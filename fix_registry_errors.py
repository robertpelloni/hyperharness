import sys
import re

with open("tools/registry.go", "r") as f:
    content = f.read()

# Delete duplicate GetInt in registry.go
pattern = """// GetInt parses a generic JSON number or float into an int.
func GetInt(val interface{}, defaultVal int) int {
	if val == nil {
		return defaultVal
	}
	switch v := val.(type) {
	case float64:
		return int(v)
	case json.Number:
		if i, err := v.Int64(); err == nil {
			return int(i)
		}
	}
	return defaultVal
}"""

content = content.replace(pattern, "")

with open("tools/registry.go", "w") as f:
    f.write(content)

with open("tools/goose_opencode_kimi_parity.go", "r") as f:
    content = f.read()

content = content.replace('GetInt(end, startLine+100)', 'GetIntDef(end, startLine+100)')

with open("tools/goose_opencode_kimi_parity.go", "w") as f:
    f.write(content)
