"use client";

import { useEffect, useState } from "react";
import api from "@/api";
import useAPI from "@/api/hooks";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { EmptyState } from "@/components/ui/empty-state";
import { Input } from "@/components/ui/input";
import { PageHeader } from "@/components/ui/page-header";
import Spinner from "@/components/ui/spinner";
import { Notifier } from "@/components/ui/sonner";
import Icon from "@/components/ui/icon";
import { useAppSelector } from "@/store/hooks";
import { getErrorMessage } from "@/lib/errors";
import { useRouter } from "next/navigation";
import { AIConnectionWarning } from "@/components/ai-connection-warning";

function StepIndicator() {
  return (
    // Steps set as folios: § number, running head, current step inked in.
    <ol className="flex items-center gap-4">
      <li className="press-label flex items-center gap-2 text-muted-foreground">
        <Icon type="check" className="!text-sm text-success-dark" />
        Connect GitHub
      </li>
      <li aria-hidden className="h-px w-6 bg-border" />
      <li className="press-label flex items-center gap-2 text-primary">
        § 02 Select repository
      </li>
    </ol>
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
  const [searchQuery, setSearchQuery] = useState("");

  useEffect(() => {
    if (currentWorkspace?.id) {
      void loadRepos(currentWorkspace.id);
    }
  }, [currentWorkspace?.id, loadRepos]);

  useEffect(() => {
    if (oauthResponse.result?.authorizeUrl) {
      window.location.href = oauthResponse.result.authorizeUrl;
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
      const response = await api.github.connectRepository(currentWorkspace.id, {
        githubRepoId: repository.id,
        fullName: repository.fullName,
        description: repository.description,
        defaultBranch: repository.defaultBranch,
      });
      Notifier.success({
        title: "Repository connected",
        message: `${repository.fullName} has been connected successfully.`,
      });
      router.push(`/repositories/connect/${response.data.id}`);
    } catch (error: unknown) {
      Notifier.error({
        title: "Unable to connect repository",
        message: getErrorMessage(error),
      });
    } finally {
      setConnectingRepoId(null);
    }
  };

  const filteredRepos = reposResponse.result?.filter((repo) =>
    repo.fullName.toLowerCase().includes(searchQuery.toLowerCase()),
  );

  return (
    <div className="flex flex-col gap-8">
      <PageHeader
        eyebrow="Connect"
        title="Select a repository"
        description="Pick the repository you want Abstractions to track. We read your code to understand what's shipping — we never store it, and we never write to your repo."
        showBackButton
        handleBack={() => router.push("/repositories")}
      >
        <StepIndicator />
      </PageHeader>

      <AIConnectionWarning />

      {reposResponse.isFetching || !currentWorkspace?.id ? (
        <div className="flex min-h-[240px] flex-col items-center justify-center gap-4">
          <Spinner />
          <p className="press-label text-muted-foreground">Setting the run…</p>
        </div>
      ) : reposResponse.result && reposResponse.result.length > 0 ? (
        <div className="flex flex-col gap-4">
          <Input
            id="repo-search"
            placeholder="Search repositories..."
            value={searchQuery}
            onChange={(value) => setSearchQuery(value)}
          />
          {filteredRepos && filteredRepos.length > 0 ? (
            <div className="grid gap-4">
              {filteredRepos.map((repo) => (
                <div
                  key={repo.id}
                  className="press-card flex items-start justify-between gap-4 overflow-hidden border border-border bg-card p-5"
                >
                  <div className="flex min-w-0 items-start gap-3">
                    <Icon
                      type="deployed_code"
                      className="mt-1.5 shrink-0 !text-base text-primary"
                    />
                    <div className="flex min-w-0 flex-col gap-1.5">
                      <div className="flex flex-wrap items-center gap-2.5">
                        <h3 className="truncate text-xl font-medium">
                          {repo.fullName}
                        </h3>
                        {repo.private && (
                          <Badge size="sm" variant="secondary">
                            <Icon type="lock" className="!text-[10px]" />
                            Private
                          </Badge>
                        )}
                      </div>
                      {repo.description && (
                        <p className="line-clamp-2 max-w-[68ch] text-sm leading-relaxed text-muted-foreground">
                          {repo.description}
                        </p>
                      )}
                      <span className="mt-0.5 font-mono text-xs text-muted-foreground">
                        {repo.defaultBranch}
                      </span>
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
          ) : (
            <p className="py-8 text-center text-sm text-muted-foreground">
              No repositories match your search.
            </p>
          )}
        </div>
      ) : reposResponse.error ? (
        <EmptyState
          icon="link_off"
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
          icon="search_off"
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
