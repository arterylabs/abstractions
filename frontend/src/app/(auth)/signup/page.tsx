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

  const handleSignupWithGitHub = () => {
    window.location.href = `${process.env.NEXT_PUBLIC_API_URL}/auth/github`;
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
        title: "Error signing up",
        message: signupResponse?.error.message,
      });
    }
  }, [signupResponse.result, signupResponse.error]);

  return (
    <div className="flex flex-col gap-6">
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-xl">Create Account</CardTitle>
          <CardDescription>
            Sign up with your Google or GitHub account
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid gap-6">
            <div className="flex flex-col gap-4">
              <Button
                variant="outline"
                className="w-full bg-white text-[#3c4043] hover:bg-gray-50 border-[#dadce0] hover:border-[#dadce0]"
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
              <Button
                className="w-full bg-[#24292f] text-white hover:bg-[#24292f]/90"
                onClick={handleSignupWithGitHub}
              >
                <svg
                  className="h-5 w-5"
                  viewBox="0 0 98 96"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    fillRule="evenodd"
                    clipRule="evenodd"
                    d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z"
                    fill="currentColor"
                  />
                </svg>
                Sign up with GitHub
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
                />
                <Input
                  label="Email"
                  id="email"
                  type="email"
                  value={email}
                  onChange={(v) => setEmail(v)}
                  placeholder="gordon@cometlist.net"
                  required
                />
                <Input
                  label="Password"
                  id="password"
                  type="password"
                  value={password}
                  onChange={(v) => setPassword(v)}
                  required
                  toggleVisibility
                />
                <Button
                  type="submit"
                  className="w-full"
                  loading={signupResponse.isFetching}
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
