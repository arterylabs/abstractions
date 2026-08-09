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
import { Input } from "@/components/ui/input";
import { Notifier } from "@/components/ui/sonner";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { FormEvent, useEffect, useState } from "react";
import * as yup from "yup";

const resetPasswordSchema = yup.object({
  password: yup.string().min(8).required(),
  passwordConfirmation: yup
    .string()
    .oneOf([yup.ref("password")], "Passwords must match")
    .required(),
});

export default function ResetPassword() {
  const [password, setPassword] = useState("");
  const [passwordConfirmation, setPasswordConfirmation] = useState("");
  const router = useRouter();

  const queryParams = useSearchParams();

  const token = queryParams.get("token");
  const id = queryParams.get("id");

  const [resetPasswordResponse, sendResetPasswordRequest] = useAPI(() =>
    api.auth.resetPasswordConfirm(token || "", password, id || ""),
  );

  const handlePasswordReset = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    e.stopPropagation();

    try {
      await resetPasswordSchema.validate({ password, passwordConfirmation });
    } catch (error: unknown) {
      Notifier.error({
        title: "Validation Error",
        message: (error as yup.ValidationError).message,
      });
      return;
    }

    try {
      await sendResetPasswordRequest();
    } catch (error: unknown) {
      Notifier.error({
        title: "Error Resetting Your Password",
        message: (error as Error).message,
      });
    }

    // signupRequest.run(() => api.auth.signup({ email, password, name }));
  };

  useEffect(() => {
    if (resetPasswordResponse?.result) {
      Notifier.success({
        title: "Password Reset Successful",
        message: "Login again using your new password",
      });
      router.push("/login");
    }
  }, [resetPasswordResponse.result]);

  return (
    <div className="flex flex-col gap-6">
      <Card>
        <CardHeader className="items-center gap-2 border-b border-border p-6 text-center">
          <CardTitle className="press-display text-3xl">
            Reset Password
          </CardTitle>
          <CardDescription>Set your new password here</CardDescription>
        </CardHeader>
        <CardContent className="p-6">
          <div className="grid gap-6">
            <form onSubmit={handlePasswordReset}>
              <div className="grid gap-6">
                <Input
                  label="Password"
                  id="password"
                  type="password"
                  value={password}
                  onChange={(v) => setPassword(v)}
                  required
                  toggleVisibility
                />
                <Input
                  label="Confirm Password"
                  id="confirm_password"
                  type="password"
                  value={passwordConfirmation}
                  onChange={(v) => setPasswordConfirmation(v)}
                  required
                  toggleVisibility
                />
                <Button
                  type="submit"
                  className="w-full"
                  loading={resetPasswordResponse.isFetching}
                >
                  Reset Password
                </Button>
              </div>
            </form>
            <div className="text-center text-sm">
              Just remembered your password? <Link href="/login">Login</Link>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
