package routes

import (
	"abstractions/clients/email"
	"abstractions/config"
	"abstractions/db"
	s "abstractions/server"
	"abstractions/server/handlers"
	"abstractions/server/middlewares"
	authorization_service "abstractions/services/authorization"
	"fmt"
	"strconv"

	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func ConfigureRoutes(server *s.Server) {
	authHandler := handlers.NewAuthHandler(server)
	usersHandler := handlers.NewUsersHandler(server)
	subscriptionHandler := handlers.NewSubscriptionHandler(server)
	adminHandler := handlers.NewAdminHandler(server)
	workspaceHandler := handlers.NewWorkspaceHandler(server)
	invitationHandler := handlers.NewInvitationHandler(server)
	integrationHandler := handlers.NewIntegrationHandler(server)
	repositoryHandler := handlers.NewRepositoryHandler(server)

	// Create a shared asynq client for enqueueing tasks from the HTTP server
	redisAddr := server.Config.Environment.RedisHost + ":" + strconv.Itoa(server.Config.Environment.RedisPort)
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	webhookHandler := handlers.NewWebhookHandler(server, asynqClient)
	unsubscribeHandler := handlers.NewUnsubscribeHandler(server.DB)

	server.Echo.Use(echomw.Logger())
	server.Echo.Use(echomw.CORS())
	server.Echo.Use(db.WrapInTransaction(server.DB))
	server.Echo.Use(middlewares.InitializeServices()) // Add services middleware after transaction middleware

	server.Echo.GET("/email", func(c echo.Context) error {
		emailHtml, emailText, err := email.BuildEmail(email.SendGenericEmailParams{
			Title:         "Welcome to Abstractions",
			BodyHTML:      "<p>Thank you for signing up to Abstractions. We are excited to have you on board.</p>",
			BodyPlainText: "Thank you for signing up to Abstractions. We are excited to have you on board.",
			CTAText:       "Get Started",
			CTALink:       "https://tryabstractions.com",
			Reason:        "Welcome",
		})

		fmt.Println(emailHtml, emailText, err)

		return c.HTML(200, emailHtml)
	})

	server.Echo.Use(middlewares.AuthMiddleware)

	// Webhook routes (public — verified by HMAC signature, not by session token)
	server.Echo.POST("/webhooks/github", webhookHandler.ReceiveGitHubWebhook)

	// Unsubscribe route (public — verified by signed JWT in query param)
	server.Echo.GET("/unsubscribe", unsubscribeHandler.Unsubscribe)

	if !config.IsProduction(server.Config) {
		server.Echo.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	// public routes
	server.Echo.GET("/auth/google", authHandler.GoogleAuthorize)
	server.Echo.GET("/auth/google/callback", authHandler.GoogleCallback)
	server.Echo.POST("/auth/signup", authHandler.Signup)
	server.Echo.POST("/auth/login", authHandler.Login)
	server.Echo.POST("/users/verify", authHandler.VerifyEmail)
	server.Echo.POST("/users/password/reset/confirm", authHandler.ResetPasswordConfirm)
	// authenticated routes
	server.Echo.POST("/users/verify/resend", authHandler.ResendVerificationEmail)
	server.Echo.POST("/users/password/reset", authHandler.ResetPassword)
	server.Echo.GET("/users/me", usersHandler.GetCurrentUser)
	server.Echo.PATCH("/users/me", usersHandler.UpdateCurrentUser)
	server.Echo.DELETE("/users/me", usersHandler.DeleteCurrentUser)
	server.Echo.POST("/users/me/avatar", usersHandler.UpdateUserAvatar)
	server.Echo.PUT("/users/me/password", usersHandler.UpdateUserPassword)
	server.Echo.POST("/users/me/workspaces/switch", usersHandler.SwitchWorkspace)

	// Regular user workspace subscription endpoint
	server.Echo.GET("/workspaces/:workspace_id/subscription", subscriptionHandler.GetActiveWorkspaceSubscription)

	// User workspace routes - basic routes without auth
	server.Echo.GET("/workspaces", workspaceHandler.GetUserWorkspaces)
	server.Echo.POST("/workspaces", workspaceHandler.CreateWorkspace)

	// Protected workspace routes
	server.Echo.PUT("/workspaces/:workspace_id", workspaceHandler.UpdateWorkspace, middlewares.RequirePermissions(authorization_service.WorkspaceUpdate))
	server.Echo.DELETE("/workspaces/:workspace_id", workspaceHandler.DeleteWorkspace, middlewares.RequirePermissions(authorization_service.WorkspaceDelete))
	server.Echo.POST("/workspaces/:workspace_id/logo", workspaceHandler.UploadLogo, middlewares.RequirePermissions(authorization_service.WorkspaceUpdate))

	// Workspace member management routes
	server.Echo.GET("/workspaces/:workspace_id/users", workspaceHandler.GetWorkspaceUsers)
	server.Echo.PUT("/workspaces/:workspace_id/users", workspaceHandler.UpdateWorkspaceUserRole, middlewares.RequirePermissions(authorization_service.WorkspaceManageUsers))
	server.Echo.DELETE("/workspaces/:workspace_id/users", workspaceHandler.RemoveWorkspaceUser, middlewares.RequirePermissions(authorization_service.WorkspaceManageUsers))

	server.Echo.GET("/workspaces/:workspace_id/roles", workspaceHandler.GetRoles)

	// Workspace Invitation Routes
	server.Echo.GET("/workspaces/:workspace_id/invitations", invitationHandler.GetWorkspaceInvitations, middlewares.RequirePermissions(authorization_service.WorkspaceManageUsers))
	server.Echo.POST("/workspaces/:workspace_id/invitations", invitationHandler.SendInvitations, middlewares.RequirePermissions(authorization_service.WorkspaceManageUsers))
	server.Echo.POST("/workspaces/:workspace_id/invitations/:invitation_id/resend", invitationHandler.ResendInvitation, middlewares.RequirePermissions(authorization_service.WorkspaceManageUsers))
	server.Echo.POST("/workspaces/:workspace_id/invitations/:invitation_id/revoke", invitationHandler.RevokeInvitation, middlewares.RequirePermissions(authorization_service.WorkspaceManageUsers))

	// User Invitation Routes - no special permissions needed
	server.Echo.GET("/users/me/invitations", invitationHandler.GetUserPendingInvitations)
	server.Echo.POST("/users/me/invitations/:invitation_id/accept", invitationHandler.AcceptInvitation)
	server.Echo.POST("/users/me/invitations/:invitation_id/ignore", invitationHandler.IgnoreInvitation)

	// Integration Routes
	server.Echo.GET("/integrations/types", integrationHandler.GetIntegrationTypes)
	server.Echo.POST("/workspaces/:workspace_id/integrations/:type/setup", integrationHandler.IntegrationSetup, middlewares.RequirePermissions(authorization_service.WorkspaceManageUsers))
	server.Echo.GET("/integrations/callback/:type", integrationHandler.IntegrationCallback)
	server.Echo.GET("/workspaces/:workspace_id/integrations", integrationHandler.GetIntegrations)
	server.Echo.GET("/workspaces/:workspace_id/integrations/:type", integrationHandler.GetIntegrationByType)
	server.Echo.DELETE("/workspaces/:workspace_id/integrations/:integration_id", integrationHandler.DeleteIntegration, middlewares.RequirePermissions(authorization_service.WorkspaceManageUsers))
	server.Echo.PUT("/workspaces/:workspace_id/integrations/:integration_id", integrationHandler.UpdateIntegration, middlewares.RequirePermissions(authorization_service.WorkspaceManageUsers))

	// Repository Routes
	server.Echo.GET("/workspaces/:workspace_id/github/repos", repositoryHandler.ListGitHubRepos)
	server.Echo.GET("/workspaces/:workspace_id/repositories", repositoryHandler.ListRepositories)
	server.Echo.POST("/workspaces/:workspace_id/repositories", repositoryHandler.ConnectRepository, middlewares.RequirePermissions(authorization_service.WorkspaceManageUsers))

	// admin routes
	server.Echo.Use(middlewares.AdminAuthMiddleware, middlewares.RequirePermissions(authorization_service.WorkspaceManageUsers))
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

	// Workspace Admin Routes
	server.Echo.GET("/admin/workspaces", adminHandler.GetWorkspaces)
	server.Echo.POST("/admin/workspaces", adminHandler.CreateWorkspace)
	server.Echo.GET("/admin/workspaces/:workspace_id", adminHandler.GetWorkspace)
	server.Echo.PUT("/admin/workspaces/:workspace_id", adminHandler.UpdateWorkspace)
	server.Echo.DELETE("/admin/workspaces", adminHandler.BulkDeleteWorkspaces)
	server.Echo.POST("/admin/workspaces/recover", adminHandler.BulkRecoverWorkspaces)

	// Workspace Subscription Routes
	server.Echo.GET("/admin/workspaces/:workspace_id/subscription", adminHandler.GetWorkspaceSubscription)
	server.Echo.POST("/admin/workspaces/:workspace_id/subscription", adminHandler.CreateWorkspaceSubscription)
	server.Echo.PUT("/admin/workspaces/:workspace_id/subscription", adminHandler.UpdateWorkspaceSubscription)
	server.Echo.POST("/admin/workspaces/:workspace_id/subscription/disable", adminHandler.DisableWorkspaceSubscription)

	// Workspace User Management Routes
	server.Echo.GET("/admin/workspaces/:workspace_id/users", adminHandler.GetWorkspaceUsers)
	server.Echo.POST("/admin/workspaces/:workspace_id/users", adminHandler.AddUserToWorkspace)
	server.Echo.PUT("/admin/workspaces/:workspace_id/users/:userId", adminHandler.UpdateWorkspaceUserRole)
	server.Echo.DELETE("/admin/workspaces/:workspace_id/users/:userId", adminHandler.RemoveUserFromWorkspace)

	server.Echo.GET("/admin/stats", adminHandler.GetStats)
	server.Echo.GET("/admin/historical-data", adminHandler.GetHistoricalData)
}
