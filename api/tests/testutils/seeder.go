package testutils

import (
	"abstractions/db/models"
	"time"

	"github.com/google/uuid"
	pgvector "github.com/pgvector/pgvector-go"
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

// CreateRepository creates a test repository in the given workspace.
func (s *Seeder) CreateRepository(workspaceID, userID uuid.UUID, opts ...func(*models.Repository)) *models.Repository {
	repo := &models.Repository{
		ID:            uuid.New(),
		WorkspaceID:   workspaceID,
		UserID:        userID,
		GitHubRepoID:  int64(uuid.New().ID()),
		FullName:      "test-org/test-repo-" + uuid.New().String(),
		DefaultBranch: "main",
		Status:        models.RepositoryStatusPending,
	}
	for _, opt := range opts {
		opt(repo)
	}
	if err := s.db.Create(repo).Error; err != nil {
		panic(err)
	}
	return repo
}

// CreateFileIndex creates a test file index entry for the given repository.
func (s *Seeder) CreateFileIndex(repoID uuid.UUID, filePath string, opts ...func(*models.FileIndex)) *models.FileIndex {
	now := time.Now()
	fi := &models.FileIndex{
		ID:              uuid.New(),
		RepoID:          repoID,
		FilePath:        filePath,
		GitSHA:          "abc123",
		DomainTags:      models.StringArray{},
		LastProcessedAt: &now,
		PineconeIDs:     models.StringArray{},
	}
	for _, opt := range opts {
		opt(fi)
	}
	if err := s.db.Create(fi).Error; err != nil {
		panic(err)
	}
	return fi
}

// CreateNewsletter creates a test newsletter for the given repository.
func (s *Seeder) CreateNewsletter(repoID uuid.UUID, opts ...func(*models.Newsletter)) *models.Newsletter {
	n := &models.Newsletter{
		ID:          uuid.New(),
		RepoID:      repoID,
		Topics:      models.JSONBField{},
		HTMLContent: "<p>Test newsletter content</p>",
	}
	for _, opt := range opts {
		opt(n)
	}
	if err := s.db.Create(n).Error; err != nil {
		panic(err)
	}
	return n
}

// CreateTeamMember creates a test team member for the given repository.
func (s *Seeder) CreateTeamMember(repoID uuid.UUID, email string, opts ...func(*models.TeamMember)) *models.TeamMember {
	tm := &models.TeamMember{
		ID:         uuid.New(),
		RepoID:     repoID,
		Email:      email,
		Subscribed: true,
		AddedAt:    time.Now(),
	}
	for _, opt := range opts {
		opt(tm)
	}
	if err := s.db.Create(tm).Error; err != nil {
		panic(err)
	}
	return tm
}

// CreateChunkEmbedding creates a test chunk embedding for the given repository.
func (s *Seeder) CreateChunkEmbedding(repoID uuid.UUID, filePath string, chunkIndex int, embedding []float32, opts ...func(*models.ChunkEmbedding)) *models.ChunkEmbedding {
	vec := make([]float32, 1536)
	copy(vec, embedding)

	ce := &models.ChunkEmbedding{
		ID:         uuid.New(),
		RepoID:     repoID,
		FilePath:   filePath,
		ChunkIndex: chunkIndex,
		Summary:    "Test chunk summary",
		DomainTags: models.StringArray{},
		Embedding:  pgvector.NewVector(vec),
		CreatedAt:  time.Now(),
	}
	for _, opt := range opts {
		opt(ce)
	}
	if err := s.db.Create(ce).Error; err != nil {
		panic(err)
	}
	return ce
}
