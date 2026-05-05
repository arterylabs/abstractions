"use client";

import { useEffect, useCallback, useMemo } from "react";
import api from "@/api";
import useAPI from "@/api/hooks";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import SlackIntegrationCard from "@/components/dashboard/SlackIntegrationCard";
import { IntegrationType } from "@/types";
import { fetchIntegrationTypes } from "@/store/modules/integrations";
import { PageHeader } from "@/components/ui/page-header";

export default function DashboardHome() {
  const dispatch = useAppDispatch();
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );
  const { types: integrationTypes, typesStatus } = useAppSelector(
    (state) => state.integrations,
  );

  const [integrationResponse, getSlackIntegration] = useAPI(
    api.integrations.getSlackIntegration,
  );

  // Load Slack integration instance
  const loadIntegration = useCallback(async () => {
    if (currentWorkspace?.id) {
      await getSlackIntegration(currentWorkspace.id);
    }
  }, [getSlackIntegration, currentWorkspace?.id]);

  // Fetch integration types metadata (once)
  useEffect(() => {
    if (typesStatus === "idle") {
      dispatch(fetchIntegrationTypes());
    }
  }, [dispatch, typesStatus]);

  useEffect(() => {
    loadIntegration();
  }, [loadIntegration]);

  const handleConnectSlack = async (redirectTo?: string) => {
    if (currentWorkspace?.id) {
      try {
        const response = await api.integrations.setupIntegration(
          currentWorkspace.id,
          IntegrationType.SLACK,
          redirectTo,
        );
        window.location.href = response.data.authorize_url;
      } catch (error) {
        console.error("Failed to setup integration:", error);
      }
    }
  };

  // Derive Slack integration type metadata from store
  const slackTypeMeta = useMemo(
    () => integrationTypes?.find((t) => t.type === IntegrationType.SLACK),
    [integrationTypes],
  );

  return (
    <>
      <PageHeader
        title="👋 Welcome to Abstractions."
        description="You're all set up and ready to go! Here's a summary of what's happening in your Slack Connect channels"
      ></PageHeader>
      <div className="mt-12">
        <SlackIntegrationCard
          integration={integrationResponse.result}
          // Combine loading states: integration fetch + types fetch
          isLoading={
            integrationResponse.isFetching || typesStatus === "loading"
          }
          onConnect={handleConnectSlack}
          onIntegrationUpdated={loadIntegration}
          meta={slackTypeMeta}
        />
      </div>
    </>
  );
}
