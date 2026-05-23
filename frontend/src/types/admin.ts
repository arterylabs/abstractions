import {
  User,
  SubscriptionPlanId,
  SubscriptionStatus,
  WorkspaceUser,
} from "./user";
import { Workspace } from "./workspace";
import { PaginatedResponse, PaginationParams } from "./pagination";

export interface AdminStats {
  usersCount: number;
  activeUsersCount: number;
  workspacesCount: number;
  repositoriesCount: number;
  newslettersCount: number;
}

// Analytics data interfaces
export interface AnalyticsDataPoint {
  label: string; // The formatted date/time label
  value: number; // The actual timestamp for sorting
  count: number; // The count for this data point
}

export interface AnalyticsDataResponse {
  dataPoints: AnalyticsDataPoint[];
  metric: string; // e.g., "users_created", "workspaces_created"
  dimension: string; // e.g., "date", "hour", "week"
  granularity: string; // e.g., "day", "hour", "week"
  startDate: string;
  endDate: string;
  totalCount: number; // Sum of all counts
}

export interface AdminGetAnalyticsDataRequest {
  fromDate: string; // ISO date string
  toDate: string; // ISO date string
  metric:
    | "users_created"
    | "workspaces_created"
    | "users_active"
    | "subscriptions_created"
    | "invitations_sent"
    | "repositories_created"
    | "newsletters_created";
  dimension: "date" | "hour" | "week" | "month" | "status" | "type";
  granularity: "minute" | "hour" | "day" | "week" | "month" | "year";
}

// Backward compatibility - keeping old interfaces for migration
export interface HistoricalDataPoint {
  date: string;
  count: number;
}

export interface HistoricalDataResponse {
  dataPoints: HistoricalDataPoint[];
  period: string;
  metric: string;
  startDate: string;
  endDate: string;
}

export interface AdminGetHistoricalDataRequest {
  metric: "users_created" | "workspaces_created";
  days?: number;
}

export interface AdminGetUsersRequest {
  page: number;
  perPage: number;
  query: string;
  includeDeleted?: boolean;
}

export interface AdminGetUsersResponse {
  results: User[];
  pageInfo: {
    page: number;
    perPage: number;
    totalRecords: number;
    totalPages: number;
    hasNextPage: boolean;
    hasPreviousPage: boolean;
  };
}

export interface AdminCreateUserRequest {
  email: string;
  name: string;
}

export interface AdminUpdateUserRequest {
  email?: string;
  name?: string;
}

export interface AdminDeleteUsersRequest {
  userIds: string[];
}

export interface AdminDeleteUsersResponse {
  deletedIds: string[];
}

export interface AdminRecoverUsersRequest {
  userIds: string[];
}

export interface AdminRecoverUsersResponse {
  recoveredIds: string[];
}

export interface AdminSubscriptionResponse {
  id: string;
  workspaceId: string;
  planId: SubscriptionPlanId;
  status: SubscriptionStatus;
  expiresAt: string;
  limit: number;
  user: User;
}

// Enhanced workspace interface for admin views
export interface AdminWorkspace extends Workspace {
  updatedAt: string;
  deletedAt?: string;
  deletedBy?: string;
  deletedByName?: string;
  usersCount: number;
  billingEmail: string;
  createdByName?: string;
  updatedByName?: string;
  subscription?: {
    id: string;
    planId: string;
    status: string;
    expiresAt: string;
    limit: number;
  };
}

// Workspace-related admin interfaces
export interface AdminGetWorkspacesRequest extends PaginationParams {
  includeDeleted?: boolean;
}

export type AdminGetWorkspacesResponse = PaginatedResponse<AdminWorkspace>;

export interface AdminCreateWorkspaceRequest {
  name: string;
}

export interface AdminUpdateWorkspaceRequest {
  name?: string;
  logo?: string;
  website?: string;
}

export type AdminGetWorkspaceUsersRequest = PaginationParams;

export type AdminGetWorkspaceUsersResponse = PaginatedResponse<WorkspaceUser>;

export interface AdminAddUserToWorkspaceRequest {
  userEmail: string;
  role: string;
}

export interface AdminUpdateWorkspaceUserRoleRequest {
  role: string;
}

export interface AdminRemoveUserFromWorkspaceRequest {
  workspaceId: string;
  userId: string;
}

export interface AdminBulkDeleteWorkspacesRequest {
  workspaceIds: string[];
}

export interface AdminBulkDeleteWorkspacesResponse {
  deletedIds: string[];
}

export interface AdminBulkRecoverWorkspacesRequest {
  workspaceIds: string[];
}

export interface AdminBulkRecoverWorkspacesResponse {
  recoveredIds: string[];
}

// User verification and impersonation
export interface AdminVerifyUsersRequest {
  userIds: string[];
}

export interface AdminVerifyUsersResponse {
  verifiedIds: string[];
}

export interface AdminImpersonateUserRequest {
  id: string;
}

export interface AdminImpersonateUserResponse {
  token: string;
  expiresAt: string;
}

export interface AdminExitImpersonationResponse {
  token: string;
  expiresAt: string;
}

// Subscription management
export interface AdminCreateWorkspaceSubscriptionRequest {
  workspaceId: string;
  planId: SubscriptionPlanId;
  status: SubscriptionStatus;
  expiresAt: string;
  limit: number;
}

export interface AdminUpdateWorkspaceSubscriptionRequest {
  planId?: SubscriptionPlanId;
  status?: SubscriptionStatus;
  expiresAt?: string;
  limit?: number;
}

export interface AdminDisableSubscriptionsRequest {
  workspaceIds: string[];
  disabledReason: string;
}
