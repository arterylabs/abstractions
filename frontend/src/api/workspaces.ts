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
};

export default workspaces;
