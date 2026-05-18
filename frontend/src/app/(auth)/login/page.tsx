"use client";

import api from "@/api";
import useAPI from "@/api/hooks";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
} from "@/components/ui/card";
import { GitHubOAuthButton } from "@/components/ui/github-oauth-button";
import { GoogleOAuthButton } from "@/components/ui/google-oauth-button";
import { Input } from "@/components/ui/input";
import { Notifier } from "@/components/ui/sonner";
import { getOAuthErrorMessage } from "@/lib/auth-errors";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useSearchParams } from "next/navigation";
import { FormEvent, useEffect, useState } from "react";
import * as yup from "yup";

const loginFormSchema = yup.object({
  email: yup.string().email().required(),
  password: yup.string().required(),
});

export default function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const router = useRouter();
  const searchParams = useSearchParams();
  const refPath = searchParams.get("ref");

  const [loginResponse, sendLoginRequest] = useAPI(api.auth.login);

  const handleLoginWithGoogle = () => {
    window.location.href = `${process.env.NEXT_PUBLIC_API_URL}/auth/google?ref=${refPath}&page=login`;
  };

  const handleLoginWithGitHub = () => {
    window.location.href = `${process.env.NEXT_PUBLIC_API_URL}/auth/github?ref=${refPath}&page=login`;
  };

  useEffect(() => {
    const error = searchParams.get("error");
    const errorCode = searchParams.get("error_code");

    if (!error) return;

    Notifier.error({
      title: "Sign in failed",
      message: getOAuthErrorMessage(errorCode),
    });

    // Remove error params from URL so the toast doesn't re-trigger on refresh
    const params = new URLSearchParams(searchParams.toString());
    params.delete("error");
    params.delete("error_code");
    const newPath = params.toString()
      ? `/login?${params.toString()}`
      : "/login";
    router.replace(newPath);
  }, []);

  const handleLogin = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    e.stopPropagation();

    try {
      await loginFormSchema.validate({ email, password });
      await sendLoginRequest(email, password);
    } catch (error: unknown) {
      Notifier.error({
        title: "Error logging in",
        message: (error as Error).message,
      });
      return;
    }
  };

  useEffect(() => {
    if (loginResponse.result) {
      Notifier.success({
        title: "Login successful",
        message: "Welcome back!",
      });
      api.utils.setToken(
        loginResponse.result.token,
        loginResponse.result.expiresAt,
      );
      // @ts-expect-error this supports strings
      window.location = refPath || "/dashboard";
    }
    if (loginResponse.error) {
      Notifier.error({
        title: "Error logging in",
        message: loginResponse.error.message,
      });
    }
  }, [loginResponse.result, loginResponse.error]);

  return (
    <div className="flex flex-col gap-6">
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-xl">Sign In</CardTitle>
          <CardDescription>
            Login with your Google or GitHub account
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid gap-6">
            <div className="flex flex-col gap-4">
              <GoogleOAuthButton
                onClick={handleLoginWithGoogle}
                label="Login with Google"
              />
              <GitHubOAuthButton
                onClick={handleLoginWithGitHub}
                label="Login with GitHub"
              />
            </div>
            <div className="relative text-center text-sm after:absolute after:inset-0 after:top-1/2 after:z-0 after:flex after:items-center after:border-t after:border-border">
              <span className="relative z-10 bg-background px-2 text-muted-foreground">
                Or continue with
              </span>
            </div>
            <form onSubmit={handleLogin}>
              <div className="grid gap-6">
                <Input
                  label="Email"
                  id="email"
                  type="email"
                  value={email}
                  onChange={(v) => setEmail(v)}
                  placeholder="gordon@cometlist.net"
                  required
                  data-testid="login-email-input"
                />
                <Input
                  label="Password"
                  id="password"
                  type="password"
                  value={password}
                  onChange={(v) => setPassword(v)}
                  required
                  labelAction={
                    <Link href="/forgot" className="text-xs" tabIndex={-1}>
                      Forgot password?
                    </Link>
                  }
                  toggleVisibility
                  data-testid="login-password-input"
                />
                <Button
                  type="submit"
                  className="w-full"
                  loading={loginResponse.isFetching}
                  data-testid="login-submit-button"
                >
                  Login
                </Button>
              </div>
            </form>
            <div className="text-center text-sm">
              Don&apos;t have an account? <Link href="/signup">Sign up</Link>
            </div>
          </div>
        </CardContent>
      </Card>
      <div className="text-balance text-center text-xs text-muted-foreground">
        By clicking continue, you agree to our{" "}
        <Link href="/terms">Terms of Service</Link> and{" "}
        <Link href="/privacy">Privacy Policy</Link>.
      </div>
    </div>
  );
}
