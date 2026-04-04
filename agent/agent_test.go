package agent

import "testing"

func TestBuildOpenAIToolsUsesExactSchemas(t *testing.T) {
	agent := NewAgent()
	tools := agent.buildOpenAITools()
	if len(tools) == 0 {
		t.Fatal("expected tools")
	}
	foundRead := false
	for _, tool := range tools {
		if tool.Function == nil {
			continue
		}
		if tool.Function.Name == "read" {
			foundRead = true
			if tool.Function.Parameters == nil {
				t.Fatal("read tool missing parameters")
			}
		}
	}
	if !foundRead {
		t.Fatal("expected read tool in OpenAI tool list")
	}
}
