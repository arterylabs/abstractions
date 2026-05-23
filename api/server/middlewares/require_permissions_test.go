package middlewares_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"abstractions/db"
	"abstractions/db/models"
	"abstractions/server/middlewares"
	authorization_service "abstractions/services/authorization"
	user_service "abstractions/services/user"
	"abstractions/tests/testutils"
	"abstractions/utils/context_utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// RequirePermissionsTestSuite tests the RequirePermissions middleware end-to-end
// using a real test database and the full middleware stack.
type RequirePermissionsTestSuite struct {
	testutils.TestSuite
	seeder      *testutils.Seeder
	userService *user_service.Service
}

func TestRequirePermissions(t *testing.T) {
	suite.Run(t, new(RequirePermissionsTestSuite))
}

func (s *RequirePermissionsTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.userService = user_service.NewUserService(s.DB)
}

// setupEcho builds a test Echo instance with the production middleware stack
// (transaction wrapper → services initializer → auth) and registers a single
// route at the given method/path protected by the supplied permissions.
func (s *RequirePermissionsTestSuite) setupEcho(method, path string, permissions ...authorization_service.Permission) *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	// Replicate the production middleware order.
	e.Use(db.WrapInTransaction(s.DB))
	e.Use(middlewares.InitializeServices())
	e.Use(middlewares.AuthMiddleware)

	stub := func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	}
	e.Add(method, path, stub, middlewares.RequirePermissions(permissions...))
	return e
}

// do fires an HTTP request against the given Echo instance and returns the recorder.
func (s *RequirePermissionsTestSuite) do(e *echo.Echo, method, url, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, url, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

// errBody extracts the "error" field from a JSON error response body.
func errBody(rec *httptest.ResponseRecorder) string {
	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		return ""
	}
	if v, ok := body["error"].(string); ok {
		return v
	}
	return ""
}

// createUserWithToken creates a verified user + active session and returns
// the user model together with a ready-to-use Bearer token.
func (s *RequirePermissionsTestSuite) createUserWithToken() (*models.User, string) {
	user, err := s.userService.CreateVerifiedUser(user_service.CreateVerifiedUserParams{
		Email: "user-" + uuid.New().String() + "@test.com",
		Name:  "Test User",
	})
	s.Require().NoError(err)

	session, err := s.userService.CreateUserSession(user.ID, &user_service.CreateSessionParams{
		UserAgent: "test-agent",
		IP:        "127.0.0.1",
	})
	s.Require().NoError(err)
	return user, session.Token
}

// ── Middleware error path tests ──────────────────────────────────────────────

// TestNoAuthHeader verifies 401 when the Authorization header is absent.
func (s *RequirePermissionsTestSuite) TestNoAuthHeader() {
	workspace := s.seeder.CreateWorkspace()
	e := s.setupEcho(http.MethodPut, "/workspaces/:workspace_id", authorization_service.WorkspaceUpdate)

	rec := s.do(e, http.MethodPut, "/workspaces/"+workspace.ID.String(), "")
	assert.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}

// TestInvalidToken verifies 401 when a malformed / fake token is supplied.
func (s *RequirePermissionsTestSuite) TestInvalidToken() {
	workspace := s.seeder.CreateWorkspace()
	e := s.setupEcho(http.MethodPut, "/workspaces/:workspace_id", authorization_service.WorkspaceUpdate)

	rec := s.do(e, http.MethodPut, "/workspaces/"+workspace.ID.String(), "not-a-valid-jwt")
	assert.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}

// TestNoWorkspaceIDParam verifies that the middleware is a no-op when the route
// does not have a :workspace_id parameter, so the stub handler returns 200.
func (s *RequirePermissionsTestSuite) TestNoWorkspaceIDParam() {
	_, token := s.createUserWithToken()
	e := s.setupEcho(http.MethodGet, "/users/me", authorization_service.WorkspaceUpdate)

	rec := s.do(e, http.MethodGet, "/users/me", token)
	assert.Equal(s.T(), http.StatusOK, rec.Code)
}

// TestInvalidWorkspaceUUID verifies 400 when the workspace_id is not a valid UUID.
func (s *RequirePermissionsTestSuite) TestInvalidWorkspaceUUID() {
	_, token := s.createUserWithToken()
	e := s.setupEcho(http.MethodPut, "/workspaces/:workspace_id", authorization_service.WorkspaceUpdate)

	rec := s.do(e, http.MethodPut, "/workspaces/not-a-uuid", token)
	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
	assert.Equal(s.T(), "invalid workspace_id format", errBody(rec))
}

// TestUserNotInWorkspace verifies 403 with the "no access" message when the
// authenticated user has no role in the target workspace.
func (s *RequirePermissionsTestSuite) TestUserNotInWorkspace() {
	_, token := s.createUserWithToken()
	workspace := s.seeder.CreateWorkspace()
	e := s.setupEcho(http.MethodPut, "/workspaces/:workspace_id", authorization_service.WorkspaceUpdate)

	rec := s.do(e, http.MethodPut, "/workspaces/"+workspace.ID.String(), token)
	assert.Equal(s.T(), http.StatusForbidden, rec.Code)
	assert.Equal(s.T(), "user has no access to this workspace", errBody(rec))
}

// TestInsufficientPermissions verifies 403 with the "required permissions"
// message when the user has a workspace role but it lacks the needed permission.
func (s *RequirePermissionsTestSuite) TestInsufficientPermissions() {
	user, token := s.createUserWithToken()
	workspace := s.seeder.CreateWorkspace()
	s.seeder.CreateUserWorkspaceRelationship(user, workspace, "member")

	e := s.setupEcho(http.MethodPut, "/workspaces/:workspace_id", authorization_service.WorkspaceUpdate)

	rec := s.do(e, http.MethodPut, "/workspaces/"+workspace.ID.String(), token)
	assert.Equal(s.T(), http.StatusForbidden, rec.Code)
	assert.Equal(s.T(), "user does not have required permissions", errBody(rec))
}

// TestSufficientPermissions verifies the middleware passes through (200) when
// the user's role satisfies the required permissions.
func (s *RequirePermissionsTestSuite) TestSufficientPermissions() {
	user, token := s.createUserWithToken()
	workspace := s.seeder.CreateWorkspace()
	s.seeder.CreateUserWorkspaceRelationship(user, workspace, "owner")

	e := s.setupEcho(http.MethodPut, "/workspaces/:workspace_id", authorization_service.WorkspaceUpdate)

	rec := s.do(e, http.MethodPut, "/workspaces/"+workspace.ID.String(), token)
	assert.Equal(s.T(), http.StatusOK, rec.Code)
}

// TestGodModeBypassesPermissions verifies that a god-mode user can access any
// permission-protected route even without a workspace role.
func (s *RequirePermissionsTestSuite) TestGodModeBypassesPermissions() {
	user, token := s.createUserWithToken()
	s.Require().NoError(s.DB.Model(user).Update("god_mode", true).Error)

	workspace := s.seeder.CreateWorkspace()
	// Intentionally no workspace relationship for the god-mode user.
	e := s.setupEcho(http.MethodDelete, "/workspaces/:workspace_id", authorization_service.WorkspaceDelete)

	rec := s.do(e, http.MethodDelete, "/workspaces/"+workspace.ID.String(), token)
	assert.Equal(s.T(), http.StatusOK, rec.Code)
}

// TestNoRequiredPermission_AnyMemberPasses verifies that RequirePermissions()
// with no arguments only requires workspace membership (any role passes).
func (s *RequirePermissionsTestSuite) TestNoRequiredPermission_AnyMemberPasses() {
	user, token := s.createUserWithToken()
	workspace := s.seeder.CreateWorkspace()
	s.seeder.CreateUserWorkspaceRelationship(user, workspace, "member")

	e := s.setupEcho(http.MethodGet, "/workspaces/:workspace_id/newsletters" /* no permissions */)

	rec := s.do(e, http.MethodGet, "/workspaces/"+workspace.ID.String()+"/newsletters", token)
	assert.Equal(s.T(), http.StatusOK, rec.Code)
}

// TestNoRequiredPermission_NonMemberBlocked verifies that even a no-arg
// RequirePermissions blocks users who have no workspace role.
func (s *RequirePermissionsTestSuite) TestNoRequiredPermission_NonMemberBlocked() {
	_, token := s.createUserWithToken()
	workspace := s.seeder.CreateWorkspace()
	// No workspace relationship.
	e := s.setupEcho(http.MethodGet, "/workspaces/:workspace_id/newsletters" /* no permissions */)

	rec := s.do(e, http.MethodGet, "/workspaces/"+workspace.ID.String()+"/newsletters", token)
	assert.Equal(s.T(), http.StatusForbidden, rec.Code)
	assert.Equal(s.T(), "user has no access to this workspace", errBody(rec))
}

// ── Permission-matrix tests (role × permission) ──────────────────────────────

type permTestCase struct {
	role        string
	permission  authorization_service.Permission
	wantAllowed bool
}

// permMatrix is the complete expected permission matrix.
var permMatrix = []permTestCase{
	// WorkspaceUpdate
	{"owner", authorization_service.WorkspaceUpdate, true},
	{"admin", authorization_service.WorkspaceUpdate, true},
	{"member", authorization_service.WorkspaceUpdate, false},
	// WorkspaceDelete (owner-only)
	{"owner", authorization_service.WorkspaceDelete, true},
	{"admin", authorization_service.WorkspaceDelete, false},
	{"member", authorization_service.WorkspaceDelete, false},
	// WorkspaceManageUsers
	{"owner", authorization_service.WorkspaceManageUsers, true},
	{"admin", authorization_service.WorkspaceManageUsers, true},
	{"member", authorization_service.WorkspaceManageUsers, false},
	// WorkspaceManageIntegrations
	{"owner", authorization_service.WorkspaceManageIntegrations, true},
	{"admin", authorization_service.WorkspaceManageIntegrations, true},
	{"member", authorization_service.WorkspaceManageIntegrations, false},
	// WorkspaceManageRepositories
	{"owner", authorization_service.WorkspaceManageRepositories, true},
	{"admin", authorization_service.WorkspaceManageRepositories, true},
	{"member", authorization_service.WorkspaceManageRepositories, false},
}

// TestPermissionMatrix runs through every (role, permission) combination and
// asserts that the middleware produces the expected HTTP status.
func (s *RequirePermissionsTestSuite) TestPermissionMatrix() {
	for _, tc := range permMatrix {
		tc := tc // capture loop variable
		s.Run(tc.role+"_"+tc.permission.String(), func() {
			user, token := s.createUserWithToken()
			workspace := s.seeder.CreateWorkspace()
			s.seeder.CreateUserWorkspaceRelationship(user, workspace, tc.role)

			e := s.setupEcho(http.MethodPut, "/workspaces/:workspace_id", tc.permission)
			rec := s.do(e, http.MethodPut, "/workspaces/"+workspace.ID.String(), token)

			if tc.wantAllowed {
				assert.Equal(s.T(), http.StatusOK, rec.Code,
					"expected %s to be allowed for role=%s perm=%s", http.StatusOK, tc.role, tc.permission)
			} else {
				assert.Equal(s.T(), http.StatusForbidden, rec.Code,
					"expected %s to be forbidden for role=%s perm=%s", http.StatusForbidden, tc.role, tc.permission)
				assert.Equal(s.T(), "user does not have required permissions", errBody(rec))
			}
		})
	}
}

// TestMultiplePermissionsAllRequired verifies that ALL supplied permissions
// must be held — a user missing even one is rejected.
func (s *RequirePermissionsTestSuite) TestMultiplePermissionsAllRequired() {
	user, token := s.createUserWithToken()
	workspace := s.seeder.CreateWorkspace()
	// admin has WorkspaceManageUsers but NOT WorkspaceDelete
	s.seeder.CreateUserWorkspaceRelationship(user, workspace, "admin")

	e := s.setupEcho(http.MethodPut, "/workspaces/:workspace_id",
		authorization_service.WorkspaceManageUsers,
		authorization_service.WorkspaceDelete,
	)
	rec := s.do(e, http.MethodPut, "/workspaces/"+workspace.ID.String(), token)
	assert.Equal(s.T(), http.StatusForbidden, rec.Code)
	assert.Equal(s.T(), "user does not have required permissions", errBody(rec))
}

// TestContextPassthrough verifies that after middleware succeeds, the Echo
// context still has the DB transaction set (i.e. middleware doesn't replace it).
func (s *RequirePermissionsTestSuite) TestContextPassthrough() {
	user, token := s.createUserWithToken()
	workspace := s.seeder.CreateWorkspace()
	s.seeder.CreateUserWorkspaceRelationship(user, workspace, "owner")

	e := echo.New()
	e.HideBanner = true
	e.Use(db.WrapInTransaction(s.DB))
	e.Use(middlewares.InitializeServices())
	e.Use(middlewares.AuthMiddleware)

	var txFound bool
	e.Add(http.MethodPut, "/workspaces/:workspace_id", func(c echo.Context) error {
		txFound = context_utils.GetTransactionFromContext(c) != nil
		return c.JSON(http.StatusOK, nil)
	}, middlewares.RequirePermissions(authorization_service.WorkspaceUpdate))

	s.do(e, http.MethodPut, "/workspaces/"+workspace.ID.String(), token)
	assert.True(s.T(), txFound, "DB transaction should remain in context after RequirePermissions passes")
}
