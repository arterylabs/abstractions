import { createSlice, createAsyncThunk, PayloadAction } from "@reduxjs/toolkit";
import api from "@/api";
import { IntegrationTypeMetadata } from "@/types";
import { getErrorMessage } from "@/lib/errors";

export interface IntegrationsState {
  types: IntegrationTypeMetadata[] | null;
  typesStatus: "idle" | "loading" | "success" | "failed";
  typesError: string | null;
}

const initialState: IntegrationsState = {
  types: null,
  typesStatus: "idle",
  typesError: null,
};

export const fetchIntegrationTypes = createAsyncThunk(
  "integrations/fetchTypes",
  async (_, { rejectWithValue }) => {
    try {
      const response = await api.integrations.getIntegrationTypes();
      return response.data as IntegrationTypeMetadata[];
    } catch (err: unknown) {
      return rejectWithValue(getErrorMessage(err));
    }
  },
);

const integrationsSlice = createSlice({
  name: "integrations",
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchIntegrationTypes.pending, (state) => {
        state.typesStatus = "loading";
        state.typesError = null;
      })
      .addCase(
        fetchIntegrationTypes.fulfilled,
        (state, action: PayloadAction<IntegrationTypeMetadata[]>) => {
          state.types = action.payload;
          state.typesStatus = "success";
        },
      )
      .addCase(fetchIntegrationTypes.rejected, (state, action) => {
        state.typesStatus = "failed";
        state.typesError = action.payload as string;
      });
  },
});

export default integrationsSlice.reducer;
