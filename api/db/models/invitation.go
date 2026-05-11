package models

import (
	"time"

	"github.com/google/uuid"
)

type InvitationState string

const (
	InvitationStatePending  InvitationState = "pending"
	InvitationStateExpired  InvitationState = "expired"
	InvitationStateAccepted InvitationState = "accepted"
	InvitationStateIgnored  InvitationState = "ignored"
)

func (s InvitationState) String() string {
	return string(s)
}

type Invitation struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	WorkspaceID   uuid.UUID  `json:"workspace_id" gorm:"type:uuid"`
	InvitedBy     uuid.UUID  `json:"invited_by" gorm:"type:uuid"`
	Email         string     `json:"email" gorm:"type:text"`
	Role          string     `json:"role" gorm:"type:text"`
	ExpiresAt     time.Time  `json:"expires_at" gorm:"type:timestamp"`
	AcceptedAt    *time.Time `json:"accepted_at,omitempty" gorm:"type:timestamp"`
	IgnoredAt     *time.Time `json:"ignored_at,omitempty" gorm:"type:timestamp"`
	CreatedAt     time.Time  `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"type:timestamp"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" gorm:"type:timestamptz"`
	DeletedBy     *uuid.UUID `json:"deleted_by,omitempty" gorm:"type:uuid"`
	InvitedByUser *User      `gorm:"-" json:"invitedByUser,omitempty"`
	Workspace     *Workspace `gorm:"-" json:"workspace,omitempty"`
}

func (i *Invitation) GetState() InvitationState {
	if i.AcceptedAt != nil {
		return InvitationStateAccepted
	}
	if i.IgnoredAt != nil {
		return InvitationStateIgnored
	}
	if i.ExpiresAt.Before(time.Now()) {
		return InvitationStateExpired
	}
	return InvitationStatePending
}

func (i *Invitation) IsPendingOrExpired() bool {
	state := i.GetState()
	return state == InvitationStatePending || state == InvitationStateExpired
}
