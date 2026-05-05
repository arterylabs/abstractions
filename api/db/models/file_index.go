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

// FileIndex represents an indexed file within a connected repository.
type FileIndex struct {
	ID              uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	RepoID          uuid.UUID   `json:"repoId" gorm:"type:uuid;not null"`
	FilePath        string      `json:"filePath" gorm:"not null"`
	GitSHA          string      `json:"gitSha" gorm:"column:git_sha;not null;default:''"`
	DomainTags      StringArray `json:"domainTags" gorm:"column:domain_tags;type:text[]"`
	LastProcessedAt *time.Time  `json:"lastProcessedAt"`
	PineconeIDs     StringArray `json:"pineconeIds" gorm:"column:pinecone_ids;type:text[]"`
	CreatedAt       time.Time   `json:"createdAt"`
	UpdatedAt       time.Time   `json:"updatedAt"`
}

func (FileIndex) TableName() string {
	return "file_index"
}
