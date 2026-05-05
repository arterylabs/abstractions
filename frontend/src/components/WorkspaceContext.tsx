"use client";

import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { syncWorkspaces } from "@/store/modules/auth";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import Spinner from "./ui/spinner";

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

  // If workspaces failed to load, could add error handling here
  if (workspacesStatus === "failed") {
    return null;
  }

  // If workspaces exist, render children
  return <>{children}</>;
}
