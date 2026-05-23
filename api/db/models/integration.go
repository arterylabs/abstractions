package models

import (
	"abstractions/services/integrations/enums"
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
	ResourceType    *string                   `json:"resource_type"`
	ResourceID      *uuid.UUID                `json:"resource_id"`
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

// SetConfiguration marshals a configuration struct into the Configuration field.
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

// GetConfiguration unmarshals the Configuration field into the provided struct.
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
