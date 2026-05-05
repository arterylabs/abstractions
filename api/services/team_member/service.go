package team_member_service

import (
	"abstractions/db/models"
	"abstractions/repositories"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const FreeTierMemberLimit = 10

var (
	ErrMemberNotFound   = errors.New("team member not found")
	ErrDuplicateEmail   = errors.New("a member with this email already exists")
	ErrFreeTierLimitHit = errors.New("free tier allows a maximum of 10 team members; upgrade to add more")
)

type Service struct {
	repo *repositories.TeamMemberRepository
	DB   *gorm.DB
}

func NewTeamMemberService(db *gorm.DB) *Service {
	return &Service{
		repo: repositories.NewTeamMemberRepository(db),
		DB:   db,
	}
}

// AddMember adds a new subscriber to a repository's newsletter.
// Returns ErrDuplicateEmail if the email is already registered for this repo.
// Returns ErrFreeTierLimitHit when the free-tier cap (10 members) is exceeded
// and the workspace has no active subscription.
func (s *Service) AddMember(ctx context.Context, repoID uuid.UUID, email string, hasActiveSub bool) (*models.TeamMember, error) {
	existing, err := s.repo.GetByRepoAndEmail(ctx, repoID, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrDuplicateEmail
	}

	if !hasActiveSub {
		members, err := s.repo.ListByRepoID(ctx, repoID)
		if err != nil {
			return nil, err
		}
		if len(members) >= FreeTierMemberLimit {
			return nil, ErrFreeTierLimitHit
		}
	}

	tm := &models.TeamMember{
		RepoID:     repoID,
		Email:      email,
		Subscribed: true,
	}
	if err := s.repo.Create(ctx, tm); err != nil {
		return nil, err
	}
	return tm, nil
}

// ListMembers returns all team members for the given repository.
func (s *Service) ListMembers(ctx context.Context, repoID uuid.UUID) ([]*models.TeamMember, error) {
	return s.repo.ListByRepoID(ctx, repoID)
}

// RemoveMember deletes a team member by ID.
// Returns ErrMemberNotFound if no member with that ID exists.
func (s *Service) RemoveMember(ctx context.Context, memberID uuid.UUID) error {
	existing, err := s.repo.GetByID(ctx, memberID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrMemberNotFound
	}
	return s.repo.Delete(ctx, memberID)
}
