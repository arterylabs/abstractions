"use client";

import api from "@/api";
import useAPI from "@/api/hooks";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import Icon from "@/components/ui/icon";
import { Input } from "@/components/ui/input";
import { Notifier } from "@/components/ui/sonner";
import Link from "next/link";
import { FormEvent, useEffect, useState } from "react";
import * as yup from "yup";

const forgotPasswordSchema = yup.object({
  email: yup.string().email().required(),
});

export default function ForgotPassword() {
  const [email, setEmail] = useState("");

  const [resetEmailSent, setResetEmailSent] = useState(false);
  const [forgotPasswordResponse, sendForgotPasswordEmail] = useAPI(() =>
    api.auth.resetPassword(email),
  );

  const handleForgotPassword = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    e.stopPropagation();

    try {
      forgotPasswordSchema.validateSync({ email });
    } catch (error: unknown) {
      Notifier.error({
        title: "Error",
        message: (error as Error).message,
      });
      return;
    }

    try {
      await sendForgotPasswordEmail();
    } catch (error: unknown) {
      Notifier.error({
        title: "Error sending password reset email",
        message: (error as Error).message,
      });
    }
  };

  useEffect(() => {
    if (forgotPasswordResponse?.result) {
      Notifier.success("Password Reset Email Sent");
      setResetEmailSent(true);
    }
  }, [forgotPasswordResponse?.result]);

  return (
    <div className="flex flex-col gap-6">
      {resetEmailSent ? (
        <Card>
          <CardHeader className="text-center">
            <div className="text-3xl">
              <Icon type="check_circle" className="text-primary" />{" "}
            </div>
            <CardTitle className="text-xl">Check Your Inbox</CardTitle>
            <CardDescription>
              We&lsquo;ve sent an email to <strong>{email}</strong> with the
              link to reset your password
            </CardDescription>
          </CardHeader>
        </Card>
      ) : (
        <Card>
          <CardHeader className="text-center">
            <CardTitle className="text-xl">Forgot Password?</CardTitle>
            <CardDescription>
              We&lsquo;ll send you an email with the link to reset your password
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid gap-6">
              <form onSubmit={handleForgotPassword}>
                <div className="grid gap-6">
                  <Input
                    label="Email"
                    id="email"
                    type="email"
                    value={email}
                    onChange={(v) => setEmail(v)}
                    placeholder="gordon@cometlist.net"
                    required
                  />
                  <Button
                    type="submit"
                    className="w-full"
                    loading={forgotPasswordResponse.isFetching}
                  >
                    Send Password Reset Email
                  </Button>
                </div>
              </form>
              <div className="text-center text-sm">
                Just remembered? <Link href="/login">Login</Link>
              </div>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
