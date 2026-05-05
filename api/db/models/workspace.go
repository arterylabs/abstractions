package models

import (
	"time"

	"github.com/google/uuid"
)

type Workspace struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name         string     `json:"name"`
	Logo         *string    `json:"logo"`
	Website      *string    `json:"website"`
	About        *string    `json:"about"`
	BillingEmail string     `json:"billing_email"`
	CreatedBy    uuid.UUID  `json:"created_by"`
	UpdatedBy    uuid.UUID  `json:"updated_by"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
	DeletedBy    *uuid.UUID `json:"deleted_by"`
}

type UserWorkspaceRelationship struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID      uuid.UUID `json:"user_id"`
	WorkspaceID uuid.UUID `json:"workspace_id"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
