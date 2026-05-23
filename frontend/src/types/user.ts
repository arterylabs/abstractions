import { UserRole } from "./roles";

export interface ImpersonationInfo {
  isImpersonated: boolean;
  impersonatedBy?: string;
  impersonatedAt?: string;
  impersonatorUser?: ImpersonatorUser;
}

export interface ImpersonatorUser {
  id: string;
  name: string;
  email: string;
  avatar?: string;
}

export type User = {
  id: string;
  email: string;
  createdAt: string;
  updatedAt: string;
  name: string;
  avatar?: string;
  verifiedAt?: string;
  deletedAt?: string;
  deletedBy?: string;
  deletedByName?: string;
  defaultWorkspaceId?: string;
  defaultWorkspaceName?: string;
  workspacesCount?: number;
  authenticators?: AuthenticatorType[];
  workspaceRole?: string;
  godMode?: boolean;
  impersonation?: ImpersonationInfo;
};

export interface WorkspaceUser extends User {
  role: UserRole; // Assuming role is a string, you might want to define a more specific type
}

export type Token = {
  id: string;
  userId: string;
  token: string;
  expiresAt: string;
  createdAt: string;
  type: TokenType;
};

export type TokenType = "email_verification" | "password_reset";

export type Subscription = {
  id: string;
  userId: string;
  planId: SubscriptionPlanId;
  status: SubscriptionStatus;
  expiresAt: string;
  createdAt: string;
  updatedAt: string;
  deletedAt?: string;
  deletedBy?: string;
  limit?: number;
};

export type SubscriptionStatus = "active" | "inactive" | "trialing";
export type SubscriptionPlanId = "lifetime" | "free";

export type Session = {
  id: string;
  userId: string;
  token: string;
  expiresAt: string;
  createdAt: string;
  updatedBy?: string;
  userAgent: string;
  ip: string;
  impersonatedAt?: string;
  impersonatedBy?: string;
};

export type OnboardingGuide = {
  id: string;
  userId: string;
  steps: JSONObject; // Assuming steps is a JSONB type, you might want to define a more specific type
};

export type Authenticator = {
  id: string;
  userId: string;
  type: AuthenticatorType;
  token: string;
  refreshToken?: string;
  scopes: string;
  createdAt: string;
  updatedBy?: string;
};

export type JSONObject = Record<string, unknown>;

export type AuthenticatorType = "google" | "github" | "password";

export type InvitationResponseStatus = "success" | "error" | "partial_success";

export interface InvitationResult {
  email: string;
  role: string;
  success: boolean;
  error?: string;
}

export interface SendInvitationsResponse {
  status: InvitationResponseStatus;
  results: InvitationResult[];
}
