"use client";

import { Badge } from "@/components/ui/badge";
import Icon from "@/components/ui/icon";
import Tooltip from "@/components/ui/tooltip";
import { formatDateTime, formatDaysAgo } from "@/lib/format";
import {
  describeRecipients,
  formatLocalSendWindow,
  formatSendWindow,
} from "@/lib/repository";
import { daysUntilNextSend, getBrowserTimezone } from "@/lib/timezone";
import { Repository } from "@/types";
import { RepositoryStatusBadge } from "./RepositoryStatusBadge";

/** A single term in the meta rail: a running head over its value. */
function Fact({
  label,
  children,
}: {
  label: string;
  children: React.ReactNode;
}) {
  return (
    <div className="flex min-w-0 flex-col gap-1 text-right">
      <dt className="press-label text-muted-foreground/60">{label}</dt>
      <dd className="truncate font-mono text-xs text-foreground">{children}</dd>
    </div>
  );
}

/** A meta value that is not set yet, set in the same measure but quieter. */
function Unset({ children }: { children: React.ReactNode }) {
  return <span className="text-muted-foreground">{children}</span>;
}

/**
 * Right rail of standing repository facts — status badges and the colophon
 * terms that used to sit under the masthead. Read-only; each fact is edited
 * in Settings / Newsletters.
 */
export function RepositoryMetaColumn({
  repository,
}: {
  repository: Repository;
}) {
  const daysToSend = daysUntilNextSend(
    repository.sendDay,
    repository.sendHour,
    new Date(),
  );

  return (
    <div className="flex w-full flex-col items-end gap-5">
      <div className="flex flex-wrap items-center justify-end gap-2">
        <Badge variant={repository.newsletterEnabled ? "success" : "secondary"}>
          <Icon type="mail" className="text-[10px]" />
          {repository.newsletterEnabled ? "Newsletters On" : "Newsletters Off"}
        </Badge>
        <RepositoryStatusBadge
          status={repository.status}
          indexedAt={repository.indexedAt}
          indexError={repository.indexError}
        />
      </div>

      <dl className="flex w-full flex-col gap-4">
        <Fact label="Branch">{repository.defaultBranch}</Fact>

        <Fact label="Next issue">
          {repository.newsletterEnabled ? (
            <Tooltip
              value={`${formatLocalSendWindow(
                repository.sendDay,
                repository.sendHour,
                getBrowserTimezone(),
              )} · ${
                daysToSend === 0
                  ? "today"
                  : `in ${daysToSend} day${daysToSend === 1 ? "" : "s"}`
              }`}
            >
              <span className="cursor-help">
                {formatSendWindow(repository.sendDay, repository.sendHour)}
              </span>
            </Tooltip>
          ) : (
            <Unset>Paused</Unset>
          )}
        </Fact>

        <Fact label="Recipients">{describeRecipients(repository)}</Fact>

        <Fact label="Last issue">
          {repository.lastNewsletterSentAt ? (
            <Tooltip
              value={formatDateTime(repository.lastNewsletterSentAt) ?? ""}
            >
              <span className="cursor-help">
                {formatDaysAgo(repository.lastNewsletterSentAt)}
              </span>
            </Tooltip>
          ) : (
            <Unset>Never</Unset>
          )}
        </Fact>

        <Fact label="Indexed">
          {repository.indexedAt ? (
            <Tooltip value={formatDateTime(repository.indexedAt) ?? ""}>
              <span className="cursor-help">
                {formatDaysAgo(repository.indexedAt)}
              </span>
            </Tooltip>
          ) : (
            <Unset>Never</Unset>
          )}
        </Fact>
      </dl>
    </div>
  );
}
