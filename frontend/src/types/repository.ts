export interface Repository {
  id: string;
  workspaceId: string;
  userId: string;
  githubRepoId: number;
  fullName: string;
  defaultBranch: string;
  webhookId: number | null;
  status: "pending" | "indexing" | "indexed" | "error";
  indexedAt: string | null;
  sendDay: number;
  sendHour: number;
  lastNewsletterSentAt: string | null;
  createdAt: string;
  updatedAt: string;
}
