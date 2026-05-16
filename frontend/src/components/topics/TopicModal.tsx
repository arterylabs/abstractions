"use client";

import { useState, useEffect } from "react";
import api from "@/api";
import { type NewsletterTopic } from "@/api/newsletters";
import { Button } from "@/components/ui/button";
import { SimpleDialog } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Notifier } from "@/components/ui/sonner";
import { getErrorMessage } from "@/lib/errors";

interface TopicModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  repoId: string;
  topic?: NewsletterTopic;
  onSaved: () => void;
}

export function TopicModal({
  open,
  onOpenChange,
  repoId,
  topic,
  onSaved,
}: TopicModalProps) {
  const isEditing = !!topic;
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    if (open) {
      setTitle(topic?.title ?? "");
      setDescription(topic?.description ?? "");
    }
  }, [open, topic]);

  const handleSave = async () => {
    if (!title.trim()) return;

    setIsSaving(true);
    try {
      if (isEditing) {
        await api.newsletters.topics.update(repoId, topic.id, {
          title: title.trim(),
          description: description.trim(),
        });
        Notifier.success({ title: "Topic updated" });
      } else {
        await api.newsletters.topics.create(repoId, {
          title: title.trim(),
          description: description.trim(),
        });
        Notifier.success({ title: "Topic added to queue" });
      }
      onSaved();
      onOpenChange(false);
    } catch (error) {
      Notifier.error({
        title: isEditing ? "Failed to update topic" : "Failed to create topic",
        message: getErrorMessage(error),
      });
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <SimpleDialog
      open={open}
      onOpenChange={onOpenChange}
      title={isEditing ? "Edit topic" : "Add topic"}
      description={
        isEditing
          ? "Update the topic details."
          : "Add a custom topic to your newsletter queue."
      }
      footer={
        <>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button
            onClick={() => void handleSave()}
            loading={isSaving}
            disabled={!title.trim()}
          >
            {isEditing ? "Save changes" : "Add topic"}
          </Button>
        </>
      }
    >
      <div className="space-y-4">
        <div className="space-y-1.5">
          <Label htmlFor="topic-title">Title</Label>
          <Input
            id="topic-title"
            value={title}
            onChange={(value) => setTitle(value)}
            placeholder="e.g. New authentication flow"
            maxLength={200}
          />
        </div>
        <div className="space-y-1.5">
          <Label htmlFor="topic-description">
            Description{" "}
            <span className="text-muted-foreground font-normal">
              (optional)
            </span>
          </Label>
          <Textarea
            id="topic-description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Briefly describe what this topic should cover..."
            rows={3}
          />
        </div>
      </div>
    </SimpleDialog>
  );
}
