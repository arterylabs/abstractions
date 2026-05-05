"use client";
import { Button } from "@/components/ui/button";
import { useRouter } from "next/navigation";

export default function AbstractionsLanding() {
  const router = useRouter();

  const navigateTo = (url: string) => {
    router.push(url);
  };

  return (
    <div className="flex flex-col items-center justify-center h-screen gap-6">
      Abstractions landing page
      <Button
        onClick={() => {
          navigateTo("/login");
        }}
      >
        Login
      </Button>
      <Button
        onClick={() => {
          navigateTo("/signup");
        }}
      >
        Sign Up
      </Button>
      <Button
        onClick={() => {
          navigateTo("/dashboard");
        }}
      >
        Dashboard
      </Button>
    </div>
  );
}
