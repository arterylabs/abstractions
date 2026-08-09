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
      "Add your model provider API keys (Anthropic, OpenAI, Grok, and more) and select the models Abstractions will use to index your codebase and generate newsletters.",
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
    <div className="flex flex-col gap-8">
      <PageHeader
        eyebrow="Setup"
        title="Getting started"
        description="Abstractions turns your GitHub activity into automated newsletters for your team. Connect your tools, index a repository, and start delivering changelogs with zero manual effort."
      >
        {!loading && status && (
          <span
            className={`press-label ${allDone ? "text-primary" : "text-muted-foreground"}`}
          >
            {allDone
              ? "Ready to press"
              : `${completedCount} of ${steps.length} complete`}
          </span>
        )}
      </PageHeader>

      {loading && (
        <div className="flex justify-center py-8">
          <Spinner className="h-6 w-6" />
        </div>
      )}

      {/* The press run — each step opens on its own rule, which inks in once done */}
      {!loading && (
        <ol className="flex flex-col">
          {steps.map((step) => {
            const done = status ? status[step.doneKey] : false;

            return (
              <li
                key={step.number}
                className={`group flex flex-col gap-3 border-t-2 py-7 transition-colors duration-300 ${
                  done
                    ? "border-primary/30"
                    : "border-border hover:border-primary"
                }`}
              >
                <div className="flex items-center gap-3">
                  <span className="press-label text-primary">
                    § {String(step.number).padStart(2, "0")}
                  </span>
                  {done && (
                    <span className="press-label flex items-center gap-1.5 text-muted-foreground">
                      <Icon
                        type="check"
                        className="!text-sm text-success-dark"
                      />
                      Done
                    </span>
                  )}
                </div>

                <h3
                  className={`max-w-[22ch] text-2xl font-medium ${done ? "text-muted-foreground" : ""}`}
                >
                  {step.title}
                </h3>

                <p className="max-w-[68ch] text-[0.95rem] leading-[1.75] text-muted-foreground">
                  {step.description}
                </p>

                {!done && (
                  <Link
                    href={step.href}
                    className="press-label mt-1 inline-flex w-fit items-center gap-2 text-primary hover:underline"
                  >
                    {step.cta}
                    <Icon type="arrow_forward" className="!text-sm" />
                  </Link>
                )}
              </li>
            );
          })}
        </ol>
      )}
    </div>
  );
}
