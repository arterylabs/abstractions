package models

import (
	"abstractions/services/integrations/enums"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Integration struct {
	ID              uuid.UUID                 `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Type            enums.IntegrationTypeEnum `json:"type"`
	ServiceTitle    *string                   `json:"service_title"`
	ServiceIcon     *string                   `json:"service_logo"`
	WorkspaceID     uuid.UUID                 `json:"workspace_id"`
	Token           string                    `json:"token"`
	RefreshToken    *string                   `json:"refresh_token"`
	TokenExpiresAt  *time.Time                `json:"token_expires_at"`
	Scopes          *string                   `json:"scopes"`
	ServiceTenantID *string                   `json:"service_tenant_id"`
	ServiceUserID   *string                   `json:"service_user_id"`
	Configuration   JSONBField                `json:"configuration" gorm:"type:jsonb"`
	CreatedBy       uuid.UUID                 `json:"created_by"`
	UpdatedBy       uuid.UUID                 `json:"updated_by"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
	DeletedAt       *time.Time                `json:"deleted_at"`
	DeletedBy       *uuid.UUID                `json:"deleted_by"`
}

// JSONBField is a custom type to handle JSONB fields
type JSONBField map[string]interface{}

// Value implements the driver.Valuer interface for database storage
func (j JSONBField) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface for database retrieval
func (j *JSONBField) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	result := make(map[string]interface{})
	err := json.Unmarshal(bytes, &result)
	*j = result
	return err
}

// SetConfiguration marshals a configuration struct into the Configuration field
func (i *Integration) SetConfiguration(config interface{}) error {
	if config == nil {
		i.Configuration = nil
		return nil
	}

	bytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}

	i.Configuration = result
	return nil
}

// GetConfiguration unmarshals the Configuration field into the provided struct
func (i *Integration) GetConfiguration(config interface{}) error {
	if i.Configuration == nil {
		return nil
	}

	bytes, err := json.Marshal(i.Configuration)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, config)
}

// Example configuration structs for different integration types

// SlackConfiguration represents configuration for Slack integrations
type SlackConfiguration struct {
	AwayMessage string `json:"away_message"`
}
