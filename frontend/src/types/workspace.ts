export type Workspace = {
  id: string;
  createdAt: string;
  updatedAt: string;
  name: string;
  logo?: string | null;
  website?: string | null;
  about?: string | null;
  billingEmail?: string;
  userRole?: string;
};

export type AIModel = {
  id: string;
  displayName: string;
  // Provider ID (e.g. "anthropic", "openai", "grok"). Not a fixed union — providers are
  // registry-driven on the backend and fetched via getAIProviders.
  provider: string;
  supportsTextGeneration: boolean;
  supportsEmbeddings: boolean;
  pricePerMillionTokens: number;
};

// AIProvider describes a supported model provider. Everything provider-specific in the UI (icon,
// key placeholder, key rows) is driven off this list rather than hardcoded.
export type AIProvider = {
  id: string;
  displayName: string;
  iconPath: string;
  keyPlaceholder: string;
  supportsEmbeddings: boolean;
};

// SystemPrompts holds optional workspace enhancements layered on top of built-in system prompts.
// Null/missing means no enhancement (built-in default only). Values never replace the default.
export type SystemPrompts = {
  summarizeCode?: string | null;
  summarizeFile?: string | null;
  topicGenerator?: string | null;
  newsletterEntry?: string | null;
};

// DefaultSystemPrompts holds the built-in system prompts (read-only reference from the API).
export type DefaultSystemPrompts = {
  summarizeCode: string;
  summarizeFile: string;
  topicGenerator: string;
  newsletterEntry: string;
};

export type WorkspaceAISettings = {
  // Last four characters of each provider's configured key, keyed by provider ID; null when unset.
  apiKeyLastFour?: Record<string, string | null>;
  indexingModel?: string | null;
  newsletterModel?: string | null;
  embeddingsModel?: string | null;
  // Stored enhancements only; null/missing means the workspace uses the built-in default alone.
  systemPrompts?: SystemPrompts;
  // Built-in defaults (always applied; enhancements are appended on top).
  defaultSystemPrompts?: DefaultSystemPrompts;
};
