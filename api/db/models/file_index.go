package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// StringArray is a custom type for PostgreSQL text[] columns.
type StringArray []string

func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return "{}", nil
	}
	result := "{"
	for i, v := range s {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf("%q", v)
	}
	result += "}"
	return result, nil
}

func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = StringArray{}
		return nil
	}
	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	// Parse PostgreSQL array literal {a,b,c}
	if len(str) < 2 || str[0] != '{' || str[len(str)-1] != '}' {
		*s = StringArray{}
		return nil
	}
	inner := str[1 : len(str)-1]
	if inner == "" {
		*s = StringArray{}
		return nil
	}
	var result []string
	if err := json.Unmarshal([]byte(`[`+inner+`]`), &result); err != nil {
		// fallback: simple comma split
		result = splitPGArray(inner)
	}
	*s = result
	return nil
}

// splitPGArray splits a raw PostgreSQL array interior on commas, stripping quotes.
func splitPGArray(s string) []string {
	var result []string
	current := ""
	inQuote := false
	for _, c := range s {
		switch {
		case c == '"':
			inQuote = !inQuote
		case c == ',' && !inQuote:
			result = append(result, current)
			current = ""
		default:
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

// TODO - move this to a shared package we will need it in multiple places
// JSONB is a custom type for PostgreSQL jsonb columns backed by json.RawMessage.
type JSONB json.RawMessage

func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return string(j), nil
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*j = append((*j)[0:0], v...)
	case string:
		*j = JSONB(v)
	default:
		return fmt.Errorf("unsupported type for JSONB: %T", value)
	}
	return nil
}

func (j JSONB) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return []byte("null"), nil
	}
	return json.RawMessage(j).MarshalJSON()
}

func (j *JSONB) UnmarshalJSON(data []byte) error {
	*j = append((*j)[0:0], data...)
	return nil
}

// RepositoryFile represents an indexed file within a connected repository.
type RepositoryFile struct {
	ID              uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	RepoID          uuid.UUID   `json:"repoId" gorm:"type:uuid;not null"`
	Path            string      `json:"path" gorm:"column:path;not null"`
	GitSHA          string      `json:"gitSha" gorm:"column:git_sha;not null;default:''"`
	Tags            StringArray `json:"tags" gorm:"column:tags;type:text[];not null;default:'{}'"`
	ProcessedAt     *time.Time  `json:"processedAt" gorm:"column:processed_at"`
	Size            int         `json:"size" gorm:"column:size;not null;default:0"`
	Language        string      `json:"language" gorm:"column:language;not null;default:''"`
	LastCommitAt    *time.Time  `json:"lastCommitAt" gorm:"column:last_commit_at"`
	ServiceMetadata JSONB       `json:"serviceMetadata" gorm:"column:service_metadata;type:jsonb"`
	DeletedAt       *time.Time  `json:"deletedAt" gorm:"type:timestamptz;column:deleted_at"`
	DeletedBy       *uuid.UUID  `json:"deletedBy" gorm:"type:uuid;column:deleted_by"`
	CreatedAt       time.Time   `json:"createdAt"`
	UpdatedAt       time.Time   `json:"updatedAt"`
}

func (RepositoryFile) TableName() string {
	return "repository_files"
}
