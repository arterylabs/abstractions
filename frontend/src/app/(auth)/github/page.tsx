"use client";

import { useEffect, useMemo, useState } from "react";
import api, { axiosInstance } from "@/api";
import useAPI from "@/api/hooks";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import Spinner from "@/components/ui/spinner";
import { Notifier } from "@/components/ui/sonner";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { syncCurrentUser, syncWorkspaces } from "@/store/modules/auth";
import { getErrorMessage } from "@/lib/errors";
import { useRouter, useSearchParams } from "next/navigation";

export default function GitHubConnectPage() {
  const dispatch = useAppDispatch();
  const router = useRouter();
  const searchParams = useSearchParams();
  const code = searchParams.get("code");
  const state = searchParams.get("state");

  const currentWorkspace = useAppSelector(
    (store) => store.auth.currentWorkspace,
  );
  const currentUserStatus = useAppSelector(
    (store) => store.auth.currentUserStatus,
  );
  const workspacesStatus = useAppSelector(
    (store) => store.auth.workspacesStatus,
  );

  const [oauthResponse, initiateOAuth] = useAPI(api.github.initiateOAuth);
  const [callbackError, setCallbackError] = useState<string | null>(null);
  const [isConnecting, setIsConnecting] = useState(false);

  useEffect(() => {
    if (currentUserStatus === "idle") {
      dispatch(syncCurrentUser());
    }
  }, [currentUserStatus, dispatch]);

  useEffect(() => {
    if (currentUserStatus === "success" && workspacesStatus === "idle") {
      dispatch(syncWorkspaces());
    }
  }, [currentUserStatus, dispatch, workspacesStatus]);

  useEffect(() => {
    if (!code) {
      return;
    }

    let isMounted = true;

    const connectGitHub = async () => {
      setIsConnecting(true);
      setCallbackError(null);

      try {
        await axiosInstance.get("/integrations/callback/github", {
          params: {
            code,
            state: state ?? undefined,
          },
        });

        if (isMounted) {
          router.push("/dashboard/repositories/connect");
        }
      } catch (error: unknown) {
        if (isMounted) {
          setCallbackError(getErrorMessage(error));
        }
      } finally {
        if (isMounted) {
          setIsConnecting(false);
        }
      }
    };

    void connectGitHub();

    return () => {
      isMounted = false;
    };
  }, [code, router, state]);

  useEffect(() => {
    if (oauthResponse.result?.url) {
      window.location.href = oauthResponse.result.url;
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

  const isLoadingWorkspace = useMemo(
    () =>
      !code && (currentUserStatus === "idle" || workspacesStatus === "loading"),
    [code, currentUserStatus, workspacesStatus],
  );

  const handleConnect = async () => {
    if (!currentWorkspace?.id) {
      setCallbackError("We could not find a workspace to connect.");
      return;
    }

    await initiateOAuth(currentWorkspace.id);
  };

  return (
    <div className="flex flex-col gap-6">
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-xl">Connect GitHub</CardTitle>
          <CardDescription>
            Connect your GitHub account to start syncing repositories.
          </CardDescription>
        </CardHeader>
        <CardContent>
          {code && !callbackError ? (
            <div className="flex flex-col items-center gap-4 py-6 text-center">
              <Spinner />
              <p className="text-sm text-muted-foreground">
                {isConnecting ? "Connecting GitHub..." : "Redirecting..."}
              </p>
            </div>
          ) : isLoadingWorkspace ? (
            <div className="flex flex-col items-center gap-4 py-6 text-center">
              <Spinner />
              <p className="text-sm text-muted-foreground">
                Loading workspace...
              </p>
            </div>
          ) : (
            <div className="grid gap-4">
              <Button
                className="w-full"
                onClick={handleConnect}
                loading={oauthResponse.isFetching}
                disabled={!currentWorkspace?.id}
              >
                Connect GitHub
              </Button>
              {callbackError && (
                <p className="text-sm text-destructive text-center">
                  {callbackError}
                </p>
              )}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
