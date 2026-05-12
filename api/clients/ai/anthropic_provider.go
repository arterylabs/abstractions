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
	maxTokensSummary    = 800
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

// SummarizeCode produces a concise summary, domain tags, and symbol metadata for a code chunk.
// chunkContent must never be logged.
func (p *AnthropicProvider) SummarizeCode(ctx context.Context, filePath, language, chunkContent string, startLine, endLine int, chunkType string) (*SummarizeResult, error) {
	prompt := buildSummarizePrompt(filePath, language, chunkContent, startLine, endLine, chunkType)
	raw, err := p.complete(ctx, modelHaiku, prompt, maxTokensSummary)
	if err != nil {
		return nil, fmt.Errorf("anthropic: SummarizeCode: %w", err)
	}
	return parseSummarizeResult(raw), nil
}

// SummarizeFile produces a single file-level summary (chunk_index = -1) from the
// concatenated chunk summaries for a file.
func (p *AnthropicProvider) SummarizeFile(ctx context.Context, filePath, language, allChunkSummaries string) (*SummarizeResult, error) {
	prompt := buildSummarizeFilePrompt(filePath, language, allChunkSummaries)
	raw, err := p.complete(ctx, modelHaiku, prompt, maxTokensSummary)
	if err != nil {
		return nil, fmt.Errorf("anthropic: SummarizeFile: %w", err)
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

func buildSummarizePrompt(filePath, language, chunkContent string, startLine, endLine int, chunkType string) string {
	return fmt.Sprintf(`You are a senior software engineer analysing a code chunk.
File: %s
Language: %s
Lines: %d-%d
Chunk type: %s

If this chunk contains only import statements, package declarations, require() calls, or comments with no executable logic, respond with exactly: {"skip": true}

Otherwise, produce a JSON response with exactly four keys:
- "summary": describe (a) what the symbol does in one sentence, then (b) walk through the key internal logic step-by-step — conditions, loops, data transformations, key function calls, error handling, side effects. Be specific and concrete (2-5 sentences total).
- "domain_tags": an array of 1-5 short domain labels (e.g. ["auth","jwt","middleware"]).
- "symbol_name": the primary symbol name this chunk defines (function, class, type, etc.) — empty string if none.
- "symbol_type": the kind of symbol ("function", "class", "interface", "type", "method", "variable", or "").

Code:
<code>
%s
</code>

Respond with valid JSON only. No markdown fences, no extra text.`, filePath, language, startLine, endLine, chunkType, chunkContent)
}

func buildSummarizeFilePrompt(filePath, language, allChunkSummaries string) string {
	return fmt.Sprintf(`You are a senior software engineer writing a high-level file summary.
File: %s
Language: %s

The following are per-chunk summaries for this file:
%s

Produce a JSON response with exactly four keys:
- "summary": a 2-4 sentence plain-English overview of the entire file's purpose and key responsibilities.
- "domain_tags": an array of 1-5 short domain labels that best describe the file.
- "symbol_name": empty string (file-level row).
- "symbol_type": empty string.

Respond with valid JSON only. No markdown fences, no extra text.`, filePath, language, allChunkSummaries)
}

func buildTopicsPrompt(req *TopicGenerationRequest) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Repository: %s\n\n", req.RepoFullName))
	sb.WriteString("File summaries (analyze holistically to understand the codebase's features, modules, and how files connect):\n")
	for path, summary := range req.FileSummaries {
		sb.WriteString(fmt.Sprintf("- %s: %s\n", path, summary))
	}
	if len(req.PreviousTopics) > 0 {
		sb.WriteString("\nPreviously covered topics (do NOT generate duplicates or close variants of these):\n")
		for _, t := range req.PreviousTopics {
			sb.WriteString(fmt.Sprintf("- %s\n", t))
		}
	}
	sb.WriteString(`
First, reason about what this codebase does: identify its key features, modules, recent changes, and how the files connect together.
Then generate 3-5 newsletter topics that would genuinely interest engineers working on this repository. Prefer topics grounded in specific files or code changes.

Respond with a JSON array where each element has:
- "title": short, specific topic title
- "description": 2-3 sentences explaining the change/feature and why it matters
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
Write a focused, technically rich HTML newsletter section for an engineering audience.
Your audience: engineers working on the same codebase who may not be familiar with this specific area.
Assume they are smart but have limited context on this subsystem.

Structure:
- <h2> for the topic title
- Opening paragraph: what this area of the codebase does and why it exists
- Body (2-4 paragraphs): key implementation details, design decisions, interesting patterns, how it connects to the rest of the system
- <code> for inline identifiers, <pre><code> for short code snippets
- Closing paragraph: what engineers outside this area should know or watch out for

Target length: ~1000 words. Be specific and concrete — name actual files, functions, and patterns.
Return valid HTML only. No markdown.`)
	return sb.String()
}

// --- Response parsers ---

type summarizeJSON struct {
	Summary    string   `json:"summary"`
	DomainTags []string `json:"domain_tags"`
	SymbolName string   `json:"symbol_name"`
	SymbolType string   `json:"symbol_type"`
	Skip       bool     `json:"skip"`
}

func parseSummarizeResult(raw string) *SummarizeResult {
	// Strip markdown code fences if present
	cleaned := strings.TrimSpace(raw)
	if strings.HasPrefix(cleaned, "```") {
		lines := strings.SplitN(cleaned, "\n", 2)
		if len(lines) == 2 {
			cleaned = lines[1]
		}
		if idx := strings.LastIndex(cleaned, "```"); idx != -1 {
			cleaned = cleaned[:idx]
		}
		cleaned = strings.TrimSpace(cleaned)
	}

	var parsed summarizeJSON
	if err := json.Unmarshal([]byte(cleaned), &parsed); err != nil {
		return &SummarizeResult{Summary: strings.TrimSpace(raw)}
	}
	// Model signalled this chunk should be skipped (e.g. import-only block)
	if parsed.Skip {
		return nil
	}
	return &SummarizeResult{
		Summary:    parsed.Summary,
		DomainTags: parsed.DomainTags,
		SymbolName: parsed.SymbolName,
		SymbolType: parsed.SymbolType,
	}
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
