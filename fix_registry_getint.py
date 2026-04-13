import sys
import glob
import re

for filepath in glob.glob("tools/*.go"):
    if filepath == "tools/helpers.go":
        continue

    with open(filepath, "r") as f:
        content = f.read()

    # We renamed toInt to GetInt but toInt took two arguments (val, def). So we need to change those specific ones to GetIntDef
    content = re.sub(r'GetInt\((args\[[^\]]+\]),\s*(\d+)\)', r'GetIntDef(\1, \2)', content)

    # Some files might have `GetInt(val, def)` left over
    content = re.sub(r'GetInt\(([^,]+),\s*(\d+)\)', r'GetIntDef(\1, \2)', content)

    with open(filepath, "w") as f:
        f.write(content)
