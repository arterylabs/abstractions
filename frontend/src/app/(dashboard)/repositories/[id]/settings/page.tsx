"use client";

import { useEffect, useMemo, useState } from "react";
import api from "@/api";
import useAPI from "@/api/hooks";
import { Button } from "@/components/ui/button";
import { FloatingBar } from "@/components/ui/floating-bar";
import { PageHeader } from "@/components/ui/page-header";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import Spinner from "@/components/ui/spinner";
import { Notifier } from "@/components/ui/sonner";
import { getErrorMessage } from "@/lib/errors";
import { isValidEmail } from "@/lib/validation";
import {
  getDayOptions,
  isoWeekdayToName,
  localHourToUtc,
  utcHourToLocal,
} from "@/lib/timezone";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import {
  deleteRepository,
  updateRepositorySettings,
} from "@/store/modules/repositories";
import { useParams, useRouter } from "next/navigation";
import { DeleteRepositoryDialog } from "./DeleteRepositoryDialog";

const getLocalHourFromUtc = (utcHour: number) => {
  const now = new Date();
  return new Date(
    Date.UTC(now.getFullYear(), now.getMonth(), now.getDate(), utcHour),
  ).getHours();
};

const HOUR_OPTIONS = Array.from({ length: 24 }, (_, localHour) => ({
  value: String(localHour),
  label: utcHourToLocal(localHourToUtc(localHour)),
}));

export default function RepositorySettingsPage() {
  const dispatch = useAppDispatch();
  const router = useRouter();
  const params = useParams<{ id: string }>();
  const repoId = params.id;
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );

  const [repoResponse, loadRepository] = useAPI(api.repositories.get);
  const [sendDay, setSendDay] = useState("");
  const [sendHour, setSendHour] = useState("");
  const [newsletterEnabled, setNewsletterEnabled] = useState(true);
  const [customEmails, setCustomEmails] = useState<string[]>([]);
  const [emailInputValue, setEmailInputValue] = useState("");
  const [emailInputError, setEmailInputError] = useState<string | null>(null);
  const [sendToWorkspaceMembers, setSendToWorkspaceMembers] = useState(true);
  const [isSavingSchedule, setIsSavingSchedule] = useState(false);

  // Baseline values for dirty detection
  const [savedSendDay, setSavedSendDay] = useState("");
  const [savedSendHour, setSavedSendHour] = useState("");
  const [savedNewsletterEnabled, setSavedNewsletterEnabled] = useState(true);
  const [savedCustomEmails, setSavedCustomEmails] = useState<string[]>([]);
  const [savedSendToWorkspaceMembers, setSavedSendToWorkspaceMembers] =
    useState(true);
  const [isTriggeringIndex, setIsTriggeringIndex] = useState(false);
  const [isResetDialogOpen, setIsResetDialogOpen] = useState(false);
  const [confirmationValue, setConfirmationValue] = useState("");
  const [isDeleting, setIsDeleting] = useState(false);

  useEffect(() => {
    if (!currentWorkspace?.id || !repoId) {
      return;
    }

    void loadRepository(currentWorkspace.id, repoId);
  }, [currentWorkspace?.id, loadRepository, repoId]);

  useEffect(() => {
    if (!repoResponse.result) {
      return;
    }

    const repo = repoResponse.result;
    const day = String(repo.sendDay);
    const hour = String(getLocalHourFromUtc(repo.sendHour));
    const enabled = repo.newsletterEnabled;
    const emails = repo.emails ?? [];
    const sendToMembers = repo.sendToWorkspaceMembers ?? true;

    setSendDay(day);
    setSendHour(hour);
    setNewsletterEnabled(enabled);
    setCustomEmails(emails);
    setSendToWorkspaceMembers(sendToMembers);

    setSavedSendDay(day);
    setSavedSendHour(hour);
    setSavedNewsletterEnabled(enabled);
    setSavedCustomEmails(emails);
    setSavedSendToWorkspaceMembers(sendToMembers);
  }, [repoResponse.result]);

  const dayOptions = useMemo(() => getDayOptions(), []);
  const repository = repoResponse.result;

  const refreshRepository = async () => {
    if (currentWorkspace?.id) {
      await loadRepository(currentWorkspace.id, repoId);
    }
  };

  const handleTriggerIndex = async () => {
    if (!currentWorkspace?.id || !repository) {
      return;
    }

    setIsTriggeringIndex(true);

    try {
      await api.repositories.triggerIndex(currentWorkspace.id, repository.id);
      await refreshRepository();
      Notifier.success({
        title: "Indexing queued",
        message: "Repository indexing has been queued successfully.",
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

  const handleSaveSchedule = async () => {
    if (!currentWorkspace?.id || !repository || !sendDay || !sendHour) {
      return;
    }

    setIsSavingSchedule(true);

    try {
      const result = await dispatch(
        updateRepositorySettings({
          workspaceId: currentWorkspace.id,
          repoId: repository.id,
          data: {
            sendDay: Number(sendDay),
            sendHour: localHourToUtc(Number(sendHour)),
            newsletterEnabled,
            emails: customEmails,
            sendToWorkspaceMembers,
          },
        }),
      );

      if ("error" in result) {
        throw new Error(result.payload as string);
      }

      await refreshRepository();
      setSavedSendDay(sendDay);
      setSavedSendHour(sendHour);
      setSavedNewsletterEnabled(newsletterEnabled);
      setSavedCustomEmails([...customEmails]);
      setSavedSendToWorkspaceMembers(sendToWorkspaceMembers);
      Notifier.success({
        title: "Settings updated",
        message: "Repository settings saved successfully.",
      });
    } catch (error: unknown) {
      Notifier.error({
        title: "Unable to save settings",
        message: getErrorMessage(error),
      });
    } finally {
      setIsSavingSchedule(false);
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

  const isDirty =
    sendDay !== savedSendDay ||
    sendHour !== savedSendHour ||
    newsletterEnabled !== savedNewsletterEnabled ||
    JSON.stringify(customEmails) !== JSON.stringify(savedCustomEmails) ||
    sendToWorkspaceMembers !== savedSendToWorkspaceMembers;

  const handleDiscardChanges = () => {
    setSendDay(savedSendDay);
    setSendHour(savedSendHour);
    setNewsletterEnabled(savedNewsletterEnabled);
    setCustomEmails([...savedCustomEmails]);
    setSendToWorkspaceMembers(savedSendToWorkspaceMembers);
  };

  if (!repository || repoResponse.isFetching) {
    return (
      <div className="flex min-h-[240px] flex-col items-center justify-center gap-4">
        <Spinner />
        <p className="text-sm text-muted-foreground">
          Loading repository settings...
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <PageHeader
        size="h2"
        title="Settings"
        description="Manage newsletter delivery, indexing, and repository data."
      />

      <Card>
        <CardHeader>
          <CardTitle>Send Schedule</CardTitle>
          <CardDescription>
            Choose when your repository newsletter should be sent.
          </CardDescription>
        </CardHeader>
        <CardContent className="grid gap-4 md:grid-cols-2">
          <div className="grid gap-2">
            <label className="text-sm font-medium">Day</label>
            <Select value={sendDay} onValueChange={setSendDay}>
              <SelectTrigger>
                <SelectValue placeholder="Select a day" />
              </SelectTrigger>
              <SelectContent>
                {dayOptions.map((option) => (
                  <SelectItem key={option.value} value={String(option.value)}>
                    {option.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div className="grid gap-2">
            <label className="text-sm font-medium">Time</label>
            <Select value={sendHour} onValueChange={setSendHour}>
              <SelectTrigger>
                <SelectValue placeholder="Select a time" />
              </SelectTrigger>
              <SelectContent>
                {HOUR_OPTIONS.map((option) => (
                  <SelectItem key={option.value} value={option.value}>
                    {option.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div className="md:col-span-2 rounded-lg border p-4">
            <p className="text-sm text-muted-foreground">
              Currently set for {isoWeekdayToName(Number(sendDay))} at{" "}
              {sendHour
                ? HOUR_OPTIONS.find((option) => option.value === sendHour)
                    ?.label
                : "-"}
            </p>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Recipients</CardTitle>
          <CardDescription>
            Add extra email addresses to receive newsletters. Useful for group
            inboxes or external stakeholders.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
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
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Newsletters</CardTitle>
          <CardDescription>
            Enable or disable automated newsletters for this repository.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between rounded-lg border p-4">
            <div className="space-y-1">
              <p className="text-sm font-medium">Enable newsletters</p>
              <p className="text-sm text-muted-foreground">
                When enabled, newsletters will be sent according to the schedule
                above.
              </p>
            </div>
            <Switch
              checked={newsletterEnabled}
              onCheckedChange={(checked) => {
                setNewsletterEnabled(checked);
              }}
            />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Indexing</CardTitle>
          <CardDescription>
            Manually trigger a re-index of this repository.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col gap-4 rounded-lg border p-4">
            <p className="text-sm text-muted-foreground">
              Re-indexing will update the repository file index and re-process
              all files. This may take a few minutes.
            </p>
            <div>
              <Button
                variant="outline"
                onClick={() => void handleTriggerIndex()}
                loading={isTriggeringIndex}
              >
                Trigger Index
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Danger Zone</CardTitle>
          <CardDescription>
            Permanently delete this repository and all its data.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col gap-4 rounded-lg border border-destructive/30 bg-destructive/5 p-4">
            <p className="text-sm text-muted-foreground">
              Deleting the repository removes all indexed code, upcoming topics,
              and stored newsletter data. This action cannot be undone.
            </p>
            <div>
              <Button
                variant="destructive"
                onClick={() => setIsResetDialogOpen(true)}
              >
                Delete Repository
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      <DeleteRepositoryDialog
        open={isResetDialogOpen}
        onOpenChange={setIsResetDialogOpen}
        repositoryFullName={repository.fullName}
        confirmationValue={confirmationValue}
        onConfirmationChange={setConfirmationValue}
        onConfirm={() => void handleDeleteRepository()}
        isDeleting={isDeleting}
      />

      <FloatingBar
        visible={isDirty}
        content={
          <p className="text-sm text-primary-foreground/80">
            You have unsaved changes
          </p>
        }
        cancelButton={
          <Button
            variant="ghost"
            size="sm"
            className="text-primary-foreground hover:text-primary-foreground hover:bg-primary-foreground/10"
            onClick={handleDiscardChanges}
          >
            Discard
          </Button>
        }
        saveButton={
          <Button
            variant="secondary"
            size="sm"
            onClick={() => void handleSaveSchedule()}
            loading={isSavingSchedule}
          >
            Save Changes
          </Button>
        }
      />
    </div>
  );
}
