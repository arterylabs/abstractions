"use client";

import { FormEvent, useEffect, useState } from "react";
import * as yup from "yup";
import { Button } from "@/components/ui/button";
import { ImageUploadButton } from "@/components/ui/image-upload-button";
import { Input } from "@/components/ui/input";
import { Notifier } from "@/components/ui/sonner";
import { Textarea } from "@/components/ui/textarea";
import { SectionMark } from "@/components/ui/section-mark";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import {
  uploadWorkspaceLogo,
  updateWorkspace,
} from "@/store/modules/workspace";
import { getInitials } from "@/lib/user";

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
        title: "Settings saved",
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

  const handleLogoFileSelect = async (file: File) => {
    if (!currentWorkspace) return;
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
        title: "Logo updated",
      });
    } catch (error: unknown) {
      Notifier.error({
        title: "Failed to update logo",
        message: (error as Error).message,
      });
    } finally {
      setIsUploading(false);
    }
  };

  if (!currentWorkspace) return null;

  return (
    <section className="flex flex-col gap-5">
      <SectionMark label="General" />
      <div>
        <form onSubmit={handleWorkspaceUpdate} className="space-y-6">
          <div className="grid grid-cols-1 border border-border bg-card p-6 md:grid-cols-3">
            <div className="col-span-1 md:col-span-2">
              <div className="mb-5 flex max-w-md flex-col gap-5">
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
                  data-testid="workspace-website-input"
                />

                <div className="space-y-2">
                  <label
                    htmlFor="about"
                    className="press-label text-muted-foreground"
                  >
                    About (optional)
                  </label>
                  <Textarea
                    id="about"
                    placeholder="Tell us about your workspace"
                    value={about}
                    onChange={(e) => setAbout(e.target.value)}
                    className="min-h-[100px]"
                    data-testid="workspace-about-input"
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
            <div className="col-span-1 border-t md:border-t-0 md:border-l pt-6 md:pt-0">
              <ImageUploadButton
                imageUrl={currentWorkspace.logo}
                fallback={getInitials(currentWorkspace.name)}
                onFileSelect={handleLogoFileSelect}
                isUploading={isUploading}
                uploadLabel="Upload Logo"
                changeLabel="Change Logo"
              />
            </div>
          </div>
        </form>
      </div>
    </section>
  );
}
