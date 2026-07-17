import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import invitations, { Invitation } from "@/api/invitations";
import { PageInfo, PaginationParams } from "@/types/pagination";
import { getErrorMessage } from "@/lib/errors";

interface InvitationsState {
  invitations: Invitation[];
  loading: boolean;
  error: string | null;
  pageInfo: PageInfo | null;
  search: string;
}

const initialState: InvitationsState = {
  invitations: [],
  loading: false,
  error: null,
  pageInfo: null,
  search: "",
};

export const sendInvitations = createAsyncThunk(
  "invitations/sendInvitations",
  async (
    {
      workspaceId,
      emails,
      role,
    }: {
      workspaceId: string;
      emails: string[];
      role: string;
    },
    { rejectWithValue },
  ) => {
    try {
      const response = await invitations.sendInvitations(workspaceId, {
        invitations: emails.map((email) => ({ email, role })),
      });

      const data = response.data;
      const errors: { [email: string]: string } = {};
      const successfulEmails: string[] = [];
      let successCount = 0;

      data.results.forEach((result) => {
        if (result.success) {
          successCount++;
          successfulEmails.push(result.email);
        } else if (result.error) {
          errors[result.email] = result.error;
        }
      });

      return {
        errors: Object.keys(errors).length > 0 ? errors : undefined,
        successCount,
        successfulEmails,
      };
    } catch (err: unknown) {
      return rejectWithValue(getErrorMessage(err));
    }
  },
);

export const fetchInvitations = createAsyncThunk(
  "invitations/fetchInvitations",
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
      const response = await invitations.getWorkspaceInvitations(
        workspaceId,
        params,
      );
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

export const resendInvitation = createAsyncThunk(
  "invitations/resendInvitation",
  async (
    {
      workspaceId,
      invitationId,
    }: {
      workspaceId: string;
      invitationId: string;
    },
    { rejectWithValue },
  ) => {
    try {
      await invitations.resendInvitation(workspaceId, invitationId);
      return invitationId;
    } catch (err: unknown) {
      return rejectWithValue(getErrorMessage(err));
    }
  },
);

export const revokeInvitation = createAsyncThunk(
  "invitations/revokeInvitation",
  async (
    {
      workspaceId,
      invitationId,
    }: {
      workspaceId: string;
      invitationId: string;
    },
    { rejectWithValue },
  ) => {
    try {
      await invitations.revokeInvitation(workspaceId, invitationId);
      return invitationId;
    } catch (err: unknown) {
      return rejectWithValue(getErrorMessage(err));
    }
  },
);

const invitationsSlice = createSlice({
  name: "invitations",
  initialState,
  reducers: {
    setSearch: (state, action) => {
      state.search = action.payload;
    },
    resetInvitations: () => {
      return initialState;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchInvitations.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchInvitations.fulfilled, (state, action) => {
        const { results, pageInfo, append } = action.payload;
        state.loading = false;
        state.error = null;
        if (append && state.invitations) {
          state.invitations = [...state.invitations, ...results];
        } else {
          state.invitations = results;
        }
        state.pageInfo = pageInfo;
      })
      .addCase(fetchInvitations.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      })
      .addCase(resendInvitation.fulfilled, (state, action) => {
        if (!state.invitations) return;

        const index = state.invitations.findIndex(
          (item) => item.id === action.payload,
        );
        if (index !== -1) {
          state.invitations[index] = {
            ...state.invitations[index],
            expiresAt: new Date(
              Date.now() + 7 * 24 * 60 * 60 * 1000,
            ).toISOString(),
          };
        }
      })
      .addCase(revokeInvitation.fulfilled, (state, action) => {
        if (!state.invitations) return;

        const index = state.invitations.findIndex(
          (item) => item.id === action.payload,
        );
        if (index !== -1) {
          state.invitations[index] = {
            ...state.invitations[index],
            state: "ignored",
          };
        }
      });
  },
});

export const { setSearch } = invitationsSlice.actions;

export default invitationsSlice.reducer;
