package services

import (
	"abstractions/config"
	"abstractions/repositories"
	analytics_service "abstractions/services/analytics"
	authenticator_service "abstractions/services/authenticator"
	authorization_service "abstractions/services/authorization"
	integration_service "abstractions/services/integrations"
	invitation_service "abstractions/services/invitation"
	repository_service "abstractions/services/repository"
	session_service "abstractions/services/session"
	subscription_service "abstractions/services/subscription"
	test_service "abstractions/services/test"
	token_service "abstractions/services/token"
	user_service "abstractions/services/user"
	workspace_service "abstractions/services/workspace"
	workspace_configuration_service "abstractions/services/workspace_configuration"

	"gorm.io/gorm"
)

type Services struct {
	Analytics              *analytics_service.Service
	Authenticator          *authenticator_service.Service
	Authorization          *authorization_service.Service
	Integration            *integration_service.Service
	Invitation             *invitation_service.Service
	Repository             *repository_service.Service
	Session                *session_service.Service
	Subscription           *subscription_service.Service
	Test                   *test_service.Service
	Token                  *token_service.Service
	User                   *user_service.Service
	Workspace              *workspace_service.Service
	WorkspaceConfiguration *workspace_configuration_service.Service
}

func New(db *gorm.DB) *Services {
	encryptionKey := ""
	if config.Environment != nil {
		encryptionKey = config.Environment.EncryptionKey32Bit
	}
	return &Services{
		Analytics:              analytics_service.NewService(db),
		Authenticator:          authenticator_service.NewService(db),
		Authorization:          authorization_service.NewService(db),
		Integration:            integration_service.NewService(db),
		Invitation:             invitation_service.NewService(db),
		Repository:             repository_service.NewService(db),
		Session:                session_service.NewService(db),
		Subscription:           subscription_service.NewService(db),
		Test:                   test_service.NewService(db),
		Token:                  token_service.NewService(db),
		User:                   user_service.NewService(db),
		Workspace:              workspace_service.NewService(db),
		WorkspaceConfiguration: workspace_configuration_service.NewService(repositories.NewWorkspaceConfigurationRepository(db), encryptionKey),
	}
}
