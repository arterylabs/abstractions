"use client";

import { useRef, useCallback } from "react";
import { Button } from "@/components/ui/button";
import Icon from "@/components/ui/icon";

interface FileUploaderProps {
  onFileSelect: (file: File) => void;
  disabled?: boolean;
  children?: React.ReactNode;
}

export default function FileUploader({
  onFileSelect,
  disabled = false,
  children,
}: FileUploaderProps) {
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFileSelect = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const file = e.target.files?.[0];
      if (!file) return;

      onFileSelect(file);
      // Reset the input so the same file can be selected again if needed
      if (fileInputRef.current) {
        fileInputRef.current.value = "";
      }
    },
    [onFileSelect],
  );

  const handleClick = useCallback(() => {
    fileInputRef.current?.click();
  }, []);

  return (
    <>
      {children ? (
        <div
          onClick={handleClick}
          className={
            disabled ? "pointer-events-none opacity-50" : "cursor-pointer"
          }
        >
          {children}
        </div>
      ) : (
        <Button
          onClick={handleClick}
          disabled={disabled}
          className="flex items-center gap-2"
        >
          <Icon type="upload" className="h-4 w-4" />
          Upload File
        </Button>
      )}
      <input
        type="file"
        ref={fileInputRef}
        className="hidden"
        onChange={handleFileSelect}
      />
    </>
  );
}
