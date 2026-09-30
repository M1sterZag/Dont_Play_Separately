package teams_repository

import (
	"time"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	"github.com/google/uuid"
)

type TeamModel struct {
	ID               uuid.UUID
	Version          int
	OwnerID          uuid.UUID
	GameID           int
	PlatformID       int
	Title            string
	Description      *string
	IsRatingRequired bool
	DesiredRating    *string
	ContactLink      *string
	SlotsTotal       int
	SlotsTaken       int
	CreatedAt        time.Time
	IsActive         bool
}

func TeamDomainFromModel(teamModel TeamModel) domain.Team {
	return domain.NewTeam(
		teamModel.ID,
		teamModel.Version,
		teamModel.OwnerID,
		teamModel.GameID,
		teamModel.PlatformID,
		teamModel.Title,
		teamModel.Description,
		teamModel.IsRatingRequired,
		teamModel.DesiredRating,
		teamModel.ContactLink,
		teamModel.SlotsTotal,
		teamModel.SlotsTaken,
		teamModel.CreatedAt,
		teamModel.IsActive,
	)
}

func TeamDomainsFromModels(teamModels []TeamModel) []domain.Team {
	teamDomains := make([]domain.Team, len(teamModels))
	for i, model := range teamModels {
		teamDomains[i] = TeamDomainFromModel(model)
	}

	return teamDomains
}

type TeamMemberModel struct {
	TeamID   uuid.UUID
	UserID   uuid.UUID
	Nickname string
}

func TeamMemberDomainFromModel(teamMemberModel TeamMemberModel) domain.TeamMember {
	return domain.TeamMember{
		TeamID:   teamMemberModel.TeamID,
		UserID:   teamMemberModel.UserID,
		Nickname: teamMemberModel.Nickname,
	}
}

func TeamMemberDomainsFromModels(teamMemberModels []TeamMemberModel) []domain.TeamMember {
	teamMembersDomains := make([]domain.TeamMember, len(teamMemberModels))
	for i, model := range teamMemberModels {
		teamMembersDomains[i] = TeamMemberDomainFromModel(model)
	}

	return teamMembersDomains
}
