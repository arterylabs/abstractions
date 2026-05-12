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
import { useRouter } from "next/navigation";
import { FormEvent, useEffect, useState } from "react";
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

  const handleSignupWithGoogle = () => {
    window.location.href = `${process.env.NEXT_PUBLIC_API_URL}/auth/google`;
  };

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
        title: "Account created successfully",
        message: "Logging you in...",
      });
      api.utils.setToken(
        signupResponse.result.token,
        signupResponse.result.expiresAt,
      );
      router.push("/dashboard");
    }
    if (signupResponse?.error) {
      Notifier.error({
        title: "Let's try that again.",
        message: signupResponse?.error.message,
      });
    }
  }, [signupResponse.result, signupResponse.error]);

  return (
    <div className="flex flex-col gap-6">
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-xl">Create Account</CardTitle>
          <CardDescription>Sign up with your Google account</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid gap-6">
            <div className="flex flex-col gap-4">
              <Button
                variant="outline"
                className="w-full"
                onClick={handleSignupWithGoogle}
                style={{ fontFamily: "'Roboto', sans-serif" }}
              >
                <svg
                  className="h-5 w-5"
                  viewBox="0 0 100 100"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <g clipPath="url(#clip0_101_21)">
                    <rect width="100" height="100" fill="white" />
                    <path
                      fillRule="evenodd"
                      clipRule="evenodd"
                      d="M98 51.1364C98 47.5909 97.6818 44.1818 97.0909 40.9091H50V60.25H76.9091C75.75 66.5 72.2273 71.7955 66.9318 75.3409V87.8864H83.0909C92.5455 79.1818 98 66.3637 98 51.1364Z"
                      fill="#4285F4"
                    />
                    <path
                      fillRule="evenodd"
                      clipRule="evenodd"
                      d="M50 100C63.5 100 74.8182 95.5227 83.0909 87.8864L66.9318 75.3409C62.4546 78.3409 56.7273 80.1136 50 80.1136C36.9773 80.1136 25.9546 71.3182 22.0228 59.5H5.31821V72.4545C13.5455 88.7955 30.4546 100 50 100Z"
                      fill="#34A853"
                    />
                    <path
                      fillRule="evenodd"
                      clipRule="evenodd"
                      d="M22.0227 59.5C21.0227 56.5 20.4545 53.2955 20.4545 50C20.4545 46.7046 21.0227 43.5 22.0227 40.5V27.5455H5.31818C1.93182 34.2955 0 41.9318 0 50C0 58.0682 1.93182 65.7046 5.31818 72.4546L22.0227 59.5Z"
                      fill="#FBBC05"
                    />
                    <path
                      fillRule="evenodd"
                      clipRule="evenodd"
                      d="M50 19.8864C57.3409 19.8864 63.9318 22.4091 69.1137 27.3636L83.4546 13.0227C74.7955 4.95455 63.4773 0 50 0C30.4546 0 13.5455 11.2045 5.31821 27.5455L22.0228 40.5C25.9546 28.6818 36.9773 19.8864 50 19.8864Z"
                      fill="#EA4335"
                    />
                  </g>
                  <defs>
                    <clipPath id="clip0_101_21">
                      <rect width="100" height="100" fill="white" />
                    </clipPath>
                  </defs>
                </svg>
                Login with Google
              </Button>
            </div>
            <div className="relative text-center text-sm after:absolute after:inset-0 after:top-1/2 after:z-0 after:flex after:items-center after:border-t after:border-border">
              <span className="relative z-10 bg-background px-2 text-muted-foreground">
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
      <div className="text-balance text-center text-xs text-muted-foreground">
        By clicking continue, you agree to our{" "}
        <Link href="/terms">Terms of Service</Link> and{" "}
        <Link href="/privacy">Privacy Policy</Link>.
      </div>
    </div>
  );
}
