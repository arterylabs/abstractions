package models

// WorkspaceWithRole extends Workspace to include the user role information from UserWorkspaceRelationship
type WorkspaceWithRole struct {
	Workspace
	UserRole string `json:"userRole" gorm:"column:role"`
}
