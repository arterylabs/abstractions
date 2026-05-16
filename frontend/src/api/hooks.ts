import { AxiosResponse, isAxiosError } from "axios";
import { useState, useCallback } from "react";

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

export default useAPI;
