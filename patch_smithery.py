import sys

with open("tools/goose_opencode_kimi_parity.go", "r") as f:
    content = f.read()

# Add imports
if '"github.com/robertpelloni/hyperharness/internal/extensions"' not in content:
    content = content.replace(
        '"github.com/robertpelloni/hyperharness/foundation/pi"',
        '"github.com/robertpelloni/hyperharness/foundation/pi"\n\t"github.com/robertpelloni/hyperharness/internal/extensions"'
    )

old_install = """		Execute: func(args map[string]interface{}) (string, error) {
			serverName, _ := args["server_name"].(string)
			if serverName == "" {
				return "", fmt.Errorf("server_name is required")
			}
			return fmt.Sprintf("Smithery: Installing MCP server '%s'... (MCP registry integration required)", serverName), nil
		},"""

new_install = """		Execute: func(args map[string]interface{}) (string, error) {
			serverName, _ := args["server_name"].(string)
			if serverName == "" {
				return "", fmt.Errorf("server_name is required")
			}

			configRaw, _ := args["config"].(map[string]interface{})
			if configRaw == nil {
				configRaw = make(map[string]interface{})
			}

			// Use the extensions manager to install
			extMgr, err := extensions.NewManager("")
			if err != nil {
				return "", fmt.Errorf("failed to init extension manager: %w", err)
			}

			if err := extMgr.InstallFromSmithery(serverName, configRaw); err != nil {
				return "", fmt.Errorf("smithery install failed: %w", err)
			}

			return fmt.Sprintf("Smithery: Successfully installed MCP server '%s'", serverName), nil
		},"""

content = content.replace(old_install, new_install)

with open("tools/goose_opencode_kimi_parity.go", "w") as f:
    f.write(content)
print("Patched smithery tool")
