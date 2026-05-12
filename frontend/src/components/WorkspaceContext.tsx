"use client";

import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { syncWorkspaces } from "@/store/modules/auth";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import Spinner from "./ui/spinner";
import { Button } from "./ui/button";

export default function WorkspaceContext({
  children,
}: {
  children: React.ReactNode;
}) {
  const dispatch = useAppDispatch();
  const router = useRouter();

  const workspacesStatus = useAppSelector(
    (state) => state.auth.workspacesStatus,
  );
  const workspaces = useAppSelector((state) => state.auth.workspaces);

  useEffect(() => {
    if (workspacesStatus === "idle") {
      dispatch(syncWorkspaces());
    }
  }, [dispatch, workspacesStatus]);

  // If workspaces are loaded and there are no workspaces, redirect to create
  useEffect(() => {
    if (
      workspacesStatus === "success" &&
      (!workspaces || workspaces.length === 0)
    ) {
      router.push("/workspaces/create");
    }
  }, [workspaces, workspacesStatus, router]);

  // Show loading spinner while fetching workspaces
  if (workspacesStatus === "loading" || workspacesStatus === "idle") {
    return (
      <div className="h-screen w-screen flex items-center justify-center">
        <Spinner />
      </div>
    );
  }

  // If workspaces failed to load, show an error with a retry option
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

  // If workspaces exist, render children
  return <>{children}</>;
}
