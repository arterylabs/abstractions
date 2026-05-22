"use client";

import { isAxiosError } from "axios";
import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Notifier } from "@/components/ui/sonner";
import Spinner from "@/components/ui/spinner";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useAppSelector } from "@/store/hooks";
import { useAPI } from "@/api/hooks";
import api from "@/api";
import { AIModel, WorkspaceAISettings } from "@/types/workspace";
import { Label } from "@/components/ui/label";
import { PageHeader } from "@/components/ui/page-header";
import { askConfirmation } from "@/components/ui/confirmation";

const PROVIDER_ICON: Record<string, string> = {
  anthropic: "/assets/img/claude.svg",
  openai: "/assets/img/openai.svg",
};

function ModelRow({
  id,
  label,
  description,
  value,
  onChange,
  models,
}: {
  id: string;
  label: string;
  description: string;
  value: string;
  onChange: (v: string) => void;
  models: AIModel[];
}) {
  return (
    <div className="flex items-start justify-between gap-8 py-4 border-b last:border-b-0">
      <div className="space-y-0.5 flex-1 min-w-0">
        <Label htmlFor={id} className="text-sm font-medium">
          {label}
        </Label>
        <p className="text-xs text-muted-foreground">{description}</p>
      </div>
      <div className="w-72 shrink-0">
        <Select value={value} onValueChange={onChange}>
          <SelectTrigger id={id} className="w-full">
            <SelectValue placeholder="Select model" />
          </SelectTrigger>
          <SelectContent>
            {models.map((m) => (
              <SelectItem key={m.id} value={m.id}>
                <div className="flex items-center gap-2">
                  <img
                    src={PROVIDER_ICON[m.provider]}
                    alt=""
                    className="w-4 h-4"
                  />
                  {m.displayName}
                  <span className="text-muted-foreground text-xs">
                    ${m.pricePerMillionTokens}/M
                  </span>
                </div>
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
    </div>
  );
}

function ProviderKeyRow({
  provider,
  label,
  placeholder,
  lastFour,
  value,
  onChange,
  onRemove,
}: {
  provider: string;
  label: string;
  placeholder: string;
  lastFour?: string;
  value: string;
  onChange: (v: string) => void;
  onRemove: () => void;
}) {
  return (
    <div className="flex items-start justify-between gap-8 py-4 border-b last:border-b-0">
      <div className="space-y-1 flex-1 min-w-0">
        <div className="flex items-center gap-2">
          <img src={PROVIDER_ICON[provider]} alt={label} className="w-4 h-4" />
          <span className="text-sm font-medium">{label}</span>
        </div>
        {lastFour ? (
          <div className="flex items-center gap-2 mt-1">
            <span className="font-mono text-xs bg-muted px-2 py-1 rounded text-muted-foreground">
              ••••••••••••{lastFour}
            </span>
            <Button
              variant="link"
              size="sm"
              className="text-xs text-destructive p-0 h-auto"
              onClick={onRemove}
            >
              Remove
            </Button>
          </div>
        ) : (
          <p className="text-xs text-muted-foreground">No key configured.</p>
        )}
      </div>
      <div className="w-72 shrink-0">
        <Input
          id={`${provider}-key`}
          type="password"
          toggleVisibility
          label={lastFour ? "Replace key" : "Add key"}
          placeholder={placeholder}
          value={value}
          onChange={onChange}
        />
      </div>
    </div>
  );
}

export default function AISettingsPage() {
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );

  const [anthropicKey, setAnthropicKey] = useState("");
  const [openAIKey, setOpenAIKey] = useState("");
  const [indexingModel, setIndexingModel] = useState("");
  const [newsletterModel, setNewsletterModel] = useState("");
  const [embeddingsModel, setEmbeddingsModel] = useState("");
  const [isSaving, setIsSaving] = useState(false);

  const [modelsResponse, fetchModels] = useAPI<AIModel[]>(
    api.workspaces.getAIModels,
  );
  const [settingsResponse, fetchSettings] = useAPI<WorkspaceAISettings>(
    api.workspaces.getWorkspaceAISettings,
  );

  const allModels = modelsResponse.result ?? [];

  useEffect(() => {
    fetchModels();
  }, []);

  useEffect(() => {
    if (currentWorkspace?.id) {
      fetchSettings(currentWorkspace.id);
    }
  }, [currentWorkspace?.id]);

  useEffect(() => {
    if (settingsResponse.result) {
      const s = settingsResponse.result;
      setIndexingModel(s.indexingModel ?? "claude-haiku-4-5-20251001");
      setNewsletterModel(s.newsletterModel ?? "claude-sonnet-4-6");
      setEmbeddingsModel(s.embeddingsModel ?? "text-embedding-3-small");
    }
  }, [settingsResponse.result]);

  const getActiveProviders = () => {
    const models = modelsResponse.result ?? [];
    return new Set<string>(
      [indexingModel, newsletterModel, embeddingsModel]
        .map((modelId) => models.find((m) => m.id === modelId)?.provider)
        .filter(Boolean) as string[],
    );
  };

  const handleSave = async () => {
    if (!currentWorkspace?.id) return;

    const savedSettings = settingsResponse.result;
    const requiredProviders = getActiveProviders();
    const missingKeys: string[] = [];
    if (
      requiredProviders.has("anthropic") &&
      !savedSettings?.anthropicApiKeyLastFour &&
      !anthropicKey
    ) {
      missingKeys.push("Anthropic");
    }
    if (
      requiredProviders.has("openai") &&
      !savedSettings?.openAIApiKeyLastFour &&
      !openAIKey
    ) {
      missingKeys.push("OpenAI");
    }
    if (missingKeys.length > 0) {
      Notifier.error(
        `API key required for: ${missingKeys.join(", ")}. Add a key or switch to a different model.`,
      );
      return;
    }

    setIsSaving(true);
    try {
      await api.workspaces.updateWorkspaceAISettings(currentWorkspace.id, {
        anthropicApiKey: anthropicKey || undefined,
        openAIApiKey: openAIKey || undefined,
        indexingModel,
        newsletterModel,
        embeddingsModel,
      });
      setAnthropicKey("");
      setOpenAIKey("");
      await fetchSettings(currentWorkspace.id);
      Notifier.success("AI settings saved.");
    } catch (err: unknown) {
      const message = isAxiosError(err)
        ? (err.response?.data?.error ?? "Failed to save AI settings.")
        : "Failed to save AI settings.";
      Notifier.error(message);
    } finally {
      setIsSaving(false);
    }
  };

  const handleRemoveKey = (keyName: "anthropicApiKey" | "openAIApiKey") => {
    if (!currentWorkspace?.id) return;

    const providerName = keyName === "anthropicApiKey" ? "anthropic" : "openai";
    const providerLabel =
      keyName === "anthropicApiKey" ? "Anthropic" : "OpenAI";

    const featureLabels: Record<string, string> = {
      indexingModel: "file indexing",
      newsletterModel: "newsletter generation",
      embeddingsModel: "semantic search",
    };
    const modelSelections: Record<string, string> = {
      indexingModel,
      newsletterModel,
      embeddingsModel,
    };

    const affectedFeatures = Object.entries(modelSelections)
      .filter(([, modelId]) => {
        const model = allModels.find((m) => m.id === modelId);
        return model?.provider === providerName;
      })
      .map(([key]) => featureLabels[key]);

    const featureList =
      affectedFeatures.length > 0
        ? ` Features that depend on ${providerLabel} models — ${affectedFeatures.join(", ")} — will stop working until a new key is provided.`
        : "";

    askConfirmation({
      title: `Remove ${providerLabel} API key?`,
      content: `Are you sure you want to remove the ${providerLabel} API key?${featureList}`,
      variant: "destructive",
      okText: "Remove key",
      okCallback: async () => {
        try {
          await api.workspaces.updateWorkspaceAISettings(currentWorkspace.id!, {
            [keyName]: "",
          });
          await fetchSettings(currentWorkspace.id!);
          Notifier.success("API key removed.");
        } catch {
          Notifier.error("Failed to remove API key.");
        }
      },
    });
  };

  if (settingsResponse.isFetching || modelsResponse.isFetching) {
    return (
      <div className="flex items-center gap-2 text-sm text-muted-foreground">
        <Spinner />
        Loading AI settings…
      </div>
    );
  }

  const settings = settingsResponse.result;
  const textGenerationModels = allModels.filter(
    (m) => m.supportsTextGeneration,
  );
  const embeddingModels = allModels.filter((m) => m.supportsEmbeddings);
  const activeProviders = getActiveProviders();

  return (
    <div className="space-y-6">
      <PageHeader
        title="AI"
        description="Configure AI models and API keys for your workspace."
      />

      <div className="space-y-6">
        {/* Models */}
        <div className="rounded-lg border">
          <div className="px-5 py-4 border-b">
            <h2 className="text-sm font-semibold">Models</h2>
            <p className="text-xs text-muted-foreground mt-0.5">
              Choose models for each task. Faster models work well for indexing
              — stronger models give better newsletter output.
            </p>
          </div>
          <div className="px-5">
            <ModelRow
              id="indexing-model"
              label="Indexing Model"
              description="Generates summaries of files, functions, and symbols. A faster, cheaper model is recommended."
              value={indexingModel}
              onChange={setIndexingModel}
              models={textGenerationModels}
            />
            <ModelRow
              id="newsletter-model"
              label="Newsletter Model"
              description="Drafts newsletter topics and compiles content from indexed codebase activity."
              value={newsletterModel}
              onChange={setNewsletterModel}
              models={textGenerationModels}
            />
            <ModelRow
              id="embeddings-model"
              label="Embeddings Model"
              description="Converts code and text into vectors for semantic search."
              value={embeddingsModel}
              onChange={setEmbeddingsModel}
              models={embeddingModels}
            />
          </div>
        </div>

        {/* API Keys */}
        <div className="rounded-lg border">
          <div className="px-5 py-4 border-b">
            <h2 className="text-sm font-semibold">API Keys</h2>
            <p className="text-xs text-muted-foreground mt-0.5">
              Keys are required for each provider used by your selected models.
            </p>
          </div>
          <div className="px-5">
            {activeProviders.size === 0 && (
              <p className="text-sm text-muted-foreground py-4">
                Select models above to configure API keys.
              </p>
            )}
            {activeProviders.has("anthropic") && (
              <ProviderKeyRow
                provider="anthropic"
                label="Anthropic API Key"
                placeholder="sk-ant-…"
                lastFour={settings?.anthropicApiKeyLastFour ?? undefined}
                value={anthropicKey}
                onChange={setAnthropicKey}
                onRemove={() => handleRemoveKey("anthropicApiKey")}
              />
            )}
            {activeProviders.has("openai") && (
              <ProviderKeyRow
                provider="openai"
                label="OpenAI API Key"
                placeholder="sk-…"
                lastFour={settings?.openAIApiKeyLastFour ?? undefined}
                value={openAIKey}
                onChange={setOpenAIKey}
                onRemove={() => handleRemoveKey("openAIApiKey")}
              />
            )}
          </div>
        </div>
      </div>

      <div className="flex justify-end">
        <Button onClick={handleSave} disabled={isSaving}>
          {isSaving ? <Spinner /> : null}
          Save changes
        </Button>
      </div>
    </div>
  );
}
