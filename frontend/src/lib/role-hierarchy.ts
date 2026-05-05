// Role hierarchy utility functions for frontend validation

export type WorkspaceRole = "owner" | "admin" | "member";

// Role hierarchy levels (higher number = more permissions)
function getRoleLevel(role: WorkspaceRole): number {
  switch (role) {
    case "owner":
      return 3;
    case "admin":
      return 2;
    case "member":
      return 1;
    default:
      return 0;
  }
}

/**
 * Checks if a user with requestingRole can grant targetRole to another user
 * A user can only grant roles that are strictly lower than their own
 */
export function canGrantRole(
  requestingRole: WorkspaceRole,
  targetRole: WorkspaceRole,
): boolean {
  const requestingLevel = getRoleLevel(requestingRole);
  const targetLevel = getRoleLevel(targetRole);

  // Owner (level 3) can grant Admin (level 2) and Member (level 1)
  // Admin (level 2) can grant Member (level 1) only
  // Member (level 1) cannot grant any roles
  return requestingLevel > targetLevel;
}

/**
 * Checks if a user with requestingRole can modify the role of userRole
 * A user can only modify users with roles strictly lower than their own
 */
export function canModifyUserRole(
  requestingRole: WorkspaceRole,
  userRole: WorkspaceRole,
): boolean {
  const requestingLevel = getRoleLevel(requestingRole);
  const userLevel = getRoleLevel(userRole);

  // Owner can modify Admin and Member
  // Admin can modify Member only
  // Member cannot modify anyone
  return requestingLevel > userLevel;
}

/**
 * Returns the available roles that a user can grant to others
 */
export function getGrantableRoles(
  requestingRole: WorkspaceRole,
): WorkspaceRole[] {
  const allRoles: WorkspaceRole[] = ["owner", "admin", "member"];
  return allRoles.filter((role) => canGrantRole(requestingRole, role));
}

/**
 * Checks if a role change operation is valid
 */
export function isValidRoleChange(
  requestingRole: WorkspaceRole,
  currentUserRole: WorkspaceRole,
  newRole: WorkspaceRole,
): { valid: boolean; reason?: string } {
  // Can't modify someone with equal or higher role
  if (!canModifyUserRole(requestingRole, currentUserRole)) {
    return {
      valid: false,
      reason: "You don't have permission to modify this user's role",
    };
  }

  // Can't grant a role higher than or equal to your own
  if (!canGrantRole(requestingRole, newRole)) {
    return {
      valid: false,
      reason: "You don't have permission to assign this role",
    };
  }

  return { valid: true };
}
