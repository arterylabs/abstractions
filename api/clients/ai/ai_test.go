package ai_client

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type inMemoryCache struct {
	mu    sync.Mutex
	store map[string]string
}

func newInMemoryCache() *inMemoryCache {
	return &inMemoryCache{store: map[string]string{}}
}

func (c *inMemoryCache) Get(ctx context.Context, key string) (string, bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.store[key]
	return value, ok, nil
}

func (c *inMemoryCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = value
	return nil
}

func (c *inMemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
	return nil
}

// mockServer returns an httptest.Server that responds with the given body and status code.
func mockServer(statusCode int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(body))
	}))
}

// anthropicSuccessBody builds a minimal Anthropic Messages API response body.
func anthropicSuccessBody(text string) string {
	resp := anthropicAPIResponse{
		Content: []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		}{{Type: "text", Text: text}},
	}
	b, _ := json.Marshal(resp)
	return string(b)
}

// testTextProviders builds the full set of text transports (one per registry provider) pointed at
// the given HTTP client, so any model in the registry routes correctly during tests.
func testTextProviders(httpClient *http.Client) map[Provider]textProvider {
	return map[Provider]textProvider{
		ProviderAnthropic: newAnthropicProvider("test-key", httpClient),
		ProviderOpenAI:    newOpenAIProvider("test-openai-key", httpClient),
		ProviderGrok:      newOpenAICompatibleProvider(grokAPIBase, "test-grok-key", httpClient),
	}
}

// compositeWithTestServer creates a compositeProvider that points all HTTP requests at the given test server.
// Uses Anthropic models by default so text generation routes to anthropicProvider.
func compositeWithTestServer(server *httptest.Server) *compositeProvider {
	transport := &rewriteTransport{target: server.URL, inner: server.Client().Transport}
	httpClient := &http.Client{Transport: transport}
	return &compositeProvider{
		textProviders:   testTextProviders(httpClient),
		embeddings:      newOpenAIProvider("test-openai-key", httpClient),
		indexingModel:   defaultIndexingModel,
		newsletterModel: defaultNewsletterModel,
		embeddingsModel: defaultEmbeddingsModel,
	}
}

// rewriteTransport rewrites all requests to point at the target server URL (ignoring the original host).
type rewriteTransport struct {
	target string
	inner  http.RoundTripper
}

func (t *rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.URL.Scheme = "http"
	req.URL.Host = strings.TrimPrefix(t.target, "http://")
	req.URL.Path = "/"
	return t.inner.RoundTrip(req)
}

// --- Anthropic routing tests ---

func TestAnthropicProvider_SummarizeCode_ParsesJSON(t *testing.T) {
	responseText := `{"summary":"Handles user authentication via JWT.","domain_tags":["auth","jwt"]}`
	server := mockServer(http.StatusOK, anthropicSuccessBody(responseText))
	defer server.Close()

	provider := compositeWithTestServer(server)
	result, err := provider.SummarizeCode(context.Background(), "api/auth.go", "Go", "func ValidateToken(...) {}", 1, 1, "function")
	require.NoError(t, err)
	assert.Equal(t, "Handles user authentication via JWT.", result.Summary)
	assert.Equal(t, []string{"auth", "jwt"}, result.DomainTags)
}

func TestAnthropicProvider_SummarizeCode_FallsBackOnBadJSON(t *testing.T) {
	server := mockServer(http.StatusOK, anthropicSuccessBody("This code validates tokens."))
	defer server.Close()

	provider := compositeWithTestServer(server)
	result, err := provider.SummarizeCode(context.Background(), "api/auth.go", "Go", "func ValidateToken(...) {}", 1, 1, "function")
	require.NoError(t, err)
	assert.Contains(t, result.Summary, "This code validates tokens.")
}

func TestAnthropicProvider_SummarizeCode_ReturnsErrorOnAPIError(t *testing.T) {
	errBody := `{"error":{"type":"authentication_error","message":"Invalid API key"}}`
	server := mockServer(http.StatusOK, errBody)
	defer server.Close()

	provider := compositeWithTestServer(server)
	_, err := provider.SummarizeCode(context.Background(), "api/auth.go", "Go", "func ValidateToken(...) {}", 1, 1, "function")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "authentication_error")
}

func TestAnthropicProvider_GenerateTopics_ParsesJSON(t *testing.T) {
	topicsText := `[{"title":"JWT Auth","description":"New token validation.","related_files":["api/auth.go"]}]`
	server := mockServer(http.StatusOK, anthropicSuccessBody(topicsText))
	defer server.Close()

	provider := compositeWithTestServer(server)
	req := &TopicGenerationRequest{
		RepoFullName:  "octocat/myrepo",
		FileSummaries: map[string]string{"api/auth.go": "Handles JWT validation."},
	}
	topics, err := provider.GenerateTopics(context.Background(), req)
	require.NoError(t, err)
	require.Len(t, topics, 1)
	assert.Equal(t, "JWT Auth", topics[0].Title)
	assert.Equal(t, []string{"api/auth.go"}, topics[0].RelatedFiles)
}

func TestAnthropicProvider_GenerateNewsletterEntry_ReturnsHTMLContent(t *testing.T) {
	htmlContent := "<h2>Auth Refactor</h2><p>We improved JWT handling.</p>"
	server := mockServer(http.StatusOK, anthropicSuccessBody(htmlContent))
	defer server.Close()

	provider := compositeWithTestServer(server)
	req := &NewsletterEntryRequest{
		RepoFullName:  "octocat/myrepo",
		Topic:         Topic{Title: "Auth Refactor", Description: "JWT improvements"},
		FileSummaries: map[string]string{"api/auth.go": "Validates tokens"},
	}
	entry, err := provider.GenerateNewsletterEntry(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, "Auth Refactor", entry.Title)
	assert.Equal(t, htmlContent, entry.Content)
}

// --- OpenAI routing tests ---

func openAIChatSuccessBody(text string) string {
	resp := openAIChatResponse{
		Choices: []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		}{{Message: struct {
			Content string `json:"content"`
		}{Content: text}, FinishReason: "stop"}},
	}
	b, _ := json.Marshal(resp)
	return string(b)
}

func TestOpenAIProvider_SummarizeCode_ParsesJSON(t *testing.T) {
	responseText := `{"summary":"Handles user authentication via JWT.","domain_tags":["auth","jwt"]}`
	server := mockServer(http.StatusOK, openAIChatSuccessBody(responseText))
	defer server.Close()

	transport := &rewriteTransport{target: server.URL, inner: server.Client().Transport}
	httpClient := &http.Client{Transport: transport}
	provider := &compositeProvider{
		textProviders:   testTextProviders(httpClient),
		embeddings:      newOpenAIProvider("test-openai-key", httpClient),
		indexingModel:   "gpt-4.1-mini",
		newsletterModel: "gpt-4.1",
		embeddingsModel: defaultEmbeddingsModel,
	}

	result, err := provider.SummarizeCode(context.Background(), "api/auth.go", "Go", "func ValidateToken(...) {}", 1, 1, "function")
	require.NoError(t, err)
	assert.Equal(t, "Handles user authentication via JWT.", result.Summary)
	assert.Equal(t, []string{"auth", "jwt"}, result.DomainTags)
}

func TestOpenAIProvider_GenerateEmbedding_ReturnsCorrectDimension(t *testing.T) {
	server := mockServer(http.StatusOK, `{"data":[{"index":0,"embedding":[0.1,0.2,0.3]}]}`)
	defer server.Close()

	provider := compositeWithTestServer(server)
	vec, err := provider.GenerateEmbedding(context.Background(), "some text")
	require.NoError(t, err)
	assert.Len(t, vec, 3)
}

func TestComposite_ReturnsErrorForUnknownModel(t *testing.T) {
	provider := &compositeProvider{
		textProviders: map[Provider]textProvider{
			ProviderAnthropic: newAnthropicProvider("key", nil),
			ProviderOpenAI:    newOpenAIProvider("key", nil),
		},
		indexingModel:   "nonexistent-model-xyz",
		newsletterModel: defaultNewsletterModel,
		embeddingsModel: defaultEmbeddingsModel,
	}
	_, err := provider.SummarizeCode(context.Background(), "file.go", "Go", "func foo() {}", 1, 1, "function")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown model")
}

func TestComposite_ReturnsErrorWhenAnthropicKeyMissing(t *testing.T) {
	provider := &compositeProvider{
		textProviders: map[Provider]textProvider{
			ProviderAnthropic: newAnthropicProvider("", nil),
		},
		indexingModel:   defaultIndexingModel,
		newsletterModel: defaultNewsletterModel,
		embeddingsModel: defaultEmbeddingsModel,
	}
	// Should fail at the HTTP level, not at routing — empty key is still routed to anthropic.
	// We just verify it doesn't panic and errors cleanly.
	_, err := provider.SummarizeCode(context.Background(), "file.go", "Go", "func foo() {}", 1, 1, "function")
	require.Error(t, err)
}

func TestComposite_ReturnsErrorWhenOpenAINotConfigured(t *testing.T) {
	provider := &compositeProvider{
		textProviders: map[Provider]textProvider{
			ProviderAnthropic: newAnthropicProvider("key", nil),
		},
		embeddings:      nil,
		indexingModel:   defaultIndexingModel,
		newsletterModel: defaultNewsletterModel,
		embeddingsModel: defaultEmbeddingsModel,
	}
	_, err := provider.GenerateEmbeddings(context.Background(), []string{"text"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no OpenAI provider configured")
}

// --- Utility tests ---

func TestBuildSummarizeUserPrompt_ContainsFilePathAndLanguage(t *testing.T) {
	prompt := buildSummarizeUserPrompt("cmd/main.go", "Go", "package main", 1, 1, "file")
	assert.Contains(t, prompt, "cmd/main.go")
	assert.Contains(t, prompt, "Go")
}

func TestBuildSummarizeUserPrompt_ContainsChunkContent(t *testing.T) {
	prompt := buildSummarizeUserPrompt("file.go", "Go", "func hello() {}", 3, 5, "function")
	assert.Contains(t, prompt, "func hello() {}")
}

func TestNormalizeCodeContent_StripsCommonIndent(t *testing.T) {
	code := "    func foo() {\n        return 1\n    }"
	result := normalizeCodeContent(code)
	assert.Equal(t, "func foo() {\n    return 1\n}", result)
}

func TestNormalizeCodeContent_PreservesRelativeIndent(t *testing.T) {
	code := "    def foo():\n        x = 1\n        return x"
	result := normalizeCodeContent(code)
	assert.Equal(t, "def foo():\n    x = 1\n    return x", result)
}

func TestNormalizeCodeContent_CollapsesExtraBlankLines(t *testing.T) {
	code := "func a() {}\n\n\n\nfunc b() {}"
	result := normalizeCodeContent(code)
	assert.Equal(t, "func a() {}\n\n\nfunc b() {}", result)
}

func TestBuildTopicsPrompt_TruncatesLongSummaries(t *testing.T) {
	longSummary := strings.Repeat("x", maxTopicSummaryLength+100)
	req := &TopicGenerationRequest{
		RepoFullName:  "owner/repo",
		FileSummaries: map[string]string{"main.go": longSummary},
	}
	prompt := buildTopicsPrompt(req)
	assert.Contains(t, prompt, "...")
	assert.NotContains(t, prompt, longSummary)
}

func TestBuildTopicsPrompt_ContainsRepoAndSummaries(t *testing.T) {
	req := &TopicGenerationRequest{
		RepoFullName:  "owner/repo",
		FileSummaries: map[string]string{"main.go": "entry point"},
	}
	prompt := buildTopicsPrompt(req)
	assert.Contains(t, prompt, "owner/repo")
	assert.Contains(t, prompt, "main.go")
	assert.Contains(t, prompt, "entry point")
}

func TestFactory_ReturnsCompositeProvider(t *testing.T) {
	provider := NewAIProvider()
	assert.NotNil(t, provider)
	_, ok := provider.(*compositeProvider)
	assert.True(t, ok, "NewAIProvider should return *compositeProvider")
}

// --- Provider registry / Grok tests ---

func TestProviderRegistry_IncludesGrok(t *testing.T) {
	descriptor, ok := ProviderByID(ProviderGrok)
	require.True(t, ok, "Grok must be registered in Providers")
	assert.Equal(t, ProviderGrok, descriptor.ID)
	assert.NotEmpty(t, descriptor.DisplayName)
	assert.NotEmpty(t, descriptor.ConfigKey)
	assert.False(t, descriptor.SupportsEmbeddings, "Grok has no embeddings API")
	require.NotNil(t, descriptor.NewTextProvider)

	// Grok's transport must be an OpenAI-compatible provider pointed at the xAI base URL.
	built, ok := descriptor.NewTextProvider("grok-key", nil).(*openAIProvider)
	require.True(t, ok, "Grok text provider should be an *openAIProvider")
	assert.Equal(t, grokAPIBase, built.baseURL)
}

// TestGrokModel_RoutesThroughOpenAICompatibleTransport proves that a Grok model's text generation is
// dispatched to the OpenAI-compatible Chat Completions endpoint (not the Anthropic transport).
func TestGrokModel_RoutesThroughOpenAICompatibleTransport(t *testing.T) {
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(openAIChatSuccessBody(`{"summary":"grok summary","domain_tags":["grok"]}`)))
	}))
	defer server.Close()

	// Point the Grok transport's baseURL directly at the test server so the request path is real.
	provider := &compositeProvider{
		textProviders: map[Provider]textProvider{
			ProviderGrok: newOpenAICompatibleProvider(server.URL, "grok-key", server.Client()),
		},
		indexingModel:   "grok-3",
		newsletterModel: "grok-3",
		embeddingsModel: defaultEmbeddingsModel,
	}

	result, err := provider.SummarizeCode(context.Background(), "api/auth.go", "Go", "func Foo() {}", 1, 1, "function")
	require.NoError(t, err)
	assert.Equal(t, "grok summary", result.Summary)
	assert.Equal(t, "/chat/completions", gotPath, "Grok should route to the OpenAI-compatible chat completions endpoint")
}

// --- contentBudgetChars tests ---

func TestContentBudgetChars_UnknownModel_ReturnsMaxInt(t *testing.T) {
	result := contentBudgetChars("not-a-real-model", 500, 50)
	assert.Equal(t, math.MaxInt, result)
}

func TestContentBudgetChars_KnownModel_ReturnsCorrectBudget(t *testing.T) {
	// claude-haiku-4-5-20251001 has ContextWindow = 200_000
	// availableTokens = 200000 - 800 - 50 = 199150; chars = 199150 * 4
	result := contentBudgetChars("claude-haiku-4-5-20251001", 800, 50)
	assert.Equal(t, (200_000-800-50)*4, result)
}

func TestContentBudgetChars_AvailableTokensNegative_ReturnsZero(t *testing.T) {
	// text-embedding-3-small has ContextWindow = 8_191
	// 8191 - 8191 - 100 = -100 → clamped to 0
	result := contentBudgetChars("text-embedding-3-small", 8_191, 100)
	assert.Equal(t, 0, result)
}

// --- truncateContent tests ---

func TestTruncateContent_ShortContent_Unchanged(t *testing.T) {
	result := truncateContent("hello", 100)
	assert.Equal(t, "hello", result)
}

func TestTruncateContent_ExactLength_Unchanged(t *testing.T) {
	result := truncateContent("hello", 5)
	assert.Equal(t, "hello", result)
}

func TestTruncateContent_LongContent_TruncatedWithSuffix(t *testing.T) {
	result := truncateContent("abcdefghij", 5)
	assert.Equal(t, "abcde\n[... truncated to fit model context window ...]", result)
}

func TestTruncateContent_ZeroMax_TruncatesImmediately(t *testing.T) {
	result := truncateContent("hello", 0)
	assert.Equal(t, "\n[... truncated to fit model context window ...]", result)
}

// --- truncateSummarySlice tests ---

func TestTruncateSummarySlice_EmptySlice_ReturnedUnchanged(t *testing.T) {
	result := truncateSummarySlice([]string{}, 1000)
	assert.Empty(t, result)
}

func TestTruncateSummarySlice_NormalCase_EachEntryTruncated(t *testing.T) {
	summaries := []string{"aaaaaa", "bbbbbb"}
	// budget 4 / 2 entries = 2 chars each
	result := truncateSummarySlice(summaries, 4)
	require.Len(t, result, 2)
	assert.Equal(t, "aa\n[... truncated to fit model context window ...]", result[0])
	assert.Equal(t, "bb\n[... truncated to fit model context window ...]", result[1])
}

func TestTruncateSummarySlice_ZeroBudget_ReturnsEmptySlice(t *testing.T) {
	result := truncateSummarySlice([]string{"hello", "world"}, 0)
	assert.Empty(t, result)
}

// --- truncateSummaryMap tests ---

func TestTruncateSummaryMap_EmptyMap_ReturnedUnchanged(t *testing.T) {
	result := truncateSummaryMap(map[string]string{}, 1000)
	assert.Empty(t, result)
}

func TestTruncateSummaryMap_NormalCase_ValuesAreTruncated(t *testing.T) {
	summaries := map[string]string{"file.go": "abcdefghij"}
	// budget 4 / 1 entry = 4 chars
	result := truncateSummaryMap(summaries, 4)
	assert.Equal(t, "abcd\n[... truncated to fit model context window ...]", result["file.go"])
}

func TestTruncateSummaryMap_ZeroBudget_ReturnsEmptyMap(t *testing.T) {
	result := truncateSummaryMap(map[string]string{"file.go": "some summary"}, 0)
	assert.Empty(t, result)
}

// --- stripMarkdownFences tests ---

func TestStripMarkdownFences_NoFences_Unchanged(t *testing.T) {
	result := stripMarkdownFences("plain text")
	assert.Equal(t, "plain text", result)
}

func TestStripMarkdownFences_LeadingFenceWithLanguage_Stripped(t *testing.T) {
	result := stripMarkdownFences("```json\n{\"key\":\"value\"}\n```")
	assert.Equal(t, `{"key":"value"}`, result)
}

func TestStripMarkdownFences_LeadingFenceNoLanguage_Stripped(t *testing.T) {
	result := stripMarkdownFences("```\nsome content\n```")
	assert.Equal(t, "some content", result)
}

func TestStripMarkdownFences_TrailingFenceOnly_Stripped(t *testing.T) {
	result := stripMarkdownFences("some content\n```")
	assert.Equal(t, "some content", result)
}

// --- parseSummarizeResult tests ---

func TestParseSummarizeResult_ValidJSON_ReturnsParsedResult(t *testing.T) {
	raw := `{"summary":"does auth","domain_tags":["auth"],"symbol_name":"ValidateToken","symbol_type":"function"}`
	result := parseSummarizeResult(raw)
	require.NotNil(t, result)
	assert.Equal(t, "does auth", result.Summary)
	assert.Equal(t, []string{"auth"}, result.DomainTags)
	assert.Equal(t, "ValidateToken", result.SymbolName)
	assert.Equal(t, "function", result.SymbolType)
}

func TestParseSummarizeResult_SkipTrue_ReturnsNil(t *testing.T) {
	result := parseSummarizeResult(`{"skip":true}`)
	assert.Nil(t, result)
}

func TestParseSummarizeResult_InvalidJSON_FallsBackToRawString(t *testing.T) {
	result := parseSummarizeResult("not json at all")
	require.NotNil(t, result)
	assert.Equal(t, "not json at all", result.Summary)
}

func TestParseSummarizeResult_JSONWrappedInFences_StripsAndParses(t *testing.T) {
	raw := "```json\n{\"summary\":\"handles auth\",\"domain_tags\":[\"auth\"]}\n```"
	result := parseSummarizeResult(raw)
	require.NotNil(t, result)
	assert.Equal(t, "handles auth", result.Summary)
	assert.Equal(t, []string{"auth"}, result.DomainTags)
}

// --- parseTopics tests ---

func TestParseTopics_ValidJSON_ReturnsParsedTopics(t *testing.T) {
	raw := `[{"title":"How auth works?","description":"JWT flow","related_files":["auth.go"]}]`
	topics := parseTopics(raw)
	require.Len(t, topics, 1)
	assert.Equal(t, "How auth works?", topics[0].Title)
	assert.Equal(t, "JWT flow", topics[0].Description)
	assert.Equal(t, []string{"auth.go"}, topics[0].RelatedFiles)
}

func TestParseTopics_InvalidJSON_ReturnsNil(t *testing.T) {
	topics := parseTopics("not json")
	assert.Nil(t, topics)
}

func TestParseTopics_EmptyArray_ReturnsEmptySlice(t *testing.T) {
	topics := parseTopics(`[]`)
	assert.Empty(t, topics)
}

// --- ModelByID tests ---

func TestModelByID_KnownModel_ReturnsModel(t *testing.T) {
	m, ok := ModelByID("claude-sonnet-4-6")
	assert.True(t, ok)
	assert.Equal(t, "claude-sonnet-4-6", m.ID)
	assert.Equal(t, ProviderAnthropic, m.Provider)
}

func TestModelByID_UnknownModel_ReturnsFalse(t *testing.T) {
	_, ok := ModelByID("totally-fake-model")
	assert.False(t, ok)
}

func TestPromptCacheKey_SameRequestSameKey(t *testing.T) {
	req := completionRequest{
		model:        "claude-sonnet-4-6",
		systemPrompt: "You are a coding assistant",
		userPrompt:   "Summarize this file",
		maxTokens:    800,
	}
	assert.Equal(t, promptCacheKey(req), promptCacheKey(req))
}

func TestPromptCacheKey_DifferentRequestDifferentKey(t *testing.T) {
	reqA := completionRequest{
		model:        "claude-sonnet-4-6",
		systemPrompt: "system",
		userPrompt:   "same prompt",
		maxTokens:    800,
	}
	reqB := completionRequest{
		model:        "gpt-4.1-mini",
		systemPrompt: "system",
		userPrompt:   "same prompt",
		maxTokens:    800,
	}
	assert.NotEqual(t, promptCacheKey(reqA), promptCacheKey(reqB))
}

func TestCompositeProvider_Complete_UsesPromptCache(t *testing.T) {
	var callCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(anthropicSuccessBody(`{"summary":"from provider","domain_tags":["cache"]}`)))
	}))
	defer server.Close()

	transport := &rewriteTransport{target: server.URL, inner: server.Client().Transport}
	httpClient := &http.Client{Transport: transport}
	provider := &compositeProvider{
		textProviders:   testTextProviders(httpClient),
		embeddings:      newOpenAIProvider("test-openai-key", httpClient),
		indexingModel:   defaultIndexingModel,
		newsletterModel: defaultNewsletterModel,
		embeddingsModel: defaultEmbeddingsModel,
		cache:           newInMemoryCache(),
		cacheTTL:        7 * 24 * time.Hour,
	}

	first, err := provider.SummarizeCode(context.Background(), "api/auth.go", "Go", "func ValidateToken() {}", 1, 1, "function")
	require.NoError(t, err)
	second, err := provider.SummarizeCode(context.Background(), "api/auth.go", "Go", "func ValidateToken() {}", 1, 1, "function")
	require.NoError(t, err)

	require.NotNil(t, first)
	require.NotNil(t, second)
	assert.Equal(t, first.Summary, second.Summary)
	assert.Equal(t, 1, callCount, "second identical request should be served from cache")
}

// --- TextGenerationModels / EmbeddingModels tests ---

func TestTextGenerationModels_ReturnsOnlyTextGenModels(t *testing.T) {
	models := TextGenerationModels()
	assert.NotEmpty(t, models)
	for _, m := range models {
		assert.True(t, m.SupportsTextGeneration, "expected SupportsTextGeneration for %s", m.ID)
	}
}

func TestEmbeddingModels_ReturnsOnlyEmbeddingModels(t *testing.T) {
	models := EmbeddingModels()
	assert.NotEmpty(t, models)
	for _, m := range models {
		assert.True(t, m.SupportsEmbeddings, "expected SupportsEmbeddings for %s", m.ID)
	}
}

// --- buildSummarizeFileUserPrompt tests ---

func TestBuildSummarizeFileUserPrompt_ContainsFileAndLanguage(t *testing.T) {
	prompt := buildSummarizeFileUserPrompt("src/auth.go", "Go", []string{"validates JWT tokens"})
	assert.Contains(t, prompt, "src/auth.go")
	assert.Contains(t, prompt, "Go")
	assert.Contains(t, prompt, "validates JWT tokens")
}

func TestBuildSummarizeFileUserPrompt_NumbersChunkSummaries(t *testing.T) {
	prompt := buildSummarizeFileUserPrompt("file.go", "Go", []string{"first chunk", "second chunk"})
	assert.Contains(t, prompt, "1. first chunk")
	assert.Contains(t, prompt, "2. second chunk")
}

// --- buildNewsletterEntryPrompt tests ---

func TestBuildNewsletterEntryPrompt_ContainsRepoTopicAndDescription(t *testing.T) {
	req := &NewsletterEntryRequest{
		RepoFullName:  "owner/repo",
		Topic:         Topic{Title: "Auth Flow", Description: "JWT details"},
		FileSummaries: map[string]string{"auth.go": "handles JWT"},
	}
	prompt := buildNewsletterEntryPrompt(req)
	assert.Contains(t, prompt, "owner/repo")
	assert.Contains(t, prompt, "Auth Flow")
	assert.Contains(t, prompt, "JWT details")
}

func TestBuildNewsletterEntryPrompt_ContainsFileSummaries(t *testing.T) {
	req := &NewsletterEntryRequest{
		RepoFullName:  "owner/repo",
		Topic:         Topic{Title: "Topic", Description: "Desc"},
		FileSummaries: map[string]string{"main.go": "entry point logic"},
	}
	prompt := buildNewsletterEntryPrompt(req)
	assert.Contains(t, prompt, "main.go")
	assert.Contains(t, prompt, "entry point logic")
}

// --- normalizeCodeContent additional edge cases ---

func TestNormalizeCodeContent_EmptyString_ReturnsEmpty(t *testing.T) {
	result := normalizeCodeContent("")
	assert.Equal(t, "", result)
}

func TestNormalizeCodeContent_CRLFLineEndings_Normalized(t *testing.T) {
	code := "func foo() {\r\n    return 1\r\n}"
	result := normalizeCodeContent(code)
	assert.NotContains(t, result, "\r")
	assert.Equal(t, "func foo() {\n    return 1\n}", result)
}

func TestNormalizeCodeContent_LeadingAndTrailingBlankLines_Trimmed(t *testing.T) {
	result := normalizeCodeContent("\n\nfunc foo() {}\n\n")
	assert.Equal(t, "func foo() {}", result)
}

func TestNormalizeCodeContent_TrailingWhitespace_Stripped(t *testing.T) {
	result := normalizeCodeContent("func foo() {}   \nfunc bar() {}   ")
	assert.Equal(t, "func foo() {}\nfunc bar() {}", result)
}

func TestNormalizeCodeContent_OnlyBlankLines_ReturnsEmpty(t *testing.T) {
	result := normalizeCodeContent("\n\n\n")
	assert.Equal(t, "", result)
}

// --- composeSystemPrompt / system prompt enhancements ---

func TestComposeSystemPrompt_EmptyUsesBaseOnly(t *testing.T) {
	assert.Equal(t, "default", composeSystemPrompt("default", ""))
	assert.Equal(t, "default", composeSystemPrompt("default", "   "))
	assert.Equal(t, "default", composeSystemPrompt("default", "\n\t"))
}

func TestComposeSystemPrompt_NonEmptyAppendsEnhancement(t *testing.T) {
	got := composeSystemPrompt("default", "custom voice")
	assert.True(t, strings.HasPrefix(got, "default"))
	assert.Contains(t, got, "custom voice")
	assert.Contains(t, got, "Additional workspace guidelines")
	// Enhancement is appended, not a replacement.
	assert.NotEqual(t, "custom voice", got)

	// Leading/trailing whitespace on enhancement is trimmed before append.
	got = composeSystemPrompt("default", "  kept padding  ")
	assert.True(t, strings.HasSuffix(got, "kept padding"))
	assert.NotContains(t, got, "  kept padding  ")
}

func TestDefaultSystemPrompts_MatchesBuiltins(t *testing.T) {
	defaults := DefaultSystemPrompts()
	assert.Equal(t, summarizeCodeSystemPrompt, defaults.SummarizeCode)
	assert.Equal(t, summarizeFileSystemPrompt, defaults.SummarizeFile)
	assert.Equal(t, topicGeneratorSystemPrompt, defaults.TopicGenerator)
	assert.Equal(t, newsletterEntrySystemPrompt, defaults.NewsletterEntry)
	// Every default must instruct the model about optional workspace guidelines.
	for _, p := range []string{
		defaults.SummarizeCode,
		defaults.SummarizeFile,
		defaults.TopicGenerator,
		defaults.NewsletterEntry,
	} {
		assert.Contains(t, p, "additional workspace guidelines")
		assert.Contains(t, p, "structure and format rules")
	}
}

func TestCompositeProvider_SystemPromptEnhancement_IsAppended(t *testing.T) {
	var capturedSystem string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body anthropicAPIRequest
		_ = json.NewDecoder(r.Body).Decode(&body)
		if len(body.System) > 0 {
			capturedSystem = body.System[0].Text
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(anthropicSuccessBody(`{"summary":"ok","domain_tags":[],"symbol_name":"x","symbol_type":"function"}`)))
	}))
	defer server.Close()

	provider := compositeWithTestServer(server)
	provider.systemPrompts = SystemPromptOverrides{
		SummarizeCode: "Prefer security-focused language.",
	}

	_, err := provider.SummarizeCode(context.Background(), "a.go", "Go", "func x() {}", 1, 1, "function")
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(capturedSystem, summarizeCodeSystemPrompt))
	assert.Contains(t, capturedSystem, "Prefer security-focused language.")
	assert.Contains(t, capturedSystem, "Additional workspace guidelines")
}

func TestCompositeProvider_EmptySystemPromptEnhancement_UsesDefault(t *testing.T) {
	var capturedSystem string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body anthropicAPIRequest
		_ = json.NewDecoder(r.Body).Decode(&body)
		if len(body.System) > 0 {
			capturedSystem = body.System[0].Text
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(anthropicSuccessBody(`{"summary":"ok","domain_tags":[],"symbol_name":"x","symbol_type":"function"}`)))
	}))
	defer server.Close()

	provider := compositeWithTestServer(server)
	provider.systemPrompts = SystemPromptOverrides{
		SummarizeCode: "   ",
	}

	_, err := provider.SummarizeCode(context.Background(), "a.go", "Go", "func x() {}", 1, 1, "function")
	require.NoError(t, err)
	assert.Equal(t, summarizeCodeSystemPrompt, capturedSystem)
}

func TestCompositeProvider_NewsletterEntryEnhancement_DoesNotAffectSummarize(t *testing.T) {
	var capturedSystem string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body anthropicAPIRequest
		_ = json.NewDecoder(r.Body).Decode(&body)
		if len(body.System) > 0 {
			capturedSystem = body.System[0].Text
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(anthropicSuccessBody(`{"summary":"ok","domain_tags":[],"symbol_name":"x","symbol_type":"function"}`)))
	}))
	defer server.Close()

	provider := compositeWithTestServer(server)
	provider.systemPrompts = SystemPromptOverrides{
		NewsletterEntry: "ONLY NEWSLETTER CUSTOM",
	}

	_, err := provider.SummarizeCode(context.Background(), "a.go", "Go", "func x() {}", 1, 1, "function")
	require.NoError(t, err)
	assert.Equal(t, summarizeCodeSystemPrompt, capturedSystem)
	assert.NotContains(t, capturedSystem, "ONLY NEWSLETTER CUSTOM")
}
