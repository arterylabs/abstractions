package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ConnectRepositoryRequest struct {
	GitHubRepoID  int64  `json:"githubRepoId"`
	FullName      string `json:"fullName"`
	DefaultBranch string `json:"defaultBranch"`
}

func (r ConnectRepositoryRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.GitHubRepoID, validation.Required),
		validation.Field(&r.FullName, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.DefaultBranch, validation.Length(0, 255)),
	)
}
