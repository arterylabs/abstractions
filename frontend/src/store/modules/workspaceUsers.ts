import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import { WorkspaceUser } from "@/types/user";
import { PaginationParams, PageInfo } from "@/types/pagination";
import api from "@/api";
import { UserRole } from "@/types/roles";
import { getErrorMessage } from "@/lib/errors";

interface WorkspaceUsersState {
  users: WorkspaceUser[] | null;
  loading: boolean;
  error: string | null;
  pageInfo: PageInfo | null;
}

const initialState: WorkspaceUsersState = {
  users: null,
  loading: false,
  error: null,
  pageInfo: null,
};

export const fetchWorkspaceUsers = createAsyncThunk(
  "workspaceUsers/fetchUsers",
  async (
    {
      workspaceId,
      params,
      append = false,
    }: {
      workspaceId: string;
      params: PaginationParams;
      append?: boolean;
    },
    { rejectWithValue },
  ) => {
    try {
      const response = await api.workspaces.getWorkspaceUsers(
        workspaceId,
        params,
      );
      console.log("Fetched workspace users:", response);
      return {
        results: response.data.results,
        pageInfo: response.data.pageInfo,
        append,
      };
    } catch (err: unknown) {
      return rejectWithValue(getErrorMessage(err));
    }
  },
);

export const updateWorkspaceUserRole = createAsyncThunk(
  "workspaceUsers/updateRole",
  async (
    {
      workspaceId,
      userId,
      role,
    }: {
      workspaceId: string;
      userId: string;
      role: string;
    },
    { rejectWithValue },
  ) => {
    try {
      await api.workspaces.updateWorkspaceUserRole(workspaceId, {
        userId,
        role,
      });
      return { userId, role: <UserRole>role };
    } catch (err: unknown) {
      return rejectWithValue(getErrorMessage(err));
    }
  },
);

export const removeWorkspaceUser = createAsyncThunk(
  "workspaceUsers/removeUser",
  async (
    { workspaceId, userId }: { workspaceId: string; userId: string },
    { rejectWithValue },
  ) => {
    try {
      await api.workspaces.removeWorkspaceUser(workspaceId, { userId });
      return userId;
    } catch (err: unknown) {
      return rejectWithValue(getErrorMessage(err));
    }
  },
);

const workspaceUsersSlice = createSlice({
  name: "workspaceUsers",
  initialState,
  reducers: {
    resetWorkspaceUsers: () => {
      return initialState;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchWorkspaceUsers.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchWorkspaceUsers.fulfilled, (state, action) => {
        const { results, pageInfo, append } = action.payload;
        state.loading = false;
        state.error = null;
        if (append && state.users) {
          state.users = [...state.users, ...results];
        } else {
          state.users = results;
        }
        state.pageInfo = pageInfo;
      })
      .addCase(fetchWorkspaceUsers.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      })
      .addCase(updateWorkspaceUserRole.fulfilled, (state, action) => {
        const { userId, role } = action.payload;
        if (!state.users) return;

        const userIndex = state.users.findIndex((user) => user.id === userId);
        if (userIndex !== -1) {
          state.users[userIndex] = {
            ...state.users[userIndex],
            role,
          };
        }
      })
      .addCase(removeWorkspaceUser.fulfilled, (state, action) => {
        if (!state.users) return;

        state.users = state.users.filter((user) => user.id !== action.payload);
        if (state.pageInfo) {
          state.pageInfo.totalRecords -= 1;
          state.pageInfo.pageRecords -= 1;
        }
      });
  },
});

export const { resetWorkspaceUsers } = workspaceUsersSlice.actions;
export default workspaceUsersSlice.reducer;
