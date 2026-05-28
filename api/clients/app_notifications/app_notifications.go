package app_notifications

import (
	"abstractions/config"
	"abstractions/logger"
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// block kit helpers

type headerBlock struct {
	Type string    `json:"type"`
	Text textBlock `json:"text"`
}

type sectionBlock struct {
	Type   string      `json:"type"`
	Text   *textBlock  `json:"text,omitempty"`
	Fields []textBlock `json:"fields,omitempty"`
}

type contextBlock struct {
	Type     string      `json:"type"`
	Elements []textBlock `json:"elements"`
}

type textBlock struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Emoji bool   `json:"emoji,omitempty"`
}

type dividerBlock struct {
	Type string `json:"type"`
}

type slackPayload struct {
	Blocks []any `json:"blocks"`
}

func mrkdwn(text string) textBlock {
	return textBlock{Type: "mrkdwn", Text: text}
}

func plainText(text string) textBlock {
	return textBlock{Type: "plain_text", Text: text, Emoji: true}
}

func ts() string {
	return time.Now().UTC().Format("Jan 2, 2006 at 15:04 UTC")
}

func sendSlack(ctx context.Context, payload slackPayload) {
	url := config.Environment.SlackAppEventsWebhookURL
	if url == "" || config.Environment.IsDevelopment() {
		return
	}

	_, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post(url)

	if err != nil {
		logger.Error(ctx, fmt.Sprintf("app_notifications: slack send failed: %v", err))
	}
}

// NewSignupParams holds the data for a new user signup notification.
type NewSignupParams struct {
	Name  string
	Email string
}

// SendNewSignupNotification sends a rich Slack card when a new user signs up.
func SendNewSignupNotification(ctx context.Context, params NewSignupParams) {
	sendSlack(ctx, slackPayload{
		Blocks: []any{
			headerBlock{
				Type: "header",
				Text: plainText(":wave: New signup"),
			},
			sectionBlock{
				Type: "section",
				Fields: []textBlock{
					mrkdwn(fmt.Sprintf("*Name*\n%s", params.Name)),
					mrkdwn(fmt.Sprintf("*Email*\n%s", params.Email)),
				},
			},
			contextBlock{
				Type:     "context",
				Elements: []textBlock{mrkdwn(ts())},
			},
		},
	})
}

// RepositoryConnectedParams holds the data for a repository connect notification.
type RepositoryConnectedParams struct {
	RepoName      string
	ServiceType   string
	WorkspaceName string
}

// SendRepositoryConnectedNotification sends a minimal Slack card when a repository is connected.
func SendRepositoryConnectedNotification(ctx context.Context, params RepositoryConnectedParams) {
	sendSlack(ctx, slackPayload{
		Blocks: []any{
			sectionBlock{
				Type: "section",
				Text: &textBlock{Type: "mrkdwn", Text: ":file_folder: *Repository connected*"},
				Fields: []textBlock{
					mrkdwn(fmt.Sprintf("*Repository*\n%s", params.RepoName)),
					mrkdwn(fmt.Sprintf("*Service*\n%s", params.ServiceType)),
					mrkdwn(fmt.Sprintf("*Workspace*\n%s", params.WorkspaceName)),
				},
			},
			contextBlock{
				Type:     "context",
				Elements: []textBlock{mrkdwn(ts())},
			},
		},
	})
}

// GitHubIntegrationConnectedParams holds the data for a GitHub integration connect notification.
type GitHubIntegrationConnectedParams struct {
	WorkspaceName string
	UserEmail     string
}

// SendGitHubIntegrationConnectedNotification sends a minimal Slack card when GitHub is connected to a workspace.
func SendGitHubIntegrationConnectedNotification(ctx context.Context, params GitHubIntegrationConnectedParams) {
	sendSlack(ctx, slackPayload{
		Blocks: []any{
			sectionBlock{
				Type: "section",
				Text: &textBlock{Type: "mrkdwn", Text: ":github: *GitHub integration connected*"},
				Fields: []textBlock{
					mrkdwn(fmt.Sprintf("*Workspace*\n%s", params.WorkspaceName)),
					mrkdwn(fmt.Sprintf("*By*\n%s", params.UserEmail)),
				},
			},
			contextBlock{
				Type:     "context",
				Elements: []textBlock{mrkdwn(ts())},
			},
		},
	})
}

// WorkspaceCreatedParams holds the data for a workspace creation notification.
type WorkspaceCreatedParams struct {
	WorkspaceName string
	OwnerEmail    string
}

// SendWorkspaceCreatedNotification sends a minimal Slack card when a workspace is created.
func SendWorkspaceCreatedNotification(ctx context.Context, params WorkspaceCreatedParams) {
	sendSlack(ctx, slackPayload{
		Blocks: []any{
			sectionBlock{
				Type: "section",
				Text: &textBlock{Type: "mrkdwn", Text: ":office: *Workspace created*"},
				Fields: []textBlock{
					mrkdwn(fmt.Sprintf("*Workspace*\n%s", params.WorkspaceName)),
					mrkdwn(fmt.Sprintf("*Owner*\n%s", params.OwnerEmail)),
				},
			},
			contextBlock{
				Type:     "context",
				Elements: []textBlock{mrkdwn(ts())},
			},
		},
	})
}

// NewsletterSentParams holds the data for a newsletter sent notification.
type NewsletterSentParams struct {
	RepoName      string
	WorkspaceName string
}

// SendNewsletterSentNotification sends a minimal Slack card when a newsletter is dispatched.
func SendNewsletterSentNotification(ctx context.Context, params NewsletterSentParams) {
	sendSlack(ctx, slackPayload{
		Blocks: []any{
			sectionBlock{
				Type: "section",
				Text: &textBlock{Type: "mrkdwn", Text: ":newspaper: *Newsletter sent*"},
				Fields: []textBlock{
					mrkdwn(fmt.Sprintf("*Repository*\n%s", params.RepoName)),
					mrkdwn(fmt.Sprintf("*Workspace*\n%s", params.WorkspaceName)),
				},
			},
			contextBlock{
				Type:     "context",
				Elements: []textBlock{mrkdwn(ts())},
			},
		},
	})
}
