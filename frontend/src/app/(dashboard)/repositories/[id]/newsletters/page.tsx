"use client";

import { useEffect, useState } from "react";
import { isAxiosError } from "axios";
import api from "@/api";
import useAPI from "@/api/hooks";
import { NewsletterCard } from "@/components/newsletters/NewsletterCard";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { SimpleDialog } from "@/components/ui/dialog";
import { EmptyState } from "@/components/ui/empty-state";
import { Switch } from "@/components/ui/switch";
import Spinner from "@/components/ui/spinner";
import { Notifier } from "@/components/ui/sonner";
import { getErrorMessage } from "@/lib/errors";
import { isValidEmail } from "@/lib/validation";
import { useParams, useRouter } from "next/navigation";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { updateRepositorySettings } from "@/store/modules/repositories";
import { localHourToUtc } from "@/lib/timezone";

const topicToLabel = (topic: unknown) => {
  if (typeof topic === "string") {
    return topic;
  }

  if (topic && typeof topic === "object") {
    const topicRecord = topic as Record<string, unknown>;
    if (typeof topicRecord.title === "string") {
      return topicRecord.title;
    }
    if (typeof topicRecord.name === "string") {
      return topicRecord.name;
    }
    return JSON.stringify(topicRecord);
  }

  return "Untitled topic";
};

export default function RepositoryNewslettersPage() {
  const dispatch = useAppDispatch();
  const router = useRouter();
  const params = useParams<{ id: string }>();
  const repoId = params.id;
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );
  const repositoryFromStore = useAppSelector((state) =>
    state.repositories.items.find((r) => r.id === repoId),
  );

  const [newslettersResponse, loadNewsletters] = useAPI(api.newsletters.list);
  const [topicsResponse, loadTopics] = useAPI(api.newsletters.topics.list);
  const [isGeneratingDraft, setIsGeneratingDraft] = useState(false);
  const [isTriggeringNewsletter, setIsTriggeringNewsletter] = useState(false);

  // Recipients modal state
  const [isRecipientsModalOpen, setIsRecipientsModalOpen] = useState(false);
  const [customEmails, setCustomEmails] = useState<string[]>([]);
  const [emailInputValue, setEmailInputValue] = useState("");
  const [emailInputError, setEmailInputError] = useState<string | null>(null);
  const [sendToWorkspaceMembers, setSendToWorkspaceMembers] = useState(true);
  const [isSavingRecipients, setIsSavingRecipients] = useState(false);

  useEffect(() => {
    if (!repoId) {
      return;
    }

    void loadNewsletters(repoId);
    void loadTopics(repoId);
  }, [loadNewsletters, loadTopics, repoId]);

  // Sync recipients state from store whenever the repository changes
  useEffect(() => {
    if (!repositoryFromStore) return;
    setCustomEmails(repositoryFromStore.emails ?? []);
    setSendToWorkspaceMembers(
      repositoryFromStore.sendToWorkspaceMembers ?? true,
    );
  }, [repositoryFromStore]);

  const handleGenerateDraft = async () => {
    setIsGeneratingDraft(true);

    try {
      await api.newsletters.generatePreview(repoId);
      await Promise.all([loadNewsletters(repoId), loadTopics(repoId)]);
    } catch (error: unknown) {
      Notifier.error({
        title: "Unable to generate draft",
        message: getErrorMessage(error),
      });
    } finally {
      setIsGeneratingDraft(false);
    }
  };

  const handleTriggerNewsletter = async () => {
    setIsTriggeringNewsletter(true);

    try {
      await api.newsletters.triggerNewsletter(repoId);
      await Promise.all([loadNewsletters(repoId), loadTopics(repoId)]);
      Notifier.success({
        title: "Newsletter queued",
        message: "Newsletter generation has been queued successfully.",
      });
    } catch (error: unknown) {
      const isRateLimited =
        isAxiosError(error) && error.response?.status === 429;
      Notifier.error({
        title: "Unable to trigger newsletter",
        message: isRateLimited
          ? "Newsletter already sent in the last 24 hours."
          : getErrorMessage(error),
      });
    } finally {
      setIsTriggeringNewsletter(false);
    }
  };

  const handleEmailInputKeyDown = (
    e: React.KeyboardEvent<HTMLInputElement>,
  ) => {
    if (e.key !== "Enter" && e.key !== ",") return;
    e.preventDefault();
    const value = emailInputValue.trim().replace(/,$/, "");
    if (!value) return;
    if (!isValidEmail(value)) {
      setEmailInputError("Invalid email address");
      return;
    }
    if (customEmails.includes(value)) {
      setEmailInputError("Email already added");
      return;
    }
    setCustomEmails((prev) => [...prev, value]);
    setEmailInputValue("");
    setEmailInputError(null);
  };

  const handleRemoveCustomEmail = (email: string) => {
    setCustomEmails((prev) => prev.filter((e) => e !== email));
  };

  const handleSaveRecipients = async () => {
    if (!currentWorkspace?.id || !repositoryFromStore) return;
    setIsSavingRecipients(true);
    try {
      const result = await dispatch(
        updateRepositorySettings({
          workspaceId: currentWorkspace.id,
          repoId,
          data: {
            sendDay: repositoryFromStore.sendDay,
            sendHour: localHourToUtc(repositoryFromStore.sendHour),
            newsletterEnabled: repositoryFromStore.newsletterEnabled,
            emails: customEmails,
            sendToWorkspaceMembers,
          },
        }),
      );
      if ("error" in result) throw new Error(result.payload as string);
      Notifier.success({
        title: "Recipients updated",
        message: "Newsletter recipients saved successfully.",
      });
      setIsRecipientsModalOpen(false);
    } catch (error: unknown) {
      Notifier.error({
        title: "Unable to save recipients",
        message: getErrorMessage(error),
      });
    } finally {
      setIsSavingRecipients(false);
    }
  };

  const openRecipientsModal = () => {
    // Reset to current saved state when opening
    setCustomEmails(repositoryFromStore?.emails ?? []);
    setSendToWorkspaceMembers(
      repositoryFromStore?.sendToWorkspaceMembers ?? true,
    );
    setEmailInputValue("");
    setEmailInputError(null);
    setIsRecipientsModalOpen(true);
  };

  const configuredEmails = repositoryFromStore?.emails ?? [];

  return (
    <div className="space-y-6">
      <div className="space-y-2">
        <h2 className="text-xl font-medium">Newsletters</h2>
        <p className="text-sm text-muted-foreground">
          Generate drafts, send newsletters, and review repository history.
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Recipients</CardTitle>
          <CardDescription>
            Who receives newsletters for this repository.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {configuredEmails.length === 0 ? (
            <>
              <p className="text-sm text-muted-foreground">All Team Members</p>
              <div className="flex flex-wrap gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={openRecipientsModal}
                >
                  Configure Recipients
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() =>
                    router.push(`/repositories/${repoId}/settings`)
                  }
                >
                  Invite Members
                </Button>
              </div>
            </>
          ) : (
            <>
              <div className="flex flex-wrap gap-2">
                {configuredEmails.map((email) => (
                  <span
                    key={email}
                    className="inline-flex items-center gap-1 rounded-full border bg-secondary px-3 py-1 text-sm text-secondary-foreground"
                  >
                    {email}
                  </span>
                ))}
              </div>
              <Button variant="outline" size="sm" onClick={openRecipientsModal}>
                Configure Recipients
              </Button>
            </>
          )}
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="gap-3 sm:flex-row sm:items-start sm:justify-between">
          <div>
            <CardTitle>Next Newsletter</CardTitle>
            <CardDescription>
              {topicsResponse.result?.length
                ? `${topicsResponse.result.length} topic${topicsResponse.result.length === 1 ? "" : "s"} queued for the next draft.`
                : "No pending topics queued yet."}
            </CardDescription>
          </div>
          <div className="flex flex-wrap gap-2">
            <Button
              variant="outline"
              onClick={() => void handleTriggerNewsletter()}
              loading={isTriggeringNewsletter}
            >
              Send Newsletter Now
            </Button>
            <Button
              onClick={() => void handleGenerateDraft()}
              loading={isGeneratingDraft}
            >
              Generate Draft
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          {topicsResponse.isFetching && !topicsResponse.result ? (
            <div className="flex items-center gap-3 text-sm text-muted-foreground">
              <Spinner className="h-5 w-5" />
              Loading queued topics...
            </div>
          ) : topicsResponse.result?.length ? (
            <ul className="list-disc space-y-2 pl-5 text-sm text-muted-foreground">
              {topicsResponse.result.map((topic) => (
                <li key={topic.id}>
                  <span className="font-medium text-foreground">
                    {topicToLabel(topic)}
                  </span>{" "}
                  <span>{topic.sentAt ? "Sent" : "Pending"}</span>
                </li>
              ))}
            </ul>
          ) : (
            <p className="text-sm text-muted-foreground">
              No topics generated yet.
            </p>
          )}
        </CardContent>
      </Card>

      {newslettersResponse.isFetching && !newslettersResponse.result ? (
        <div className="flex min-h-[240px] flex-col items-center justify-center gap-4">
          <Spinner />
          <p className="text-sm text-muted-foreground">
            Loading newsletters...
          </p>
        </div>
      ) : newslettersResponse.result &&
        newslettersResponse.result.length > 0 ? (
        <div className="grid gap-4">
          {newslettersResponse.result.map((newsletter) => (
            <NewsletterCard
              key={newsletter.id}
              newsletter={newsletter}
              onClick={() =>
                router.push(
                  `/repositories/${repoId}/newsletters/${newsletter.id}`,
                )
              }
            />
          ))}
        </div>
      ) : (
        <EmptyState
          title="No newsletters yet"
          description="Connected repositories will show newsletter history here once your first send goes out."
          cta={
            <Button
              variant="outline"
              onClick={() => router.push("/repositories")}
            >
              Back to Repositories
            </Button>
          }
        />
      )}

      <SimpleDialog
        open={isRecipientsModalOpen}
        onOpenChange={(open) => {
          if (!open) {
            // Reset edits on close without saving
            setCustomEmails(repositoryFromStore?.emails ?? []);
            setSendToWorkspaceMembers(
              repositoryFromStore?.sendToWorkspaceMembers ?? true,
            );
            setEmailInputValue("");
            setEmailInputError(null);
          }
          setIsRecipientsModalOpen(open);
        }}
        title="Configure Recipients"
        description="Add extra email addresses to receive newsletters. Useful for group inboxes or external stakeholders."
        footer={
          <>
            <Button
              variant="outline"
              onClick={() => setIsRecipientsModalOpen(false)}
            >
              Cancel
            </Button>
            <Button
              onClick={() => void handleSaveRecipients()}
              loading={isSavingRecipients}
            >
              Save
            </Button>
          </>
        }
      >
        <div className="space-y-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Add email</label>
            <input
              type="email"
              className="flex h-9 w-full rounded-md border border-input bg-accent px-3 py-1 text-sm shadow-sm transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
              placeholder="group@example.com"
              value={emailInputValue}
              onChange={(e) => {
                setEmailInputValue(e.target.value);
                setEmailInputError(null);
              }}
              onKeyDown={handleEmailInputKeyDown}
            />
            <p className="text-xs text-muted-foreground">
              Press Enter or comma to add an email.
            </p>
            {emailInputError && (
              <p className="text-xs text-destructive">{emailInputError}</p>
            )}
          </div>

          {customEmails.length > 0 && (
            <div className="flex flex-wrap gap-2">
              {customEmails.map((email) => (
                <span
                  key={email}
                  className="inline-flex items-center gap-1 rounded-full border bg-secondary px-3 py-1 text-sm text-secondary-foreground"
                >
                  {email}
                  <button
                    type="button"
                    onClick={() => handleRemoveCustomEmail(email)}
                    className="ml-1 text-muted-foreground hover:text-foreground"
                    aria-label={`Remove ${email}`}
                  >
                    ×
                  </button>
                </span>
              ))}
            </div>
          )}

          {customEmails.length > 0 && (
            <div className="flex items-center justify-between rounded-lg border p-4">
              <div className="space-y-1">
                <p className="text-sm font-medium">Send to workspace members</p>
                <p className="text-sm text-muted-foreground">
                  When disabled, newsletters are sent only to the custom emails
                  above.
                </p>
              </div>
              <Switch
                checked={sendToWorkspaceMembers}
                onCheckedChange={setSendToWorkspaceMembers}
              />
            </div>
          )}
        </div>
      </SimpleDialog>
    </div>
  );
}
