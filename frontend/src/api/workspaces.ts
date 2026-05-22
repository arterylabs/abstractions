import { PaginatedResponse, PaginationParams } from "@/types/pagination";
import { axiosInstance } from "./index";
import { WorkspaceUser } from "@/types/user";
import { AxiosResponse } from "axios";
import { Role } from "@/types/roles";

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

interface WorkspaceConnectionStatus {
  configured: boolean;
  valid: boolean;
  error?: string;
}

export interface WorkspaceHealthResponse {
  anthropic: WorkspaceConnectionStatus;
  openAI: WorkspaceConnectionStatus;
  anthropicRequired: boolean;
  openAIRequired: boolean;
}

interface UpdateWorkspaceAISettingsParams {
  anthropicApiKey?: string | null;
  openAIApiKey?: string | null;
  indexingModel?: string | null;
  newsletterModel?: string | null;
  embeddingsModel?: string | null;
}

const workspaces = {
  getWorkspaces: async () => {
    return axiosInstance.get("/workspaces");
  },

  createWorkspace: async (data: CreateWorkspaceParams) => {
    return axiosInstance.post("/workspaces", data);
  },

  updateWorkspace: async (id: string, data: UpdateWorkspaceParams) => {
    // Map camelCase to snake_case for backend
    const payload = {
      name: data.name,
      logo: data.logo,
      website: data.website,
      about: data.about,
      billing_email: data.billingEmail,
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

  getWorkspaceSubscription: async (id: string) => {
    return axiosInstance.get(`/workspaces/${id}/subscription`);
  },

  getWorkspaceRoles: async (id: string): Promise<AxiosResponse<Role[]>> => {
    return axiosInstance.get(`/workspaces/${id}/roles`);
  },

  getAIModels: async () => {
    return axiosInstance.get("/ai/models");
  },

  getWorkspaceAISettings: async (id: string) => {
    return axiosInstance.get(`/workspaces/${id}/ai-settings`);
  },

  getWorkspaceHealth: async (
    id: string,
  ): Promise<AxiosResponse<WorkspaceHealthResponse>> => {
    return axiosInstance.get(`/workspaces/${id}/health`);
  },

  updateWorkspaceAISettings: async (
    id: string,
    data: UpdateWorkspaceAISettingsParams,
  ) => {
    return axiosInstance.put(`/workspaces/${id}/ai-settings`, data);
  },
};

export default workspaces;
