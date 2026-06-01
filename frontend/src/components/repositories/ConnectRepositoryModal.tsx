"use client";

import { useCallback, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
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
import { Notifier } from "@/components/ui/sonner";
import Spinner from "@/components/ui/spinner";
import { useAppSelector } from "@/store/hooks";

type Step = "checking" | "oauth";

interface ConnectRepositoryModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function ConnectRepositoryModal({
  open,
  onOpenChange,
}: ConnectRepositoryModalProps) {
  const router = useRouter();
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );

  const [step, setStep] = useState<Step>("checking");

  const [oauthResponse, initiateOAuth] = useAPI(api.github.initiateOAuth);

  useEffect(() => {
    if (oauthResponse.result?.authorizeUrl) {
      window.location.href = oauthResponse.result.authorizeUrl;
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
      await api.github.listRepos(currentWorkspace.id);
      onOpenChange(false);
      router.push("/repositories/connect");
    } catch {
      setStep("oauth");
    }
  }, [currentWorkspace?.id, onOpenChange, router]);

  useEffect(() => {
    if (open) {
      void checkGitHubAndLoadRepos();
    }
  }, [open, checkGitHubAndLoadRepos]);

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
                Authorize access to your GitHub account so we can read your
                repositories.
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-4 py-2">
              <span className="inline-flex items-center gap-1.5 rounded-full bg-primary/10 px-2.5 py-0.5 text-xs font-medium text-primary">
                Step 1 of 2
              </span>

              <div className="space-y-3">
                <p className="text-sm font-medium">How it works</p>
                <ol className="list-none space-y-2 text-sm">
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
                      We index your codebase and generate engineering
                      newsletters so your team always knows what shipped —
                      without anyone having to write it.
                    </span>
                  </li>
                </ol>
              </div>

              <Separator />

              <div className="space-y-2">
                <p className="text-sm font-medium">Is it safe?</p>
                <ul className="space-y-1.5 text-sm">
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
                      We do not store your actual code. Only embeddings and a
                      short summary of what your code does are saved. Your code
                      is only processed to generate these and is never stored.
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
      </DialogContent>
    </Dialog>
  );
}
