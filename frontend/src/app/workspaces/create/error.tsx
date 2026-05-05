"use client";

import { useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <div className="flex items-center justify-center min-h-screen bg-slate-50">
      <Card className="w-full max-w-md p-6">
        <div className="space-y-6 text-center">
          <h2 className="text-2xl font-semibold">Something went wrong!</h2>
          <p className="text-gray-600">
            {error.message || "Failed to create workspace"}
          </p>
          <div className="flex justify-center gap-4">
            <Button variant="outline" onClick={reset}>
              Try again
            </Button>
            <Button onClick={() => (window.location.href = "/dashboard")}>
              Go to Dashboard
            </Button>
          </div>
        </div>
      </Card>
    </div>
  );
}
