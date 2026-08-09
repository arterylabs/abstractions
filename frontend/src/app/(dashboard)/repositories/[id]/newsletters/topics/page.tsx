"use client";

import { useEffect, useState, useCallback } from "react";
import { useParams } from "next/navigation";
import { useAppSelector } from "@/store/hooks";
import api from "@/api";
import { type NewsletterTopic, type ReorderDirection } from "@/api/newsletters";
import useAPI from "@/api/hooks";
import { TopicCard } from "@/components/topics/TopicCard";
import { TopicModal } from "@/components/topics/TopicModal";
import { Button } from "@/components/ui/button";
import { askConfirmation } from "@/components/ui/confirmation";
import { EmptyState } from "@/components/ui/empty-state";
import Icon from "@/components/ui/icon";
import { PageHeader } from "@/components/ui/page-header";
import { SectionMark } from "@/components/ui/section-mark";
import Spinner from "@/components/ui/spinner";
import { Notifier } from "@/components/ui/sonner";
import { getErrorMessage } from "@/lib/errors";

export default function RepositoryTopicsPage() {
  const params = useParams<{ id: string }>();
  const repoId = params.id;
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );

  const [topicsResponse, loadTopics] = useAPI(api.newsletters.topics.list);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingTopic, setEditingTopic] = useState<NewsletterTopic | undefined>(
    undefined,
  );

  const refresh = useCallback(() => {
    if (!currentWorkspace?.id) return;
    void loadTopics(currentWorkspace.id, repoId);
  }, [loadTopics, repoId, currentWorkspace?.id]);

  useEffect(() => {
    if (!repoId || !currentWorkspace?.id) return;
    void loadTopics(currentWorkspace.id, repoId);
  }, [loadTopics, repoId, currentWorkspace?.id]);

  const topics = topicsResponse.result ?? [];

  const sentTopics = topics.filter((t) => t.sentAt !== null);
  const pendingTopics = topics.filter((t) => t.sentAt === null);

  const lastSent = sentTopics.length > 0 ? sentTopics[0] : null;
  const nextUp = pendingTopics.length > 0 ? pendingTopics[0] : null;
  const comingUp = pendingTopics.slice(1);

  const handleReorderTopic = async (
    topic: NewsletterTopic,
    direction: ReorderDirection,
  ) => {
    try {
      await api.newsletters.topics.reorder(
        currentWorkspace!.id,
        repoId,
        topic.id,
        direction,
      );
      refresh();
    } catch (error) {
      Notifier.error({
        title: "Failed to reorder topic",
        message: getErrorMessage(error),
      });
    }
  };

  const handleDelete = (topic: NewsletterTopic) => {
    askConfirmation({
      title: "Delete topic",
      content: `Are you sure you want to delete "${topic.title}"? This cannot be undone.`,
      variant: "destructive",
      okText: "Delete",
      okCallback: async () => {
        await api.newsletters.topics.delete(
          currentWorkspace!.id,
          repoId,
          topic.id,
        );
        Notifier.success({ title: "Topic deleted" });
        refresh();
      },
    });
  };

  const handleEdit = (topic: NewsletterTopic) => {
    setEditingTopic(topic);
    setIsModalOpen(true);
  };

  const handleCreate = () => {
    setEditingTopic(undefined);
    setIsModalOpen(true);
  };

  if (topicsResponse.isFetching && topics.length === 0) {
    return (
      <div className="flex min-h-[240px] flex-col items-center justify-center gap-4">
        <Spinner />
        <p className="press-label text-muted-foreground">Setting copy…</p>
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-8">
      <PageHeader
        size="h2"
        rule={false}
        title="Topics"
        description="What goes in your next newsletter. Topics are generated from your codebase — you can reorder, edit, or add your own."
      >
        <Button size="sm" onClick={handleCreate}>
          <Icon type="add" className="!text-base" />
          Add topic
        </Button>
      </PageHeader>

      {topics.length === 0 ? (
        <EmptyState
          icon="topic"
          title="Your queue is empty"
          description="Topics are generated automatically before each newsletter send. You can also add your own."
          cta={
            <Button size="sm" onClick={handleCreate}>
              <Icon type="add" className="!text-base" />
              Add topic
            </Button>
          }
        />
      ) : (
        <div className="flex flex-col gap-8">
          {lastSent && (
            <section className="flex flex-col gap-4">
              <SectionMark index="00" label="Last sent" />
              <TopicCard
                topic={lastSent}
                isFirst={true}
                isLast={true}
                isSent={true}
                onMoveUp={() => void handleReorderTopic(lastSent, "up")}
                onMoveDown={() => void handleReorderTopic(lastSent, "down")}
                onMoveToTop={() => void handleReorderTopic(lastSent, "top")}
                onEdit={() => handleEdit(lastSent)}
                onDelete={() => handleDelete(lastSent)}
              />
            </section>
          )}

          {nextUp && (
            <section className="flex flex-col gap-4">
              <SectionMark index="01" label="Up next" />
              <TopicCard
                topic={nextUp}
                isFirst={true}
                isLast={comingUp.length === 0}
                isSent={false}
                onMoveUp={() => void handleReorderTopic(nextUp, "up")}
                onMoveDown={() => void handleReorderTopic(nextUp, "down")}
                onMoveToTop={() => void handleReorderTopic(nextUp, "top")}
                onEdit={() => handleEdit(nextUp)}
                onDelete={() => handleDelete(nextUp)}
              />
            </section>
          )}

          {comingUp.length > 0 && (
            <section className="flex flex-col gap-4">
              <SectionMark index="02" label="Coming up" />
              <div className="flex flex-col gap-4">
                {comingUp.map((topic, idx) => (
                  <TopicCard
                    key={topic.id}
                    topic={topic}
                    isFirst={idx === 0}
                    isLast={idx === comingUp.length - 1}
                    isSent={false}
                    onMoveUp={() => void handleReorderTopic(topic, "up")}
                    onMoveDown={() => void handleReorderTopic(topic, "down")}
                    onMoveToTop={() => void handleReorderTopic(topic, "top")}
                    onEdit={() => handleEdit(topic)}
                    onDelete={() => handleDelete(topic)}
                  />
                ))}
              </div>
            </section>
          )}
        </div>
      )}

      <TopicModal
        open={isModalOpen}
        onOpenChange={setIsModalOpen}
        workspaceId={currentWorkspace?.id ?? ""}
        repoId={repoId}
        topic={editingTopic}
        onSaved={refresh}
      />
    </div>
  );
}
