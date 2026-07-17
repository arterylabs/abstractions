package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"abstractions/db/models"
	"abstractions/server"
	"abstractions/server/routes"
	authorization_service "abstractions/services/authorization"
	github_integration "abstractions/services/integrations/core/github"
	user_service "abstractions/services/user"
	"abstractions/tests/testutils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// SecurityAccessControlTestSuite exercises workspace-scoped access control on
// routes that take a workspace_id from the URL. It asserts that a user who is
// NOT a member of a workspace cannot read that workspace's data (cross-tenant
// IDOR), while a legitimate member still can.
type SecurityAccessControlTestSuite struct {
	testutils.TestSuite
	seeder      *testutils.Seeder
	userService *user_service.Service
}

func TestSecurityAccessControlSuite(t *testing.T) {
	suite.Run(t, new(SecurityAccessControlTestSuite))
}

func (s *SecurityAccessControlTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.userService = user_service.NewService(s.DB)
}

// router builds an Echo instance wired with the real production route table
// (routes.ConfigureRoutes), backed by the test database. This ensures the tests
// exercise the actual middleware attached to each route in routes.go.
func (s *SecurityAccessControlTestSuite) router() *echo.Echo {
	s.T().Helper()
	srv := &server.Server{
		Echo:   echo.New(),
		DB:     s.DB,
		Config: s.Config,
	}
	routes.ConfigureRoutes(srv)
	return srv.Echo
}

// tokenFor creates a verified user and returns a usable Bearer token.
func (s *SecurityAccessControlTestSuite) tokenFor(email string) (*models.User, string) {
	s.T().Helper()
	user, err := s.userService.CreateVerifiedUser(s.Ctx, user_service.CreateVerifiedUserParams{
		Email: email,
		Name:  "Test User",
	})
	s.Require().NoError(err)

	session, err := s.userService.CreateUserSession(s.Ctx, user.ID, &user_service.CreateSessionParams{
		UserAgent: "test-agent",
		IP:        "127.0.0.1",
	})
	s.Require().NoError(err)
	return user, session.Token
}

func (s *SecurityAccessControlTestSuite) get(e *echo.Echo, url, token string) *httptest.ResponseRecorder {
	s.T().Helper()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

// seedWorkspaceWithData creates a workspace owned by owner, containing a
// repository and a GitHub integration.
func (s *SecurityAccessControlTestSuite) seedWorkspaceWithData(owner *models.User) (*models.Workspace, *models.Repository) {
	s.T().Helper()
	ws := s.seeder.CreateWorkspace()
	s.seeder.CreateUserWorkspaceRelationship(owner, ws, authorization_service.RoleOwner.String())
	repo := s.seeder.CreateRepository(ws.ID, owner.ID)

	integration := &models.Integration{
		ID:          uuid.New(),
		Type:        github_integration.IntegrationTypeGitHub,
		WorkspaceID: ws.ID,
		Token:       "encrypted-token",
		CreatedBy:   owner.ID,
		UpdatedBy:   owner.ID,
	}
	s.Require().NoError(s.DB.Create(integration).Error)

	return ws, repo
}

// TestNonMemberCannotReadRepositories is the core cross-tenant IDOR test.
// A user with no membership in the target workspace must not be able to read
// its repositories, a single repository, or its integrations.
func (s *SecurityAccessControlTestSuite) TestNonMemberCannotReadWorkspaceData() {
	owner, _ := s.tokenFor("owner-" + uuid.NewString() + "@example.com")
	ws, repo := s.seedWorkspaceWithData(owner)

	// Attacker is a fully authenticated user, but NOT a member of ws.
	_, attackerToken := s.tokenFor("attacker-" + uuid.NewString() + "@example.com")

	e := s.router()

	cases := []struct {
		name string
		url  string
	}{
		{"list repositories", "/workspaces/" + ws.ID.String() + "/repositories"},
		{"get repository", "/workspaces/" + ws.ID.String() + "/repositories/" + repo.ID.String()},
		{"get repository status", "/workspaces/" + ws.ID.String() + "/repositories/" + repo.ID.String() + "/status"},
		{"get repository stats", "/workspaces/" + ws.ID.String() + "/repositories/" + repo.ID.String() + "/stats"},
		{"list integrations", "/workspaces/" + ws.ID.String() + "/integrations"},
		{"get integration by type", "/workspaces/" + ws.ID.String() + "/integrations/github"},
	}

	for _, tc := range cases {
		rec := s.get(e, tc.url, attackerToken)
		assert.Equalf(s.T(), http.StatusForbidden, rec.Code,
			"non-member should be denied (403) on %q, got %d: %s", tc.name, rec.Code, rec.Body.String())
	}
}

// TestMemberCanReadWorkspaceData is the positive control: a legitimate member
// must retain access after the access-control fix.
func (s *SecurityAccessControlTestSuite) TestMemberCanReadWorkspaceData() {
	owner, ownerToken := s.tokenFor("owner-" + uuid.NewString() + "@example.com")
	ws, repo := s.seedWorkspaceWithData(owner)

	e := s.router()

	for _, url := range []string{
		"/workspaces/" + ws.ID.String() + "/repositories",
		"/workspaces/" + ws.ID.String() + "/repositories/" + repo.ID.String(),
		"/workspaces/" + ws.ID.String() + "/integrations",
	} {
		rec := s.get(e, url, ownerToken)
		assert.Equalf(s.T(), http.StatusOK, rec.Code,
			"member should be allowed (200) on %q, got %d: %s", url, rec.Code, rec.Body.String())
	}
}

// TestSwitchDefaultWorkspaceRequiresMembership verifies that a user cannot set
// their default workspace to a workspace they are not a member of. The current
// implementation only checks that the workspace exists, so this fails until a
// membership check is added.
func (s *SecurityAccessControlTestSuite) TestSwitchDefaultWorkspaceRequiresMembership() {
	// A workspace owned by someone else.
	other, _ := s.tokenFor("other-" + uuid.NewString() + "@example.com")
	foreignWS := s.seeder.CreateWorkspace()
	s.seeder.CreateUserWorkspaceRelationship(other, foreignWS, authorization_service.RoleOwner.String())

	// The attacker is a member of nothing relevant here.
	attacker, _ := s.tokenFor("attacker-" + uuid.NewString() + "@example.com")

	err := s.userService.UpdateDefaultWorkspace(s.Ctx, attacker.ID, foreignWS.ID)
	assert.Error(s.T(), err, "setting default workspace to a non-member workspace must be rejected")
}

// TestSwitchDefaultWorkspaceAllowsMember is the positive control for the check above.
func (s *SecurityAccessControlTestSuite) TestSwitchDefaultWorkspaceAllowsMember() {
	user, _ := s.tokenFor("member-" + uuid.NewString() + "@example.com")
	ws := s.seeder.CreateWorkspace()
	s.seeder.CreateUserWorkspaceRelationship(user, ws, authorization_service.RoleMember.String())

	err := s.userService.UpdateDefaultWorkspace(s.Ctx, user.ID, ws.ID)
	assert.NoError(s.T(), err, "a member must be able to set that workspace as default")
}
