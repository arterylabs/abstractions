package jobs

import (
	ai_client "abstractions/clients/ai"
	"abstractions/config"
	"abstractions/logger"
	"abstractions/repositories"
	workspace_configuration_service "abstractions/services/workspace_configuration"
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// newWorkspaceAIProvider creates an AIProvider resolved against the workspace configuration for
// the given workspace. It resolves the configured models (or their defaults) and validates that
// the API keys required by those models are available — either from the workspace configuration
// or from the system environment. Returns an error only when a key that is actually needed by
// the configured models is missing.
func newWorkspaceAIProvider(ctx context.Context, db *gorm.DB, workspaceID uuid.UUID) (ai_client.AIProvider, error) {
	encryptionKey := ""
	if config.Environment != nil {
		encryptionKey = config.Environment.EncryptionKey32Bit
	}
	cfgRepo := repositories.NewWorkspaceConfigurationRepository(db)
	cfgSvc := workspace_configuration_service.NewService(cfgRepo, encryptionKey)

	indexingModel, _ := cfgSvc.Get(ctx, workspaceID, workspace_configuration_service.KeyIndexingModel)
	newsletterModel, _ := cfgSvc.Get(ctx, workspaceID, workspace_configuration_service.KeyNewsletterModel)
	embeddingsModel, _ := cfgSvc.Get(ctx, workspaceID, workspace_configuration_service.KeyEmbeddingsModel)

	summarizeCodePrompt, _ := cfgSvc.Get(ctx, workspaceID, workspace_configuration_service.KeySystemPromptSummarizeCode)
	summarizeFilePrompt, _ := cfgSvc.Get(ctx, workspaceID, workspace_configuration_service.KeySystemPromptSummarizeFile)
	topicGeneratorPrompt, _ := cfgSvc.Get(ctx, workspaceID, workspace_configuration_service.KeySystemPromptTopicGenerator)
	newsletterEntryPrompt, _ := cfgSvc.Get(ctx, workspaceID, workspace_configuration_service.KeySystemPromptNewsletterEntry)

	// Resolve effective model IDs (apply defaults for empty values) so we can determine
	// which provider keys are actually required before validating.
	effectiveIndexing := indexingModel
	if effectiveIndexing == "" {
		effectiveIndexing = ai_client.DefaultIndexingModel
	}
	effectiveNewsletter := newsletterModel
	if effectiveNewsletter == "" {
		effectiveNewsletter = ai_client.DefaultNewsletterModel
	}
	effectiveEmbeddings := embeddingsModel
	if effectiveEmbeddings == "" {
		effectiveEmbeddings = ai_client.DefaultEmbeddingsModel
	}

	logger.Debug(ctx, fmt.Sprintf(
		"workspace_ai: resolved models workspaceId=%s indexingModel=%s newsletterModel=%s embeddingsModel=%s",
		workspaceID, effectiveIndexing, effectiveNewsletter, effectiveEmbeddings,
	))

	// Build the set of providers required by the configured models.
	// We track the first model that requires each provider for a clear error message.
	requiredProviders := map[ai_client.Provider]string{}
	for _, modelID := range []string{effectiveIndexing, effectiveNewsletter, effectiveEmbeddings} {
		if info, ok := ai_client.ModelByID(modelID); ok {
			if _, alreadyRequired := requiredProviders[info.Provider]; !alreadyRequired {
				requiredProviders[info.Provider] = modelID
			}
		}
	}

	// Resolve each provider's key from the workspace configuration (which already applies the
	// system-level env fallback via the registry), then validate that every required provider has
	// one. Build the Keys map for the factory in the same pass.
	keys := make(map[ai_client.Provider]string, len(ai_client.Providers))
	for _, descriptor := range ai_client.Providers {
		key, _ := cfgSvc.Get(ctx, workspaceID, descriptor.ConfigKey)
		keys[descriptor.ID] = key

		logger.Debug(ctx, fmt.Sprintf(
			"workspace_ai: key availability workspaceId=%s provider=%s hasKey=%v",
			workspaceID, descriptor.ID, key != "",
		))

		if model, needed := requiredProviders[descriptor.ID]; needed && key == "" {
			return nil, fmt.Errorf("no %s API key configured for workspace %s (required by model %q)", descriptor.DisplayName, workspaceID, model)
		}
	}

	return ai_client.NewAIProviderWithKeys(ai_client.AIProviderParams{
		Keys:            keys,
		IndexingModel:   indexingModel,
		NewsletterModel: newsletterModel,
		EmbeddingsModel: embeddingsModel,
		SystemPrompts: ai_client.SystemPromptOverrides{
			SummarizeCode:   summarizeCodePrompt,
			SummarizeFile:   summarizeFilePrompt,
			TopicGenerator:  topicGeneratorPrompt,
			NewsletterEntry: newsletterEntryPrompt,
		},
	}), nil
}
