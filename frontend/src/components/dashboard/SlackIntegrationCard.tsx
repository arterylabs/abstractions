import IntegrationCard from "./IntegrationCard";
import { Integration, IntegrationType, IntegrationTypeMetadata } from "@/types";

interface SlackIntegrationCardProps {
  integration?: Integration;
  isLoading: boolean;
  onConnect?: (redirectTo?: string) => Promise<void>;
  onIntegrationUpdated?: () => void | Promise<void>;
  meta?: IntegrationTypeMetadata;
}

export default function SlackIntegrationCard(props: SlackIntegrationCardProps) {
  return <IntegrationCard integrationType={IntegrationType.SLACK} {...props} />;
}
