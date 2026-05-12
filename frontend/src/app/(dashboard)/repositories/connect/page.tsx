"use client";

import { useEffect, useState } from "react";
import api from "@/api";
import useAPI from "@/api/hooks";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { EmptyState } from "@/components/ui/empty-state";
import { PageHeader } from "@/components/ui/page-header";
import Spinner from "@/components/ui/spinner";
import { Notifier } from "@/components/ui/sonner";
import Icon from "@/components/ui/icon";
import { useAppSelector } from "@/store/hooks";
import { getErrorMessage } from "@/lib/errors";
import { useRouter } from "next/navigation";

function StepIndicator() {
  return (
    <div className="flex items-center gap-2 text-sm">
      <div className="flex items-center gap-1.5">
        <span className="flex h-5 w-5 items-center justify-center rounded-full bg-success text-xs font-semibold text-success-foreground">
          <Icon type="check" className="text-xs" />
        </span>
        <span className="text-muted-foreground">Connect GitHub</span>
      </div>
      <span className="text-muted-foreground">→</span>
      <div className="flex items-center gap-1.5">
        <span className="flex h-5 w-5 items-center justify-center rounded-full bg-primary text-xs font-semibold text-primary-foreground">
          2
        </span>
        <span className="font-medium">Select Repository</span>
      </div>
    </div>
  );
}

export default function ConnectPage() {
  const router = useRouter();
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );
  const [reposResponse, loadRepos] = useAPI(api.github.listRepos);
  const [oauthResponse, initiateOAuth] = useAPI(api.github.initiateOAuth);
  const [connectingRepoId, setConnectingRepoId] = useState<number | null>(null);

  useEffect(() => {
    if (currentWorkspace?.id) {
      void loadRepos(currentWorkspace.id);
    }
  }, [currentWorkspace?.id, loadRepos]);

  useEffect(() => {
    if (oauthResponse.result?.authorize_url) {
      window.location.href = oauthResponse.result.authorize_url;
    }
  }, [oauthResponse.result]);

  const handleConnectRepository = async (repoId: number) => {
    if (!currentWorkspace?.id || !reposResponse.result) {
      return;
    }

    const repository = reposResponse.result.find((item) => item.id === repoId);
    if (!repository) {
      return;
    }

    setConnectingRepoId(repoId);

    try {
      await api.github.connectRepository(currentWorkspace.id, {
        githubRepoId: repository.id,
        fullName: repository.fullName,
        description: repository.description,
        defaultBranch: repository.defaultBranch,
      });
      router.push("/repositories");
    } catch (error: unknown) {
      Notifier.error({
        title: "Unable to connect repository",
        message: getErrorMessage(error),
      });
    } finally {
      setConnectingRepoId(null);
    }
  };

  return (
    <div className="space-y-6">
      <PageHeader
        title="Select a Repository"
        showBackButton
        handleBack={() => router.push("/repositories")}
      >
        <StepIndicator />
      </PageHeader>

      <p className="text-sm text-muted-foreground">
        Choose the repository you want to index and generate newsletters for. We
        only request read-only access.
      </p>

      {reposResponse.isFetching || !currentWorkspace?.id ? (
        <div className="flex min-h-[240px] flex-col items-center justify-center gap-4">
          <Spinner />
          <p className="text-sm text-muted-foreground">
            Loading repositories...
          </p>
        </div>
      ) : reposResponse.result && reposResponse.result.length > 0 ? (
        <div className="grid gap-3">
          {reposResponse.result.map((repo) => (
            <div
              key={repo.id}
              className="flex items-start justify-between gap-4 rounded-lg border bg-card p-4"
            >
              <div className="flex min-w-0 items-start gap-3">
                <Icon
                  type="deployed_code"
                  className="mt-0.5 shrink-0 text-xl text-muted-foreground"
                />
                <div className="min-w-0 space-y-1">
                  <div className="flex flex-wrap items-center gap-2">
                    <span className="truncate font-medium">
                      {repo.fullName}
                    </span>
                    {repo.private && <Badge variant="secondary">Private</Badge>}
                  </div>
                  {repo.description && (
                    <p className="text-sm text-muted-foreground line-clamp-2">
                      {repo.description}
                    </p>
                  )}
                  <div className="flex items-center gap-1 text-xs text-muted-foreground">
                    <Icon type="account_tree" className="text-xs" />
                    <span>{repo.defaultBranch}</span>
                  </div>
                </div>
              </div>
              <Button
                className="shrink-0"
                onClick={() => void handleConnectRepository(repo.id)}
                loading={connectingRepoId === repo.id}
                disabled={
                  connectingRepoId !== null && connectingRepoId !== repo.id
                }
              >
                Connect
              </Button>
            </div>
          ))}
        </div>
      ) : reposResponse.error ? (
        <EmptyState
          title="GitHub not connected"
          description="Connect your GitHub account to access your repositories and start syncing code."
          cta={
            <Button
              onClick={() =>
                void initiateOAuth(currentWorkspace.id, {
                  createRepository: true,
                })
              }
              loading={oauthResponse.isFetching}
            >
              Connect GitHub
            </Button>
          }
        />
      ) : (
        <EmptyState
          title="No repositories found"
          description="We couldn't find any GitHub repositories available to connect."
          cta={
            <Button
              variant="outline"
              onClick={() => void loadRepos(currentWorkspace.id)}
            >
              Refresh
            </Button>
          }
        />
      )}
    </div>
  );
}
