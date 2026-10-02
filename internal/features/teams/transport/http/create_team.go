package teams_transport_http

import (
	"net/http"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	core_logger "github.com/M1sterZag/Dont_Play_Separately/internal/core/logger"
	core_http_middleware "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/middleware"
	core_http_request "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/request"
	core_http_response "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/response"
)

type CreateTeamResponse TeamDTOResponse

// CreateTeam creates a new team.
// @Summary Create team
// @Description Creates a team for the authenticated user as its owner. The owner automatically becomes a member.
// @Tags teams
// @Accept json
// @Produce json
// @Param request body CreateTeamRequest true "Create team payload"
// @Success 201 {object} CreateTeamResponse "Created"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /teams [post]
func (h *TeamsHTTPHandler) CreateTeam(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var createTeamRequest CreateTeamRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &createTeamRequest); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request body")
		return
	}

	userID, ok := core_http_middleware.UserIDFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(core_errors.ErrUnauthenticated, "failed to get user_id (context)")
		return
	}

	teamDomain := domain.Team{
		GameID:           createTeamRequest.GameID,
		PlatformID:       createTeamRequest.PlatformID,
		Title:            createTeamRequest.Title,
		Description:      createTeamRequest.Description,
		IsRatingRequired: createTeamRequest.IsRatingRequired,
		DesiredRating:    createTeamRequest.DesiredRating,
		ContactLink:      createTeamRequest.ContactLink,
		SlotsTotal:       createTeamRequest.SlotsTotal,
	}

	team, err := h.teamsService.CreateTeam(ctx, userID, teamDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create team")
		return
	}

	response := CreateTeamResponse(teamDTOFromDomain(team))

	responseHandler.JSONResponse(response, http.StatusCreated)
}
