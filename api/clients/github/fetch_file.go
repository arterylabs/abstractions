package github_client

import (
	"context"
	"fmt"

	"github.com/google/go-github/v57/github"
)

// GitHubFileMetadata holds the metadata fields from a GitHub Contents API response.
// The actual file content (Content field) is intentionally excluded — it must never
// be persisted to the database. This struct is safe to store as service_metadata.
type GitHubFileMetadata struct {
	Type            string `json:"type,omitempty"`
	Encoding        string `json:"encoding,omitempty"`
	Size            int    `json:"size,omitempty"`
	Name            string `json:"name,omitempty"`
	Path            string `json:"path,omitempty"`
	SHA             string `json:"sha,omitempty"`
	URL             string `json:"url,omitempty"`
	GitURL          string `json:"gitUrl,omitempty"`
	HTMLURL         string `json:"htmlUrl,omitempty"`
	DownloadURL     string `json:"downloadUrl,omitempty"`
	SubmoduleGitURL string `json:"submoduleGitUrl,omitempty"`
	Target          string `json:"target,omitempty"`
}

// GetFileContent fetches the raw decoded content of a single file from GitHub.
// It also returns a GitHubFileMetadata struct populated from the Contents API response,
// with the file content field stripped out. The content is returned in-memory and is
// never persisted by this function.
// ref can be a branch name, tag, or commit SHA; pass empty string for default branch.
func GetFileContent(ctx context.Context, accessToken, owner, repo, path, ref string) ([]byte, *GitHubFileMetadata, error) {
	client := NewClient(ctx, accessToken)

	var opts *github.RepositoryContentGetOptions
	if ref != "" {
		opts = &github.RepositoryContentGetOptions{Ref: ref}
	}

	fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, opts)
	if err != nil {
		return nil, nil, fmt.Errorf("github_client: get file content %s: %w", path, err)
	}
	if fileContent == nil {
		return nil, nil, fmt.Errorf("github_client: path %s is a directory, not a file", path)
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return nil, nil, fmt.Errorf("github_client: decode file content %s: %w", path, err)
	}

	meta := metadataFromContent(fileContent)
	return []byte(content), meta, nil
}

// metadataFromContent extracts all non-content fields from a RepositoryContent response.
func metadataFromContent(rc *github.RepositoryContent) *GitHubFileMetadata {
	if rc == nil {
		return nil
	}
	return &GitHubFileMetadata{
		Type:            rc.GetType(),
		Encoding:        rc.GetEncoding(),
		Size:            rc.GetSize(),
		Name:            rc.GetName(),
		Path:            rc.GetPath(),
		SHA:             rc.GetSHA(),
		URL:             rc.GetURL(),
		GitURL:          rc.GetGitURL(),
		HTMLURL:         rc.GetHTMLURL(),
		DownloadURL:     rc.GetDownloadURL(),
		SubmoduleGitURL: rc.GetSubmoduleGitURL(),
		Target:          rc.GetTarget(),
	}
}
