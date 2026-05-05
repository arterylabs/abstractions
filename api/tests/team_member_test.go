package tests

import (
	"abstractions/repositories"
	team_member_service "abstractions/services/team_member"
	"abstractions/tests/testutils"
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ---- team member test suite ----

type TeamMemberTestSuite struct {
	testutils.TestSuite
	svc    *team_member_service.Service
	tmRepo *repositories.TeamMemberRepository
	seeder *testutils.Seeder
}

func (s *TeamMemberTestSuite) SetupTest() {
	s.svc = team_member_service.NewTeamMemberService(s.DB)
	s.tmRepo = repositories.NewTeamMemberRepository(s.DB)
	s.seeder = testutils.NewSeeder(s.DB)
}

func TestTeamMember(t *testing.T) {
	suite.Run(t, new(TeamMemberTestSuite))
}

// TestAddMember_Success adds a member on free tier (< 10 members).
func (s *TeamMemberTestSuite) TestAddMember_Success() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	member, err := s.svc.AddMember(context.Background(), repo.ID, "alice@example.com", false)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "alice@example.com", member.Email)
	assert.True(s.T(), member.Subscribed)
}

// TestAddMember_DuplicateEmail returns ErrDuplicateEmail for the same email on the same repo.
func (s *TeamMemberTestSuite) TestAddMember_DuplicateEmail() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	s.seeder.CreateTeamMember(repo.ID, "bob@example.com")

	_, err := s.svc.AddMember(context.Background(), repo.ID, "bob@example.com", false)
	assert.ErrorIs(s.T(), err, team_member_service.ErrDuplicateEmail)
}

// TestAddMember_FreeTierLimitHit returns ErrFreeTierLimitHit when there are already 10 members and no active sub.
func (s *TeamMemberTestSuite) TestAddMember_FreeTierLimitHit() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	// Seed exactly 10 members
	for i := 0; i < team_member_service.FreeTierMemberLimit; i++ {
		s.seeder.CreateTeamMember(repo.ID, fmt.Sprintf("member%d@example.com", i))
	}

	_, err := s.svc.AddMember(context.Background(), repo.ID, "overflow@example.com", false)
	assert.ErrorIs(s.T(), err, team_member_service.ErrFreeTierLimitHit)
}

// TestAddMember_FreeTierLimitBypassedWithSub allows >10 members when workspace has an active subscription.
func (s *TeamMemberTestSuite) TestAddMember_FreeTierLimitBypassedWithSub() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	for i := 0; i < team_member_service.FreeTierMemberLimit; i++ {
		s.seeder.CreateTeamMember(repo.ID, fmt.Sprintf("member%d@example.com", i))
	}

	// hasActiveSub = true → limit not checked
	member, err := s.svc.AddMember(context.Background(), repo.ID, "pro@example.com", true)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "pro@example.com", member.Email)
}

// TestRemoveMember_Success removes an existing member.
func (s *TeamMemberTestSuite) TestRemoveMember_Success() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)
	tm := s.seeder.CreateTeamMember(repo.ID, "carol@example.com")

	err := s.svc.RemoveMember(context.Background(), tm.ID)
	assert.NoError(s.T(), err)

	// Confirm deleted
	found, err := s.tmRepo.GetByID(context.Background(), tm.ID)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), found)
}

// TestRemoveMember_NotFound returns ErrMemberNotFound for an unknown ID.
func (s *TeamMemberTestSuite) TestRemoveMember_NotFound() {
	err := s.svc.RemoveMember(context.Background(), uuid.New())
	assert.ErrorIs(s.T(), err, team_member_service.ErrMemberNotFound)
}

// TestListMembers_ReturnsScopedToRepo ensures list is scoped to the given repo.
func (s *TeamMemberTestSuite) TestListMembers_ReturnsScopedToRepo() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repoA := s.seeder.CreateRepository(workspace.ID, user.ID)
	repoB := s.seeder.CreateRepository(workspace.ID, user.ID)

	s.seeder.CreateTeamMember(repoA.ID, "a1@example.com")
	s.seeder.CreateTeamMember(repoA.ID, "a2@example.com")
	s.seeder.CreateTeamMember(repoB.ID, "b1@example.com")

	members, err := s.svc.ListMembers(context.Background(), repoA.ID)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), members, 2)
	for _, m := range members {
		assert.Equal(s.T(), repoA.ID, m.RepoID)
	}
}
