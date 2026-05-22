package jobs

import (
	ai_client "abstractions/clients/ai"
	"abstractions/config"
	"abstractions/repositories"
	workspace_configuration_service "abstractions/services/workspace_configuration"
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// newWorkspaceAIProvider creates an AIProvider resolved against the workspace configuration for the
// given workspace. Returns an error if no Anthropic API key is available (workspace or system default).
func newWorkspaceAIProvider(ctx context.Context, db *gorm.DB, workspaceID uuid.UUID) (ai_client.AIProvider, error) {
	encryptionKey := ""
	if config.Environment != nil {
		encryptionKey = config.Environment.EncryptionKey32Bit
	}
	cfgRepo := repositories.NewWorkspaceConfigurationRepository(db)
	cfgSvc := workspace_configuration_service.NewService(cfgRepo, encryptionKey)

	anthropicKey, _ := cfgSvc.Get(ctx, workspaceID, workspace_configuration_service.KeyAnthropicAPIKey)
	openAIKey, _ := cfgSvc.Get(ctx, workspaceID, workspace_configuration_service.KeyOpenAIAPIKey)
	indexingModel, _ := cfgSvc.Get(ctx, workspaceID, workspace_configuration_service.KeyIndexingModel)
	newsletterModel, _ := cfgSvc.Get(ctx, workspaceID, workspace_configuration_service.KeyNewsletterModel)
	embeddingsModel, _ := cfgSvc.Get(ctx, workspaceID, workspace_configuration_service.KeyEmbeddingsModel)

	if anthropicKey == "" {
		return nil, fmt.Errorf("no Anthropic API key configured for workspace %s", workspaceID)
	}

	return ai_client.NewAIProviderWithKeys(ai_client.AIProviderParams{
		AnthropicKey:    anthropicKey,
		OpenAIKey:       openAIKey,
		IndexingModel:   indexingModel,
		NewsletterModel: newsletterModel,
		EmbeddingsModel: embeddingsModel,
	}), nil
}
