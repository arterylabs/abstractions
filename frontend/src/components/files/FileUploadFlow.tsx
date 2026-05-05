"use client";

import { useState, useCallback, useEffect } from "react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Notifier } from "@/components/ui/sonner";
import { Progress } from "@/components/ui/progress";
import api from "@/api";
import type { File as UploadedFile } from "@/api/files";
import { getErrorMessage } from "@/lib/errors";

interface FileUploadFlowProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  selectedFile?: File | null;
  onComplete?: () => void;
  maxWidth?: string;
  maxHeight?: string;
}

export default function FileUploadFlow({
  open,
  onOpenChange,
  selectedFile,
  onComplete,
  maxWidth = "640px",
  maxHeight = "420px",
}: FileUploadFlowProps) {
  const [isUploading, setIsUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [uploadedFile, setUploadedFile] = useState<UploadedFile | null>(null);
  const [pendingFile, setPendingFile] = useState<File | null>(
    selectedFile || null,
  );

  useEffect(() => {
    if (!open) {
      setIsUploading(false);
      setUploadProgress(0);
      setUploadedFile(null);
      setPendingFile(selectedFile || null);
    }
  }, [open, selectedFile]);

  const handleFileUpload = useCallback(
    async (file: File) => {
      setIsUploading(true);
      setUploadProgress(0);
      try {
        const formData = new FormData();
        formData.append("file", file);

        const response = await api.files.uploadFile(formData, {
          onUploadProgress: (progressEvent) => {
            const percentCompleted = Math.round(
              (progressEvent.loaded * 100) / (progressEvent.total || file.size),
            );
            setUploadProgress(percentCompleted);
          },
        });

        setUploadedFile(response.data);
        setUploadProgress(100);
        Notifier.success("File uploaded successfully");
        onComplete?.();
      } catch (error: unknown) {
        console.error("Upload failed:", error);
        Notifier.error(getErrorMessage(error));
        setUploadProgress(0);
      } finally {
        setIsUploading(false);
      }
    },
    [onComplete],
  );

  useEffect(() => {
    if (selectedFile && open) {
      setPendingFile(selectedFile);
      handleFileUpload(selectedFile);
    }
  }, [selectedFile, open, handleFileUpload]);

  const handleManualFileSelection = useCallback(
    (event: React.ChangeEvent<HTMLInputElement>) => {
      const file = event.target.files?.[0];
      if (!file) return;
      setPendingFile(file);
    },
    [],
  );

  const handleUploadClick = useCallback(() => {
    if (!pendingFile) {
      Notifier.error("Please select a file first");
      return;
    }
    handleFileUpload(pendingFile);
  }, [pendingFile, handleFileUpload]);

  const handleClose = useCallback(() => {
    if (!isUploading) {
      onOpenChange(false);
    }
  }, [isUploading, onOpenChange]);

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent
        className="w-full flex flex-col"
        style={{ maxWidth, maxHeight, height: "100%" }}
      >
        <DialogHeader>
          <DialogTitle>Upload File</DialogTitle>
        </DialogHeader>

        <div className="space-y-4 py-2">
          {!selectedFile && (
            <input
              type="file"
              onChange={handleManualFileSelection}
              disabled={isUploading}
              className="w-full text-sm"
            />
          )}

          {pendingFile && (
            <p className="text-sm text-muted-foreground">
              Selected: {pendingFile.name}
            </p>
          )}

          {isUploading && (
            <div className="space-y-2">
              <Progress value={uploadProgress} className="w-full" />
              <p className="text-xs text-muted-foreground">
                {uploadProgress}% complete
              </p>
            </div>
          )}

          {uploadedFile && !isUploading && (
            <p className="text-sm text-green-600">
              Uploaded successfully: {uploadedFile.Name}
            </p>
          )}
        </div>

        <DialogFooter>
          <Button
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={isUploading}
          >
            Close
          </Button>
          <Button
            onClick={handleUploadClick}
            disabled={isUploading || !pendingFile}
          >
            {isUploading ? "Uploading..." : "Upload"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
