import { JSONObject } from "./user";

// Integration type enum matching the backend
export enum IntegrationType {
  SLACK = "slack",
  // Add more integration types as they become available
  // DISCORD = "discord",
  // TEAMS = "teams",
  // etc.
}

export interface Integration {
  id: string;
  type: IntegrationType;
  service_title?: string;
  service_icon?: string;
  workspace_id: string;
  token_expires_at?: string;
  scopes?: string;
  service_tenant_id?: string;
  service_user_id?: string;
  resourceType?: string;
  resourceId?: string;
  configuration: JSONObject;
  created_at: string;
  updated_at: string;
}

// Metadata for supported integration types
export interface IntegrationTypeMetadata {
  type: IntegrationType;
  name: string;
  icon: string; // URL or path to icon
  description: string;
}

// Slack-specific configuration interface
export interface SlackConfiguration {
  team_name?: string;
  team_id?: string;
  user_id?: string;
  user_name?: string;
  scopes?: string[];
  away_message?: string;
}

// Future integration configurations can be added here
// export interface DiscordConfiguration {
//   guild_name?: string;
//   guild_id?: string;
//   // ... other Discord-specific fields
// }

// Union type for all configuration types
export type IntegrationConfiguration = SlackConfiguration;

// Helper functions
export const getIntegrationDisplayName = (type: IntegrationType): string => {
  switch (type) {
    case IntegrationType.SLACK:
      return "Slack";
    default:
      return type;
  }
};

export const getIntegrationIcon = (type: IntegrationType): string => {
  switch (type) {
    case IntegrationType.SLACK:
      return "/assets/img/slack.svg";
    default:
      return "/assets/img/integration-default.svg";
  }
};

export const isValidIntegrationType = (
  type: string,
): type is IntegrationType => {
  return Object.values(IntegrationType).includes(type as IntegrationType);
};
