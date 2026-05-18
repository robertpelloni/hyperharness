package lsp

import (
	"encoding/json"
	"testing"
)

func TestPosition(t *testing.T) {
	pos := Position{Line: 10, Character: 5}
	data, err := json.Marshal(pos)
	if err != nil {
		t.Fatal(err)
	}
	var decoded Position
	json.Unmarshal(data, &decoded)
	if decoded.Line != 10 || decoded.Character != 5 {
		t.Errorf("position roundtrip: %+v", decoded)
	}
}

func TestRange(t *testing.T) {
	r := Range{
		Start: Position{Line: 1, Character: 0},
		End:   Position{Line: 1, Character: 10},
	}
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	var decoded Range
	json.Unmarshal(data, &decoded)
	if decoded.Start.Line != 1 {
		t.Errorf("range start: %+v", decoded.Start)
	}
}

func TestDiagnostic(t *testing.T) {
	diag := Diagnostic{
		Range:    Range{Start: Position{Line: 5}, End: Position{Line: 5, Character: 20}},
		Severity: 1,
		Message:  "syntax error",
		Source:   "gopls",
	}
	data, err := json.Marshal(diag)
	if err != nil {
		t.Fatal(err)
	}
	var decoded Diagnostic
	json.Unmarshal(data, &decoded)
	if decoded.Message != "syntax error" {
		t.Errorf("diagnostic message: %s", decoded.Message)
	}
	if decoded.Severity != 1 {
		t.Errorf("diagnostic severity: %d", decoded.Severity)
	}
}

func TestLocation(t *testing.T) {
	loc := Location{
		URI:   "file:///test.go",
		Range: Range{Start: Position{Line: 1}, End: Position{Line: 1, Character: 10}},
	}
	data, _ := json.Marshal(loc)
	var decoded Location
	json.Unmarshal(data, &decoded)
	if decoded.URI != "file:///test.go" {
		t.Errorf("location URI: %s", decoded.URI)
	}
}

func TestDocumentSymbol(t *testing.T) {
	sym := DocumentSymbol{
		Name:           "MyFunc",
		Kind:           12, // Function
		Range:          Range{Start: Position{Line: 10}, End: Position{Line: 20}},
		SelectionRange: Range{Start: Position{Line: 10, Character: 5}, End: Position{Line: 10, Character: 11}},
	}
	data, _ := json.Marshal(sym)
	var decoded DocumentSymbol
	json.Unmarshal(data, &decoded)
	if decoded.Name != "MyFunc" {
		t.Errorf("symbol name: %s", decoded.Name)
	}
}

func TestCompletionItem(t *testing.T) {
	item := CompletionItem{
		Label:      "fmt.Println",
		Kind:       3, // Function
		Detail:     "func Println(a ...interface{}) (n int, err error)",
		InsertText: "fmt.Println(${1:args})",
	}
	data, _ := json.Marshal(item)
	var decoded CompletionItem
	json.Unmarshal(data, &decoded)
	if decoded.Label != "fmt.Println" {
		t.Errorf("completion label: %s", decoded.Label)
	}
}

func TestTextEdit(t *testing.T) {
	edit := TextEdit{
		Range:   Range{Start: Position{Line: 1}, End: Position{Line: 1, Character: 5}},
		NewText: "replacement",
	}
	data, _ := json.Marshal(edit)
	var decoded TextEdit
	json.Unmarshal(data, &decoded)
	if decoded.NewText != "replacement" {
		t.Errorf("text edit new text: %s", decoded.NewText)
	}
}

func TestWorkspaceEdit(t *testing.T) {
	edit := WorkspaceEdit{
		Changes: map[string][]TextEdit{
			"file:///test.go": {
				{NewText: "hello", Range: Range{Start: Position{Line: 1}, End: Position{Line: 1, Character: 5}}},
			},
		},
	}
	data, _ := json.Marshal(edit)
	var decoded WorkspaceEdit
	json.Unmarshal(data, &decoded)
	if len(decoded.Changes) != 1 {
		t.Errorf("workspace edit changes: %d", len(decoded.Changes))
	}
}

func TestRPCRequest(t *testing.T) {
	req := RPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params:  map[string]any{"processId": 0},
	}
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"processId":0}}` {
		// JSON ordering may vary, just check it's valid
		var decoded RPCRequest
		json.Unmarshal(data, &decoded)
		if decoded.Method != "initialize" {
			t.Errorf("method: %s", decoded.Method)
		}
	}
}

func TestRPCResponse(t *testing.T) {
	resp := RPCResponse{
		JSONRPC: "2.0",
		ID:      1,
		Result:  json.RawMessage(`{"capabilities":{}}`),
	}
	data, _ := json.Marshal(resp)
	var decoded RPCResponse
	json.Unmarshal(data, &decoded)
	if decoded.ID != 1 {
		t.Errorf("response id: %d", decoded.ID)
	}
}

func TestRPCError(t *testing.T) {
	resp := RPCResponse{
		JSONRPC: "2.0",
		ID:      2,
		Error:   &RPCError{Code: -32601, Message: "method not found"},
	}
	data, _ := json.Marshal(resp)
	var decoded RPCResponse
	json.Unmarshal(data, &decoded)
	if decoded.Error == nil {
		t.Fatal("should have error")
	}
	if decoded.Error.Code != -32601 {
		t.Errorf("error code: %d", decoded.Error.Code)
	}
}

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("client should not be nil")
	}
	if client.State() != ClientStopped {
		t.Errorf("initial state should be stopped: %d", client.State())
	}
}

func TestClientConfig(t *testing.T) {
	config := ClientConfig{
		Command: "gopls",
		Args:    []string{"serve"},
		RootURI: "file:///project",
		Env:     map[string]string{"GOFLAGS": "-mod=mod"},
	}
	if config.Command != "gopls" {
		t.Errorf("command: %s", config.Command)
	}
	if len(config.Args) != 1 {
		t.Errorf("args: %v", config.Args)
	}
	if config.RootURI != "file:///project" {
		t.Errorf("root URI: %s", config.RootURI)
	}
}

func TestDefaultCapabilities(t *testing.T) {
	caps := defaultCapabilities()
	if caps == nil {
		t.Error("capabilities should not be nil")
	}
	td, ok := caps["textDocument"]
	if !ok {
		t.Error("should have textDocument capability")
	}
	_ = td
}

func TestFindHeaderEnd(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"Content-Length: 100\r\n\r\n", 23},
		{"Content-Length: 100\r\n", -1},
		{"no header here", -1},
	}

	for _, tc := range tests {
		result := findHeaderEnd([]byte(tc.input))
		if result != tc.expected {
			t.Errorf("findHeaderEnd(%q) = %d, want %d", tc.input, result, tc.expected)
		}
	}
}

func TestOnNotification(t *testing.T) {
	client := NewClient()
	client.OnNotification("textDocument/publishDiagnostics", func(params any) {})

	// Verify handler was registered
	client.mu.Lock()
	_, ok := client.handlers["textDocument/publishDiagnostics"]
	client.mu.Unlock()
	if !ok {
		t.Error("handler should be registered")
	}
}

func TestRPCNotification(t *testing.T) {
	notif := RPCNotification{
		JSONRPC: "2.0",
		Method:  "textDocument/didOpen",
		Params: map[string]any{
			"textDocument": map[string]any{
				"uri":        "file:///test.go",
				"languageId": "go",
			},
		},
	}
	data, _ := json.Marshal(notif)
	var decoded RPCNotification
	json.Unmarshal(data, &decoded)
	if decoded.Method != "textDocument/didOpen" {
		t.Errorf("method: %s", decoded.Method)
	}
}
