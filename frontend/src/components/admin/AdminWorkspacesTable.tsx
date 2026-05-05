"use client";

import { useState } from "react";
import { Input } from "../ui/input";
import { Button } from "../ui/button";
import { Badge } from "../ui/badge";
import { Checkbox } from "../ui/checkbox";
import { DataTable, DataTableColumn } from "../ui/data-table";
import { AdminWorkspace } from "@/types/admin";
import Icon from "../ui/icon";
import { Notifier } from "../ui/sonner";
import { askConfirmation } from "../ui/confirmation";
import { format } from "date-fns";

interface AdminWorkspacesTableProps {
  workspaces: AdminWorkspace[];
  loading: boolean;
  onSearch: (query: string) => void;
  onViewDetails: (workspaceId: string) => void;
  onDelete: (workspaceIds: string[]) => void;
  onRecover?: (workspaceIds: string[]) => void;
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
  pagination,
}: AdminWorkspacesTableProps) {
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedWorkspaces, setSelectedWorkspaces] = useState<string[]>([]);

  const handleSearch = (value: string) => {
    setSearchQuery(value);
    onSearch(value);
  };

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

  const copyToClipboard = async (text: string, type: string) => {
    try {
      await navigator.clipboard.writeText(text);
      Notifier.success(`${type} copied to clipboard`);
    } catch (err) {
      console.error("Failed to copy: ", err);
      Notifier.error(`Failed to copy ${type}`);
    }
  };

  const formatDateTime = (dateString: string) => {
    try {
      return format(new Date(dateString), "PPp"); // e.g., "Apr 29, 2023, 3:45 PM"
    } catch {
      return "Invalid date";
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
      className: "min-w-[250px]",
      render: (_, workspace) => (
        <div className="flex items-center gap-3">
          <div className="flex flex-col">
            <span className="font-medium flex items-center gap-2">
              {workspace.name}
              <div
                className="p-1 h-8 w-8 flex items-center justify-center rounded hover:bg-accent text-xs text-muted-foreground hover:text-foreground font-mono cursor-pointer"
                onClick={() => copyToClipboard(workspace.id, "Workspace ID")}
              >
                <Icon type="content_copy" />
              </div>
            </span>
            {workspace.website && (
              <a
                href={workspace.website}
                target="_blank"
                rel="noopener noreferrer"
                className="text-sm text-blue-600 hover:underline"
              >
                {workspace.website}
              </a>
            )}
            {workspace.billingEmail && (
              <span className="text-sm text-muted-foreground">
                {workspace.billingEmail}
              </span>
            )}
          </div>
        </div>
      ),
    },
    {
      key: "usersCount",
      label: "Users",
      order: 2,
      className: "min-w-[80px]",
      render: (_, workspace) => (
        <div className="flex items-center justify-center">
          <Badge variant="outline" className="font-mono">
            {workspace.usersCount ?? 0}
          </Badge>
        </div>
      ),
    },
    {
      key: "status",
      label: "Status",
      order: 3,
      className: "min-w-[120px]",
      render: (_, workspace) => (
        <div className="flex flex-col gap-1">
          <Badge variant={workspace.deletedAt ? "destructive" : "default"}>
            {workspace.deletedAt ? "Deleted" : "Active"}
          </Badge>
          {workspace.deletedByName && (
            <span className="text-xs text-muted-foreground">
              Deleted by: {workspace.deletedByName}
            </span>
          )}
        </div>
      ),
    },
    {
      key: "subscription",
      label: "Subscription",
      order: 4,
      className: "min-w-[120px]",
      render: (_, workspace) => (
        <div className="flex flex-col gap-1">
          {workspace.subscription ? (
            <>
              <Badge
                variant={
                  workspace.subscription.status === "active"
                    ? "default"
                    : "secondary"
                }
              >
                {workspace.subscription.planId}
              </Badge>
              <span className="text-xs text-muted-foreground">
                {workspace.subscription.status}
              </span>
            </>
          ) : (
            <Badge variant="outline">None</Badge>
          )}
        </div>
      ),
    },
    {
      key: "timestamps",
      label: "Timestamps",
      order: 5,
      className: "min-w-[200px]",
      render: (_, workspace) => (
        <div className="flex flex-col gap-1 text-xs">
          <div>
            <span className="text-muted-foreground">Created:</span>{" "}
            <span>{formatDateTime(workspace.createdAt)}</span>
          </div>
          <div>
            <span className="text-muted-foreground">Updated:</span>{" "}
            <span>{formatDateTime(workspace.updatedAt)}</span>
          </div>
          {workspace.deletedAt && (
            <div>
              <span className="text-muted-foreground">Deleted:</span>{" "}
              <span>{formatDateTime(workspace.deletedAt)}</span>
            </div>
          )}
        </div>
      ),
    },
    {
      key: "actions",
      label: "Actions",
      order: 6,
      className: "min-w-[200px]",
      render: (_, workspace) => (
        <div className="flex gap-2 flex-wrap">
          <Button
            variant="outline"
            size="sm"
            onClick={() => onViewDetails(workspace.id)}
          >
            View
          </Button>
          {workspace.deletedAt ? (
            onRecover && (
              <Button
                variant="outline"
                size="sm"
                onClick={() => onRecover([workspace.id])}
              >
                Recover
              </Button>
            )
          ) : (
            <Button
              variant="outline"
              size="sm"
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
              Delete
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
        <div className="w-1/3">
          <Input
            id="search-workspaces"
            placeholder="Search workspaces by name..."
            value={searchQuery}
            onChange={handleSearch}
          />
        </div>

        {selectedWorkspaces.length > 0 && (
          <div className="flex gap-2">
            {onRecover && (
              <Button variant="outline" onClick={handleBulkRecover} size="sm">
                Recover Selected ({selectedWorkspaces.length})
              </Button>
            )}
            <Button variant="destructive" onClick={handleBulkDelete} size="sm">
              Delete Selected ({selectedWorkspaces.length})
            </Button>
          </div>
        )}
      </div>

      {/* Table */}
      <DataTable
        data={workspaces}
        columns={columns}
        loading={loading}
        emptyMessage="No workspaces found"
        loadingMessage="Loading workspaces..."
        rowClassName={(workspace) =>
          workspace.deletedAt ? "bg-red-50 border-red-200" : ""
        }
      />

      {/* Pagination */}
      {pagination && (
        <div className="flex justify-between items-center">
          <p className="text-sm text-gray-600">
            Showing{" "}
            {Math.min(
              (pagination.page - 1) * pagination.perPage + 1,
              pagination.total,
            )}{" "}
            to{" "}
            {Math.min(pagination.page * pagination.perPage, pagination.total)}{" "}
            of {pagination.total} workspaces
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
