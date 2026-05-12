"use client";

import { useEffect } from "react";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { format, parseISO } from "date-fns";
import api from "@/api";
import { useAPI } from "@/hooks/use-api";
import Spinner from "@/components/ui/spinner";
import { Alert } from "@/components/ui/alert";
import { HistoricalDataPoint } from "@/types/admin";

interface HistoricalChartProps {
  metric: "users_created" | "workspaces_created";
  title: string;
  days?: number;
  color?: string;
}

export default function HistoricalChart({
  metric,
  title,
  days = 30,
  color = "#3b82f6",
}: HistoricalChartProps) {
  const [historicalDataResponse, getHistoricalData] = useAPI(() =>
    api.admin.getHistoricalData({ metric, days }),
  );

  useEffect(() => {
    getHistoricalData();
  }, [getHistoricalData, metric, days]);

  if (historicalDataResponse.isFetching) {
    return (
      <div className="bg-white p-6 rounded-lg shadow">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">{title}</h3>
        <div className="flex items-center justify-center h-64">
          <Spinner />
        </div>
      </div>
    );
  }

  if (historicalDataResponse.error) {
    return (
      <div className="bg-white p-6 rounded-lg shadow">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">{title}</h3>
        <Alert variant="destructive">
          <p>
            Failed to load chart data: {historicalDataResponse.error.message}
          </p>
        </Alert>
      </div>
    );
  }

  const data = historicalDataResponse.result;

  if (!data) {
    return null;
  }

  // Transform data for Recharts
  const chartData = data.dataPoints.map((point: HistoricalDataPoint) => ({
    date: point.date,
    count: point.count,
    formattedDate: format(parseISO(point.date), "MMM dd"),
  }));

  // Calculate total for the period
  const total = data.dataPoints.reduce((sum, point) => sum + point.count, 0);

  return (
    <div className="bg-white p-6 rounded-lg shadow">
      <div className="mb-4">
        <h3 className="text-lg font-semibold text-gray-900">{title}</h3>
        <p className="text-sm text-gray-600">
          Total for last {days} days: {total}
        </p>
      </div>

      <div className="h-64">
        <ResponsiveContainer width="100%" height="100%">
          <LineChart data={chartData}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis
              dataKey="formattedDate"
              tick={{ fontSize: 12 }}
              interval="preserveStartEnd"
            />
            <YAxis tick={{ fontSize: 12 }} />
            <Tooltip
              labelFormatter={(label) => `Date: ${label}`}
              formatter={(value) => [value, "Count"]}
            />
            <Line
              type="monotone"
              dataKey="count"
              stroke={color}
              strokeWidth={2}
              dot={{ r: 3 }}
              activeDot={{ r: 5 }}
            />
          </LineChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}
