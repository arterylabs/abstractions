"use client";

import { useState } from "react";
import {
  Dialog,
  DialogPortal,
  DialogOverlay,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog";
import * as DialogPrimitive from "@radix-ui/react-dialog";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import auth from "@/api/auth";
import Icon from "./ui/icon";

const RESEND_COOLDOWN_MS = 60_000;

export default function EmailVerificationModal({
  open,
  email,
}: {
  open: boolean;
  email?: string;
}) {
  const [isSending, setIsSending] = useState(false);
  const [lastSentAt, setLastSentAt] = useState<number | null>(null);
  const [rateLimitHit, setRateLimitHit] = useState(false);

  const secondsSinceLastSend = lastSentAt
    ? Math.floor((Date.now() - lastSentAt) / 1000)
    : null;
  const cooldownRemaining =
    secondsSinceLastSend !== null
      ? Math.max(
          0,
          Math.ceil((RESEND_COOLDOWN_MS - (Date.now() - lastSentAt!)) / 1000),
        )
      : 0;
  const isOnCooldown = lastSentAt !== null && cooldownRemaining > 0;

  async function handleResend() {
    if (isOnCooldown) return;
    setIsSending(true);
    setRateLimitHit(false);
    try {
      await auth.resendVerificationEmail();
      setLastSentAt(Date.now());
    } catch (err: unknown) {
      const status = (err as { response?: { status?: number } })?.response
        ?.status;
      if (status === 429) {
        setRateLimitHit(true);
        setLastSentAt(Date.now());
      }
    } finally {
      setIsSending(false);
    }
  }

  return (
    <Dialog open={open} onOpenChange={() => {}}>
      <DialogPortal>
        <DialogOverlay />
        <DialogPrimitive.Content
          data-testid="email-verification-modal"
          onInteractOutside={(e) => e.preventDefault()}
          onEscapeKeyDown={(e) => e.preventDefault()}
          className={cn(
            "bg-background fixed top-[50%] left-[50%] z-50 grid w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] gap-4 rounded-lg border p-6 shadow-lg duration-200 sm:max-w-lg",
            "data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95",
            "data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95",
          )}
        >
          <div className="flex items-center justify-center flex-col">
            <Icon
              type="email"
              className="mx-auto !text-4xl text-primary mb-4"
            />
            <DialogTitle className=" mb-4">
              Verify your email address
            </DialogTitle>
            <DialogDescription className="mb-4 text-center">
              We&apos;ve sent a verification link to{" "}
              {email ? (
                <span className="font-medium text-foreground">{email}</span>
              ) : (
                "your email address"
              )}
              . Click the link to activate your account and unlock all features.
            </DialogDescription>
            <div className="flex flex-col gap-3 text-sm mb-4 text-center">
              {(lastSentAt || rateLimitHit) && (
                <p className="text-muted-foreground">
                  We&apos;ve already sent a verification email. Please check
                  your inbox
                  {isOnCooldown
                    ? ` — you can request another in ${cooldownRemaining}s`
                    : ""}
                  .
                </p>
              )}
              <a
                href="https://mail.google.com/mail/u/0/#search/subject%3A(verify+your+email)"
                target="_blank"
                rel="noopener noreferrer"
                className="text-primary underline underline-offset-4"
              >
                Use Gmail? Here&apos;s a quick link to find the email
              </a>
            </div>

            <DialogFooter className="mt-4">
              <Button
                data-testid="resend-verification-button"
                variant="outline"
                onClick={handleResend}
                loading={isSending}
                disabled={isOnCooldown || isSending}
              >
                {isOnCooldown
                  ? `Resend in ${cooldownRemaining}s`
                  : "Resend verification email"}
              </Button>
            </DialogFooter>
          </div>
        </DialogPrimitive.Content>
      </DialogPortal>
    </Dialog>
  );
}
