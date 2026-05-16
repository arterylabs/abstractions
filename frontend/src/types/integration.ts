import { JSONObject } from "./user";

export enum IntegrationType {
  GITHUB = "github",
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

export interface IntegrationTypeMetadata {
  type: IntegrationType;
  name: string;
  icon: string;
  description: string;
}

export type IntegrationConfiguration = Record<string, unknown>;

export const getIntegrationDisplayName = (type: IntegrationType): string => {
  switch (type) {
    case IntegrationType.GITHUB:
      return "GitHub";
    default:
      return type;
  }
};

export const getIntegrationIcon = (type: IntegrationType): string => {
  switch (type) {
    case IntegrationType.GITHUB:
      return "/assets/img/github.svg";
    default:
      return "/assets/img/integration-default.svg";
  }
};

export const isValidIntegrationType = (
  type: string,
): type is IntegrationType => {
  return Object.values(IntegrationType).includes(type as IntegrationType);
};
