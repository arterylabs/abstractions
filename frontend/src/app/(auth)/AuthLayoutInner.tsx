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
    if (api.utils.getToken()) {
      router.replace("/dashboard");
    }
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
        <div className="flex min-h-svh flex-col items-center justify-center gap-6 bg-background p-6 md:p-10">
          <div className="flex w-full max-w-sm flex-col gap-6">
            <Link
              href="/"
              className="flex items-center gap-2 self-center font-medium"
            >
              <Logo width="250" />
            </Link>
            {children}
          </div>
        </div>
      </AuthLayoutRedirect>
    </Provider>
  );
}
