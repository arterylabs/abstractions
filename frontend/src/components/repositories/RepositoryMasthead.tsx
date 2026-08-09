"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import api from "@/api";
import { Button } from "@/components/ui/button";
import { SimpleDropdownMenu } from "@/components/ui/dropdown-menu";
import Icon from "@/components/ui/icon";
import { PageHeader } from "@/components/ui/page-header";
import { Notifier } from "@/components/ui/sonner";
import { getErrorMessage } from "@/lib/errors";
import { getRepositoryURL } from "@/lib/repository";
import { useAppDispatch } from "@/store/hooks";
import { fetchRepository } from "@/store/modules/repositories";
import { Repository } from "@/types";
import { SendNewsletterButton } from "./SendNewsletterButton";

/**
 * The repository masthead: who this publication is and the actions that apply
 * wherever you are inside it. Standing facts (status, schedule, recipients)
 * live in the right meta column on the repository layout.
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

  return (
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
  );
}
