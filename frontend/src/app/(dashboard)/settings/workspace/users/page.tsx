"use client";

import { useEffect, useState, useCallback } from "react";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { SearchBar } from "@/components/ui/search-bar";
import { SectionMark } from "@/components/ui/section-mark";
import { UsersList } from "@/components/workspace/users-list";
import { Notifier } from "@/components/ui/sonner";
import { hasPermission } from "@/store/modules/auth";
import {
  fetchWorkspaceUsers,
  updateWorkspaceUserRole,
  removeWorkspaceUser,
  resetWorkspaceUsers,
} from "@/store/modules/workspaceUsers";

type WorkspaceUserRole = "owner" | "admin" | "member";

export default function WorkspaceUsersPage() {
  const dispatch = useAppDispatch();
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );
  const { currentUser } = useAppSelector((state) => state.auth);
  const {
    users,
    status: workspaceUsersStatus,
    pageInfo,
  } = useAppSelector((state) => state.workspaceUsers);

  const [currentPage, setCurrentPage] = useState(1);
  const [searchQuery, setSearchQuery] = useState("");

  const perPage = 10;

  // Check if current user has permission to manage users
  const currentUserRole = currentUser?.workspaceRole;
  const canManageUsers = useAppSelector((state) =>
    currentUserRole
      ? hasPermission(state, currentUserRole, "workspace:users:manage")
      : false,
  );

  const loadUsers = useCallback(
    async (pageNumber: number, search?: string, append = false) => {
      if (!currentWorkspace?.id) return;
      dispatch(
        fetchWorkspaceUsers({
          workspaceId: currentWorkspace.id,
          params: {
            page: pageNumber,
            perPage,
            query: search,
          },
          append,
        }),
      );
    },
    [currentWorkspace?.id, dispatch, perPage],
  );

  useEffect(() => {
    if (currentWorkspace?.id) {
      loadUsers(1);
    }

    return () => {
      dispatch(resetWorkspaceUsers());
    };
  }, [currentWorkspace?.id, dispatch, loadUsers]);

  const handleSearch = useCallback(
    (query: string) => {
      setSearchQuery(query);
      setCurrentPage(1);
      loadUsers(1, query);
    },
    [loadUsers],
  );

  const handleLoadMore = useCallback(() => {
    const nextPage = currentPage + 1;
    setCurrentPage(nextPage);
    loadUsers(nextPage, searchQuery, true);
  }, [loadUsers, searchQuery, currentPage]);

  const handleUpdateRole = async (
    userId: string,
    newRole: WorkspaceUserRole,
  ) => {
    if (!currentWorkspace?.id || !canManageUsers) return;

    try {
      await dispatch(
        updateWorkspaceUserRole({
          workspaceId: currentWorkspace.id,
          userId,
          role: newRole,
        }),
      ).unwrap();

      Notifier.success({
        title: "Role updated",
      });
    } catch (error) {
      Notifier.error({
        title: "Couldn't update role",
        message:
          typeof error === "string" ? error : "Failed to update user role",
      });
    }
  };

  const handleRemoveUser = async (userId: string) => {
    if (!currentWorkspace?.id || !canManageUsers) return;

    try {
      await dispatch(
        removeWorkspaceUser({
          workspaceId: currentWorkspace.id,
          userId,
        }),
      ).unwrap();

      Notifier.success({
        title: "User removed",
      });
    } catch (error) {
      Notifier.error({
        title: "Couldn't remove user",
        message: typeof error === "string" ? error : "Failed to remove user",
      });
    }
  };

  return (
    <section className="flex flex-col gap-5">
      <SectionMark label="Members" />
      <div>
        <div className="flex flex-col gap-5">
          <SearchBar placeholder="Search users..." onSearch={handleSearch} />
          <UsersList
            users={users || []}
            loading={workspaceUsersStatus === "loading"}
            pageInfo={
              pageInfo ? { hasNextPage: !!pageInfo.hasNextPage } : undefined
            }
            onLoadMore={handleLoadMore}
            onUpdateRole={handleUpdateRole}
            onRemoveUser={handleRemoveUser}
            workspaceId={currentWorkspace?.id || ""}
          />
        </div>
      </div>
    </section>
  );
}
