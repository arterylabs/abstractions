package tests

import (
	"abstractions/db"
	"abstractions/db/models"
	"abstractions/server/handlers"
	"abstractions/server/middlewares"
	"abstractions/tests/testutils"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WebhookHandlerTestSuite struct {
	testutils.TestSuite
	seeder   *testutils.Seeder
	enqueuer *mockTaskEnqueuer
}

func TestWebhookHandler(t *testing.T) {
	suite.Run(t, new(WebhookHandlerTestSuite))
}

func (s *WebhookHandlerTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.enqueuer = &mockTaskEnqueuer{}
}

func (s *WebhookHandlerTestSuite) buildWebhookEcho() *echo.Echo {
	e := echo.New()
	e.Use(db.WrapInTransaction(s.DB))
	e.Use(middlewares.InitializeServices())
	handler := handlers.NewWebhookHandler(s.Server, s.enqueuer)
	e.POST("/webhooks/github", handler.ReceiveGitHubWebhook)
	return e
}

func wbSignBody(body []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func wbPushBody(after string, commits []map[string]interface{}) []byte {
	b, _ := json.Marshal(map[string]interface{}{"after": after, "commits": commits})
	return b
}

func wbWithWebhook(webhookID int64, secret string) func(*models.Repository) {
	return func(r *models.Repository) {
		r.WebhookID = &webhookID
		r.WebhookSecret = &secret
	}
}

func wbRequest(hookID int64, event string, body []byte, secret string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/webhooks/github", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Hook-ID", fmt.Sprintf("%d", hookID))
	req.Header.Set("X-GitHub-Event", event)
	req.Header.Set("X-Hub-Signature-256", wbSignBody(body, secret))
	return req
}

// TestHMACVerification_ValidSignature: correctly signed request returns 204.
func (s *WebhookHandlerTestSuite) TestHMACVerification_ValidSignature() {
	ws := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	const secret = "test-secret-abc"
	const hookID = int64(111111)
	s.seeder.CreateRepository(ws.ID, user.ID, wbWithWebhook(hookID, secret))

	body := wbPushBody("abc123", nil)
	rec := httptest.NewRecorder()
	s.buildWebhookEcho().ServeHTTP(rec, wbRequest(hookID, "push", body, secret))
	assert.Equal(s.T(), http.StatusOK, rec.Code)
}

// TestHMACVerification_InvalidSignature: wrong secret returns 401.
func (s *WebhookHandlerTestSuite) TestHMACVerification_InvalidSignature() {
	ws := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	const secret = "correct-secret"
	const hookID = int64(222222)
	s.seeder.CreateRepository(ws.ID, user.ID, wbWithWebhook(hookID, secret))

	body := wbPushBody("abc123", nil)
	req := httptest.NewRequest(http.MethodPost, "/webhooks/github", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Hook-ID", fmt.Sprintf("%d", hookID))
	req.Header.Set("X-GitHub-Event", "push")
	req.Header.Set("X-Hub-Signature-256", wbSignBody(body, "wrong-secret"))

	rec := httptest.NewRecorder()
	s.buildWebhookEcho().ServeHTTP(rec, req)
	assert.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}

// TestHMACVerification_TamperedBody: body tampered after signing returns 401.
func (s *WebhookHandlerTestSuite) TestHMACVerification_TamperedBody() {
	ws := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	const secret = "another-secret"
	const hookID = int64(333333)
	s.seeder.CreateRepository(ws.ID, user.ID, wbWithWebhook(hookID, secret))

	originalBody := wbPushBody("abc123", nil)
	sig := wbSignBody(originalBody, secret)

	tamperedBody := wbPushBody("HACKED", nil)
	req := httptest.NewRequest(http.MethodPost, "/webhooks/github", bytes.NewReader(tamperedBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Hook-ID", fmt.Sprintf("%d", hookID))
	req.Header.Set("X-GitHub-Event", "push")
	req.Header.Set("X-Hub-Signature-256", sig)

	rec := httptest.NewRecorder()
	s.buildWebhookEcho().ServeHTTP(rec, req)
	assert.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}

// TestPushEvent_DeduplicatesFiles: a file touched in two commits produces one UpdateFileJob.
func (s *WebhookHandlerTestSuite) TestPushEvent_DeduplicatesFiles() {
	ws := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	const secret = "ded-secret"
	const hookID = int64(444444)
	s.seeder.CreateRepository(ws.ID, user.ID, wbWithWebhook(hookID, secret))

	commits := []map[string]interface{}{
		{"added": []string{}, "modified": []string{"main.go"}, "removed": []string{}},
		{"added": []string{}, "modified": []string{"main.go", "util.go"}, "removed": []string{}},
	}
	body := wbPushBody("sha1", commits)
	s.enqueuer.enqueuedTasks = nil

	rec := httptest.NewRecorder()
	s.buildWebhookEcho().ServeHTTP(rec, wbRequest(hookID, "push", body, secret))

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	assert.Len(s.T(), s.enqueuer.enqueuedTasks, 2)
	for _, t := range s.enqueuer.enqueuedTasks {
		assert.Equal(s.T(), handlers.UpdateFileJobID, t.Type())
	}
}

// TestPushEvent_IgnorePatternsFiltered: node_modules and go.sum are not enqueued.
func (s *WebhookHandlerTestSuite) TestPushEvent_IgnorePatternsFiltered() {
	ws := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	const secret = "ign-secret"
	const hookID = int64(555555)
	s.seeder.CreateRepository(ws.ID, user.ID, wbWithWebhook(hookID, secret))

	commits := []map[string]interface{}{
		{
			"added":    []string{"node_modules/react/index.js", "go.sum"},
			"modified": []string{"main.go"},
			"removed":  []string{},
		},
	}
	body := wbPushBody("sha2", commits)
	s.enqueuer.enqueuedTasks = nil

	rec := httptest.NewRecorder()
	s.buildWebhookEcho().ServeHTTP(rec, wbRequest(hookID, "push", body, secret))

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	assert.Len(s.T(), s.enqueuer.enqueuedTasks, 1)
	var payload handlers.UpdateFilePayload
	_ = json.Unmarshal(s.enqueuer.enqueuedTasks[0].Payload(), &payload)
	assert.Equal(s.T(), "main.go", payload.FilePath)
}

// TestPushEvent_RemovedFilesEnqueueDeleteJobs: removed files produce DeleteFileJob.
func (s *WebhookHandlerTestSuite) TestPushEvent_RemovedFilesEnqueueDeleteJobs() {
	ws := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	const secret = "del-secret"
	const hookID = int64(666666)
	s.seeder.CreateRepository(ws.ID, user.ID, wbWithWebhook(hookID, secret))

	commits := []map[string]interface{}{
		{"added": []string{}, "modified": []string{}, "removed": []string{"old.go"}},
	}
	body := wbPushBody("sha3", commits)
	s.enqueuer.enqueuedTasks = nil

	rec := httptest.NewRecorder()
	s.buildWebhookEcho().ServeHTTP(rec, wbRequest(hookID, "push", body, secret))

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	assert.Len(s.T(), s.enqueuer.enqueuedTasks, 1)
	assert.Equal(s.T(), handlers.DeleteFileJobID, s.enqueuer.enqueuedTasks[0].Type())
}

// TestUnknownWebhookID_Returns200: unknown hook IDs return 200 to prevent GitHub retries.
func (s *WebhookHandlerTestSuite) TestUnknownWebhookID_Returns200() {
	body := wbPushBody("sha4", nil)
	req := httptest.NewRequest(http.MethodPost, "/webhooks/github", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Hook-ID", "99999999")
	req.Header.Set("X-GitHub-Event", "push")
	req.Header.Set("X-Hub-Signature-256", wbSignBody(body, "any"))

	rec := httptest.NewRecorder()
	s.buildWebhookEcho().ServeHTTP(rec, req)
	assert.Equal(s.T(), http.StatusOK, rec.Code)
}
