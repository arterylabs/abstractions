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
  sentAt: string | null;
  createdAt: string;
}

export interface CreateNewsletterTopicBody {
  title: string;
  description?: string;
  relatedFiles?: string[];
}

const newsletters = {
  list: async (repoId: string): Promise<AxiosResponse<NewsletterSummary[]>> => {
    return axiosInstance.get(`/repositories/${repoId}/newsletters`);
  },

  generatePreview: async (
    repoId: string,
  ): Promise<AxiosResponse<{ message: string }>> => {
    return axiosInstance.post(`/repositories/${repoId}/newsletters/preview`);
  },

  triggerNewsletter: async (
    repoId: string,
  ): Promise<AxiosResponse<{ message: string }>> => {
    return axiosInstance.post(`/repositories/${repoId}/trigger-newsletter`);
  },

  get: async (
    repoId: string,
    newsletterId: string,
  ): Promise<AxiosResponse<NewsletterDetail>> => {
    return axiosInstance.get(
      `/repositories/${repoId}/newsletters/${newsletterId}`,
    );
  },

  topics: {
    list: async (repoId: string): Promise<AxiosResponse<NewsletterTopic[]>> => {
      return axiosInstance.get(`/repositories/${repoId}/newsletter-topics`);
    },

    create: async (
      repoId: string,
      body: CreateNewsletterTopicBody,
    ): Promise<AxiosResponse<NewsletterTopic>> => {
      return axiosInstance.post(
        `/repositories/${repoId}/newsletter-topics`,
        body,
      );
    },

    delete: async (
      repoId: string,
      topicId: string,
    ): Promise<AxiosResponse<void>> => {
      return axiosInstance.delete(
        `/repositories/${repoId}/newsletter-topics/${topicId}`,
      );
    },
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

  updateContent: async (
    repoId: string,
    newsletterId: string,
    htmlContent: string,
  ): Promise<AxiosResponse<NewsletterDetail>> => {
    return axiosInstance.patch(
      `/repositories/${repoId}/newsletters/${newsletterId}`,
      { htmlContent },
    );
  },

  send: async (
    repoId: string,
    newsletterId: string,
  ): Promise<AxiosResponse<NewsletterDetail>> => {
    return axiosInstance.post(
      `/repositories/${repoId}/newsletters/${newsletterId}/send`,
    );
  },
};

export default newsletters;
