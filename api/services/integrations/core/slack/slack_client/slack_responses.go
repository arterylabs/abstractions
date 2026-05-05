package slack_client

type SlackTeamInfoResponse struct {
	OK   bool `json:"ok"`
	Team struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Domain      string `json:"domain"`
		EmailDomain string `json:"email_domain"`
		Icon        struct {
			Image34  *string `json:"image_34"`
			Image44  *string `json:"image_44"`
			Image68  *string `json:"image_68"`
			Image88  *string `json:"image_88"`
			Image102 *string `json:"image_102"`
			Image132 *string `json:"image_132"`
		} `json:"icon"`
	} `json:"team"`
	User struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
}

type SlackAuthTestResponse struct {
	OK     bool   `json:"ok"`
	URL    string `json:"url"`
	Team   string `json:"team"`
	User   string `json:"user"`
	TeamID string `json:"team_id"`
	UserID string `json:"user_id"`
}
