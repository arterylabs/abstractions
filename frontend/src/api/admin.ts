import { axiosInstance } from ".";
import {
  AdminStats,
  AdminGetUsersRequest,
  AdminGetUsersResponse,
  AdminCreateUserRequest,
  AdminUpdateUserRequest,
  AdminDeleteUsersRequest,
  AdminDeleteUsersResponse,
  AdminRecoverUsersRequest,
  AdminRecoverUsersResponse,
  AdminSubscriptionResponse,
  AdminWorkspace,
  AdminGetWorkspacesRequest,
  AdminGetWorkspacesResponse,
  AdminCreateWorkspaceRequest,
  AdminUpdateWorkspaceRequest,
  AdminGetWorkspaceUsersRequest,
  AdminGetWorkspaceUsersResponse,
  AdminAddUserToWorkspaceRequest,
  AdminUpdateWorkspaceUserRoleRequest,
  AdminVerifyUsersRequest,
  AdminVerifyUsersResponse,
  AdminImpersonateUserRequest,
  AdminImpersonateUserResponse,
  AdminExitImpersonationResponse,
  AdminCreateWorkspaceSubscriptionRequest,
  AdminUpdateWorkspaceSubscriptionRequest,
  AdminDisableSubscriptionsRequest,
  AnalyticsDataResponse,
  AdminGetAnalyticsDataRequest,
  HistoricalDataResponse,
  AdminGetHistoricalDataRequest,
} from "@/types/admin";
import { User } from "@/types/user";
import { Workspace } from "@/types/workspace";

const admin = {
  // Stats
  getStats: async () => {
    return axiosInstance.get<AdminStats>("/admin/stats");
  },

  // Analytics data
  getAnalyticsData: async (params: AdminGetAnalyticsDataRequest) => {
    return axiosInstance.get<AnalyticsDataResponse>("/admin/historical-data", {
      params,
    });
  },

  // Historical data (backward compatibility)
  getHistoricalData: async (params: AdminGetHistoricalDataRequest) => {
    const analyticsParams: AdminGetAnalyticsDataRequest = {
      fromDate: new Date(Date.now() - (params.days || 30) * 24 * 60 * 60 * 1000)
        .toISOString()
        .split("T")[0],
      toDate: new Date().toISOString().split("T")[0],
      metric: params.metric,
      dimension: "date",
      granularity: "day",
    };
    return axiosInstance.get<HistoricalDataResponse>("/admin/historical-data", {
      params: analyticsParams,
    });
  },

  // User management
  getUsers: async ({
    page,
    perPage,
    query,
    includeDeleted,
  }: AdminGetUsersRequest) => {
    return axiosInstance.get<AdminGetUsersResponse>("/admin/users", {
      params: { page, perPage, query, includeDeleted },
    });
  },

  getUser: async (id: string) => {
    return axiosInstance.get<User>(`/admin/users/${id}`);
  },

  createUser: async (data: AdminCreateUserRequest) => {
    return axiosInstance.post<User>("/admin/users", data);
  },

  updateUser: async (id: string, data: AdminUpdateUserRequest) => {
    return axiosInstance.put<User>(`/admin/users/${id}`, data);
  },

  verifyUsers: async (data: AdminVerifyUsersRequest) => {
    return axiosInstance.post<AdminVerifyUsersResponse>(
      "/admin/users/verify",
      data,
    );
  },

  deleteUsers: async (request: AdminDeleteUsersRequest) => {
    return axiosInstance.delete<AdminDeleteUsersResponse>("/admin/users", {
      data: request,
    });
  },

  recoverUsers: async (request: AdminRecoverUsersRequest) => {
    return axiosInstance.post<AdminRecoverUsersResponse>(
      "/admin/users/recover",
      request,
    );
  },

  impersonateUser: async (data: AdminImpersonateUserRequest) => {
    return axiosInstance.post<AdminImpersonateUserResponse>(
      "/admin/users/impersonate",
      data,
    );
  },

  exitImpersonation: async () => {
    return axiosInstance.post<AdminExitImpersonationResponse>(
      "/admin/users/exit_impersonation",
    );
  },

  // Workspace management
  getWorkspaces: async (params: AdminGetWorkspacesRequest) => {
    return axiosInstance.get<AdminGetWorkspacesResponse>("/admin/workspaces", {
      params,
    });
  },

  getWorkspace: async (workspaceId: string) => {
    return axiosInstance.get<AdminWorkspace>(
      `/admin/workspaces/${workspaceId}`,
    );
  },

  createWorkspace: async (data: AdminCreateWorkspaceRequest) => {
    return axiosInstance.post<Workspace>("/admin/workspaces", data);
  },

  updateWorkspace: async (
    workspaceId: string,
    data: AdminUpdateWorkspaceRequest,
  ) => {
    return axiosInstance.put<Workspace>(
      `/admin/workspaces/${workspaceId}`,
      data,
    );
  },

  bulkDeleteWorkspaces: async (workspaceIds: string[]) => {
    return axiosInstance.delete("/admin/workspaces", {
      data: { workspaceIds },
    });
  },

  bulkRecoverWorkspaces: async (workspaceIds: string[]) => {
    return axiosInstance.post("/admin/workspaces/recover", {
      workspaceIds,
    });
  },

  // Workspace users
  getWorkspaceUsers: async (
    workspaceId: string,
    params: AdminGetWorkspaceUsersRequest,
  ) => {
    return axiosInstance.get<AdminGetWorkspaceUsersResponse>(
      `/admin/workspaces/${workspaceId}/users`,
      { params },
    );
  },

  addUserToWorkspace: async (
    workspaceId: string,
    data: AdminAddUserToWorkspaceRequest,
  ): Promise<void> => {
    await axiosInstance.post(`/admin/workspaces/${workspaceId}/users`, data);
  },

  updateWorkspaceUserRole: async (
    workspaceId: string,
    userId: string,
    data: AdminUpdateWorkspaceUserRoleRequest,
  ): Promise<void> => {
    await axiosInstance.put(
      `/admin/workspaces/${workspaceId}/users/${userId}`,
      data,
    );
  },

  removeUserFromWorkspace: async (
    workspaceId: string,
    userId: string,
  ): Promise<void> => {
    await axiosInstance.delete(
      `/admin/workspaces/${workspaceId}/users/${userId}`,
    );
  },

  // Get user workspaces with roles
  getUserWorkspacesWithRole: async (userId: string) => {
    return axiosInstance.get(`/admin/users/${userId}/workspaces`);
  },

  // Workspace subscriptions
  getWorkspaceSubscription: async (workspaceId: string) => {
    return axiosInstance.get<AdminSubscriptionResponse>(
      `/admin/workspaces/${workspaceId}/subscription`,
    );
  },

  createWorkspaceSubscription: async (
    workspaceId: string,
    data: AdminCreateWorkspaceSubscriptionRequest,
  ): Promise<AdminSubscriptionResponse> => {
    const response = await axiosInstance.post<AdminSubscriptionResponse>(
      `/admin/workspaces/${workspaceId}/subscription`,
      data,
    );
    return response.data;
  },

  updateWorkspaceSubscription: async (
    workspaceId: string,
    data: AdminUpdateWorkspaceSubscriptionRequest,
  ): Promise<AdminSubscriptionResponse> => {
    const response = await axiosInstance.put<AdminSubscriptionResponse>(
      `/admin/workspaces/${workspaceId}/subscription`,
      data,
    );
    return response.data;
  },

  disableWorkspaceSubscription: async (
    workspaceId: string,
    data: { disabledReason: string },
  ): Promise<void> => {
    await axiosInstance.post(
      `/admin/workspaces/${workspaceId}/subscription/disable`,
      data,
    );
  },

  // Bulk subscription management
  createSubscription: async (
    data: AdminCreateWorkspaceSubscriptionRequest,
  ): Promise<AdminSubscriptionResponse> => {
    const response = await axiosInstance.post<AdminSubscriptionResponse>(
      "/admin/subscriptions",
      data,
    );
    return response.data;
  },

  updateSubscription: async (
    id: string,
    data: AdminUpdateWorkspaceSubscriptionRequest,
  ): Promise<AdminSubscriptionResponse> => {
    const response = await axiosInstance.put<AdminSubscriptionResponse>(
      `/admin/subscriptions/${id}`,
      data,
    );
    return response.data;
  },

  disableSubscriptions: async (
    data: AdminDisableSubscriptionsRequest,
  ): Promise<void> => {
    await axiosInstance.post("/admin/subscriptions/disable", data);
  },
};

export default admin;
