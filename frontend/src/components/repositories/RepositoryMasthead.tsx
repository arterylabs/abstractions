"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import api from "@/api";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { SimpleDropdownMenu } from "@/components/ui/dropdown-menu";
import Icon from "@/components/ui/icon";
import { PageHeader } from "@/components/ui/page-header";
import { Notifier } from "@/components/ui/sonner";
import Tooltip from "@/components/ui/tooltip";
import { getErrorMessage } from "@/lib/errors";
import { formatDateTime, formatDaysAgo } from "@/lib/format";
import {
  describeRecipients,
  formatLocalSendWindow,
  formatSendWindow,
  getRepositoryURL,
} from "@/lib/repository";
import { daysUntilNextSend, getBrowserTimezone } from "@/lib/timezone";
import { useAppDispatch } from "@/store/hooks";
import { fetchRepository } from "@/store/modules/repositories";
import { Repository } from "@/types";
import { RepositoryStatusBadge } from "./RepositoryStatusBadge";
import { SendNewsletterButton } from "./SendNewsletterButton";

/** A single term in the colophon strip: a running head over its value. */
function Fact({
  label,
  children,
}: {
  label: string;
  children: React.ReactNode;
}) {
  return (
    <div className="flex min-w-0 flex-col gap-1">
      <dt className="press-label text-muted-foreground/60">{label}</dt>
      <dd className="truncate font-mono text-xs text-foreground">{children}</dd>
    </div>
  );
}

/** A colophon value that is not set yet, set in the same measure but quieter. */
function Unset({ children }: { children: React.ReactNode }) {
  return <span className="text-muted-foreground">{children}</span>;
}

/**
 * The repository masthead: who this publication is, the standing facts about
 * it, and the actions that apply wherever you are inside it. Every tab hangs
 * beneath this, so the facts here are read-only — each one is edited in exactly
 * one place, which the overflow menu and the Settings tab lead to.
 */
export function RepositoryMasthead({
  repository,
  workspaceId,
  onEditDescription,
}: {
  repository: Repository;
  workspaceId: string;
  onEditDescription: () => void;
}) {
  const router = useRouter();
  const dispatch = useAppDispatch();
  const [isTriggeringIndex, setIsTriggeringIndex] = useState(false);

  const handleTriggerIndex = async () => {
    setIsTriggeringIndex(true);

    try {
      await api.repositories.triggerIndex(workspaceId, repository.id);
      await dispatch(fetchRepository({ workspaceId, repoId: repository.id }));
      Notifier.success({ title: "Indexing queued" });
    } catch (error: unknown) {
      Notifier.error({
        title: "Unable to trigger index",
        message: getErrorMessage(error),
      });
    } finally {
      setIsTriggeringIndex(false);
    }
  };

  const repositoryURL = getRepositoryURL(repository);
  const daysToSend = daysUntilNextSend(
    repository.sendDay,
    repository.sendHour,
    new Date(),
  );

  return (
    <div className="flex flex-col gap-6">
      <PageHeader
        eyebrow="Repository"
        rule={false}
        display={false}
        title={
          <span className="flex items-center gap-2.5">
            <Icon
              type="deployed_code"
              className="shrink-0 !text-lg text-primary"
            />
            {repository.fullName}
          </span>
        }
        description={
          repository.description ? (
            repository.description
          ) : (
            <button
              className="press-label cursor-pointer text-muted-foreground/50 transition-colors hover:text-muted-foreground"
              onClick={onEditDescription}
            >
              Add description
            </button>
          )
        }
      >
        <div className="flex items-center gap-2">
          <SendNewsletterButton
            workspaceId={workspaceId}
            repository={repository}
          />
          <SimpleDropdownMenu
            label={repository.name ?? repository.fullName}
            options={[
              {
                label: repository.description
                  ? "Edit description"
                  : "Add description",
                icon: <Icon type="edit_note" className="!text-base" />,
                onClick: onEditDescription,
              },
              {
                label: isTriggeringIndex ? "Re-indexing…" : "Re-index now",
                icon: <Icon type="sync" className="!text-base" />,
                disabled: isTriggeringIndex,
                onClick: () => void handleTriggerIndex(),
              },
              ...(repositoryURL
                ? [
                    {
                      label: "Open on GitHub",
                      icon: <Icon type="open_in_new" className="!text-base" />,
                      onClick: () =>
                        window.open(
                          repositoryURL,
                          "_blank",
                          "noopener,noreferrer",
                        ),
                    },
                  ]
                : []),
              {
                label: "Settings",
                icon: <Icon type="settings" className="!text-base" />,
                onClick: () =>
                  router.push(`/repositories/${repository.id}/settings`),
              },
            ]}
          >
            <Button variant="outline" size="sm" aria-label="Repository actions">
              <Icon type="more_horiz" className="!text-base" />
            </Button>
          </SimpleDropdownMenu>
        </div>
      </PageHeader>

      <div className="flex flex-wrap items-center gap-2">
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

      {/* Colophon — the standing facts, set as small print between rules */}
      <dl className="grid grid-cols-2 gap-x-6 gap-y-4 border-y border-border py-4 sm:grid-cols-3 lg:grid-cols-5">
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
