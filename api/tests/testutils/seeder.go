package testutils

import (
	"abstractions/db/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Seeder provides utilities for seeding test data
type Seeder struct {
	db *gorm.DB
}

// NewSeeder creates a new seeder instance
func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{db: db}
}

// CreateUser creates a test user
func (s *Seeder) CreateUser(opts ...func(*models.User)) *models.User {
	user := &models.User{
		ID:    uuid.New(),
		Name:  "Test User",
		Email: "test-" + uuid.New().String() + "@example.com",
	}

	// Apply any options
	for _, opt := range opts {
		opt(user)
	}

	if err := s.db.Create(user).Error; err != nil {
		panic(err)
	}

	return user
}

// CreateWorkspace creates a test workspace
func (s *Seeder) CreateWorkspace(opts ...func(*models.Workspace)) *models.Workspace {
	workspace := &models.Workspace{
		ID:   uuid.New(),
		Name: "Test Workspace",
	}

	// Apply any options
	for _, opt := range opts {
		opt(workspace)
	}

	if err := s.db.Create(workspace).Error; err != nil {
		panic(err)
	}

	return workspace
}

// CreateUserWorkspaceRelationship creates a test user workspace relationship
func (s *Seeder) CreateUserWorkspaceRelationship(user *models.User, workspace *models.Workspace, role string) *models.UserWorkspaceRelationship {
	rel := &models.UserWorkspaceRelationship{
		ID:          uuid.New(),
		UserID:      user.ID,
		WorkspaceID: workspace.ID,
		Role:        role,
	}

	if err := s.db.Create(rel).Error; err != nil {
		panic(err)
	}

	return rel
}

// CreateAuthenticator creates a test authenticator for a user
func (s *Seeder) CreateAuthenticator(user *models.User, authType models.AuthenticatorType, token string) *models.Authenticator {
	authenticator := &models.Authenticator{
		ID:     uuid.New(),
		UserID: user.ID,
		Type:   authType,
		Token:  token,
	}

	if err := s.db.Create(authenticator).Error; err != nil {
		panic(err)
	}

	return authenticator
}

// WithEmail is an option for CreateUser
func WithEmail(email string) func(*models.User) {
	return func(u *models.User) {
		u.Email = email
	}
}

// WithName returns a function that sets the name of a workspace
func WithName(name string) func(*models.Workspace) {
	return func(w *models.Workspace) {
		w.Name = name
	}
}

// WithUserName returns a function that sets the name of a user
func WithUserName(name string) func(*models.User) {
	return func(u *models.User) {
		u.Name = name
	}
}
