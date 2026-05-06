import { AxiosResponse } from "axios";
import { Repository } from "@/types";
import { axiosInstance } from "./index";

export interface UpdateSettingsParams {
  sendDay: number;
  sendHour: number;
}

const repositories = {
  list: async (workspaceId: string): Promise<AxiosResponse<Repository[]>> => {
    return axiosInstance.get(`/workspaces/${workspaceId}/repositories`);
  },

  get: async (
    workspaceId: string,
    repoId: string,
  ): Promise<AxiosResponse<Repository>> => {
    return axiosInstance.get(
      `/workspaces/${workspaceId}/repositories/${repoId}`,
    );
  },

  updateSettings: async (
    workspaceId: string,
    repoId: string,
    data: UpdateSettingsParams,
  ): Promise<AxiosResponse<Repository>> => {
    return axiosInstance.put(
      `/workspaces/${workspaceId}/repositories/${repoId}`,
      data,
    );
  },

  deleteData: async (
    workspaceId: string,
    repoId: string,
  ): Promise<AxiosResponse<Repository>> => {
    return axiosInstance.delete(
      `/workspaces/${workspaceId}/repositories/${repoId}/data`,
    );
  },
};

export default repositories;
