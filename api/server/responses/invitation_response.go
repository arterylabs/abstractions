package responses

import (
	"abstractions/db/models"
	"abstractions/server/requests"
	"time"

	"github.com/google/uuid"
)

type UserBasicInfo struct {
	ID     uuid.UUID `json:"id"`
	Email  string    `json:"email"`
	Name   string    `json:"name"`
	Avatar *string   `json:"avatar,omitempty"`
}

type InvitationResponse struct {
	ID            uuid.UUID      `json:"id"`
	WorkspaceID   uuid.UUID      `json:"workspaceId"`
	WorkspaceName string         `json:"workspaceName"`
	InvitedBy     uuid.UUID      `json:"invitedBy"`
	Email         string         `json:"email"`
	Role          string         `json:"role"`
	State         string         `json:"state"`
	ExpiresAt     time.Time      `json:"expiresAt"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	InvitedByUser *UserBasicInfo `json:"invitedByUser,omitempty"`
}

type InvitationResult struct {
	Email   string `json:"email"`
	Role    string `json:"role"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type InvitationResponseStatus string

const (
	InvitationResponseStatusSuccess        InvitationResponseStatus = "success"
	InvitationResponseStatusError          InvitationResponseStatus = "error"
	InvitationResponseStatusPartialSuccess InvitationResponseStatus = "partial_success"
)

type SendInvitationsResponse struct {
	Status  InvitationResponseStatus `json:"status"`
	Results []InvitationResult       `json:"results"`
}

type PaginatedInvitationsResponse struct {
	Results  []InvitationResponse `json:"results"`
	PageInfo PageInfo             `json:"pageInfo"`
}

func NewInvitationResponse(invitation *models.Invitation) *InvitationResponse {
	workspaceName := ""
	if invitation.Workspace != nil {
		workspaceName = invitation.Workspace.Name
	}

	response := &InvitationResponse{
		ID:            invitation.ID,
		WorkspaceID:   invitation.WorkspaceID,
		WorkspaceName: workspaceName,
		InvitedBy:     invitation.InvitedBy,
		Email:         invitation.Email,
		Role:          invitation.Role,
		State:         invitation.GetState().String(),
		ExpiresAt:     invitation.ExpiresAt,
		CreatedAt:     invitation.CreatedAt,
		UpdatedAt:     invitation.UpdatedAt,
	}

	if invitation.InvitedByUser != nil {
		response.InvitedByUser = &UserBasicInfo{
			ID:     invitation.InvitedByUser.ID,
			Name:   invitation.InvitedByUser.Name,
			Email:  invitation.InvitedByUser.Email,
			Avatar: invitation.InvitedByUser.Avatar,
		}
	}

	return response
}

func NewInvitationListResponse(invitations []models.Invitation) []InvitationResponse {
	response := make([]InvitationResponse, len(invitations))
	for i, invitation := range invitations {
		response[i] = *NewInvitationResponse(&invitation)
	}
	return response
}

// NewSendInvitationsResponse creates a new SendInvitationsResponse from a list of invitations and errors
func NewSendInvitationsResponse(invitations []requests.InvitationInput, errors []error) SendInvitationsResponse {
	results := make([]InvitationResult, len(invitations))
	successCount := 0
	failureCount := 0

	for i, invite := range invitations {
		results[i] = InvitationResult{
			Email:   invite.Email,
			Role:    invite.Role,
			Success: errors[i] == nil,
		}
		if errors[i] != nil {
			results[i].Error = errors[i].Error()
			failureCount++
		} else {
			successCount++
		}
	}

	var status InvitationResponseStatus
	if successCount == 0 {
		status = InvitationResponseStatusError
	} else if failureCount == 0 {
		status = InvitationResponseStatusSuccess
	} else {
		status = InvitationResponseStatusPartialSuccess
	}

	return SendInvitationsResponse{
		Status:  status,
		Results: results,
	}
}

func NewPaginatedInvitationsResponse(invitations []models.Invitation, page, perPage, totalRecords int) *PaginatedInvitationsResponse {
	response := make([]InvitationResponse, len(invitations))
	for i, invitation := range invitations {
		response[i] = *NewInvitationResponse(&invitation)
	}

	return &PaginatedInvitationsResponse{
		Results:  response,
		PageInfo: GetPageInfo(page, perPage, len(response), totalRecords),
	}
}
