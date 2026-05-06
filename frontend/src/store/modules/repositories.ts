import { createAsyncThunk, createSlice, PayloadAction } from "@reduxjs/toolkit";
import api from "@/api";
import { getErrorMessage } from "@/lib/errors";
import { Repository } from "@/types";
import { UpdateSettingsParams } from "@/api/repositories";

export interface RepositoriesState {
  items: Repository[];
  status: "idle" | "loading" | "success" | "failed";
  error: string | null;
}

const initialState: RepositoriesState = {
  items: [],
  status: "idle",
  error: null,
};

export const fetchRepositories = createAsyncThunk(
  "repositories/fetchRepositories",
  async (workspaceId: string, { rejectWithValue }) => {
    try {
      const response = await api.repositories.list(workspaceId);
      return response.data as Repository[];
    } catch (err: unknown) {
      return rejectWithValue(getErrorMessage(err));
    }
  },
);

export const updateRepositorySettings = createAsyncThunk(
  "repositories/updateRepositorySettings",
  async (
    {
      workspaceId,
      repoId,
      data,
    }: {
      workspaceId: string;
      repoId: string;
      data: UpdateSettingsParams;
    },
    { rejectWithValue },
  ) => {
    try {
      const response = await api.repositories.updateSettings(
        workspaceId,
        repoId,
        data,
      );
      return response.data as Repository;
    } catch (err: unknown) {
      return rejectWithValue(getErrorMessage(err));
    }
  },
);

export const deleteRepositoryData = createAsyncThunk(
  "repositories/deleteRepositoryData",
  async (
    {
      workspaceId,
      repoId,
    }: {
      workspaceId: string;
      repoId: string;
    },
    { rejectWithValue },
  ) => {
    try {
      const response = await api.repositories.deleteData(workspaceId, repoId);
      return {
        repoId,
        repository: response.data as Repository,
      };
    } catch (err: unknown) {
      return rejectWithValue(getErrorMessage(err));
    }
  },
);

const repositoriesSlice = createSlice({
  name: "repositories",
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchRepositories.pending, (state) => {
        state.status = "loading";
        state.error = null;
      })
      .addCase(
        fetchRepositories.fulfilled,
        (state, action: PayloadAction<Repository[]>) => {
          state.items = action.payload;
          state.status = "success";
        },
      )
      .addCase(fetchRepositories.rejected, (state, action) => {
        state.status = "failed";
        state.error = action.payload as string;
      })
      .addCase(updateRepositorySettings.pending, (state) => {
        state.status = "loading";
        state.error = null;
      })
      .addCase(
        updateRepositorySettings.fulfilled,
        (state, action: PayloadAction<Repository>) => {
          state.items = state.items.map((repository) =>
            repository.id === action.payload.id ? action.payload : repository,
          );
          state.status = "success";
        },
      )
      .addCase(updateRepositorySettings.rejected, (state, action) => {
        state.status = "failed";
        state.error = action.payload as string;
      })
      .addCase(deleteRepositoryData.pending, (state) => {
        state.status = "loading";
        state.error = null;
      })
      .addCase(
        deleteRepositoryData.fulfilled,
        (
          state,
          action: PayloadAction<{ repoId: string; repository?: Repository }>,
        ) => {
          const updatedRepository = action.payload.repository;

          if (updatedRepository) {
            state.items = state.items.map((repository) =>
              repository.id === action.payload.repoId
                ? updatedRepository
                : repository,
            );
          }
          state.status = "success";
        },
      )
      .addCase(deleteRepositoryData.rejected, (state, action) => {
        state.status = "failed";
        state.error = action.payload as string;
      });
  },
});

export default repositoriesSlice.reducer;
