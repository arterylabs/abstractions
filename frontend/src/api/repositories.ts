import { AxiosResponse } from "axios";
import { Repository } from "@/types";
import { axiosInstance } from "./index";

export interface UpdateSettingsParams {
  sendDay: number;
  sendHour: number;
  newsletterEnabled: boolean;
  emails?: string[];
  sendToWorkspaceMembers?: boolean;
}

type RepositoryResponse = {
  id: string;
  workspaceId: string;
  userId: string;
  defaultBranch: string;
  status: Repository["status"];
  sendDay: number;
  sendHour: number;
  createdAt: string;
  updatedAt: string;
  indexedAt?: string | null;
  description?: string;
  newsletterEnabled?: boolean;
  sendToWorkspaceMembers?: boolean;
  serviceType?: string;
  serviceId?: string;
  name?: string;
  fullName?: string;
  lastNewsletterSentAt?: string | null;
  newsletterSentAt?: string | null;
  emails?: string[];
  recipients?: string[];
};

const normalizeRepository = (repository: RepositoryResponse): Repository => ({
  id: repository.id,
  workspaceId: repository.workspaceId,
  userId: repository.userId,
  serviceType: repository.serviceType,
  serviceId: repository.serviceId,
  name: repository.name,
  fullName: repository.fullName ?? repository.name ?? "",
  description: repository.description ?? "",
  defaultBranch: repository.defaultBranch,
  status: repository.status,
  indexedAt: repository.indexedAt ?? null,
  sendDay: repository.sendDay,
  sendHour: repository.sendHour,
  newsletterEnabled: repository.newsletterEnabled ?? true,
  lastNewsletterSentAt:
    repository.lastNewsletterSentAt ?? repository.newsletterSentAt ?? null,
  emails: repository.emails ?? repository.recipients ?? [],
  sendToWorkspaceMembers: repository.sendToWorkspaceMembers ?? true,
  createdAt: repository.createdAt,
  updatedAt: repository.updatedAt,
});

const repositories = {
  list: async (workspaceId: string): Promise<AxiosResponse<Repository[]>> => {
    const response = await axiosInstance.get<RepositoryResponse[]>(
      `/workspaces/${workspaceId}/repositories`,
    );
    return {
      ...response,
      data: response.data.map(normalizeRepository),
    };
  },

  get: async (
    workspaceId: string,
    repoId: string,
  ): Promise<AxiosResponse<Repository>> => {
    const response = await axiosInstance.get<RepositoryResponse>(
      `/workspaces/${workspaceId}/repositories/${repoId}`,
    );
    return {
      ...response,
      data: normalizeRepository(response.data),
    };
  },

  updateSettings: async (
    workspaceId: string,
    repoId: string,
    data: UpdateSettingsParams,
  ): Promise<AxiosResponse<Repository>> => {
    const response = await axiosInstance.put<RepositoryResponse>(
      `/workspaces/${workspaceId}/repositories/${repoId}`,
      data,
    );
    return {
      ...response,
      data: normalizeRepository(response.data),
    };
  },

  deleteData: async (
    workspaceId: string,
    repoId: string,
  ): Promise<AxiosResponse<void>> => {
    return axiosInstance.delete(
      `/workspaces/${workspaceId}/repositories/${repoId}/data`,
    );
  },

  triggerIndex: async (
    workspaceId: string,
    repoId: string,
  ): Promise<AxiosResponse<{ message: string }>> => {
    return axiosInstance.post(
      `/workspaces/${workspaceId}/repositories/${repoId}/trigger-index`,
    );
  },
};

export default repositories;
