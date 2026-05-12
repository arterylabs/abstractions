"use client";

import React from "react";
import { useAppSelector } from "@/store/hooks";
import { Button } from "@/components/ui/button";
import { useAPI } from "@/hooks/use-api";
import api from "@/api";
import Icon from "@/components/ui/icon";
import { Notifier } from "@/components/ui/sonner";

export default function ImpersonationBanner() {
  const currentUser = useAppSelector((state) => state.auth.currentUser);
  const [exitImpersonationResponse, exitImpersonation] = useAPI(
    api.admin.exitImpersonation,
  );

  // Handle exit impersonation response
  React.useEffect(() => {
    if (exitImpersonationResponse.result) {
      // Get the stored impersonator token (admin's token) from localStorage
      const impersonatorToken = localStorage.getItem("impersonatorToken");

      if (impersonatorToken) {
        // Restore the original god mode user's token
        api.utils.setToken(impersonatorToken, "");

        // Remove the stored impersonator token
        localStorage.removeItem("impersonatorToken");
      }

      Notifier.success("Exited impersonation successfully");

      // Redirect to god mode users page
      window.location.href = "/god-mode/users";
    }

    if (exitImpersonationResponse.error) {
      Notifier.error("Failed to exit impersonation");
    }
  }, [exitImpersonationResponse.result, exitImpersonationResponse.error]);

  // Add top margin to body and adjust sidebar position when impersonation is active
  React.useEffect(() => {
    const isImpersonated = currentUser?.impersonation?.isImpersonated;

    if (isImpersonated) {
      // Add margin to body for regular content
      document.body.style.marginTop = "30px";

      // Adjust fixed positioned elements (sidebar container)
      const sidebarContainer = document.querySelector("#app-sidebar");
      if (sidebarContainer && sidebarContainer instanceof HTMLElement) {
        sidebarContainer.style.top = "30px";
        sidebarContainer.style.height = "calc(100vh - 30px)";
      }
    } else {
      // Reset body margin
      document.body.style.marginTop = "";

      // Reset sidebar position
      const sidebarContainer = document.querySelector("#app-sidebar");
      if (sidebarContainer && sidebarContainer instanceof HTMLElement) {
        sidebarContainer.style.top = "";
        sidebarContainer.style.height = "";
      }
    }

    // Cleanup function to remove styles when component unmounts
    return () => {
      document.body.style.marginTop = "";
      const sidebarContainer = document.querySelector("#app-sidebar");
      if (sidebarContainer && sidebarContainer instanceof HTMLElement) {
        sidebarContainer.style.top = "";
        sidebarContainer.style.height = "";
      }
    };
  }, [currentUser?.impersonation?.isImpersonated]);

  // Don't show banner if not impersonated
  if (!currentUser?.impersonation?.isImpersonated) {
    return null;
  }

  const handleExitImpersonation = async () => {
    try {
      // Get the stored impersonator token (admin's token) from localStorage
      const impersonatorToken = localStorage.getItem("impersonatorToken");

      if (!impersonatorToken) {
        Notifier.error(
          "Cannot exit impersonation: Original admin token not found",
        );
        return;
      }

      // Store current impersonated token
      const currentImpersonatedToken = api.utils.getToken();

      // Switch to the admin token BEFORE calling the API
      api.utils.setToken(impersonatorToken, "");

      try {
        // Now call the useAPI hook's method - it will use the admin token
        await exitImpersonation();
      } catch (error) {
        // If the call fails, restore the impersonated token
        if (currentImpersonatedToken) {
          api.utils.setToken(currentImpersonatedToken, "");
        }
        throw error;
      }
    } catch (error) {
      console.error("Failed to exit impersonation:", error);
      Notifier.error("Failed to exit impersonation");
    }
  };

  const impersonatorName =
    currentUser.impersonation.impersonatorUser?.name || "Unknown Admin";

  return (
    <div className="fixed top-0 left-0 right-0 z-50 border-b border-destructive-dark bg-destructive text-destructive-foreground shadow-sm h-[30px]">
      <div className="flex items-center justify-between px-4 h-full">
        <div className="flex items-center gap-2 text-sm font-medium">
          <Icon type="warning" className="h-3 w-3 flex-shrink-0" />
          <span className="hidden sm:inline">Impersonation Active:</span>
          <span className="font-semibold">{currentUser.name}</span>
          <span className="hidden md:inline">by {impersonatorName}</span>
        </div>
        <Button
          onClick={handleExitImpersonation}
          disabled={exitImpersonationResponse.isFetching}
          variant="outline"
          size="sm"
          className="h-6 px-2 py-0 text-xs "
        >
          <Icon type="logout" className="h-3 w-3 mr-1" />
          {exitImpersonationResponse.isFetching
            ? "Exiting..."
            : "Exit Impersonation"}
        </Button>
      </div>
    </div>
  );
}
