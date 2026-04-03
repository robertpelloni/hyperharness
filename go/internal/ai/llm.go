package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type LLMResponse struct {
	Content string
	Usage   struct {
		InputTokens  int
		OutputTokens int
	}
	Provider string
	Model    string
}

type Provider interface {
	GenerateText(ctx context.Context, model string, messages []Message) (*LLMResponse, error)
}

type OpenAIProvider struct {
	APIKey  string
	BaseURL string
}

func (p *OpenAIProvider) GenerateText(ctx context.Context, model string, messages []Message) (*LLMResponse, error) {
	if p.BaseURL == "" {
		p.BaseURL = "https://api.openai.com/v1/chat/completions"
	}

	reqBody, _ := json.Marshal(map[string]interface{}{
		"model":    model,
		"messages": messages,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", p.BaseURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.APIKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OpenAI API error: %s - %s", resp.Status, string(body))
	}

	var payload struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
		} `json:"usage"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	if len(payload.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned from OpenAI")
	}

	return &LLMResponse{
		Content:  payload.Choices[0].Message.Content,
		Provider: "openai",
		Model:    model,
		Usage: struct {
			InputTokens  int
			OutputTokens int
		}{
			InputTokens:  payload.Usage.PromptTokens,
			OutputTokens: payload.Usage.CompletionTokens,
		},
	}, nil
}

type AnthropicProvider struct {
	APIKey  string
	BaseURL string
}

func (p *AnthropicProvider) GenerateText(ctx context.Context, model string, messages []Message) (*LLMResponse, error) {
	if p.BaseURL == "" {
		p.BaseURL = "https://api.anthropic.com/v1/messages"
	}

	reqBody, _ := json.Marshal(map[string]interface{}{
		"model":      model,
		"max_tokens": 4096,
		"messages":   messages,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", p.BaseURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Anthropic API error: %s - %s", resp.Status, string(body))
	}

	var payload struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	if len(payload.Content) == 0 {
		return nil, fmt.Errorf("no content returned from Anthropic")
	}

	return &LLMResponse{
		Content:  payload.Content[0].Text,
		Provider: "anthropic",
		Model:    model,
		Usage: struct {
			InputTokens  int
			OutputTokens int
		}{
			InputTokens:  payload.Usage.InputTokens,
			OutputTokens: payload.Usage.OutputTokens,
		},
	}, nil
}

// AutoRoute selects the best available provider based on environment variables
// This acts as a lightweight fallback router when the main TypeScript Core is unavailable
func AutoRoute(ctx context.Context, messages []Message) (*LLMResponse, error) {
	if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		provider := &AnthropicProvider{APIKey: key}
		return provider.GenerateText(ctx, "claude-3-5-sonnet-20241022", messages)
	}

	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		provider := &OpenAIProvider{APIKey: key}
		return provider.GenerateText(ctx, "gpt-4o", messages)
	}

	return nil, fmt.Errorf("no suitable LLM provider configured in the environment (missing ANTHROPIC_API_KEY or OPENAI_API_KEY)")
}
