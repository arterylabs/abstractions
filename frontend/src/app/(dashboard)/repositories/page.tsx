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
import { isoWeekdayToName } from "@/lib/timezone";
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
      className="press-card group relative cursor-pointer overflow-hidden border border-border bg-card p-5"
      onClick={() => router.push(`/repositories/${repository.id}`)}
    >
      {/* Top row: name set as a headline + status stamps */}
      <div className="flex items-start justify-between gap-4">
        <div className="flex min-w-0 items-start gap-3">
          <Icon
            type="deployed_code"
            className="mt-1 shrink-0 !text-base text-primary"
          />
          <div className="min-w-0">
            <h3 className="truncate text-2xl font-medium transition-colors group-hover:text-primary">
              {repository.fullName}
            </h3>
            {repository.description ? (
              <p className="mt-1.5 line-clamp-2 max-w-[70ch] text-sm leading-relaxed text-muted-foreground">
                {repository.description}
              </p>
            ) : (
              <button
                className="press-label mt-2 cursor-pointer text-muted-foreground/50 transition-colors hover:text-muted-foreground"
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

      {/* Colophon — the repository's terms, set as small print above a rule */}
      <dl className="mt-5 grid grid-cols-2 gap-x-6 gap-y-4 border-t border-border pt-4 sm:grid-cols-4">
        <div className="flex flex-col gap-1">
          <dt className="press-label text-muted-foreground/60">Branch</dt>
          <dd className="truncate font-mono text-xs text-foreground">
            {repository.defaultBranch}
          </dd>
        </div>

        <div className="flex flex-col gap-1">
          <dt className="press-label text-muted-foreground/60">Indexed</dt>
          <dd className="font-mono text-xs text-foreground">
            {repository.indexedAt ? (
              <Tooltip value={formatDateTime(repository.indexedAt) ?? ""}>
                <span className="cursor-help">
                  {formatDaysAgo(repository.indexedAt)}
                </span>
              </Tooltip>
            ) : (
              <span className="text-muted-foreground">Never</span>
            )}
          </dd>
        </div>

        <div className="flex flex-col gap-1">
          <dt className="press-label text-muted-foreground/60">Ships</dt>
          <dd className="font-mono text-xs text-foreground">
            {`${isoWeekdayToName(repository.sendDay)} ${new Date(0, 0, 1, repository.sendHour).toLocaleTimeString("en-US", { hour: "numeric", minute: "2-digit", hour12: true })} UTC`}
          </dd>
        </div>

        <div className="flex flex-col gap-1">
          <dt className="press-label text-muted-foreground/60">Last issue</dt>
          <dd className="font-mono text-xs text-foreground">
            {repository.lastNewsletterSentAt ? (
              <Tooltip
                value={formatDateTime(repository.lastNewsletterSentAt) ?? ""}
              >
                <span className="cursor-help">
                  {formatDaysAgo(repository.lastNewsletterSentAt)}
                </span>
              </Tooltip>
            ) : (
              <span className="text-muted-foreground">Never sent</span>
            )}
          </dd>
        </div>
      </dl>
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
    <div className="flex flex-col gap-8">
      <PageHeader
        eyebrow="Workspace"
        title="Repositories"
        description="Every repository Abstractions watches, and when its next issue ships."
      >
        <Button onClick={() => setIsModalOpen(true)}>
          <Icon type="add" className="!text-base" />
          Connect repository
        </Button>
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
          <p className="press-label text-muted-foreground">Setting the run…</p>
        </div>
      ) : items.length === 0 ? (
        <EmptyState
          icon="hub"
          title="Connect your first repository"
          description="Connect a repository to start tracking what your team ships — and turn that into automated context for everyone."
          cta={
            <Button onClick={() => setIsModalOpen(true)}>
              Connect repository
              <Icon type="arrow_forward" className="!text-base" />
            </Button>
          }
        />
      ) : (
        <div className="grid gap-4">
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
