package authorization_service

type Permission string

func (p Permission) String() string {
	return string(p)
}

// Workspace permissions
const (
	WorkspaceUpdate Permission = "workspaces:update"
	WorkspaceDelete Permission = "workspaces:delete"

	// Workspace user management
	WorkspaceManageUsers Permission = "workspace:users:manage"

	// Workspace settings (logo, subscription, etc)
	WorkspaceManageSubscription Permission = "workspace:subscription:manage"

	WorkspaceManageIntegrations Permission = "workspace:integrations:manage"
	WorkspaceManageRepositories Permission = "workspace:repositories:manage"
)
