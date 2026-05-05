"use client";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { syncCurrentUser } from "@/store/modules/auth";
import { useRouter, usePathname } from "next/navigation";
import { useEffect } from "react";
import Spinner from "./ui/spinner";

export default function AuthWrapper({
  children,
  className = "",
}: {
  children: React.ReactNode;
  className?: string;
}) {
  const dispatch = useAppDispatch();
  const router = useRouter();
  const pathname = usePathname();

  useEffect(() => {
    dispatch(syncCurrentUser());
  }, []);

  const currentUserStatus = useAppSelector(
    (state) => state.auth.currentUserStatus,
  );

  useEffect(() => {
    if (currentUserStatus === "failed") {
      const loginUrl = `/login?ref=${encodeURIComponent(pathname)}`;
      router.push(loginUrl);
    }
  }, [currentUserStatus, pathname, router]);

  if (currentUserStatus === "success") {
    return children;
  }

  return (
    <div className={`flex items-center justify-center ${className}`}>
      <Spinner />
    </div>
  );
}
