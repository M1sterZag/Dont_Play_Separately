package teams_service

import (
	"context"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	"github.com/google/uuid"
)

type TeamsRepository interface {
	CreateTeam(ctx context.Context, team domain.Team) (domain.Team, error)
	GetTeamByID(ctx context.Context, teamID uuid.UUID) (domain.Team, error)
	ListTeams(ctx context.Context, filter domain.TeamFilter) ([]domain.Team, error)
	PatchTeam(ctx context.Context, teamID uuid.UUID, patch domain.Team) (domain.Team, error)
	DeleteTeam(ctx context.Context, teamID uuid.UUID) error
	JoinTeam(ctx context.Context, teamID, userID uuid.UUID) error
	LeaveTeam(ctx context.Context, teamID, userID uuid.UUID) error
	RemoveMember(ctx context.Context, teamID, userID uuid.UUID) error
	ListMembers(ctx context.Context, teamID uuid.UUID) ([]domain.TeamMember, error)
}

type TeamsService struct {
	teamsRepository TeamsRepository
}

func NewTeamsService(teamsRepository TeamsRepository) *TeamsService {
	return &TeamsService{
		teamsRepository: teamsRepository,
	}
}
