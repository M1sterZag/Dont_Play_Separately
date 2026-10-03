package teams_transport_http

import (
	"net/http"

	core_logger "github.com/M1sterZag/Dont_Play_Separately/internal/core/logger"
	core_http_request "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/request"
	core_http_response "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/response"
)

type GetTeamResponse TeamDTOResponse

// GetTeamByID returns a team with its members.
// @Summary Get team
// @Description Returns a team by ID together with its members.
// @Tags teams
// @Produce json
// @Param team_id path string true "Team ID" Format(uuid) example(330ca7c9-80bf-4808-a78b-d38cewer87b5)
// @Success 200 {object} GetTeamResponse "OK"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 404 {object} core_http_response.ErrorResponse "Not found"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /teams/{team_id} [get]
func (h *TeamsHTTPHandler) GetTeamByID(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	teamID, err := core_http_request.GetUUIDPathParam(r, "team_id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get team_id (path)")
		return
	}

	team, err := h.teamsService.GetTeamByID(ctx, teamID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get team")
		return
	}

	teamMembers, err := h.teamsService.ListMembers(ctx, teamID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get team members")
		return
	}

	response := GetTeamResponse(teamDTOFromDomain(team, teamMembers))

	responseHandler.JSONResponse(response, http.StatusOK)
}
