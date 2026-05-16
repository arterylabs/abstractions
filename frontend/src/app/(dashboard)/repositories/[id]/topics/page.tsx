"use client";

import { useEffect, useState, useCallback } from "react";
import { useParams } from "next/navigation";
import api from "@/api";
import { type NewsletterTopic } from "@/api/newsletters";
import useAPI from "@/api/hooks";
import { TopicCard } from "@/components/topics/TopicCard";
import { TopicModal } from "@/components/topics/TopicModal";
import { Button } from "@/components/ui/button";
import { askConfirmation } from "@/components/ui/confirmation";
import { EmptyState } from "@/components/ui/empty-state";
import Icon from "@/components/ui/icon";
import Spinner from "@/components/ui/spinner";
import { Notifier } from "@/components/ui/sonner";
import { getErrorMessage } from "@/lib/errors";

export default function RepositoryTopicsPage() {
  const params = useParams<{ id: string }>();
  const repoId = params.id;

  const [topicsResponse, loadTopics] = useAPI(api.newsletters.topics.list);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingTopic, setEditingTopic] = useState<NewsletterTopic | undefined>(
    undefined,
  );

  const refresh = useCallback(() => {
    void loadTopics(repoId);
  }, [loadTopics, repoId]);

  useEffect(() => {
    if (!repoId) return;
    void loadTopics(repoId);
  }, [loadTopics, repoId]);

  const topics = topicsResponse.result ?? [];

  const sentTopics = topics.filter((t) => t.sentAt !== null);
  const pendingTopics = topics.filter((t) => t.sentAt === null);

  const lastSent = sentTopics.length > 0 ? sentTopics[0] : null;
  const nextUp = pendingTopics.length > 0 ? pendingTopics[0] : null;
  const comingUp = pendingTopics.slice(1);

  const handleMoveUp = async (topic: NewsletterTopic) => {
    try {
      await api.newsletters.topics.reorder(repoId, topic.id, "up");
      refresh();
    } catch (error) {
      Notifier.error({
        title: "Failed to reorder topic",
        message: getErrorMessage(error),
      });
    }
  };

  const handleMoveDown = async (topic: NewsletterTopic) => {
    try {
      await api.newsletters.topics.reorder(repoId, topic.id, "down");
      refresh();
    } catch (error) {
      Notifier.error({
        title: "Failed to reorder topic",
        message: getErrorMessage(error),
      });
    }
  };

  const handleMoveToTop = async (topic: NewsletterTopic) => {
    try {
      await api.newsletters.topics.reorder(repoId, topic.id, "top");
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
        await api.newsletters.topics.delete(repoId, topic.id);
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
        <p className="text-sm text-muted-foreground">Loading topics...</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-base font-semibold">Topics</h2>
          <p className="text-sm text-muted-foreground mt-0.5">
            Manage the topics queue for your newsletter.
          </p>
        </div>
        <Button size="sm" onClick={handleCreate}>
          <Icon type="add" className="!text-base" />
          Add topic
        </Button>
      </div>

      {topics.length === 0 ? (
        <EmptyState
          title="No topics yet"
          description="Topics are generated automatically before each newsletter send. You can also add your own."
          cta={
            <Button size="sm" onClick={handleCreate}>
              <Icon type="add" className="!text-base" />
              Add topic
            </Button>
          }
        />
      ) : (
        <div className="space-y-6">
          {lastSent && (
            <section className="space-y-3">
              <h3 className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">
                Last sent
              </h3>
              <TopicCard
                topic={lastSent}
                isFirst={true}
                isLast={true}
                isSent={true}
                onMoveUp={() => void handleMoveUp(lastSent)}
                onMoveDown={() => void handleMoveDown(lastSent)}
                onMoveToTop={() => void handleMoveToTop(lastSent)}
                onEdit={() => handleEdit(lastSent)}
                onDelete={() => handleDelete(lastSent)}
              />
            </section>
          )}

          {nextUp && (
            <section className="space-y-3">
              <h3 className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">
                Up next
              </h3>
              <TopicCard
                topic={nextUp}
                isFirst={true}
                isLast={comingUp.length === 0}
                isSent={false}
                onMoveUp={() => void handleMoveUp(nextUp)}
                onMoveDown={() => void handleMoveDown(nextUp)}
                onMoveToTop={() => void handleMoveToTop(nextUp)}
                onEdit={() => handleEdit(nextUp)}
                onDelete={() => handleDelete(nextUp)}
              />
            </section>
          )}

          {comingUp.length > 0 && (
            <section className="space-y-3">
              <h3 className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">
                Coming up
              </h3>
              <div className="space-y-3">
                {comingUp.map((topic, idx) => (
                  <TopicCard
                    key={topic.id}
                    topic={topic}
                    isFirst={idx === 0}
                    isLast={idx === comingUp.length - 1}
                    isSent={false}
                    onMoveUp={() => void handleMoveUp(topic)}
                    onMoveDown={() => void handleMoveDown(topic)}
                    onMoveToTop={() => void handleMoveToTop(topic)}
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
        repoId={repoId}
        topic={editingTopic}
        onSaved={refresh}
      />
    </div>
  );
}
