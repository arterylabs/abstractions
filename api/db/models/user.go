package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                 uuid.UUID  `json:"id" gorm:"primaryKey;autoIncrement;"`
	Email              string     `json:"email" gorm:"type:text"`
	Name               string     `json:"name" gorm:"type:text"`
	Avatar             *string    `json:"avatar" gorm:"type:text"`
	DefaultWorkspaceID *uuid.UUID `json:"default_workspace_id" gorm:"type:uuid;references:workspaces(id)"`
	CreatedAt          time.Time  `json:"created_at" gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt          time.Time  `json:"updated_at" gorm:"type:timestamp;autoUpdateTime"`
	VerifiedAt         *time.Time `json:"verified_at" gorm:"type:timestamp"`
	DeletedAt          *time.Time `json:"deleted_at" gorm:"type:timestamp;index"`
	DeletedBy          *uuid.UUID `json:"deleted_by" gorm:"type:uuid"`
	GodMode            bool       `json:"god_mode" gorm:"type:boolean"`
}
