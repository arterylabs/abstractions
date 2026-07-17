"use client";

import { isAxiosError } from "axios";
import { useEffect, useMemo, useState } from "react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Notifier } from "@/components/ui/sonner";
import Spinner from "@/components/ui/spinner";
import Icon from "@/components/ui/icon";
import Tooltip from "@/components/ui/tooltip";
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
import {
  WorkspaceHealthResponse,
  WorkspaceProviderHealth,
} from "@/api/workspaces";
import {
  AIModel,
  AIProvider,
  SystemPrompts,
  WorkspaceAISettings,
} from "@/types/workspace";
import { Label } from "@/components/ui/label";
import { PageHeader } from "@/components/ui/page-header";
import { FloatingBar } from "@/components/ui/floating-bar";
import { askConfirmation } from "@/components/ui/confirmation";
import { cn } from "@/lib/utils";

type SystemPromptKey = keyof Required<SystemPrompts>;

const SYSTEM_PROMPT_FIELDS: {
  key: SystemPromptKey;
  label: string;
  description: string;
  placeholder: string;
}[] = [
  {
    key: "summarizeCode",
    label: "Code chunk summary",
    description:
      "Optional guidance layered on top of the built-in indexing instructions for each code chunk.",
    placeholder:
      "e.g. Prefer security and auth concerns. Call out concurrency risks when present.",
  },
  {
    key: "summarizeFile",
    label: "File summary",
    description:
      "Optional guidance layered on top of the built-in whole-file summary instructions.",
    placeholder:
      "e.g. Emphasize how this file connects to neighbouring modules.",
  },
  {
    key: "topicGenerator",
    label: "Topic generation",
    description:
      "Optional guidance layered on top of the built-in topic selection criteria.",
    placeholder:
      "e.g. Favour reliability and operational topics over pure architecture deep-dives.",
  },
  {
    key: "newsletterEntry",
    label: "Newsletter writing",
    description:
      "Optional guidance layered on top of the built-in newsletter voice and style.",
    placeholder:
      "e.g. Keep a dry, slightly irreverent tone. Avoid sports metaphors.",
  },
];

const emptyPromptDrafts = (): Record<SystemPromptKey, string> => ({
  summarizeCode: "",
  summarizeFile: "",
  topicGenerator: "",
  newsletterEntry: "",
});

function ModelRow({
  id,
  label,
  description,
  value,
  onChange,
  models,
  providerIcons,
  warning,
}: {
  id: string;
  label: string;
  description: string;
  value: string;
  onChange: (v: string) => void;
  models: AIModel[];
  providerIcons: Record<string, string>;
  warning?: string;
}) {
  return (
    <div className="flex items-start justify-between gap-8 py-4 border-b last:border-b-0">
      <div className="space-y-0.5 flex-1 min-w-0">
        <Label htmlFor={id} className="text-sm font-medium">
          {label}
        </Label>
        <p className="text-xs text-muted-foreground">{description}</p>
      </div>
      <div className="w-72 shrink-0 space-y-1.5">
        <Select value={value} onValueChange={onChange}>
          <SelectTrigger id={id} className="w-full">
            <SelectValue placeholder="Select model" />
          </SelectTrigger>
          <SelectContent>
            {models.map((m) => (
              <SelectItem key={m.id} value={m.id}>
                <div className="flex items-center gap-2">
                  {providerIcons[m.provider] ? (
                    <img
                      src={providerIcons[m.provider]}
                      alt=""
                      className="w-4 h-4"
                    />
                  ) : null}
                  {m.displayName}
                  <span className="text-muted-foreground text-xs">
                    ${m.pricePerMillionTokens}/M
                  </span>
                </div>
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        {warning ? (
          <p className="flex items-center gap-1 text-xs text-warning-dark">
            <Icon type="warning" className="text-sm" />
            {warning}
          </p>
        ) : null}
      </div>
    </div>
  );
}

// ProviderStatusIndicator is the small badge overlaid on the provider logo: a green check or a
// red error icon (with the failure reason in a tooltip) once the health check has run.
function ProviderStatusIndicator({
  health,
}: {
  health?: WorkspaceProviderHealth;
}) {
  if (!health?.configured) return null;
  if (health.valid) {
    return (
      <span className="absolute -bottom-1.5 -right-1.5 flex items-center justify-center rounded-full bg-background">
        <Icon
          type="check_circle"
          className="text-xl leading-none text-success-dark"
        />
      </span>
    );
  }
  return (
    <Tooltip value={health.error || "API key is invalid."}>
      <span className="absolute -bottom-1.5 -right-1.5 flex items-center justify-center rounded-full bg-background">
        <Icon
          type="error"
          className="text-xl leading-none text-destructive-dark"
        />
      </span>
    </Tooltip>
  );
}

function ProviderCard({
  provider,
  lastFour,
  health,
  checking,
  required,
  onSaveKey,
  onRemove,
}: {
  provider: AIProvider;
  lastFour?: string;
  health?: WorkspaceProviderHealth;
  checking: boolean;
  required: boolean;
  // onSaveKey persists the key and refreshes settings + health; throws on failure.
  onSaveKey: (key: string) => Promise<void>;
  onRemove: () => void;
}) {
  const [formOpen, setFormOpen] = useState(false);
  const [draft, setDraft] = useState("");
  const [isSaving, setIsSaving] = useState(false);

  const configured = health?.configured ?? Boolean(lastFour);

  const handleSave = async () => {
    setIsSaving(true);
    try {
      await onSaveKey(draft);
      setDraft("");
      setFormOpen(false);
    } catch (err: unknown) {
      const fallback = `Failed to save ${provider.displayName} API key.`;
      const message = isAxiosError(err)
        ? (err.response?.data?.error ?? fallback)
        : fallback;
      Notifier.error(message);
    } finally {
      setIsSaving(false);
    }
  };

  // While the health check runs, a spinner takes the badge's place.
  let statusBadge = checking ? (
    <span className="flex h-6 items-center">
      <Spinner className="w-4 h-4" />
    </span>
  ) : (
    <Badge variant="outline" size="sm" className="text-muted-foreground">
      Not connected
    </Badge>
  );
  if (!checking && configured) {
    // Without health data (check failed) we know a key exists but not whether it works.
    statusBadge = !health ? (
      <Badge variant="secondary" size="sm">
        Connected
      </Badge>
    ) : health.valid ? (
      <Badge variant="success" size="sm">
        Connected
      </Badge>
    ) : (
      <Badge variant="destructive" size="sm">
        Invalid key
      </Badge>
    );
  }

  return (
    <Card
      className={cn("flex flex-col gap-4 p-5", !configured && "bg-muted/30")}
    >
      <div className="flex items-start justify-between gap-3">
        <div className="relative shrink-0">
          <div
            className={cn(
              "flex size-12 items-center justify-center rounded-lg border bg-background",
              !configured && "opacity-70",
            )}
          >
            <img
              src={provider.iconPath}
              alt={provider.displayName}
              className={cn("w-7 h-7", !configured && "opacity-60 grayscale")}
            />
          </div>
          {!checking && <ProviderStatusIndicator health={health} />}
        </div>
        {statusBadge}
      </div>

      <div className="space-y-1 flex-1">
        <div className="text-sm font-semibold">{provider.displayName}</div>
        {configured && lastFour ? (
          <span className="inline-block font-mono text-xs bg-muted px-2 py-1 rounded text-muted-foreground">
            ••••••••••••{lastFour}
          </span>
        ) : (
          <p className="text-xs text-muted-foreground">
            No API key configured.
          </p>
        )}
        {required && !configured ? (
          <p className="text-xs text-warning-dark">
            Required by your selected models.
          </p>
        ) : null}
      </div>

      {formOpen ? (
        <div className="space-y-2">
          <Input
            id={`${provider.id}-key`}
            type="password"
            toggleVisibility
            label={configured ? "Update key" : "Add key"}
            placeholder={provider.keyPlaceholder}
            value={draft}
            onChange={setDraft}
            disabled={isSaving}
          />
          <div className="flex items-center gap-2">
            <Button
              size="sm"
              onClick={handleSave}
              disabled={isSaving || !draft}
            >
              {isSaving ? <Spinner className="w-4 h-4" /> : null}
              Save key
            </Button>
            <Button
              size="sm"
              variant="ghost"
              disabled={isSaving}
              onClick={() => {
                setDraft("");
                setFormOpen(false);
              }}
            >
              Cancel
            </Button>
          </div>
        </div>
      ) : (
        <div className="flex items-center gap-2">
          <Button size="sm" variant="outline" onClick={() => setFormOpen(true)}>
            {configured ? "Update key" : "Add key"}
          </Button>
          {configured ? (
            <Button
              size="sm"
              variant="ghost"
              className="text-destructive-dark hover:text-destructive-dark"
              onClick={onRemove}
            >
              Remove
            </Button>
          ) : null}
        </div>
      )}
    </Card>
  );
}

export default function AISettingsPage() {
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );

  const [indexingModel, setIndexingModel] = useState("");
  const [newsletterModel, setNewsletterModel] = useState("");
  const [embeddingsModel, setEmbeddingsModel] = useState("");
  const [isSaving, setIsSaving] = useState(false);
  const [promptDrafts, setPromptDrafts] =
    useState<Record<SystemPromptKey, string>>(emptyPromptDrafts);

  // Baseline values for dirty detection (mirrors repository settings)
  const [savedIndexingModel, setSavedIndexingModel] = useState("");
  const [savedNewsletterModel, setSavedNewsletterModel] = useState("");
  const [savedEmbeddingsModel, setSavedEmbeddingsModel] = useState("");
  const [savedPromptDrafts, setSavedPromptDrafts] =
    useState<Record<SystemPromptKey, string>>(emptyPromptDrafts);

  const [modelsResponse, fetchModels] = useAPI<AIModel[]>(
    api.workspaces.getAIModels,
  );
  const [providersResponse, fetchProviders] = useAPI<AIProvider[]>(
    api.workspaces.getAIProviders,
  );
  const [settingsResponse, fetchSettings] = useAPI<
    WorkspaceAISettings,
    [string]
  >(api.workspaces.getWorkspaceAISettings);
  const [healthResponse, checkHealth] = useAPI<
    WorkspaceHealthResponse,
    [string]
  >(api.workspaces.getWorkspaceHealth);

  const allModels = useMemo(
    () => modelsResponse.result ?? [],
    [modelsResponse.result],
  );
  const allProviders = useMemo(
    () => providersResponse.result ?? [],
    [providersResponse.result],
  );
  const providersById = useMemo(
    () => new Map(allProviders.map((p) => [p.id, p])),
    [allProviders],
  );
  const providerIcons = useMemo(
    () =>
      Object.fromEntries(allProviders.map((p) => [p.id, p.iconPath])) as Record<
        string,
        string
      >,
    [allProviders],
  );

  const settingsLoaded = Boolean(settingsResponse.result);

  useEffect(() => {
    fetchModels();
    fetchProviders();
  }, []);

  useEffect(() => {
    if (currentWorkspace?.id) {
      fetchSettings(currentWorkspace.id);
      checkHealth(currentWorkspace.id);
    }
  }, [currentWorkspace?.id]);

  useEffect(() => {
    if (settingsResponse.result) {
      const s = settingsResponse.result;
      const indexing = s.indexingModel ?? "claude-haiku-4-5-20251001";
      const newsletter = s.newsletterModel ?? "claude-sonnet-4-6";
      const embeddings = s.embeddingsModel ?? "text-embedding-3-small";

      setIndexingModel(indexing);
      setNewsletterModel(newsletter);
      setEmbeddingsModel(embeddings);
      setSavedIndexingModel(indexing);
      setSavedNewsletterModel(newsletter);
      setSavedEmbeddingsModel(embeddings);

      // Drafts hold enhancements only — never the full built-in system prompt.
      const enhancements = s.systemPrompts;
      const prompts = {
        summarizeCode: enhancements?.summarizeCode ?? "",
        summarizeFile: enhancements?.summarizeFile ?? "",
        topicGenerator: enhancements?.topicGenerator ?? "",
        newsletterEntry: enhancements?.newsletterEntry ?? "",
      };
      setPromptDrafts(prompts);
      setSavedPromptDrafts(prompts);
    }
  }, [settingsResponse.result]);

  // Provider IDs used by the currently selected models — these need API keys.
  const activeProviders = useMemo(() => {
    return new Set<string>(
      [indexingModel, newsletterModel, embeddingsModel]
        .map((modelId) => allModels.find((m) => m.id === modelId)?.provider)
        .filter(Boolean) as string[],
    );
  }, [indexingModel, newsletterModel, embeddingsModel, allModels]);

  const lastFourFor = (providerId: string) =>
    settingsResponse.result?.apiKeyLastFour?.[providerId] ?? undefined;

  const healthFor = (providerId: string) =>
    healthResponse.result?.providers?.[providerId];

  const refreshSettingsAndHealth = async () => {
    if (!currentWorkspace?.id) return;
    await Promise.all([
      fetchSettings(currentWorkspace.id),
      checkHealth(currentWorkspace.id),
    ]);
  };

  const handleSaveKey = async (providerId: string, key: string) => {
    if (!currentWorkspace?.id) return;
    await api.workspaces.updateWorkspaceAISettings(currentWorkspace.id, {
      apiKeys: { [providerId]: key },
    });
    Notifier.success(
      `${providersById.get(providerId)?.displayName ?? providerId} API key saved.`,
    );
    await refreshSettingsAndHealth();
  };

  const isPromptCustom = (key: SystemPromptKey) => {
    const stored = savedPromptDrafts[key];
    return Boolean(stored && stored.trim().length > 0);
  };

  const isPromptDirty = (key: SystemPromptKey) => {
    if (!settingsLoaded) return false;
    return promptDrafts[key] !== savedPromptDrafts[key];
  };

  const modelsDirty =
    indexingModel !== savedIndexingModel ||
    newsletterModel !== savedNewsletterModel ||
    embeddingsModel !== savedEmbeddingsModel;

  const promptsDirty = SYSTEM_PROMPT_FIELDS.some((f) => isPromptDirty(f.key));
  const isDirty = modelsDirty || promptsDirty;

  const handleClearPrompt = (key: SystemPromptKey) => {
    setPromptDrafts((prev) => ({ ...prev, [key]: "" }));
  };

  const handleDiscardChanges = () => {
    setIndexingModel(savedIndexingModel);
    setNewsletterModel(savedNewsletterModel);
    setEmbeddingsModel(savedEmbeddingsModel);
    setPromptDrafts({ ...savedPromptDrafts });
  };

  const handleSave = async () => {
    if (!currentWorkspace?.id || !settingsLoaded) return;

    // Every active provider must have a configured key.
    const missing = [...activeProviders].filter((id) => !lastFourFor(id));
    if (missing.length > 0) {
      const names = missing.map(
        (id) => providersById.get(id)?.displayName ?? id,
      );
      Notifier.error(
        `API key required for: ${names.join(", ")}. Connect the provider above or switch to a different model.`,
      );
      return;
    }

    const systemPrompts: SystemPrompts = {};
    for (const { key } of SYSTEM_PROMPT_FIELDS) {
      // Empty / whitespace-only clears the enhancement (built-in default only).
      systemPrompts[key] =
        promptDrafts[key].trim() === "" ? "" : promptDrafts[key];
    }

    setIsSaving(true);
    try {
      await api.workspaces.updateWorkspaceAISettings(currentWorkspace.id, {
        indexingModel,
        newsletterModel,
        embeddingsModel,
        systemPrompts,
      });
      await refreshSettingsAndHealth();
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

  const handleRemoveKey = (providerId: string) => {
    if (!currentWorkspace?.id) return;

    const providerLabel =
      providersById.get(providerId)?.displayName ?? providerId;

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
        return model?.provider === providerId;
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
            apiKeys: { [providerId]: "" },
          });
          await refreshSettingsAndHealth();
          Notifier.success("API key removed.");
        } catch {
          Notifier.error("Failed to remove API key.");
        }
      },
    });
  };

  // Inline warning for a model row when its provider's key is missing or failing.
  const warningForModel = (modelId: string) => {
    const providerId = allModels.find((m) => m.id === modelId)?.provider;
    if (!providerId) return undefined;
    const providerLabel =
      providersById.get(providerId)?.displayName ?? providerId;
    const health = healthFor(providerId);
    if (health) {
      if (!health.configured)
        return `${providerLabel} is not connected — add an API key above.`;
      if (!healthResponse.isFetching && !health.valid)
        return `${providerLabel} API key is invalid — replace it above.`;
      return undefined;
    }
    if (!lastFourFor(providerId))
      return `${providerLabel} is not connected — add an API key above.`;
    return undefined;
  };

  if (
    settingsResponse.isFetching ||
    modelsResponse.isFetching ||
    providersResponse.isFetching
  ) {
    return (
      <div className="flex items-center gap-2 text-sm text-muted-foreground">
        <Spinner />
        Loading AI settings…
      </div>
    );
  }

  const textGenerationModels = allModels.filter(
    (m) => m.supportsTextGeneration,
  );
  const embeddingModels = allModels.filter((m) => m.supportsEmbeddings);

  return (
    <div className="space-y-6">
      <PageHeader
        title="AI"
        description="Configure AI models and API keys for your workspace."
      />

      <div className="space-y-6">
        {/* Providers */}
        <div className="space-y-3">
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {allProviders.map((provider) => (
              <ProviderCard
                key={provider.id}
                provider={provider}
                lastFour={lastFourFor(provider.id)}
                health={healthFor(provider.id)}
                checking={healthResponse.isFetching}
                required={activeProviders.has(provider.id)}
                onSaveKey={(key) => handleSaveKey(provider.id, key)}
                onRemove={() => handleRemoveKey(provider.id)}
              />
            ))}
          </div>
        </div>

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
              providerIcons={providerIcons}
              warning={warningForModel(indexingModel)}
            />
            <ModelRow
              id="newsletter-model"
              label="Newsletter Model"
              description="Drafts newsletter topics and compiles content from indexed codebase activity."
              value={newsletterModel}
              onChange={setNewsletterModel}
              models={textGenerationModels}
              providerIcons={providerIcons}
              warning={warningForModel(newsletterModel)}
            />
            <ModelRow
              id="embeddings-model"
              label="Embeddings Model"
              description="Converts code and text into vectors for semantic search."
              value={embeddingsModel}
              onChange={setEmbeddingsModel}
              models={embeddingModels}
              providerIcons={providerIcons}
              warning={warningForModel(embeddingsModel)}
            />
          </div>
        </div>

        {/* Prompt enhancements — layered on top of fixed built-in system prompts */}
        <div className="rounded-lg border">
          <div className="px-5 py-4 border-b space-y-2">
            <h2 className="text-sm font-semibold">Prompt enhancements</h2>
            <p className="text-xs text-muted-foreground">
              Add optional guidance for indexing, topic selection, and
              newsletter writing. Built-in instructions always run first; your
              text is appended on top for voice, focus, and selection criteria.
            </p>
            <p className="flex items-start gap-1.5 text-xs text-muted-foreground">
              <Icon type="info" className="text-sm shrink-0 mt-0.5" />
              <span>
                Response structure and output format are fixed by the system and
                cannot be changed here. Leave a field empty to use the built-in
                prompt only.
              </span>
            </p>
          </div>
          <div className="divide-y">
            {SYSTEM_PROMPT_FIELDS.map(
              ({ key, label, description, placeholder }) => {
                const custom = isPromptCustom(key);
                const dirty = isPromptDirty(key);
                const hasDraft = promptDrafts[key].trim().length > 0;
                return (
                  <div key={key} className="px-5 py-4 space-y-3">
                    <div className="flex items-start justify-between gap-4">
                      <div className="space-y-0.5 min-w-0">
                        <div className="flex items-center gap-2">
                          <Label
                            htmlFor={`prompt-${key}`}
                            className="text-sm font-medium"
                          >
                            {label}
                          </Label>
                          {custom ? (
                            <Badge variant="secondary" size="sm">
                              Enhanced
                            </Badge>
                          ) : (
                            <Badge
                              variant="outline"
                              size="sm"
                              className="text-muted-foreground"
                            >
                              Default only
                            </Badge>
                          )}
                          {dirty ? (
                            <Badge variant="outline" size="sm">
                              Unsaved
                            </Badge>
                          ) : null}
                        </div>
                        <p className="text-xs text-muted-foreground">
                          {description}
                        </p>
                      </div>
                      <Button
                        type="button"
                        size="sm"
                        variant="ghost"
                        disabled={
                          !settingsLoaded || isSaving || (!hasDraft && !custom)
                        }
                        onClick={() => handleClearPrompt(key)}
                      >
                        Clear
                      </Button>
                    </div>
                    <Textarea
                      id={`prompt-${key}`}
                      value={promptDrafts[key]}
                      onChange={(e) =>
                        setPromptDrafts((prev) => ({
                          ...prev,
                          [key]: e.target.value,
                        }))
                      }
                      rows={4}
                      placeholder={placeholder}
                      className="font-mono text-xs leading-relaxed"
                      disabled={isSaving || !settingsLoaded}
                    />
                  </div>
                );
              },
            )}
          </div>
        </div>
      </div>

      <FloatingBar
        visible={isDirty}
        content={
          <p className="text-sm text-primary-foreground/80">
            You have unsaved changes
          </p>
        }
        cancelButton={
          <Button
            variant="ghost"
            size="sm"
            className="text-primary-foreground hover:text-primary-foreground hover:bg-primary-foreground/10"
            onClick={handleDiscardChanges}
            disabled={isSaving}
          >
            Discard
          </Button>
        }
        saveButton={
          <Button
            variant="secondary"
            size="sm"
            onClick={() => void handleSave()}
            loading={isSaving}
            disabled={!settingsLoaded}
          >
            Save Changes
          </Button>
        }
      />
    </div>
  );
}
