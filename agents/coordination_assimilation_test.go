package agents

import (
	"context"
	"errors"
	"strings"
	"testing"
)

type recordingProvider struct {
	chatFn    func(ctx context.Context, messages []Message, tools []Tool) (Message, error)
	streamFn  func(ctx context.Context, messages []Message, tools []Tool, chunkChan chan<- string) error
	modelName string
}

func (p *recordingProvider) Chat(ctx context.Context, messages []Message, tools []Tool) (Message, error) {
	if p.chatFn != nil {
		return p.chatFn(ctx, messages, tools)
	}
	return Message{Role: RoleAssistant, Content: "ok"}, nil
}

func (p *recordingProvider) Stream(ctx context.Context, messages []Message, tools []Tool, chunkChan chan<- string) error {
	if p.streamFn != nil {
		return p.streamFn(ctx, messages, tools, chunkChan)
	}
	return nil
}

func (p *recordingProvider) GetModelName() string {
	return p.modelName
}

func TestDirectorHandleInputRejectsNilReceiverAndMissingProvider(t *testing.T) {
	var nilDirector *Director
	if _, err := nilDirector.HandleInput(context.Background(), "hello"); err == nil || !strings.Contains(err.Error(), "director is required") {
		t.Fatalf("expected nil director error, got %v", err)
	}

	d := &Director{}
	if _, err := d.HandleInput(context.Background(), "hello"); err == nil || !strings.Contains(err.Error(), "provider is required") {
		t.Fatalf("expected missing provider error, got %v", err)
	}
}

func TestDirectorInjectSystemContextInitializesHistory(t *testing.T) {
	d := &Director{}
	d.InjectSystemContext("ctx")
	if len(d.History) != 1 || d.History[0].Role != RoleSystem || d.History[0].Content != "ctx" {
		t.Fatalf("expected system history initialization, got %#v", d.History)
	}
}

func TestDirectorHandleInputStoresPlanAndAppendsResponse(t *testing.T) {
	provider := &recordingProvider{
		chatFn: func(ctx context.Context, messages []Message, tools []Tool) (Message, error) {
			if len(messages) < 2 {
				t.Fatalf("expected prior history plus planning context, got %#v", messages)
			}
			return Message{Role: RoleAssistant, Content: "director response"}, nil
		},
	}
	d := NewDirector(provider)
	d.WorkingDir = t.TempDir()

	resp, err := d.HandleInput(context.Background(), "analyze architecture")
	if err != nil {
		t.Fatalf("handle input failed: %v", err)
	}
	if !strings.Contains(resp, "[Director Plan]") || !strings.Contains(resp, "director response") {
		t.Fatalf("unexpected response: %q", resp)
	}
	if _, ok := d.State["lastPlan"]; !ok {
		t.Fatalf("expected lastPlan to be stored, got %#v", d.State)
	}
	if len(d.History) < 3 || d.History[len(d.History)-1].Content != "director response" {
		t.Fatalf("expected assistant response appended to history, got %#v", d.History)
	}
}

func TestDisclosureProxyRequiresBaseProvider(t *testing.T) {
	var proxy *DisclosureProxy
	if _, err := proxy.Chat(context.Background(), nil, nil); err == nil || !strings.Contains(err.Error(), "base provider is required") {
		t.Fatalf("expected missing provider error for nil proxy, got %v", err)
	}

	proxy = &DisclosureProxy{}
	if err := proxy.Stream(context.Background(), nil, nil, nil); err == nil || !strings.Contains(err.Error(), "base provider is required") {
		t.Fatalf("expected missing provider error for stream, got %v", err)
	}
	if name := proxy.GetModelName(); name != "" {
		t.Fatalf("expected empty model name without base provider, got %q", name)
	}
}

func TestDisclosureProxyFormatsMinimalToolsAndForwardsThem(t *testing.T) {
	var seenTools []Tool
	provider := &recordingProvider{
		modelName: "demo-model",
		chatFn: func(ctx context.Context, messages []Message, tools []Tool) (Message, error) {
			seenTools = append([]Tool(nil), tools...)
			return Message{Role: RoleAssistant, Content: "ok"}, nil
		},
	}
	proxy := NewDisclosureProxy(provider, []Tool{{Name: "apply_search_replace"}, {Name: "bash"}})

	msg, err := proxy.Chat(context.Background(), []Message{{Role: RoleUser, Content: "hi"}}, []Tool{{Name: "ignored"}})
	if err != nil {
		t.Fatalf("chat failed: %v", err)
	}
	if msg.Content != "ok" {
		t.Fatalf("unexpected chat response: %#v", msg)
	}
	if len(seenTools) != 3 {
		t.Fatalf("expected 3 visible tools, got %#v", seenTools)
	}
	if seenTools[0].Name != "apply_search_replace" || seenTools[1].Name != "search_tools" || seenTools[2].Name != "auto_call_tool" {
		t.Fatalf("unexpected visible tools: %#v", seenTools)
	}
	if proxy.GetModelName() != "demo-model" {
		t.Fatalf("unexpected model name: %q", proxy.GetModelName())
	}
}

func TestDisclosureProxyPropagatesProviderErrors(t *testing.T) {
	want := errors.New("stream failed")
	provider := &recordingProvider{
		streamFn: func(ctx context.Context, messages []Message, tools []Tool, chunkChan chan<- string) error {
			return want
		},
	}
	proxy := NewDisclosureProxy(provider, nil)
	if err := proxy.Stream(context.Background(), nil, nil, nil); !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}
