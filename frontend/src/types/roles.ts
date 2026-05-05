export enum UserRole {
  ADMIN = "admin",
  OWNER = "owner",
  MEMBER = "member",
}
export type UserRoleType = keyof typeof UserRole;

export type Role = {
  id: string;
  name: string;
  description: string;
  permissions: string[];
};
