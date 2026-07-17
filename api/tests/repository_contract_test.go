package tests

import (
	"abstractions/db/models"
	"abstractions/repositories"
	"abstractions/services/integrations/enums"
	"abstractions/tests/testutils"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RepositoryContractTestSuite struct {
	testutils.TestSuite
	seeder *testutils.Seeder
}

func TestRepositoryContract(t *testing.T) {
	suite.Run(t, new(RepositoryContractTestSuite))
}

func (s *RepositoryContractTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
}

func (s *RepositoryContractTestSuite) TestUserRepositoryCRUDAndDeletes() {
	repo := repositories.NewUserRepository(s.DB)
	user := &models.User{
		ID:    uuid.New(),
		Name:  "Repository User",
		Email: "repo-user@example.com",
	}

	assert.NoError(s.T(), repo.Create(s.Ctx, user))

	got, err := repo.GetByID(s.Ctx, user.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user.Email, got.Email)

	got.Name = "Updated Repository User"
	assert.NoError(s.T(), repo.Update(s.Ctx, got))

	updated, err := repo.GetByID(s.Ctx, user.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "Updated Repository User", updated.Name)

	deletedBy := uuid.New()
	assert.NoError(s.T(), repo.Delete(s.Ctx, user.ID, &deletedBy))

	deleted, err := repo.GetByID(s.Ctx, user.ID)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), deleted)

	includingDeleted, err := repo.GetByIDsIncludingDeleted(s.Ctx, []uuid.UUID{user.ID})
	assert.NoError(s.T(), err)
	assert.Len(s.T(), includingDeleted, 1)
	assert.NotNil(s.T(), includingDeleted[0].DeletedAt)
	assert.Equal(s.T(), deletedBy, *includingDeleted[0].DeletedBy)

	assert.NoError(s.T(), repo.PermanentlyDelete(s.Ctx, user.ID))
	afterHardDelete, err := repo.GetByIDsIncludingDeleted(s.Ctx, []uuid.UUID{user.ID})
	assert.NoError(s.T(), err)
	assert.Empty(s.T(), afterHardDelete)
}

func (s *RepositoryContractTestSuite) TestRepositorySoftAndPermanentDelete() {
	user := s.seeder.CreateUser()
	workspace := s.seeder.CreateWorkspace()
	repoModel := s.seeder.CreateRepository(workspace.ID, user.ID)
	repo := repositories.NewRepositoryRepository(s.DB)

	assert.NoError(s.T(), repo.Delete(s.Ctx, repoModel.ID, user.ID))

	deleted, err := repo.GetByID(s.Ctx, repoModel.ID)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), deleted)

	includingDeleted, err := repo.GetByIDIncludingDeleted(s.Ctx, repoModel.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), includingDeleted)
	assert.NotNil(s.T(), includingDeleted.DeletedAt)

	assert.NoError(s.T(), repo.PermanentlyDelete(s.Ctx, repoModel.ID))
	afterHardDelete, err := repo.GetByIDIncludingDeleted(s.Ctx, repoModel.ID)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), afterHardDelete)
}

func (s *RepositoryContractTestSuite) TestIntegrationTypedUpdatesAndDeletes() {
	user := s.seeder.CreateUser()
	workspace := s.seeder.CreateWorkspace()
	repo := repositories.NewIntegrationRepository(s.DB)
	serviceTitle := "GitHub"
	integration := &models.Integration{
		ID:           uuid.New(),
		WorkspaceID:  workspace.ID,
		Type:         enums.IntegrationTypeEnum("github"),
		ServiceTitle: &serviceTitle,
		Token:        "token",
		CreatedBy:    user.ID,
		UpdatedBy:    user.ID,
	}

	assert.NoError(s.T(), repo.Create(s.Ctx, integration))

	assert.NoError(s.T(), repo.UpdateTokens(s.Ctx, integration.ID, "new-token", nil, nil, user.ID))
	updated, err := repo.GetByID(s.Ctx, integration.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "new-token", updated.Token)

	assert.NoError(s.T(), repo.Delete(s.Ctx, integration.ID, user.ID))
	deleted, err := repo.GetByID(s.Ctx, integration.ID)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), deleted)

	assert.NoError(s.T(), repo.PermanentlyDelete(s.Ctx, integration.ID))
	includingDeleted, err := repo.GetByIDIncludingDeleted(s.Ctx, integration.ID)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), includingDeleted)
}

func (s *RepositoryContractTestSuite) TestInvitationStateTransitionsUseTypedPersistence() {
	user := s.seeder.CreateUser()
	workspace := s.seeder.CreateWorkspace()
	repo := repositories.NewInvitationRepository(s.DB)
	invitation := &models.Invitation{
		ID:          uuid.New(),
		WorkspaceID: workspace.ID,
		Email:       "invitee@example.com",
		Role:        "member",
		InvitedBy:   user.ID,
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}

	assert.NoError(s.T(), repo.Create(s.Ctx, invitation))
	assert.NoError(s.T(), repo.Accept(s.Ctx, invitation.ID))

	got, err := repo.GetByID(s.Ctx, invitation.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), got.AcceptedAt)
}

func (s *RepositoryContractTestSuite) TestNewsletterPreloadsTopic() {
	user := s.seeder.CreateUser()
	workspace := s.seeder.CreateWorkspace()
	repository := s.seeder.CreateRepository(workspace.ID, user.ID)
	topic := s.seeder.CreateNewsletterTopic(repository.ID, nil)
	newsletter := s.seeder.CreateNewsletter(repository.ID, func(n *models.Newsletter) {
		n.TopicID = &topic.ID
	})
	repo := repositories.NewNewsletterRepository(s.DB)

	got, err := repo.GetByID(context.Background(), newsletter.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), got)
	assert.NotNil(s.T(), got.Topic)
	assert.Equal(s.T(), topic.ID, got.Topic.ID)
}
