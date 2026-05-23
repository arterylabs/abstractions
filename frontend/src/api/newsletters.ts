import { AxiosResponse } from "axios";
import { axiosInstance } from "./index";

export interface NewsletterSummary {
  id: string;
  repoId: string;
  subject: string;
  sentAt: string | null;
  recipientCount: number;
  openCount: number;
  clickCount: number;
  openRate: number;
  clickRate: number;
  createdAt: string;
}

export interface NewsletterDetail extends NewsletterSummary {
  topicId: string | null;
  topics: NewsletterTopic[];
  htmlContent: string;
}

export interface NewsletterTopic {
  id: string;
  repoId: string;
  title: string;
  description: string;
  relatedFiles: string[];
  position: number;
  sentAt: string | null;
  requestedByUserId: string | null;
  requestedByName: string | null;
  requestedByAvatar: string | null;
  newsletterId: string | null;
  createdAt: string;
}

export interface CreateNewsletterTopicBody {
  title: string;
  description?: string;
  relatedFiles?: string[];
}

export interface UpdateNewsletterTopicBody {
  title: string;
  description?: string;
  relatedFiles?: string[];
}

export type ReorderDirection = "up" | "down" | "top";

const newsletters = {
  list: async (
    workspaceId: string,
    repoId: string,
  ): Promise<AxiosResponse<NewsletterSummary[]>> => {
    return axiosInstance.get(
      `/workspaces/${workspaceId}/repositories/${repoId}/newsletters`,
    );
  },

  triggerNewsletter: async (
    workspaceId: string,
    repoId: string,
  ): Promise<AxiosResponse<{ message: string }>> => {
    return axiosInstance.post(
      `/workspaces/${workspaceId}/repositories/${repoId}/trigger-newsletter`,
    );
  },

  get: async (
    workspaceId: string,
    repoId: string,
    newsletterId: string,
  ): Promise<AxiosResponse<NewsletterDetail>> => {
    return axiosInstance.get(
      `/workspaces/${workspaceId}/repositories/${repoId}/newsletters/${newsletterId}`,
    );
  },

  topics: {
    list: async (
      workspaceId: string,
      repoId: string,
    ): Promise<AxiosResponse<NewsletterTopic[]>> => {
      return axiosInstance.get(
        `/workspaces/${workspaceId}/repositories/${repoId}/newsletter-topics`,
      );
    },

    create: async (
      workspaceId: string,
      repoId: string,
      body: CreateNewsletterTopicBody,
    ): Promise<AxiosResponse<NewsletterTopic>> => {
      return axiosInstance.post(
        `/workspaces/${workspaceId}/repositories/${repoId}/newsletter-topics`,
        body,
      );
    },

    update: async (
      workspaceId: string,
      repoId: string,
      topicId: string,
      body: UpdateNewsletterTopicBody,
    ): Promise<AxiosResponse<NewsletterTopic>> => {
      return axiosInstance.put(
        `/workspaces/${workspaceId}/repositories/${repoId}/newsletter-topics/${topicId}`,
        body,
      );
    },

    reorder: async (
      workspaceId: string,
      repoId: string,
      topicId: string,
      direction: ReorderDirection,
    ): Promise<AxiosResponse<void>> => {
      return axiosInstance.put(
        `/workspaces/${workspaceId}/repositories/${repoId}/newsletter-topics/${topicId}/reorder`,
        { direction },
      );
    },

    delete: async (
      workspaceId: string,
      repoId: string,
      topicId: string,
    ): Promise<AxiosResponse<void>> => {
      return axiosInstance.delete(
        `/workspaces/${workspaceId}/repositories/${repoId}/newsletter-topics/${topicId}`,
      );
    },
  },
};

export default newsletters;
