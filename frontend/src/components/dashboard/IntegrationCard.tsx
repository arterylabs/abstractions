/* eslint-disable @next/next/no-img-element */
"use client";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import Spinner from "@/components/ui/spinner";
import { SimpleDropdownMenu } from "@/components/ui/dropdown-menu";
import Image from "next/image";
import {
  Integration,
  getIntegrationIcon,
  IntegrationType,
  IntegrationTypeMetadata,
} from "@/types";
import { usePathname } from "next/navigation";
import Timestamp from "../ui/timestamp";
import Icon from "../ui/icon";
import { askConfirmation } from "@/components/ui/confirmation";
import api from "@/api";
import { Notifier } from "@/components/ui/sonner";
import { useAppSelector } from "@/store/hooks";

interface IntegrationCardProps {
  integrationType: IntegrationType;
  integration?: Integration;
  isLoading: boolean;
  onIntegrationUpdated?: () => void;
  meta?: IntegrationTypeMetadata;
  onConnect?: (redirectTo?: string) => void; // optional custom connect handler, falls back to setupIntegration
}

export default function IntegrationCard({
  integrationType,
  integration,
  isLoading,
  onIntegrationUpdated,
  meta,
  onConnect,
}: IntegrationCardProps) {
  const pathname = usePathname();
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );

  const handleConnect = async () => {
    if (onConnect) return onConnect(pathname);
    if (!currentWorkspace?.id) return;
    try {
      const response = await api.integrations.setupIntegration(
        currentWorkspace.id,
        integrationType,
        pathname,
      );
      window.location.href = response.data.authorize_url;
    } catch (error) {
      console.error("Failed to setup integration:", error);
      Notifier.error({
        title: "Error",
        message: `Failed to connect ${meta?.name || "integration"}.`,
      });
    }
  };

  const handleDelete = async () => {
    if (!integration || !currentWorkspace?.id) return;

    askConfirmation({
      title: `Delete ${meta?.name || "Integration"}`,
      content: `Are you sure you want to delete this ${
        meta?.name || "integration"
      }? This action cannot be undone. You'll need to reconnect to set up the integration again.`,
      variant: "destructive",
      okText: "Delete",
      cancelText: "Cancel",
      okCallback: async () => {
        try {
          await api.integrations.deleteIntegration(
            currentWorkspace.id,
            integration.id,
          );
          Notifier.success({
            title: "Integration deleted",
            message: `${meta?.name || "Integration"} deleted successfully.`,
          });
          onIntegrationUpdated?.();
        } catch (error) {
          console.error("Failed to delete integration:", error);
          Notifier.error({
            title: "Error",
            message: `Failed to delete ${
              meta?.name || "integration"
            }. Please try again.`,
          });
        }
      },
    });
  };

  const handleReconnect = async () => {
    if (!currentWorkspace?.id) return;

    askConfirmation({
      title: `Reconnect ${meta?.name || "Integration"}`,
      content: `This will disconnect your current ${
        meta?.name || "integration"
      } connection and allow you to connect a new one. Do you want to continue?`,
      variant: "primary",
      okText: "Reconnect",
      cancelText: "Cancel",
      okCallback: async () => {
        try {
          if (integration) {
            const response = await api.integrations.setupIntegration(
              currentWorkspace.id,
              integrationType,
              pathname,
            );
            // Trigger a refresh so parent re-fetches integrations list immediately
            onIntegrationUpdated?.();
            window.location.href = response.data.authorize_url;
          }
        } catch (error) {
          console.error("Failed to reconnect integration:", error);
          Notifier.error({
            title: "Error",
            message: `Failed to reconnect ${
              meta?.name || "integration"
            }. Please try again.`,
          });
        }
      },
    });
  };

  if (isLoading) {
    return (
      <div className="text-sm text-muted-foreground">
        {meta?.name || "Integration"}
        <Card className="mt-2">
          <CardContent className="flex items-center justify-center p-6">
            <Spinner />
          </CardContent>
        </Card>
      </div>
    );
  }

  if (!integration) {
    return (
      <Card className="mt-2">
        <CardContent className="flex items-start justify-between p-6">
          <div className="space-y-4 w-full">
            {meta?.icon && (
              <Image
                src={meta.icon}
                alt={meta?.name || "Integration"}
                width={32}
                height={32}
              />
            )}
            <div className="space-y-2">
              <span className="text-lg font-semibold">
                {meta?.name || integrationType}
              </span>
              {meta?.description && (
                <p className="text-sm text-muted-foreground">
                  {meta.description}
                </p>
              )}
            </div>
          </div>
          <Button onClick={handleConnect}>Connect</Button>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="mt-2">
      <CardContent className="flex items-start justify-between p-6">
        <div className="space-y-4 w-full">
          <Image
            src={getIntegrationIcon(integration.type)}
            alt={meta?.name || "Integration Logo"}
            width={32}
            height={32}
          />
          <span className="text-lg font-semibold flex items-center gap-2">
            {integration.service_icon ? (
              <img
                src={integration.service_icon}
                className="h-4 w-4 rounded"
                alt={`${meta?.name || "Integration"} Instance`}
              />
            ) : null}
            {integration.service_title || meta?.name || integrationType}
          </span>
          <div className="flex items-center gap-2">
            <span className="text-xs text-muted-foreground">
              Connected <Timestamp timestamp={integration.created_at} />
            </span>
          </div>
        </div>
        <Badge variant="success" className="mr-2">
          Connected
        </Badge>
        <SimpleDropdownMenu
          options={[
            {
              label: "Reconnect",
              onClick: handleReconnect,
              icon: <Icon type="refresh" className="mr-2 h-4 w-4" />,
              variant: "default",
            },
            {
              label: "Delete",
              onClick: handleDelete,
              icon: <Icon type="delete" className="mr-2 h-4 w-4" />,
              variant: "destructive",
            },
          ]}
        >
          <Button variant="outline" size="sm">
            <Icon type="more_vert" />
          </Button>
        </SimpleDropdownMenu>
      </CardContent>
    </Card>
  );
}
