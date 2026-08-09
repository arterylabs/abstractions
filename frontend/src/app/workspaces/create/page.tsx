"use client";

import { Card, CardContent, CardHeader } from "@/components/ui/card";
import WorkspaceCreateForm from "./WorkspaceCreateForm";
import Logo from "@/components/ui/logo";
import Link from "next/link";
import { Provider } from "react-redux";
import Authenticated from "@/components/authenticated";
import { store } from "@/store";
import { PageHeader } from "@/components/ui/page-header";
import Icon from "@/components/ui/icon";

export default function WorkspaceCreatePage() {
  return (
    <Provider store={store}>
      <Authenticated>
        <div className="press-grain relative flex min-h-svh flex-col items-center justify-center gap-6 overflow-hidden p-6 md:p-10">
          <div
            aria-hidden
            className="press-wash mask-fade pointer-events-none absolute inset-x-0 -top-40 -z-10 h-[42rem]"
          />
          <div className="flex w-full max-w-3xl flex-col gap-6">
            <div className="relative flex w-full items-center justify-center">
              <Link
                href="/dashboard"
                className="press-label absolute left-0 flex items-center gap-1.5 text-muted-foreground transition-colors hover:text-foreground"
              >
                <Icon type="arrow_back" className="!text-sm" />
                Back
              </Link>
              <Link href="/" className="flex items-center gap-2 font-medium">
                <Logo width="250" />
              </Link>
            </div>
            <div className="flex flex-col gap-6">
              <Card>
                <CardHeader className="border-b border-border">
                  <PageHeader
                    title="Create your workspace"
                    description="Your workspace is where your team, repositories, and newsletters live."
                    size="h2"
                    eyebrow="New workspace"
                    rule={false}
                  />
                </CardHeader>
                <CardContent>
                  <WorkspaceCreateForm />
                </CardContent>
              </Card>
            </div>
          </div>
        </div>
      </Authenticated>
    </Provider>
  );
}
