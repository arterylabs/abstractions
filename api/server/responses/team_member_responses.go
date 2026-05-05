package responses

import (
	"abstractions/db/models"
	"time"
)

type TeamMemberResponse struct {
	ID         string `json:"id"`
	RepoID     string `json:"repoId"`
	Email      string `json:"email"`
	Subscribed bool   `json:"subscribed"`
	AddedAt    string `json:"addedAt"`
}

func NewTeamMemberResponse(tm *models.TeamMember) TeamMemberResponse {
	return TeamMemberResponse{
		ID:         tm.ID.String(),
		RepoID:     tm.RepoID.String(),
		Email:      tm.Email,
		Subscribed: tm.Subscribed,
		AddedAt:    tm.AddedAt.UTC().Format(time.RFC3339),
	}
}

func NewTeamMemberResponses(members []*models.TeamMember) []TeamMemberResponse {
	result := make([]TeamMemberResponse, len(members))
	for i, m := range members {
		result[i] = NewTeamMemberResponse(m)
	}
	return result
}
