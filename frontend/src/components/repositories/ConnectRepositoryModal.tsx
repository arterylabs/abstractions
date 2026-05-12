"use client";

import { useCallback, useEffect, useState } from "react";
import api from "@/api";
import useAPI from "@/api/hooks";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import { Notifier } from "@/components/ui/sonner";
import Icon from "@/components/ui/icon";
import Spinner from "@/components/ui/spinner";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { fetchRepositories } from "@/store/modules/repositories";
import { getErrorMessage } from "@/lib/errors";
import { GitHubRepo } from "@/api/github";

type Step = "checking" | "oauth" | "selecting";

interface ConnectRepositoryModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function ConnectRepositoryModal({
  open,
  onOpenChange,
}: ConnectRepositoryModalProps) {
  const dispatch = useAppDispatch();
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );

  const [step, setStep] = useState<Step>("checking");
  const [repos, setRepos] = useState<GitHubRepo[]>([]);
  const [connectingRepoId, setConnectingRepoId] = useState<number | null>(null);

  const [oauthResponse, initiateOAuth] = useAPI(api.github.initiateOAuth);

  useEffect(() => {
    if (oauthResponse.result?.authorize_url) {
      window.location.href = oauthResponse.result.authorize_url;
    }
  }, [oauthResponse.result]);

  useEffect(() => {
    if (oauthResponse.error) {
      Notifier.error({
        title: "Unable to connect GitHub",
        message: oauthResponse.error.message,
      });
    }
  }, [oauthResponse.error]);

  const checkGitHubAndLoadRepos = useCallback(async () => {
    if (!currentWorkspace?.id) return;
    setStep("checking");
    try {
      const result = await api.github.listRepos(currentWorkspace.id);
      setRepos(result.data);
      setStep("selecting");
    } catch {
      setStep("oauth");
    }
  }, [currentWorkspace?.id]);

  useEffect(() => {
    if (open) {
      void checkGitHubAndLoadRepos();
    }
  }, [open, checkGitHubAndLoadRepos]);

  const handleConnectRepository = async (repo: GitHubRepo) => {
    if (!currentWorkspace?.id) return;
    setConnectingRepoId(repo.id);
    try {
      await api.github.connectRepository(currentWorkspace.id, {
        githubRepoId: repo.id,
        fullName: repo.fullName,
        description: repo.description,
        defaultBranch: repo.defaultBranch,
      });
      dispatch(fetchRepositories(currentWorkspace.id));
      onOpenChange(false);
    } catch (error: unknown) {
      Notifier.error({
        title: "Unable to connect repository",
        message: getErrorMessage(error),
      });
    } finally {
      setConnectingRepoId(null);
    }
  };

  const handleOAuth = async () => {
    if (!currentWorkspace?.id) return;
    await initiateOAuth(currentWorkspace.id, { createRepository: true });
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-lg">
        {step === "checking" && (
          <>
            <DialogHeader>
              <DialogTitle>Connect a Repository</DialogTitle>
              <DialogDescription>
                Checking your GitHub connection…
              </DialogDescription>
            </DialogHeader>
            <div className="flex min-h-[160px] items-center justify-center">
              <Spinner />
            </div>
          </>
        )}

        {step === "oauth" && (
          <>
            <DialogHeader>
              <DialogTitle>Connect a GitHub Repository</DialogTitle>
              <DialogDescription>
                <span className="inline-flex items-center gap-1.5 rounded-full bg-primary/10 px-2.5 py-0.5 text-xs font-medium text-primary">
                  Step 1 of 2
                </span>{" "}
                Authorize access to your GitHub account so we can read your
                repositories.
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-4 py-2">
              <div className="space-y-3">
                <p className="text-sm font-medium">How it works</p>
                <ol className="list-none space-y-2 text-sm text-muted-foreground">
                  <li className="flex gap-3">
                    <span className="flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-primary/10 text-xs font-semibold text-primary">
                      1
                    </span>
                    <span>
                      Click <strong>Connect to GitHub</strong> to authorize read
                      access to your repositories.
                    </span>
                  </li>
                  <li className="flex gap-3">
                    <span className="flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-primary/10 text-xs font-semibold text-primary">
                      2
                    </span>
                    <span>
                      Select the repository you want to connect from the list.
                    </span>
                  </li>
                  <li className="flex gap-3">
                    <span className="flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-primary/10 text-xs font-semibold text-primary">
                      3
                    </span>
                    <span>
                      We index your code and start generating automated
                      newsletters for your team.
                    </span>
                  </li>
                </ol>
              </div>

              <Separator />

              <div className="space-y-2">
                <p className="text-sm font-medium">Our security promise</p>
                <ul className="space-y-1.5 text-sm text-muted-foreground">
                  <li className="flex gap-2">
                    <span className="mt-0.5 text-success">✓</span>
                    <span>
                      We request <strong>read-only</strong> access — we never
                      push code or modify your repository.
                    </span>
                  </li>
                  <li className="flex gap-2">
                    <span className="mt-0.5 text-success">✓</span>
                    <span>
                      Your OAuth token is encrypted at rest and never exposed
                      outside our servers.
                    </span>
                  </li>
                  <li className="flex gap-2">
                    <span className="mt-0.5 text-success">✓</span>
                    <span>
                      You can revoke access at any time from your GitHub account
                      settings.
                    </span>
                  </li>
                </ul>
              </div>
            </div>

            <DialogFooter className="flex-col-reverse gap-2 sm:flex-row">
              <Button
                variant="outline"
                onClick={() => onOpenChange(false)}
                disabled={oauthResponse.isFetching}
              >
                Cancel
              </Button>
              <Button
                onClick={() => void handleOAuth()}
                loading={oauthResponse.isFetching}
                disabled={!currentWorkspace?.id}
              >
                Connect to GitHub
              </Button>
            </DialogFooter>
          </>
        )}

        {step === "selecting" && (
          <>
            <DialogHeader>
              <DialogTitle>Connect a Repository</DialogTitle>
              <DialogDescription>
                <span className="inline-flex items-center gap-1.5 rounded-full bg-success/10 px-2.5 py-0.5 text-xs font-medium text-success">
                  GitHub connected
                </span>{" "}
                Choose which repository to add.
              </DialogDescription>
            </DialogHeader>

            <div className="max-h-[400px] overflow-y-auto pr-1 -mr-1">
              {repos.length === 0 ? (
                <p className="py-8 text-center text-sm text-muted-foreground">
                  No repositories found on your GitHub account.
                </p>
              ) : (
                <div className="space-y-2">
                  {repos.map((repo) => (
                    <div
                      key={repo.id}
                      className="flex items-start justify-between gap-3 rounded-lg border p-3"
                    >
                      <div className="flex min-w-0 items-start gap-2.5">
                        <Icon
                          type="deployed_code"
                          className="mt-0.5 shrink-0 text-base text-muted-foreground"
                        />
                        <div className="min-w-0 space-y-0.5">
                          <div className="flex flex-wrap items-center gap-1.5">
                            <span className="truncate text-sm font-medium">
                              {repo.fullName}
                            </span>
                            {repo.private && (
                              <Badge variant="secondary" className="text-xs">
                                Private
                              </Badge>
                            )}
                          </div>
                          {repo.description && (
                            <p className="truncate text-xs text-muted-foreground">
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
                        size="sm"
                        className="shrink-0"
                        onClick={() => void handleConnectRepository(repo)}
                        loading={connectingRepoId === repo.id}
                        disabled={
                          connectingRepoId !== null &&
                          connectingRepoId !== repo.id
                        }
                      >
                        Connect
                      </Button>
                    </div>
                  ))}
                </div>
              )}
            </div>

            <DialogFooter>
              <Button
                variant="outline"
                onClick={() => onOpenChange(false)}
                disabled={connectingRepoId !== null}
              >
                Cancel
              </Button>
            </DialogFooter>
          </>
        )}
      </DialogContent>
    </Dialog>
  );
}
