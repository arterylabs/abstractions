"use client";

import api from "@/api";
import { Notifier } from "@/components/ui/sonner";
import { useRouter } from "next/navigation";

/** Returns a logout handler that clears the auth token and redirects to home. */
export function useLogout() {
  const router = useRouter();
  return () => {
    api.utils.removeToken();
    Notifier.success({
      title: "Logged out",
      message: "You have been logged out successfully",
    });
    router.push("/");
  };
}
