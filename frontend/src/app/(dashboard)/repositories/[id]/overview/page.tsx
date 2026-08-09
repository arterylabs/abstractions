"use client";

import { useEffect } from "react";
import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import api from "@/api";
import useAPI from "@/api/hooks";
import { useNewsletterPolling } from "@/hooks/use-newsletter-polling";
import { Badge } from "@/components/ui/badge";
import { SectionMark } from "@/components/ui/section-mark";
import Spinner from "@/components/ui/spinner";
import { StatusBadge } from "@/components/ui/status-badge";
import { formatDaysAgo } from "@/lib/format";
import {
  describeRecipients,
  formatLocalSendWindow,
  formatSendWindow,
} from "@/lib/repository";
import { getBrowserTimezone } from "@/lib/timezone";
import { useAppSelector } from "@/store/hooks";

/* The overview is the front page: it reports the state of the repository and
   sends you to the section that owns each thing. Nothing here is editable —
   topics, delivery, and issues are all edited under Newsletters. */

/** A quiet in-panel placeholder. The dashed EmptyState is reserved for whole
    pages; inside a panel it would fight the panel's own rule. */
function PanelEmpty({
  children,
  href,
  action,
}: {
  children: React.ReactNode;
  href: string;
  action: string;
}) {
  return (
    <div className="flex flex-col items-start gap-3 py-2">
      <p className="text-sm leading-relaxed text-muted-foreground">
        {children}
      </p>
      <Link href={href} className="press-label text-primary hover:underline">
        {action} →
      </Link>
    </div>
  );
}

function Panel({
  label,
  href,
  action,
  children,
}: {
  label: string;
  /** Where this panel's content is actually managed. */
  href?: string;
  action?: string;
  children: React.ReactNode;
}) {
  return (
    <section className="flex flex-col gap-4">
      <SectionMark
        label={label}
        action={
          href && action ? (
            <Link
              href={href}
              className="press-label text-muted-foreground transition-colors hover:text-primary"
            >
              {action} →
            </Link>
          ) : undefined
        }
      />
      <div className="flex-1 border border-border bg-card p-5">{children}</div>
    </section>
  );
}

function Loading({ what }: { what: string }) {
  return (
    <div className="flex items-center gap-2.5 py-2">
      <Spinner className="h-4 w-4" />
      <span className="press-label text-muted-foreground">{what}</span>
    </div>
  );
}

/** A labelled line inside a panel — the terms of delivery, set as a list. */
function Term({
  label,
  children,
}: {
  label: string;
  children: React.ReactNode;
}) {
  return (
    <div className="flex flex-col gap-1">
      <dt className="press-label text-muted-foreground/60">{label}</dt>
      <dd className="text-sm text-foreground">{children}</dd>
    </div>
  );
}

export default function RepositoryOverviewPage() {
  const params = useParams<{ id: string }>();
  const router = useRouter();
  const repoId = params.id;

  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );
  const repository = useAppSelector((state) =>
    state.repositories.items.find((r) => r.id === repoId),
  );

  const [newslettersResponse] = useNewsletterPolling(
    currentWorkspace?.id,
    repoId,
  );
  const [topicsResponse, loadTopics] = useAPI(api.newsletters.topics.list);
  const [statsResponse, loadStats] = useAPI(api.repositories.getStats);

  useEffect(() => {
    if (!repoId || !currentWorkspace?.id) return;
    void loadTopics(currentWorkspace.id, repoId);
    void loadStats(currentWorkspace.id, repoId);
  }, [repoId, currentWorkspace?.id, loadTopics, loadStats]);

  const newsletters = newslettersResponse.result ?? [];
  const topics = topicsResponse.result ?? [];
  const stats = statsResponse.result;

  const pendingTopics = topics.filter((t) => !t.sentAt);
  const queuedPreview = pendingTopics.slice(0, 5);
  const sentNewsletters = newsletters.filter((n) => n.sentAt).slice(0, 5);

  return (
    <div className="flex flex-col gap-9">
      <div className="grid grid-cols-1 items-stretch gap-9 lg:grid-cols-[1.4fr_1fr]">
        <Panel
          label="Up next"
          href={`/repositories/${repoId}/newsletters/topics`}
          action="Manage topics"
        >
          {topicsResponse.isFetching && topics.length === 0 ? (
            <Loading what="Setting copy…" />
          ) : pendingTopics.length === 0 ? (
            <PanelEmpty
              href={`/repositories/${repoId}/newsletters/topics`}
              action="Queue a topic"
            >
              Nothing is queued for the next issue yet. Add a topic and it will
              go out with the next run.
            </PanelEmpty>
          ) : (
            <div className="flex flex-col gap-3">
              <p className="press-label text-muted-foreground">
                {pendingTopics.length}{" "}
                {pendingTopics.length === 1 ? "topic" : "topics"} queued
              </p>
              <ul className="flex flex-col divide-y divide-border">
                {queuedPreview.map((topic) => (
                  <li
                    key={topic.id}
                    className="flex items-baseline gap-3 py-2.5 text-sm first:pt-0 last:pb-0"
                  >
                    <span aria-hidden className="text-[0.6rem] text-primary">
                      ◆
                    </span>
                    <span className="line-clamp-1">{topic.title}</span>
                  </li>
                ))}
              </ul>
              {pendingTopics.length > queuedPreview.length && (
                <p className="press-label text-muted-foreground/60">
                  + {pendingTopics.length - queuedPreview.length} more queued
                </p>
              )}
            </div>
          )}
        </Panel>

        <Panel
          label="Delivery"
          href={`/repositories/${repoId}/newsletters/schedule`}
          action="Change"
        >
          {repository && (
            <dl className="flex flex-col gap-4">
              <Term label="Schedule">
                {repository.newsletterEnabled ? (
                  <span className="flex flex-col gap-0.5">
                    {formatSendWindow(repository.sendDay, repository.sendHour)}
                    <span className="text-xs text-muted-foreground">
                      {formatLocalSendWindow(
                        repository.sendDay,
                        repository.sendHour,
                        getBrowserTimezone(),
                      )}
                    </span>
                  </span>
                ) : (
                  <span className="text-muted-foreground">
                    Newsletters are paused — issues only go out when you send
                    them.
                  </span>
                )}
              </Term>

              <Term label="Recipients">
                <span className="flex flex-col gap-2">
                  {describeRecipients(repository)}
                  {repository.emails.length > 0 && (
                    <span className="flex flex-wrap gap-2">
                      {repository.emails.slice(0, 4).map((email) => (
                        <Badge
                          key={email}
                          variant="outline"
                          size="sm"
                          className="normal-case tracking-normal"
                        >
                          {email}
                        </Badge>
                      ))}
                      {repository.emails.length > 4 && (
                        <Badge variant="secondary" size="sm">
                          +{repository.emails.length - 4}
                        </Badge>
                      )}
                    </span>
                  )}
                </span>
              </Term>
            </dl>
          )}
        </Panel>
      </div>

      <Panel label="Index">
        {statsResponse.isFetching && !stats ? (
          <Loading what="Reading the source…" />
        ) : !stats || stats.fileCount === 0 ? (
          <div className="flex items-center gap-3">
            {repository && <StatusBadge status={repository.status} />}
            <span className="text-sm text-muted-foreground">
              {repository?.status === "indexing"
                ? "Indexing in progress…"
                : "No files indexed yet."}
            </span>
          </div>
        ) : (
          <div className="grid grid-cols-1 gap-8 sm:grid-cols-[auto_1fr] sm:gap-12">
            <div className="flex flex-col gap-1">
              <span className="font-mono text-4xl tabular-nums text-primary">
                {stats.fileCount.toLocaleString()}
              </span>
              <span className="press-label text-muted-foreground">
                Files indexed
              </span>
            </div>

            {stats.languageBreakdown.length > 0 && (
              <div className="flex flex-col gap-2.5">
                {stats.languageBreakdown.slice(0, 6).map((lang) => (
                  <div key={lang.language} className="flex flex-col gap-1">
                    <div className="flex items-baseline justify-between gap-4 font-mono text-xs">
                      <span>{lang.language || "Unknown"}</span>
                      <span className="text-muted-foreground">
                        {lang.fileCount} · {Math.round(lang.percentage)}%
                      </span>
                    </div>
                    {/* A measure rule rather than a rounded meter */}
                    <div className="h-px w-full bg-border">
                      <div
                        className="h-px bg-primary"
                        style={{ width: `${lang.percentage}%` }}
                      />
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}
      </Panel>

      <Panel
        label="Back issues"
        href={`/repositories/${repoId}/newsletters`}
        action="All issues"
      >
        {newslettersResponse.isFetching && newsletters.length === 0 ? (
          <Loading what="Pulling the archive…" />
        ) : sentNewsletters.length === 0 ? (
          <PanelEmpty
            href={`/repositories/${repoId}/newsletters`}
            action="Go to newsletters"
          >
            Nothing has shipped yet. Once the first issue goes out it will be
            filed here.
          </PanelEmpty>
        ) : (
          <ul className="divide-y divide-border">
            {sentNewsletters.map((newsletter, i) => (
              <li key={newsletter.id}>
                <button
                  type="button"
                  className="group flex w-full items-baseline justify-between gap-4 py-3 text-left first:pt-0 last:pb-0"
                  onClick={() =>
                    router.push(
                      `/repositories/${repoId}/newsletters/${newsletter.id}`,
                    )
                  }
                >
                  <span className="flex min-w-0 items-baseline gap-3">
                    <span className="press-label shrink-0 text-muted-foreground/50">
                      {String(sentNewsletters.length - i).padStart(2, "0")}
                    </span>
                    <span className="truncate text-sm transition-colors group-hover:text-primary">
                      {newsletter.subject || "Newsletter"}
                    </span>
                  </span>
                  <span className="shrink-0 font-mono text-xs text-muted-foreground">
                    {formatDaysAgo(newsletter.sentAt ?? "") ?? ""}
                  </span>
                </button>
              </li>
            ))}
          </ul>
        )}
      </Panel>
    </div>
  );
}
