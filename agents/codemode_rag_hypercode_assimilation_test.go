package agents

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestResolveCodeModeSupportsAliases(t *testing.T) {
	tsExec, err := resolveCodeMode(" TS ")
	if err != nil {
		t.Fatalf("expected ts alias to resolve: %v", err)
	}
	if tsExec.fileName != codeModeTypeScriptFile || tsExec.command != "bun" || len(tsExec.args) != 1 || tsExec.args[0] != "run" {
		t.Fatalf("unexpected ts execution config: %#v", tsExec)
	}

	shellExec, err := resolveCodeMode("bash")
	if err != nil {
		t.Fatalf("expected bash alias to resolve: %v", err)
	}
	if shellExec.fileName != codeModeShellFile || shellExec.command != "pwsh" {
		t.Fatalf("unexpected shell execution config: %#v", shellExec)
	}

	if _, err := resolveCodeMode("python"); err == nil || !strings.Contains(err.Error(), "missing language execution binding") {
		t.Fatalf("expected unsupported language error, got %v", err)
	}
}

func TestRAGContracts(t *testing.T) {
	service := NewDocumentIntakeService()
	if service.vectorStore != defaultRAGVectorStore {
		t.Fatalf("expected default vector store %q, got %q", defaultRAGVectorStore, service.vectorStore)
	}
	if err := service.Ingest("doc.md"); err != nil {
		t.Fatalf("expected ingest to accept filepath, got %v", err)
	}
	if err := service.Ingest("   "); err == nil || !strings.Contains(err.Error(), "filepath is required") {
		t.Fatalf("expected filepath validation error, got %v", err)
	}

	embeddings, err := (&EmbeddingService{}).Compute("hello world")
	if err != nil {
		t.Fatalf("expected embedding compute to succeed, got %v", err)
	}
	if len(embeddings) != len(defaultEmbeddingVector) {
		t.Fatalf("unexpected embedding length: %#v", embeddings)
	}
	if &embeddings[0] == &defaultEmbeddingVector[0] {
		t.Fatal("expected embedding result to be copied, not aliased")
	}
	if _, err := (&EmbeddingService{}).Compute("   "); err == nil || !strings.Contains(err.Error(), "text is required") {
		t.Fatalf("expected empty text error, got %v", err)
	}
}

func TestHyperCodeRequestAndResponseHelpers(t *testing.T) {
	payload, err := buildHyperCodeChatRequest([]Message{{Role: RoleUser, Content: "hello"}})
	if err != nil {
		t.Fatalf("build request failed: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("failed to decode request payload: %v", err)
	}
	if decoded["message"] != "" {
		t.Fatalf("expected empty message field, got %#v", decoded)
	}
	history, ok := decoded["history"].([]any)
	if !ok || len(history) != 1 {
		t.Fatalf("expected one history message, got %#v", decoded)
	}
	msg, err := parseHyperCodeChatResponse(hyperCodeChatResponse{Success: true, Data: struct {
		Content  string `json:"content"`
		Provider string `json:"provider"`
		Model    string `json:"model"`
	}{Content: "ok"}})
	if err != nil || msg.Content != "ok" || msg.Role != RoleAssistant {
		t.Fatalf("unexpected parsed response: %#v %v", msg, err)
	}
	if _, err := parseHyperCodeChatResponse(hyperCodeChatResponse{Success: false, Error: "boom"}); err == nil || !strings.Contains(err.Error(), "HyperCode rejected chat: boom") {
		t.Fatalf("expected rejection error, got %v", err)
	}
	if normalizeHyperCodeBaseURL(" ") != defaultHyperCodeBaseURL {
		t.Fatalf("expected default base URL fallback")
	}
	if normalizeHyperCodeBaseURL("http://x/") != "http://x" {
		t.Fatalf("expected trailing slash trimming")
	}
}

func TestHyperCodeProviderChatAndStream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/agent/chat" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"content":"hello from hypercode"}}`))
	}))
	defer server.Close()

	provider := &HyperCodeControlPlaneProvider{BaseURL: server.URL}
	msg, err := provider.Chat(context.Background(), []Message{{Role: RoleUser, Content: "hi"}}, nil)
	if err != nil {
		t.Fatalf("chat failed: %v", err)
	}
	if msg.Content != "hello from hypercode" {
		t.Fatalf("unexpected chat content: %#v", msg)
	}

	chunks := make(chan string, 1)
	if err := provider.Stream(context.Background(), []Message{{Role: RoleUser, Content: "hi"}}, nil, chunks); err != nil {
		t.Fatalf("stream failed: %v", err)
	}
	chunk, ok := <-chunks
	if !ok || chunk != "hello from hypercode" {
		t.Fatalf("unexpected streamed chunk: %q ok=%v", chunk, ok)
	}
}

func TestHyperCodeProviderValidationAndErrors(t *testing.T) {
	var provider *HyperCodeControlPlaneProvider
	if _, err := provider.Chat(context.Background(), nil, nil); err == nil || !strings.Contains(err.Error(), "provider is required") {
		t.Fatalf("expected nil provider error, got %v", err)
	}

	provider = NewHyperCodeProvider()
	if err := provider.Stream(context.Background(), nil, nil, nil); err == nil || !strings.Contains(err.Error(), "chunkChan is required") {
		t.Fatalf("expected nil chunk channel error, got %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusBadGateway)
	}))
	defer server.Close()
	provider.BaseURL = server.URL
	if _, err := provider.Chat(context.Background(), nil, nil); err == nil || !strings.Contains(err.Error(), "HyperCode API error") {
		t.Fatalf("expected HTTP error, got %v", err)
	}

	badJSONServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`not-json`))
	}))
	defer badJSONServer.Close()
	provider.BaseURL = badJSONServer.URL
	if _, err := provider.Chat(context.Background(), nil, nil); err == nil || !strings.Contains(err.Error(), "failed to parse HyperCode response") {
		t.Fatalf("expected parse error, got %v", err)
	}

	rejectServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"success":false,"error":"denied"}`))
	}))
	defer rejectServer.Close()
	provider.BaseURL = rejectServer.URL
	if _, err := provider.Chat(context.Background(), nil, nil); err == nil || !strings.Contains(err.Error(), "HyperCode rejected chat: denied") {
		t.Fatalf("expected rejection error, got %v", err)
	}

	unreachable := &HyperCodeControlPlaneProvider{BaseURL: "http://127.0.0.1:1"}
	_, err := unreachable.Chat(context.Background(), nil, nil)
	if err == nil || !strings.Contains(err.Error(), "failed to contact HyperCode Control Plane") {
		t.Fatalf("expected wrapped transport error, got %v", err)
	}
}
