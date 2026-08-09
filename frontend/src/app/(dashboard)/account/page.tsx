"use client";

import { Button } from "@/components/ui/button";
import { ImageUploadButton } from "@/components/ui/image-upload-button";
import { Input } from "@/components/ui/input";
import { Notifier } from "@/components/ui/sonner";
import { SectionMark } from "@/components/ui/section-mark";
import { FormEvent, useEffect, useState } from "react";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { updateCurrentUser, uploadUserAvatar } from "@/store/modules/auth";
import { JSONObject } from "@/types/user";

export default function AccountSettings() {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [isUploading, setIsUploading] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const dispatch = useAppDispatch();
  const currentUser = useAppSelector((state) => state.auth.currentUser);

  useEffect(() => {
    if (currentUser) {
      setName(currentUser.name);
      setEmail(currentUser.email);
    }
  }, [currentUser]);

  const handleProfileUpdate = async (e: FormEvent) => {
    e.preventDefault();
    setIsSaving(true);
    try {
      await dispatch(updateCurrentUser({ name }));
      Notifier.success({
        title: "Profile saved",
      });
    } catch (error: unknown) {
      Notifier.error({
        title: "Failed to update profile",
        message: (error as Error).message,
      });
    } finally {
      setIsSaving(false);
    }
  };

  const handleAvatarFileSelect = async (file: File) => {
    setIsUploading(true);
    try {
      const result = await dispatch(uploadUserAvatar(file));
      if ((result as JSONObject).error) {
        throw new Error(result.payload as string);
      }
      Notifier.success({
        title: "Avatar updated",
      });
    } catch (error: unknown) {
      Notifier.error({
        title: "Failed to update avatar",
        message: (error as Error).message,
      });
    } finally {
      setIsUploading(false);
    }
  };

  if (!currentUser) return null;

  return (
    <section className="flex flex-col gap-5">
      <SectionMark label="Profile" />
      <form
        onSubmit={handleProfileUpdate}
        className="grid grid-cols-1 gap-8 border border-border bg-card p-6 md:grid-cols-[1fr_auto]"
      >
        <div className="flex max-w-md flex-col gap-5">
          <Input
            id="name"
            label="Name"
            value={name}
            onChange={(v) => setName(v)}
            data-testid="account-name-input"
          />
          <Input
            id="email"
            type="email"
            label="Email"
            value={email}
            disabled
            onChange={(v) => setEmail(v)}
            data-testid="account-email-input"
          />
          <Button
            type="submit"
            loading={isSaving}
            className="w-fit"
            data-testid="account-save-button"
          >
            Save changes
          </Button>
        </div>

        <div className="border-border md:border-l md:pl-8">
          <ImageUploadButton
            imageUrl={currentUser.avatar}
            fallback={name.charAt(0)}
            onFileSelect={handleAvatarFileSelect}
            isUploading={isUploading}
            uploadLabel="Upload avatar"
            changeLabel="Change avatar"
          />
        </div>
      </form>
    </section>
  );
}
