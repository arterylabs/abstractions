"use client";

import { useState, useCallback, useEffect } from "react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { SubscriptionPlanId, SubscriptionStatus } from "@/types/user";
import api from "@/api";
import useAPI from "@/api/hooks";
import Spinner from "@/components/ui/spinner";
import { Notifier } from "@/components/ui/sonner";
import { askConfirmation } from "@/components/ui/confirmation";
import { format, formatISO, parseISO } from "date-fns";

interface WorkspaceSubscriptionTabProps {
  workspaceId: string | null;
}

export default function WorkspaceSubscriptionTab({
  workspaceId,
}: WorkspaceSubscriptionTabProps) {
  const [subscriptionResponse, getSubscription] = useAPI(
    api.admin.getWorkspaceSubscription,
  );

  const [showSubscriptionForm, setShowSubscriptionForm] = useState(false);
  const [subscriptionFormData, setSubscriptionFormData] = useState({
    planId: "BASIC" as SubscriptionPlanId,
    expiresAt: "",
    status: "active" as SubscriptionStatus,
  });

  const loadSubscription = useCallback(async () => {
    if (workspaceId) {
      await getSubscription(workspaceId);
    }
  }, [getSubscription, workspaceId]);

  useEffect(() => {
    loadSubscription();
  }, [loadSubscription]);

  const handleCreateSubscription = async () => {
    if (!workspaceId) return;

    try {
      // Convert datetime-local to ISO string using date-fns
      console.log("Original datetime value:", subscriptionFormData.expiresAt);
      const expiresAtISO = formatISO(new Date(subscriptionFormData.expiresAt));
      console.log("Converted to ISO with date-fns:", expiresAtISO);

      await api.admin.createWorkspaceSubscription(workspaceId, {
        planId: subscriptionFormData.planId,
        expiresAt: expiresAtISO,
        workspaceId: "",
        status: subscriptionFormData.status,
        limit: 0,
      });
      Notifier.success("Subscription created successfully");
      setShowSubscriptionForm(false);
      await getSubscription(workspaceId); // Reload subscription data
    } catch (error) {
      console.error("Failed to create subscription:", error);
      Notifier.error("Failed to create subscription");
    }
  };

  const handleDisableSubscription = async () => {
    if (!workspaceId || !subscription) return;

    askConfirmation({
      title: "Disable Subscription",
      content:
        "Are you sure you want to disable this subscription? This action will affect the workspace's access to premium features.",
      variant: "destructive",
      okText: "Disable",
      cancelText: "Cancel",
      okCallback: async () => {
        try {
          await api.admin.disableWorkspaceSubscription(workspaceId, {
            disabledReason: "Disabled by admin",
          });
          Notifier.success("Subscription disabled successfully");
          await getSubscription(workspaceId); // Reload subscription data
        } catch {
          Notifier.error("Failed to disable subscription");
        }
      },
    });
  };

  const formatDate = (dateString: string | undefined) => {
    if (!dateString) return "Invalid date";
    try {
      return format(new Date(dateString), "PPpp"); // e.g., "Apr 29, 2023 at 3:45 PM"
    } catch {
      return "Invalid date";
    }
  };

  // Helper function to convert ISO datetime to datetime-local format
  const formatForDatetimeLocal = (isoString: string | undefined): string => {
    if (!isoString) return "";
    try {
      const date = parseISO(isoString);
      // Format as YYYY-MM-DDTHH:mm for datetime-local input
      return format(date, "yyyy-MM-dd'T'HH:mm");
    } catch {
      return "";
    }
  };

  const subscription = subscriptionResponse.result;

  return (
    <div className="space-y-6">
      <h3 className="text-lg font-medium">Subscription Details</h3>

      {subscriptionResponse.isFetching ? (
        <div className="flex items-center justify-center py-12">
          <Spinner />
        </div>
      ) : subscription ? (
        <div className="space-y-6">
          <div className="bg-accent/50 rounded-lg p-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-2">
                <Label className="text-sm font-medium text-foreground">
                  Plan
                </Label>
                <p className="text-base font-medium">{subscription.planId}</p>
              </div>
              <div className="space-y-2">
                <Label className="text-sm font-medium text-foreground">
                  Status
                </Label>
                <Badge
                  variant={
                    subscription.status === "active" ? "default" : "destructive"
                  }
                  className="text-sm"
                >
                  {subscription.status}
                </Badge>
              </div>
              <div className="space-y-2">
                <Label className="text-sm font-medium text-foreground">
                  Expires At
                </Label>
                <p className="text-sm text-muted-foreground">
                  {subscription.expiresAt
                    ? formatDate(subscription.expiresAt)
                    : "Never"}
                </p>
              </div>
              <div className="space-y-2">
                <Label className="text-sm font-medium text-foreground">
                  Limit
                </Label>
                <p className="text-sm text-muted-foreground">
                  {subscription.limit}
                </p>
              </div>
            </div>
          </div>

          <div className="flex space-x-3">
            <Button
              variant="outline"
              size="sm"
              onClick={() => {
                // Populate form with existing subscription data
                setSubscriptionFormData({
                  planId: subscription.planId,
                  expiresAt: formatForDatetimeLocal(subscription.expiresAt),
                  status: subscription.status,
                });
                setShowSubscriptionForm(true);
              }}
              className="h-9"
            >
              Edit Subscription
            </Button>
            <Button
              variant="destructive"
              size="sm"
              onClick={handleDisableSubscription}
              className="h-9"
            >
              Disable
            </Button>
          </div>
        </div>
      ) : (
        <div className="text-center py-12">
          <p className="text-muted-foreground mb-4">No active subscription</p>
          <Button
            variant="outline"
            onClick={() => setShowSubscriptionForm(true)}
          >
            Create Subscription
          </Button>
        </div>
      )}

      {/* Subscription Form */}
      {showSubscriptionForm && (
        <div className="mt-6 p-6 border rounded-lg bg-accent/30">
          <h4 className="text-lg font-medium mb-4">
            {subscription ? "Edit Subscription" : "Create Subscription"}
          </h4>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="planId">Plan</Label>
              <Select
                value={subscriptionFormData.planId}
                onValueChange={(value) =>
                  setSubscriptionFormData((prev) => ({
                    ...prev,
                    planId: value as SubscriptionPlanId,
                  }))
                }
              >
                <SelectTrigger className="w-full h-10">
                  <SelectValue placeholder="Select plan" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="BASIC">Basic</SelectItem>
                  <SelectItem value="PRO">Pro</SelectItem>
                  <SelectItem value="ENTERPRISE">Enterprise</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="status">Status</Label>
              <Select
                value={subscriptionFormData.status}
                onValueChange={(value) =>
                  setSubscriptionFormData((prev) => ({
                    ...prev,
                    status: value as SubscriptionStatus,
                  }))
                }
              >
                <SelectTrigger className="w-full h-10">
                  <SelectValue placeholder="Select status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="active">Active</SelectItem>
                  <SelectItem value="inactive">Inactive</SelectItem>
                  <SelectItem value="trialing">Trialing</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="expiresAt">Expires At *</Label>
              <Input
                id="expiresAt"
                type="datetime-local"
                value={subscriptionFormData.expiresAt}
                onChange={(value) =>
                  setSubscriptionFormData((prev) => ({
                    ...prev,
                    expiresAt: value,
                  }))
                }
                required
                className="h-10"
              />
            </div>
            <div className="flex space-x-3 pt-2">
              <Button
                size="sm"
                onClick={handleCreateSubscription}
                disabled={!subscriptionFormData.expiresAt}
                className="h-9"
              >
                {subscription ? "Update" : "Create"}
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setShowSubscriptionForm(false)}
                className="h-9"
              >
                Cancel
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
