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
        <div className="flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10">
          <div className="flex w-full max-w-3xl flex-col gap-6">
            <div className="relative flex w-full items-center justify-center">
              <Link
                href="/dashboard"
                className="absolute left-0 flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground transition-colors"
              >
                <Icon type="arrow_back" className="text-base" />
                Back
              </Link>
              <Link href="/" className="flex items-center gap-2 font-medium">
                <Logo width="250" />
              </Link>
            </div>
            <div className="flex flex-col gap-6">
              <Card>
                <CardHeader>
                  <PageHeader
                    title="Create your workspace"
                    description="Your workspace is where your team, repositories, and newsletters live."
                    size="h2"
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
