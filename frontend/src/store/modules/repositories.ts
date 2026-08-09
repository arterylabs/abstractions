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
  async (
    {
      workspaceId,
      query,
    }: {
      workspaceId: string;
      query?: string;
    },
    { rejectWithValue },
  ) => {
    try {
      const response = await api.repositories.list(workspaceId, { query });
      return response.data as Repository[];
    } catch (err: unknown) {
      return rejectWithValue(getErrorMessage(err));
    }
  },
);

export const fetchRepository = createAsyncThunk(
  "repositories/fetchRepository",
  async (
    { workspaceId, repoId }: { workspaceId: string; repoId: string },
    { rejectWithValue },
  ) => {
    try {
      const response = await api.repositories.get(workspaceId, repoId);
      return response.data as Repository;
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

export const updateRepositoryDescription = createAsyncThunk(
  "repositories/updateRepositoryDescription",
  async (
    {
      workspaceId,
      repoId,
      description,
    }: {
      workspaceId: string;
      repoId: string;
      description: string;
    },
    { rejectWithValue },
  ) => {
    try {
      const response = await api.repositories.updateDescription(
        workspaceId,
        repoId,
        description,
      );
      return response.data as Repository;
    } catch (err: unknown) {
      return rejectWithValue(getErrorMessage(err));
    }
  },
);

export const deleteRepository = createAsyncThunk(
  "repositories/deleteRepository",
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
      await api.repositories.deleteRepository(workspaceId, repoId);
      return { repoId };
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
    const upsertRepository = (
      state: RepositoriesState,
      repository: Repository,
    ) => {
      const index = state.items.findIndex((item) => item.id === repository.id);
      if (index === -1) {
        state.items.push(repository);
        return;
      }
      state.items[index] = repository;
    };

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
        fetchRepository.fulfilled,
        (state, action: PayloadAction<Repository>) => {
          // Only upsert the individual item; do not touch the global list status
          // so ongoing fetchRepositories operations are not interrupted.
          upsertRepository(state, action.payload);
        },
      )
      .addCase(fetchRepository.rejected, () => {
        // Single-repo fetch failure is silently ignored at the slice level.
        // The banner component handles its own error display via useAPI state.
      })
      .addCase(
        updateRepositorySettings.fulfilled,
        (state, action: PayloadAction<Repository>) => {
          upsertRepository(state, action.payload);
          state.status = "success";
        },
      )
      .addCase(updateRepositorySettings.rejected, (state, action) => {
        state.status = "failed";
        state.error = action.payload as string;
      })
      .addCase(updateRepositoryDescription.pending, (state) => {
        state.error = null;
      })
      .addCase(
        updateRepositoryDescription.fulfilled,
        (state, action: PayloadAction<Repository>) => {
          upsertRepository(state, action.payload);
        },
      )
      .addCase(updateRepositoryDescription.rejected, (state, action) => {
        state.error = action.payload as string;
      })
      .addCase(deleteRepository.pending, (state) => {
        state.status = "loading";
        state.error = null;
      })
      .addCase(
        deleteRepository.fulfilled,
        (state, action: PayloadAction<{ repoId: string }>) => {
          state.items = state.items.filter(
            (repository) => repository.id !== action.payload.repoId,
          );
          state.status = "success";
        },
      )
      .addCase(deleteRepository.rejected, (state, action) => {
        state.status = "failed";
        state.error = action.payload as string;
      });
  },
});

export default repositoriesSlice.reducer;
