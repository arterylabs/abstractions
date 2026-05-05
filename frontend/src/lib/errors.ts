import { isAxiosError } from "axios";

/**
 * Extracts a user-facing error message from any thrown value.
 * Use in createAsyncThunk catch blocks with rejectWithValue:
 *
 *   } catch (err: unknown) {
 *     return rejectWithValue(getErrorMessage(err));
 *   }
 */
export function getErrorMessage(err: unknown): string {
  if (isAxiosError(err)) {
    return (
      err.response?.data?.error ?? err.response?.data?.message ?? err.message
    );
  }
  if (err instanceof Error) return err.message;
  return "An unexpected error occurred";
}
