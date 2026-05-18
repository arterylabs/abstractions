"use client";

import { useEffect, useRef, useState } from "react";
import { useSearchParams } from "next/navigation";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { Notifier } from "@/components/ui/sonner";
import api from "@/api";
import useAPI from "@/api/hooks";
import { fetchIntegrationTypes } from "@/store/modules/integrations";
import { PageHeader } from "@/components/ui/page-header";
import { Integration, IntegrationType } from "@/types/integration";
import { IntegrationCard } from "@/components/integrations/IntegrationCard";

export default function IntegrationsPage() {
  const searchParams = useSearchParams();
  const dispatch = useAppDispatch();
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );
  const { types: integrationTypes, typesStatus } = useAppSelector(
    (state) => state.integrations,
  );

  const [connectingType, setConnectingType] = useState<IntegrationType | null>(
    null,
  );
  const [disconnectingId, setDisconnectingId] = useState<string | null>(null);
  const oauthNotifiedRef = useRef(false);

  const [integrationResponse, getIntegrations] = useAPI(
    api.integrations.getWorkspaceIntegrations,
  );

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
    if (oauthNotifiedRef.current) return;

    const success = searchParams.get("success");
    const error = searchParams.get("error");
    const type = searchParams.get("type");

    if (success === "true" && type) {
      oauthNotifiedRef.current = true;
      Notifier.success({
        title: "Integration Connected",
        message: `Successfully connected ${type.charAt(0).toUpperCase() + type.slice(1)} integration!`,
      });
      if (currentWorkspace?.id) {
        getIntegrations(currentWorkspace.id);
      }
    } else if (error) {
      oauthNotifiedRef.current = true;
      const errorMessages: Record<string, string> = {
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

  const connectedIntegrations: Integration[] =
    integrationResponse.result?.results ?? [];

  const handleConnect = async (type: IntegrationType) => {
    if (!currentWorkspace?.id) return;
    setConnectingType(type);
    try {
      const response = await api.integrations.setupIntegration(
        currentWorkspace.id,
        type,
        "/integrations",
      );
      window.location.href = response.data.authorize_url;
    } catch {
      Notifier.error({ title: "Failed to initiate connection" });
      setConnectingType(null);
    }
  };

  const handleDisconnect = async (integrationId: string) => {
    if (!currentWorkspace?.id) return;
    setDisconnectingId(integrationId);
    try {
      await api.integrations.deleteIntegration(
        currentWorkspace.id,
        integrationId,
      );
      Notifier.success({ title: "Integration disconnected" });
      getIntegrations(currentWorkspace.id);
    } catch {
      Notifier.error({ title: "Failed to disconnect integration" });
    } finally {
      setDisconnectingId(null);
    }
  };

  return (
    <>
      <PageHeader
        title="Integrations"
        description="Connect your workspace to external services and tools."
        className="space-y-8"
      />

      <div className="grid gap-4 mt-4">
        {(integrationTypes ?? []).map((integType) => {
          const connected = connectedIntegrations.find(
            (i) => i.type === integType.type,
          );
          const isConnecting = connectingType === integType.type;
          const isDisconnecting = connected
            ? disconnectingId === connected.id
            : false;

          return (
            <IntegrationCard
              key={integType.type}
              type={integType.type as IntegrationType}
              name={integType.name}
              description={integType.description}
              connected={connected}
              isConnecting={isConnecting}
              isDisconnecting={isDisconnecting}
              onConnect={() => handleConnect(integType.type as IntegrationType)}
              onDisconnect={() => connected && handleDisconnect(connected.id)}
            />
          );
        })}
      </div>
    </>
  );
}
