package newsletter

import (
	"abstractions/db/models"
	"abstractions/repositories"
	repository_service "abstractions/services/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrNewsletterNotFound      = errors.New("newsletter not found")
	ErrNewsletterNotResendable = errors.New("only newsletters with status 'error' can be resent")
	ErrTopicNotFound           = errors.New("topic not found")
	ErrTopicAlreadySent        = errors.New("topic already sent")
)

type Service struct {
	newsletterRepo      *repositories.NewsletterRepository
	newsletterTopicRepo *repositories.NewsletterTopicRepository
	repositoryService   *repository_service.Service
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		newsletterRepo:      repositories.NewNewsletterRepository(db),
		newsletterTopicRepo: repositories.NewNewsletterTopicRepository(db),
		repositoryService:   repository_service.NewService(db),
	}
}

func (s *Service) Count(ctx context.Context) (int, error) {
	return s.newsletterRepo.Count(ctx)
}

func (s *Service) QueueNewsletter(ctx context.Context, workspaceID, repoID uuid.UUID) (*models.Repository, *models.Newsletter, error) {
	repo, err := s.repositoryService.GetRepositoryInWorkspace(ctx, workspaceID, repoID)
	if err != nil {
		return nil, nil, err
	}

	newsletter := &models.Newsletter{
		ID:     uuid.New(),
		RepoID: repo.ID,
		Status: models.NewsletterStatusQueued,
	}

	if nextTopic, err := s.newsletterTopicRepo.GetNextForRepo(ctx, repo.ID); err != nil {
		return nil, nil, fmt.Errorf("newsletter_service: queue_newsletter: get_next_topic: %w", err)
	} else if nextTopic != nil {
		newsletter.TopicID = &nextTopic.ID
	}

	if err := s.newsletterRepo.Create(ctx, newsletter); err != nil {
		return nil, nil, fmt.Errorf("newsletter_service: queue_newsletter: create: %w", err)
	}
	return repo, newsletter, nil
}

func (s *Service) ListNewsletters(ctx context.Context, workspaceID, repoID uuid.UUID) ([]*models.Newsletter, error) {
	repo, err := s.repositoryService.GetRepositoryInWorkspace(ctx, workspaceID, repoID)
	if err != nil {
		return nil, err
	}
	return s.newsletterRepo.ListByRepoID(ctx, repo.ID)
}

func (s *Service) GetNewsletter(ctx context.Context, workspaceID, repoID, newsletterID uuid.UUID) (*models.Newsletter, error) {
	repo, err := s.repositoryService.GetRepositoryInWorkspace(ctx, workspaceID, repoID)
	if err != nil {
		return nil, err
	}
	newsletter, err := s.newsletterRepo.GetByID(ctx, newsletterID)
	if err != nil {
		return nil, fmt.Errorf("newsletter_service: get_newsletter: %w", err)
	}
	if newsletter == nil || newsletter.RepoID != repo.ID {
		return nil, ErrNewsletterNotFound
	}
	return newsletter, nil
}

func (s *Service) ResetNewsletterForResend(ctx context.Context, workspaceID, repoID, newsletterID uuid.UUID) (*models.Repository, *models.Newsletter, error) {
	repo, err := s.repositoryService.GetRepositoryInWorkspace(ctx, workspaceID, repoID)
	if err != nil {
		return nil, nil, err
	}
	newsletter, err := s.newsletterRepo.GetByID(ctx, newsletterID)
	if err != nil {
		return nil, nil, fmt.Errorf("newsletter_service: reset_newsletter_for_resend: get: %w", err)
	}
	if newsletter == nil || newsletter.RepoID != repo.ID {
		return nil, nil, ErrNewsletterNotFound
	}
	if newsletter.Status != models.NewsletterStatusError {
		return nil, nil, ErrNewsletterNotResendable
	}

	newsletter.Status = models.NewsletterStatusQueued
	newsletter.ErrorMessage = nil
	newsletter.HTMLContent = ""
	newsletter.TopicID = nil
	newsletter.SentAt = nil
	newsletter.RecipientCount = 0
	if err := s.newsletterRepo.Update(ctx, newsletter); err != nil {
		return nil, nil, fmt.Errorf("newsletter_service: reset_newsletter_for_resend: update: %w", err)
	}
	return repo, newsletter, nil
}

func (s *Service) ListNewsletterTopics(ctx context.Context, workspaceID, repoID uuid.UUID) ([]*repositories.NewsletterTopicWithNewsletterID, error) {
	repo, err := s.repositoryService.GetRepositoryInWorkspace(ctx, workspaceID, repoID)
	if err != nil {
		return nil, err
	}
	return s.newsletterTopicRepo.ListByRepoID(ctx, repo.ID)
}

func (s *Service) CreateNewsletterTopic(ctx context.Context, workspaceID, repoID uuid.UUID, title, description string, relatedFiles []string, createdBy *uuid.UUID) (*models.NewsletterTopic, error) {
	repo, err := s.repositoryService.GetRepositoryInWorkspace(ctx, workspaceID, repoID)
	if err != nil {
		return nil, err
	}
	relatedJSON, err := json.Marshal(relatedFiles)
	if err != nil {
		return nil, fmt.Errorf("newsletter_service: create_topic: encode_related_files: %w", err)
	}
	maxPos, err := s.newsletterTopicRepo.GetMaxPendingPosition(ctx, repo.ID)
	if err != nil {
		return nil, fmt.Errorf("newsletter_service: create_topic: max_position: %w", err)
	}
	topic := &models.NewsletterTopic{
		ID:           uuid.New(),
		RepoID:       repo.ID,
		Title:        title,
		Description:  description,
		RelatedFiles: models.JSONBRaw(relatedJSON),
		Position:     maxPos + 1,
		CreatedBy:    createdBy,
	}
	if err := s.newsletterTopicRepo.Create(ctx, topic); err != nil {
		return nil, fmt.Errorf("newsletter_service: create_topic: create: %w", err)
	}
	return topic, nil
}

func (s *Service) UpdateNewsletterTopic(ctx context.Context, workspaceID, repoID, topicID uuid.UUID, title, description string, relatedFiles []string) (*models.NewsletterTopic, error) {
	topic, err := s.getTopicInRepository(ctx, workspaceID, repoID, topicID)
	if err != nil {
		return nil, err
	}
	if topic.SentAt != nil {
		return nil, ErrTopicAlreadySent
	}
	relatedJSON, err := json.Marshal(relatedFiles)
	if err != nil {
		return nil, fmt.Errorf("newsletter_service: update_topic: encode_related_files: %w", err)
	}
	if err := s.newsletterTopicRepo.UpdateTopic(ctx, topicID, title, description, models.JSONBRaw(relatedJSON)); err != nil {
		return nil, fmt.Errorf("newsletter_service: update_topic: update: %w", err)
	}
	topic.Title = title
	topic.Description = description
	topic.RelatedFiles = models.JSONBRaw(relatedJSON)
	return topic, nil
}

func (s *Service) ReorderNewsletterTopic(ctx context.Context, workspaceID, repoID, topicID uuid.UUID, direction string) error {
	topic, err := s.getTopicInRepository(ctx, workspaceID, repoID, topicID)
	if err != nil {
		return err
	}
	if topic.SentAt != nil {
		return ErrTopicAlreadySent
	}
	if err := s.newsletterTopicRepo.ReorderTopic(ctx, repoID, topicID, direction); err != nil {
		return fmt.Errorf("newsletter_service: reorder_topic: %w", err)
	}
	return nil
}

func (s *Service) DeleteNewsletterTopic(ctx context.Context, workspaceID, repoID, topicID uuid.UUID) error {
	topic, err := s.getTopicInRepository(ctx, workspaceID, repoID, topicID)
	if err != nil {
		return err
	}
	if topic.SentAt != nil {
		return ErrTopicAlreadySent
	}
	if err := s.newsletterTopicRepo.Delete(ctx, topicID); err != nil {
		return fmt.Errorf("newsletter_service: delete_topic: %w", err)
	}
	return nil
}

func (s *Service) getTopicInRepository(ctx context.Context, workspaceID, repoID, topicID uuid.UUID) (*models.NewsletterTopic, error) {
	repo, err := s.repositoryService.GetRepositoryInWorkspace(ctx, workspaceID, repoID)
	if err != nil {
		return nil, err
	}
	topic, err := s.newsletterTopicRepo.GetByID(ctx, topicID)
	if err != nil {
		return nil, fmt.Errorf("newsletter_service: get_topic: %w", err)
	}
	if topic == nil || topic.RepoID != repo.ID {
		return nil, ErrTopicNotFound
	}
	return topic, nil
}
