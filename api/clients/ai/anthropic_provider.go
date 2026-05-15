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
	modelHaiku          = "claude-haiku-4-5-20251001"
	modelSonnet         = "claude-sonnet-4-6"
	maxTokensSummary    = 800
	maxTokensTopics     = 2048
	maxTokensNewsletter = 4096

	openAIAPIBase        = "https://api.openai.com/v1"
	openAIEmbeddingModel = "text-embedding-3-small"
)

// AnthropicProvider is the Anthropic implementation of AIProvider.
// Text generation uses Anthropic's Messages API; embeddings use OpenAI's Embeddings API
// because Anthropic does not offer a public embeddings endpoint.
type AnthropicProvider struct {
	apiKey     string
	httpClient *http.Client
	openAIKey  string
}

// NewAnthropicProvider creates a new Anthropic provider.
// openAIKey is required for GenerateEmbedding; all other methods use apiKey (Anthropic).
func NewAnthropicProvider(apiKey, openAIKey string) *AnthropicProvider {
	return &AnthropicProvider{apiKey: apiKey, httpClient: &http.Client{}, openAIKey: openAIKey}
}

// NewAnthropicProviderWithHTTPClient creates an Anthropic provider with a custom HTTP client (used in tests).
func NewAnthropicProviderWithHTTPClient(apiKey, openAIKey string, client *http.Client) *AnthropicProvider {
	return &AnthropicProvider{apiKey: apiKey, httpClient: client, openAIKey: openAIKey}
}

// SummarizeCode produces a concise plain-text summary, domain tags, and symbol metadata for a code chunk.
// startLine/endLine and chunkType are passed to build a more context-aware prompt.
// chunkContent must never be logged.
func (p *AnthropicProvider) SummarizeCode(ctx context.Context, filePath, language, chunkContent string, startLine, endLine int, chunkType string) (*SummarizeResult, error) {
	prompt := buildSummarizePrompt(filePath, language, chunkContent, startLine, endLine, chunkType)
	raw, err := p.complete(ctx, modelHaiku, prompt, maxTokensSummary)
	if err != nil {
		return nil, fmt.Errorf("anthropic: SummarizeCode: %w", err)
	}
	return parseSummarizeResult(raw), nil
}

// SummarizeFile produces a holistic plain-text summary for an entire file using the
// already-generated per-chunk summaries. Raw file content is never passed here.
func (p *AnthropicProvider) SummarizeFile(ctx context.Context, filePath, language string, chunkSummaries []string) (*SummarizeResult, error) {
	prompt := buildSummarizeFilePrompt(filePath, language, chunkSummaries)
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

// GenerateEmbedding returns a vector embedding for the given text using OpenAI's
// text-embedding-3-small model (1536 dimensions). An error is returned if the
// OpenAI API key is not configured.
func (p *AnthropicProvider) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	if p.openAIKey == "" {
		return nil, fmt.Errorf("openai: embedding: API key not configured (set OPENAI_API_KEY)")
	}

	payload := openAIEmbeddingRequest{
		Model: openAIEmbeddingModel,
		Input: text,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("openai: embedding: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, openAIAPIBase+"/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("openai: embedding: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+p.openAIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("openai: embedding: do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("openai: embedding: read response: %w", err)
	}

	var result openAIEmbeddingResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("openai: embedding: unmarshal response (status %d): %w", resp.StatusCode, err)
	}

	if result.Error != nil {
		return nil, fmt.Errorf("openai: embedding: API error: %s", result.Error.Message)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openai: embedding: unexpected HTTP status %d", resp.StatusCode)
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("openai: embedding: no data in response")
	}

	return result.Data[0].Embedding, nil
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

type openAIEmbeddingRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type openAIEmbeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
	Error *struct {
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
	return fmt.Sprintf(`You are a senior software engineer analysing a single code chunk extracted from a source file.
File: %s
Language: %s
Chunk type: %s
Lines: %d-%d

IMPORTANT: If this chunk contains ONLY import statements, package declarations, require calls, module-level comments, or blank lines with no executable logic, respond with exactly: {"skip": true}

Otherwise, your task:
1. Identify the PRIMARY symbol — the top-level function, method, class, interface, variable, or constant this chunk defines. If there are multiple, pick the most significant one.
2. Write a detailed 3-6 sentence plain-English summary that covers ALL of the following:
   - What the symbol is and its single-sentence purpose.
   - Step-by-step logic: the key operations performed in order — conditions checked, loops, branching, important internal function/method calls, and data transformations.
   - Inputs and outputs: what parameters it accepts, what it returns or mutates.
   - Notable behaviour: error handling strategies, side effects, external dependencies called, concurrency concerns, or edge cases handled.
   Do NOT produce a vague high-level description. Describe the actual logic inside the symbol.
3. List 1-5 short domain tags that categorise this code (e.g. ["auth","jwt","middleware"]).

Respond with a single JSON object containing exactly these four keys:
- "symbol_name": the name of the primary symbol as a plain string (e.g. "ValidateToken")
- "symbol_type": one of: function, method, class, interface, struct, variable, constant, type, module, other
- "summary": detailed plain-text summary string (no markdown, no code fences)
- "domain_tags": a JSON array of short lowercase strings

Code:
<code>
%s
</code>

Respond with valid JSON only. Do NOT wrap in markdown code fences. No extra text.`, filePath, language, chunkType, startLine, endLine, chunkContent)
}

func buildSummarizeFilePrompt(filePath, language string, chunkSummaries []string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`You are a senior software engineer summarising an entire source file.
File: %s
Language: %s

Below are summaries of the individual code chunks in this file (from top to bottom):
`, filePath, language))

	for i, s := range chunkSummaries {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, s))
	}

	sb.WriteString(`
Your task:
Write a 2-5 sentence plain-English overview of what this file does as a whole — its purpose, the main abstractions it defines, and how it fits into the broader codebase. Then list 1-5 domain tags.

Respond with a single JSON object containing exactly these four keys:
- "symbol_name": the file name (just the base filename, e.g. "auth.go")
- "symbol_type": "file"
- "summary": a plain-text summary string (no markdown, no code fences)
- "domain_tags": a JSON array of short lowercase strings

Respond with valid JSON only. Do NOT wrap in markdown code fences. No extra text.`)
	return sb.String()
}

func buildTopicsPrompt(req *TopicGenerationRequest) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Repository: %s\n\n", req.RepoFullName))

	sb.WriteString("File summaries (one per file):\n")
	for path, summary := range req.FileSummaries {
		sb.WriteString(fmt.Sprintf("- %s: %s\n", path, summary))
	}

	if len(req.PreviousTopics) > 0 {
		sb.WriteString("\nTopics already covered in previous newsletters (DO NOT repeat these):\n")
		for _, t := range req.PreviousTopics {
			sb.WriteString(fmt.Sprintf("- %s\n", t))
		}
	}

	sb.WriteString(`
Your task:
1. Read all the file summaries above and identify the key features, modules, and subsystems this codebase implements. Understand how the files connect and relate to each other.
2. Based on that holistic understanding, propose 3-5 newsletter topics that would genuinely interest engineers working on this repository. Focus on areas of meaningful complexity, interesting design decisions, or important workflows — not trivial details.
3. Each topic must be meaningfully different from any previously covered topic listed above.

Respond with a JSON array where each element has:
- "title": a short, engaging question (10 words or fewer) phrased like "How does X work?", "What happens when Y does Z?", "Why does X use Y?", or "How is Z structured?" — use the actual names of systems or components from the codebase
- "description": 2-3 sentences describing the feature, module, or workflow and why it is interesting
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
Your audience: engineers working on this codebase who may not be familiar with this specific subsystem. Assume they are smart and have built things intentionally — do not lecture or suggest that existing choices are wrong outside the dedicated improvements section.

Structure your response as follows:
- <h2>: the topic title exactly as given
- Opening <p>: what this part of the codebase does and why it exists
- Body (2-4 <p> tags): how it currently works — key data flows, implementation details, design patterns, and how it connects to the rest of the system. Name actual files, functions, and patterns from the summaries above. Use <code> for inline identifiers and <pre><code> for short illustrative snippets.
- <h3>Suggested Improvements</h3>: open with a short <p> disclaimer — "These suggestions are AI-generated and may not make sense for your specific codebase — treat them as conversation starters, not prescriptions." — then 2-3 <p> tags each describing one concrete, actionable improvement.

Focus entirely on explaining what IS. Do not moralize, editorialize, or suggest improvements anywhere except the Suggested Improvements section.
Target length: approximately 1000 words.
Return valid HTML only. No markdown, no code fences, no surrounding html/body tags.`)
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

// stripMarkdownFences removes leading/trailing markdown code fences (e.g. ```json ... ```)
// that models sometimes add despite being instructed not to.
func stripMarkdownFences(s string) string {
	s = strings.TrimSpace(s)
	// Remove opening fence like ```json or ```
	if strings.HasPrefix(s, "```") {
		newline := strings.Index(s, "\n")
		if newline >= 0 {
			s = s[newline+1:]
		}
	}
	// Remove closing fence
	if strings.HasSuffix(s, "```") {
		s = s[:strings.LastIndex(s, "```")]
	}
	return strings.TrimSpace(s)
}

func parseSummarizeResult(raw string) *SummarizeResult {
	cleaned := stripMarkdownFences(raw)
	var parsed summarizeJSON
	if err := json.Unmarshal([]byte(cleaned), &parsed); err != nil {
		// Fallback: treat the entire raw text as the summary.
		return &SummarizeResult{Summary: strings.TrimSpace(raw)}
	}
	// Model signalled that this chunk has no meaningful logic (e.g. imports only).
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
