import { axiosInstance } from "./client";

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
