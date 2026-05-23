"use client";

import api from "@/api";
import Spinner from "@/components/ui/spinner";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect } from "react";

export default function Authenticated() {
  const queryParams = useSearchParams();

  const refPath = queryParams.get("ref");
  const router = useRouter();

  useEffect(() => {
    const code = queryParams.get("code") || "";
    if (!code) {
      router.push("/login?error=true");
      return;
    }

    api.auth
      .exchangeCode(code)
      .then((res) => {
        const { token, expiresAt } = res.data;
        api.utils.setToken(token, expiresAt);
        if (refPath) {
          router.push(refPath);
        } else {
          window.location.href = "/dashboard";
        }
      })
      .catch(() => {
        router.push("/login?error=true");
      });
  }, []);

  return (
    <div className="flex flex-col gap-4 items-center justify-center">
      <Spinner />
      <div className="opacity-60">Signing in...</div>
    </div>
  );
}
