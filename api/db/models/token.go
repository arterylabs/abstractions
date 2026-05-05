package models

import (
	"time"

	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeEmailVerification TokenType = "email_verification"
	TokenTypePasswordReset     TokenType = "password_reset"
)

type Token struct {
	ID        uuid.UUID  `json:"id" gorm:"primaryKey;autoIncrement;"`
	UserID    uuid.UUID  `json:"user_id" gorm:"type:uuid"`
	Token     string     `json:"token" gorm:"type:text"`
	ExpiresAt time.Time  `json:"expires_at" gorm:"type:timestamp"`
	UsedAt    *time.Time `json:"used_at" gorm:"type:timestamp"`
	CreatedAt time.Time  `json:"created_at" gorm:"type:timestamp;autoCreateTime"`
	Type      TokenType  `json:"type" gorm:"type:text"`
}
