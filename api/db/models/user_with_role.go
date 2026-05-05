package models

// UserWithRole extends User to include the role information from UserWorkspaceRelationship
type UserWithRole struct {
	User
	Role string `json:"role" gorm:"column:role"`
}
