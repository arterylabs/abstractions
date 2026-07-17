"use client";

import { useMemo } from "react";
import { WorkspaceUser } from "@/types/user";
import { SearchBar } from "@/components/ui/search-bar";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import Icon from "@/components/ui/icon";
import Spinner from "@/components/ui/spinner";
import { useAppSelector } from "@/store/hooks";

interface WorkspaceUsersListProps {
  searchQuery: string;
  onSearchChange: (query: string) => void;
  users: WorkspaceUser[];
  loading: boolean;
  onRoleUpdate: (userId: string, newRole: string) => void;
  onRefresh: () => void;
  copyToClipboard: (text: string, type: string) => void;
}

export default function WorkspaceUsersList({
  searchQuery,
  onSearchChange,
  users,
  loading,
  onRoleUpdate,
  onRefresh,
  copyToClipboard,
}: WorkspaceUsersListProps) {
  const roles = useAppSelector((state) => state.auth.roles);

  // Helper function to get role display name
  const getRoleName = (roleId: string) => {
    if (!roles || roles.length === 0) {
      return roleId; // Fallback to roleId if roles not loaded
    }
    const role = roles.find((r) => r.id === roleId);
    return role?.name || roleId;
  };

  // Memoized filtered users to prevent re-rendering
  const filteredUsers = useMemo(() => {
    if (!searchQuery.trim()) return users;

    const query = searchQuery.toLowerCase();
    return users.filter(
      (user) =>
        user.name.toLowerCase().includes(query) ||
        user.email.toLowerCase().includes(query),
    );
  }, [users, searchQuery]);

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <Button variant="outline" size="sm" onClick={onRefresh} className="h-9">
          Refresh
        </Button>
      </div>

      {/* Search Bar with debouncing handled in parent component */}
      <div className="space-y-4">
        <SearchBar
          placeholder="Search users by name or email..."
          value={searchQuery}
          onChange={onSearchChange}
          debounceMs={300}
          className="h-9"
        />
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-12">
          <Spinner />
        </div>
      ) : (
        <div className="space-y-3">
          {filteredUsers.map((user) => (
            <div
              key={user.id}
              className="p-4 border rounded-lg bg-accent/30 hover:bg-accent/50 transition-colors"
            >
              <div className="flex items-start justify-between">
                <div className="flex items-center space-x-4">
                  <div className="h-10 w-10 bg-muted rounded-full flex items-center justify-center">
                    <span className="text-sm font-medium text-muted-foreground">
                      {user.name.charAt(0).toUpperCase()}
                    </span>
                  </div>
                  <div className="flex-1">
                    <div className="flex items-center gap-2">
                      <p className="font-medium text-foreground">{user.name}</p>
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-6 w-6 p-0"
                        onClick={() => copyToClipboard(user.id, "User ID")}
                      >
                        <Icon type="content_copy" className="text-xs" />
                      </Button>
                    </div>
                    <p className="text-sm text-muted-foreground">
                      {user.email}
                    </p>
                    <div className="flex items-center gap-2 mt-2">
                      {/* Role Dropdown */}
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button
                            variant="outline"
                            size="sm"
                            className="h-7 px-3 text-xs font-medium"
                          >
                            {getRoleName(user.role)}
                            <Icon type="expand_more" className="ml-1 text-xs" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="start">
                          {roles?.map((role) => (
                            <DropdownMenuItem
                              key={role.id}
                              onClick={() => onRoleUpdate(user.id, role.id)}
                              className="flex flex-col items-start"
                              disabled={role.id === user.role}
                            >
                              <span className="font-medium">{role.name}</span>
                              <span className="text-xs text-muted-foreground">
                                {role.description}
                              </span>
                            </DropdownMenuItem>
                          ))}
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          ))}
          {filteredUsers.length === 0 && !loading && (
            <div className="text-center py-12">
              <p className="text-muted-foreground">
                {searchQuery.trim()
                  ? "No users found matching your search"
                  : "No users found in this workspace"}
              </p>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
