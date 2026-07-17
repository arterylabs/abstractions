"use client";

import { useState, useEffect, useCallback } from "react";
import { Button } from "@/components/ui/button";
import { Alert } from "@/components/ui/alert";
import { PageHeader } from "@/components/ui/page-header";
import AdminUsersTable from "@/components/admin/AdminUsersTable";
import UserDetailsModal from "@/components/admin/UserDetailsModal";
import { User } from "@/types/user";
import { AdminGetUsersRequest } from "@/types/admin";
import api from "@/api";
import useAPI from "@/api/hooks";
import { Notifier } from "@/components/ui/sonner";
import Icon from "@/components/ui/icon";

export default function UsersPage() {
  const [usersResponse, getUsers] = useAPI(api.admin.getUsers);
  const [, deleteUsers] = useAPI(api.admin.deleteUsers);
  const [, verifyUsers] = useAPI(api.admin.verifyUsers);
  const [impersonateResponse, impersonateUser] = useAPI(
    api.admin.impersonateUser,
  );

  const [selectedUserId, setSelectedUserId] = useState<string | null>(null);
  const [showUserDetails, setShowUserDetails] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const [users, setUsers] = useState<User[]>([]);
  const [totalUsers, setTotalUsers] = useState(0);
  const [includeDeleted, setIncludeDeleted] = useState(false);

  const perPage = 20;

  const loadUsers = useCallback(async () => {
    const params: AdminGetUsersRequest = {
      page: currentPage,
      perPage,
      query: searchQuery,
      includeDeleted,
    };
    await getUsers(params);
  }, [getUsers, currentPage, perPage, searchQuery, includeDeleted]);

  useEffect(() => {
    loadUsers();
  }, [loadUsers]);

  useEffect(() => {
    if (usersResponse.result) {
      setUsers(usersResponse.result.results || []);
      setTotalUsers(usersResponse.result.pageInfo?.totalRecords || 0);
    }
  }, [usersResponse.result]);

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

  const handleViewDetails = (userId: string) => {
    setSelectedUserId(userId);
    setShowUserDetails(true);
  };

  const handleCloseDetails = () => {
    setShowUserDetails(false);
    setSelectedUserId(null);
  };

  const handleDelete = async (userIds: string[]) => {
    try {
      await deleteUsers({ userIds });
      Notifier.success(`Successfully deleted ${userIds.length} user(s)`);
      loadUsers();
    } catch {
      Notifier.error("Failed to delete users");
    }
  };

  const handleVerify = async (userIds: string[]) => {
    try {
      await verifyUsers({ userIds });
      Notifier.success(`Successfully verified ${userIds.length} user(s)`);
      loadUsers();
    } catch {
      Notifier.error("Failed to verify users");
    }
  };

  const handleImpersonate = async (userId: string) => {
    try {
      await impersonateUser({ id: userId });
    } catch {
      Notifier.error("Failed to impersonate user");
    }
  };

  useEffect(() => {
    if (impersonateResponse.result?.token) {
      const currentToken = api.utils.getToken();
      if (currentToken) {
        localStorage.setItem("impersonatorToken", currentToken);
      }

      api.utils.setToken(
        impersonateResponse.result.token,
        impersonateResponse.result.expiresAt,
      );
      window.location.href = "/dashboard";

      Notifier.success("Successfully started impersonation session");
    }

    if (impersonateResponse.error) {
      Notifier.error("Failed to impersonate user");
    }
  }, [impersonateResponse.result, impersonateResponse.error]);

  const handleUserUpdated = () => {
    loadUsers();
  };

  return (
    <div className="space-y-6">
      <PageHeader
        title="Users"
        description="Search, verify, and manage every account on the platform"
      >
        <Button
          variant="outline"
          size="sm"
          onClick={() => loadUsers()}
          disabled={usersResponse.isFetching}
        >
          <Icon type="refresh" className="text-sm" />
          Refresh
        </Button>
      </PageHeader>

      {usersResponse.error && (
        <Alert variant="destructive">
          <p>Failed to load users: {usersResponse.error.message}</p>
        </Alert>
      )}

      <AdminUsersTable
        users={users}
        loading={usersResponse.isFetching}
        onSearch={handleSearch}
        onViewDetails={handleViewDetails}
        onDelete={handleDelete}
        onVerify={handleVerify}
        onImpersonate={handleImpersonate}
        includeDeleted={includeDeleted}
        onIncludeDeletedChange={handleIncludeDeletedChange}
        pagination={{
          page: currentPage,
          perPage,
          total: totalUsers,
          onPageChange: handlePageChange,
        }}
      />

      <UserDetailsModal
        userId={selectedUserId}
        isOpen={showUserDetails}
        onClose={handleCloseDetails}
        onUserUpdated={handleUserUpdated}
      />
    </div>
  );
}
