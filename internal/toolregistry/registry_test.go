package toolregistry

import "testing"

func TestNewToolRegistry(t *testing.T) {
	tr := NewToolRegistry()
	if tr.Count() != 0 {
		t.Error("new registry should be empty")
	}
}

func TestRegisterAndGet(t *testing.T) {
	tr := NewToolRegistry()
	tr.Register(ToolInfo{Name: "Bash", Description: "Execute bash", Category: "exec"})
	tool, ok := tr.Get("bash")
	if !ok {
		t.Fatal("should find bash")
	}
	if tool.Name != "Bash" {
		t.Errorf("name: %s", tool.Name)
	}
}

func TestCaseInsensitiveGet(t *testing.T) {
	tr := NewToolRegistry()
	tr.Register(ToolInfo{Name: "WebSearch"})
	_, ok := tr.Get("websearch")
	if !ok {
		t.Error("should be case-insensitive")
	}
	_, ok = tr.Get("WEBSEARCH")
	if !ok {
		t.Error("should be case-insensitive")
	}
}

func TestRegisterBatch(t *testing.T) {
	tr := NewToolRegistry()
	err := tr.RegisterBatch([]ToolInfo{
		{Name: "tool1", Description: "First"},
		{Name: "tool2", Description: "Second"},
		{Name: "tool3", Description: "Third"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if tr.Count() != 3 {
		t.Errorf("count: %d", tr.Count())
	}
}

func TestUnregister(t *testing.T) {
	tr := NewToolRegistry()
	tr.Register(ToolInfo{Name: "temp"})
	if !tr.Unregister("temp") {
		t.Error("should exist")
	}
	if tr.Unregister("temp") {
		t.Error("should no longer exist")
	}
}

func TestList(t *testing.T) {
	tr := NewToolRegistry()
	tr.Register(ToolInfo{Name: "Zebra"})
	tr.Register(ToolInfo{Name: "Alpha"})
	tr.Register(ToolInfo{Name: "Middle"})

	list := tr.List()
	if len(list) != 3 {
		t.Fatal("should have 3 tools")
	}
	if list[0].Name != "Alpha" {
		t.Errorf("should be sorted, first: %s", list[0].Name)
	}
}

func TestListByServer(t *testing.T) {
	tr := NewToolRegistry()
	tr.Register(ToolInfo{Name: "a", ServerName: "mcp1"})
	tr.Register(ToolInfo{Name: "b", ServerName: "mcp2"})
	tr.Register(ToolInfo{Name: "c", ServerName: "mcp1"})

	tools := tr.ListByServer("mcp1")
	if len(tools) != 2 {
		t.Errorf("mcp1 tools: %d", len(tools))
	}
}

func TestListAlwaysOn(t *testing.T) {
	tr := NewToolRegistry()
	tr.Register(ToolInfo{Name: "bash", AlwaysOn: true})
	tr.Register(ToolInfo{Name: "rare", AlwaysOn: false})

	tools := tr.ListAlwaysOn()
	if len(tools) != 1 {
		t.Errorf("always-on: %d", len(tools))
	}
}

func TestSearch(t *testing.T) {
	tr := NewToolRegistry()
	tr.Register(ToolInfo{Name: "web_search", Description: "Search the web", Tags: []string{"search"}})
	tr.Register(ToolInfo{Name: "web_fetch", Description: "Fetch web pages"})
	tr.Register(ToolInfo{Name: "bash", Description: "Execute commands"})

	results := tr.Search("web", 10)
	if len(results) != 2 {
		t.Errorf("web results: %d", len(results))
	}

	results = tr.Search("search", 10)
	if len(results) < 1 {
		t.Error("should find search-tagged tool")
	}
}

func TestSetAlwaysOn(t *testing.T) {
	tr := NewToolRegistry()
	tr.Register(ToolInfo{Name: "tool", AlwaysOn: false})
	tr.SetAlwaysOn("tool", true)
	tool, _ := tr.Get("tool")
	if !tool.AlwaysOn {
		t.Error("should be always-on")
	}
}

func TestStats(t *testing.T) {
	tr := NewToolRegistry()
	tr.Register(ToolInfo{Name: "a", Category: "exec", Source: "native"})
	tr.Register(ToolInfo{Name: "b", Category: "exec", Source: "mcp"})
	tr.Register(ToolInfo{Name: "c", Category: "search", Source: "native", AlwaysOn: true})

	stats := tr.Stats()
	if stats.Total != 3 {
		t.Errorf("total: %d", stats.Total)
	}
	if stats.ByCategory["exec"] != 2 {
		t.Errorf("exec category: %d", stats.ByCategory["exec"])
	}
	if stats.BySource["native"] != 2 {
		t.Errorf("native source: %d", stats.BySource["native"])
	}
	if stats.AlwaysOnCount != 1 {
		t.Errorf("always-on: %d", stats.AlwaysOnCount)
	}
}

func TestClear(t *testing.T) {
	tr := NewToolRegistry()
	tr.Register(ToolInfo{Name: "x"})
	tr.Clear()
	if tr.Count() != 0 {
		t.Error("should be empty after clear")
	}
}

func TestEmptyNameError(t *testing.T) {
	tr := NewToolRegistry()
	err := tr.Register(ToolInfo{Name: ""})
	if err == nil {
		t.Error("should reject empty name")
	}
}
