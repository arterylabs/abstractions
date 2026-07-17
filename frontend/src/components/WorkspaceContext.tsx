"use client";

import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { syncWorkspaces } from "@/store/modules/auth";
import { useEffect } from "react";
import Spinner from "./ui/spinner";
import { Button } from "./ui/button";

export default function WorkspaceContext({
  children,
}: {
  children: React.ReactNode;
}) {
  const dispatch = useAppDispatch();

  const workspacesStatus = useAppSelector(
    (state) => state.auth.workspacesStatus,
  );

  useEffect(() => {
    if (workspacesStatus === "idle") {
      dispatch(syncWorkspaces());
    }
  }, [dispatch, workspacesStatus]);

  if (workspacesStatus === "loading" || workspacesStatus === "idle") {
    return (
      <div className="flex h-screen w-screen items-center justify-center bg-background">
        <Spinner />
      </div>
    );
  }

  if (workspacesStatus === "failed") {
    return (
      <div className="flex h-screen w-screen flex-col items-center justify-center gap-4 bg-background">
        <p className="text-sm text-muted-foreground">
          Failed to load your workspaces.
        </p>
        <Button variant="outline" onClick={() => dispatch(syncWorkspaces())}>
          Retry
        </Button>
      </div>
    );
  }

  return <>{children}</>;
}
