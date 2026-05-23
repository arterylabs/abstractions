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

	anthropicKey, _ := cfgSvc.Get(ctx, workspaceID, workspace_configuration_service.KeyAnthropicAPIKey)
	openAIKey, _ := cfgSvc.Get(ctx, workspaceID, workspace_configuration_service.KeyOpenAIAPIKey)
	indexingModel, _ := cfgSvc.Get(ctx, workspaceID, workspace_configuration_service.KeyIndexingModel)
	newsletterModel, _ := cfgSvc.Get(ctx, workspaceID, workspace_configuration_service.KeyNewsletterModel)
	embeddingsModel, _ := cfgSvc.Get(ctx, workspaceID, workspace_configuration_service.KeyEmbeddingsModel)

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

	// Apply system-level fallback keys (same logic as NewAIProviderWithKeys) so that
	// validation accounts for env-level defaults before rejecting the configuration.
	effectiveAnthropicKey := anthropicKey
	effectiveOpenAIKey := openAIKey
	if config.Environment != nil && !config.Environment.DisableDefaultAIUsage {
		if effectiveAnthropicKey == "" {
			effectiveAnthropicKey = config.Environment.AnthropicAPIKey
		}
		if effectiveOpenAIKey == "" {
			effectiveOpenAIKey = config.Environment.OpenAIAPIKey
		}
	}

	logger.Debug(ctx, fmt.Sprintf(
		"workspace_ai: key availability workspaceId=%s hasAnthropicKey=%v hasOpenAIKey=%v",
		workspaceID, effectiveAnthropicKey != "", effectiveOpenAIKey != "",
	))

	// Validate that every required provider has an available key.
	if model, needed := requiredProviders[ai_client.ProviderAnthropic]; needed && effectiveAnthropicKey == "" {
		return nil, fmt.Errorf("no Anthropic API key configured for workspace %s (required by model %q)", workspaceID, model)
	}
	if model, needed := requiredProviders[ai_client.ProviderOpenAI]; needed && effectiveOpenAIKey == "" {
		return nil, fmt.Errorf("no OpenAI API key configured for workspace %s (required by model %q)", workspaceID, model)
	}

	return ai_client.NewAIProviderWithKeys(ai_client.AIProviderParams{
		AnthropicKey:    anthropicKey,
		OpenAIKey:       openAIKey,
		IndexingModel:   indexingModel,
		NewsletterModel: newsletterModel,
		EmbeddingsModel: embeddingsModel,
	}), nil
}
