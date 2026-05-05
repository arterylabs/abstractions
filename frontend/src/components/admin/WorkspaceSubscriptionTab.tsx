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
import { useAPI } from "@/hooks/use-api";
import Spinner from "@/components/ui/spinner";
import { Notifier } from "@/components/ui/sonner";
import { askConfirmation } from "@/components/ui/confirmation";
import { format, formatISO, parseISO } from "date-fns";
import { formatDateTime } from "@/lib/format";
import { getErrorMessage } from "@/lib/errors";

interface WorkspaceSubscriptionTabProps {
  userId: string | null;
}

export default function WorkspaceSubscriptionTab({
  userId,
}: WorkspaceSubscriptionTabProps) {
  const [subscriptionResponse, getSubscription] = useAPI(
    api.admin.getUserSubscription,
  );

  const [showSubscriptionForm, setShowSubscriptionForm] = useState(false);
  const [subscriptionFormData, setSubscriptionFormData] = useState({
    planId: "BASIC" as SubscriptionPlanId,
    expiresAt: "",
    status: "active" as SubscriptionStatus,
  });

  const loadSubscription = useCallback(async () => {
    if (userId) {
      await getSubscription(userId);
    }
  }, [getSubscription, userId]);

  useEffect(() => {
    loadSubscription();
  }, [loadSubscription]);

  const handleCreateSubscription = async () => {
    if (!userId) return;

    try {
      const expiresAtISO = formatISO(new Date(subscriptionFormData.expiresAt));
      await api.admin.createUserSubscription(userId, {
        planId: subscriptionFormData.planId,
        expiresAt: expiresAtISO,
        userId: userId,
        status: subscriptionFormData.status,
        limit: 0,
      });
      Notifier.success("Subscription created successfully");
      setShowSubscriptionForm(false);
      await getSubscription(userId);
    } catch (error) {
      console.error("Failed to create subscription:", error);
      Notifier.error(getErrorMessage(error));
    }
  };

  const handleDisableSubscription = async () => {
    if (!userId || !subscription) return;

    askConfirmation({
      title: "Disable Subscription",
      content: "Are you sure you want to disable this subscription?",
      variant: "destructive",
      okText: "Disable",
      cancelText: "Cancel",
      okCallback: async () => {
        try {
          await api.admin.disableUserSubscription(userId, {
            disabledReason: "Disabled by admin",
          });
          Notifier.success("Subscription disabled successfully");
          await getSubscription(userId);
        } catch (error) {
          Notifier.error(getErrorMessage(error));
        }
      },
    });
  };

  const formatForDatetimeLocal = (isoString: string | undefined): string => {
    if (!isoString) return "";
    try {
      const date = parseISO(isoString);
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
          <div className="bg-gray-50/50 rounded-lg p-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-2">
                <Label className="text-sm font-medium text-gray-700">
                  Plan
                </Label>
                <p className="text-base font-medium">{subscription.planId}</p>
              </div>
              <div className="space-y-2">
                <Label className="text-sm font-medium text-gray-700">
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
                <Label className="text-sm font-medium text-gray-700">
                  Expires At
                </Label>
                <p className="text-sm text-gray-600">
                  {subscription.expiresAt
                    ? (formatDateTime(subscription.expiresAt, "PPpp") ??
                      "Invalid date")
                    : "Never"}
                </p>
              </div>
              <div className="space-y-2">
                <Label className="text-sm font-medium text-gray-700">
                  Limit
                </Label>
                <p className="text-sm text-gray-600">{subscription.limit}</p>
              </div>
            </div>
          </div>

          <div className="flex space-x-3">
            <Button
              variant="outline"
              size="sm"
              onClick={() => {
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
          <p className="text-gray-500 mb-4">No active subscription</p>
          <Button
            variant="outline"
            onClick={() => setShowSubscriptionForm(true)}
          >
            Create Subscription
          </Button>
        </div>
      )}

      {showSubscriptionForm && (
        <div className="mt-6 p-6 border rounded-lg bg-gray-50/30">
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
