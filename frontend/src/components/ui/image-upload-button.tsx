"use client";

import { useRef } from "react";
import { Avatar, AvatarFallback, AvatarImage } from "./avatar";
import { Button } from "./button";

interface ImageUploadButtonProps {
  /** Current image URL to display in the avatar */
  imageUrl?: string | null;
  /** Fallback content shown when there is no image */
  fallback?: React.ReactNode;
  /** Called with the selected File when the user picks one */
  onFileSelect: (file: File) => void;
  /** Whether an upload is in progress (shows loading state) */
  isUploading?: boolean;
  /** Button label when no image is set */
  uploadLabel?: string;
  /** Button label when an image is already set */
  changeLabel?: string;
  /** Apply rounded-lg to the avatar (useful for logos vs circular avatars) */
  rounded?: boolean;
  /** Accepted MIME types for the file input (default: "image/*") */
  accept?: string;
}

export function ImageUploadButton({
  imageUrl,
  fallback,
  onFileSelect,
  isUploading = false,
  uploadLabel = "Upload Image",
  changeLabel = "Change Image",
  rounded = false,
  accept = "image/*",
}: ImageUploadButtonProps) {
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleClick = () => fileInputRef.current?.click();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) onFileSelect(file);
    // Reset input value so selecting the same file triggers onChange again
    e.target.value = "";
  };

  const avatarClass = `h-20 w-20 cursor-pointer${rounded ? " rounded-lg" : ""}`;
  const fallbackClass = rounded ? "rounded-lg" : "";

  return (
    <div className="flex flex-col items-center gap-4">
      <Avatar className={avatarClass} onClick={handleClick}>
        <AvatarImage src={imageUrl ?? undefined} />
        <AvatarFallback className={fallbackClass}>{fallback}</AvatarFallback>
      </Avatar>
      <input
        type="file"
        ref={fileInputRef}
        className="hidden"
        accept={accept}
        onChange={handleChange}
      />
      <Button
        type="button"
        variant="outline"
        size="sm"
        onClick={handleClick}
        loading={isUploading}
      >
        {imageUrl ? changeLabel : uploadLabel}
      </Button>
    </div>
  );
}
