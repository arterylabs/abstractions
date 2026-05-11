package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID             uuid.UUID  `json:"id" gorm:"primaryKey;autoIncrement;"`
	UserID         uuid.UUID  `json:"user_id" gorm:"type:uuid"`
	Token          string     `json:"token" gorm:"type:text"`
	ExpiresAt      time.Time  `json:"expires_at" gorm:"type:timestamp"`
	RevokedAt      *time.Time `json:"revoked_at" gorm:"type:timestamp"`
	CreatedAt      time.Time  `json:"created_at" gorm:"type:timestamp;autoCreateTime"`
	UpdatedBy      *uuid.UUID `json:"updated_by" gorm:"type:uuid"`
	UserAgent      string     `json:"user_agent" gorm:"type:text"`
	IP             string     `json:"ip" gorm:"type:text"`
	ImpersonatedAt *time.Time `json:"impersonated_at" gorm:"type:timestamp"`
	ImpersonatedBy *uuid.UUID `json:"impersonated_by" gorm:"type:uuid"`
}

func (s *Session) IsValid() bool {
	now := time.Now()
	// Check if session is not expired and not invalidated
	return s.ExpiresAt.After(now) && s.RevokedAt == nil
}
