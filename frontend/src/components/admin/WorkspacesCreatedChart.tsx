"use client";

import { useEffect, useCallback, useMemo } from "react";
import AnalyticsChart from "./AnalyticsChart";
import api from "@/api";
import useAPI from "@/api/hooks";
import Spinner from "@/components/ui/spinner";
import { Alert } from "@/components/ui/alert";

interface WorkspacesCreatedChartProps {
  className?: string;
}

export default function WorkspacesCreatedChart({
  className,
}: WorkspacesCreatedChartProps) {
  const dateParams = useMemo(
    () => ({
      fromDate: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000)
        .toISOString()
        .split("T")[0],
      toDate: new Date().toISOString().split("T")[0],
    }),
    [],
  );

  const serviceMethod = useCallback(
    () =>
      api.admin.getAnalyticsData({
        ...dateParams,
        metric: "workspaces_created",
        dimension: "date",
        granularity: "day",
      }),
    [dateParams],
  );

  const [response, fetchData] = useAPI(serviceMethod);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  if (response.isFetching) {
    return (
      <div className={`bg-white p-6 rounded-lg shadow ${className || ""}`}>
        <div className="flex items-center justify-center h-64">
          <Spinner />
        </div>
      </div>
    );
  }

  if (response.error) {
    return (
      <div className={className}>
        <Alert variant="destructive">
          <p>
            Failed to load workspaces created chart: {response.error.message}
          </p>
        </Alert>
      </div>
    );
  }

  if (!response.result) {
    return null;
  }

  return (
    <AnalyticsChart
      title="Workspaces Created (Last 30 Days)"
      data={response.result.dataPoints}
      color="#10b981"
      className={className}
    />
  );
}
