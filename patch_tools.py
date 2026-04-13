import sys

with open("tools/crush_parity.go", "r") as f:
    content = f.read()

bad = """			workingDir, _ := args["working_dir"].(string)
			runInBackground, _ := args["run_in_background"].(bool)
			description, _ := args["description"].(string)
			timeoutSeconds := GetIntDef(args["timeout"], 0)"""

good = """			workingDir, _ := args["working_dir"].(string)
			runInBackground := GetBool(args, "run_in_background")
			description, _ := args["description"].(string)
			timeoutSeconds := GetIntDef(args["timeout"], 0)"""

content = content.replace(bad, good)

with open("tools/crush_parity.go", "w") as f:
    f.write(content)
