"use client";

import { useCallback, useEffect } from "react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import Icon from "@/components/ui/icon";
import api from "@/api";
import useAPI from "@/api/hooks";
import Spinner from "@/components/ui/spinner";
import { useAppSelector, useAppDispatch } from "@/store/hooks";
import { fetchRoles } from "@/store/modules/auth";
import { format } from "date-fns";
import { copyToClipboard } from "@/lib/clipboard";

interface UserWorkspacesTabProps {
  userId: string | null;
}

interface WorkspaceWithRole {
  id: string;
  name: string;
  logo?: string;
  website?: string;
  about?: string;
  billingEmail: string;
  createdAt: string;
  updatedAt: string;
  userRole: string;
}

export default function UserWorkspacesTab({ userId }: UserWorkspacesTabProps) {
  const [workspacesResponse, getUserWorkspaces] = useAPI(
    api.admin.getUserWorkspacesWithRole,
  );

  // Get roles from store and dispatch
  const dispatch = useAppDispatch();
  const roles = useAppSelector((state) => state.auth.roles);
  const rolesStatus = useAppSelector((state) => state.auth.rolesStatus);

  const loadUserWorkspaces = useCallback(async () => {
    if (userId) {
      await getUserWorkspaces(userId);
    }
  }, [getUserWorkspaces, userId]);

  useEffect(() => {
    loadUserWorkspaces();
  }, [loadUserWorkspaces]);

  // Helper function to get role display name
  const getRoleName = (roleId: string) => {
    if (!roles || roles.length === 0) {
      return roleId; // Fallback to roleId if roles not loaded
    }
    const role = roles.find((r) => r.id === roleId);
    return role?.name || roleId;
  };

  const formatDate = (dateString: string | undefined) => {
    if (!dateString) return "N/A";
    try {
      return format(new Date(dateString), "PPp"); // e.g., "Apr 29, 2023, 3:45 PM"
    } catch {
      return "Invalid date";
    }
  };

  // Load roles for the first workspace (since roles are workspace-scoped but generally the same)
  useEffect(() => {
    if (
      workspacesResponse.result &&
      workspacesResponse.result.length > 0 &&
      rolesStatus === "idle"
    ) {
      const firstWorkspaceId = workspacesResponse.result[0].id;
      dispatch(fetchRoles(firstWorkspaceId));
    }
  }, [workspacesResponse.result, rolesStatus, dispatch]);

  if (workspacesResponse.isFetching) {
    return (
      <div className="flex items-center justify-center py-12">
        <Spinner />
      </div>
    );
  }

  const workspaces = workspacesResponse.result || [];

  if (workspaces.length === 0) {
    return (
      <div className="text-center py-12">
        <p className="text-muted-foreground">
          User is not a member of any workspaces
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="space-y-4">
        {workspaces.map((workspace: WorkspaceWithRole) => (
          <div
            key={workspace.id}
            className="border border-border rounded-lg p-4 space-y-4"
          >
            {/* Workspace Header */}
            <div className="flex items-start justify-between">
              <div className="flex items-center gap-3">
                {workspace.logo ? (
                  // eslint-disable-next-line @next/next/no-img-element
                  <img
                    src={workspace.logo}
                    alt={workspace.name}
                    width={40}
                    height={40}
                    className="w-10 h-10 rounded-lg object-cover"
                  />
                ) : (
                  <div className="w-10 h-10 rounded-lg bg-accent flex items-center justify-center">
                    <span className="text-sm font-medium text-muted-foreground">
                      {workspace.name.charAt(0).toUpperCase()}
                    </span>
                  </div>
                )}
                <div>
                  <div className="flex items-center gap-2">
                    <h4 className="text-lg font-semibold text-foreground">
                      {workspace.name}
                    </h4>
                    <Button
                      variant="ghost"
                      size="sm"
                      className="h-6 px-2 text-xs text-muted-foreground hover:text-foreground"
                      onClick={() =>
                        copyToClipboard(workspace.id, "Workspace ID")
                      }
                    >
                      <Icon type="content_copy" className="text-xs" />
                    </Button>
                  </div>
                  {workspace.website && (
                    <p className="text-sm text-muted-foreground">
                      {workspace.website}
                    </p>
                  )}
                </div>
              </div>
            </div>

            {/* Workspace Metadata */}
            <Badge variant="outline" className="text-sm px-3 py-1 mr-3">
              {getRoleName(workspace.userRole)}
            </Badge>
            <span className="text-xs text-muted-foreground">
              Created {formatDate(workspace.createdAt)}
            </span>
          </div>
        ))}
      </div>
      {/* Summary */}
      <div className="pt-4 border-t">
        <p className="text-sm text-muted-foreground">
          Total workspaces:{" "}
          <span className="font-medium">{workspaces.length}</span>
        </p>
      </div>
    </div>
  );
}
