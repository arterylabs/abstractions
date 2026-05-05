"use client";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Notifier } from "@/components/ui/sonner";
import { FormEvent, useEffect, useRef, useState } from "react";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { updateCurrentUser, uploadUserAvatar } from "@/store/modules/auth";
import { JSONObject } from "@/types/user";

export default function AccountSettings() {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [isUploading, setIsUploading] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
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
        title: "Profile updated",
        message: "Your profile has been updated successfully",
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

  const handleAvatarClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setIsUploading(true);
      try {
        const result = await dispatch(uploadUserAvatar(file));
        if ((result as JSONObject).error) {
          throw new Error(result.payload as string);
        }
        Notifier.success({
          title: "Image Updated",
          message: "Your avatar has been updated successfully",
        });
      } catch (error: unknown) {
        Notifier.error({
          title: "Failed to update avatar",
          message: (error as Error).message,
        });
      } finally {
        setIsUploading(false);
      }
    }
  };

  if (!currentUser) return null;

  return (
    <div className="space-y-8">
      <div>
        <form onSubmit={handleProfileUpdate} className="space-y-6">
          <div className="grid columns-3 grid-cols-3">
            <div className="col-span-2">
              <div className="space-y-4 max-w-md mb-4">
                <div>
                  <Input
                    id="name"
                    label="Name"
                    value={name}
                    onChange={(v) => setName(v)}
                  />
                </div>
                <div>
                  <Input
                    id="email"
                    type="email"
                    label="Email"
                    value={email}
                    disabled
                    onChange={(v) => setEmail(v)}
                  />
                </div>
              </div>
              <Button type="submit" loading={isSaving}>
                Save Changes
              </Button>
            </div>
            <div className="col-span-1 border-l">
              <div className="flex flex-col items-center gap-4">
                <Avatar
                  className="h-20 w-20 cursor-pointer"
                  onClick={handleAvatarClick}
                >
                  <AvatarImage src={currentUser.avatar} />
                  <AvatarFallback>{name.charAt(0)}</AvatarFallback>
                </Avatar>
                <input
                  type="file"
                  ref={fileInputRef}
                  className="hidden"
                  accept="image/*"
                  onChange={handleFileChange}
                />
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={handleAvatarClick}
                  loading={isUploading}
                >
                  {isUploading
                    ? "Uploading..."
                    : currentUser.avatar
                      ? "Change Avatar"
                      : "Upload Avatar"}
                </Button>
              </div>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
}
