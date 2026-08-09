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
// and returns only the files that pass the ignore filter (built-in + .gitignore + .indexignore).
// branch is the branch name to resolve (empty uses the repository default branch).
// The second return value is the commit SHA at which the tree was resolved; callers should use
// this as the ref when fetching individual file contents via the GitHub Contents API.
func GetFilteredFileTree(ctx context.Context, accessToken, owner, repo, branch string) ([]FileTreeItem, string, error) {
	client := NewClient(ctx, accessToken)

	if branch == "" {
		branch = "HEAD"
	}

	// Resolve branch → commit SHA + tree SHA.
	branchRef, _, err := client.Repositories.GetBranch(ctx, owner, repo, branch, 3)
	if err != nil {
		// Fallback: repository default branch when HEAD or explicit branch fails.
		repoInfo, _, repoErr := client.Repositories.Get(ctx, owner, repo)
		if repoErr != nil {
			return nil, "", fmt.Errorf("github_client: get repo: %w", err)
		}
		defaultBranch := repoInfo.GetDefaultBranch()
		if branch == "HEAD" || branch == defaultBranch {
			branchRef, _, err = client.Repositories.GetBranch(ctx, owner, repo, defaultBranch, 3)
			if err != nil {
				return nil, "", fmt.Errorf("github_client: get branch %s: %w", defaultBranch, err)
			}
		} else {
			return nil, "", fmt.Errorf("github_client: get branch %s: %w", branch, err)
		}
	}

	commitSHA := branchRef.GetCommit().GetSHA()
	treeSHA := branchRef.GetCommit().GetCommit().GetTree().GetSHA()
	if treeSHA == "" {
		return nil, "", fmt.Errorf("github_client: could not resolve tree SHA for %s/%s", owner, repo)
	}

	tree, _, err := client.Git.GetTree(ctx, owner, repo, treeSHA, true)
	if err != nil {
		return nil, "", fmt.Errorf("github_client: get tree: %w", err)
	}

	// Collect ignore-file paths present in the tree so we can fetch their contents.
	var gitignorePaths []string
	for _, entry := range tree.Entries {
		if entry.GetType() != "blob" {
			continue
		}
		path := entry.GetPath()
		base := path
		if i := strings.LastIndex(path, "/"); i >= 0 {
			base = path[i+1:]
		}
		if base == ".gitignore" || base == ".indexignore" {
			gitignorePaths = append(gitignorePaths, path)
		}
	}

	var matcher *GitIgnoreMatcher
	for _, ignorePath := range gitignorePaths {
		content, _, fetchErr := GetFileContent(ctx, accessToken, owner, repo, ignorePath, commitSHA)
		if fetchErr != nil {
			continue
		}
		parsed := ParseGitIgnore(string(content))
		matcher = matcher.Merge(parsed)
	}

	var items []FileTreeItem
	for _, entry := range tree.Entries {
		if entry.GetType() != "blob" {
			continue
		}
		path := entry.GetPath()
		if ShouldIgnoreWithMatcher(path, entry.GetSize(), matcher) {
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

// FilterFileTreeWithMatcher applies built-in + custom ignore matchers.
func FilterFileTreeWithMatcher(items []FileTreeItem, matchers ...*GitIgnoreMatcher) []FileTreeItem {
	var result []FileTreeItem
	for _, item := range items {
		if !ShouldIgnoreWithMatcher(item.Path, item.Size, matchers...) {
			result = append(result, item)
		}
	}
	return result
}
