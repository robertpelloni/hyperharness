import sys
import re
import os
import glob

def delete_redundant_funcs(filepath):
    with open(filepath, "r") as f:
        lines = f.readlines()

    out_lines = []
    skip = False
    brace_count = 0

    for line in lines:
        if line.startswith("func GetStr(m map[string]interface{}, key string) string {") and "tools/helpers.go" not in filepath:
            skip = True
            brace_count = 1
            continue
        elif line.startswith("func GetInt(m map[string]interface{}, key string) int {") and "tools/helpers.go" not in filepath:
            skip = True
            brace_count = 1
            continue
        elif line.startswith("func GetInt(val interface{}, def int) int {") and "tools/helpers.go" not in filepath:
            skip = True
            brace_count = 1
            continue
        elif line.startswith("func TruncateString(s string, maxLen int) string {") and "tools/helpers.go" not in filepath:
            skip = True
            brace_count = 1
            continue
        elif line.startswith("// GetStr") or line.startswith("// GetInt") or line.startswith("// TruncateString") or line.startswith("// toInt"):
            # skip comment immediately preceding
            if not skip and "tools/helpers.go" not in filepath:
                # We'll just filter it
                pass

        if skip:
            brace_count += line.count('{')
            brace_count -= line.count('}')
            if brace_count <= 0:
                skip = False
            continue

        out_lines.append(line)

    with open(filepath, "w") as f:
        f.writelines(out_lines)

for f in glob.glob("tools/*.go"):
    if f != "tools/helpers.go":
        delete_redundant_funcs(f)
