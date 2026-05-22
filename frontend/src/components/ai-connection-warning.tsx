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

  useEffect(() => {
    if (currentWorkspace?.id) {
      void checkHealth(currentWorkspace.id);
    }
  }, [currentWorkspace?.id, checkHealth]);

  if (
    healthResponse.isFetching ||
    !healthResponse.result ||
    healthResponse.error
  ) {
    return null;
  }

  const { anthropic, openAI, anthropicRequired, openAIRequired } =
    healthResponse.result;
  const issues: string[] = [];

  if (anthropicRequired && !anthropic.configured) {
    issues.push("Anthropic API key is not configured");
  } else if (anthropicRequired && !anthropic.valid) {
    issues.push(`Anthropic API key is invalid: ${anthropic.error}`);
  }

  if (openAIRequired && !openAI.configured) {
    issues.push("OpenAI API key is not configured");
  } else if (openAIRequired && !openAI.valid) {
    issues.push(`OpenAI API key is invalid: ${openAI.error}`);
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
