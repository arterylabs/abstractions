"use client";

import { useVirtualizer } from "@tanstack/react-virtual";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { Table, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import Icon from "@/components/ui/icon";
import { Skeleton } from "@/components/ui/skeleton";

// Configurable chunk size for fetching data
export const GRID_DATA_CHUNK_SIZE = 50;

/**
 * Interface representing the column structure for the grid component.
 * - name: The name/key of the column
 * - width: Optional initial width for the column
 * Additional properties can be added in the future
 */
interface Column {
  name: string;
  width?: number;
  isAI?: boolean; // Marks the column as an AI-generated column
  // Future properties can be added here
}

/**
 * Grid component that renders a virtualized data table using the Table component.
 * It displays tabular data with virtualization for performance with large datasets.
 *
 * @param height - The height of the grid container in pixels.
 * @param columns - Array of column definitions with at least a name property.
 * @param rows - Array of arrays representing the row data.
 * @param autoHeight - Whether to automatically size the grid based on content.
 * @param rowHeight - The height of each row in pixels (default: 40).
 * @param totalRecords - Total number of records (for skeleton loaders).
 */
export default function Grid({
  height,
  columns,
  rows,
  autoHeight = false,
  rowHeight = 40,
  totalRecords,
  trailingHeader,
  trailingColumnWidth = 160,
}: {
  height?: number | string;
  columns: Column[];
  rows: string[][];
  autoHeight?: boolean;
  rowHeight?: number;
  totalRecords?: number;
  // Optional trailing header cell (e.g., an Add button) rendered after the rightmost column
  trailingHeader?: React.ReactNode;
  // Width for the trailing header/body spacer to keep header/body aligned
  trailingColumnWidth?: number;
}) {
  const parentRef = useRef<HTMLDivElement>(null);

  // State for column widths
  const [columnWidths, setColumnWidths] = useState<Record<string, number>>(() =>
    columns.reduce(
      (acc, col) => {
        acc[col.name] = col.width || 150; // Default width of 150px
        return acc;
      },
      {} as Record<string, number>,
    ),
  );

  // Keep widths in sync if columns prop changes (e.g., AI columns appended)
  useEffect(() => {
    setColumnWidths((prev) => {
      const next: Record<string, number> = { ...prev };
      // Ensure all incoming columns have a width
      for (const col of columns) {
        if (next[col.name] == null) {
          next[col.name] = col.width || 150;
        }
      }
      // Optionally remove widths for columns that no longer exist
      for (const key of Object.keys(next)) {
        if (!columns.some((c) => c.name === key)) {
          delete next[key];
        }
      }
      return next;
    });
  }, [columns]);

  // State for resize tracking
  const [isResizing, setIsResizing] = useState<string | null>(null);
  const [startX, setStartX] = useState(0);
  const [startWidth, setStartWidth] = useState(0);

  // Resize handlers
  const handleResizeStart = useCallback(
    (columnName: string, e: React.MouseEvent) => {
      setIsResizing(columnName);
      setStartX(e.clientX);
      setStartWidth(columnWidths[columnName]);
      e.preventDefault();
    },
    [columnWidths],
  );

  const handleResizeMove = useCallback(
    (e: MouseEvent) => {
      if (!isResizing) return;

      const deltaX = e.clientX - startX;
      const newWidth = Math.max(50, startWidth + deltaX); // Minimum width of 50px

      setColumnWidths((prev) => ({
        ...prev,
        [isResizing]: newWidth,
      }));
    },
    [isResizing, startX, startWidth],
  );

  const handleResizeEnd = useCallback(() => {
    setIsResizing(null);
  }, []);

  // Add global mouse event listeners for resizing
  useEffect(() => {
    if (isResizing) {
      document.addEventListener("mousemove", handleResizeMove);
      document.addEventListener("mouseup", handleResizeEnd);
      return () => {
        document.removeEventListener("mousemove", handleResizeMove);
        document.removeEventListener("mouseup", handleResizeEnd);
      };
    }
  }, [isResizing, handleResizeMove, handleResizeEnd]);

  // Memoize row data by transforming array of arrays to array of objects
  const rowData = useMemo(() => {
    if (!rows || !columns) return [];
    return rows.map((row) =>
      columns.reduce(
        (obj, column, index) => {
          obj[column.name] = row[index] || "";
          return obj;
        },
        {} as Record<string, unknown>,
      ),
    );
  }, [rows, columns]);

  // Virtualizer setup - use totalRecords if provided, otherwise fallback to rowData.length
  const virtualizer = useVirtualizer({
    count: totalRecords ?? rowData.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => rowHeight,
    overscan: 5,
  });

  const containerHeight = autoHeight ? "auto" : height;

  return (
    <div
      ref={parentRef}
      style={{ height: containerHeight }}
      className="overflow-auto border rounded-md relative"
    >
      <Table>
        <TableHeader className="sticky top-0 bg-background z-10">
          <TableRow className="flex w-full">
            {columns.map((column) => (
              <TableHead
                key={column.name}
                className="p-2 align-middle whitespace-nowrap border-r text-sm relative select-none flex-shrink-0 flex-grow-0 overflow-hidden text-ellipsis"
                style={{ width: columnWidths[column.name], display: "block" }}
              >
                <div className="flex items-center gap-1">
                  {column.isAI && (
                    <Icon
                      type="auto_awesome"
                      className="text-base text-primary"
                    />
                  )}
                  <span title={column.name} className="truncate">
                    {column.name}
                  </span>
                </div>
                <div
                  className="absolute right-0 top-0 h-full w-1 cursor-col-resize hover:bg-blue-500/50 active:bg-blue-500"
                  onMouseDown={(e) => handleResizeStart(column.name, e)}
                />
              </TableHead>
            ))}
            {trailingHeader && (
              <TableHead
                className="p-0 whitespace-nowrap text-sm relative select-none flex-shrink-0 flex-grow-0 overflow-hidden text-ellipsis"
                style={{ width: trailingColumnWidth, display: "block" }}
              >
                {trailingHeader}
              </TableHead>
            )}
          </TableRow>
        </TableHeader>
      </Table>
      <div style={{ height: virtualizer.getTotalSize() }} className="relative">
        {virtualizer.getVirtualItems().map((virtualItem) => {
          const row = rowData[virtualItem.index];
          const isRowLoaded = virtualItem.index < rowData.length;

          return (
            <div
              key={virtualItem.key}
              className="absolute w-full border-b hover:bg-muted/50 transition-colors"
              style={{
                transform: `translateY(${virtualItem.start}px)`,
                height: `${rowHeight}px`,
                display: "flex",
              }}
            >
              {columns.map((column) => (
                <div
                  key={column.name}
                  className="p-2 align-middle whitespace-nowrap border-r text-sm flex-shrink-0 flex-grow-0 overflow-hidden text-ellipsis"
                  style={{
                    width: columnWidths[column.name],
                  }}
                >
                  {isRowLoaded ? (
                    (row[column.name] as string)
                  ) : (
                    <Skeleton className="h-4 w-full" />
                  )}
                </div>
              ))}
              {trailingHeader && (
                <div
                  className="p-2 align-middle whitespace-nowrap text-sm flex-shrink-0 flex-grow-0 overflow-hidden text-ellipsis"
                  style={{ width: trailingColumnWidth }}
                />
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}
