package ai_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const openAIAPIBase = "https://api.openai.com/v1"

// openAIProvider handles all HTTP communication with OpenAI's APIs and any OpenAI-compatible API
// (e.g. xAI's Grok) via a configurable baseURL. It supports both text generation (chat completions)
// and embedding generation. It is a transport layer — model selection and prompt orchestration live
// in compositeProvider.
type openAIProvider struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	// provider identifies which registered Provider this transport instance serves (OpenAI or
	// Grok), so doWithRetry can pick the right shared rate limiter (see ratelimit.go: limiterFor)
	// and MaxRetries setting even though both use the same OpenAI-compatible wire format.
	provider Provider
}

// newOpenAIProvider builds a transport pointed at the canonical OpenAI API.
func newOpenAIProvider(apiKey string, httpClient *http.Client) *openAIProvider {
	return newOpenAICompatibleProvider(openAIAPIBase, apiKey, httpClient)
}

// newOpenAICompatibleProvider builds a transport pointed at any OpenAI-compatible API. Grok reuses
// this with baseURL "https://api.x.ai/v1". baseURL must not include a trailing slash.
func newOpenAICompatibleProvider(baseURL, apiKey string, httpClient *http.Client) *openAIProvider {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	if baseURL == "" {
		baseURL = openAIAPIBase
	}
	return &openAIProvider{apiKey: apiKey, baseURL: baseURL, httpClient: httpClient, provider: providerForBaseURL(baseURL)}
}

// providerForBaseURL maps a transport's base URL back to its registered Provider ID, so
// doWithRetry can select the right rate limiter. Defaults to ProviderOpenAI for any unrecognised
// base URL (e.g. a test server), which is a safe, permissive default.
func providerForBaseURL(baseURL string) Provider {
	if baseURL == grokAPIBase {
		return ProviderGrok
	}
	return ProviderOpenAI
}

// complete sends a chat completion request to the OpenAI-compatible Chat Completions API at the
// provider's baseURL. When systemPrompt is non-empty it is sent as a system message before the
// user message.
func (o *openAIProvider) complete(ctx context.Context, req completionRequest) (string, error) {
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

	buildReq := func() (*http.Request, error) {
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, o.baseURL+"/chat/completions", bytes.NewReader(body))
		if err != nil {
			return nil, fmt.Errorf("openai: chat: create request: %w", err)
		}
		httpReq.Header.Set("Authorization", "Bearer "+o.apiKey)
		httpReq.Header.Set("Content-Type", "application/json")
		return httpReq, nil
	}

	maxRetries := 5
	if descriptor, found := ProviderByID(o.provider); found && descriptor.MaxRetries > 0 {
		maxRetries = descriptor.MaxRetries
	}
	statusCode, respBody, err := doWithRetry(ctx, o.httpClient, o.provider, estimateTokens(req.maxTokens, len(req.systemPrompt)+len(req.userPrompt)),
		buildReq, parseOpenAIRateLimit, retryConfig{maxRetries: maxRetries})
	if err != nil {
		if IsRateLimitError(err) {
			return "", err
		}
		return "", fmt.Errorf("openai: chat: do request: %w", err)
	}

	var result openAIChatResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("openai: chat: unmarshal response (status %d): %w", statusCode, err)
	}

	if result.Error != nil {
		return "", fmt.Errorf("openai API error %s: %s", result.Error.Type, result.Error.Message)
	}

	if statusCode != http.StatusOK {
		return "", fmt.Errorf("openai: chat: unexpected HTTP status %d", statusCode)
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

	buildReq := func() (*http.Request, error) {
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, o.baseURL+"/embeddings", bytes.NewReader(body))
		if err != nil {
			return nil, fmt.Errorf("openai: embeddings: create request: %w", err)
		}
		httpReq.Header.Set("Authorization", "Bearer "+o.apiKey)
		httpReq.Header.Set("Content-Type", "application/json")
		return httpReq, nil
	}

	maxRetries := 5
	if descriptor, found := ProviderByID(o.provider); found && descriptor.MaxRetries > 0 {
		maxRetries = descriptor.MaxRetries
	}
	estimatedTokens := 0
	for _, in := range inputs {
		estimatedTokens += len(in) / 4
	}
	statusCode, respBody, err := doWithRetry(ctx, o.httpClient, o.provider, estimatedTokens,
		buildReq, parseOpenAIRateLimit, retryConfig{maxRetries: maxRetries})
	if err != nil {
		if IsRateLimitError(err) {
			return nil, err
		}
		return nil, fmt.Errorf("openai: embeddings: do request: %w", err)
	}

	var result openAIEmbeddingResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("openai: embeddings: unmarshal response (status %d): %w", statusCode, err)
	}

	if result.Error != nil {
		return nil, fmt.Errorf("openai: embeddings: API error: %s", result.Error.Message)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("openai: embeddings: unexpected HTTP status %d", statusCode)
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
	return validateOpenAICompatibleKey(openAIAPIBase, "OpenAI", apiKey, httpClient)
}

// ValidateGrokAPIKey checks whether apiKey is accepted by xAI's Grok API, which exposes an
// OpenAI-compatible /models endpoint. Pass nil for httpClient to use http.DefaultClient.
func ValidateGrokAPIKey(apiKey string, httpClient *http.Client) error {
	return validateOpenAICompatibleKey(grokAPIBase, "Grok", apiKey, httpClient)
}

// validateOpenAICompatibleKey validates a key against any OpenAI-compatible /models endpoint.
// providerLabel is used only in the user-facing error message.
func validateOpenAICompatibleKey(baseURL, providerLabel, apiKey string, httpClient *http.Client) error {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	req, err := http.NewRequest(http.MethodGet, baseURL+"/models", nil)
	if err != nil {
		return fmt.Errorf("openai-compatible: validate key: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("openai-compatible: validate key: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("invalid %s API key", providerLabel)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("openai-compatible: validate key: unexpected status %d", resp.StatusCode)
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
