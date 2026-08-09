"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { PageHeader } from "@/components/ui/page-header";
import { SectionMark } from "@/components/ui/section-mark";

import { Notifier } from "@/components/ui/sonner";
import { getErrorMessage } from "@/lib/errors";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import {
  deleteRepository,
  fetchRepository,
} from "@/store/modules/repositories";
import api from "@/api";
import { useParams, useRouter } from "next/navigation";
import { DeleteRepositoryDialog } from "./DeleteRepositoryDialog";
import { Section } from "../settings-primitives";

/* Repository-level settings only — newsletter delivery (schedule,
   recipients) lives under the Newsletters tab, next to the topics and
   issues it governs. */
export default function RepositorySettingsPage() {
  const dispatch = useAppDispatch();
  const router = useRouter();
  const params = useParams<{ id: string }>();
  const repoId = params.id;
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );
  const repository = useAppSelector((state) =>
    state.repositories.items.find((r) => r.id === repoId),
  );

  const [isTriggeringIndex, setIsTriggeringIndex] = useState(false);
  const [isResetDialogOpen, setIsResetDialogOpen] = useState(false);
  const [confirmationValue, setConfirmationValue] = useState("");
  const [isDeleting, setIsDeleting] = useState(false);

  const handleTriggerIndex = async () => {
    if (!currentWorkspace?.id || !repository) {
      return;
    }

    setIsTriggeringIndex(true);

    try {
      await api.repositories.triggerIndex(currentWorkspace.id, repository.id);
      await dispatch(
        fetchRepository({
          workspaceId: currentWorkspace.id,
          repoId: repository.id,
        }),
      );
      Notifier.success({
        title: "Indexing queued",
      });
    } catch (error: unknown) {
      Notifier.error({
        title: "Unable to trigger index",
        message: getErrorMessage(error),
      });
    } finally {
      setIsTriggeringIndex(false);
    }
  };

  const handleDeleteRepository = async () => {
    if (!currentWorkspace?.id || !repository) {
      return;
    }

    setIsDeleting(true);

    try {
      const result = await dispatch(
        deleteRepository({
          workspaceId: currentWorkspace.id,
          repoId: repository.id,
        }),
      );

      if ("error" in result) {
        throw new Error(result.payload as string);
      }

      Notifier.success({
        title: "Repository deleted",
        message: "The repository and all its data have been removed.",
      });
      router.push("/repositories");
    } catch (error: unknown) {
      Notifier.error({
        title: "Unable to delete repository",
        message: getErrorMessage(error),
      });
    } finally {
      setIsDeleting(false);
    }
  };

  if (!repository) {
    return null;
  }

  return (
    <div className="flex flex-col gap-9">
      <PageHeader
        size="h2"
        rule={false}
        title="Settings"
        description="Control how this repository is indexed, or remove it entirely."
      />

      <Section
        label="Indexing"
        description="Abstractions reads your default branch to work out what changed and what is worth writing about. Re-indexing refreshes that picture and re-processes every file, which can take a few minutes."
        action={
          <Button
            variant="outline"
            size="sm"
            onClick={() => void handleTriggerIndex()}
            loading={isTriggeringIndex}
          >
            Re-index now
          </Button>
        }
      >
        <dl className="grid grid-cols-2 gap-x-6 gap-y-4 sm:grid-cols-3">
          <div className="flex flex-col gap-1">
            <dt className="press-label text-muted-foreground/60">Branch</dt>
            <dd className="font-mono text-xs">{repository.defaultBranch}</dd>
          </div>
          <div className="flex flex-col gap-1">
            <dt className="press-label text-muted-foreground/60">Status</dt>
            <dd className="font-mono text-xs capitalize">
              {repository.status}
            </dd>
          </div>
          <div className="flex flex-col gap-1">
            <dt className="press-label text-muted-foreground/60">
              Last indexed
            </dt>
            <dd className="font-mono text-xs">
              {repository.indexedAt ? (
                new Date(repository.indexedAt).toLocaleString()
              ) : (
                <span className="text-muted-foreground">Never</span>
              )}
            </dd>
          </div>
        </dl>
      </Section>

      <section className="flex flex-col gap-4">
        <SectionMark label="Danger zone" />
        <div className="flex flex-col gap-4 border border-destructive/30 bg-destructive/5 p-5">
          <p className="max-w-[68ch] text-sm leading-relaxed text-muted-foreground">
            Deleting the repository removes all indexed code, upcoming topics,
            and stored newsletter data. This action cannot be undone.
          </p>
          <div>
            <Button
              variant="destructive"
              onClick={() => setIsResetDialogOpen(true)}
            >
              Delete repository
            </Button>
          </div>
        </div>
      </section>

      <DeleteRepositoryDialog
        open={isResetDialogOpen}
        onOpenChange={setIsResetDialogOpen}
        repositoryFullName={repository.fullName}
        confirmationValue={confirmationValue}
        onConfirmationChange={setConfirmationValue}
        onConfirm={() => void handleDeleteRepository()}
        isDeleting={isDeleting}
      />
    </div>
  );
}
