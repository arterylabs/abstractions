"use client";

import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { FloatingBar } from "@/components/ui/floating-bar";
import { Input } from "@/components/ui/input";
import { PageHeader } from "@/components/ui/page-header";
import { Notifier } from "@/components/ui/sonner";
import { getErrorMessage } from "@/lib/errors";
import { isValidEmail } from "@/lib/validation";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { updateRepositorySettings } from "@/store/modules/repositories";
import { useParams } from "next/navigation";
import { Section, SwitchRow } from "../../settings-primitives";

export default function NewsletterRecipientsPage() {
  const dispatch = useAppDispatch();
  const params = useParams<{ id: string }>();
  const repoId = params.id;
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );
  const repository = useAppSelector((state) =>
    state.repositories.items.find((r) => r.id === repoId),
  );

  const [customEmails, setCustomEmails] = useState<string[]>([]);
  const [emailInputValue, setEmailInputValue] = useState("");
  const [emailInputError, setEmailInputError] = useState<string | null>(null);
  const [sendToWorkspaceMembers, setSendToWorkspaceMembers] = useState(true);
  const [isSaving, setIsSaving] = useState(false);

  // Baseline values for dirty detection
  const [savedCustomEmails, setSavedCustomEmails] = useState<string[]>([]);
  const [savedSendToWorkspaceMembers, setSavedSendToWorkspaceMembers] =
    useState(true);

  useEffect(() => {
    if (!repository) {
      return;
    }

    const emails = repository.emails ?? [];
    const sendToMembers = repository.sendToWorkspaceMembers ?? true;

    setCustomEmails(emails);
    setSendToWorkspaceMembers(sendToMembers);

    setSavedCustomEmails(emails);
    setSavedSendToWorkspaceMembers(sendToMembers);
  }, [repository]);

  const handleSave = async () => {
    if (!currentWorkspace?.id || !repository) {
      return;
    }

    setIsSaving(true);

    try {
      const result = await dispatch(
        updateRepositorySettings({
          workspaceId: currentWorkspace.id,
          repoId: repository.id,
          data: {
            sendDay: repository.sendDay,
            sendHour: repository.sendHour,
            newsletterEnabled: repository.newsletterEnabled,
            emails: customEmails,
            sendToWorkspaceMembers,
          },
        }),
      );

      if ("error" in result) {
        throw new Error(result.payload as string);
      }

      setSavedCustomEmails([...customEmails]);
      setSavedSendToWorkspaceMembers(sendToWorkspaceMembers);
      Notifier.success({
        title: "Recipients saved",
      });
    } catch (error: unknown) {
      Notifier.error({
        title: "Unable to save recipients",
        message: getErrorMessage(error),
      });
    } finally {
      setIsSaving(false);
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
    JSON.stringify(customEmails) !== JSON.stringify(savedCustomEmails) ||
    sendToWorkspaceMembers !== savedSendToWorkspaceMembers;

  const handleDiscardChanges = () => {
    setCustomEmails([...savedCustomEmails]);
    setSendToWorkspaceMembers(savedSendToWorkspaceMembers);
  };

  if (!repository) {
    return null;
  }

  return (
    <div className="flex flex-col gap-9">
      <PageHeader
        size="h2"
        rule={false}
        title="Recipients"
        description="Control who receives each issue."
      />

      <Section
        label="Recipients"
        description="Every issue goes to your workspace members by default. Add extra addresses — teammates outside the workspace, group inboxes, or anyone who should stay in the loop."
      >
        <Input
          id="repository-recipient-email"
          type="email"
          label="Add email"
          labelClassName="press-label text-muted-foreground"
          inputClassName="font-mono"
          placeholder="group@example.com"
          value={emailInputValue}
          onChange={(value) => {
            setEmailInputValue(value);
            setEmailInputError(null);
          }}
          onKeyDown={handleEmailInputKeyDown}
        />
        <p className="-mt-3 text-xs text-muted-foreground">
          Press Enter or comma to add an email.
        </p>
        {emailInputError && (
          <p className="-mt-3 text-xs text-destructive">{emailInputError}</p>
        )}

        {customEmails.length > 0 && (
          <div className="flex flex-wrap gap-2">
            {customEmails.map((email) => (
              <span
                key={email}
                className="inline-flex items-center gap-1 rounded-sm border border-border bg-accent px-2.5 py-1 font-mono text-xs text-foreground"
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
          <SwitchRow
            title="Send to workspace members"
            description="When disabled, newsletters are sent only to the custom emails above."
            checked={sendToWorkspaceMembers}
            onCheckedChange={setSendToWorkspaceMembers}
          />
        )}
      </Section>

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
            onClick={() => void handleSave()}
            loading={isSaving}
          >
            Save Changes
          </Button>
        }
      />
    </div>
  );
}
