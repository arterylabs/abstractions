package services

import (
	analytics_service "abstractions/services/analytics"
	authenticator_service "abstractions/services/authenticator"
	authorization_service "abstractions/services/authorization"
	integration_service "abstractions/services/integrations"
	invitation_service "abstractions/services/invitation"
	session_service "abstractions/services/session"
	subscription_service "abstractions/services/subscription"
	token_service "abstractions/services/token"
	user_service "abstractions/services/user"
	workspace_service "abstractions/services/workspace"

	"gorm.io/gorm"
)

type Services struct {
	Analytics     *analytics_service.Service
	Authenticator *authenticator_service.Service
	Authorization *authorization_service.Service
	Integration   *integration_service.Service
	Invitation    *invitation_service.Service
	Session       *session_service.Service
	Subscription  *subscription_service.Service
	Token         *token_service.Service
	User          *user_service.Service
	Workspace     *workspace_service.Service
}

func New(db *gorm.DB) *Services {
	return &Services{
		Analytics:     analytics_service.NewAnalyticsService(db),
		Authenticator: authenticator_service.NewAuthenticatorService(db),
		Authorization: authorization_service.NewAuthorizationService(db),
		Integration:   integration_service.NewIntegrationService(db),
		Invitation:    invitation_service.NewInvitationService(db),
		Session:       session_service.NewSessionService(db),
		Subscription:  subscription_service.NewSubscriptionService(db),
		Token:         token_service.NewTokenService(db),
		User:          user_service.NewUserService(db),
		Workspace:     workspace_service.NewWorkspaceService(db),
	}
}
