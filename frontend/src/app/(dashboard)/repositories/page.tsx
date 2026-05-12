"use client";

import { useCallback, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import Icon from "@/components/ui/icon";
import { EmptyState } from "@/components/ui/empty-state";
import { PageHeader } from "@/components/ui/page-header";
import Spinner from "@/components/ui/spinner";
import { formatDateTime, formatDaysAgo } from "@/lib/format";
import Tooltip from "@/components/ui/tooltip";
import { isoWeekdayToName, utcHourToLocal } from "@/lib/timezone";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { fetchRepositories } from "@/store/modules/repositories";
import { Repository } from "@/types";
import { ConnectRepositoryModal } from "@/components/repositories/ConnectRepositoryModal";

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
  indexed: { label: "Indexed", variant: "success" },
  error: { label: "Error", variant: "destructive" },
};

export function RepositoryStatusBadge({
  status,
  indexedAt,
}: {
  status: Repository["status"];
  indexedAt?: string | null;
}) {
  const config = STATUS_CONFIG[status];

  const badge = (
    <Badge variant={config.variant} className={config.className}>
      {config.label}
    </Badge>
  );

  if (status === "indexed" && indexedAt) {
    return (
      <Tooltip value={`Last indexed ${formatDateTime(indexedAt)}`}>
        {badge}
      </Tooltip>
    );
  }

  return badge;
}

function RepositoryCard({ repository }: { repository: Repository }) {
  const router = useRouter();

  return (
    <div
      className="group cursor-pointer rounded-lg border bg-card p-4 transition-colors hover:bg-accent"
      onClick={() => router.push(`/repositories/${repository.id}`)}
    >
      {/* Top row: icon + name + badges */}
      <div className="flex items-start justify-between gap-3">
        <div className="flex min-w-0 items-start gap-3">
          <Icon
            type="deployed_code"
            className="mt-0.5 shrink-0 text-xl text-muted-foreground"
          />
          <div className="min-w-0">
            <p className="truncate font-semibold transition-colors group-hover:text-primary">
              {repository.fullName}
            </p>
            {repository.description && (
              <p className="mt-0.5 line-clamp-2 text-sm text-muted-foreground">
                {repository.description}
              </p>
            )}
          </div>
        </div>
        <div className="flex shrink-0 items-center gap-2">
          <Badge
            variant={repository.newsletterEnabled ? "success" : "secondary"}
          >
            {repository.newsletterEnabled
              ? "Newsletters On"
              : "Newsletters Off"}
          </Badge>
          <RepositoryStatusBadge
            status={repository.status}
            indexedAt={repository.indexedAt}
          />
        </div>
      </div>

      {/* Bottom row: metadata chips */}
      <div className="mt-3 flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-muted-foreground">
        <span className="flex items-center gap-1">
          <Icon type="account_tree" className="text-xs" />
          {repository.defaultBranch}
        </span>
        {repository.indexedAt ? (
          <Tooltip value={formatDateTime(repository.indexedAt) ?? ""}>
            <span className="flex cursor-help items-center gap-1">
              <Icon type="history" className="text-xs" />
              Indexed {formatDaysAgo(repository.indexedAt)}
            </span>
          </Tooltip>
        ) : (
          <span className="flex items-center gap-1">
            <Icon type="history" className="text-xs" />
            Not indexed yet
          </span>
        )}
        <span className="flex items-center gap-1">
          <Icon type="schedule" className="text-xs" />
          {isoWeekdayToName(repository.sendDay)} at{" "}
          {utcHourToLocal(repository.sendHour)}
        </span>
        {repository.lastNewsletterSentAt ? (
          <Tooltip
            value={formatDateTime(repository.lastNewsletterSentAt) ?? ""}
          >
            <span className="flex cursor-help items-center gap-1">
              <Icon type="mail" className="text-xs" />
              Last sent {formatDaysAgo(repository.lastNewsletterSentAt)}
            </span>
          </Tooltip>
        ) : (
          <span className="flex items-center gap-1">
            <Icon type="mail" className="text-xs" />
            Never sent
          </span>
        )}
      </div>
    </div>
  );
}

export default function RepositoriesPage() {
  const dispatch = useAppDispatch();
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );
  const { items, status } = useAppSelector((state) => state.repositories);
  const [isModalOpen, setIsModalOpen] = useState(false);

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
        <Button onClick={() => setIsModalOpen(true)}>Connect Repository</Button>
      </PageHeader>

      <ConnectRepositoryModal
        open={isModalOpen}
        onOpenChange={setIsModalOpen}
      />

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
            <Button onClick={() => setIsModalOpen(true)}>
              Connect Repository
            </Button>
          }
        />
      ) : (
        <div className="grid gap-3">
          {items.map((repository) => (
            <RepositoryCard key={repository.id} repository={repository} />
          ))}
        </div>
      )}
    </div>
  );
}
