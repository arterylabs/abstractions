package models

import (
	"time"

	"github.com/google/uuid"
)

type WorkspaceConfiguration struct {
	ID          uuid.UUID `json:"id"          gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	WorkspaceID uuid.UUID `json:"workspaceId"  gorm:"column:workspace_id"`
	Key         string    `json:"key"          gorm:"column:key"`
	Value       string    `json:"-"            gorm:"column:value"`
	Encrypted   bool      `json:"encrypted"    gorm:"column:encrypted"`
	CreatedAt   time.Time `json:"createdAt"    gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updatedAt"    gorm:"column:updated_at"`
	CreatedBy   uuid.UUID `json:"createdBy"    gorm:"column:created_by"`
	UpdatedBy   uuid.UUID `json:"updatedBy"    gorm:"column:updated_by"`
}
