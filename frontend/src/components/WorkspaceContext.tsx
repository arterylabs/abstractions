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
      <div className="h-screen w-screen flex items-center justify-center">
        <Spinner />
      </div>
    );
  }

  if (workspacesStatus === "failed") {
    return (
      <div className="h-screen w-screen flex flex-col items-center justify-center gap-4">
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
