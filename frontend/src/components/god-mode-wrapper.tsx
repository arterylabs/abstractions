"use client";
import { useAppSelector } from "@/store/hooks";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import AuthWrapper from "./authenticated";

export default function GodModeWrapper({
  children,
  className = "",
}: {
  children: React.ReactNode;
  className?: string;
}) {
  const router = useRouter();
  const currentUser = useAppSelector((state) => state.auth.currentUser);
  const currentUserStatus = useAppSelector(
    (state) => state.auth.currentUserStatus,
  );

  useEffect(() => {
    // Only check god mode after user authentication is successful
    if (
      currentUserStatus === "success" &&
      currentUser &&
      !currentUser.godMode
    ) {
      // Redirect to 404 to hide the existence of god mode pages
      router.push("/404");
    }
  }, [currentUserStatus, currentUser, router]);

  // Use AuthWrapper to handle authentication first
  return (
    <AuthWrapper className={className}>
      {/* Only render children if user is authenticated and has god mode */}
      {currentUserStatus === "success" && currentUser?.godMode
        ? children
        : null}
    </AuthWrapper>
  );
}
