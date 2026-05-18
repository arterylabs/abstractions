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
  provider: "anthropic" | "openai";
  supportsTextGeneration: boolean;
  supportsEmbeddings: boolean;
  pricePerMillionTokens: number;
};

export type WorkspaceAISettings = {
  anthropicApiKeyLastFour?: string | null;
  openAIApiKeyLastFour?: string | null;
  indexingModel?: string | null;
  newsletterModel?: string | null;
  embeddingsModel?: string | null;
};
