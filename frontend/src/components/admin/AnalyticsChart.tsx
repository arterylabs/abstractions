"use client";

import { useMemo } from "react";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";

interface AnalyticsDataPoint {
  label: string;
  value: number;
  count: number;
}

interface AnalyticsChartProps {
  title: string;
  data: AnalyticsDataPoint[];
  color?: string;
  className?: string;
  height?: number;
}

interface ChartDataPoint {
  label: string;
  count: number;
  date: string;
}

export default function AnalyticsChart({
  title,
  data,
  color = "#3b82f6",
  className = "",
  height = 300,
}: AnalyticsChartProps) {
  const chartData: ChartDataPoint[] = useMemo(() => {
    return data.map((point) => ({
      label: point.label,
      count: point.count,
      date: point.label,
    }));
  }, [data]);

  const totalCount = useMemo(() => {
    return data.reduce((sum, point) => sum + point.count, 0);
  }, [data]);

  const formatTooltipLabel = (label: string) => {
    // If it's a date format like "2025-06-09", format it nicely
    if (label.match(/^\d{4}-\d{2}-\d{2}$/)) {
      return new Date(label).toLocaleDateString("en-US", {
        weekday: "long",
        year: "numeric",
        month: "long",
        day: "numeric",
      });
    }
    return label;
  };

  const formatXAxisLabel = (label: string) => {
    // If it's a date format like "2025-06-09", show short format
    if (label.match(/^\d{4}-\d{2}-\d{2}$/)) {
      return new Date(label).toLocaleDateString("en-US", {
        month: "short",
        day: "numeric",
      });
    }
    return label;
  };

  return (
    <Card className={className}>
      <CardHeader>
        <CardTitle>{title}</CardTitle>
        <p className="text-sm text-muted-foreground">
          Total: <span className="font-medium">{totalCount}</span>
        </p>
      </CardHeader>

      <CardContent>
        <ResponsiveContainer width="100%" height={height}>
          <LineChart data={chartData}>
            <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
            <XAxis
              dataKey="date"
              tickFormatter={formatXAxisLabel}
              stroke="#6b7280"
              fontSize={12}
              tick={{ fill: "#6b7280" }}
            />
            <YAxis
              stroke="#6b7280"
              fontSize={12}
              tick={{ fill: "#6b7280" }}
              allowDecimals={false}
            />
            <Tooltip
              labelFormatter={formatTooltipLabel}
              formatter={(value: number) => [value, "Count"]}
              contentStyle={{
                backgroundColor: "#ffffff",
                border: "1px solid #e5e7eb",
                borderRadius: "6px",
                boxShadow: "0 4px 6px -1px rgba(0, 0, 0, 0.1)",
              }}
            />
            <Line
              type="monotone"
              dataKey="count"
              stroke={color}
              strokeWidth={2}
              dot={{ fill: color, strokeWidth: 1, r: 1 }}
              activeDot={{ r: 2, stroke: color, strokeWidth: 2 }}
            />
          </LineChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
}
