import { PaginatedResponse, PaginationParams } from "@/types/pagination";
import { axiosInstance } from "./index";
import { AxiosResponse } from "axios";
import { JSONObject } from "@/types/user";
import { Integration, IntegrationType, IntegrationTypeMetadata } from "@/types";

const integrations = {
  // Get integration types metadata (name, icon, description)
  getIntegrationTypes: async (): Promise<
    AxiosResponse<IntegrationTypeMetadata[]>
  > => {
    return axiosInstance.get("/integrations/types");
  },

  // Get all integrations for a workspace with pagination, optionally filtered by resource
  getWorkspaceIntegrations: async (
    workspaceId: string,
    params?: PaginationParams & {
      resource_type?: string;
      resource_id?: string;
    },
  ): Promise<AxiosResponse<PaginatedResponse<Integration>>> => {
    return axiosInstance.get(`/workspaces/${workspaceId}/integrations`, {
      params,
    });
  },

  // Get specific integration by ID
  getIntegration: async (
    workspaceId: string,
    integrationId: string,
  ): Promise<AxiosResponse<Integration>> => {
    return axiosInstance.get(
      `/workspaces/${workspaceId}/integrations/${integrationId}`,
    );
  },

  // Get specific integration by type for a workspace
  getIntegrationByType: async (
    workspaceId: string,
    type: IntegrationType,
  ): Promise<AxiosResponse<Integration>> => {
    return axiosInstance.get(`/workspaces/${workspaceId}/integrations/${type}`);
  },

  // Get Slack integration specifically
  getSlackIntegration: async (
    workspaceId: string,
  ): Promise<AxiosResponse<Integration>> => {
    return axiosInstance.get(
      `/workspaces/${workspaceId}/integrations/${IntegrationType.SLACK}`,
    );
  },

  // Setup integration and get authorization URL
  setupIntegration: async (
    workspaceId: string,
    type: IntegrationType,
    redirectTo?: string,
    resourceType?: string,
    resourceId?: string,
  ): Promise<AxiosResponse<{ authorize_url: string }>> => {
    const body: {
      redirectTo?: string;
      resourceType?: string;
      resourceId?: string;
    } = {};
    if (redirectTo) body.redirectTo = redirectTo;
    if (resourceType) body.resourceType = resourceType;
    if (resourceId) body.resourceId = resourceId;
    return axiosInstance.post(
      `/workspaces/${workspaceId}/integrations/${type}/setup`,
      body,
    );
  },

  // Delete an integration
  deleteIntegration: async (workspaceId: string, integrationId: string) => {
    return axiosInstance.delete(
      `/workspaces/${workspaceId}/integrations/${integrationId}`,
    );
  },

  // Update an integration
  updateIntegration: async (
    workspaceId: string,
    integrationId: string,
    data: {
      service_title?: string;
      configuration?: JSONObject;
    },
  ) => {
    return axiosInstance.put(
      `/workspaces/${workspaceId}/integrations/${integrationId}`,
      data,
    );
  },
};

export default integrations;

// Re-export types for convenience
export type {
  Integration,
  IntegrationType,
  SlackConfiguration,
  IntegrationConfiguration,
  IntegrationTypeMetadata,
} from "@/types";
