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

// RepositoryIndexingFailedParams holds the data for a repository indexing failure notification.
type RepositoryIndexingFailedParams struct {
	RepoName      string
	WorkspaceName string
	ErrorMessage  string
}

// SendRepositoryIndexingFailedNotification sends a Slack alert when repository indexing fails.
func SendRepositoryIndexingFailedNotification(ctx context.Context, params RepositoryIndexingFailedParams) {
	errMsg := params.ErrorMessage
	if len(errMsg) > 500 {
		errMsg = errMsg[:500] + "..."
	}
	sendSlack(ctx, slackPayload{
		Blocks: []any{
			sectionBlock{
				Type: "section",
				Text: &textBlock{Type: "mrkdwn", Text: ":warning: *Repository indexing failed*"},
				Fields: []textBlock{
					mrkdwn(fmt.Sprintf("*Repository*\n%s", params.RepoName)),
					mrkdwn(fmt.Sprintf("*Workspace*\n%s", params.WorkspaceName)),
					mrkdwn(fmt.Sprintf("*Error*\n%s", errMsg)),
				},
			},
			contextBlock{
				Type:     "context",
				Elements: []textBlock{mrkdwn(ts())},
			},
		},
	})
}

// ErrorParams holds the data for a generic error notification.
type ErrorParams struct {
	// Source is a short description of where the error occurred (e.g. "GET /foo" or "job:index_repository").
	Source  string
	Message string
}

// SendErrorNotification sends a Slack alert for a non-panic error.
func SendJobFailureNotification(ctx context.Context, params ErrorParams) {
	msg := params.Message
	if len(msg) > 500 {
		msg = msg[:500] + "..."
	}
	sendSlack(ctx, slackPayload{
		Blocks: []any{
			sectionBlock{
				Type: "section",
				Text: &textBlock{Type: "mrkdwn", Text: ":x: *Job failed*"},
				Fields: []textBlock{
					mrkdwn(fmt.Sprintf("*Source*\n%s", params.Source)),
					mrkdwn(fmt.Sprintf("*Error*\n%s", msg)),
				},
			},
			contextBlock{
				Type:     "context",
				Elements: []textBlock{mrkdwn(ts())},
			},
		},
	})
}

// PanicParams holds the data for a panic/crash notification.
type PanicParams struct {
	// Source is a short description of where the panic occurred (e.g. "GET /foo" or "job:index_repository").
	Source  string
	Message string
	Stack   string
}

// SendPanicNotification sends a Slack alert when a panic is recovered.
func SendPanicNotification(ctx context.Context, params PanicParams) {
	stack := params.Stack
	if len(stack) > 800 {
		stack = stack[:800] + "..."
	}
	sendSlack(ctx, slackPayload{
		Blocks: []any{
			headerBlock{
				Type: "header",
				Text: plainText(":rotating_light: Panic recovered"),
			},
			sectionBlock{
				Type: "section",
				Fields: []textBlock{
					mrkdwn(fmt.Sprintf("*Source*\n%s", params.Source)),
					mrkdwn(fmt.Sprintf("*Message*\n%s", params.Message)),
				},
			},
			sectionBlock{
				Type: "section",
				Text: &textBlock{Type: "mrkdwn", Text: fmt.Sprintf("*Stack trace*\n```%s```", stack)},
			},
			contextBlock{
				Type:     "context",
				Elements: []textBlock{mrkdwn(ts())},
			},
		},
	})
}

// NewsletterFailedParams holds the data for a newsletter failure notification.
type NewsletterFailedParams struct {
	RepoName      string
	WorkspaceName string
	ErrorMessage  string
}

// SendNewsletterFailedNotification sends a Slack alert when newsletter generation or sending fails.
func SendNewsletterFailedNotification(ctx context.Context, params NewsletterFailedParams) {
	errMsg := params.ErrorMessage
	if len(errMsg) > 500 {
		errMsg = errMsg[:500] + "..."
	}
	sendSlack(ctx, slackPayload{
		Blocks: []any{
			sectionBlock{
				Type: "section",
				Text: &textBlock{Type: "mrkdwn", Text: ":warning: *Newsletter failed*"},
				Fields: []textBlock{
					mrkdwn(fmt.Sprintf("*Repository*\n%s", params.RepoName)),
					mrkdwn(fmt.Sprintf("*Workspace*\n%s", params.WorkspaceName)),
					mrkdwn(fmt.Sprintf("*Error*\n%s", errMsg)),
				},
			},
			contextBlock{
				Type:     "context",
				Elements: []textBlock{mrkdwn(ts())},
			},
		},
	})
}
