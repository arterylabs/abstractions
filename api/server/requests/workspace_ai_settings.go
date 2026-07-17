package requests

import (
	ai_client "abstractions/clients/ai"
	"fmt"
	"unicode/utf8"
)

// SystemPromptsUpdate holds optional system prompt enhancements. A nil field is ignored; an empty
// string clears the stored enhancement (built-in default only). Non-empty values are appended on
// top of the built-in system prompt — they never replace it.
type SystemPromptsUpdate struct {
	SummarizeCode   *string `json:"summarizeCode"`
	SummarizeFile   *string `json:"summarizeFile"`
	TopicGenerator  *string `json:"topicGenerator"`
	NewsletterEntry *string `json:"newsletterEntry"`
}

// UpdateWorkspaceAISettings updates a workspace's AI configuration. APIKeys is keyed by provider ID
// (e.g. "anthropic", "openai", "grok"); a nil entry is ignored, an empty-string entry clears the
// stored key. Model fields and system prompt fields are ignored when nil.
type UpdateWorkspaceAISettings struct {
	APIKeys         map[string]*string   `json:"apiKeys"`
	IndexingModel   *string              `json:"indexingModel"`
	NewsletterModel *string              `json:"newsletterModel"`
	EmbeddingsModel *string              `json:"embeddingsModel"`
	SystemPrompts   *SystemPromptsUpdate `json:"systemPrompts"`
}

func (r *UpdateWorkspaceAISettings) Validate() error {
	if r.SystemPrompts == nil {
		return nil
	}
	for name, value := range map[string]*string{
		"summarizeCode":   r.SystemPrompts.SummarizeCode,
		"summarizeFile":   r.SystemPrompts.SummarizeFile,
		"topicGenerator":  r.SystemPrompts.TopicGenerator,
		"newsletterEntry": r.SystemPrompts.NewsletterEntry,
	} {
		if value == nil {
			continue
		}
		if utf8.RuneCountInString(*value) > ai_client.MaxSystemPromptLength {
			return fmt.Errorf("systemPrompts.%s exceeds maximum length of %d characters", name, ai_client.MaxSystemPromptLength)
		}
	}
	return nil
}
