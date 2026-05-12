"use client";

import React, { useEffect } from "react";
import { InfiniteScroll } from "@/components/ui/infinite-scroll";
import { Card, CardContent } from "@/components/ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import Spinner from "@/components/ui/spinner";
import { User } from "@/types/user";
import { useAppSelector, useAppDispatch } from "@/store/hooks";
import { fetchRoles, hasPermission } from "@/store/modules/auth";
import Icon from "../ui/icon";
import { askConfirmation } from "@/components/ui/confirmation";
import {
  canModifyUserRole,
  isValidRoleChange,
  getGrantableRoles,
} from "@/lib/role-hierarchy";
import { getInitials } from "@/lib/user";

type WorkspaceUserRole = "owner" | "admin" | "member";

interface WorkspaceUser extends User {
  role: WorkspaceUserRole;
}

interface UsersListProps {
  users: WorkspaceUser[];
  loading: boolean;
  pageInfo?: {
    hasNextPage: boolean;
  };
  onLoadMore: () => void;
  onUpdateRole: (userId: string, newRole: WorkspaceUserRole) => void;
  onRemoveUser: (userId: string) => void;
  workspaceId: string;
}

export function UsersList({
  users,
  loading,
  pageInfo,
  onLoadMore,
  onUpdateRole,
  onRemoveUser,
  workspaceId,
}: UsersListProps) {
  const dispatch = useAppDispatch();
  const roles = useAppSelector((state) => state.auth.roles);
  const rolesStatus = useAppSelector((state) => state.auth.rolesStatus);
  const currentUser = useAppSelector((state) => state.auth.currentUser);

  // Check if current user has permission to manage users
  const currentUserRole = currentUser?.workspaceRole;
  const canManageUsers = useAppSelector((state) =>
    currentUserRole
      ? hasPermission(state, currentUserRole, "workspace:users:manage")
      : false,
  );

  useEffect(() => {
    if (workspaceId && rolesStatus === "idle") {
      dispatch(fetchRoles(workspaceId));
    }
  }, [workspaceId, rolesStatus, dispatch]);

  const getRoleName = (roleId: string) => {
    const role = roles?.find((r) => r.id === roleId);
    return role?.name || roleId;
  };

  const handleRoleChange = (userId: string, newRole: WorkspaceUserRole) => {
    const user = users.find((u) => u.id === userId);
    const newRoleName = getRoleName(newRole);

    // Validate the role change
    if (currentUserRole) {
      const validation = isValidRoleChange(
        currentUserRole as WorkspaceUserRole,
        user?.role as WorkspaceUserRole,
        newRole,
      );

      if (!validation.valid) {
        askConfirmation({
          title: "Cannot Change Role",
          content:
            validation.reason ||
            "You don't have permission to perform this action",
          variant: "destructive",
          okText: "OK",
          okCallback: () => {},
        });
        return;
      }
    }

    askConfirmation({
      title: "Change User Role",
      content: `Are you sure you want to change ${user?.name}'s role to ${newRoleName}?`,
      variant: "primary",
      okText: "Change Role",
      okCallback: () => onUpdateRole(userId, newRole),
    });
  };

  const handleRemoveUser = (userId: string) => {
    const user = users.find((u) => u.id === userId);

    // Validate if current user can remove this user
    if (currentUserRole && user?.role) {
      if (
        !canModifyUserRole(
          currentUserRole as WorkspaceUserRole,
          user.role as WorkspaceUserRole,
        )
      ) {
        askConfirmation({
          title: "Cannot Remove User",
          content: "You don't have permission to remove this user",
          variant: "destructive",
          okText: "OK",
          okCallback: () => {},
        });
        return;
      }
    }

    askConfirmation({
      title: "Remove User",
      content: `Are you sure you want to remove ${user?.name} from this workspace? This action cannot be undone.`,
      variant: "destructive",
      okText: "Remove User",
      okCallback: () => onRemoveUser(userId),
    });
  };

  if (loading && users.length === 0) {
    return (
      <div className="flex h-[200px] w-full items-center justify-center">
        <Spinner className="h-8 w-8" />
      </div>
    );
  }

  if (!users?.length) {
    return (
      <div className="flex h-[200px] w-full items-center justify-center text-muted-foreground">
        No users found
      </div>
    );
  }

  return (
    <InfiniteScroll
      loading={loading}
      hasMore={!!pageInfo?.hasNextPage}
      onLoadMore={onLoadMore}
      className="flex flex-col gap-4"
    >
      {users.map((user) => (
        <Card key={user.id} className="flex flex-col">
          <CardContent className="flex items-start justify-between gap-4 p-6">
            <div className="flex items-start gap-3">
              <Avatar className="h-12 w-12">
                <AvatarImage src={user.avatar} alt={user.name} />
                <AvatarFallback>{getInitials(user.name)}</AvatarFallback>
              </Avatar>
              <div className="flex-1 min-w-0">
                <div className="font-medium truncate">{user.name}</div>
                <div className="text-sm text-muted-foreground truncate">
                  {user.email}
                </div>
                <div className="flex items-center text-sm text-muted-foreground mt-1">
                  <Icon type="assignment_ind" className="inline mr-1" />
                  {getRoleName(user.role)}
                </div>
              </div>
            </div>
            {currentUser?.id !== user.id ? (
              <div className="flex gap-2 mt-2">
                {/* Role Dropdown - shows current role and allows changes if user has permission */}
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button
                      variant="outline"
                      size="sm"
                      className="flex-1"
                      disabled={
                        !canManageUsers ||
                        !currentUserRole ||
                        !canModifyUserRole(
                          currentUserRole as WorkspaceUserRole,
                          user.role as WorkspaceUserRole,
                        )
                      }
                    >
                      {getRoleName(user.role)}
                      <Icon type="unfold_more" className="ml-auto text-2xl" />
                    </Button>
                  </DropdownMenuTrigger>
                  {canManageUsers && (
                    <DropdownMenuContent align="end">
                      {(() => {
                        // Get the roles that the current user can grant
                        const grantableRoles = currentUserRole
                          ? getGrantableRoles(
                              currentUserRole as WorkspaceUserRole,
                            )
                          : [];

                        // Filter roles to only show those the current user can grant
                        const availableRoles =
                          roles?.filter((role) =>
                            grantableRoles.includes(
                              role.id as WorkspaceUserRole,
                            ),
                          ) || [];

                        return availableRoles.map((role) => (
                          <DropdownMenuItem
                            key={role.id}
                            onClick={() =>
                              handleRoleChange(
                                user.id,
                                role.id as WorkspaceUserRole,
                              )
                            }
                            className="flex flex-col items-start"
                            disabled={role.id === user.role}
                          >
                            <span className="font-medium">{role.name}</span>
                            <span className="text-xs text-muted-foreground">
                              {role.description}
                            </span>
                          </DropdownMenuItem>
                        ));
                      })()}
                    </DropdownMenuContent>
                  )}
                </DropdownMenu>

                {/* Actions Dropdown - contains remove user and other actions */}
                {canManageUsers &&
                  currentUserRole &&
                  canModifyUserRole(
                    currentUserRole as WorkspaceUserRole,
                    user.role as WorkspaceUserRole,
                  ) && (
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="outline" size="sm" className="px-2">
                          <Icon type="more_vert" className="h-4 w-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem
                          onClick={() => handleRemoveUser(user.id)}
                          className="text-destructive focus:text-destructive"
                        >
                          <Icon type="person_remove" className="mr-2 h-4 w-4" />
                          Remove User
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  )}
              </div>
            ) : (
              <div className="text-sm text-muted-foreground mt-2">You</div>
            )}
          </CardContent>
        </Card>
      ))}
    </InfiniteScroll>
  );
}
