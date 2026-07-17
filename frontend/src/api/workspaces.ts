import { PaginatedResponse, PaginationParams } from "@/types/pagination";
import { axiosInstance } from "./client";
import { WorkspaceUser } from "@/types/user";
import { AxiosResponse } from "axios";
import { Role } from "@/types/roles";
import { AIModel, AIProvider, WorkspaceAISettings } from "@/types/workspace";

interface CreateWorkspaceParams {
  name: string;
  logo?: string;
  website?: string;
  about?: string;
}

interface UpdateWorkspaceParams {
  name: string;
  logo?: string;
  website?: string;
  about?: string;
  billingEmail?: string;
}

interface UpdateWorkspaceUserRoleParams {
  userId: string;
  role: string;
}

interface RemoveWorkspaceUserParams {
  userId: string;
}

export interface WorkspaceProviderHealth {
  configured: boolean;
  valid: boolean;
  error?: string;
  required: boolean;
}

export interface WorkspaceHealthResponse {
  // Health status keyed by provider ID (e.g. "anthropic", "openai", "grok").
  providers: Record<string, WorkspaceProviderHealth>;
}

export interface WorkspaceGettingStartedResponse {
  aiConfigured: boolean;
  githubConnected: boolean;
  newsletterSent: boolean;
}

interface UpdateWorkspaceAISettingsParams {
  // Provider API keys keyed by provider ID; empty string clears a key, omitted keys are untouched.
  apiKeys?: Record<string, string | null>;
  indexingModel?: string | null;
  newsletterModel?: string | null;
  embeddingsModel?: string | null;
  // System prompt enhancements; empty string clears an enhancement (built-in default only).
  // Non-empty values are appended on top of the default — they never replace it. Omitted fields are untouched.
  systemPrompts?: {
    summarizeCode?: string | null;
    summarizeFile?: string | null;
    topicGenerator?: string | null;
    newsletterEntry?: string | null;
  };
}

const workspaces = {
  getWorkspaces: async () => {
    return axiosInstance.get("/workspaces");
  },

  createWorkspace: async (data: CreateWorkspaceParams) => {
    return axiosInstance.post("/workspaces", data);
  },

  updateWorkspace: async (id: string, data: UpdateWorkspaceParams) => {
    const payload = {
      name: data.name,
      logo: data.logo,
      website: data.website,
      about: data.about,
      billingEmail: data.billingEmail,
    };
    return axiosInstance.put(`/workspaces/${id}`, payload);
  },

  uploadLogo: async (id: string, logo: FormData) => {
    return axiosInstance.post(`/workspaces/${id}/logo`, logo, {
      headers: {
        "content-type": "multipart/form-data",
      },
    });
  },

  deleteWorkspace: async (id: string) => {
    return axiosInstance.delete(`/workspaces/${id}`);
  },

  // Workspace Users
  getWorkspaceUsers: async (
    id: string,
    params?: PaginationParams,
  ): Promise<AxiosResponse<PaginatedResponse<WorkspaceUser>>> => {
    return axiosInstance.get(`/workspaces/${id}/users`, { params });
  },

  updateWorkspaceUserRole: async (
    id: string,
    data: UpdateWorkspaceUserRoleParams,
  ) => {
    return axiosInstance.put(`/workspaces/${id}/users`, data);
  },

  removeWorkspaceUser: async (id: string, data: RemoveWorkspaceUserParams) => {
    return axiosInstance.delete(`/workspaces/${id}/users`, { data });
  },

  getWorkspaceRoles: async (id: string): Promise<AxiosResponse<Role[]>> => {
    return axiosInstance.get(`/workspaces/${id}/roles`);
  },

  getAIModels: async (): Promise<AxiosResponse<AIModel[]>> => {
    return axiosInstance.get("/ai/models");
  },

  getAIProviders: async (): Promise<AxiosResponse<AIProvider[]>> => {
    return axiosInstance.get("/ai/providers");
  },

  getWorkspaceAISettings: async (
    id: string,
  ): Promise<AxiosResponse<WorkspaceAISettings>> => {
    return axiosInstance.get(`/workspaces/${id}/ai_settings`);
  },

  getWorkspaceHealth: async (
    id: string,
  ): Promise<AxiosResponse<WorkspaceHealthResponse>> => {
    return axiosInstance.get(`/workspaces/${id}/health`);
  },

  getGettingStarted: async (
    id: string,
  ): Promise<AxiosResponse<WorkspaceGettingStartedResponse>> => {
    return axiosInstance.get(`/workspaces/${id}/getting_started`);
  },

  updateWorkspaceAISettings: async (
    id: string,
    data: UpdateWorkspaceAISettingsParams,
  ) => {
    return axiosInstance.put(`/workspaces/${id}/ai_settings`, data);
  },
};

export default workspaces;
