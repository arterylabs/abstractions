package tests

import (
	"abstractions/repositories"
	"abstractions/tests/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TeamMemberRepositoryTestSuite struct {
	testutils.TestSuite
	seeder *testutils.Seeder
	repo   *repositories.TeamMemberRepository
}

func TestTeamMemberRepository(t *testing.T) {
	suite.Run(t, new(TeamMemberRepositoryTestSuite))
}

func (s *TeamMemberRepositoryTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.repo = repositories.NewTeamMemberRepository(s.DB)
}

// TestCreate_GetByID creates a team member and retrieves it by ID.
func (s *TeamMemberRepositoryTestSuite) TestCreate_GetByID() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repository := s.seeder.CreateRepository(workspace.ID, user.ID)

	tm := s.seeder.CreateTeamMember(repository.ID, "dev@example.com")
	assert.NotNil(s.T(), tm)

	got, err := s.repo.GetByID(s.Ctx, tm.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), got)
	assert.Equal(s.T(), "dev@example.com", got.Email)
	assert.True(s.T(), got.Subscribed)
}

// TestGetByRepoAndEmail finds a member by repository and email address.
func (s *TeamMemberRepositoryTestSuite) TestGetByRepoAndEmail() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repository := s.seeder.CreateRepository(workspace.ID, user.ID)

	s.seeder.CreateTeamMember(repository.ID, "alice@example.com")

	got, err := s.repo.GetByRepoAndEmail(s.Ctx, repository.ID, "alice@example.com")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), got)
	assert.Equal(s.T(), "alice@example.com", got.Email)
}

// TestListByRepoID_ReturnsOnlyTargetRepo ensures per-repository isolation.
func (s *TeamMemberRepositoryTestSuite) TestListByRepoID_ReturnsOnlyTargetRepo() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo1 := s.seeder.CreateRepository(workspace.ID, user.ID)
	repo2 := s.seeder.CreateRepository(workspace.ID, user.ID)

	s.seeder.CreateTeamMember(repo1.ID, "a@example.com")
	s.seeder.CreateTeamMember(repo1.ID, "b@example.com")
	s.seeder.CreateTeamMember(repo2.ID, "c@example.com")

	members, err := s.repo.ListByRepoID(s.Ctx, repo1.ID)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), members, 2)
}

// TestUpdateSubscription toggles a member's subscription flag.
func (s *TeamMemberRepositoryTestSuite) TestUpdateSubscription() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repository := s.seeder.CreateRepository(workspace.ID, user.ID)

	tm := s.seeder.CreateTeamMember(repository.ID, "sub@example.com")
	assert.True(s.T(), tm.Subscribed)

	err := s.repo.UpdateSubscription(s.Ctx, tm.ID, false)
	assert.NoError(s.T(), err)

	got, err := s.repo.GetByID(s.Ctx, tm.ID)
	assert.NoError(s.T(), err)
	assert.False(s.T(), got.Subscribed)
}

// TestDelete removes a team member from the database.
func (s *TeamMemberRepositoryTestSuite) TestDelete() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repository := s.seeder.CreateRepository(workspace.ID, user.ID)

	tm := s.seeder.CreateTeamMember(repository.ID, "remove@example.com")

	err := s.repo.Delete(s.Ctx, tm.ID)
	assert.NoError(s.T(), err)

	got, err := s.repo.GetByID(s.Ctx, tm.ID)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), got)
}

// TestListSubscribedByRepoID returns only subscribed members.
func (s *TeamMemberRepositoryTestSuite) TestListSubscribedByRepoID() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repository := s.seeder.CreateRepository(workspace.ID, user.ID)

	s.seeder.CreateTeamMember(repository.ID, "sub1@example.com")
	s.seeder.CreateTeamMember(repository.ID, "sub2@example.com")
	tm3 := s.seeder.CreateTeamMember(repository.ID, "unsub@example.com")

	err := s.repo.UpdateSubscription(s.Ctx, tm3.ID, false)
	assert.NoError(s.T(), err)

	subscribed, err := s.repo.ListSubscribedByRepoID(s.Ctx, repository.ID)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), subscribed, 2)
}
