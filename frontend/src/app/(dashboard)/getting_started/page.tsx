"use client";

import { useEffect } from "react";
import Link from "next/link";
import { useAppSelector } from "@/store/hooks";
import api from "@/api";
import useAPI from "@/api/hooks";
import { PageHeader } from "@/components/ui/page-header";
import Spinner from "@/components/ui/spinner";
import Icon from "@/components/ui/icon";

const steps = [
  {
    number: 1,
    title: "Connect your AI providers",
    description:
      "Add your Anthropic or OpenAI API keys and select the models Abstractions will use to index your codebase and generate newsletters.",
    href: "/ai",
    cta: "Configure AI",
    doneKey: "aiConfigured" as const,
  },
  {
    number: 2,
    title: "Connect your GitHub repository",
    description:
      "Link your GitHub account and choose a repository. Abstractions will index your commits and pull requests to build context for each newsletter.",
    href: "/integrations",
    cta: "Connect GitHub",
    doneKey: "githubConnected" as const,
  },
  {
    number: 3,
    title: "Send your first newsletter",
    description:
      "Generate a newsletter from your recent activity and send it to your subscribers. This is your team's changelog, automated.",
    href: "/repositories",
    cta: "Go to repositories",
    doneKey: "newsletterSent" as const,
  },
];

export default function GettingStartedPage() {
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );

  const [statusResponse, fetchStatus] = useAPI(
    api.workspaces.getGettingStarted,
  );

  useEffect(() => {
    if (currentWorkspace?.id) {
      fetchStatus(currentWorkspace.id);
    }
  }, [fetchStatus, currentWorkspace?.id]);

  const status = statusResponse.result;
  const loading = statusResponse.isFetching;

  const completedCount = status
    ? steps.filter((s) => status[s.doneKey]).length
    : 0;
  const allDone = completedCount === steps.length;

  return (
    <>
      <PageHeader
        title="Getting started"
        description="Follow these steps to set up Abstractions for your workspace."
      />

      <div className="mt-6 space-y-4">
        <div className="rounded-lg border bg-card p-6">
          <p className="text-sm text-muted-foreground leading-relaxed">
            Abstractions turns your GitHub activity into automated newsletters
            for your team and stakeholders. Connect your tools, index your
            repository, and start delivering changelogs with zero manual effort.
          </p>
          {!loading && status && (
            <p className="mt-3 text-sm font-medium">
              {allDone ? (
                <span className="text-green-600">
                  All steps complete — you&apos;re ready to go!
                </span>
              ) : (
                <span className="text-muted-foreground">
                  {completedCount} of {steps.length} steps completed
                </span>
              )}
            </p>
          )}
        </div>

        {loading && (
          <div className="flex justify-center py-8">
            <Spinner className="w-6 h-6" />
          </div>
        )}

        {!loading &&
          steps.map((step) => {
            const done = status ? status[step.doneKey] : false;

            return (
              <div
                key={step.number}
                className={`rounded-lg border p-6 transition-colors ${
                  done ? "bg-muted/40" : "bg-card"
                }`}
              >
                <div className="flex items-start gap-4">
                  <div
                    className={`flex h-8 w-8 shrink-0 items-center justify-center rounded-full text-sm font-semibold ${
                      done
                        ? "bg-green-100 text-green-700"
                        : "bg-primary/10 text-primary"
                    }`}
                  >
                    {done ? (
                      <Icon type="check" className="!text-base" />
                    ) : (
                      step.number
                    )}
                  </div>

                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 flex-wrap">
                      <h3
                        className={`text-sm font-semibold ${done ? "text-muted-foreground line-through" : ""}`}
                      >
                        {step.title}
                      </h3>
                      {done && (
                        <span className="text-xs text-green-600 font-medium">
                          Done
                        </span>
                      )}
                    </div>
                    <p className="mt-1 text-sm text-muted-foreground">
                      {step.description}
                    </p>

                    {!done && (
                      <Link
                        href={step.href}
                        className="mt-3 inline-flex items-center gap-1 text-sm font-medium text-primary hover:underline"
                      >
                        {step.cta}
                        <Icon type="arrow_forward" className="!text-sm" />
                      </Link>
                    )}
                  </div>
                </div>
              </div>
            );
          })}
      </div>
    </>
  );
}
