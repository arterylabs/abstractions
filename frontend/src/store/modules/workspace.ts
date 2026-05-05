import { createAsyncThunk } from "@reduxjs/toolkit";
import { Workspace } from "@/types/workspace";
import api from "@/api";
import { getErrorMessage } from "@/lib/errors";

interface CreateWorkspaceResponse {
  workspace: Workspace;
}

interface UpdateWorkspaceParams {
  name: string;
  website?: string;
  about?: string;
  billingEmail?: string;
}

export const updateWorkspace = createAsyncThunk(
  "workspace/updateWorkspace",
  async (
    {
      workspaceId,
      data,
    }: {
      workspaceId: string;
      data: UpdateWorkspaceParams;
    },
    { rejectWithValue },
  ) => {
    try {
      const response = await api.workspaces.updateWorkspace(workspaceId, data);
      return response.data as Workspace;
    } catch (err: unknown) {
      return rejectWithValue(getErrorMessage(err));
    }
  },
);

export const uploadWorkspaceLogo = createAsyncThunk(
  "workspace/uploadLogo",
  async (
    {
      workspaceId,
      logo,
    }: {
      workspaceId: string;
      logo: File;
    },
    { rejectWithValue },
  ) => {
    try {
      const formData = new FormData();
      formData.append("file", logo);
      const response = await api.workspaces.uploadLogo(workspaceId, formData);
      return response.data as Workspace;
    } catch (err: unknown) {
      return rejectWithValue(getErrorMessage(err));
    }
  },
);

export const createWorkspace = createAsyncThunk(
  "workspace/createWorkspace",
  async (data: { name: string }, { rejectWithValue }) => {
    try {
      const response = await api.workspaces.createWorkspace(data);
      return response.data as CreateWorkspaceResponse;
    } catch (err: unknown) {
      return rejectWithValue(getErrorMessage(err));
    }
  },
);
