package repositories

import (
	"abstractions/db/models"
	"abstractions/utils/pagination_utils"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrPendingInvitationExists = errors.New("a pending invitation already exists for this email in this workspace")
	ErrInvalidInvitation       = errors.New("invalid or non-existent invitation")
)

type InvitationRepository struct {
	db *gorm.DB
}

func NewInvitationRepository(db *gorm.DB) *InvitationRepository {
	return &InvitationRepository{db: db}
}

func (r *InvitationRepository) Create(ctx context.Context, invitation *models.Invitation) error {
	// Check for existing pending invitation
	var existingInvitation models.Invitation
	err := r.db.WithContext(ctx).
		Where("workspace_id = ? AND email = ? AND accepted_at IS NULL AND ignored_at IS NULL AND expires_at > ?",
			invitation.WorkspaceID, invitation.Email, time.Now()).
		First(&existingInvitation).Error

	if err == nil {
		return ErrPendingInvitationExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("invitation_repository: create: check_existing: %w", err)
	}

	return r.db.WithContext(ctx).Create(invitation).Error
}

func (r *InvitationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Invitation, error) {
	var invitation models.Invitation
	if err := r.db.WithContext(ctx).First(&invitation, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("invitation_repository: get_by_id: %w", err)
	}
	return &invitation, nil
}

func (r *InvitationRepository) Update(ctx context.Context, invitation *models.Invitation) error {
	return r.db.WithContext(ctx).Save(invitation).Error
}

func (r *InvitationRepository) GetPendingByEmail(ctx context.Context, email string) ([]models.Invitation, error) {
	var invitations []models.Invitation
	err := r.db.WithContext(ctx).
		Where("email = ? AND accepted_at IS NULL AND ignored_at IS NULL AND expires_at > ?", email, time.Now()).
		Find(&invitations).Error
	return invitations, err
}

func (r *InvitationRepository) GetByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]models.Invitation, error) {
	var invitations []models.Invitation
	err := r.db.WithContext(ctx).
		Where("workspace_id = ?", workspaceID).
		Find(&invitations).Error
	return invitations, err
}

func (r *InvitationRepository) ResendInvitation(ctx context.Context, id uuid.UUID, invitedByUserID uuid.UUID, newExpiresAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&models.Invitation{}).
		Where("id = ? AND accepted_at IS NULL AND ignored_at IS NULL", id).
		Updates(map[string]interface{}{
			"invited_by_user_id": invitedByUserID,
			"expires_at":         newExpiresAt,
		}).Error
}

func (r *InvitationRepository) Accept(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Invitation{}).
		Where("id = ? AND accepted_at IS NULL AND ignored_at IS NULL", id).
		Update("accepted_at", time.Now()).
		Error
}

func (r *InvitationRepository) Ignore(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Invitation{}).
		Where("id = ? AND accepted_at IS NULL AND ignored_at IS NULL", id).
		Update("ignored_at", time.Now()).
		Error
}

func (r *InvitationRepository) BeginTx(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Begin()
}

func (r *InvitationRepository) WithTx(tx *gorm.DB) *InvitationRepository {
	return &InvitationRepository{
		db: tx,
	}
}

func (r *InvitationRepository) GetPendingByEmailPaginated(ctx context.Context, email string, params pagination_utils.PaginationParams) ([]models.Invitation, int, error) {
	var invitations []models.Invitation
	var totalCount int64

	query := r.db.WithContext(ctx).
		Table("invitations").Model(&models.Invitation{}).
		Where("email = ? AND accepted_at IS NULL AND ignored_at IS NULL AND expires_at > NOW()", email)

	// Apply search if provided
	if params.Query != "" {
		query = query.Where("email LIKE ?", "%"+params.Query+"%")
	}

	// Get total count before pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("invitation_repository: get_pending_by_email_paginated: count: %w", err)
	}

	// Apply ordering
	if params.Sort != "" {
		sortOrder := pagination_utils.ParseSortOrder(params.SortOrder)
		switch params.Sort {
		case "created_at":
			query = query.Order("created_at " + sortOrder.String())
		default:
			query = query.Order("created_at DESC")
		}
	} else {
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	if err := query.Offset(params.GetOffset()).Limit(params.GetLimit()).Find(&invitations).Error; err != nil {
		return nil, 0, fmt.Errorf("invitation_repository: get_pending_by_email_paginated: find: %w", err)
	}

	return invitations, int(totalCount), nil
}

func (r *InvitationRepository) GetByWorkspacePaginated(ctx context.Context, workspaceID uuid.UUID, params pagination_utils.PaginationParams) ([]models.Invitation, int, error) {
	var invitations []models.Invitation
	var totalCount int64

	query := r.db.WithContext(ctx).Table("invitations").Model(&models.Invitation{}).
		Where("workspace_id = ?", workspaceID)

	// Apply search if provided
	if params.Query != "" {
		query = query.Where("email LIKE ?", "%"+params.Query+"%")
	}

	// Get total count before pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("invitation_repository: get_by_workspace_paginated: count: %w", err)
	}

	// Apply ordering
	if params.Sort != "" {
		sortOrder := pagination_utils.ParseSortOrder(params.SortOrder)
		switch params.Sort {
		case "created_at":
			query = query.Order("created_at " + sortOrder.String())
		case "email":
			query = query.Order("email " + sortOrder.String())
		case "role":
			query = query.Order("role " + sortOrder.String())
		case "state":
			query = query.Order("accepted_at IS NOT NULL DESC")
		default:
			query = query.Order("created_at DESC")
		}
	} else {
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	if err := query.Table("invitations").Offset(params.GetOffset()).Limit(params.GetLimit()).Find(&invitations).Error; err != nil {
		return nil, 0, fmt.Errorf("invitation_repository: get_by_workspace_paginated: find: %w", err)
	}

	return invitations, int(totalCount), nil
}

func (r *InvitationRepository) CountByWorkspace(ctx context.Context, workspaceID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Invitation{}).
		Where("workspace_id = ?", workspaceID).
		Count(&count).Error
	return count, err
}

func (r *InvitationRepository) CountPendingByEmail(ctx context.Context, email string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Invitation{}).
		Where("email = ? AND accepted_at IS NULL AND ignored_at IS NULL AND expires_at > ?", email, time.Now()).
		Count(&count).Error
	return count, err
}
