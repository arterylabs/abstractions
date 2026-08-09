"use client";

import api from "@/api";
import useAPI from "@/api/hooks";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { GitHubOAuthButton } from "@/components/ui/github-oauth-button";
import { GoogleOAuthButton } from "@/components/ui/google-oauth-button";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Notifier } from "@/components/ui/sonner";
import { getOAuthErrorMessage } from "@/lib/auth-errors";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { FormEvent, useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";
import * as yup from "yup";

const signupRequestSchema = yup.object().shape({
  email: yup.string().email().required(),
  password: yup.string().min(8).required(),
  name: yup.string().required(),
});

export default function Signup() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [name, setName] = useState("");
  const router = useRouter();
  const searchParams = useSearchParams();

  const handleSignupWithGoogle = () => {
    window.location.href = `${process.env.NEXT_PUBLIC_API_URL}/auth/google?page=signup`;
  };

  const handleSignupWithGitHub = () => {
    window.location.href = `${process.env.NEXT_PUBLIC_API_URL}/auth/github?page=signup`;
  };

  useEffect(() => {
    const error = searchParams.get("error");
    const errorCode = searchParams.get("error_code");

    if (!error) return;

    Notifier.error({
      title: "Sign up failed",
      message: getOAuthErrorMessage(errorCode),
    });

    // Remove error params from URL so the toast doesn't re-trigger on refresh
    const params = new URLSearchParams(searchParams.toString());
    params.delete("error");
    params.delete("error_code");
    const newPath = params.toString()
      ? `/signup?${params.toString()}`
      : "/signup";
    router.replace(newPath);
  }, []);

  const [signupResponse, sendSignupRequest] = useAPI(() =>
    api.auth.signup({ email, password, name }),
  );

  const handleSignup = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    e.stopPropagation();

    try {
      await signupRequestSchema.validate({ email, password, name });
    } catch (error: unknown) {
      Notifier.error({
        title: "Error",
        message: (error as Error).message,
      });
      return;
    }

    try {
      await sendSignupRequest();
    } catch (error: unknown) {
      Notifier.error({
        title: "Error creating your account",
        message: (error as Error).message,
      });
    }
  };

  useEffect(() => {
    if (signupResponse?.result) {
      Notifier.success({
        title: "Account created",
        message: "Signing you in now...",
      });
      api.utils.setToken(
        signupResponse.result.token,
        signupResponse.result.expiresAt,
      );
      router.push("/dashboard");
    }
    if (signupResponse?.error) {
      Notifier.error({
        title: "Error signing up",
        message: signupResponse?.error.message,
      });
    }
  }, [signupResponse.result, signupResponse.error]);

  return (
    <div className="flex flex-col gap-6">
      <Card>
        <CardHeader className="items-center gap-2 border-b border-border p-6 text-center">
          <CardTitle className="press-display text-3xl">
            Create Account
          </CardTitle>
          <CardDescription>
            Sign up with your Google or GitHub account
          </CardDescription>
        </CardHeader>
        <CardContent className="p-6">
          <div className="grid gap-6">
            <div className="flex flex-col gap-4">
              <GoogleOAuthButton
                onClick={handleSignupWithGoogle}
                label="Sign up with Google"
              />
              <GitHubOAuthButton
                onClick={handleSignupWithGitHub}
                label="Sign up with GitHub"
              />
            </div>
            <div className="relative text-center after:absolute after:inset-0 after:top-1/2 after:z-0 after:flex after:items-center after:border-t after:border-border">
              <span className="press-label relative z-10 bg-card px-3 text-muted-foreground">
                Or continue with
              </span>
            </div>
            <form onSubmit={handleSignup}>
              <div className="grid gap-6">
                <Input
                  label="Name"
                  id="name"
                  type="text"
                  value={name}
                  onChange={(v) => setName(v)}
                  placeholder="Gordon Clark"
                  required
                  data-testid="signup-name-input"
                />
                <Input
                  label="Email"
                  id="email"
                  type="email"
                  value={email}
                  onChange={(v) => setEmail(v)}
                  placeholder="gordon@cometlist.net"
                  required
                  data-testid="signup-email-input"
                />
                <Input
                  label="Password"
                  id="password"
                  type="password"
                  value={password}
                  onChange={(v) => setPassword(v)}
                  required
                  toggleVisibility
                  data-testid="signup-password-input"
                />
                <Button
                  type="submit"
                  className="w-full"
                  loading={signupResponse.isFetching}
                  data-testid="signup-submit-button"
                >
                  Create Account
                </Button>
              </div>
            </form>
            <div className="text-center text-sm">
              Already have an account? <Link href="/login">Login</Link>
            </div>
          </div>
        </CardContent>
      </Card>
      <div className="text-center text-xs leading-relaxed text-balance text-muted-foreground">
        By clicking continue, you agree to our{" "}
        <Link href="/terms">Terms of Service</Link> and{" "}
        <Link href="/privacy">Privacy Policy</Link>.
      </div>
    </div>
  );
}
