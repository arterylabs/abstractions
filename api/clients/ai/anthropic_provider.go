package ai_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	anthropicAPIBase    = "https://api.anthropic.com/v1"
	anthropicVersion    = "2023-06-01"
	modelHaiku          = "claude-3-haiku-20240307"
	modelSonnet         = "claude-3-5-sonnet-20241022"
	maxTokensSummary    = 512
	maxTokensTopics     = 2048
	maxTokensNewsletter = 4096
)

// AnthropicProvider is the Anthropic implementation of AIProvider.
type AnthropicProvider struct {
	apiKey     string
	httpClient *http.Client
}

// NewAnthropicProvider creates a new Anthropic provider.
func NewAnthropicProvider(apiKey string) *AnthropicProvider {
	return &AnthropicProvider{apiKey: apiKey, httpClient: &http.Client{}}
}

// NewAnthropicProviderWithHTTPClient creates an Anthropic provider with a custom HTTP client (used in tests).
func NewAnthropicProviderWithHTTPClient(apiKey string, client *http.Client) *AnthropicProvider {
	return &AnthropicProvider{apiKey: apiKey, httpClient: client}
}

// SummarizeCode produces a concise summary and domain tags for a code chunk.
// chunkContent must never be logged.
func (p *AnthropicProvider) SummarizeCode(ctx context.Context, filePath, language, chunkContent string) (*SummarizeResult, error) {
	prompt := buildSummarizePrompt(filePath, language, chunkContent)
	raw, err := p.complete(ctx, modelHaiku, prompt, maxTokensSummary)
	if err != nil {
		return nil, fmt.Errorf("anthropic: SummarizeCode: %w", err)
	}
	return parseSummarizeResult(raw), nil
}

// GenerateTopics analyses per-file summaries and proposes newsletter topics.
func (p *AnthropicProvider) GenerateTopics(ctx context.Context, req *TopicGenerationRequest) ([]Topic, error) {
	prompt := buildTopicsPrompt(req)
	raw, err := p.complete(ctx, modelSonnet, prompt, maxTokensTopics)
	if err != nil {
		return nil, fmt.Errorf("anthropic: GenerateTopics: %w", err)
	}
	return parseTopics(raw), nil
}

// GenerateNewsletterEntry writes the HTML content for one newsletter topic.
func (p *AnthropicProvider) GenerateNewsletterEntry(ctx context.Context, req *NewsletterEntryRequest) (*NewsletterEntry, error) {
	prompt := buildNewsletterEntryPrompt(req)
	raw, err := p.complete(ctx, modelSonnet, prompt, maxTokensNewsletter)
	if err != nil {
		return nil, fmt.Errorf("anthropic: GenerateNewsletterEntry: %w", err)
	}
	return &NewsletterEntry{Title: req.Topic.Title, Content: raw}, nil
}

// GenerateEmbedding returns a 1536-dimensional embedding vector.
// This is a stub implementation — replace with a real endpoint call when available.
func (p *AnthropicProvider) GenerateEmbedding(_ context.Context, _ string) ([]float32, error) {
	return make([]float32, 1536), nil
}

// --- HTTP helpers ---

type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	Messages  []anthropicMessage `json:"messages"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

func (p *AnthropicProvider) complete(ctx context.Context, model, userPrompt string, maxTokens int) (string, error) {
	payload := anthropicRequest{
		Model:     model,
		MaxTokens: maxTokens,
		Messages:  []anthropicMessage{{Role: "user", Content: userPrompt}},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, anthropicAPIBase+"/messages", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", anthropicVersion)
	req.Header.Set("content-type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	var result anthropicResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("unmarshal response (status %d): %w", resp.StatusCode, err)
	}

	if result.Error != nil {
		return "", fmt.Errorf("anthropic API error %s: %s", result.Error.Type, result.Error.Message)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected HTTP status %d", resp.StatusCode)
	}

	for _, c := range result.Content {
		if c.Type == "text" {
			return c.Text, nil
		}
	}
	return "", fmt.Errorf("no text content in response")
}

// --- Prompt builders ---

func buildSummarizePrompt(filePath, language, chunkContent string) string {
	return fmt.Sprintf(`You are a senior software engineer analysing a code chunk.
File: %s
Language: %s

Produce a JSON response with exactly two keys:
- "summary": a 1-3 sentence plain-English explanation of what this code does and why it matters.
- "domain_tags": an array of 1-5 short domain labels (e.g. ["auth","jwt","middleware"]).

Code:
<code>
%s
</code>

Respond with valid JSON only. No markdown, no extra text.`, filePath, language, chunkContent)
}

func buildTopicsPrompt(req *TopicGenerationRequest) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Repository: %s\n\nFile summaries:\n", req.RepoFullName))
	for path, summary := range req.FileSummaries {
		sb.WriteString(fmt.Sprintf("- %s: %s\n", path, summary))
	}
	sb.WriteString(`
Based on these summaries, generate 3-5 newsletter topics that would interest engineers working on this repository.
Respond with a JSON array where each element has:
- "title": short topic title
- "description": 2-3 sentences describing the change/feature
- "related_files": array of relevant file paths from the summaries above

Respond with valid JSON only. No markdown, no extra text.`)
	return sb.String()
}

func buildNewsletterEntryPrompt(req *NewsletterEntryRequest) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Repository: %s\nTopic: %s\nDescription: %s\n\nRelated file summaries:\n",
		req.RepoFullName, req.Topic.Title, req.Topic.Description))
	for path, summary := range req.FileSummaries {
		sb.WriteString(fmt.Sprintf("- %s: %s\n", path, summary))
	}
	sb.WriteString(`
Write a concise, engaging HTML newsletter section for an engineering team.
Use <h2> for the title, <p> for paragraphs, <code> for inline code, <pre><code> for code blocks.
Keep it under 300 words. Be specific. Return valid HTML only, no markdown.`)
	return sb.String()
}

// --- Response parsers ---

type summarizeJSON struct {
	Summary    string   `json:"summary"`
	DomainTags []string `json:"domain_tags"`
}

func parseSummarizeResult(raw string) *SummarizeResult {
	var parsed summarizeJSON
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return &SummarizeResult{Summary: strings.TrimSpace(raw)}
	}
	return &SummarizeResult{Summary: parsed.Summary, DomainTags: parsed.DomainTags}
}

type topicJSON struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	RelatedFiles []string `json:"related_files"`
}

func parseTopics(raw string) []Topic {
	var parsed []topicJSON
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return nil
	}
	topics := make([]Topic, len(parsed))
	for i, t := range parsed {
		topics[i] = Topic{Title: t.Title, Description: t.Description, RelatedFiles: t.RelatedFiles}
	}
	return topics
}
