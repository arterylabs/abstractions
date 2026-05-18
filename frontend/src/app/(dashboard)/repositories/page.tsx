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
import { UpdateDescriptionModal } from "@/components/repositories/UpdateDescriptionModal";
import { RepositoryStatusBadge } from "@/components/repositories/RepositoryStatusBadge";

function RepositoryCard({
  repository,
  onAddDescription,
}: {
  repository: Repository;
  onAddDescription: (repository: Repository) => void;
}) {
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
            {repository.description ? (
              <p className="mt-0.5 line-clamp-2 text-sm text-muted-foreground">
                {repository.description}
              </p>
            ) : (
              <button
                className="mt-0.5 text-sm text-muted-foreground/50 hover:text-muted-foreground transition-colors cursor-pointer"
                onClick={(e) => {
                  e.stopPropagation();
                  onAddDescription(repository);
                }}
              >
                Add description
              </button>
            )}
          </div>
        </div>
        <div className="flex shrink-0 items-center gap-2">
          <Badge
            variant={repository.newsletterEnabled ? "success" : "secondary"}
          >
            <Icon type="mail" className="text-[10px]" />
            {repository.newsletterEnabled
              ? "Newsletters On"
              : "Newsletters Off"}
          </Badge>
          <RepositoryStatusBadge
            status={repository.status}
            indexedAt={repository.indexedAt}
            indexError={repository.indexError}
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
  const [descriptionTarget, setDescriptionTarget] = useState<Repository | null>(
    null,
  );

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
      !items.some(
        (repository) =>
          repository.status === "indexing" || repository.status === "pending",
      )
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
      {descriptionTarget && (
        <UpdateDescriptionModal
          repository={descriptionTarget}
          open={!!descriptionTarget}
          onOpenChange={(open) => !open && setDescriptionTarget(null)}
        />
      )}

      {status === "loading" && items.length === 0 ? (
        <div className="flex min-h-[240px] flex-col items-center justify-center gap-4">
          <Spinner />
          <p className="text-sm text-muted-foreground">
            Loading repositories...
          </p>
        </div>
      ) : items.length === 0 ? (
        <EmptyState
          icon="hub"
          title="Connect your first repository"
          description="Connect a repository to start tracking what your team ships — and turn that into automated context for everyone."
          cta={
            <Button variant="link" onClick={() => setIsModalOpen(true)}>
              Connect Repository
            </Button>
          }
        />
      ) : (
        <div className="grid gap-3">
          {items.map((repository) => (
            <RepositoryCard
              key={repository.id}
              repository={repository}
              onAddDescription={setDescriptionTarget}
            />
          ))}
        </div>
      )}
    </div>
  );
}
