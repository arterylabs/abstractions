import { JSONObject } from "./user";

export enum IntegrationType {
  GITHUB = "github",
}

export interface Integration {
  id: string;
  type: IntegrationType;
  serviceTitle?: string;
  serviceIcon?: string;
  workspaceId: string;
  tokenExpiresAt?: string;
  scopes?: string;
  serviceTenantId?: string;
  serviceUserId?: string;
  resourceType?: string;
  resourceId?: string;
  configuration: JSONObject;
  createdAt: string;
  updatedAt: string;
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
