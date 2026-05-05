package repositories

import (
	"abstractions/db/models"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamMemberRepository struct {
	db *gorm.DB
}

func NewTeamMemberRepository(db *gorm.DB) *TeamMemberRepository {
	return &TeamMemberRepository{db: db}
}

func (r *TeamMemberRepository) Create(ctx context.Context, tm *models.TeamMember) error {
	return r.db.WithContext(ctx).Create(tm).Error
}

func (r *TeamMemberRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.TeamMember, error) {
	var tm models.TeamMember
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&tm).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tm, nil
}

func (r *TeamMemberRepository) GetByRepoAndEmail(ctx context.Context, repoID uuid.UUID, email string) (*models.TeamMember, error) {
	var tm models.TeamMember
	if err := r.db.WithContext(ctx).Where("repo_id = ? AND email = ?", repoID, email).First(&tm).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tm, nil
}

func (r *TeamMemberRepository) ListByRepoID(ctx context.Context, repoID uuid.UUID) ([]*models.TeamMember, error) {
	var members []*models.TeamMember
	if err := r.db.WithContext(ctx).Where("repo_id = ?", repoID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *TeamMemberRepository) ListSubscribedByRepoID(ctx context.Context, repoID uuid.UUID) ([]*models.TeamMember, error) {
	var members []*models.TeamMember
	if err := r.db.WithContext(ctx).Where("repo_id = ? AND subscribed = TRUE", repoID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *TeamMemberRepository) UpdateSubscription(ctx context.Context, id uuid.UUID, subscribed bool) error {
	return r.db.WithContext(ctx).Model(&models.TeamMember{}).Where("id = ?", id).Update("subscribed", subscribed).Error
}

func (r *TeamMemberRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.TeamMember{}).Error
}
