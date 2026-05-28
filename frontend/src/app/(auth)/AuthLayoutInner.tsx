"use client";

import Logo from "@/components/ui/logo";
import { store } from "@/store";
import Link from "next/link";
import { Provider } from "react-redux";

export default function AuthLayoutInner({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <Provider store={store}>
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
    </Provider>
  );
}
