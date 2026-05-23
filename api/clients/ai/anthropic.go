package ai_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	anthropicAPIBase = "https://api.anthropic.com/v1"
	anthropicVersion = "2023-06-01"
)

// anthropicProvider handles all HTTP communication with the Anthropic Messages API.
// It is a transport layer — model selection and prompt orchestration live in compositeProvider.
type anthropicProvider struct {
	apiKey     string
	httpClient *http.Client
}

func newAnthropicProvider(apiKey string, httpClient *http.Client) *anthropicProvider {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &anthropicProvider{apiKey: apiKey, httpClient: httpClient}
}

// complete sends a text completion request to the Anthropic Messages API. When systemPrompt is
// non-empty it is attached as a cacheable system block using the prompt-caching beta, which can
// cut input-token costs by ~90% on repeated calls with the same static instructions.
func (a *anthropicProvider) complete(ctx context.Context, req completionRequest) (string, error) {
	payload := anthropicAPIRequest{
		Model:     req.model,
		MaxTokens: req.maxTokens,
		Messages:  []anthropicMessage{{Role: "user", Content: req.userPrompt}},
	}
	if req.systemPrompt != "" {
		payload.System = []anthropicSystemBlock{{
			Type:         "text",
			Text:         req.systemPrompt,
			CacheControl: &anthropicCacheControl{Type: "ephemeral"},
		}}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("anthropic: marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, anthropicAPIBase+"/messages", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("anthropic: create request: %w", err)
	}
	httpReq.Header.Set("x-api-key", a.apiKey)
	httpReq.Header.Set("anthropic-version", anthropicVersion)
	httpReq.Header.Set("content-type", "application/json")
	if req.systemPrompt != "" {
		httpReq.Header.Set("anthropic-beta", "prompt-caching-2024-07-31")
	}

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("anthropic: do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("anthropic: read response: %w", err)
	}

	var result anthropicAPIResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("anthropic: unmarshal response (status %d): %w", resp.StatusCode, err)
	}

	if result.Error != nil {
		return "", fmt.Errorf("anthropic API error %s: %s", result.Error.Type, result.Error.Message)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("anthropic: unexpected HTTP status %d", resp.StatusCode)
	}

	if result.Usage != nil {
		log.Printf("anthropic: tokens model=%s input=%d output=%d cache_create=%d cache_read=%d",
			req.model, result.Usage.InputTokens, result.Usage.OutputTokens,
			result.Usage.CacheCreationInputTokens, result.Usage.CacheReadInputTokens)
	}

	for _, contentBlock := range result.Content {
		if contentBlock.Type == "text" {
			return contentBlock.Text, nil
		}
	}
	return "", fmt.Errorf("anthropic: no text content in response")
}

// ValidateAnthropicAPIKey checks whether apiKey is accepted by Anthropic's API.
// Pass nil for httpClient to use http.DefaultClient.
func ValidateAnthropicAPIKey(apiKey string, httpClient *http.Client) error {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	req, err := http.NewRequest(http.MethodGet, anthropicAPIBase+"/models", nil)
	if err != nil {
		return fmt.Errorf("anthropic: validate key: %w", err)
	}
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", anthropicVersion)
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("anthropic: validate key: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("invalid Anthropic API key")
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("anthropic: validate key: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// --- Anthropic HTTP types ---

type anthropicSystemBlock struct {
	Type         string                 `json:"type"`
	Text         string                 `json:"text"`
	CacheControl *anthropicCacheControl `json:"cache_control,omitempty"`
}

type anthropicCacheControl struct {
	Type string `json:"type"`
}

type anthropicAPIRequest struct {
	Model     string                 `json:"model"`
	MaxTokens int                    `json:"max_tokens"`
	System    []anthropicSystemBlock `json:"system,omitempty"`
	Messages  []anthropicMessage     `json:"messages"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicAPIResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Usage *struct {
		InputTokens              int `json:"input_tokens"`
		OutputTokens             int `json:"output_tokens"`
		CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
		CacheReadInputTokens     int `json:"cache_read_input_tokens"`
	} `json:"usage"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}
