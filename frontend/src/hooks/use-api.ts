import { AxiosResponse } from "axios";
import { useState, useCallback } from "react";
import { isAxiosError } from "axios";

/**
 * useAPI — lightweight manual-trigger data fetching hook.
 *
 * Returns [response, trigger] where trigger fires the request on demand.
 * This intentionally differs from auto-firing query libraries:
 *
 * React Query evaluation (2026-03): The [response, trigger] pattern maps
 * to useMutation, not useQuery, since calls are imperative (user actions,
 * form submits). React Query adds value for auto-firing, cached reads —
 * most useAPI call sites here are write/trigger patterns. Adopting React
 * Query would require rearchitecting call sites, not just swapping internals.
 * Decision: keep useAPI as-is; revisit if read-heavy pages need caching.
 */

export interface ApiResponse<T> {
  result?: T | undefined;
  error?: Error | undefined;
  isFetching: boolean;
  status?: number;
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any -- service methods have varied signatures; generics can't express rest params
type AnyFn = (...args: any[]) => Promise<AxiosResponse<unknown>>;
type UseApiResponse<T> = [
  ApiResponse<T>,
  (...args: Parameters<AnyFn>) => Promise<void>,
];

export function useAPI<T>(
  // eslint-disable-next-line @typescript-eslint/no-explicit-any -- intentional: hook accepts any api.* method
  serviceMethod: (...args: any[]) => Promise<AxiosResponse<T>>,
): UseApiResponse<T> {
  const [result, setResult] = useState<T | undefined>(undefined);
  const [isFetching, setFetching] = useState(false);
  const [error, setError] = useState<Error | undefined>(undefined);
  const [status, setStatus] = useState<number | undefined>(undefined);

  const apiCallback = useCallback(
    // eslint-disable-next-line @typescript-eslint/no-explicit-any -- mirrors serviceMethod signature above
    async (...args: any[]) => {
      setError(undefined);
      setStatus(undefined);
      setFetching(true);
      try {
        const response = await serviceMethod(...args);
        setResult(response.data);
        setStatus(response.status);
      } catch (err: unknown) {
        if (isAxiosError(err)) {
          setError(Error(err.response?.data?.error ?? err.message));
          setStatus(err.response?.status);
        } else {
          setError(
            err instanceof Error ? err : Error("An unexpected error occurred"),
          );
        }
      } finally {
        setFetching(false);
      }
    },
    [serviceMethod],
  );

  return [{ result, error, isFetching, status }, apiCallback];
}
