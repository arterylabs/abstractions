"use client";

import { useState, useCallback, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import api from "@/api";
import useAPI from "@/api/hooks";
import Spinner from "@/components/ui/spinner";
import { Notifier } from "@/components/ui/sonner";
import { format } from "date-fns";
import Icon from "@/components/ui/icon";
import { UserAvatar } from "@/components/ui/user-avatar";
import { copyToClipboard } from "@/lib/clipboard";
import { askConfirmation } from "@/components/ui/confirmation";

interface UserDetailsTabProps {
  userId: string | null;
  onUserUpdated: () => void;
}

export default function UserDetailsTab({
  userId,
  onUserUpdated,
}: UserDetailsTabProps) {
  const [userResponse, getUser] = useAPI(api.admin.getUser);
  const [updateResponse, updateUser] = useAPI(api.admin.updateUser);
  const [verifyResponse, verifyUsers] = useAPI(api.admin.verifyUsers);
  const [impersonateResponse, impersonateUser] = useAPI(
    api.admin.impersonateUser,
  );
  const [deleteResponse, deleteUsers] = useAPI(api.admin.deleteUsers);
  const [recoverResponse, recoverUsers] = useAPI(api.admin.recoverUsers);

  const [formData, setFormData] = useState({
    name: "",
    email: "",
  });
  const [isEditing, setIsEditing] = useState(false);

  const loadUser = useCallback(async () => {
    if (userId) {
      await getUser(userId);
    }
  }, [getUser, userId]);

  useEffect(() => {
    loadUser();
  }, [loadUser]);

  useEffect(() => {
    if (userResponse.result) {
      setFormData({
        name: userResponse.result.name || "",
        email: userResponse.result.email || "",
      });
    }
  }, [userResponse.result]);

  const handleSave = async () => {
    if (!userId) return;

    try {
      await updateUser(userId, formData);
    } catch {
      Notifier.error("Failed to update user");
    }
  };

  const handleVerifyUser = async () => {
    if (!userId) return;

    try {
      await verifyUsers({ userIds: [userId] });
    } catch {
      Notifier.error("Failed to verify user");
    }
  };

  const handleImpersonateUser = async () => {
    if (!userId) return;

    try {
      await impersonateUser({ id: userId });
    } catch {
      Notifier.error("Failed to impersonate user");
    }
  };

  const handleRecoverUser = async () => {
    if (!userId) return;

    askConfirmation({
      title: "Recover User",
      content: `Are you sure you want to recover "${user?.name}"? This will restore the user account.`,
      variant: "primary",
      okText: "Recover User",
      cancelText: "Cancel",
      okCallback: async () => {
        try {
          await recoverUsers({ userIds: [userId] });
        } catch {
          Notifier.error("Failed to recover user");
        }
      },
    });
  };

  const handleDeleteUser = async () => {
    if (!userId) return;

    askConfirmation({
      title: "Delete User",
      content: `Are you sure you want to delete "${user?.name}"? This action cannot be undone and will permanently remove the user from the system.`,
      variant: "destructive",
      okText: "Delete User",
      cancelText: "Cancel",
      okCallback: async () => {
        try {
          await deleteUsers({ userIds: [userId] });
        } catch {
          Notifier.error("Failed to delete user");
        }
      },
    });
  };

  // Handle update response
  useEffect(() => {
    if (updateResponse.result) {
      Notifier.success("User updated successfully");
      setIsEditing(false);
      onUserUpdated();
    }

    if (updateResponse.error) {
      Notifier.error("Failed to update user");
    }
  }, [updateResponse.result, updateResponse.error, onUserUpdated]);

  // Handle verify response
  useEffect(() => {
    if (verifyResponse.result) {
      Notifier.success("User verified successfully");
      loadUser(); // Refresh user data
      onUserUpdated();
    }

    if (verifyResponse.error) {
      Notifier.error("Failed to verify user");
    }
  }, [verifyResponse.result, verifyResponse.error, loadUser, onUserUpdated]);

  // Handle impersonate response
  useEffect(() => {
    if (impersonateResponse.result?.token) {
      // Store the current god mode user's token before impersonating
      const currentToken = api.utils.getToken();
      if (currentToken) {
        localStorage.setItem("impersonatorToken", currentToken);
      }

      // Set the impersonation token and redirect
      api.utils.setToken(
        impersonateResponse.result.token,
        impersonateResponse.result.expiresAt,
      );
      Notifier.success("Impersonating user successfully");
      window.location.href = "/";
    }

    if (impersonateResponse.error) {
      Notifier.error("Failed to impersonate user");
    }
  }, [impersonateResponse.result, impersonateResponse.error]);

  // Handle recover response
  useEffect(() => {
    if (recoverResponse.result) {
      Notifier.success("User recovered successfully");
      loadUser(); // Refresh user data
      onUserUpdated();
    }

    if (recoverResponse.error) {
      Notifier.error("Failed to recover user");
    }
  }, [recoverResponse.result, recoverResponse.error, loadUser, onUserUpdated]);

  // Handle delete response
  useEffect(() => {
    if (deleteResponse.result) {
      Notifier.success("User deleted successfully");
      onUserUpdated();
    }

    if (deleteResponse.error) {
      Notifier.error("Failed to delete user");
    }
  }, [deleteResponse.result, deleteResponse.error, onUserUpdated]);

  const formatDate = (dateString: string | undefined) => {
    if (!dateString) return "N/A";
    try {
      return format(new Date(dateString), "PPpp"); // e.g., "Apr 29, 2023 at 3:45 PM"
    } catch {
      return "Invalid date";
    }
  };

  const user = userResponse.result;

  if (userResponse.isFetching) {
    return (
      <div className="flex items-center justify-center py-12">
        <Spinner />
      </div>
    );
  }

  if (!user) {
    return (
      <div className="text-center py-12">
        <p className="text-muted-foreground">User not found</p>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      {/* User Profile Section */}
      <div className="flex flex-col gap-8">
        <div>
          <h3 className="text-lg font-medium mb-4">Profile Information</h3>

          {/* Avatar and Name Section */}
          <div className="flex items-center gap-4 mb-6">
            <UserAvatar
              name={user.name}
              imageUrl={user.avatar}
              className="h-16 w-16"
            />
            <div>
              <h4 className="text-xl font-medium">{user.name}</h4>
              <p className="text-sm text-muted-foreground">{user.email}</p>
            </div>
          </div>

          <div className="grid grid-cols-1 gap-6">
            <div className="space-y-2">
              <Label htmlFor="name">Full Name</Label>
              <Input
                id="name"
                value={formData.name}
                onChange={(value) => setFormData({ ...formData, name: value })}
                disabled={!isEditing}
                className="w-full"
                inputClassName="h-10"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="email">Email Address</Label>
              <Input
                id="email"
                type="email"
                value={formData.email}
                onChange={(value) => setFormData({ ...formData, email: value })}
                disabled={!isEditing}
                className="w-full"
                inputClassName="h-10"
              />
            </div>
          </div>
        </div>

        {/* User Status Section */}
        <div>
          <h3 className="text-lg font-medium mb-4">Account Status</h3>
          <div className="flex gap-3">
            {user.deletedAt && (
              <Badge variant="destructive" className="text-sm px-3 py-1">
                Deleted
              </Badge>
            )}
            <Badge
              variant={
                user.verifiedAt && user.verifiedAt !== ""
                  ? "default"
                  : "destructive"
              }
              className="text-sm px-3 py-1"
            >
              {user.verifiedAt && user.verifiedAt !== ""
                ? "Verified"
                : "Unverified"}
            </Badge>
            {user.godMode && (
              <Badge
                variant="outline"
                className="border-violet/40 bg-violet/10 text-violet-light"
              >
                Admin
              </Badge>
            )}
          </div>
        </div>

        {/* Account Metadata Section */}
        <div>
          <h3 className="text-lg font-medium mb-4">Account Information</h3>
          <div className="grid grid-cols-1 gap-6">
            <div className="space-y-2">
              <Label className="text-sm font-medium text-foreground">
                User ID
              </Label>
              <div className="flex items-center gap-2">
                <p className="text-sm text-muted-foreground bg-accent px-3 py-2 rounded-md font-mono flex-1">
                  {user.id}
                </p>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => copyToClipboard(user.id, "User ID")}
                  className="h-8 w-8 p-0"
                >
                  <Icon type="content_copy" className="text-sm" />
                </Button>
              </div>
            </div>
            <div className="space-y-2">
              <Label className="text-sm font-medium text-foreground">
                Created At
              </Label>
              <p className="text-sm text-muted-foreground bg-accent px-3 py-2 rounded-md">
                {formatDate(user.createdAt)}
              </p>
            </div>
            {user.updatedAt && (
              <div className="space-y-2">
                <Label className="text-sm font-medium text-foreground">
                  Updated At
                </Label>
                <p className="text-sm text-muted-foreground bg-accent px-3 py-2 rounded-md">
                  {formatDate(user.updatedAt)}
                </p>
              </div>
            )}
            {user.verifiedAt && (
              <div className="space-y-2">
                <Label className="text-sm font-medium text-foreground">
                  Verified At
                </Label>
                <p className="text-sm text-muted-foreground bg-accent px-3 py-2 rounded-md">
                  {formatDate(user.verifiedAt)}
                </p>
              </div>
            )}
            {user.deletedAt && (
              <>
                <div className="space-y-2">
                  <Label className="text-sm font-medium text-foreground">
                    Deleted At
                  </Label>
                  <p className="text-sm text-muted-foreground bg-accent px-3 py-2 rounded-md">
                    {formatDate(user.deletedAt)}
                  </p>
                </div>
                {user.deletedByName && (
                  <div className="space-y-2">
                    <Label className="text-sm font-medium text-foreground">
                      Deleted By
                    </Label>
                    <p className="text-sm text-muted-foreground bg-accent px-3 py-2 rounded-md">
                      {user.deletedByName}
                    </p>
                  </div>
                )}
              </>
            )}
            {user.defaultWorkspaceId && (
              <div className="space-y-2">
                <Label className="text-sm font-medium text-foreground">
                  Default Workspace ID
                </Label>
                <div className="flex items-center gap-2">
                  <p className="text-sm text-muted-foreground bg-accent px-3 py-2 rounded-md font-mono flex-1">
                    {user.defaultWorkspaceId}
                  </p>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() =>
                      copyToClipboard(
                        user.defaultWorkspaceId!,
                        "Default Workspace ID",
                      )
                    }
                    className="h-8 w-8 p-0"
                  >
                    <Icon type="content_copy" className="text-sm" />
                  </Button>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Actions Section */}
      <div className="space-y-4 border-t border-border pt-6">
        <div className="flex flex-wrap gap-2">
          {(!user.verifiedAt || user.verifiedAt === "") && (
            <Button
              onClick={handleVerifyUser}
              disabled={verifyResponse.isFetching}
              variant="outline"
              size="sm"
            >
              {verifyResponse.isFetching ? <Spinner /> : "Verify user"}
            </Button>
          )}
          <Button
            onClick={handleImpersonateUser}
            disabled={impersonateResponse.isFetching || !!user.deletedAt}
            variant="outline"
            size="sm"
          >
            {impersonateResponse.isFetching ? <Spinner /> : "Impersonate"}
          </Button>
          {user.deletedAt ? (
            <Button
              onClick={handleRecoverUser}
              disabled={recoverResponse.isFetching}
              variant="outline"
              size="sm"
            >
              {recoverResponse.isFetching ? <Spinner /> : "Recover user"}
            </Button>
          ) : (
            <Button
              onClick={handleDeleteUser}
              disabled={deleteResponse.isFetching}
              variant="destructive"
              size="sm"
            >
              {deleteResponse.isFetching ? <Spinner /> : "Delete user"}
            </Button>
          )}
        </div>

        <div className="flex justify-end gap-2">
          {isEditing ? (
            <>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setIsEditing(false)}
              >
                Cancel
              </Button>
              <Button
                size="sm"
                onClick={handleSave}
                disabled={updateResponse.isFetching}
              >
                {updateResponse.isFetching ? <Spinner /> : "Save"}
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
