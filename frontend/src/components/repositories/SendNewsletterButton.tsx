"use client";

import { useState } from "react";
import { isAxiosError } from "axios";
import api from "@/api";
import { Button } from "@/components/ui/button";
import { SimpleDialog } from "@/components/ui/dialog";
import Icon from "@/components/ui/icon";
import { Notifier } from "@/components/ui/sonner";
import { getErrorMessage } from "@/lib/errors";
import { Repository } from "@/types";

// The API rate-limits a repository to one newsletter per day, so a send inside
// this window is confirmed before it is spent.
const RECENT_SEND_WINDOW_MS = 24 * 60 * 60 * 1000;

/**
 * The repository's primary action. It lives in the masthead so it is reachable
 * from every tab rather than only from the newsletter archive.
 */
export function SendNewsletterButton({
  workspaceId,
  repository,
  onSent,
}: {
  workspaceId: string;
  repository: Repository;
  /** Called after a send is accepted, so the caller can refresh its lists. */
  onSent?: () => void | Promise<void>;
}) {
  const [isSending, setIsSending] = useState(false);
  const [isConfirmOpen, setIsConfirmOpen] = useState(false);

  const send = async () => {
    setIsSending(true);

    try {
      await api.newsletters.triggerNewsletter(workspaceId, repository.id);
      await onSent?.();
      Notifier.info({
        title: "Newsletter is being generated",
        message: "Your newsletter is being generated and will be sent shortly.",
      });
    } catch (error: unknown) {
      const isRateLimited =
        isAxiosError(error) && error.response?.status === 429;
      Notifier.error({
        title: "Unable to send newsletter",
        message: isRateLimited
          ? "Newsletter already sent in the last 24 hours."
          : getErrorMessage(error),
      });
    } finally {
      setIsSending(false);
    }
  };

  const handleClick = () => {
    const sentAt = repository.lastNewsletterSentAt;
    const sentRecently =
      sentAt && Date.now() - new Date(sentAt).getTime() < RECENT_SEND_WINDOW_MS;

    if (sentRecently) {
      setIsConfirmOpen(true);
      return;
    }

    void send();
  };

  return (
    <>
      <Button size="sm" onClick={handleClick} loading={isSending}>
        <Icon type="mail" className="!text-base" />
        Send newsletter
      </Button>

      <SimpleDialog
        open={isConfirmOpen}
        onOpenChange={setIsConfirmOpen}
        title="Send another newsletter?"
        description="You've already sent a newsletter in the last 24 hours. Do you want to send another one?"
        footer={
          <>
            <Button variant="outline" onClick={() => setIsConfirmOpen(false)}>
              Cancel
            </Button>
            <Button
              onClick={() => {
                setIsConfirmOpen(false);
                void send();
              }}
              loading={isSending}
            >
              Send Anyway
            </Button>
          </>
        }
      >
        <></>
      </SimpleDialog>
    </>
  );
}
