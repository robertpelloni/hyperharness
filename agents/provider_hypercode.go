package agents

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultHyperCodeBaseURL = "http://127.0.0.1:4000"

type hyperCodePayloadMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type hyperCodeChatResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Content  string `json:"content"`
		Provider string `json:"provider"`
		Model    string `json:"model"`
	} `json:"data"`
	Error string `json:"error"`
}

type HyperCodeControlPlaneProvider struct {
	BaseURL string
}

func NewHyperCodeProvider() *HyperCodeControlPlaneProvider {
	return &HyperCodeControlPlaneProvider{BaseURL: defaultHyperCodeBaseURL}
}

func (p *HyperCodeControlPlaneProvider) Chat(ctx context.Context, messages []Message, tools []Tool) (Message, error) {
	if p == nil {
		return Message{}, fmt.Errorf("provider is required")
	}
	reqBody, err := buildHyperCodeChatRequest(messages)
	if err != nil {
		return Message{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, normalizeHyperCodeBaseURL(p.BaseURL)+"/api/agent/chat", bytes.NewBuffer(reqBody))
	if err != nil {
		return Message{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return Message{}, fmt.Errorf("failed to contact HyperCode Control Plane: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return Message{}, fmt.Errorf("HyperCode API error: %s - %s", resp.Status, string(body))
	}

	var result hyperCodeChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Message{}, fmt.Errorf("failed to parse HyperCode response: %w", err)
	}
	return parseHyperCodeChatResponse(result)
}

func (p *HyperCodeControlPlaneProvider) Stream(ctx context.Context, messages []Message, tools []Tool, chunkChan chan<- string) error {
	if chunkChan == nil {
		return fmt.Errorf("chunkChan is required")
	}
	msg, err := p.Chat(ctx, messages, tools)
	if err != nil {
		return err
	}
	chunkChan <- msg.Content
	close(chunkChan)
	return nil
}

func (p *HyperCodeControlPlaneProvider) GetModelName() string {
	return "hypercode-router-active"
}

func buildHyperCodeChatRequest(messages []Message) ([]byte, error) {
	history := make([]hyperCodePayloadMessage, 0, len(messages))
	for _, msg := range messages {
		history = append(history, hyperCodePayloadMessage{Role: string(msg.Role), Content: msg.Content})
	}
	return json.Marshal(map[string]interface{}{"message": "", "history": history})
}

func parseHyperCodeChatResponse(result hyperCodeChatResponse) (Message, error) {
	if !result.Success {
		return Message{}, fmt.Errorf("HyperCode rejected chat: %s", result.Error)
	}
	return Message{Role: RoleAssistant, Content: result.Data.Content}, nil
}

func normalizeHyperCodeBaseURL(baseURL string) string {
	trimmed := strings.TrimSpace(baseURL)
	if trimmed == "" {
		return defaultHyperCodeBaseURL
	}
	return strings.TrimRight(trimmed, "/")
}
