"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Checkbox } from "@/components/ui/checkbox";
import { Switch } from "@/components/ui/switch";
import { DataTable, DataTableColumn } from "@/components/ui/data-table";
import { SearchBar } from "@/components/ui/search-bar";
import Timestamp from "@/components/ui/timestamp";
import { User } from "@/types/user";
import { UserAvatar } from "@/components/ui/user-avatar";
import Icon from "@/components/ui/icon";
import { askConfirmation } from "@/components/ui/confirmation";
import { copyToClipboard } from "@/lib/clipboard";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

interface AdminUsersTableProps {
  users: User[];
  loading: boolean;
  onSearch: (query: string) => void;
  onViewDetails: (userId: string) => void;
  onDelete: (userIds: string[]) => void;
  onVerify: (userIds: string[]) => void;
  onImpersonate: (userId: string) => void;
  includeDeleted?: boolean;
  onIncludeDeletedChange?: (includeDeleted: boolean) => void;
  pagination?: {
    page: number;
    perPage: number;
    total: number;
    onPageChange: (page: number) => void;
  };
}

export default function AdminUsersTable({
  users,
  loading,
  onSearch,
  onViewDetails,
  onDelete,
  onVerify,
  onImpersonate,
  includeDeleted = false,
  onIncludeDeletedChange,
  pagination,
}: AdminUsersTableProps) {
  const [selectedUsers, setSelectedUsers] = useState<string[]>([]);

  const handleSelectAll = (checked: boolean) => {
    setSelectedUsers(checked ? users.map((user) => user.id) : []);
  };

  const handleSelectUser = (userId: string, checked: boolean) => {
    setSelectedUsers((prev) =>
      checked ? [...prev, userId] : prev.filter((id) => id !== userId),
    );
  };

  const handleBulkVerify = () => {
    if (selectedUsers.length > 0) {
      onVerify(selectedUsers);
      setSelectedUsers([]);
    }
  };

  const handleBulkDelete = () => {
    if (selectedUsers.length > 0) {
      askConfirmation({
        title: "Delete Users",
        content: `Are you sure you want to delete ${selectedUsers.length} user(s)? This action cannot be undone.`,
        variant: "destructive",
        okText: "Delete",
        cancelText: "Cancel",
        okCallback: () => {
          onDelete(selectedUsers);
          setSelectedUsers([]);
        },
      });
    }
  };

  const columns: DataTableColumn<User>[] = [
    {
      key: "select",
      label: "",
      order: 0,
      className: "w-12",
      headerRender: () => (
        <Checkbox
          checked={selectedUsers.length === users.length && users.length > 0}
          onCheckedChange={(checked) => handleSelectAll(checked as boolean)}
        />
      ),
      render: (_, user) => (
        <Checkbox
          checked={selectedUsers.includes(user.id)}
          onCheckedChange={(checked) =>
            handleSelectUser(user.id, checked as boolean)
          }
        />
      ),
    },
    {
      key: "user",
      label: "User",
      order: 1,
      className: "min-w-[220px]",
      render: (_, user) => (
        <div className="flex items-center gap-3">
          <UserAvatar
            name={user.name}
            imageUrl={user.avatar}
            className="h-8 w-8"
          />
          <div className="flex min-w-0 flex-col">
            <span className="flex items-center gap-2 font-medium">
              <span className="truncate">{user.name}</span>
              {user.godMode && (
                <Badge
                  variant="outline"
                  size="sm"
                  className="border-violet/40 bg-violet/10 text-violet-light"
                >
                  Admin
                </Badge>
              )}
              <button
                type="button"
                className="flex h-7 w-7 shrink-0 items-center justify-center rounded text-muted-foreground hover:bg-accent hover:text-foreground"
                onClick={() => copyToClipboard(user.id, "User ID")}
                aria-label="Copy user ID"
              >
                <Icon type="content_copy" className="text-sm" />
              </button>
            </span>
            <span className="truncate text-sm text-muted-foreground">
              {user.email}
            </span>
          </div>
        </div>
      ),
    },
    {
      key: "workspace",
      label: "Default workspace",
      order: 3,
      className: "min-w-[160px]",
      render: (_, user) => (
        <div className="flex flex-col gap-0.5">
          {user.defaultWorkspaceName ? (
            <>
              <span className="text-sm">{user.defaultWorkspaceName}</span>
              {user.defaultWorkspaceId && (
                <button
                  type="button"
                  className="flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground"
                  onClick={() =>
                    copyToClipboard(user.defaultWorkspaceId!, "Workspace ID")
                  }
                >
                  <Icon type="content_copy" className="text-xs" />
                  {user.defaultWorkspaceId.slice(0, 8)}…
                </button>
              )}
            </>
          ) : (
            <span className="text-sm text-muted-foreground">No workspace</span>
          )}
        </div>
      ),
    },
    {
      key: "workspacesCount",
      label: "Workspaces",
      order: 4,
      className: "min-w-[100px]",
      render: (_, user) => (
        <Badge variant="outline" className="font-mono tabular-nums">
          {user.workspacesCount ?? 0}
        </Badge>
      ),
    },
    {
      key: "status",
      label: "Status",
      order: 5,
      className: "min-w-[120px]",
      render: (_, user) => (
        <div className="flex flex-wrap gap-1">
          <Badge variant={user.verifiedAt ? "success" : "warning"} size="sm">
            {user.verifiedAt ? "Verified" : "Unverified"}
          </Badge>
          {user.deletedAt && (
            <Badge variant="destructive" size="sm">
              Deleted
            </Badge>
          )}
        </div>
      ),
    },
    {
      key: "created",
      label: "Created",
      order: 6,
      className: "min-w-[120px]",
      render: (_, user) =>
        user.createdAt ? (
          <Timestamp
            timestamp={user.createdAt}
            format="time-ago"
            className="text-xs text-muted-foreground"
          />
        ) : (
          <span className="text-xs text-muted-foreground">—</span>
        ),
    },
    {
      key: "actions",
      label: "",
      order: 8,
      className: "w-[100px]",
      render: (_, user) => (
        <div className="flex items-center justify-end gap-1">
          <Button
            variant="outline"
            size="sm"
            onClick={() => onViewDetails(user.id)}
          >
            View
          </Button>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button
                variant="ghost"
                size="sm"
                className="h-8 w-8 p-0"
                aria-label="More actions"
              >
                <Icon type="more_vert" className="text-base" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem
                onClick={() => onImpersonate(user.id)}
                disabled={!!user.deletedAt}
              >
                <Icon type="switch_account" className="text-base" />
                Impersonate
              </DropdownMenuItem>
              {!user.verifiedAt && !user.deletedAt && (
                <DropdownMenuItem onClick={() => onVerify([user.id])}>
                  <Icon type="verified" className="text-base" />
                  Verify
                </DropdownMenuItem>
              )}
              {!user.deletedAt && (
                <>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem
                    className="text-destructive focus:text-destructive"
                    onClick={() => {
                      askConfirmation({
                        title: "Delete User",
                        content: `Are you sure you want to delete "${user.name}"? This action cannot be undone.`,
                        variant: "destructive",
                        okText: "Delete",
                        cancelText: "Cancel",
                        okCallback: () => onDelete([user.id]),
                      });
                    }}
                  >
                    <Icon type="delete" className="text-base" />
                    Delete
                  </DropdownMenuItem>
                </>
              )}
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      ),
    },
  ];

  return (
    <div className="space-y-4">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex flex-1 flex-col gap-3 sm:flex-row sm:items-center">
          <SearchBar
            placeholder="Search by name or email…"
            onSearch={onSearch}
            className="max-w-sm"
          />
          <label className="flex cursor-pointer select-none items-center gap-2 text-sm text-muted-foreground">
            <Switch
              checked={includeDeleted}
              onCheckedChange={onIncludeDeletedChange}
            />
            Include deleted
          </label>
        </div>

        {selectedUsers.length > 0 && (
          <div className="flex flex-wrap gap-2">
            <Button variant="outline" onClick={handleBulkVerify} size="sm">
              Verify ({selectedUsers.length})
            </Button>
            <Button variant="destructive" onClick={handleBulkDelete} size="sm">
              Delete ({selectedUsers.length})
            </Button>
          </div>
        )}
      </div>

      <DataTable
        columns={columns}
        data={users}
        loading={loading}
        emptyMessage="No users found"
        rowClassName={(user) =>
          user.deletedAt ? "bg-destructive/5 opacity-80" : ""
        }
      />

      {pagination && pagination.total > 0 && (
        <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
          <p className="text-sm text-muted-foreground">
            Showing{" "}
            {Math.min(
              (pagination.page - 1) * pagination.perPage + 1,
              pagination.total,
            )}
            –{Math.min(pagination.page * pagination.perPage, pagination.total)}{" "}
            of {pagination.total}
          </p>
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => pagination.onPageChange(pagination.page - 1)}
              disabled={pagination.page === 1}
            >
              Previous
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={() => pagination.onPageChange(pagination.page + 1)}
              disabled={
                pagination.page * pagination.perPage >= pagination.total
              }
            >
              Next
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}
