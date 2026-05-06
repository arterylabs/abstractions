import { AxiosResponse } from "axios";
import { axiosInstance } from "./index";

export interface TeamMember {
  id: string;
  repoId: string;
  email: string;
  subscribed: boolean;
  addedAt: string;
}

const teamMembers = {
  list: async (repoId: string): Promise<AxiosResponse<TeamMember[]>> => {
    return axiosInstance.get(`/repositories/${repoId}/members`);
  },

  add: async (
    repoId: string,
    data: { email: string },
  ): Promise<AxiosResponse<TeamMember>> => {
    return axiosInstance.post(`/repositories/${repoId}/members`, data);
  },

  remove: async (repoId: string, memberId: string): Promise<AxiosResponse> => {
    return axiosInstance.delete(`/repositories/${repoId}/members/${memberId}`);
  },
};

export default teamMembers;
