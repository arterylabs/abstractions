package routes

import (
	"abstractions/config"
	"abstractions/db"
	s "abstractions/server"
	"abstractions/server/handlers"
	"abstractions/server/middlewares"
	authorization_service "abstractions/services/authorization"
	"abstractions/workers/jobs"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"golang.org/x/time/rate"
)

// authRateLimiter returns a per-IP token-bucket rate limiter middleware.
// requestsPerMinute controls the sustained rate; burst is set to 1/4 of that.
// When skip is true (non-production environments) the middleware is a no-op.
func authRateLimiter(requestsPerMinute int, skip bool) echo.MiddlewareFunc {
	store := echomw.NewRateLimiterMemoryStoreWithConfig(
		echomw.RateLimiterMemoryStoreConfig{
			Rate:      rate.Limit(float64(requestsPerMinute) / 60.0),
			Burst:     requestsPerMinute / 4,
			ExpiresIn: 3 * time.Minute,
		},
	)
	return echomw.RateLimiterWithConfig(echomw.RateLimiterConfig{
		Skipper: func(c echo.Context) bool { return skip },
		Store:   store,
		IdentifierExtractor: func(c echo.Context) (string, error) {
			return c.RealIP(), nil
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return echo.NewHTTPError(http.StatusTooManyRequests, "too many requests")
		},
		DenyHandler: func(c echo.Context, identifier string, err error) error {
			return echo.NewHTTPError(http.StatusTooManyRequests, "too many requests, please slow down")
		},
	})
}

func ConfigureRoutes(server *s.Server) {
	jobClient := jobs.NewClient(config.Environment.RedisHost + ":" + strconv.Itoa(config.Environment.RedisPort))

	// Instantiate all handlers up front for clarity.
	authHandler := handlers.NewAuthHandler(server)
	usersHandler := handlers.NewUsersHandler(server)
	subscriptionHandler := handlers.NewSubscriptionHandler(server)
	adminHandler := handlers.NewAdminHandler(server)
	workspaceHandler := handlers.NewWorkspaceHandler(server)
	invitationHandler := handlers.NewInvitationHandler(server)
	integrationHandler := handlers.NewIntegrationHandler(server)
	repositoryHandler := handlers.NewRepositoryHandler(server, jobClient)
	newsletterHandler := handlers.NewNewsletterHandler(server, jobClient)
	newsletterTopicHandler := handlers.NewNewsletterTopicHandler(server)

	// Global middleware stack (order matters).
	server.Echo.Use(middlewares.PanicRecoveryMiddleware)
	server.Echo.Use(echomw.Logger())
	server.Echo.Use(echomw.CORS())
	server.Echo.Use(db.WrapInTransaction(server.DB))
	server.Echo.Use(middlewares.InitializeServices())
	server.Echo.Use(middlewares.AuthMiddleware(server.Config))

	if !config.IsProduction(server.Config) {
		server.Echo.GET("/swagger/*", echoSwagger.WrapHandler)
		testHandler := handlers.NewTestHandler(server)
		server.Echo.POST("/test/cleanup", testHandler.CleanupTestData)
		server.Echo.POST("/test/verify_user", testHandler.VerifyTestUser)
	}

	// Per-route rate limiters for auth endpoints (IP-based token bucket).
	// Skipped in non-production environments to avoid friction during development.
	isDev := !config.IsProduction(server.Config)
	loginLimiter := authRateLimiter(10, isDev)    // 10 attempts/min per IP
	signupLimiter := authRateLimiter(5, isDev)    // 5 signups/min per IP
	resetLimiter := authRateLimiter(5, isDev)     // 5 reset requests/min per IP
	verifyLimiter := authRateLimiter(10, isDev)   // 10 verify attempts/min per IP
	resendLimiter := authRateLimiter(2, isDev)    // 2 resend attempts/min per IP
	exchangeLimiter := authRateLimiter(10, isDev) // 10 OAuth exchanges/min per IP

	// ── Public / unauthenticated routes ─────────────────────────────────────
	server.Echo.GET("/auth/google", authHandler.GoogleAuthorize)
	server.Echo.GET("/auth/github", authHandler.GitHubAuthorize)
	// Unified OAuth callback — flow is determined by the state.Flow field.
	server.Echo.GET("/oauth/callback/google", authHandler.GoogleOAuthCallback)
	server.Echo.GET("/oauth/callback/github", authHandler.GitHubOAuthCallback)
	server.Echo.POST("/auth/exchange", authHandler.ExchangeCode, exchangeLimiter)
	server.Echo.POST("/auth/signup", authHandler.Signup, signupLimiter)
	server.Echo.POST("/auth/login", authHandler.Login, loginLimiter)
	server.Echo.POST("/users/verify", authHandler.VerifyEmail, verifyLimiter)
	server.Echo.POST("/users/password/reset/confirm", authHandler.ResetPasswordConfirm, resetLimiter)

	// ── Authenticated routes (no workspace permission required) ──────────────
	server.Echo.POST("/users/verify/resend", authHandler.ResendVerificationEmail, resendLimiter)
	server.Echo.POST("/users/password/reset", authHandler.ResetPassword, resetLimiter)
	server.Echo.GET("/users/me", usersHandler.GetCurrentUser)
	server.Echo.PATCH("/users/me", usersHandler.UpdateCurrentUser)
	server.Echo.DELETE("/users/me", usersHandler.DeleteCurrentUser)
	server.Echo.POST("/users/me/avatar", usersHandler.UpdateUserAvatar)
	server.Echo.PUT("/users/me/password", usersHandler.UpdateUserPassword)
	server.Echo.POST("/users/me/workspaces/switch", usersHandler.SwitchWorkspace)

	server.Echo.GET("/workspaces", workspaceHandler.GetUserWorkspaces)
	server.Echo.POST("/workspaces", workspaceHandler.CreateWorkspace, middlewares.RequireVerifiedAccount)
	server.Echo.GET("/ai/models", workspaceHandler.GetAIModels)

	server.Echo.GET("/users/me/invitations", invitationHandler.GetUserPendingInvitations)
	server.Echo.POST("/users/me/invitations/:invitation_id/accept", invitationHandler.AcceptInvitation)
	server.Echo.POST("/users/me/invitations/:invitation_id/ignore", invitationHandler.IgnoreInvitation)

	server.Echo.GET("/integrations/types", integrationHandler.GetIntegrationTypes)
	server.Echo.GET("/workspaces/:workspace_id/integrations", integrationHandler.GetIntegrations)
	server.Echo.GET("/workspaces/:workspace_id/integrations/:type", integrationHandler.GetIntegrationByType)

	server.Echo.GET("/workspaces/:workspace_id/github/repos", repositoryHandler.ListGitHubRepos)
	server.Echo.GET("/workspaces/:workspace_id/repositories", repositoryHandler.ListRepositories)
	server.Echo.GET("/workspaces/:workspace_id/repositories/:repo_id", repositoryHandler.GetRepository)
	server.Echo.GET("/workspaces/:workspace_id/repositories/:repo_id/status", repositoryHandler.GetRepositoryStatus)
	server.Echo.GET("/workspaces/:workspace_id/repositories/:repo_id/stats", repositoryHandler.GetRepositoryStats)

	// ── WorkspaceUpdate permission ───────────────────────────────────────────
	wsUpdate := server.Echo.Group("", middlewares.RequirePermissions(authorization_service.WorkspaceUpdate))
	wsUpdate.PUT("/workspaces/:workspace_id", workspaceHandler.UpdateWorkspace)
	wsUpdate.POST("/workspaces/:workspace_id/logo", workspaceHandler.UploadLogo)
	wsUpdate.GET("/workspaces/:workspace_id/ai_settings", workspaceHandler.GetWorkspaceAISettings)
	wsUpdate.PUT("/workspaces/:workspace_id/ai_settings", workspaceHandler.UpdateWorkspaceAISettings)

	// ── WorkspaceDelete permission ───────────────────────────────────────────
	wsDelete := server.Echo.Group("", middlewares.RequireVerifiedAccount, middlewares.RequirePermissions(authorization_service.WorkspaceDelete))
	wsDelete.DELETE("/workspaces/:workspace_id", workspaceHandler.DeleteWorkspace)

	// ── WorkspaceManageUsers permission ──────────────────────────────────────
	wsManageUsers := server.Echo.Group("", middlewares.RequireVerifiedAccount, middlewares.RequirePermissions(authorization_service.WorkspaceManageUsers))
	wsManageUsers.PUT("/workspaces/:workspace_id/users", workspaceHandler.UpdateWorkspaceUserRole)
	wsManageUsers.DELETE("/workspaces/:workspace_id/users", workspaceHandler.RemoveWorkspaceUser)
	wsManageUsers.GET("/workspaces/:workspace_id/invitations", invitationHandler.GetWorkspaceInvitations)
	wsManageUsers.POST("/workspaces/:workspace_id/invitations", invitationHandler.SendInvitations)
	wsManageUsers.POST("/workspaces/:workspace_id/invitations/:invitation_id/resend", invitationHandler.ResendInvitation)
	wsManageUsers.POST("/workspaces/:workspace_id/invitations/:invitation_id/revoke", invitationHandler.RevokeInvitation)

	// ── WorkspaceManageIntegrations permission ───────────────────────────────
	wsManageIntegrations := server.Echo.Group("", middlewares.RequireVerifiedAccount, middlewares.RequirePermissions(authorization_service.WorkspaceManageIntegrations))
	wsManageIntegrations.POST("/workspaces/:workspace_id/integrations/:type/setup", integrationHandler.IntegrationSetup)
	wsManageIntegrations.DELETE("/workspaces/:workspace_id/integrations/:integration_id", integrationHandler.DeleteIntegration)
	wsManageIntegrations.PUT("/workspaces/:workspace_id/integrations/:integration_id", integrationHandler.UpdateIntegration)
	wsManageIntegrations.POST("/workspaces/:workspace_id/repositories", repositoryHandler.ConnectRepository)
	wsManageIntegrations.PUT("/workspaces/:workspace_id/repositories/:repo_id", repositoryHandler.UpdateRepositorySettings)
	wsManageIntegrations.PATCH("/workspaces/:workspace_id/repositories/:repo_id", repositoryHandler.UpdateRepositoryDescription)
	wsManageIntegrations.DELETE("/workspaces/:workspace_id/repositories/:repo_id", repositoryHandler.DeleteRepository)
	wsManageIntegrations.POST("/workspaces/:workspace_id/repositories/:repo_id/trigger_index", repositoryHandler.TriggerIndex)

	// ── WorkspaceManageRepositories permission ───────────────────────────────
	wsManageRepos := server.Echo.Group("", middlewares.RequireVerifiedAccount, middlewares.RequirePermissions(authorization_service.WorkspaceManageRepositories))
	wsManageRepos.POST("/workspaces/:workspace_id/repositories/:repo_id/trigger_newsletter", newsletterHandler.TriggerNewsletter)
	wsManageRepos.POST("/workspaces/:workspace_id/repositories/:repo_id/newsletters/:newsletter_id/resend", newsletterHandler.ResendNewsletter)
	wsManageRepos.POST("/workspaces/:workspace_id/repositories/:repo_id/newsletter_topics", newsletterTopicHandler.CreateNewsletterTopic)
	wsManageRepos.PUT("/workspaces/:workspace_id/repositories/:repo_id/newsletter_topics/:topic_id", newsletterTopicHandler.UpdateNewsletterTopic)
	wsManageRepos.PUT("/workspaces/:workspace_id/repositories/:repo_id/newsletter_topics/:topic_id/reorder", newsletterTopicHandler.ReorderNewsletterTopic)
	wsManageRepos.DELETE("/workspaces/:workspace_id/repositories/:repo_id/newsletter_topics/:topic_id", newsletterTopicHandler.DeleteNewsletterTopic)

	// ── Workspace membership required (read-only, no specific permission) ────
	wsMember := server.Echo.Group("", middlewares.RequirePermissions())
	wsMember.GET("/workspaces/:workspace_id/subscription", subscriptionHandler.GetActiveWorkspaceSubscription)
	wsMember.GET("/workspaces/:workspace_id/users", workspaceHandler.GetWorkspaceUsers)
	wsMember.GET("/workspaces/:workspace_id/roles", workspaceHandler.GetRoles)
	wsMember.GET("/workspaces/:workspace_id/health", workspaceHandler.GetWorkspaceHealth)
	wsMember.GET("/workspaces/:workspace_id/getting_started", workspaceHandler.GetGettingStarted)
	wsMember.GET("/workspaces/:workspace_id/repositories/:repo_id/newsletters", newsletterHandler.ListNewsletters)
	wsMember.GET("/workspaces/:workspace_id/repositories/:repo_id/newsletters/:newsletter_id", newsletterHandler.GetNewsletter)
	wsMember.GET("/workspaces/:workspace_id/repositories/:repo_id/newsletters/:newsletter_id/email_content", newsletterHandler.GetEmailContent)
	wsMember.GET("/workspaces/:workspace_id/repositories/:repo_id/newsletter_topics", newsletterTopicHandler.ListNewsletterTopics)

	// ── Admin routes (require AdminAuthMiddleware + god-mode flag) ───────────
	server.Echo.Use(middlewares.AdminAuthMiddleware)
	server.Echo.GET("/admin/users", adminHandler.GetUsers)
	server.Echo.POST("/admin/users", adminHandler.CreateUser)
	server.Echo.GET("/admin/users/:id", adminHandler.GetUser)
	server.Echo.PUT("/admin/users/:id", adminHandler.UpdateUser)
	server.Echo.POST("/admin/users/verify", adminHandler.VerifyUsers)
	server.Echo.DELETE("/admin/users", adminHandler.DeleteUsers)
	server.Echo.POST("/admin/users/recover", adminHandler.RecoverUsers)
	server.Echo.POST("/admin/users/impersonate", adminHandler.ImpersonateUser)
	server.Echo.POST("/admin/users/exit_impersonation", adminHandler.ExitImpersonation)
	server.Echo.GET("/admin/users/:user_id/workspaces", adminHandler.GetUserWorkspacesWithRole)
	server.Echo.GET("/admin/workspaces", adminHandler.GetWorkspaces)
	server.Echo.POST("/admin/workspaces", adminHandler.CreateWorkspace)
	server.Echo.GET("/admin/workspaces/:workspace_id", adminHandler.GetWorkspace)
	server.Echo.PUT("/admin/workspaces/:workspace_id", adminHandler.UpdateWorkspace)
	server.Echo.DELETE("/admin/workspaces", adminHandler.BulkDeleteWorkspaces)
	server.Echo.POST("/admin/workspaces/recover", adminHandler.BulkRecoverWorkspaces)
	server.Echo.GET("/admin/workspaces/:workspace_id/subscription", adminHandler.GetWorkspaceSubscription)
	server.Echo.POST("/admin/workspaces/:workspace_id/subscription", adminHandler.CreateWorkspaceSubscription)
	server.Echo.PUT("/admin/workspaces/:workspace_id/subscription", adminHandler.UpdateWorkspaceSubscription)
	server.Echo.POST("/admin/workspaces/:workspace_id/subscription/disable", adminHandler.DisableWorkspaceSubscription)
	server.Echo.GET("/admin/workspaces/:workspace_id/users", adminHandler.GetWorkspaceUsers)
	server.Echo.POST("/admin/workspaces/:workspace_id/users", adminHandler.AddUserToWorkspace)
	server.Echo.PUT("/admin/workspaces/:workspace_id/users/:user_id", adminHandler.UpdateWorkspaceUserRole)
	server.Echo.DELETE("/admin/workspaces/:workspace_id/users/:user_id", adminHandler.RemoveUserFromWorkspace)
	server.Echo.GET("/admin/stats", adminHandler.GetStats)
	server.Echo.GET("/admin/historical_data", adminHandler.GetHistoricalData)
}
