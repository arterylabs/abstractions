"use client";

import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { SimpleDialog } from "@/components/ui/dialog";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Notifier } from "@/components/ui/sonner";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { updateRepositoryDescription } from "@/store/modules/repositories";
import { Repository } from "@/types";

interface UpdateDescriptionModalProps {
  repository: Repository;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function UpdateDescriptionModal({
  repository,
  open,
  onOpenChange,
}: UpdateDescriptionModalProps) {
  const dispatch = useAppDispatch();
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );
  const [description, setDescription] = useState(repository.description);
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    if (open) {
      setDescription(repository.description);
    }
  }, [open, repository.description]);

  const handleSave = async () => {
    if (!currentWorkspace?.id) return;
    setIsSaving(true);
    try {
      await dispatch(
        updateRepositoryDescription({
          workspaceId: currentWorkspace.id,
          repoId: repository.id,
          description,
        }),
      ).unwrap();
      onOpenChange(false);
    } catch {
      Notifier.error({
        title: "Failed to update description",
        message: "Please try again.",
      });
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <SimpleDialog
      open={open}
      onOpenChange={onOpenChange}
      title={repository.description ? "Edit Description" : "Add Description"}
      description={`Add a short description for ${repository.fullName}.`}
      footer={
        <>
          <Button
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={isSaving}
          >
            Cancel
          </Button>
          <Button onClick={() => void handleSave()} loading={isSaving}>
            Save
          </Button>
        </>
      }
    >
      <div className="flex flex-col gap-2">
        <Label htmlFor="repo-description">Description</Label>
        <Textarea
          id="repo-description"
          placeholder="Describe what this repository is about…"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          rows={3}
          maxLength={1000}
        />
      </div>
    </SimpleDialog>
  );
}
