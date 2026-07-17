package email

import (
	"abstractions/db/models"
	"abstractions/repositories"
	"context"

	"gorm.io/gorm"
)

// FailureThreshold is the number of transient (soft) failures an address may accumulate
// before it is added to the suppression list.
const FailureThreshold = 3

// EmailLogEntry is the data recorded for a single send attempt.
type EmailLogEntry struct {
	Email        string
	Subject      string
	MessageID    string
	Status       string
	ErrorType    string
	ErrorMessage string
}

// SuppressionStore persists send logs and manages the suppression list. It is implemented
// over the root database connection (not a request transaction) because emails are sent
// from fire-and-forget goroutines that outlive the request.
type SuppressionStore interface {
	IsSuppressed(ctx context.Context, email string) (bool, error)
	Suppress(ctx context.Context, email string, reason string) error
	LogSend(ctx context.Context, entry EmailLogEntry) error
	CountSoftFailures(ctx context.Context, email string) (int, error)
}

// store is the process-wide suppression store, registered at startup. When nil the send
// path fails open: emails still send, they are simply not logged or suppressed.
var store SuppressionStore

// SetSuppressionStore registers the store used by the send path.
func SetSuppressionStore(s SuppressionStore) {
	store = s
}

type gormSuppressionStore struct {
	logRepo         *repositories.EmailLogRepository
	suppressionRepo *repositories.EmailSuppressionRepository
}

// NewGormSuppressionStore builds a SuppressionStore backed by the given (root) DB.
func NewGormSuppressionStore(db *gorm.DB) SuppressionStore {
	return &gormSuppressionStore{
		logRepo:         repositories.NewEmailLogRepository(db),
		suppressionRepo: repositories.NewEmailSuppressionRepository(db),
	}
}

func (s *gormSuppressionStore) IsSuppressed(ctx context.Context, email string) (bool, error) {
	suppression, err := s.suppressionRepo.GetActiveByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	return suppression != nil, nil
}

func (s *gormSuppressionStore) Suppress(ctx context.Context, email string, reason string) error {
	return s.suppressionRepo.Upsert(ctx, email, reason)
}

func (s *gormSuppressionStore) LogSend(ctx context.Context, entry EmailLogEntry) error {
	log := &models.EmailLog{
		Email:  entry.Email,
		Status: entry.Status,
	}
	if entry.MessageID != "" {
		log.MessageID = &entry.MessageID
	}
	if entry.Subject != "" {
		log.Subject = &entry.Subject
	}
	if entry.ErrorType != "" {
		log.ErrorType = &entry.ErrorType
	}
	if entry.ErrorMessage != "" {
		log.ErrorMessage = &entry.ErrorMessage
	}
	return s.logRepo.Create(ctx, log)
}

func (s *gormSuppressionStore) CountSoftFailures(ctx context.Context, email string) (int, error) {
	return s.logRepo.CountSoftFailuresByEmail(ctx, email)
}
