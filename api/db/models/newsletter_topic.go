package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// JSONBRaw stores any JSON value (object, array, scalar) as raw bytes in a JSONB column.
type JSONBRaw json.RawMessage

func (j JSONBRaw) Value() (driver.Value, error) {
	if j == nil {
		return "null", nil
	}
	return string(j), nil
}

func (j *JSONBRaw) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	var b []byte
	switch v := value.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("JSONBRaw: unsupported type %T", value)
	}
	raw := make(JSONBRaw, len(b))
	copy(raw, b)
	*j = raw
	return nil
}

// MarshalJSON returns the raw JSON bytes.
func (j JSONBRaw) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return json.RawMessage(j).MarshalJSON()
}

// NewsletterTopic represents a queued or sent topic associated with a repository newsletter.
type NewsletterTopic struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	RepoID       uuid.UUID  `json:"repoId" gorm:"type:uuid;not null;column:repo_id"`
	CreatedBy    *uuid.UUID `json:"createdBy" gorm:"type:uuid;column:created_by"`
	UpdatedBy    *uuid.UUID `json:"updatedBy" gorm:"type:uuid;column:updated_by"`
	Title        string     `json:"title" gorm:"not null;default:''"`
	Description  string     `json:"description" gorm:"not null;default:''"`
	RelatedFiles JSONBRaw   `json:"relatedFiles" gorm:"type:jsonb;not null;default:'[]'"`
	SentAt       *time.Time `json:"sentAt" gorm:"type:timestamptz;column:sent_at"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	DeletedAt    *time.Time `json:"deletedAt" gorm:"type:timestamptz"`
	DeletedBy    *uuid.UUID `json:"deletedBy" gorm:"type:uuid"`
}
