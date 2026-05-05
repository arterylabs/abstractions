package github_client

import (
	"context"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

// NewClient builds an authenticated go-github client using the user's OAuth access token.
func NewClient(ctx context.Context, accessToken string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
