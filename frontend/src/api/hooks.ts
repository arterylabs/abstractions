import { getErrorMessage } from "@/lib/errors";
import { AxiosResponse, isAxiosError } from "axios";
import { useState, useCallback } from "react";

export interface ApiResponse<T> {
  result?: T | undefined;
  error?: Error | undefined;
  isFetching: boolean;
  status?: number;
}

type ApiServiceMethod<TArgs extends unknown[], TResponse> = (
  ...args: TArgs
) => Promise<AxiosResponse<TResponse>>;

type UseApiResponse<TArgs extends unknown[], TResponse> = [
  ApiResponse<TResponse>,
  (...args: TArgs) => Promise<void>,
];

export function useAPI<TResponse, TArgs extends unknown[] = unknown[]>(
  serviceMethod: ApiServiceMethod<TArgs, TResponse>,
): UseApiResponse<TArgs, TResponse> {
  const [result, setResult] = useState<TResponse | undefined>(undefined);
  const [isFetching, setFetching] = useState(false);
  const [error, setError] = useState<Error | undefined>(undefined);
  const [status, setStatus] = useState<number | undefined>(undefined);

  const apiCallback = useCallback(
    async (...args: TArgs) => {
      setError(undefined);
      setStatus(undefined);
      setFetching(true);
      try {
        const response = await serviceMethod(...args);
        setResult(response.data);
        setStatus(response.status);
      } catch (err: unknown) {
        if (isAxiosError(err)) {
          setError(Error(getErrorMessage(err)));
          setStatus(err.response?.status);
        } else {
          setError(Error(getErrorMessage(err)));
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
