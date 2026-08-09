"use client";

import { isAxiosError } from "axios";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Notifier } from "@/components/ui/sonner";
import Spinner from "@/components/ui/spinner";
import Icon from "@/components/ui/icon";
import Tooltip from "@/components/ui/tooltip";
import { WorkspaceProviderHealth } from "@/api/workspaces";
import { AIProvider } from "@/types/workspace";
import { cn } from "@/lib/utils";

// Badge overlaid on the provider logo: green check or red error once health has run.
function ProviderStatusIndicator({
  health,
}: {
  health?: WorkspaceProviderHealth;
}) {
  if (!health?.configured) return null;
  if (health.valid) {
    return (
      <span className="absolute -bottom-1.5 -right-1.5 flex items-center justify-center rounded-full bg-card">
        <Icon
          type="check_circle"
          className="text-xl leading-none text-success-dark"
        />
      </span>
    );
  }
  return (
    <Tooltip value={health.error || "API key is invalid."}>
      <span className="absolute -bottom-1.5 -right-1.5 flex items-center justify-center rounded-full bg-card">
        <Icon
          type="error"
          className="text-xl leading-none text-destructive-dark"
        />
      </span>
    </Tooltip>
  );
}

export function ProviderCard({
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
  // Persists the key and refreshes settings + health; throws on failure.
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
      className={cn(
        "press-card flex flex-col gap-4 overflow-hidden p-5",
        !configured && "bg-accent-dark",
      )}
    >
      <div className="flex items-start justify-between gap-3">
        <div className="relative shrink-0">
          <div
            className={cn(
              "flex size-12 items-center justify-center rounded-sm border border-border bg-background",
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

      <div className="flex flex-1 flex-col gap-1.5">
        <h3 className="text-xl font-medium">{provider.displayName}</h3>
        {configured && lastFour ? (
          <span className="inline-block rounded-sm border border-border bg-accent px-2 py-1 font-mono text-xs text-muted-foreground">
            ••••••••••••{lastFour}
          </span>
        ) : (
          <p className="text-xs text-muted-foreground">
            No API key configured.
          </p>
        )}
        {required && !configured ? (
          <p className="font-mono text-xs text-warning-dark">
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
