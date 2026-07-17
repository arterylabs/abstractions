package enums

type IntegrationTypeEnum string

func FromString(s string) (IntegrationTypeEnum, error) {
	return IntegrationTypeEnum(s), nil
}

func (it IntegrationTypeEnum) String() string {
	return string(it)
}

type AuthType string

const (
	AuthTypeOAuth2    AuthType = "oauth2"
	AuthTypeAPIKey    AuthType = "api_key"
	AuthTypeBasicAuth AuthType = "basic_auth"
	AuthTypeNone      AuthType = "none"
)

type EndpointDetails struct {
	Path    string `json:"path"`    // e.g., /api/v1/resource
	Method  string `json:"method"`  // e.g., GET, POST, PUT, DELETE
	Public  bool   `json:"public"`  // Indicates if the endpoint is public or requires authentication
	Webhook bool   `json:"webhook"` // Indicates if the endpoint is a webhook
}
