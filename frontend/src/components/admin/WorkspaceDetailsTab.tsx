"use client";

import { useState, useCallback, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Badge } from "@/components/ui/badge";
import { AdminWorkspace } from "@/types/admin";
import api from "@/api";
import useAPI from "@/api/hooks";
import Spinner from "@/components/ui/spinner";
import { Notifier } from "@/components/ui/sonner";
import { format } from "date-fns";
import Icon from "@/components/ui/icon";
import { copyToClipboard } from "@/lib/clipboard";
import { askConfirmation } from "@/components/ui/confirmation";

interface WorkspaceDetailsTabProps {
  workspaceId: string | null;
  onWorkspaceUpdated: () => void;
  onWorkspaceDeleted?: () => void;
}

type PendingAction = "save" | "delete" | "recover";

function formatDate(dateString: string | undefined) {
  if (!dateString) return "Invalid date";
  try {
    return format(new Date(dateString), "PPpp"); // e.g., "Apr 29, 2023 at 3:45 PM"
  } catch {
    return "Invalid date";
  }
}

export default function WorkspaceDetailsTab({
  workspaceId,
  onWorkspaceUpdated,
  onWorkspaceDeleted,
}: WorkspaceDetailsTabProps) {
  const [workspaceResponse, getWorkspace] = useAPI(api.admin.getWorkspace);

  const [formData, setFormData] = useState({
    name: "",
    website: "",
    about: "",
    billingEmail: "",
  });
  const [isEditing, setIsEditing] = useState(false);
  const [pendingAction, setPendingAction] = useState<PendingAction | null>(
    null,
  );

  const loadWorkspace = useCallback(async () => {
    if (workspaceId) {
      await getWorkspace(workspaceId);
    }
  }, [getWorkspace, workspaceId]);

  useEffect(() => {
    loadWorkspace();
  }, [loadWorkspace]);

  useEffect(() => {
    if (workspaceResponse.result) {
      setFormData({
        name: workspaceResponse.result.name || "",
        website: workspaceResponse.result.website || "",
        about: workspaceResponse.result.about || "",
        billingEmail: workspaceResponse.result.billingEmail || "",
      });
    }
  }, [workspaceResponse.result]);

  const handleSave = async () => {
    if (!workspaceId) return;

    setPendingAction("save");
    try {
      await api.admin.updateWorkspace(workspaceId, formData);
      Notifier.success("Workspace updated successfully");
      setIsEditing(false);
      onWorkspaceUpdated();
    } catch {
      Notifier.error("Failed to update workspace");
    } finally {
      setPendingAction(null);
    }
  };

  const handleRecover = (workspace: AdminWorkspace) => {
    if (!workspaceId) return;

    askConfirmation({
      title: "Recover Workspace",
      content: `Are you sure you want to recover workspace "${workspace.name}"? This will restore the workspace and make it available again.`,
      variant: "primary",
      okText: "Recover",
      cancelText: "Cancel",
      okCallback: async () => {
        setPendingAction("recover");
        try {
          await api.admin.bulkRecoverWorkspaces([workspaceId]);
          Notifier.success("Workspace recovered successfully");
          await loadWorkspace();
          onWorkspaceUpdated();
        } catch {
          Notifier.error("Failed to recover workspace");
        } finally {
          setPendingAction(null);
        }
      },
    });
  };

  const handleDelete = (workspace: AdminWorkspace) => {
    if (!workspaceId) return;

    askConfirmation({
      title: "Delete Workspace",
      content: `Are you sure you want to delete workspace "${workspace.name}"? This action cannot be undone and will remove all associated data.`,
      variant: "destructive",
      okText: "Delete",
      cancelText: "Cancel",
      okCallback: async () => {
        setPendingAction("delete");
        try {
          await api.admin.bulkDeleteWorkspaces([workspaceId]);
          Notifier.success("Workspace deleted successfully");
          onWorkspaceDeleted?.();
        } catch {
          Notifier.error("Failed to delete workspace");
        } finally {
          setPendingAction(null);
        }
      },
    });
  };

  const workspace = workspaceResponse.result as AdminWorkspace | undefined;

  if (workspaceResponse.isFetching) {
    return (
      <div className="flex items-center justify-center py-12">
        <Spinner />
      </div>
    );
  }

  if (!workspace) {
    return (
      <div className="text-center py-12">
        <p className="text-muted-foreground">Workspace not found</p>
      </div>
    );
  }

  const isBusy = pendingAction !== null;

  return (
    <div className="space-y-8">
      {/* Basic Info Section */}
      <div className="flex flex-col gap-8">
        <div>
          <h3 className="text-lg font-medium mb-4">Basic Information</h3>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="name">Workspace Name</Label>
              <Input
                id="name"
                value={formData.name}
                onChange={(value) => setFormData({ ...formData, name: value })}
                disabled={!isEditing}
                className="h-10"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="billing-email">Billing Email</Label>
              <Input
                id="billing-email"
                type="email"
                value={formData.billingEmail}
                onChange={(value) =>
                  setFormData({ ...formData, billingEmail: value })
                }
                disabled={!isEditing}
                placeholder="billing@example.com"
                className="h-10"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="website">Website</Label>
              <Input
                id="website"
                value={formData.website}
                onChange={(value) =>
                  setFormData({ ...formData, website: value })
                }
                disabled={!isEditing}
                placeholder="https://example.com"
                className="h-10"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="about">About</Label>
              <Textarea
                id="about"
                value={formData.about}
                onChange={(e) =>
                  setFormData({ ...formData, about: e.target.value })
                }
                disabled={!isEditing}
                placeholder="Tell us about this workspace"
                className="min-h-[100px]"
              />
            </div>
          </div>
        </div>

        {/* Workspace Status Section */}
        <div>
          <h3 className="text-lg font-medium mb-4">Status</h3>
          <div className="flex gap-3">
            {workspace.deletedAt ? (
              <Badge variant="destructive">Deleted</Badge>
            ) : (
              <Badge variant="success">Active</Badge>
            )}
          </div>
        </div>

        {/* Metadata Section */}
        <div>
          <h3 className="text-lg font-medium mb-4">Metadata</h3>
          <div className="grid grid-cols-1 gap-4">
            <div className="space-y-2">
              <Label className="text-sm font-medium text-foreground">
                Created At
              </Label>
              <p className="text-sm text-muted-foreground bg-accent px-3 py-2 rounded-md">
                {formatDate(workspace.createdAt)}
              </p>
            </div>
            <div className="space-y-2">
              <Label className="text-sm font-medium text-foreground">
                Updated At
              </Label>
              <p className="text-sm text-muted-foreground bg-accent px-3 py-2 rounded-md">
                {formatDate(workspace.updatedAt)}
              </p>
            </div>
            <div className="space-y-2">
              <Label className="text-sm font-medium text-foreground">
                Created By
              </Label>
              <p className="text-sm text-muted-foreground bg-accent px-3 py-2 rounded-md">
                {workspace.createdByName || "Unknown"}
              </p>
            </div>
            <div className="space-y-2">
              <Label className="text-sm font-medium text-foreground">
                Updated By
              </Label>
              <p className="text-sm text-muted-foreground bg-accent px-3 py-2 rounded-md">
                {workspace.updatedByName || "Unknown"}
              </p>
            </div>
            <div className="space-y-2">
              <Label className="text-sm font-medium text-foreground">
                User Count
              </Label>
              <p className="text-sm text-muted-foreground bg-accent px-3 py-2 rounded-md">
                {workspace.usersCount || 0} users
              </p>
            </div>
            {workspace.deletedAt && (
              <>
                <div className="space-y-2">
                  <Label className="text-sm font-medium text-foreground">
                    Deleted At
                  </Label>
                  <p className="text-sm text-muted-foreground bg-accent px-3 py-2 rounded-md">
                    {formatDate(workspace.deletedAt)}
                  </p>
                </div>
                <div className="space-y-2">
                  <Label className="text-sm font-medium text-foreground">
                    Deleted By
                  </Label>
                  <p className="text-sm text-muted-foreground bg-accent px-3 py-2 rounded-md">
                    {workspace.deletedByName || "Unknown"}
                  </p>
                </div>
              </>
            )}
            <div className="space-y-2">
              <Label className="text-sm font-medium text-foreground">
                Workspace ID
              </Label>
              <div className="flex items-center gap-2 bg-accent px-3 py-2 rounded-md">
                <span className="text-sm text-muted-foreground font-mono flex-1">
                  {workspace.id}
                </span>
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-6 w-6 p-0"
                  onClick={() => copyToClipboard(workspace.id, "Workspace ID")}
                >
                  <Icon type="content_copy" className="text-sm" />
                </Button>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Actions Section */}
      <div className="flex flex-wrap items-center justify-between gap-3 border-t border-border pt-6">
        {workspace.deletedAt ? (
          <Button
            variant="outline"
            size="sm"
            onClick={() => handleRecover(workspace)}
            disabled={isBusy}
          >
            {pendingAction === "recover" ? <Spinner /> : "Recover workspace"}
          </Button>
        ) : (
          <Button
            variant="destructive"
            size="sm"
            onClick={() => handleDelete(workspace)}
            disabled={isBusy}
          >
            {pendingAction === "delete" ? <Spinner /> : "Delete workspace"}
          </Button>
        )}

        <div className="flex gap-2">
          {isEditing ? (
            <>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setIsEditing(false)}
                disabled={isBusy}
              >
                Cancel
              </Button>
              <Button size="sm" onClick={handleSave} disabled={isBusy}>
                {pendingAction === "save" ? <Spinner /> : "Save"}
              </Button>
            </>
          ) : (
            <Button size="sm" onClick={() => setIsEditing(true)}>
              Edit
            </Button>
          )}
        </div>
      </div>
    </div>
  );
}
