"use client";

import { useState, useCallback, useEffect, useMemo } from "react";
import { WorkspaceUser } from "@/types/user";
import WorkspaceUsersList from "./WorkspaceUsersList";
import api from "@/api";
import useAPI from "@/api/hooks";
import { Notifier } from "@/components/ui/sonner";
import { useAppSelector, useAppDispatch } from "@/store/hooks";
import { fetchRoles } from "@/store/modules/auth";
import { useDebounce } from "@/hooks/use-debounce";
import { copyToClipboard } from "@/lib/clipboard";

interface WorkspaceUsersTabProps {
  workspaceId: string | null;
}

export default function WorkspaceUsersTab({
  workspaceId,
}: WorkspaceUsersTabProps) {
  const [usersResponse, getWorkspaceUsers] = useAPI(
    api.admin.getWorkspaceUsers,
  );

  // Get roles from store and dispatch
  const dispatch = useAppDispatch();
  const roles = useAppSelector((state) => state.auth.roles);
  const rolesStatus = useAppSelector((state) => state.auth.rolesStatus);

  const [users, setUsers] = useState<WorkspaceUser[]>([]);
  const [usersSearchQuery, setUsersSearchQuery] = useState("");

  // Debounce the search query to avoid excessive API calls
  const debouncedSearchQuery = useDebounce(usersSearchQuery, 300);

  const loadWorkspaceUsers = useCallback(async () => {
    if (!workspaceId) return;
    await getWorkspaceUsers(workspaceId, {
      page: 1,
      perPage: 50,
      query: debouncedSearchQuery,
    });
  }, [getWorkspaceUsers, workspaceId, debouncedSearchQuery]);

  // Handle role update
  const handleRoleUpdate = useCallback(
    async (userId: string, newRole: string) => {
      if (!workspaceId) return;

      try {
        await api.admin.updateWorkspaceUserRole(workspaceId, userId, {
          role: newRole,
        });
        Notifier.success("User role updated successfully");
        // Reload users to get updated data
        await loadWorkspaceUsers();
      } catch (error) {
        console.error("Failed to update user role:", error);
        Notifier.error("Failed to update user role");
      }
    },
    [workspaceId, loadWorkspaceUsers],
  );

  // Memoized filtered users to prevent re-rendering
  const filteredUsers = useMemo(() => {
    if (!usersSearchQuery.trim()) return users;

    const query = usersSearchQuery.toLowerCase();
    return users.filter(
      (user) =>
        user.name.toLowerCase().includes(query) ||
        user.email.toLowerCase().includes(query),
    );
  }, [users, usersSearchQuery]);

  useEffect(() => {
    if (workspaceId) {
      loadWorkspaceUsers();

      // Load roles if not already loaded
      if (rolesStatus === "idle" || !roles) {
        dispatch(fetchRoles(workspaceId));
      }
    }
  }, [workspaceId, loadWorkspaceUsers, dispatch, rolesStatus, roles]);

  // Trigger search when debounced query changes
  useEffect(() => {
    if (workspaceId) {
      loadWorkspaceUsers();
    }
  }, [debouncedSearchQuery, workspaceId, loadWorkspaceUsers]);

  useEffect(() => {
    if (usersResponse.result?.results) {
      setUsers(usersResponse.result.results);
    }
  }, [usersResponse.result]);

  return (
    <WorkspaceUsersList
      searchQuery={usersSearchQuery}
      onSearchChange={setUsersSearchQuery}
      users={filteredUsers}
      loading={usersResponse.isFetching}
      onRoleUpdate={handleRoleUpdate}
      onRefresh={loadWorkspaceUsers}
      copyToClipboard={copyToClipboard}
    />
  );
}
