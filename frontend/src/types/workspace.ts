export type Workspace = {
  id: string;
  createdAt: string;
  updatedAt: string;
  name: string;
  logo?: string | null;
  website?: string | null;
  about?: string | null;
  billingEmail?: string;
  userRole?: string;
};
