package ai_client

import (
	cache_client "abstractions/clients/cache"
	"context"
	"fmt"
	"log"
	"time"
)

// completionRequest is the provider-agnostic input to a single text generation call.
// Both anthropicProvider and openAIProvider accept this type in their transport methods.
// Set nocache to true for non-deterministic operations (e.g. newsletter generation)
// that should never be served from the prompt cache.
type completionRequest struct {
	model        string
	systemPrompt string
	userPrompt   string
	maxTokens    int
	nocache      bool
}

// compositeProvider implements AIProvider by routing each operation to the correct
// provider based on the configured model's Provider field. Text generation uses either
// anthropicProvider or openAIProvider depending on which provider the model belongs to.
// Embedding generation uses openAIProvider.
type compositeProvider struct {
	anthropic       *anthropicProvider
	openai          *openAIProvider
	indexingModel   string
	newsletterModel string
	embeddingsModel string
	cache           cache_client.Client
	cacheTTL        time.Duration
}

func (c *compositeProvider) SummarizeCode(ctx context.Context, filePath, language, chunkContent string, startLine, endLine int, chunkType string) (*SummarizeResult, error) {
	budget := contentBudgetChars(c.indexingModel, maxTokensSummary, 50)
	normalizedContent := normalizeCodeContent(chunkContent)
	if len(normalizedContent) > budget {
		log.Printf("ai: SummarizeCode: truncating chunk content from %d to %d chars for model %s", len(normalizedContent), budget, c.indexingModel)
		normalizedContent = truncateContent(normalizedContent, budget)
	}
	rawResponse, err := c.complete(ctx, completionRequest{
		model:        c.indexingModel,
		systemPrompt: summarizeCodeSystemPrompt,
		userPrompt:   buildSummarizeUserPrompt(filePath, language, normalizedContent, startLine, endLine, chunkType),
		maxTokens:    maxTokensSummary,
	})
	if err != nil {
		return nil, fmt.Errorf("ai: SummarizeCode: %w", err)
	}
	return parseSummarizeResult(rawResponse), nil
}

func (c *compositeProvider) SummarizeFile(ctx context.Context, filePath, language string, chunkSummaries []string) (*SummarizeResult, error) {
	budget := contentBudgetChars(c.indexingModel, maxTokensSummary, 50)
	totalChars := 0
	for _, s := range chunkSummaries {
		totalChars += len(s)
	}
	if totalChars > budget {
		log.Printf("ai: SummarizeFile: truncating %d summary chars to %d for model %s", totalChars, budget, c.indexingModel)
		chunkSummaries = truncateSummarySlice(chunkSummaries, budget)
	}
	rawResponse, err := c.complete(ctx, completionRequest{
		model:        c.indexingModel,
		systemPrompt: summarizeFileSystemPrompt,
		userPrompt:   buildSummarizeFileUserPrompt(filePath, language, chunkSummaries),
		maxTokens:    maxTokensSummary,
	})
	if err != nil {
		return nil, fmt.Errorf("ai: SummarizeFile: %w", err)
	}
	return parseSummarizeResult(rawResponse), nil
}

func (c *compositeProvider) GenerateTopics(ctx context.Context, req *TopicGenerationRequest) ([]Topic, error) {
	budget := contentBudgetChars(c.newsletterModel, maxTokensTopics, 200)
	totalChars := 0
	for _, summary := range req.FileSummaries {
		totalChars += len(summary)
	}
	if totalChars > budget {
		log.Printf("ai: GenerateTopics: truncating %d summary chars to %d for model %s", totalChars, budget, c.newsletterModel)
		req = &TopicGenerationRequest{
			RepoFullName:   req.RepoFullName,
			FileSummaries:  truncateSummaryMap(req.FileSummaries, budget),
			PreviousTopics: req.PreviousTopics,
		}
	}
	rawResponse, err := c.complete(ctx, completionRequest{
		model:      c.newsletterModel,
		userPrompt: buildTopicsPrompt(req),
		maxTokens:  maxTokensTopics,
		nocache:    true, // newsletter generation should not be served from cache
	})
	if err != nil {
		return nil, fmt.Errorf("ai: GenerateTopics: %w", err)
	}
	return parseTopics(rawResponse), nil
}

func (c *compositeProvider) GenerateNewsletterEntry(ctx context.Context, req *NewsletterEntryRequest) (*NewsletterEntry, error) {
	budget := contentBudgetChars(c.newsletterModel, maxTokensNewsletter, 250)
	totalChars := 0
	for _, summary := range req.FileSummaries {
		totalChars += len(summary)
	}
	if totalChars > budget {
		log.Printf("ai: GenerateNewsletterEntry: truncating %d summary chars to %d for model %s", totalChars, budget, c.newsletterModel)
		req = &NewsletterEntryRequest{
			RepoFullName:  req.RepoFullName,
			Topic:         req.Topic,
			FileSummaries: truncateSummaryMap(req.FileSummaries, budget),
		}
	}
	rawResponse, err := c.complete(ctx, completionRequest{
		model:      c.newsletterModel,
		userPrompt: buildNewsletterEntryPrompt(req),
		maxTokens:  maxTokensNewsletter,
		nocache:    true, // newsletter generation should not be served from cache
	})
	if err != nil {
		return nil, fmt.Errorf("ai: GenerateNewsletterEntry: %w", err)
	}
	return &NewsletterEntry{Title: req.Topic.Title, Content: rawResponse}, nil
}

func (c *compositeProvider) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	results, err := c.GenerateEmbeddings(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	if len(results) == 0 || results[0] == nil {
		return nil, fmt.Errorf("ai: embedding: no data in response")
	}
	return results[0], nil
}

func (c *compositeProvider) GenerateEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	if c.openai == nil {
		return nil, fmt.Errorf("ai: embeddings: no OpenAI provider configured")
	}
	if len(texts) == 0 {
		return nil, nil
	}
	return c.openai.embed(ctx, c.embeddingsModel, texts)
}

// complete routes a text completion request to the appropriate provider based on the
// model's Provider field. This allows mixing Anthropic and OpenAI models across jobs —
// e.g. GPT-4.1 Mini for indexing and Claude Sonnet for newsletter generation.
func (c *compositeProvider) complete(ctx context.Context, req completionRequest) (string, error) {
	var cacheKey string
	if c.cache != nil && !req.nocache {
		cacheKey = promptCacheKey(req)
		cachedResponse, found, err := c.cache.Get(ctx, cacheKey)
		if err != nil {
			log.Printf("ai: prompt cache read failed for key %s: %v", cacheKey, err)
		} else if found {
			return cachedResponse, nil
		}
	}

	modelInfo, ok := ModelByID(req.model)
	if !ok {
		return "", fmt.Errorf("ai: unknown model %q — add it to the models list or choose a registered model", req.model)
	}

	var output string
	var err error
	switch modelInfo.Provider {
	case ProviderAnthropic:
		if c.anthropic == nil {
			return "", fmt.Errorf("ai: model %q requires an Anthropic API key", req.model)
		}
		output, err = c.anthropic.complete(ctx, req)
	case ProviderOpenAI:
		if c.openai == nil {
			return "", fmt.Errorf("ai: model %q requires an OpenAI API key", req.model)
		}
		output, err = c.openai.completeText(ctx, req)
	default:
		return "", fmt.Errorf("ai: no provider registered for model %q", req.model)
	}

	if err != nil {
		return "", err
	}

	if cacheKey != "" {
		if err := c.cache.Set(ctx, cacheKey, output, c.cacheTTL); err != nil {
			log.Printf("ai: prompt cache write failed for key %s: %v", cacheKey, err)
		}
	}

	return output, nil
}
