package authorization_service

type RoleID string

var RoleOwner RoleID = "owner"
var RoleAdmin RoleID = "admin"
var RoleMember RoleID = "member"

var DEFAULT_ROLE = RoleOwner

func (r RoleID) String() string {
	return string(r)
}

type Role struct {
	ID          RoleID       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions"`
}

var OWNER_ROLE = Role{
	ID:          RoleOwner,
	Name:        "Owner",
	Description: "Full access to all workspace features and settings.",
	Permissions: []Permission{
		WorkspaceUpdate,
		WorkspaceDelete,
		WorkspaceManageUsers,
		WorkspaceManageIntegrations,
		WorkspaceManageRepositories,
	},
}

var ADMIN_ROLE = Role{
	ID:          RoleAdmin,
	Name:        "Admin",
	Description: "Manage workspace users and settings, but cannot delete the workspace.",
	Permissions: []Permission{
		WorkspaceUpdate,
		WorkspaceManageUsers,
		WorkspaceManageIntegrations,
		WorkspaceManageRepositories,
	},
}

var MEMBER_ROLE = Role{
	ID:          RoleMember,
	Name:        "Member",
	Description: "Basic access to workspace features, can read settings and accept invitations.",
	Permissions: []Permission{},
}

var ALL_ROLES = []Role{
	OWNER_ROLE,
	ADMIN_ROLE,
	MEMBER_ROLE,
}

// GetRoleByID returns the role with the given ID, or nil if not found.
func GetRoleByID(roleID string) *Role {
	for _, role := range ALL_ROLES {
		if role.ID.String() == roleID {
			return &role
		}
	}
	return nil
}

// GetRoleByName returns the role with the given name, or nil if not found.
func GetRoleByName(roleName string) *Role {
	for _, role := range ALL_ROLES {
		if role.Name == roleName {
			return &role
		}
	}
	return nil
}

func RolePermissions() map[RoleID][]Permission {
	permissionsMap := make(map[RoleID][]Permission)
	for _, role := range ALL_ROLES {
		permissionsMap[role.ID] = role.Permissions
	}
	return permissionsMap
}

func AllRoleIDsInterface() []interface{} {
	ids := make([]interface{}, len(ALL_ROLES))
	for i, role := range ALL_ROLES {
		ids[i] = role.ID.String()
	}
	return ids
}

// Role hierarchy levels (higher number = more permissions)
func getRoleLevel(roleID string) int {
	switch RoleID(roleID) {
	case RoleOwner:
		return 3
	case RoleAdmin:
		return 2
	case RoleMember:
		return 1
	default:
		return 0
	}
}

// CanGrantRole checks if a user with requestingRole can grant targetRole to another user
func CanGrantRole(requestingRole, targetRole string) bool {
	requestingLevel := getRoleLevel(requestingRole)
	targetLevel := getRoleLevel(targetRole)

	// Members cannot grant any roles
	if requestingLevel <= 1 {
		return false
	}

	// A user can grant roles up to their own level (but not higher)
	// Owner (level 3) can grant Owner (level 3), Admin (level 2) and Member (level 1)
	// Admin (level 2) can grant Admin (level 2) and Member (level 1)
	// Member (level 1) cannot grant any roles
	return requestingLevel >= targetLevel
}

// CanModifyUserRole checks if a user with requestingRole can modify the role of userRole
func CanModifyUserRole(requestingRole, userRole string) bool {
	requestingLevel := getRoleLevel(requestingRole)
	userLevel := getRoleLevel(userRole)

	// A user can only modify users with roles strictly lower than their own
	// Owner can modify Admin and Member
	// Admin can modify Member only
	// Member cannot modify anyone
	return requestingLevel > userLevel
}
