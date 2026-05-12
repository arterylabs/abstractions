import api from "@/api";
import { User } from "@/types/user";
import { Workspace } from "@/types/workspace";
import { Role } from "@/types/roles";
import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import type { PayloadAction } from "@reduxjs/toolkit";
import { uploadWorkspaceLogo, updateWorkspace } from "./workspace";
import { getErrorMessage } from "@/lib/errors";

export interface AuthState {
  currentUser: User | null;
  currentUserStatus: "idle" | "loading" | "failed" | "success";
  currentWorkspace?: Workspace | null;
  currentWorkspaceStatus?: "idle" | "loading" | "failed" | "success";
  workspaces?: Workspace[] | null;
  workspacesStatus?: "idle" | "loading" | "failed" | "success";
  currentRole: string | null;
  roles: Role[] | null;
  rolesStatus: "idle" | "loading" | "success" | "failed";
  rolesError: string | null;
}

const initialState: AuthState = {
  currentUser: null,
  currentUserStatus: "idle",
  currentWorkspace: null,
  currentWorkspaceStatus: "idle",
  workspaces: null,
  workspacesStatus: "idle",
  currentRole: null,
  roles: null,
  rolesStatus: "idle",
  rolesError: null,
};

export const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    setCurrentWorkspace: (state, action: PayloadAction<Workspace>) => {
      state.currentWorkspace = action.payload;
      // Update current user's workspace role when workspace changes
      if (state.currentUser && action.payload.userRole) {
        state.currentUser = {
          ...state.currentUser,
          workspaceRole: action.payload.userRole,
        };
        // Also update the currentRole
        state.currentRole = action.payload.userRole;
      }
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(syncCurrentUser.pending, (state) => {
        state.currentUserStatus = "loading";
      })
      .addCase(
        syncCurrentUser.fulfilled,
        (state, action: PayloadAction<User>) => {
          state.currentUser = action.payload;
          state.currentUserStatus = "success";
          // Set current role from user's workspace role
          if (action.payload.workspaceRole) {
            state.currentRole = action.payload.workspaceRole;
          }
        },
      )
      .addCase(syncCurrentUser.rejected, (state) => {
        state.currentUserStatus = "failed";
      })
      .addCase(updateCurrentUser.fulfilled, (state, action) => {
        state.currentUser = {
          ...state.currentUser,
          ...action.payload,
        };
      })
      .addCase(uploadUserAvatar.fulfilled, (state, action) => {
        state.currentUser = {
          ...state.currentUser,
          ...action.payload,
        };
      })
      .addCase(uploadWorkspaceLogo.fulfilled, (state, action) => {
        const updatedWorkspace = action.payload;
        // Update current workspace if it matches
        if (state.currentWorkspace?.id === updatedWorkspace.id) {
          state.currentWorkspace = updatedWorkspace;
        }
        // Update workspace in workspaces array
        if (state.workspaces) {
          state.workspaces = state.workspaces.map((ws) =>
            ws.id === updatedWorkspace.id ? updatedWorkspace : ws,
          );
        }
      })
      .addCase(updateWorkspace.fulfilled, (state, action) => {
        const updatedWorkspace = action.payload;
        // Update current workspace if it matches
        if (state.currentWorkspace?.id === updatedWorkspace.id) {
          state.currentWorkspace = updatedWorkspace;
        }
        // Update workspace in workspaces array
        if (state.workspaces) {
          state.workspaces = state.workspaces.map((ws) =>
            ws.id === updatedWorkspace.id ? updatedWorkspace : ws,
          );
        }
      })
      .addCase(syncWorkspaces.pending, (state) => {
        state.workspacesStatus = "loading";
      })
      .addCase(syncWorkspaces.fulfilled, (state, action) => {
        state.workspaces = action.payload;
        state.workspacesStatus = "success";

        // Set current workspace based on user's default workspace or first available
        if (state.currentWorkspace === null && action.payload.length > 0) {
          const defaultWorkspace = state.currentUser?.defaultWorkspaceId
            ? action.payload.find(
                (ws) => ws.id === state.currentUser?.defaultWorkspaceId,
              )
            : action.payload[0];

          const selectedWorkspace = defaultWorkspace || action.payload[0];
          state.currentWorkspace = selectedWorkspace;

          // Update current user's workspace role when workspace is set
          if (state.currentUser && selectedWorkspace.userRole) {
            state.currentUser = {
              ...state.currentUser,
              workspaceRole: selectedWorkspace.userRole,
            };
            // Also set currentRole
            state.currentRole = selectedWorkspace.userRole;
          }
        }
      })
      .addCase(syncWorkspaces.rejected, (state) => {
        state.workspacesStatus = "failed";
      })
      .addCase(switchWorkspace.fulfilled, (state, action) => {
        state.currentUser = action.payload;
        // Update currentRole when switching workspace
        if (action.payload.workspaceRole) {
          state.currentRole = action.payload.workspaceRole;
        }
        // Update currentWorkspace to match the new defaultWorkspaceId
        if (state.workspaces && action.payload.defaultWorkspaceId) {
          const newWorkspace = state.workspaces.find(
            (ws) => ws.id === action.payload.defaultWorkspaceId,
          );
          if (newWorkspace) {
            state.currentWorkspace = newWorkspace;
          }
        }
      })
      .addCase(fetchRoles.pending, (state) => {
        state.rolesStatus = "loading";
        state.rolesError = null;
      })
      .addCase(fetchRoles.fulfilled, (state, action) => {
        state.roles = action.payload;
        state.rolesStatus = "success";
        state.rolesError = null;
      })
      .addCase(fetchRoles.rejected, (state, action) => {
        state.rolesStatus = "failed";
        state.rolesError = action.payload as string;
      });
  },
});

export const fetchRoles = createAsyncThunk(
  "auth/fetchRoles",
  async (workspaceId: string, { rejectWithValue }) => {
    try {
      const response = await api.workspaces.getWorkspaceRoles(workspaceId);
      return response.data;
    } catch (err: unknown) {
      return rejectWithValue(getErrorMessage(err));
    }
  },
);

export const syncCurrentUser = createAsyncThunk(
  "auth/syncCurrentUser",
  async (args, { rejectWithValue }) => {
    console.log("fetching current user", args);

    try {
      const response = await api.auth.getCurrentUser();
      return response.data as User;
    } catch (err: unknown) {
      console.log(err);
      return rejectWithValue(getErrorMessage(err));
    }
  },
);

export const updateCurrentUser = createAsyncThunk(
  "auth/updateCurrentUser",
  async (user: { name: string }, { rejectWithValue }) => {
    try {
      const response = await api.users.updateProfile(user);
      return response.data as User;
    } catch (err: unknown) {
      return rejectWithValue(getErrorMessage(err));
    }
  },
);

export const uploadUserAvatar = createAsyncThunk(
  "auth/uploadUserAvatar",
  async (avatar: File, { rejectWithValue }) => {
    try {
      const formData = new FormData();
      formData.append("file", avatar);
      const response = await api.users.uploadAvatar(formData);
      return response.data as User;
    } catch (err: unknown) {
      return rejectWithValue(getErrorMessage(err));
    }
  },
);

export const syncWorkspaces = createAsyncThunk(
  "auth/syncWorkspaces",
  async (args, { rejectWithValue }) => {
    try {
      const response = await api.workspaces.getWorkspaces();
      return response.data as Workspace[];
    } catch (err: unknown) {
      return rejectWithValue(getErrorMessage(err));
    }
  },
);

export const switchWorkspace = createAsyncThunk(
  "auth/switchWorkspace",
  async (workspaceId: string, { rejectWithValue }) => {
    try {
      const response = await api.users.switchWorkspace({ workspaceId });
      return response.data as User;
    } catch (err: unknown) {
      return rejectWithValue(getErrorMessage(err));
    }
  },
);

// Action creators are generated for each case reducer function
export const { setCurrentWorkspace } = authSlice.actions;

// Selector to check if a role has a specific permission
export const hasPermission = (
  state: { auth: AuthState },
  roleId: string,
  permission: string,
): boolean => {
  if (!state.auth.roles) return false;

  const role = state.auth.roles.find((r) => r.id === roleId);
  if (!role) return false;

  return role.permissions.includes(permission);
};

// Selector to get role details by ID
export const getRoleDetails = (
  state: { auth: AuthState },
  roleId: string,
): { name: string; description: string } | null => {
  if (!state.auth.roles) return null;

  const role = state.auth.roles.find((r) => r.id === roleId);
  if (!role) return null;

  return {
    name: role.name,
    description: role.description,
  };
};

export default authSlice.reducer;
