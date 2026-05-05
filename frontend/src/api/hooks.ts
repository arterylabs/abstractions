import { AxiosResponse } from "axios";
import { useState, useCallback } from "react";
import { getErrorMessage } from "@/lib/errors";

export interface ApiResponse<T> {
  result?: T | undefined;
  error?: Error | undefined;
  isFetching: boolean;
  status?: number;
}

type AsyncServiceMethod<TArgs extends unknown[], TData> = (
  ...args: TArgs
) => Promise<AxiosResponse<TData>>;

type UseApiResponse<TData, TArgs extends unknown[]> = [
  ApiResponse<TData>,
  (...args: TArgs) => Promise<void>,
];

export function useAPI<TData, TArgs extends unknown[]>(
  serviceMethod: AsyncServiceMethod<TArgs, TData>,
): UseApiResponse<TData, TArgs> {
  const [result, setResult] = useState<TData | undefined>(undefined);
  const [isFetching, setFetching] = useState(false);
  const [error, setError] = useState<Error | undefined>(undefined);
  const [status, setStatus] = useState<number | undefined>(undefined);

  const apiCallback = useCallback(
    async (...args: TArgs) => {
      setError(undefined);
      setStatus(undefined);
      setFetching(true);
      try {
        const result = await serviceMethod(...args);
        setResult(result.data);
        setStatus(result.status);
      } catch (error: unknown) {
        setError(new Error(getErrorMessage(error)));
      }
      setFetching(false);
    },
    [serviceMethod],
  );

  const apiResponse = {
    result,
    error,
    isFetching,
    status,
  };
  return [apiResponse, apiCallback];
}

export default useAPI;
