package github_client

import (
	"context"
	"fmt"

	"github.com/google/go-github/v57/github"
)

// GetFileContent fetches the raw decoded content of a single file from GitHub.
// The content is returned in-memory and is never persisted by this function.
// ref can be a branch name, tag, or commit SHA; pass empty string for default branch.
func GetFileContent(ctx context.Context, accessToken, owner, repo, path, ref string) ([]byte, error) {
	client := NewClient(ctx, accessToken)

	var opts *github.RepositoryContentGetOptions
	if ref != "" {
		opts = &github.RepositoryContentGetOptions{Ref: ref}
	}

	fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, opts)
	if err != nil {
		return nil, fmt.Errorf("github_client: get file content %s: %w", path, err)
	}
	if fileContent == nil {
		return nil, fmt.Errorf("github_client: path %s is a directory, not a file", path)
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return nil, fmt.Errorf("github_client: decode file content %s: %w", path, err)
	}
	return []byte(content), nil
}
