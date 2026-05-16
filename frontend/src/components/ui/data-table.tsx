"use client";

import * as React from "react";
import { cn } from "@/lib/utils";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import Spinner from "@/components/ui/spinner";

export interface DataTableColumn<TData = Record<string, unknown>> {
  key: string;
  label: string;
  order: number;
  render?: (value: unknown, row: TData) => React.ReactNode;
  headerRender?: () => React.ReactNode;
  className?: string;
  headerClassName?: string;
}

export interface DataTableProps<TData = Record<string, unknown>> {
  columns: DataTableColumn<TData>[];
  data: TData[];
  className?: string;
  loading?: boolean;
  emptyMessage?: string;
  rowClassName?: (row: TData) => string;
}

function DataTable<TData = Record<string, unknown>>({
  columns,
  data,
  className,
  loading = false,
  emptyMessage = "No data available",
  rowClassName,
}: DataTableProps<TData>) {
  // Sort columns by order
  const sortedColumns = React.useMemo(() => {
    return [...columns].sort((a, b) => a.order - b.order);
  }, [columns]);

  const renderCellContent = React.useCallback(
    (column: DataTableColumn<TData>, row: TData) => {
      const value = (row as Record<string, unknown>)[column.key];

      if (column.render) {
        return column.render(value, row);
      }

      // Default rendering for common types
      if (value === null || value === undefined) {
        return <span className="text-muted-foreground">—</span>;
      }

      if (typeof value === "boolean") {
        return value ? "Yes" : "No";
      }

      if (value instanceof Date) {
        return value.toLocaleDateString();
      }

      if (
        typeof value === "string" &&
        value.includes("T") &&
        value.includes("Z")
      ) {
        // Assume it's an ISO date string
        try {
          return new Date(value).toLocaleDateString();
        } catch {
          return value;
        }
      }

      return String(value);
    },
    [],
  );

  return (
    <div className={cn("border rounded-lg overflow-hidden", className)}>
      <Table>
        <TableHeader className="bg-muted">
          <TableRow className="border-b hover:bg-muted">
            {sortedColumns.map((column) => (
              <TableHead key={column.key} className={column.headerClassName}>
                {column.headerRender ? column.headerRender() : column.label}
              </TableHead>
            ))}
          </TableRow>
        </TableHeader>
        <TableBody className="bg-background">
          {loading ? (
            <TableRow>
              <TableCell
                colSpan={sortedColumns.length}
                className="text-center py-8"
              >
                <div className="flex justify-center">
                  <Spinner />
                </div>
              </TableCell>
            </TableRow>
          ) : data.length === 0 ? (
            <TableRow>
              <TableCell
                colSpan={sortedColumns.length}
                className="text-center py-8 text-muted-foreground"
              >
                {emptyMessage}
              </TableCell>
            </TableRow>
          ) : (
            data.map((row, rowIndex) => (
              <TableRow
                key={rowIndex}
                className={cn(
                  "bg-background border-b hover:bg-muted/20",
                  rowClassName ? rowClassName(row) : "",
                )}
              >
                {sortedColumns.map((column) => (
                  <TableCell key={column.key} className={column.className}>
                    {renderCellContent(column, row)}
                  </TableCell>
                ))}
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </div>
  );
}

export { DataTable };
