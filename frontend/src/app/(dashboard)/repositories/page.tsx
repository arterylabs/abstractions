"use client";

import Link from "next/link";
import { useCallback, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { EmptyState } from "@/components/ui/empty-state";
import { PageHeader } from "@/components/ui/page-header";
import Spinner from "@/components/ui/spinner";
import { formatDateTime } from "@/lib/format";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { fetchRepositories } from "@/store/modules/repositories";
import { Repository } from "@/types";

const STATUS_CONFIG: Record<
  Repository["status"],
  {
    label: string;
    variant: "secondary" | "warning" | "success" | "destructive";
    className?: string;
  }
> = {
  pending: { label: "Pending", variant: "secondary" },
  indexing: {
    label: "Indexing…",
    variant: "warning",
    className: "animate-pulse",
  },
  indexed: { label: "Ready", variant: "success" },
  error: { label: "Error", variant: "destructive" },
};

export function RepositoryStatusBadge({
  status,
}: {
  status: Repository["status"];
}) {
  const config = STATUS_CONFIG[status];

  return (
    <Badge variant={config.variant} className={config.className}>
      {config.label}
    </Badge>
  );
}

export default function RepositoriesPage() {
  const dispatch = useAppDispatch();
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );
  const { items, status } = useAppSelector((state) => state.repositories);

  const loadRepositories = useCallback(() => {
    if (currentWorkspace?.id) {
      dispatch(fetchRepositories(currentWorkspace.id));
    }
  }, [currentWorkspace?.id, dispatch]);

  useEffect(() => {
    loadRepositories();
  }, [loadRepositories]);

  useEffect(() => {
    if (
      !currentWorkspace?.id ||
      !items.some((repository) => repository.status === "indexing")
    ) {
      return;
    }

    const intervalId = window.setInterval(() => {
      dispatch(fetchRepositories(currentWorkspace.id));
    }, 10000);

    return () => {
      window.clearInterval(intervalId);
    };
  }, [currentWorkspace?.id, dispatch, items]);

  return (
    <div className="space-y-6">
      <PageHeader title="Repositories">
        <Button asChild>
          <Link href="/dashboard/repositories/connect">Connect Repository</Link>
        </Button>
      </PageHeader>

      {status === "loading" && items.length === 0 ? (
        <div className="flex min-h-[240px] flex-col items-center justify-center gap-4">
          <Spinner />
          <p className="text-sm text-muted-foreground">
            Loading repositories...
          </p>
        </div>
      ) : items.length === 0 ? (
        <EmptyState
          title="No repositories connected"
          description="Connect a GitHub repository to start indexing code and sending newsletters."
          cta={
            <Button asChild>
              <Link href="/dashboard/repositories/connect">
                Connect Repository
              </Link>
            </Button>
          }
        />
      ) : (
        <div className="grid gap-4">
          {items.map((repository) => (
            <Card key={repository.id}>
              <CardHeader className="gap-3 sm:flex-row sm:items-start sm:justify-between">
                <div className="space-y-2">
                  <CardTitle className="text-lg">
                    {repository.fullName}
                  </CardTitle>
                  <CardDescription>
                    Default branch: {repository.defaultBranch}
                  </CardDescription>
                </div>
                <RepositoryStatusBadge status={repository.status} />
              </CardHeader>
              <CardContent className="grid gap-2 text-sm text-muted-foreground sm:grid-cols-2">
                <div>
                  <span className="font-medium text-foreground">
                    Last indexed:
                  </span>{" "}
                  {formatDateTime(repository.indexedAt) ?? "Not indexed yet"}
                </div>
                <div>
                  <span className="font-medium text-foreground">
                    Last newsletter sent:
                  </span>{" "}
                  {formatDateTime(repository.lastNewsletterSentAt) ?? "Never"}
                </div>
                <div>
                  <span className="font-medium text-foreground">
                    Send schedule:
                  </span>{" "}
                  Day {repository.sendDay} at {repository.sendHour}:00 UTC
                </div>
                <div className="flex gap-2">
                  <Button asChild variant="outline" size="sm">
                    <Link
                      href={`/dashboard/repositories/${repository.id}/settings`}
                    >
                      Settings
                    </Link>
                  </Button>
                  <Button asChild variant="outline" size="sm">
                    <Link
                      href={`/dashboard/repositories/${repository.id}/newsletters`}
                    >
                      Newsletters
                    </Link>
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
