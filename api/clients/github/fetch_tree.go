package github_client

import (
	"context"
	"fmt"
	"strings"
)

// FileTreeItem represents a single file entry in a repository's file tree.
type FileTreeItem struct {
	Path string
	SHA  string
	Size int
}

// GetFilteredFileTree fetches the full recursive file tree for a repository via the GitHub Trees API
// and returns only the files that pass the ignore filter (see ignore_patterns.go).
// The second return value is the commit SHA at which the tree was resolved; callers should use
// this as the ref when fetching individual file contents via the GitHub Contents API.
// The returned slice is safe to process; no blob content is fetched here.
func GetFilteredFileTree(ctx context.Context, accessToken, owner, repo string) ([]FileTreeItem, string, error) {
	client := NewClient(ctx, accessToken)

	// Fetch the default branch to resolve HEAD SHA.
	branch, _, err := client.Repositories.GetBranch(ctx, owner, repo, "HEAD", 3)
	if err != nil {
		// Fallback: get repo default branch then fetch that.
		repoInfo, _, repoErr := client.Repositories.Get(ctx, owner, repo)
		if repoErr != nil {
			return nil, "", fmt.Errorf("github_client: get repo: %w", err)
		}
		defaultBranch := repoInfo.GetDefaultBranch()
		branch, _, err = client.Repositories.GetBranch(ctx, owner, repo, defaultBranch, 3)
		if err != nil {
			return nil, "", fmt.Errorf("github_client: get branch %s: %w", defaultBranch, err)
		}
	}

	commitSHA := branch.GetCommit().GetSHA()
	treeSHA := branch.GetCommit().GetCommit().GetTree().GetSHA()
	if treeSHA == "" {
		return nil, "", fmt.Errorf("github_client: could not resolve tree SHA for %s/%s", owner, repo)
	}

	tree, _, err := client.Git.GetTree(ctx, owner, repo, treeSHA, true)
	if err != nil {
		return nil, "", fmt.Errorf("github_client: get tree: %w", err)
	}

	var items []FileTreeItem
	for _, entry := range tree.Entries {
		if entry.GetType() != "blob" {
			continue
		}
		path := entry.GetPath()
		if ShouldIgnore(path, entry.GetSize()) {
			continue
		}
		items = append(items, FileTreeItem{
			Path: path,
			SHA:  entry.GetSHA(),
			Size: entry.GetSize(),
		})
	}
	return items, commitSHA, nil
}

// FilterFileTree applies the ignore rules to a pre-built list of tree items.
// Exposed as a pure function for unit testing without network calls.
func FilterFileTree(items []FileTreeItem) []FileTreeItem {
	var result []FileTreeItem
	for _, item := range items {
		if !ShouldIgnore(item.Path, item.Size) {
			result = append(result, item)
		}
	}
	return result
}

// splitFullName splits "owner/repo" into (owner, repo).
func splitFullName(fullName string) (string, string) {
	parts := strings.SplitN(fullName, "/", 2)
	if len(parts) != 2 {
		return fullName, ""
	}
	return parts[0], parts[1]
}
