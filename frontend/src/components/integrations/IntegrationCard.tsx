"use client";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import {
  Integration,
  IntegrationType,
  getIntegrationIcon,
} from "@/types/integration";

interface IntegrationCardProps {
  type: IntegrationType;
  name: string;
  description: string;
  connected?: Integration;
  isConnecting: boolean;
  isDisconnecting: boolean;
  onConnect: () => void;
  onDisconnect: () => void;
}

export function IntegrationCard({
  type,
  name,
  description,
  connected,
  isConnecting,
  isDisconnecting,
  onConnect,
  onDisconnect,
}: IntegrationCardProps) {
  return (
    <Card>
      <CardHeader className="pb-3">
        <div className="flex items-start justify-between gap-4">
          <div className="flex items-center gap-3">
            <img
              src={getIntegrationIcon(type)}
              alt={name}
              className="h-8 w-8 shrink-0"
            />
            <div>
              <CardTitle className="text-base">{name}</CardTitle>
              <CardDescription className="mt-0.5">
                {description}
              </CardDescription>
            </div>
          </div>

          {!connected ? (
            <Button
              size="sm"
              onClick={onConnect}
              loading={isConnecting}
              disabled={isConnecting}
            >
              Connect
            </Button>
          ) : (
            <div className="flex items-center gap-2 shrink-0">
              <Button
                size="sm"
                variant="outline"
                onClick={onConnect}
                loading={isConnecting}
                disabled={isConnecting || isDisconnecting}
              >
                Reconnect
              </Button>
              <AlertDialog>
                <AlertDialogTrigger asChild>
                  <Button
                    size="sm"
                    variant="destructive"
                    disabled={isDisconnecting || isConnecting}
                    loading={isDisconnecting}
                  >
                    Disconnect
                  </Button>
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>Disconnect {name}?</AlertDialogTitle>
                    <AlertDialogDescription>
                      This will remove the {name} integration from your
                      workspace. Features that depend on it may stop working
                      until you reconnect.
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                    <AlertDialogAction
                      className="bg-destructive hover:bg-destructive/90"
                      onClick={onDisconnect}
                    >
                      Disconnect
                    </AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            </div>
          )}
        </div>
      </CardHeader>

      {connected && (
        <CardContent className="pt-0">
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            {connected.serviceIcon && (
              <img
                src={connected.serviceIcon}
                alt={connected.serviceTitle ?? ""}
                className="h-5 w-5 rounded-full"
              />
            )}
            <span>
              Connected as{" "}
              <span className="font-medium text-foreground">
                {connected.serviceTitle}
              </span>
            </span>
          </div>
        </CardContent>
      )}
    </Card>
  );
}
