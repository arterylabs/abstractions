"use client";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import WorkspaceCreateForm from "./WorkspaceCreateForm";
import Logo from "@/components/ui/logo";
import Link from "next/link";
import { Provider } from "react-redux";
import Authenticated from "@/components/authenticated";
import { store } from "@/store";

export default function WorkspaceCreatePage() {
  return (
    <Provider store={store}>
      <Authenticated>
        <div className="flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10">
          <div className="flex w-full max-w-lg flex-col gap-6">
            <Link
              href="/"
              className="flex items-center gap-2 self-center font-medium"
            >
              <Logo width="250" />
            </Link>
            <div className="flex flex-col gap-6">
              <Card>
                <CardHeader className="text-center">
                  <CardTitle className="text-xl">
                    Create your workspace
                  </CardTitle>
                  <CardDescription>
                    This is where everything will belong
                  </CardDescription>
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
