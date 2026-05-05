import { PaginationParams, PaginatedResponse } from "@/types/pagination";
import { axiosInstance } from "./index";
import { SendInvitationsResponse } from "@/types/user";
import { AxiosResponse } from "axios";

interface SendInvitationsParams {
  invitations: {
    email: string;
    role: string;
  }[];
}

export interface Invitation {
  id: string;
  workspaceId: string;
  workspaceName: string;
  invitedByUserId: string;
  email: string;
  role: string;
  state: string;
  expiresAt: string;
  createdAt: string;
  updatedAt: string;
  invitedBy?: {
    name: string;
    email: string;
    avatar?: string;
  };
}

const invitations = {
  // Get pending invitations for the current user
  getUserPendingInvitations: async (params?: PaginationParams) => {
    return axiosInstance.get<PaginatedResponse<Invitation>>(
      "/users/me/invitations",
      {
        params: {
          ...params,
        },
      },
    );
  },

  // Get all invitations for a workspace with pagination
  getWorkspaceInvitations: async (
    workspaceId: string,
    params?: PaginationParams,
  ) => {
    return axiosInstance.get<PaginatedResponse<Invitation>>(
      `/workspaces/${workspaceId}/invitations`,
      {
        params: {
          ...params,
        },
      },
    );
  },

  // Send invitations to multiple email addresses
  sendInvitations: async (
    workspaceId: string,
    data: SendInvitationsParams,
  ): Promise<AxiosResponse<SendInvitationsResponse>> => {
    return axiosInstance.post(`/workspaces/${workspaceId}/invitations`, data);
  },

  // Accept a workspace invitation (for users)
  acceptInvitation: async (invitationId: string) => {
    return axiosInstance.post(`/users/me/invitations/${invitationId}/accept`);
  },

  // Ignore/reject a workspace invitation (for users)
  ignoreInvitation: async (invitationId: string) => {
    return axiosInstance.post(`/users/me/invitations/${invitationId}/ignore`);
  },

  // Revoke an invitation (for workspace admins)
  revokeInvitation: async (workspaceId: string, invitationId: string) => {
    return axiosInstance.post(
      `/workspaces/${workspaceId}/invitations/${invitationId}/revoke`,
    );
  },

  // Resend an invitation
  resendInvitation: async (workspaceId: string, invitationId: string) => {
    return axiosInstance.post(
      `/workspaces/${workspaceId}/invitations/${invitationId}/resend`,
    );
  },
};

export default invitations;
