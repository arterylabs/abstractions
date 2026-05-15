package slack_client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SlackClient struct {
	Token string
	CTX   context.Context
}

func NewSlackClient(ctx context.Context, token string) *SlackClient {
	return &SlackClient{
		Token: token,
		CTX:   ctx,
	}
}

// AuthTest retrieves Slack team information using the access token
func AuthTest(accessToken string) (*SlackAuthTestResponse, error) {
	url := "https://slack.com/api/auth.test"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("slack_client: auth_test: create_request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("slack_client: auth_test: do_request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("slack_client: auth_test: read_body: %w", err)
	}

	var slackResp SlackAuthTestResponse
	if err := json.Unmarshal(body, &slackResp); err != nil {
		return nil, fmt.Errorf("slack_client: auth_test: unmarshal: %w", err)
	}

	if !slackResp.OK {
		return nil, fmt.Errorf("slack_client: auth_test: slack returned not ok")
	}

	return &slackResp, nil
}

func GetTeamInfo(accessToken string) (*SlackTeamInfoResponse, error) {
	url := "https://slack.com/api/team.info"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("slack_client: get_team_info: create_request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("slack_client: get_team_info: do_request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("slack_client: get_team_info: read_body: %w", err)
	}

	var slackResp SlackTeamInfoResponse
	if err := json.Unmarshal(body, &slackResp); err != nil {
		return nil, fmt.Errorf("slack_client: get_team_info: unmarshal: %w", err)
	}

	if !slackResp.OK {
		return nil, fmt.Errorf("slack_client: get_team_info: slack returned not ok")
	}

	return &slackResp, nil
}

func (c *SlackClient) SendMessage(channel string, message string) error {
	// Implementation for sending a message to a Slack channel
	// This is a placeholder; actual implementation would use Slack's API
	return nil
}

func (c *SlackClient) GetTeamInfo(userID string) (string, error) {
	// Implementation for retrieving user information from Slack
	// This is a placeholder; actual implementation would use Slack's API
	return "Team Info", nil
}
