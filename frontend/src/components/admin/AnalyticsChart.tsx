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
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
  CardDescription,
} from "@/components/ui/card";

interface AnalyticsDataPoint {
  label: string;
  value: number;
  count: number;
}

interface AnalyticsChartProps {
  title: string;
  data: AnalyticsDataPoint[];
  /** CSS color or chart token value (e.g. var(--chart-1)) */
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
  color = "var(--chart-1)",
  className = "",
  height = 280,
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
      <CardHeader className="pb-2">
        <CardTitle className="text-base font-medium">{title}</CardTitle>
        <CardDescription>
          Total:{" "}
          <span className="font-medium text-foreground tabular-nums">
            {totalCount.toLocaleString()}
          </span>
        </CardDescription>
      </CardHeader>

      <CardContent>
        <ResponsiveContainer width="100%" height={height}>
          <LineChart data={chartData}>
            <CartesianGrid
              strokeDasharray="3 3"
              stroke="var(--border)"
              vertical={false}
            />
            <XAxis
              dataKey="date"
              tickFormatter={formatXAxisLabel}
              stroke="var(--muted-foreground)"
              fontSize={12}
              tick={{ fill: "var(--muted-foreground)" }}
              tickLine={false}
              axisLine={{ stroke: "var(--border)" }}
            />
            <YAxis
              stroke="var(--muted-foreground)"
              fontSize={12}
              tick={{ fill: "var(--muted-foreground)" }}
              tickLine={false}
              axisLine={false}
              allowDecimals={false}
              width={36}
            />
            <Tooltip
              labelFormatter={formatTooltipLabel}
              formatter={(value: number) => [value, "Count"]}
              contentStyle={{
                backgroundColor: "var(--popover)",
                border: "1px solid var(--border)",
                borderRadius: "var(--radius)",
                color: "var(--popover-foreground)",
                boxShadow: "0 4px 12px rgba(0, 0, 0, 0.4)",
              }}
              labelStyle={{ color: "var(--muted-foreground)" }}
              itemStyle={{ color: "var(--foreground)" }}
            />
            <Line
              type="monotone"
              dataKey="count"
              stroke={color}
              strokeWidth={2}
              dot={false}
              activeDot={{
                r: 4,
                stroke: color,
                strokeWidth: 2,
                fill: "var(--card)",
              }}
            />
          </LineChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
}
