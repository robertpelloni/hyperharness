import re

file_path = "agent/agent.go"

with open(file_path, "r") as f:
    content = f.read()

import_pattern = r'import \('
if '"github.com/robertpelloni/hyperharness/internal/providers"' not in content:
    content = re.sub(import_pattern, 'import (\n\t"github.com/robertpelloni/hyperharness/internal/providers"', content)

old_client = """	return &Agent{
		client: openai.NewClient(apiKey),"""

new_client = """	config := openai.DefaultConfig(apiKey)
	config.HTTPClient = providers.GetPooledHTTPClient()
	
	return &Agent{
		client: openai.NewClientWithConfig(config),"""

content = content.replace(old_client, new_client)

with open(file_path, "w") as f:
    f.write(content)
