"use client";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { getInitials } from "@/lib/user";

interface UserAvatarProps {
  name: string;
  imageUrl?: string | null;
  alt?: string;
  className?: string;
  fallbackClassName?: string;
}

export function UserAvatar({
  name,
  imageUrl,
  alt,
  className,
  fallbackClassName,
}: UserAvatarProps) {
  return (
    <Avatar className={className}>
      {imageUrl && <AvatarImage src={imageUrl} alt={alt ?? name} />}
      <AvatarFallback name={name} className={fallbackClassName}>
        {getInitials(name)}
      </AvatarFallback>
    </Avatar>
  );
}
