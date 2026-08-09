"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import api from "@/api";
import { type NewsletterTopic } from "@/api/newsletters";
import { type RepositoryStatus } from "@/api/repositories";
import useAPI from "@/api/hooks";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import Spinner from "@/components/ui/spinner";
import Icon from "@/components/ui/icon";
import { Notifier } from "@/components/ui/sonner";
import { useAppSelector } from "@/store/hooks";
import { getErrorMessage } from "@/lib/errors";

// How often to poll during indexing and topic generation phases.
const POLL_INTERVAL_MS = 5000;

type WizardStep = "indexing" | "topics" | "newsletter";

function StepIndicator({ step }: { step: WizardStep }) {
  const steps: { key: WizardStep; label: string }[] = [
    { key: "indexing", label: "Indexing" },
    { key: "topics", label: "Generating topics" },
    { key: "newsletter", label: "First newsletter" },
  ];

  const currentIndex = steps.findIndex((s) => s.key === step);

  return (
    // Steps set as folios along a rule: done, current (inked), still to come.
    <ol className="flex flex-wrap items-center justify-center gap-x-4 gap-y-2">
      {steps.map((s, idx) => {
        const isDone = idx < currentIndex;
        const isCurrent = idx === currentIndex;

        return (
          <li key={s.key} className="flex items-center gap-4">
            {idx > 0 && <span aria-hidden className="h-px w-6 bg-border" />}
            <span
              className={[
                "press-label flex items-center gap-1.5",
                isCurrent
                  ? "text-primary"
                  : isDone
                    ? "text-muted-foreground"
                    : "text-muted-foreground/40",
              ].join(" ")}
            >
              {isDone ? (
                <Icon type="check" className="!text-sm text-success-dark" />
              ) : (
                <span>§ {String(idx + 1).padStart(2, "0")}</span>
              )}
              {s.label}
            </span>
          </li>
        );
      })}
    </ol>
  );
}

function IndexingStep({
  status,
  repoFullName,
}: {
  status: RepositoryStatus | null;
  repoFullName: string;
}) {
  const totalFiles = status?.totalFiles ?? 0;
  const summarizedFiles = status?.summarizedFiles ?? 0;
  const progressPercentage = status?.progressPercentage ?? 0;
  const languageBreakdown = status?.languageBreakdown ?? [];
  const topLanguage = languageBreakdown[0];

  return (
    <div className="flex flex-col gap-8">
      <div className="flex flex-col items-center gap-3 text-center">
        <div className="flex items-center justify-center">
          <Spinner className="h-10 w-10" />
        </div>
        <h2 className="text-3xl font-medium">Indexing {repoFullName}</h2>
        <p className="max-w-[52ch] text-sm leading-relaxed text-muted-foreground">
          We&apos;re reading your codebase to understand what&apos;s been built.
          This usually takes a few minutes.
        </p>
      </div>

      <Card>
        <CardContent className="space-y-4 p-6">
          <div className="flex items-center justify-between font-mono text-xs">
            <span className="text-muted-foreground">
              {totalFiles > 0
                ? `${summarizedFiles.toLocaleString()} / ${totalFiles.toLocaleString()} files summarised`
                : "Waiting for files to be fetched…"}
            </span>
            <span className="tabular-nums text-foreground">
              {Math.round(progressPercentage)}%
            </span>
          </div>

          <Progress value={progressPercentage} aria-label="Indexing progress" />

          {topLanguage && (
            <p className="font-mono text-xs text-muted-foreground">
              Top language:{" "}
              <span className="font-medium text-foreground">
                {topLanguage.language}
              </span>{" "}
              — {topLanguage.fileCount} files (
              {Math.round(topLanguage.percentage)}%)
              {languageBreakdown.length > 1 &&
                ` · ${languageBreakdown.length} languages detected`}
            </p>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

function TopicsStep({ repoFullName }: { repoFullName: string }) {
  return (
    <div className="flex flex-col gap-8">
      <div className="flex flex-col items-center gap-3 text-center">
        <div className="flex items-center justify-center">
          <Spinner className="h-10 w-10" />
        </div>
        <h2 className="text-3xl font-medium">Generating newsletter topics</h2>
        <p className="max-w-[52ch] text-sm leading-relaxed text-muted-foreground">
          Analysing {repoFullName} and identifying the most interesting
          engineering updates for your team…
        </p>
      </div>

      <Card>
        <CardContent className="flex items-center gap-4 p-6">
          <div className="flex h-10 w-10 shrink-0 items-center justify-center border border-primary/30">
            <Icon type="auto_awesome" className="!text-lg text-primary" />
          </div>
          <div className="space-y-1">
            <p className="press-label text-foreground">
              Building your first newsletter queue
            </p>
            <p className="text-sm text-muted-foreground">
              We&apos;re picking out recent changes, features, and fixes worth
              sharing. This takes about a minute.
            </p>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

function NewsletterStep({
  topic,
  isSending,
  onSend,
  onSkip,
}: {
  topic: NewsletterTopic;
  isSending: boolean;
  onSend: () => void;
  onSkip: () => void;
}) {
  return (
    <div className="flex flex-col gap-8">
      <div className="flex flex-col items-center gap-3 text-center">
        <div className="flex h-12 w-12 items-center justify-center border border-success-dark/40 text-success-dark">
          <Icon type="check_circle" className="!text-2xl" />
        </div>
        <h2 className="text-3xl font-medium">
          Your codebase is ready — send your first newsletter
        </h2>
        <p className="max-w-[52ch] text-sm leading-relaxed text-muted-foreground">
          We found topics to share with your team. Here&apos;s what your first
          newsletter will cover:
        </p>
      </div>

      <Card className="border-primary/25 bg-accent-dark">
        <CardContent className="space-y-3 p-6">
          <div className="flex items-start gap-3">
            <div className="mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center border border-primary/30">
              <Icon type="topic" className="!text-base text-primary" />
            </div>
            <div className="space-y-1">
              <h3 className="text-2xl font-medium">{topic.title}</h3>
              {topic.description && (
                <p className="text-sm text-muted-foreground line-clamp-3">
                  {topic.description}
                </p>
              )}
            </div>
          </div>
        </CardContent>
      </Card>

      <div className="flex flex-col gap-3 sm:flex-row sm:justify-center">
        <Button
          size="lg"
          onClick={onSend}
          loading={isSending}
          className="gap-2"
        >
          <Icon type="mail_outline" className="!text-base" />
          Send first newsletter
        </Button>
        <Button
          size="lg"
          variant="outline"
          onClick={onSkip}
          disabled={isSending}
        >
          Skip for now
        </Button>
      </div>

      <p className="press-label text-center text-muted-foreground/70">
        You can always send a newsletter later from the repository page.
      </p>
    </div>
  );
}

export default function RepositoryConnectWizardPage() {
  const params = useParams<{ id: string }>();
  const router = useRouter();
  const repoId = params.id;

  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );

  // Load the repository so we have the full name
  const [repositoryResponse, loadRepository] = useAPI(api.repositories.get);
  const [statusResponse, loadStatus] = useAPI(api.repositories.getStatus);
  const [topicsResponse, loadTopics] = useAPI(api.newsletters.topics.list);

  const [step, setStep] = useState<WizardStep>("indexing");
  const [isSending, setIsSending] = useState(false);
  const [firstTopic, setFirstTopic] = useState<NewsletterTopic | null>(null);

  // Track whether we've already auto-advanced to avoid double transitions
  const advancedToTopicsRef = useRef(false);
  const advancedToNewsletterRef = useRef(false);

  // Load repository details on mount
  useEffect(() => {
    if (!currentWorkspace?.id || !repoId) return;
    void loadRepository(currentWorkspace.id, repoId);
  }, [currentWorkspace?.id, loadRepository, repoId]);

  // Step 1: poll indexing status
  const pollStatus = useCallback(() => {
    if (!currentWorkspace?.id || !repoId) return;
    void loadStatus(currentWorkspace.id, repoId);
  }, [currentWorkspace?.id, loadStatus, repoId]);

  useEffect(() => {
    if (step !== "indexing") return;
    pollStatus();
    const id = window.setInterval(pollStatus, POLL_INTERVAL_MS);
    return () => window.clearInterval(id);
  }, [pollStatus, step]);

  // Auto-advance from indexing → topics once indexed
  useEffect(() => {
    const liveStatus = statusResponse.result?.status;
    if (
      liveStatus === "indexed" &&
      step === "indexing" &&
      !advancedToTopicsRef.current
    ) {
      advancedToTopicsRef.current = true;
      setStep("topics");
    }
    // If indexing errored, still advance so the user isn't stuck
    if (
      liveStatus === "error" &&
      step === "indexing" &&
      !advancedToTopicsRef.current
    ) {
      advancedToTopicsRef.current = true;
      setStep("topics");
    }
  }, [statusResponse.result?.status, step]);

  // Step 2: poll topics
  const pollTopics = useCallback(() => {
    if (!currentWorkspace?.id || !repoId) return;
    void loadTopics(currentWorkspace.id, repoId);
  }, [currentWorkspace?.id, loadTopics, repoId]);

  useEffect(() => {
    if (step !== "topics") return;
    pollTopics();
    const id = window.setInterval(pollTopics, POLL_INTERVAL_MS);
    return () => window.clearInterval(id);
  }, [pollTopics, step]);

  // Auto-advance from topics → newsletter once topics exist
  useEffect(() => {
    const topics = topicsResponse.result;
    if (
      topics &&
      topics.length > 0 &&
      step === "topics" &&
      !advancedToNewsletterRef.current
    ) {
      advancedToNewsletterRef.current = true;
      setFirstTopic(topics[0]);
      setStep("newsletter");
    }
  }, [topicsResponse.result, step]);

  const handleSendNewsletter = async () => {
    if (!currentWorkspace?.id) return;
    setIsSending(true);
    try {
      await api.newsletters.triggerNewsletter(currentWorkspace.id, repoId);
      Notifier.info({
        title: "Newsletter is being generated",
        message:
          "Your first newsletter is being generated and will be sent shortly.",
      });
      router.push(`/repositories/${repoId}/newsletters`);
    } catch (error: unknown) {
      Notifier.error({
        title: "Unable to send newsletter",
        message: getErrorMessage(error),
      });
      setIsSending(false);
    }
  };

  const handleSkip = () => {
    router.push(`/repositories/${repoId}/overview`);
  };

  const repository = repositoryResponse.result;
  const repoFullName = repository?.fullName ?? "your repository";

  return (
    <div className="mx-auto flex max-w-xl flex-col gap-10 py-8">
      {/* Step indicator */}
      <div className="flex flex-col items-center gap-4">
        <StepIndicator step={step} />
      </div>

      {/* Step content */}
      {step === "indexing" && (
        <IndexingStep
          status={statusResponse.result ?? null}
          repoFullName={repoFullName}
        />
      )}

      {step === "topics" && <TopicsStep repoFullName={repoFullName} />}

      {step === "newsletter" && firstTopic && (
        <NewsletterStep
          topic={firstTopic}
          isSending={isSending}
          onSend={() => void handleSendNewsletter()}
          onSkip={handleSkip}
        />
      )}

      {/* Fallback: topics step but still no topics after initial load */}
      {step === "newsletter" && !firstTopic && (
        <div className="flex flex-col items-center gap-4 text-center">
          <p className="text-sm text-muted-foreground">
            No topics were generated yet. You can send a newsletter once topics
            are ready.
          </p>
          <Button variant="outline" onClick={handleSkip}>
            Go to repository
          </Button>
        </div>
      )}
    </div>
  );
}
