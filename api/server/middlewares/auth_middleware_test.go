package middlewares_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"abstractions/db/models"
	"abstractions/server/middlewares"
	"abstractions/utils/context_utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func errMessage(rec *httptest.ResponseRecorder) string {
	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		return ""
	}
	if v, ok := body["error"].(string); ok {
		return v
	}
	return ""
}

func TestRequireVerifiedAccount_UnverifiedForbidden(t *testing.T) {
	handler := middlewares.RequireVerifiedAccount(func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
	})

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// Unverified user in context (VerifiedAt nil).
	c = context_utils.SetUserInContext(c, models.User{
		ID:    uuid.New(),
		Email: "unverified@example.com",
	})

	err := handler(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Equal(t, "account not verified", errMessage(rec))
}

func TestRequireVerifiedAccount_VerifiedPasses(t *testing.T) {
	now := time.Now()
	handler := middlewares.RequireVerifiedAccount(func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
	})

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c = context_utils.SetUserInContext(c, models.User{
		ID:         uuid.New(),
		Email:      "verified@example.com",
		VerifiedAt: &now,
	})

	err := handler(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRequireVerifiedAccount_MissingUserForbidden(t *testing.T) {
	handler := middlewares.RequireVerifiedAccount(func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
	})

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Equal(t, "account not verified", errMessage(rec))
}

func TestAdminAuthMiddleware_NonAdminPathSkipped(t *testing.T) {
	// Path is not /admin*, so middleware is a no-op even without a user.
	handler := middlewares.AdminAuthMiddleware(func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
	})

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/workspaces/123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/workspaces/:id")

	err := handler(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAdminAuthMiddleware_AdminRouteWithoutGodModeForbidden(t *testing.T) {
	handler := middlewares.AdminAuthMiddleware(func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
	})

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/admin/users")
	c = context_utils.SetUserInContext(c, models.User{
		ID:      uuid.New(),
		Email:   "member@example.com",
		GodMode: false,
	})

	err := handler(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Equal(t, "forbidden", errMessage(rec))
}

func TestAdminAuthMiddleware_GodModeAllowed(t *testing.T) {
	handler := middlewares.AdminAuthMiddleware(func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
	})

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/admin/users")
	c = context_utils.SetUserInContext(c, models.User{
		ID:      uuid.New(),
		Email:   "god@example.com",
		GodMode: true,
	})

	err := handler(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAdminAuthMiddleware_AdminRouteMissingUserForbidden(t *testing.T) {
	handler := middlewares.AdminAuthMiddleware(func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
	})

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/admin/users")

	err := handler(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Equal(t, "forbidden", errMessage(rec))
}
