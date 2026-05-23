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

const openAIAPIBase = "https://api.openai.com/v1"

// openAIProvider handles all HTTP communication with OpenAI's APIs.
// It supports both text generation (chat completions) and embedding generation.
// It is a transport layer — model selection and prompt orchestration live in compositeProvider.
type openAIProvider struct {
	apiKey     string
	httpClient *http.Client
}

func newOpenAIProvider(apiKey string, httpClient *http.Client) *openAIProvider {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &openAIProvider{apiKey: apiKey, httpClient: httpClient}
}

// completeText sends a chat completion request to OpenAI's Chat Completions API.
// When systemPrompt is non-empty it is sent as a system message before the user message.
func (o *openAIProvider) completeText(ctx context.Context, req completionRequest) (string, error) {
	messages := make([]openAIChatMessage, 0, 2)
	if req.systemPrompt != "" {
		messages = append(messages, openAIChatMessage{Role: "system", Content: req.systemPrompt})
	}
	messages = append(messages, openAIChatMessage{Role: "user", Content: req.userPrompt})

	payload := openAIChatRequest{
		Model:     req.model,
		Messages:  messages,
		MaxTokens: req.maxTokens,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("openai: chat: marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, openAIAPIBase+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("openai: chat: create request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+o.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := o.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("openai: chat: do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("openai: chat: read response: %w", err)
	}

	var result openAIChatResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("openai: chat: unmarshal response (status %d): %w", resp.StatusCode, err)
	}

	if result.Error != nil {
		return "", fmt.Errorf("openai API error %s: %s", result.Error.Type, result.Error.Message)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("openai: chat: unexpected HTTP status %d", resp.StatusCode)
	}

	if result.Usage != nil {
		log.Printf("openai: tokens model=%s prompt=%d completion=%d total=%d",
			req.model, result.Usage.PromptTokens, result.Usage.CompletionTokens, result.Usage.TotalTokens)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("openai: chat: no choices in response")
	}
	return result.Choices[0].Message.Content, nil
}

// embed sends a batch embedding request to OpenAI's Embeddings API.
func (o *openAIProvider) embed(ctx context.Context, model string, inputs []string) ([][]float32, error) {
	payload := openAIEmbeddingRequest{
		Model: model,
		Input: inputs,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("openai: embeddings: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, openAIAPIBase+"/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("openai: embeddings: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+o.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("openai: embeddings: do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("openai: embeddings: read response: %w", err)
	}

	var result openAIEmbeddingResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("openai: embeddings: unmarshal response (status %d): %w", resp.StatusCode, err)
	}

	if result.Error != nil {
		return nil, fmt.Errorf("openai: embeddings: API error: %s", result.Error.Message)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openai: embeddings: unexpected HTTP status %d", resp.StatusCode)
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("openai: embeddings: no data in response")
	}

	// Map by index to preserve input ordering (OpenAI guarantees order but we verify defensively).
	out := make([][]float32, len(inputs))
	for _, d := range result.Data {
		if d.Index >= 0 && d.Index < len(out) {
			out[d.Index] = d.Embedding
		}
	}
	return out, nil
}

// ValidateOpenAIAPIKey checks whether apiKey is accepted by OpenAI's API.
// Pass nil for httpClient to use http.DefaultClient.
func ValidateOpenAIAPIKey(apiKey string, httpClient *http.Client) error {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	req, err := http.NewRequest(http.MethodGet, openAIAPIBase+"/models", nil)
	if err != nil {
		return fmt.Errorf("openai: validate key: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("openai: validate key: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("invalid OpenAI API key")
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("openai: validate key: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// --- OpenAI HTTP types ---

type openAIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIChatRequest struct {
	Model     string              `json:"model"`
	Messages  []openAIChatMessage `json:"messages"`
	MaxTokens int                 `json:"max_tokens"`
}

type openAIChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage *struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

type openAIEmbeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type openAIEmbeddingResponse struct {
	Data []struct {
		Index     int       `json:"index"`
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}
