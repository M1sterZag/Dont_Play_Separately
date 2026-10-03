package teams_transport_http

import (
	"time"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	core_http_types "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/types"
	"github.com/google/uuid"
)

type TeamDTOResponse struct {
	ID               uuid.UUID               `json:"id"`
	Version          int                     `json:"version"`
	OwnerID          uuid.UUID               `json:"owner_id"`
	GameID           int                     `json:"game_id"`
	PlatformID       int                     `json:"platform_id"`
	Title            string                  `json:"title"`
	Description      *string                 `json:"description,omitempty"`
	IsRatingRequired bool                    `json:"is_rating_required"`
	DesiredRating    *string                 `json:"desired_rating,omitempty"`
	ContactLink      *string                 `json:"contact_link,omitempty"`
	SlotsTotal       int                     `json:"slots_total"`
	SlotsTaken       int                     `json:"slots_taken"`
	CreatedAt        time.Time               `json:"created_at"`
	IsActive         bool                    `json:"is_active"`
	Members          []TeamMemberDTOResponse `json:"members,omitempty"`
}

type TeamMemberDTOResponse struct {
	UserID   uuid.UUID `json:"user_id"`
	Nickname string    `json:"nickname"`
}

func teamDTOFromDomain(team domain.Team, members ...[]domain.TeamMember) TeamDTOResponse {
	var teamMembers []TeamMemberDTOResponse
	if len(members) > 0 {
		teamMembers = teamMembersDTOFromDomain(members[0])
	}
	return TeamDTOResponse{
		ID:               team.ID,
		Version:          team.Version,
		OwnerID:          team.OwnerID,
		GameID:           team.GameID,
		PlatformID:       team.PlatformID,
		Title:            team.Title,
		Description:      team.Description,
		IsRatingRequired: team.IsRatingRequired,
		DesiredRating:    team.DesiredRating,
		ContactLink:      team.ContactLink,
		SlotsTotal:       team.SlotsTotal,
		SlotsTaken:       team.SlotsTaken,
		CreatedAt:        team.CreatedAt,
		IsActive:         team.IsActive,
		Members:          teamMembers,
	}
}

func teamMembersDTOFromDomain(members []domain.TeamMember) []TeamMemberDTOResponse {
	var teamMembers []TeamMemberDTOResponse
	for _, member := range members {
		teamMembers = append(teamMembers, TeamMemberDTOResponse{
			UserID:   member.UserID,
			Nickname: member.Nickname,
		})
	}
	return teamMembers
}

type CreateTeamRequest struct {
	GameID           int     `json:"game_id" validate:"gt=0"`
	PlatformID       int     `json:"platform_id" validate:"gt=0"`
	Title            string  `json:"title" validate:"required,min=1,max=100"`
	Description      *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	IsRatingRequired bool    `json:"is_rating_required"`
	DesiredRating    *string `json:"desired_rating,omitempty" validate:"omitempty,max=200"`
	ContactLink      *string `json:"contact_link,omitempty" validate:"omitempty,max=2000"`
	SlotsTotal       int     `json:"slots_total" validate:"gt=0"`
}

type PatchTeamRequest struct {
	Title            core_http_types.Nullable[string] `json:"title,omitempty" validate:"omitempty,min=1,max=100"`
	Description      core_http_types.Nullable[string] `json:"description,omitempty" validate:"omitempty,max=1000"`
	IsRatingRequired *bool                            `json:"is_rating_required,omitempty"`
	DesiredRating    core_http_types.Nullable[string] `json:"desired_rating,omitempty" validate:"omitempty,max=200"`
	ContactLink      core_http_types.Nullable[string] `json:"contact_link,omitempty" validate:"omitempty,max=2000"`
	SlotsTotal       *int                             `json:"slots_total,omitempty" validate:"omitempty,gt=0"`
}

func teamPatchFromRequest(request PatchTeamRequest) domain.TeamPatch {
	return domain.NewTeamPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.IsRatingRequired,
		request.DesiredRating.ToDomain(),
		request.ContactLink.ToDomain(),
		request.SlotsTotal,
	)
}
