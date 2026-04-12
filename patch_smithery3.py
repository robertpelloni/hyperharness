import sys
import re

with open("tools/goose_opencode_kimi_parity.go", "r") as f:
    content = f.read()

bad = """			// Use the extensions manager to install
			extMgr, err := extensions.NewManager("")
			if err != nil {
				return "", fmt.Errorf("failed to init extension manager: %w", err)
			}

			if err := extMgr.InstallFromSmithery(serverName, configRaw); err != nil {"""

good = """			// Use the extensions manager to install
			extMgr := extensions.NewManager("")

			if err := extMgr.InstallFromSmithery(serverName, configRaw); err != nil {"""

content = content.replace(bad, good)

with open("tools/goose_opencode_kimi_parity.go", "w") as f:
    f.write(content)
