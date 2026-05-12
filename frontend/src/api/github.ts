import { AxiosResponse } from "axios";
import { axiosInstance } from "./index";

export interface GitHubRepo {
  id: number;
  fullName: string;
  defaultBranch: string;
  private: boolean;
  description: string;
  htmlUrl: string;
}

export interface ConnectRepoParams {
  githubRepoId: number;
  fullName: string;
  defaultBranch: string;
}

const github = {
  initiateOAuth: async (
    workspaceId: string,
    options?: { createRepository?: boolean },
  ): Promise<AxiosResponse<{ url: string }>> => {
    const body: Record<string, unknown> = {};
    if (options?.createRepository) body.createRepository = true;
    return axiosInstance.post(
      `/workspaces/${workspaceId}/integrations/github/setup`,
      body,
    );
  },

  listRepos: async (
    workspaceId: string,
  ): Promise<AxiosResponse<GitHubRepo[]>> => {
    return axiosInstance.get(`/workspaces/${workspaceId}/github/repos`);
  },

  connectRepository: async (
    workspaceId: string,
    data: ConnectRepoParams,
  ): Promise<AxiosResponse> => {
    return axiosInstance.post(`/workspaces/${workspaceId}/repositories`, data);
  },
};

export default github;
