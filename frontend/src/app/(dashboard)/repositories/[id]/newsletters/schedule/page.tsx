"use client";

import { useEffect, useMemo, useState } from "react";
import { Button } from "@/components/ui/button";
import { FloatingBar } from "@/components/ui/floating-bar";
import { PageHeader } from "@/components/ui/page-header";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Notifier } from "@/components/ui/sonner";
import { getErrorMessage } from "@/lib/errors";
import { formatHour, formatLocalSendWindow } from "@/lib/repository";
import { getBrowserTimezone, getDayOptions } from "@/lib/timezone";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { updateRepositorySettings } from "@/store/modules/repositories";
import { useParams } from "next/navigation";
import { Section, SwitchRow } from "../../settings-primitives";

const HOUR_OPTIONS = Array.from({ length: 24 }, (_, hour) => ({
  value: String(hour),
  label: formatHour(hour),
}));

export default function NewsletterSchedulePage() {
  const dispatch = useAppDispatch();
  const params = useParams<{ id: string }>();
  const repoId = params.id;
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );
  const repository = useAppSelector((state) =>
    state.repositories.items.find((r) => r.id === repoId),
  );

  const [sendDay, setSendDay] = useState("");
  const [sendHour, setSendHour] = useState("");
  const [newsletterEnabled, setNewsletterEnabled] = useState(true);
  const [isSaving, setIsSaving] = useState(false);

  // Baseline values for dirty detection
  const [savedSendDay, setSavedSendDay] = useState("");
  const [savedSendHour, setSavedSendHour] = useState("");
  const [savedNewsletterEnabled, setSavedNewsletterEnabled] = useState(true);

  useEffect(() => {
    if (!repository) {
      return;
    }

    const day = String(repository.sendDay);
    const hour = String(repository.sendHour);
    const enabled = repository.newsletterEnabled;

    setSendDay(day);
    setSendHour(hour);
    setNewsletterEnabled(enabled);

    setSavedSendDay(day);
    setSavedSendHour(hour);
    setSavedNewsletterEnabled(enabled);
  }, [repository]);

  const dayOptions = useMemo(() => getDayOptions(), []);

  const handleSave = async () => {
    if (!currentWorkspace?.id || !repository || !sendDay || !sendHour) {
      return;
    }

    setIsSaving(true);

    try {
      const result = await dispatch(
        updateRepositorySettings({
          workspaceId: currentWorkspace.id,
          repoId: repository.id,
          data: {
            sendDay: Number(sendDay),
            sendHour: Number(sendHour),
            newsletterEnabled,
          },
        }),
      );

      if ("error" in result) {
        throw new Error(result.payload as string);
      }

      setSavedSendDay(sendDay);
      setSavedSendHour(sendHour);
      setSavedNewsletterEnabled(newsletterEnabled);
      Notifier.success({
        title: "Schedule saved",
      });
    } catch (error: unknown) {
      Notifier.error({
        title: "Unable to save schedule",
        message: getErrorMessage(error),
      });
    } finally {
      setIsSaving(false);
    }
  };

  const isDirty =
    sendDay !== savedSendDay ||
    sendHour !== savedSendHour ||
    newsletterEnabled !== savedNewsletterEnabled;

  const handleDiscardChanges = () => {
    setSendDay(savedSendDay);
    setSendHour(savedSendHour);
    setNewsletterEnabled(savedNewsletterEnabled);
  };

  const localSendWindow =
    sendDay && sendHour
      ? formatLocalSendWindow(
          Number(sendDay),
          Number(sendHour),
          getBrowserTimezone(),
        )
      : null;

  if (!repository) {
    return null;
  }

  return (
    <div className="flex flex-col gap-9">
      <PageHeader
        size="h2"
        rule={false}
        title="Schedule"
        description="Control whether newsletters go out automatically, and when."
      />

      <Section
        label="Newsletters"
        description="Automated newsletters are generated from your queued topics and sent on the schedule below. You can always send an issue by hand from the masthead."
      >
        <SwitchRow
          title="Send on a schedule"
          description="When off, nothing goes out automatically and the schedule is ignored."
          checked={newsletterEnabled}
          onCheckedChange={setNewsletterEnabled}
        />
      </Section>

      <Section
        label="Schedule"
        description="Times are stored in UTC so they hold steady across daylight saving changes."
      >
        <div className="grid gap-4 md:grid-cols-2">
          <div className="grid gap-2">
            <label className="press-label text-muted-foreground">Day</label>
            <Select
              value={sendDay}
              onValueChange={setSendDay}
              disabled={!newsletterEnabled}
            >
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
            <label className="press-label text-muted-foreground">
              Time (UTC)
            </label>
            <Select
              value={sendHour}
              onValueChange={setSendHour}
              disabled={!newsletterEnabled}
            >
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
        </div>
        <div className="border border-border p-4">
          {newsletterEnabled ? (
            <p className="text-sm text-muted-foreground">
              Issues go out every{" "}
              <span className="text-foreground">
                {dayOptions.find((o) => String(o.value) === sendDay)?.label}
              </span>{" "}
              at{" "}
              <span className="text-foreground">
                {HOUR_OPTIONS.find((o) => o.value === sendHour)?.label ?? "—"}{" "}
                UTC
              </span>
              {localSendWindow && <> — {localSendWindow} for you.</>}
            </p>
          ) : (
            <p className="text-sm text-muted-foreground">
              Scheduled sending is off. Turn it back on above to use this
              schedule.
            </p>
          )}
        </div>
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
