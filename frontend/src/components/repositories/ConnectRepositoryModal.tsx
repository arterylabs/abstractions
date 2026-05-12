"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import api from "@/api";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Notifier } from "@/components/ui/sonner";
import { useAppSelector } from "@/store/hooks";
import { getErrorMessage } from "@/lib/errors";

interface ConnectRepositoryModalProps {
  open: boolean;
  onClose: () => void;
}

export default function ConnectRepositoryModal({
  open,
  onClose,
}: ConnectRepositoryModalProps) {
  const router = useRouter();
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );
  const [isConnecting, setIsConnecting] = useState(false);

  const handleConnectWithGitHub = async () => {
    if (!currentWorkspace?.id) return;
    setIsConnecting(true);
    try {
      const response = await api.github.initiateOAuth(currentWorkspace.id, {
        createRepository: true,
      });
      // Redirect to GitHub OAuth — the callback will return to /repositories/connect
      window.location.href = response.data.url;
    } catch (error: unknown) {
      Notifier.error({
        title: "Unable to initiate GitHub connection",
        message: getErrorMessage(error),
      });
      setIsConnecting(false);
    }
  };

  const handleSelectFromConnected = () => {
    onClose();
    router.push("/dashboard/repositories/connect");
  };

  return (
    <Dialog open={open} onOpenChange={(isOpen) => !isOpen && onClose()}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Connect a Repository</DialogTitle>
          <DialogDescription>
            Connect a GitHub repository to start indexing code and generating
            newsletters.
          </DialogDescription>
        </DialogHeader>
        <div className="flex flex-col gap-3 pt-2">
          <Button
            onClick={() => void handleConnectWithGitHub()}
            loading={isConnecting}
            className="w-full"
          >
            Connect with GitHub
          </Button>
          <Button
            variant="outline"
            onClick={handleSelectFromConnected}
            className="w-full"
          >
            Select from connected repositories
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
