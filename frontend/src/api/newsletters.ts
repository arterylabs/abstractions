import { AxiosResponse } from "axios";
import { axiosInstance } from "./index";

export interface NewsletterSummary {
  id: string;
  repoId: string;
  sentAt: string | null;
  recipientCount: number;
  openCount: number;
  clickCount: number;
  openRate: number;
  clickRate: number;
  createdAt: string;
}

export interface NewsletterDetail extends NewsletterSummary {
  topics: unknown[];
  htmlContent: string;
}

export interface UpcomingNewsletter {
  nextSendAt: string;
  topics: unknown[];
}

const newsletters = {
  list: async (repoId: string): Promise<AxiosResponse<NewsletterSummary[]>> => {
    return axiosInstance.get(`/repositories/${repoId}/newsletters`);
  },

  get: async (
    repoId: string,
    newsletterId: string,
  ): Promise<AxiosResponse<NewsletterDetail>> => {
    return axiosInstance.get(
      `/repositories/${repoId}/newsletters/${newsletterId}`,
    );
  },

  getUpcoming: async (
    repoId: string,
  ): Promise<AxiosResponse<UpcomingNewsletter>> => {
    return axiosInstance.get(`/repositories/${repoId}/newsletters/upcoming`);
  },

  regenerateTopics: async (repoId: string): Promise<AxiosResponse<unknown>> => {
    return axiosInstance.post(
      `/repositories/${repoId}/newsletters/topics/regenerate`,
    );
  },

  getEmailContent: async (
    repoId: string,
    newsletterId: string,
  ): Promise<AxiosResponse<string>> => {
    return axiosInstance.get(
      `/repositories/${repoId}/newsletters/${newsletterId}/email-content`,
      { responseType: "text" },
    );
  },
};

export default newsletters;
