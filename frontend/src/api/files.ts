import { axiosInstance } from "./client";
import { AxiosRequestConfig, AxiosResponse } from "axios";

export interface File {
  id: string;
  Name: string;
  userId: string;
  mimeType: string;
  extension: string;
  sizeBytes?: number | null;
  path: string;
  /** Full HTTPS URL for the file — use this as the src for public images. */
  link: string;
  createdBy: string;
  createdAt: string;
  updatedAt: string;
  updatedBy: string;
}

const files = {
  uploadFile: async (
    file: FormData,
    config?: Pick<AxiosRequestConfig, "onUploadProgress">,
  ): Promise<AxiosResponse<File>> => {
    return axiosInstance.post(`/files`, file, {
      headers: {
        "content-type": "multipart/form-data",
      },
      ...config,
    });
  },

  getFile: async (fileId: string): Promise<AxiosResponse<File>> => {
    return axiosInstance.get(`/files/${fileId}`);
  },
};

export default files;
