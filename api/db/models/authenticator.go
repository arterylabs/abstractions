package models

import (
	"time"

	"github.com/google/uuid"
)

type AuthenticatorType string

const (
	AuthenticatorTypeGoogle   AuthenticatorType = "google"
	AuthenticatorTypePassword AuthenticatorType = "password"
)

func (a AuthenticatorType) String() string {
	return string(a)
}

type Authenticator struct {
	ID           uuid.UUID         `json:"id" gorm:"primaryKey;autoIncrement;"`
	UserID       uuid.UUID         `json:"user_id" gorm:"type:uuid"`
	Type         AuthenticatorType `json:"type" gorm:"type:text"`
	Token        string            `json:"token" gorm:"type:text"`
	RefreshToken *string           `json:"refresh_token" gorm:"type:text"`
	Scopes       string            `json:"scopes" gorm:"type:text"`
	CreatedAt    time.Time         `json:"created_at" gorm:"type:timestamp;autoCreateTime"`
	UpdatedBy    *uuid.UUID        `json:"updated_by" gorm:"type:uuid"`
}
