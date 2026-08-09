"use client";

import api from "@/api";
import useAPI from "@/api/hooks";
import Icon from "@/components/ui/icon";
import Spinner from "@/components/ui/spinner";
import { useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";

export default function VerifyEmail() {
  const [token, setToken] = useState("");
  const [verificationComplete, setVerificationComplete] = useState(false);
  const [verificationError, setVerificationError] = useState<string | null>(
    null,
  );
  const params = useSearchParams();

  const [{ result, error }, verifyUserEmail] = useAPI(() =>
    api.auth.verifyEmail(token),
  );

  useEffect(() => {
    setToken(params.get("token") || "");
  }, [params]);

  useEffect(() => {
    if (token) {
      verifyUserEmail();
    }
  }, [token]);

  useEffect(() => {
    if (result) {
      setVerificationComplete(true);
      setVerificationError(null);
    } else if (error) {
      setVerificationError(error.message);
      setVerificationComplete(false);
    }
  }, [result, error]);

  // redirect to dashboard after verification
  useEffect(() => {
    if (verificationComplete) {
      const timer = setTimeout(() => {
        window.location.href = "/dashboard";
      }, 2000); // Redirect after 2 seconds

      return () => clearTimeout(timer); // Cleanup timer on unmount
    }
  }, [verificationComplete]);

  if (verificationComplete) {
    return (
      <div className="flex flex-col gap-4 items-center justify-center text-center">
        <Icon type="check_circle" className="!text-3xl text-primary" />
        <p className="text-sm leading-relaxed text-muted-foreground">
          Your email has been successfully verified.
          <br />
          Redirecting you to the dashboard in a moment...
        </p>
      </div>
    );
  }

  if (verificationError) {
    return (
      <div className="flex gap-2 items-center justify-center text-center">
        <Icon type="error" className="text-destructive" />
        <p className="text-sm text-destructive-dark">
          Failed to verify account: {verificationError}
        </p>
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-4 items-center justify-center text-center">
      <Spinner />
      <div className="press-label text-muted-foreground">
        Verifying your account…
      </div>
    </div>
  );
}
