"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Checkbox } from "@/components/ui/checkbox";
import { Switch } from "@/components/ui/switch";
import { DataTable, DataTableColumn } from "@/components/ui/data-table";
import { SearchBar } from "@/components/ui/search-bar";
import Timestamp from "@/components/ui/timestamp";
import { AdminWorkspace } from "@/types/admin";
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

interface AdminWorkspacesTableProps {
  workspaces: AdminWorkspace[];
  loading: boolean;
  onSearch: (query: string) => void;
  onViewDetails: (workspaceId: string) => void;
  onDelete: (workspaceIds: string[]) => void;
  onRecover?: (workspaceIds: string[]) => void;
  includeDeleted?: boolean;
  onIncludeDeletedChange?: (includeDeleted: boolean) => void;
  pagination?: {
    page: number;
    perPage: number;
    total: number;
    onPageChange: (page: number) => void;
  };
}

export default function AdminWorkspacesTable({
  workspaces,
  loading,
  onSearch,
  onViewDetails,
  onDelete,
  onRecover,
  includeDeleted = false,
  onIncludeDeletedChange,
  pagination,
}: AdminWorkspacesTableProps) {
  const [selectedWorkspaces, setSelectedWorkspaces] = useState<string[]>([]);

  const handleSelectAll = (checked: boolean) => {
    setSelectedWorkspaces(
      checked ? workspaces.map((workspace) => workspace.id) : [],
    );
  };

  const handleSelectWorkspace = (workspaceId: string, checked: boolean) => {
    setSelectedWorkspaces((prev) =>
      checked
        ? [...prev, workspaceId]
        : prev.filter((id) => id !== workspaceId),
    );
  };

  const handleBulkDelete = () => {
    if (selectedWorkspaces.length > 0) {
      askConfirmation({
        title: "Delete Workspaces",
        content: `Are you sure you want to delete ${selectedWorkspaces.length} workspace(s)? This action cannot be undone.`,
        variant: "destructive",
        okText: "Delete",
        cancelText: "Cancel",
        okCallback: () => {
          onDelete(selectedWorkspaces);
          setSelectedWorkspaces([]);
        },
      });
    }
  };

  const handleBulkRecover = () => {
    if (selectedWorkspaces.length > 0 && onRecover) {
      onRecover(selectedWorkspaces);
      setSelectedWorkspaces([]);
    }
  };

  const columns: DataTableColumn<AdminWorkspace>[] = [
    {
      key: "select",
      label: "",
      order: 0,
      className: "w-12",
      headerRender: () => (
        <Checkbox
          checked={
            selectedWorkspaces.length === workspaces.length &&
            workspaces.length > 0
          }
          onCheckedChange={(checked) => handleSelectAll(checked as boolean)}
        />
      ),
      render: (_, workspace) => (
        <Checkbox
          checked={selectedWorkspaces.includes(workspace.id)}
          onCheckedChange={(checked) =>
            handleSelectWorkspace(workspace.id, checked as boolean)
          }
        />
      ),
    },
    {
      key: "workspace",
      label: "Workspace",
      order: 1,
      className: "min-w-[220px]",
      render: (_, workspace) => (
        <div className="flex min-w-0 flex-col">
          <span className="flex items-center gap-2 font-medium">
            <span className="truncate">{workspace.name}</span>
            <button
              type="button"
              className="flex h-7 w-7 shrink-0 items-center justify-center rounded text-muted-foreground hover:bg-accent hover:text-foreground"
              onClick={() => copyToClipboard(workspace.id, "Workspace ID")}
              aria-label="Copy workspace ID"
            >
              <Icon type="content_copy" className="text-sm" />
            </button>
          </span>
          {workspace.website && (
            <a
              href={workspace.website}
              target="_blank"
              rel="noopener noreferrer"
              className="truncate text-sm text-primary hover:underline"
            >
              {workspace.website}
            </a>
          )}
          {workspace.billingEmail && (
            <span className="truncate text-sm text-muted-foreground">
              {workspace.billingEmail}
            </span>
          )}
        </div>
      ),
    },
    {
      key: "usersCount",
      label: "Users",
      order: 2,
      className: "min-w-[80px]",
      render: (_, workspace) => (
        <Badge variant="outline" className="font-mono tabular-nums">
          {workspace.usersCount ?? 0}
        </Badge>
      ),
    },
    {
      key: "status",
      label: "Status",
      order: 3,
      className: "min-w-[100px]",
      render: (_, workspace) => (
        <div className="flex flex-col gap-1">
          <Badge
            variant={workspace.deletedAt ? "destructive" : "success"}
            size="sm"
          >
            {workspace.deletedAt ? "Deleted" : "Active"}
          </Badge>
          {workspace.deletedByName && (
            <span className="text-xs text-muted-foreground">
              by {workspace.deletedByName}
            </span>
          )}
        </div>
      ),
    },
    {
      key: "created",
      label: "Created",
      order: 5,
      className: "min-w-[120px]",
      render: (_, workspace) =>
        workspace.createdAt ? (
          <Timestamp
            timestamp={workspace.createdAt}
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
      order: 6,
      className: "w-[100px]",
      render: (_, workspace) => (
        <div className="flex items-center justify-end gap-1">
          <Button
            variant="outline"
            size="sm"
            onClick={() => onViewDetails(workspace.id)}
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
              {workspace.deletedAt ? (
                onRecover && (
                  <DropdownMenuItem onClick={() => onRecover([workspace.id])}>
                    <Icon type="restore" className="text-base" />
                    Recover
                  </DropdownMenuItem>
                )
              ) : (
                <DropdownMenuItem
                  className="text-destructive focus:text-destructive"
                  onClick={() => {
                    askConfirmation({
                      title: "Delete Workspace",
                      content: `Are you sure you want to delete workspace "${workspace.name}"? This action cannot be undone.`,
                      variant: "destructive",
                      okText: "Delete",
                      cancelText: "Cancel",
                      okCallback: () => onDelete([workspace.id]),
                    });
                  }}
                >
                  <Icon type="delete" className="text-base" />
                  Delete
                </DropdownMenuItem>
              )}
              <DropdownMenuSeparator />
              <DropdownMenuItem
                onClick={() => copyToClipboard(workspace.id, "Workspace ID")}
              >
                <Icon type="content_copy" className="text-base" />
                Copy ID
              </DropdownMenuItem>
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
            placeholder="Search workspaces by name…"
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

        {selectedWorkspaces.length > 0 && (
          <div className="flex flex-wrap gap-2">
            {onRecover && (
              <Button variant="outline" onClick={handleBulkRecover} size="sm">
                Recover ({selectedWorkspaces.length})
              </Button>
            )}
            <Button variant="destructive" onClick={handleBulkDelete} size="sm">
              Delete ({selectedWorkspaces.length})
            </Button>
          </div>
        )}
      </div>

      <DataTable
        data={workspaces}
        columns={columns}
        loading={loading}
        emptyMessage="No workspaces found"
        rowClassName={(workspace) =>
          workspace.deletedAt ? "bg-destructive/5 opacity-80" : ""
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
