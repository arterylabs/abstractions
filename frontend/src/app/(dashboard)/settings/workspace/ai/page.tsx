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

const PROVIDER_ICON: Record<string, string> = {
  anthropic: "/assets/img/claude.svg",
  openai: "/assets/img/openai.svg",
};

export default function WorkspaceAISettingsPage() {
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

  const handleRemoveKey = async (
    keyName: "anthropicApiKey" | "openAIApiKey",
  ) => {
    if (!currentWorkspace?.id) return;
    try {
      await api.workspaces.updateWorkspaceAISettings(currentWorkspace.id, {
        [keyName]: "",
      });
      await fetchSettings(currentWorkspace.id);
      Notifier.success("API key removed.");
    } catch {
      Notifier.error("Failed to remove API key.");
    }
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
  const allModels = modelsResponse.result ?? [];
  const textGenerationModels = allModels.filter(
    (m) => m.supportsTextGeneration,
  );
  const embeddingModels = allModels.filter((m) => m.supportsEmbeddings);
  const activeProviders = getActiveProviders();

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 gap-8 items-start">
        {/* Left: Models */}
        <div className="space-y-6 border-b pb-8 mb-4">
          <div className="space-y-1">
            <h2 className="text-sm font-semibold">Models</h2>
            <p className="text-sm text-muted-foreground">
              Choose which models to use for different tasks in your workspace.
              Faster, cheaper models work well for indexing — stronger models
              give better newsletter output.
            </p>
          </div>

          {/* Indexing Model */}
          <div className="space-y-1.5">
            <Label htmlFor="indexing-model">Indexing Model</Label>
            <p className="text-sm text-muted-foreground">
              Used when indexing your codebase to generate summaries of files,
              functions, and symbols. A faster, cheaper model is recommended
              here.
            </p>
            <Select value={indexingModel} onValueChange={setIndexingModel}>
              <SelectTrigger id="indexing-model" className="w-full">
                <SelectValue placeholder="Select model" />
              </SelectTrigger>
              <SelectContent>
                {textGenerationModels.map((m) => (
                  <SelectItem key={m.id} value={m.id}>
                    <div className="flex items-center gap-2">
                      <img
                        src={PROVIDER_ICON[m.provider]}
                        alt=""
                        className="w-4 h-4"
                      />
                      {m.displayName}
                      <span className="text-muted-foreground">
                        (${m.pricePerMillionTokens}/M tokens)
                      </span>
                    </div>
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          {/* Newsletter Model */}
          <div className="space-y-1.5">
            <Label htmlFor="newsletter-model">Newsletter Model</Label>
            <p className="text-sm text-muted-foreground">
              Used to generate newsletter topics, draft content, and compile the
              final newsletter from your indexed codebase activity.
            </p>
            <Select value={newsletterModel} onValueChange={setNewsletterModel}>
              <SelectTrigger id="newsletter-model" className="w-full">
                <SelectValue placeholder="Select model" />
              </SelectTrigger>
              <SelectContent>
                {textGenerationModels.map((m) => (
                  <SelectItem key={m.id} value={m.id}>
                    <div className="flex items-center gap-2">
                      <img
                        src={PROVIDER_ICON[m.provider]}
                        alt=""
                        className="w-4 h-4"
                      />
                      {m.displayName}
                      <span className="text-muted-foreground">
                        (${m.pricePerMillionTokens}/M tokens)
                      </span>
                    </div>
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          {/* Embeddings Model */}
          <div className="space-y-1.5">
            <Label htmlFor="embeddings-model">Embeddings Model</Label>
            <p className="text-sm text-muted-foreground">
              Used to convert code and text into vector embeddings that power
              semantic search across your indexed codebase.
            </p>
            <Select value={embeddingsModel} onValueChange={setEmbeddingsModel}>
              <SelectTrigger id="embeddings-model" className="w-full">
                <SelectValue placeholder="Select model" />
              </SelectTrigger>
              <SelectContent>
                {embeddingModels.map((m) => (
                  <SelectItem key={m.id} value={m.id}>
                    <div className="flex items-center gap-2">
                      <img
                        src={PROVIDER_ICON[m.provider]}
                        alt=""
                        className="w-4 h-4"
                      />
                      {m.displayName}
                      <span className="text-muted-foreground">
                        (${m.pricePerMillionTokens}/M tokens)
                      </span>
                    </div>
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>

        {/* Right: Model Providers */}
        <div className="space-y-6">
          <div className="space-y-1">
            <h2 className="text-sm font-semibold">Model Providers</h2>
            <p className="text-sm text-muted-foreground">
              Configure API keys for the AI providers used by your workspace.
            </p>
          </div>

          {activeProviders.size === 0 && (
            <p className="text-sm text-muted-foreground">
              Select models above to configure API keys.
            </p>
          )}

          {/* Anthropic API Key */}
          {activeProviders.has("anthropic") && (
            <div className="space-y-3">
              <div>
                <div className="flex items-center gap-2">
                  <img
                    src={PROVIDER_ICON.anthropic}
                    alt="Anthropic"
                    className="w-4 h-4"
                  />
                  <h3 className="text-sm font-medium">Anthropic API Key</h3>
                </div>
                <p className="text-sm text-muted-foreground mt-0.5">
                  Required for the Anthropic models selected above.
                </p>
              </div>
              {settings?.anthropicApiKeyLastFour && (
                <div className="flex items-center gap-3">
                  <span className="font-mono text-sm bg-muted px-2 py-1 rounded">
                    ••••••••••••{settings.anthropicApiKeyLastFour}
                  </span>
                  <Button
                    variant="link"
                    size="sm"
                    className="text-xs text-destructive p-0 h-auto"
                    onClick={() => handleRemoveKey("anthropicApiKey")}
                  >
                    Remove
                  </Button>
                </div>
              )}
              <Input
                id="anthropic-key"
                type="password"
                toggleVisibility
                label={
                  settings?.anthropicApiKeyLastFour ? "Replace key" : "Add key"
                }
                placeholder="sk-ant-…"
                value={anthropicKey}
                onChange={(value) => setAnthropicKey(value)}
              />
            </div>
          )}

          {/* OpenAI API Key */}
          {activeProviders.has("openai") && (
            <div className="space-y-3">
              <div>
                <div className="flex items-center gap-2">
                  <img
                    src={PROVIDER_ICON.openai}
                    alt="OpenAI"
                    className="w-4 h-4"
                  />
                  <h3 className="text-sm font-medium">OpenAI API Key</h3>
                </div>
                <p className="text-sm text-muted-foreground mt-0.5">
                  Required for the OpenAI models selected above.
                </p>
              </div>
              {settings?.openAIApiKeyLastFour && (
                <div className="flex items-center gap-3">
                  <span className="font-mono text-sm bg-muted px-2 py-1 rounded">
                    ••••••••••••{settings.openAIApiKeyLastFour}
                  </span>
                  <Button
                    variant="link"
                    size="sm"
                    className="text-xs text-destructive p-0 h-auto"
                    onClick={() => handleRemoveKey("openAIApiKey")}
                  >
                    Remove
                  </Button>
                </div>
              )}
              <Input
                id="openai-key"
                type="password"
                toggleVisibility
                label={
                  settings?.openAIApiKeyLastFour ? "Replace key" : "Add key"
                }
                placeholder="sk-…"
                value={openAIKey}
                onChange={(value) => setOpenAIKey(value)}
              />
            </div>
          )}
        </div>
      </div>

      <Button
        onClick={handleSave}
        disabled={isSaving}
        size="lg"
        className="w-full"
      >
        {isSaving ? <Spinner /> : null}
        Save
      </Button>
    </div>
  );
}
