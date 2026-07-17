"use client";

import { useEffect } from "react";
import api from "@/api";
import useAPI from "@/api/hooks";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { useAppSelector } from "@/store/hooks";
import { Button } from "./ui/button";
import Icon from "./ui/icon";

export function AIConnectionWarning() {
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );
  const [healthResponse, checkHealth] = useAPI(
    api.workspaces.getWorkspaceHealth,
  );
  const [providersResponse, fetchProviders] = useAPI(
    api.workspaces.getAIProviders,
  );

  useEffect(() => {
    if (currentWorkspace?.id) {
      void checkHealth(currentWorkspace.id);
    }
  }, [currentWorkspace?.id, checkHealth]);

  useEffect(() => {
    void fetchProviders();
  }, [fetchProviders]);

  if (
    healthResponse.isFetching ||
    !healthResponse.result ||
    healthResponse.error
  ) {
    return null;
  }

  const providerNames = new Map(
    (providersResponse.result ?? []).map((p) => [p.id, p.displayName]),
  );

  // Surface an issue for every required provider that is missing or has an invalid key.
  const issues: string[] = [];
  for (const [providerId, status] of Object.entries(
    healthResponse.result.providers,
  )) {
    if (!status.required) continue;
    const name = providerNames.get(providerId) ?? providerId;
    if (!status.configured) {
      issues.push(`${name} API key is not configured`);
    } else if (!status.valid) {
      issues.push(`${name} API key is invalid: ${status.error}`);
    }
  }

  if (issues.length === 0) {
    return null;
  }

  return (
    <Alert variant="warning" className="flex items-center justify-between">
      <div>
        <AlertTitle className="flex items-center gap-1">
          <Icon type="error" /> Warning
        </AlertTitle>
        <AlertDescription className="mt-2">
          <div>
            Without a working connection to your AI models, features like code
            indexing and newsletter generation will not work.
          </div>
          <ul className="">
            {issues.map((issue) => (
              <li key={issue}>{issue}</li>
            ))}
          </ul>
        </AlertDescription>
      </div>
      <Button variant="warning">Configure AI settings</Button>
    </Alert>
  );
}
