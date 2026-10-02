package teams_transport_http

import (
	"context"
	"net/http"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	core_http_server "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/server"
	"github.com/google/uuid"
)

type TeamsService interface {
	CreateTeam(ctx context.Context, ownerID uuid.UUID, team domain.Team) (domain.Team, error)
	GetTeamByID(ctx context.Context, teamID uuid.UUID) (domain.Team, error)
	ListTeams(ctx context.Context, filter domain.TeamFilter) ([]domain.Team, error)
	PatchTeam(ctx context.Context, userID, teamID uuid.UUID, patch domain.TeamPatch) (domain.Team, error)
	DeleteTeam(ctx context.Context, userID, teamID uuid.UUID) error
	JoinTeam(ctx context.Context, teamID, userID uuid.UUID) error
	LeaveTeam(ctx context.Context, teamID, userID uuid.UUID) error
	RemoveMember(ctx context.Context, teamID, userID, memberUserID uuid.UUID) error
	ListMembers(ctx context.Context, teamID uuid.UUID) ([]domain.TeamMember, error)
}

type TeamsHTTPHandler struct {
	teamsService TeamsService
}

func NewTeamsHTTPHandler(teamsService TeamsService) *TeamsHTTPHandler {
	return &TeamsHTTPHandler{
		teamsService: teamsService,
	}
}

func (h *TeamsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{Method: http.MethodPost, Path: "/teams", Handler: h.CreateTeam},
		{Method: http.MethodGet, Path: "/teams", Handler: h.ListTeams},
		{Method: http.MethodGet, Path: "/teams/{team_id}", Handler: h.GetTeamByID},
		{Method: http.MethodPut, Path: "/teams/{team_id}", Handler: h.PatchTeam},
		{Method: http.MethodDelete, Path: "/teams/{team_id}", Handler: h.DeleteTeam},
		{Method: http.MethodPost, Path: "/teams/{team_id}/join", Handler: h.JoinTeam},
		{Method: http.MethodPost, Path: "/teams/{team_id}/leave", Handler: h.LeaveTeam},
		{Method: http.MethodDelete, Path: "/teams/{team_id}/members/{user_id}", Handler: h.RemoveMember},
	}
}
