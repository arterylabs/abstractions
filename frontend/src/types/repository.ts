export interface Repository {
  id: string;
  workspaceId: string;
  userId: string;
  serviceType?: string;
  serviceId?: string;
  name?: string;
  fullName: string;
  description: string;
  defaultBranch: string;
  status: "pending" | "indexing" | "indexed" | "error";
  indexedAt: string | null;
  indexError?: string | null;
  sendDay: number;
  sendHour: number;
  newsletterEnabled: boolean;
  lastNewsletterSentAt: string | null;
  emails: string[];
  sendToWorkspaceMembers: boolean;
  createdAt: string;
  updatedAt: string;
}
