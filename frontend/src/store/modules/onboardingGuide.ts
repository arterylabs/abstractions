import api from "@/api";
import { JSONObject, OnboardingGuide } from "@/types/user";
import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import type { PayloadAction } from "@reduxjs/toolkit";
import { getErrorMessage } from "@/lib/errors";

export interface OnboardingState {
  steps?: JSONObject;
  onboardingGuideStatus: "idle" | "loading" | "failed" | "success";
}

const initialState: OnboardingState = {
  onboardingGuideStatus: "idle",
  steps: {},
};

export const onboardingSlice = createSlice({
  name: "onboarding",
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(getOnboardingGuide.pending, (state) => {
        state.onboardingGuideStatus = "loading";
      })
      .addCase(
        getOnboardingGuide.fulfilled,
        (state, action: PayloadAction<OnboardingGuide>) => {
          state.steps = action.payload.steps;
          state.onboardingGuideStatus = "success";
        },
      )
      .addCase(getOnboardingGuide.rejected, (state) => {
        state.onboardingGuideStatus = "failed";
      })
      .addCase(updateOnboardingGuide.pending, (state) => {
        state.onboardingGuideStatus = "loading";
      })
      .addCase(
        updateOnboardingGuide.fulfilled,
        (state, action: PayloadAction<OnboardingGuide>) => {
          state.steps = action.payload.steps;
          state.onboardingGuideStatus = "success";
        },
      )
      .addCase(updateOnboardingGuide.rejected, (state) => {
        state.onboardingGuideStatus = "failed";
      });
  },
});

export const getOnboardingGuide = createAsyncThunk(
  "onboarding/getOnboardingGuide",
  async (_args, { rejectWithValue }) => {
    try {
      const response = await api.users.getOnboardingGuide();
      return response.data as OnboardingGuide;
    } catch (err: unknown) {
      return rejectWithValue(getErrorMessage(err));
    }
  },
);

export const updateOnboardingGuide = createAsyncThunk(
  "onboarding/updateOnboardingGuide",
  async (steps: JSONObject, { rejectWithValue }) => {
    try {
      const response = await api.users.updateOnboardingGuide(steps);
      return response.data as OnboardingGuide;
    } catch (err: unknown) {
      return rejectWithValue(getErrorMessage(err));
    }
  },
);

// Action creators are generated for each case reducer function
export const {} = onboardingSlice.actions;

export default onboardingSlice.reducer;
