"use client";

import { FormEvent, useEffect, useRef, useState } from "react";
import * as yup from "yup";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Notifier } from "@/components/ui/sonner";
import { Textarea } from "@/components/ui/textarea";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  uploadWorkspaceLogo,
  updateWorkspace,
} from "@/store/modules/workspace";

const workspaceSchema = yup.object().shape({
  name: yup.string().required("Workspace name is required").max(64),
  website: yup.string().url("Must be a valid URL").nullable(),
  about: yup.string().max(500, "About text is too long"),
  billingEmail: yup
    .string()
    .email("Must be a valid email")
    .required("Billing email is required"),
});

export default function WorkspaceGeneralSettings() {
  const [name, setName] = useState("");
  const [website, setWebsite] = useState("");
  const [about, setAbout] = useState("");
  const [billingEmail, setBillingEmail] = useState("");
  const [isUploading, setIsUploading] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const dispatch = useAppDispatch();
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );

  useEffect(() => {
    if (currentWorkspace) {
      setName(currentWorkspace.name);
      setWebsite(currentWorkspace.website || "");
      setAbout(currentWorkspace.about || "");
      setBillingEmail(currentWorkspace.billingEmail || "");
    }
  }, [currentWorkspace]);

  const handleWorkspaceUpdate = async (e: FormEvent) => {
    e.preventDefault();
    setIsSaving(true);

    try {
      await workspaceSchema.validate({ name, website, about, billingEmail });

      const result = await dispatch(
        updateWorkspace({
          workspaceId: currentWorkspace?.id || "",
          data: {
            name,
            website: website || undefined,
            about: about || undefined,
            billingEmail,
          },
        }),
      );

      if ("error" in result) {
        throw new Error(result.payload as string);
      }

      Notifier.success({
        title: "Settings updated",
        message: "Your workspace settings have been updated successfully",
      });
    } catch (error: unknown) {
      Notifier.error({
        title: "Failed to update settings",
        message: (error as Error).message,
      });
    } finally {
      setIsSaving(false);
    }
  };

  const handleLogoClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file && currentWorkspace) {
      setIsUploading(true);
      try {
        const result = await dispatch(
          uploadWorkspaceLogo({
            workspaceId: currentWorkspace.id,
            logo: file,
          }),
        );
        if ("error" in result) {
          throw new Error(result.payload as string);
        }
        Notifier.success({
          title: "Logo Updated",
          message: "Your workspace logo has been updated successfully",
        });
      } catch (error: unknown) {
        Notifier.error({
          title: "Failed to update logo",
          message: (error as Error).message,
        });
      } finally {
        setIsUploading(false);
      }
    }
  };

  const getInitials = (name: string) => {
    return name
      .split(" ")
      .map((n) => n.charAt(0))
      .join("");
  };

  if (!currentWorkspace) return null;

  return (
    <div className="space-y-8">
      <div>
        <form onSubmit={handleWorkspaceUpdate} className="space-y-6">
          <div className="grid columns-3 grid-cols-3">
            <div className="col-span-2">
              <div className="space-y-4 max-w-md mb-4">
                <Input
                  id="workspace-name"
                  label="Workspace Name"
                  value={name}
                  onChange={(v) => setName(v)}
                  required
                  data-testid="workspace-name-input"
                />

                <Input
                  id="billing-email"
                  label="Billing Email"
                  type="email"
                  placeholder="billing@example.com"
                  value={billingEmail}
                  onChange={(v) => setBillingEmail(v)}
                  required
                  data-testid="workspace-billing-email-input"
                />

                <Input
                  id="website"
                  label="Website (optional)"
                  placeholder="https://example.com"
                  value={website}
                  onChange={(v) => setWebsite(v)}
                />

                <div className="space-y-2">
                  <label htmlFor="about" className="text-sm font-medium">
                    About (optional)
                  </label>
                  <Textarea
                    id="about"
                    placeholder="Tell us about your workspace"
                    value={about}
                    onChange={(e) => setAbout(e.target.value)}
                    className="min-h-[100px]"
                  />
                </div>
              </div>
              <Button
                type="submit"
                loading={isSaving}
                data-testid="workspace-save-button"
              >
                Save Changes
              </Button>
            </div>
            <div className="col-span-1 border-l ">
              <div className="flex flex-col items-center gap-4">
                <Avatar
                  className="h-20 w-20 cursor-pointer rounded-lg"
                  onClick={handleLogoClick}
                >
                  <AvatarImage src={currentWorkspace.logo || undefined} />
                  <AvatarFallback className="rounded-lg">
                    {getInitials(currentWorkspace.name)}
                  </AvatarFallback>
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
                  size="sm"
                  variant="outline"
                  onClick={handleLogoClick}
                  loading={isUploading}
                >
                  {isUploading
                    ? "Uploading..."
                    : currentWorkspace.logo
                      ? "Change Logo"
                      : "Upload Logo"}
                </Button>
              </div>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
}
