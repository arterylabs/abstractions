package ai_client

import (
	"context"
	"encoding/json"
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
