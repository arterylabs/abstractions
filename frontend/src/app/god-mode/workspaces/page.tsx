"use client";

import { useState, useEffect, useCallback } from "react";
import { Button } from "@/components/ui/button";
import { Alert } from "@/components/ui/alert";
import { PageHeader } from "@/components/ui/page-header";
import AdminWorkspacesTable from "@/components/admin/AdminWorkspacesTable";
import WorkspaceDetailsModal from "@/components/admin/WorkspaceDetailsModal";
import { AdminGetWorkspacesRequest, AdminWorkspace } from "@/types/admin";
import api from "@/api";
import useAPI from "@/api/hooks";
import { Notifier } from "@/components/ui/sonner";
import Icon from "@/components/ui/icon";

export default function WorkspacesPage() {
  const [workspacesResponse, getWorkspaces] = useAPI(api.admin.getWorkspaces);
  const [, bulkDeleteWorkspaces] = useAPI(api.admin.bulkDeleteWorkspaces);
  const [, bulkRecoverWorkspaces] = useAPI(api.admin.bulkRecoverWorkspaces);

  const [selectedWorkspaceId, setSelectedWorkspaceId] = useState<string | null>(
    null,
  );
  const [showWorkspaceDetails, setShowWorkspaceDetails] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const [workspaces, setWorkspaces] = useState<AdminWorkspace[]>([]);
  const [totalWorkspaces, setTotalWorkspaces] = useState(0);
  const [includeDeleted, setIncludeDeleted] = useState(false);

  const perPage = 20;

  const loadWorkspaces = useCallback(async () => {
    const params: AdminGetWorkspacesRequest = {
      page: currentPage,
      perPage,
      query: searchQuery,
      includeDeleted,
    };
    await getWorkspaces(params);
  }, [getWorkspaces, currentPage, perPage, searchQuery, includeDeleted]);

  useEffect(() => {
    loadWorkspaces();
  }, [loadWorkspaces]);

  useEffect(() => {
    if (workspacesResponse.result) {
      setWorkspaces(workspacesResponse.result.results || []);
      setTotalWorkspaces(workspacesResponse.result.pageInfo?.totalRecords || 0);
    }
  }, [workspacesResponse.result]);

  const handleSearch = (query: string) => {
    setSearchQuery(query);
    setCurrentPage(1);
  };

  const handleIncludeDeletedChange = (value: boolean) => {
    setIncludeDeleted(value);
    setCurrentPage(1);
  };

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
  };

  const handleViewDetails = (workspaceId: string) => {
    setSelectedWorkspaceId(workspaceId);
    setShowWorkspaceDetails(true);
  };

  const handleCloseDetails = () => {
    setShowWorkspaceDetails(false);
    setSelectedWorkspaceId(null);
  };

  const handleDelete = async (workspaceIds: string[]) => {
    try {
      await bulkDeleteWorkspaces(workspaceIds);
      Notifier.success(
        `Successfully deleted ${workspaceIds.length} workspace(s)`,
      );
      loadWorkspaces();
    } catch {
      Notifier.error("Failed to delete workspaces");
    }
  };

  const handleRecover = async (workspaceIds: string[]) => {
    try {
      await bulkRecoverWorkspaces(workspaceIds);
      Notifier.success(
        `Successfully recovered ${workspaceIds.length} workspace(s)`,
      );
      loadWorkspaces();
    } catch {
      Notifier.error("Failed to recover workspaces");
    }
  };

  const handleWorkspaceUpdated = () => {
    loadWorkspaces();
  };

  return (
    <div className="flex flex-col gap-8">
      <PageHeader
        title="Workspaces"
        description="Search, recover, and manage every workspace on the platform"
      >
        <Button
          variant="outline"
          size="sm"
          onClick={() => loadWorkspaces()}
          disabled={workspacesResponse.isFetching}
        >
          <Icon type="refresh" className="text-sm" />
          Refresh
        </Button>
      </PageHeader>

      {workspacesResponse.error && (
        <Alert variant="destructive">
          <p>Failed to load workspaces: {workspacesResponse.error.message}</p>
        </Alert>
      )}

      <AdminWorkspacesTable
        workspaces={workspaces}
        loading={workspacesResponse.isFetching}
        onSearch={handleSearch}
        onViewDetails={handleViewDetails}
        onDelete={handleDelete}
        onRecover={handleRecover}
        includeDeleted={includeDeleted}
        onIncludeDeletedChange={handleIncludeDeletedChange}
        pagination={{
          page: currentPage,
          perPage,
          total: totalWorkspaces,
          onPageChange: handlePageChange,
        }}
      />

      <WorkspaceDetailsModal
        workspaceId={selectedWorkspaceId}
        isOpen={showWorkspaceDetails}
        onClose={handleCloseDetails}
        onWorkspaceUpdated={handleWorkspaceUpdated}
      />
    </div>
  );
}
