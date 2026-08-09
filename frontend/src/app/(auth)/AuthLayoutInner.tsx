"use client";

import api from "@/api";
import Logo from "@/components/ui/logo";
import { store } from "@/store";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { Provider } from "react-redux";

function AuthLayoutRedirect({ children }: { children: React.ReactNode }) {
  const router = useRouter();

  useEffect(() => {
    let cancelled = false;

    const redirectAuthenticatedUser = async () => {
      if (!api.utils.getToken()) {
        return;
      }

      try {
        await api.auth.getCurrentUser();
        if (!cancelled) {
          router.replace("/dashboard");
        }
      } catch {
        api.utils.removeToken();
      }
    };

    redirectAuthenticatedUser();

    return () => {
      cancelled = true;
    };
  }, [router]);

  return <>{children}</>;
}

export default function AuthLayoutInner({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <Provider store={store}>
      <AuthLayoutRedirect>
        <div className="press-grain relative flex min-h-svh flex-col items-center justify-center gap-6 overflow-hidden bg-background p-6 md:p-10">
          {/* Same brand wash that sits behind the landing hero */}
          <div
            aria-hidden
            className="press-wash mask-fade pointer-events-none absolute inset-x-0 -top-40 -z-10 h-[42rem]"
          />

          <div className="flex w-full max-w-sm flex-col gap-8">
            {/* Masthead — the paper's name over its dateline */}
            <div className="flex flex-col items-center gap-3">
              <Link href="/" className="flex items-center gap-2">
                <Logo width="250" />
              </Link>
              <span className="press-label text-muted-foreground/70">
                Weekly dispatch — since 2025
              </span>
            </div>
            {children}
          </div>
        </div>
      </AuthLayoutRedirect>
    </Provider>
  );
}
