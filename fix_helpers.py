import sys

with open("tools/helpers.go", "r") as f:
    content = f.read()

get_int_def = """// GetIntDef extracts an int from a raw value, using a default if missing or invalid.
func GetIntDef(val interface{}, def int) int {
	if val == nil {
		return def
	}
	switch n := val.(type) {
	case float64:
		return int(n)
	case float32:
		return int(n)
	case int:
		return n
	case int64:
		return int(n)
	case json.Number:
		i, _ := n.Int64()
		return int(i)
	default:
		return def
	}
}

"""
content = content.replace("func GetInt(", get_int_def + "func GetInt(")

with open("tools/helpers.go", "w") as f:
    f.write(content)
