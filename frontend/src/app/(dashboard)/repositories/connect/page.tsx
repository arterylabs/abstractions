"use client";

import { useEffect, useState } from "react";
import api from "@/api";
import useAPI from "@/api/hooks";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { EmptyState } from "@/components/ui/empty-state";
import { PageHeader } from "@/components/ui/page-header";
import Spinner from "@/components/ui/spinner";
import { Notifier } from "@/components/ui/sonner";
import { useAppSelector } from "@/store/hooks";
import { getErrorMessage } from "@/lib/errors";
import { useRouter } from "next/navigation";

export default function ConnectPage() {
  const router = useRouter();
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );
  const [reposResponse, loadRepos] = useAPI(api.github.listRepos);
  const [connectingRepoId, setConnectingRepoId] = useState<number | null>(null);

  useEffect(() => {
    if (currentWorkspace?.id) {
      void loadRepos(currentWorkspace.id);
    }
  }, [currentWorkspace?.id, loadRepos]);

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
        defaultBranch: repository.defaultBranch,
      });
      router.push("/dashboard/repositories");
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
        handleBack={() => router.push("/dashboard/repositories")}
      />

      {reposResponse.isFetching || !currentWorkspace?.id ? (
        <div className="flex min-h-[240px] flex-col items-center justify-center gap-4">
          <Spinner />
          <p className="text-sm text-muted-foreground">
            Loading repositories...
          </p>
        </div>
      ) : reposResponse.result && reposResponse.result.length > 0 ? (
        <div className="grid gap-4">
          {reposResponse.result.map((repo) => (
            <Card key={repo.id}>
              <CardHeader className="gap-3 sm:flex-row sm:items-start sm:justify-between">
                <div className="space-y-2">
                  <CardTitle className="text-lg">{repo.fullName}</CardTitle>
                  <CardDescription>
                    {repo.description || "No description provided."}
                  </CardDescription>
                  {repo.private && <Badge variant="secondary">Private</Badge>}
                </div>
                <Button
                  onClick={() => void handleConnectRepository(repo.id)}
                  loading={connectingRepoId === repo.id}
                >
                  Connect
                </Button>
              </CardHeader>
              <CardContent>
                <a
                  href={repo.htmlUrl}
                  target="_blank"
                  rel="noreferrer"
                  className="text-sm text-primary underline-offset-4 hover:underline"
                >
                  View on GitHub
                </a>
              </CardContent>
            </Card>
          ))}
        </div>
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
