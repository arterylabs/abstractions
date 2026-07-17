"use client";

import React, { useEffect, useCallback, useRef } from "react";
import { cn } from "@/lib/utils";

export interface InfiniteScrollProps extends React.HTMLAttributes<HTMLDivElement> {
  loading: boolean;
  hasMore: boolean;
  onLoadMore: () => void;
  children: React.ReactNode;
  height?: string;
}

export function InfiniteScroll({
  loading,
  hasMore,
  onLoadMore,
  children,
  className,
  height = "calc(100vh - 250px)",
  ...props
}: InfiniteScrollProps) {
  const containerRef = useRef<HTMLDivElement>(null);

  const handleScroll = useCallback(() => {
    if (containerRef.current && !loading && hasMore) {
      const { scrollTop, scrollHeight, clientHeight } = containerRef.current;
      if (scrollHeight - scrollTop <= clientHeight * 1.5) {
        onLoadMore();
      }
    }
  }, [loading, hasMore, onLoadMore]);

  useEffect(() => {
    const currentContainer = containerRef.current;
    if (currentContainer) {
      currentContainer.addEventListener("scroll", handleScroll);
      return () => {
        currentContainer.removeEventListener("scroll", handleScroll);
      };
    }
  }, [handleScroll]);

  return (
    <div
      ref={containerRef}
      style={{ height }}
      className={cn("overflow-auto relative", className)}
      aria-busy={loading}
      role="feed"
      aria-label="Infinite scroll content"
      {...props}
    >
      {children}
      {loading && (
        <div className="flex justify-center p-4">
          <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
        </div>
      )}
    </div>
  );
}
