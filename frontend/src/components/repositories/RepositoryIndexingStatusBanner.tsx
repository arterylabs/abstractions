"use client";

import { useEffect, useRef } from "react";
import api from "@/api";
import useAPI from "@/api/hooks";
import { Card, CardContent } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import Spinner from "@/components/ui/spinner";
import { RepositoryStatusBadge } from "@/components/repositories/RepositoryStatusBadge";
import { useAppDispatch } from "@/store/hooks";
import { fetchRepository } from "@/store/modules/repositories";
import { Repository } from "@/types";

// How often to re-fetch the status endpoint while indexing is active.
// 5 s balances UI freshness against server load (each call is a read-only DB query).
const POLL_INTERVAL_MS = 5000;

const BULLET = " • ";

export function RepositoryIndexingStatusBanner({
  workspaceId,
  repository,
}: {
  workspaceId: string;
  repository: Repository;
}) {
  const dispatch = useAppDispatch();
  const [statusResponse, loadStatus] = useAPI(api.repositories.getStatus);
  const refreshedRepositoryRef = useRef(false);

  const shouldShow =
    repository.status === "pending" || repository.status === "indexing";
  const liveStatus = statusResponse.result?.status ?? repository.status;
  const isBusy = liveStatus === "pending" || liveStatus === "indexing";
  const totalFiles = statusResponse.result?.totalFiles ?? 0;
  const summarizedFiles = statusResponse.result?.summarizedFiles ?? 0;
  const progressPercentage = statusResponse.result?.progressPercentage ?? 0;
  const languageBreakdown = statusResponse.result?.languageBreakdown ?? [];
  const topLanguage = languageBreakdown[0];
  const languageCount = languageBreakdown.length;
  const fileProgressSummary =
    totalFiles > 0
      ? `${summarizedFiles.toLocaleString()} / ${totalFiles.toLocaleString()} files summarized`
      : "Waiting for files to be fetched";
  const repositoryStatsSummary = `${totalFiles.toLocaleString()} total files`;

  // Initial fetch + polling: merged into a single effect to avoid duplicated deps.
  useEffect(() => {
    if (!shouldShow) {
      return;
    }

    void loadStatus(workspaceId, repository.id);

    const intervalId = window.setInterval(() => {
      void loadStatus(workspaceId, repository.id);
    }, POLL_INTERVAL_MS);

    return () => {
      window.clearInterval(intervalId);
    };
  }, [loadStatus, repository.id, shouldShow, workspaceId]);

  useEffect(() => {
    if (liveStatus !== "indexed" || refreshedRepositoryRef.current) {
      return;
    }

    refreshedRepositoryRef.current = true;
    void dispatch(fetchRepository({ workspaceId, repoId: repository.id }));
  }, [dispatch, liveStatus, repository.id, workspaceId]);

  if (!shouldShow) {
    return null;
  }

  return (
    <Card className="border-primary/25 bg-accent-dark">
      <CardContent className="flex flex-col gap-3 p-5">
        <div className="flex flex-wrap items-start justify-between gap-3">
          <div className="flex flex-col gap-1.5">
            <p className="press-label text-primary">Indexing repository</p>
            <div className="flex items-center gap-2 text-sm leading-relaxed text-muted-foreground">
              <Spinner className="h-3.5 w-3.5" />
              <span>
                {isBusy
                  ? "Fetching all files in your repository and analyzing them... This may take a few minutes."
                  : "Your repository is being indexed. Finalizing the index and summarizing your code... This may take a few minutes."}
              </span>
            </div>
          </div>
          <RepositoryStatusBadge
            status={liveStatus}
            indexedAt={statusResponse.result?.indexedAt ?? repository.indexedAt}
            indexError={
              statusResponse.result?.indexError ?? repository.indexError
            }
          />
        </div>
        <div className="flex flex-wrap items-center justify-between gap-x-4 gap-y-1 font-mono text-xs text-muted-foreground">
          <p>
            {fileProgressSummary}
            {BULLET}
            {repositoryStatsSummary}
            {BULLET}
            {`${languageCount.toLocaleString()} languages`}
            {BULLET}
            {`Top: ${
              topLanguage
                ? `${topLanguage.language} (${topLanguage.fileCount} files, ${Math.round(topLanguage.percentage)}%)`
                : "—"
            }`}
          </p>
          <p className="text-foreground tabular-nums">
            {Math.round(progressPercentage)}%
          </p>
        </div>
        <Progress
          value={progressPercentage}
          aria-label="Summarization progress"
        />
      </CardContent>
    </Card>
  );
}
