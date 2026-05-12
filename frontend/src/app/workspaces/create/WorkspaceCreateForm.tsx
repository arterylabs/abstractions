"use client";

import { useRouter } from "next/navigation";
import { FormEvent, useState, useEffect, useRef } from "react";
import * as yup from "yup";
import api from "@/api";
import useAPI from "@/api/hooks";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Notifier } from "@/components/ui/sonner";
import { Textarea } from "@/components/ui/textarea";
import Image from "next/image";
import { uploadWorkspaceLogo } from "@/store/modules/workspace";
import { useAppDispatch } from "@/store/hooks";
import { setCurrentWorkspace, switchWorkspace } from "@/store/modules/auth";

const workspaceSchema = yup.object().shape({
  name: yup.string().required("Workspace name is required").max(64),
  website: yup.string().url("Must be a valid URL").nullable(),
  about: yup.string().max(500, "About text is too long"),
});

export default function WorkspaceCreateForm() {
  const [name, setName] = useState("");
  const [website, setWebsite] = useState("");
  const [about, setAbout] = useState("");
  const [isUploading, setIsUploading] = useState(false);
  const [logo, setLogo] = useState<File | null>(null);
  const [logoPreviewUrl, setLogoPreviewUrl] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const router = useRouter();
  const dispatch = useAppDispatch();

  // Destructure the response object from useAPI hook
  const [{ result: createWorkspaceResponse, isFetching }, sendCreateWorkspace] =
    useAPI(() =>
      api.workspaces.createWorkspace({
        name,
        website: website || undefined,
        about: about || undefined,
      }),
    );

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    e.stopPropagation();

    try {
      await workspaceSchema.validate({ name, website, about });
    } catch (error: unknown) {
      Notifier.error({
        title: "Error",
        message: (error as Error).message,
      });
      return;
    }

    sendCreateWorkspace();
  };

  // Handle workspace creation response
  useEffect(() => {
    const handleWorkspaceCreated = async () => {
      console.log("Workspace creation response:", createWorkspaceResponse);
      if (createWorkspaceResponse?.id) {
        // If logo is selected, upload it after workspace is created
        console.log("Uploading logo:", logo);
        if (logo) {
          try {
            const result = await dispatch(
              uploadWorkspaceLogo({
                workspaceId: createWorkspaceResponse.id,
                logo,
              }),
            );
            if ("error" in result) {
              throw new Error(result.payload as string);
            }
          } catch (error: unknown) {
            Notifier.error({
              title: "Error uploading logo",
              message: (error as Error).message,
            });
          }
        }

        // Set this as the current workspace
        dispatch(setCurrentWorkspace(createWorkspaceResponse));

        // Update it as the default workspace
        try {
          await dispatch(switchWorkspace(createWorkspaceResponse.id));
        } catch {
          Notifier.error({
            title: "Error setting default workspace",
            message:
              "Your workspace was created but couldn't be set as default.",
          });
        }

        Notifier.success({
          title: "Workspace created",
          message: "Redirecting to dashboard...",
        });
        router.push("/dashboard");
      }
    };

    handleWorkspaceCreated();
  }, [createWorkspaceResponse, dispatch, logo, router]);

  const handleLogoClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setIsUploading(true);
      try {
        setLogo(file);
        setLogoPreviewUrl(URL.createObjectURL(file));
      } catch {
        Notifier.error({
          title: "Failed to update logo",
          message: "An error occurred while updating the logo.",
        });
      } finally {
        setIsUploading(false);
      }
    }
  };

  // Cleanup object URL on unmount
  useEffect(() => {
    return () => {
      if (logoPreviewUrl) {
        URL.revokeObjectURL(logoPreviewUrl);
      }
    };
  }, [logoPreviewUrl]);

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div>
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
          onClick={handleLogoClick}
          loading={isUploading}
          className="w-full h-32 flex flex-col items-center justify-center gap-2"
        >
          {logoPreviewUrl ? (
            <>
              <Image
                src={logoPreviewUrl}
                alt="Selected logo"
                width={96}
                height={96}
                className="max-h-24 w-auto object-contain"
              />
              <span className="text-sm text-muted-foreground">Change Logo</span>
            </>
          ) : (
            <>
              <span className="text-2xl">+</span>
              <span className="text-sm text-muted-foreground">Upload Logo</span>
            </>
          )}
        </Button>
      </div>

      <Input
        id="workspace-name"
        label="Workspace Name"
        placeholder="Enter workspace name"
        value={name}
        onChange={(v) => setName(v)}
        required
        data-testid="create-workspace-name-input"
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

      <Button
        type="submit"
        className="w-full"
        loading={isFetching}
        data-testid="create-workspace-submit"
      >
        Create Workspace
      </Button>
    </form>
  );
}
