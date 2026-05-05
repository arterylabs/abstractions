"use client";

import { useEffect, useMemo } from "react";
import { useSearchParams } from "next/navigation";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { Notifier } from "@/components/ui/sonner";
import IntegrationCard from "@/components/dashboard/IntegrationCard";
import api from "@/api";
import useAPI from "@/api/hooks";
import { IntegrationType, IntegrationTypeMetadata } from "@/types";
import { fetchIntegrationTypes } from "@/store/modules/integrations";
import { PageHeader } from "@/components/ui/page-header";

export default function IntegrationsPage() {
  const searchParams = useSearchParams();
  const dispatch = useAppDispatch();
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );

  const [integrationResponse, getIntegrations] = useAPI(
    api.integrations.getWorkspaceIntegrations,
  );

  const { types: integrationTypes, typesStatus } = useAppSelector(
    (state) => state.integrations,
  );

  // Fetch integration types once (if not already loaded)
  useEffect(() => {
    if (typesStatus === "idle") {
      dispatch(fetchIntegrationTypes());
    }
  }, [dispatch, typesStatus]);

  useEffect(() => {
    if (currentWorkspace?.id) {
      getIntegrations(currentWorkspace.id);
    }
  }, [getIntegrations, currentWorkspace?.id]);

  useEffect(() => {
    const success = searchParams.get("success");
    const error = searchParams.get("error");
    const type = searchParams.get("type");

    if (success === "true" && type) {
      Notifier.success({
        title: "Integration Connected",
        message: `Successfully connected ${
          type.charAt(0).toUpperCase() + type.slice(1)
        } integration!`,
      });

      if (currentWorkspace?.id) {
        getIntegrations(currentWorkspace.id);
      }
    } else if (error) {
      const errorMessages: { [key: string]: string } = {
        invalid_type: "Invalid integration type",
        missing_parameters: "Missing required parameters",
        invalid_state: "Invalid OAuth state",
        unsupported_type: "Unsupported integration type",
        token_exchange_failed: "Failed to exchange OAuth token",
        encryption_failed: "Failed to encrypt integration credentials",
        user_not_found: "User authentication failed",
        no_workspace: "No workspace available",
        creation_failed: "Failed to create integration",
      };

      Notifier.error({
        title: "Integration Failed",
        message: errorMessages[error] || "An unknown error occurred",
      });
    }
  }, [searchParams, currentWorkspace?.id, getIntegrations]);

  const slackTypeMeta: IntegrationTypeMetadata | undefined = useMemo(
    () =>
      integrationTypes?.find(
        (t: IntegrationTypeMetadata) => t.type === IntegrationType.SLACK,
      ),
    [integrationTypes],
  );

  const slackIntegration = useMemo(
    () =>
      integrationResponse.result?.results.find(
        (i) => i.type === IntegrationType.SLACK,
      ),
    [integrationResponse.result?.results],
  );

  const handleIntegrationUpdated = () => {
    if (currentWorkspace?.id) {
      getIntegrations(currentWorkspace.id);
    }
  };

  return (
    <>
      <PageHeader
        title="Integrations"
        description="Connect your workspace to external services and tools."
        className="space-y-8"
      ></PageHeader>

      <div className="grid gap-6 mt-4">
        <IntegrationCard
          integrationType={IntegrationType.SLACK}
          integration={slackIntegration}
          isLoading={
            integrationResponse.isFetching || typesStatus === "loading"
          }
          onIntegrationUpdated={handleIntegrationUpdated}
          meta={slackTypeMeta}
        />
      </div>
    </>
  );
}
