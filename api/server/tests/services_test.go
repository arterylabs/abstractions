package handlers_test

import (
	"abstractions/server/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServicesMiddleware(t *testing.T) {
	// Setup Echo without transaction middleware to verify InitializeServices
	// fails with a clear middleware error before reaching the handler.
	e := echo.New()
	e.Use(middlewares.InitializeServices())

	testHandler := func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "services_refactored"})
	}

	e.GET("/test", testHandler)

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.JSONEq(t, `{"message":"database transaction not found in context"}`, rec.Body.String())
}
