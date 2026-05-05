import { JSONObject } from "@/types/user";
import { axiosInstance } from "./index";

interface ResetPasswordParams {
  email: string;
}

interface VerifyEmailParams {
  token: string;
}

interface ResetPasswordConfirmParams {
  token: string;
  password: string;
}

interface UpdateCurrentUserParams {
  name: string;
}

interface UpdatePasswordParams {
  currentPassword: string;
  newPassword: string;
}

interface SwitchWorkspaceParams {
  workspaceId: string;
}

const users = {
  getCurrentUser: async () => {
    return axiosInstance.get("/users/me");
  },
  getOnboardingGuide: async () => {
    return axiosInstance.get("/users/me/guide");
  },
  updateOnboardingGuide: async (steps: JSONObject) => {
    return axiosInstance.post("/users/me/guide", {
      steps,
    });
  },
  resendVerificationEmail: async () => {
    return axiosInstance.post("/users/verify/resend");
  },
  resetPassword: async (data: ResetPasswordParams) => {
    return axiosInstance.post("/users/password/reset", data);
  },
  verifyEmail: async (data: VerifyEmailParams) => {
    return axiosInstance.post("/users/verify", data);
  },
  resetPasswordConfirm: async (data: ResetPasswordConfirmParams) => {
    return axiosInstance.post("/users/password/reset/confirm", data);
  },
  me: async () => {
    return axiosInstance.get("/users/me");
  },
  updateProfile: async (data: UpdateCurrentUserParams) => {
    return axiosInstance.patch("/users/me", data);
  },
  deleteAccount: async () => {
    return axiosInstance.delete("/users/me");
  },
  uploadAvatar: async (data: FormData) => {
    return axiosInstance.post("/users/me/avatar", data, {
      headers: {
        "content-type": "multipart/form-data",
      },
    });
  },
  updatePassword: async (data: UpdatePasswordParams) => {
    return axiosInstance.put("/users/me/password", data);
  },
  getPendingInvitations: async () => {
    return axiosInstance.get("/users/me/invitations");
  },
  switchWorkspace: async (data: SwitchWorkspaceParams) => {
    return axiosInstance.post("/users/me/workspaces/switch", data);
  },
};

export default users;
