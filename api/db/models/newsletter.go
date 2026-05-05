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

// Newsletter represents a generated and (optionally) sent newsletter for a repository.
type Newsletter struct {
	ID             uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	RepoID         uuid.UUID  `json:"repoId" gorm:"type:uuid;not null"`
	Topics         JSONBRaw   `json:"topics" gorm:"type:jsonb;not null;default:'[]'"`
	HTMLContent    string     `json:"htmlContent" gorm:"column:html_content;not null;default:''"`
	SentAt         *time.Time `json:"sentAt"`
	RecipientCount int        `json:"recipientCount" gorm:"column:recipient_count;not null;default:0"`
	OpenCount      int        `json:"openCount" gorm:"column:open_count;not null;default:0"`
	ClickCount     int        `json:"clickCount" gorm:"column:click_count;not null;default:0"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}
