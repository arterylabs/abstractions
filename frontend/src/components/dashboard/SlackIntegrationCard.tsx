// This file has been deprecated in favor of the generic IntegrationCard component.
// Keeping a thin wrapper for backward compatibility. Prefer importing IntegrationCard directly.

import IntegrationCard from "./IntegrationCard";
import { Integration, IntegrationType, IntegrationTypeMetadata } from "@/types";

interface SlackIntegrationCardProps {
  integration?: Integration;
  isLoading: boolean;
  onConnect?: (redirectTo?: string) => void; // kept for backward compatibility but generic card auto-handles connect
  onIntegrationUpdated?: () => void;
  meta?: IntegrationTypeMetadata;
}

export default function SlackIntegrationCard(props: SlackIntegrationCardProps) {
  return (
    <IntegrationCard
      integrationType={IntegrationType.SLACK}
      integration={props.integration}
      isLoading={props.isLoading}
      onIntegrationUpdated={props.onIntegrationUpdated}
      meta={props.meta}
      onConnect={props.onConnect}
    />
  );
}
