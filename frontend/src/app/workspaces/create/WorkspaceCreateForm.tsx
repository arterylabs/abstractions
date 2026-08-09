"use client";

import { useRouter } from "next/navigation";
import { FormEvent, useState, useEffect } from "react";
import * as yup from "yup";
import api from "@/api";
import useAPI from "@/api/hooks";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Notifier } from "@/components/ui/sonner";
import { Textarea } from "@/components/ui/textarea";
import { uploadWorkspaceLogo } from "@/store/modules/workspace";
import Icon from "@/components/ui/icon";
import { useRef } from "react";
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
  const [logo, setLogo] = useState<File | null>(null);
  const [logoPreviewUrl, setLogoPreviewUrl] = useState<string | null>(null);
  const logoInputRef = useRef<HTMLInputElement>(null);
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
      if (createWorkspaceResponse?.id) {
        // If logo is selected, upload it after workspace is created
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

  const handleLogoFileSelect = (file: File) => {
    setLogo(file);
    setLogoPreviewUrl(URL.createObjectURL(file));
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
    <form onSubmit={handleSubmit}>
      <div className="grid grid-cols-3 gap-4">
        <div className="col-span-2 space-y-4 max-w-md">
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
            data-testid="create-workspace-website-input"
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
              data-testid="create-workspace-about-input"
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
        </div>

        <div className="col-span-1 border-l pl-4 flex flex-col items-center gap-3">
          <input
            type="file"
            ref={logoInputRef}
            className="hidden"
            accept="image/*"
            onChange={(e) => {
              const file = e.target.files?.[0];
              if (file) handleLogoFileSelect(file);
              e.target.value = "";
            }}
          />
          {logoPreviewUrl ? (
            <>
              <div
                className="h-24 w-24 cursor-pointer overflow-hidden rounded-sm border border-border"
                onClick={() => logoInputRef.current?.click()}
              >
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img
                  src={logoPreviewUrl}
                  alt="Workspace logo preview"
                  className="w-full h-full object-cover"
                />
              </div>
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => logoInputRef.current?.click()}
              >
                Change Logo
              </Button>
            </>
          ) : (
            <div
              className="flex h-24 w-24 cursor-pointer items-center justify-center rounded-sm border border-dashed border-border transition-colors hover:bg-accent"
              onClick={() => logoInputRef.current?.click()}
            >
              <Icon type="add" className="text-2xl text-muted-foreground" />
            </div>
          )}
        </div>
      </div>
    </form>
  );
}
