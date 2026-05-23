"use client";

import { useState } from "react";
import { Input } from "../ui/input";
import { Button } from "../ui/button";
import { Badge } from "../ui/badge";
import { Checkbox } from "../ui/checkbox";
import { Switch } from "../ui/switch";
import { DataTable, DataTableColumn } from "../ui/data-table";
import { User } from "@/types/user";
import { UserAvatar } from "../ui/user-avatar";
import Icon from "../ui/icon";
import { askConfirmation } from "../ui/confirmation";
import { copyToClipboard } from "@/lib/clipboard";

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
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedUsers, setSelectedUsers] = useState<string[]>([]);

  const handleSearch = (value: string) => {
    setSearchQuery(value);
    onSearch(value);
  };

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

  const formatDateTime = (dateString: string) => {
    try {
      return new Date(dateString).toLocaleString();
    } catch {
      return "Invalid date";
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
      className: "min-w-[200px]",
      render: (_, user) => (
        <div className="flex items-center gap-3">
          <UserAvatar
            name={user.name}
            imageUrl={user.avatar}
            className="h-8 w-8"
          />
          <div className="flex flex-col">
            <span className="font-medium flex items-center gap-2">
              {user.name}
              {user.godMode ? (
                <Badge variant="outline" className="font-mono">
                  <Icon type="star" className="text-primary" /> Admin
                </Badge>
              ) : (
                ""
              )}{" "}
              <div
                className="p-1 h-8 w-8 flex items-center justify-center rounded hover:bg-accent text-xs text-muted-foreground hover:text-foreground font-mono"
                onClick={() => copyToClipboard(user.id, "User ID")}
              >
                <Icon type="content_copy" />
              </div>
            </span>
            <span className="text-sm text-muted-foreground">{user.email}</span>
          </div>
        </div>
      ),
    },
    {
      key: "workspace",
      label: "Default Workspace",
      order: 3,
      className: "min-w-[180px]",
      render: (_, user) => (
        <div className="flex flex-col gap-1">
          {user.defaultWorkspaceName ? (
            <>
              <span className="text-sm">{user.defaultWorkspaceName}</span>
              {user.defaultWorkspaceId && (
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-auto p-0 text-xs text-muted-foreground hover:text-foreground justify-start"
                  onClick={() =>
                    copyToClipboard(user.defaultWorkspaceId!, "Workspace ID")
                  }
                >
                  <Icon type="content_copy" className="h-3 w-3 mr-1" />
                  {user.defaultWorkspaceId.slice(0, 8)}...
                </Button>
              )}
            </>
          ) : (
            <span className="text-muted-foreground text-sm">No workspace</span>
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
        <div className="flex items-center justify-center">
          <Badge variant="outline" className="font-mono">
            {user.workspacesCount ?? 0}
          </Badge>
        </div>
      ),
    },
    {
      key: "status",
      label: "Status",
      order: 5,
      className: "min-w-[120px]",
      render: (_, user) => (
        <div className="flex flex-col gap-1">
          <Badge variant={user.verifiedAt ? "default" : "destructive"}>
            {user.verifiedAt ? "Verified" : "Unverified"}
          </Badge>
          {user.godMode && <Badge variant="warning">God Mode</Badge>}
          {user.deletedAt && <Badge variant="destructive">Deleted</Badge>}
        </div>
      ),
    },
    {
      key: "timestamps",
      label: "Timestamps",
      order: 6,
      className: "min-w-[200px]",
      render: (_, user) => (
        <div className="flex flex-col gap-1 text-xs">
          <div>
            <span className="text-muted-foreground">Created:</span>{" "}
            <span>{formatDateTime(user.createdAt)}</span>
          </div>
          <div>
            <span className="text-muted-foreground">Updated:</span>{" "}
            <span>{formatDateTime(user.updatedAt)}</span>
          </div>
          {user.verifiedAt && (
            <div>
              <span className="text-muted-foreground">Verified:</span>{" "}
              <span>{formatDateTime(user.verifiedAt)}</span>
            </div>
          )}
        </div>
      ),
    },
    {
      key: "deletion",
      label: "Deletion Info",
      order: 7,
      className: "min-w-[150px]",
      render: (_, user) => (
        <div className="flex flex-col gap-1 text-xs">
          {user.deletedAt ? (
            <>
              <div>
                <span className="text-muted-foreground">Deleted:</span>{" "}
                <span>{formatDateTime(user.deletedAt)}</span>
              </div>
              {user.deletedByName && (
                <div>
                  <span className="text-muted-foreground">By:</span>{" "}
                  <span>{user.deletedByName}</span>
                </div>
              )}
            </>
          ) : (
            <span className="text-muted-foreground">Not deleted</span>
          )}
        </div>
      ),
    },
    {
      key: "actions",
      label: "Actions",
      order: 8,
      className: "min-w-[240px]",
      render: (_, user) => (
        <div className="flex gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => onViewDetails(user.id)}
          >
            View
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={() => onImpersonate(user.id)}
            disabled={!!user.deletedAt}
          >
            Impersonate
          </Button>
          {!user.verifiedAt && !user.deletedAt && (
            <Button
              variant="outline"
              size="sm"
              onClick={() => onVerify([user.id])}
            >
              Verify
            </Button>
          )}
        </div>
      ),
    },
  ];

  return (
    <div className="space-y-4">
      {/* Search and Actions */}
      <div className="flex justify-between items-center">
        <div className="flex items-center gap-4">
          <div className="w-64">
            <Input
              id="search-users"
              placeholder="Search users by name or email..."
              value={searchQuery}
              onChange={handleSearch}
            />
          </div>
          <label className="flex items-center gap-2 text-sm text-muted-foreground cursor-pointer select-none">
            <Switch
              checked={includeDeleted}
              onCheckedChange={onIncludeDeletedChange}
            />
            Include deleted
          </label>
        </div>

        {selectedUsers.length > 0 && (
          <div className="flex gap-2">
            <Button variant="outline" onClick={handleBulkVerify} size="sm">
              Verify Selected ({selectedUsers.length})
            </Button>
            <Button variant="destructive" onClick={handleBulkDelete} size="sm">
              Delete Selected ({selectedUsers.length})
            </Button>
          </div>
        )}
      </div>

      {/* Data Table */}
      <DataTable
        columns={columns}
        data={users}
        loading={loading}
        emptyMessage="No users found"
      />

      {/* Pagination */}
      {pagination && (
        <div className="flex justify-between items-center">
          <p className="text-sm text-muted-foreground">
            Showing{" "}
            {Math.min(
              (pagination.page - 1) * pagination.perPage + 1,
              pagination.total,
            )}{" "}
            to{" "}
            {Math.min(pagination.page * pagination.perPage, pagination.total)}{" "}
            of {pagination.total} users
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
