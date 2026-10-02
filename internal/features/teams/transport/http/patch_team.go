package teams_transport_http

import (
	"net/http"

	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	core_logger "github.com/M1sterZag/Dont_Play_Separately/internal/core/logger"
	core_http_middleware "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/middleware"
	core_http_request "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/request"
	core_http_response "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/response"
)

type PatchTeamResponse TeamDTOResponse

// PatchTeam partially updates a team.
// @Summary Update team
// @Description Partially updates a team. Only the owner can update; pass "description": null to clear it.
// @Tags teams
// @Accept json
// @Produce json
// @Param team_id path string true "Team ID" Format(uuid) example(330ca7c9-80bf-4808-a78b-d38cewer87b5)
// @Param request body PatchTeamRequest true "Team patch payload"
// @Success 200 {object} PatchTeamResponse "OK"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure 403 {object} core_http_response.ErrorResponse "Forbidden"
// @Failure 404 {object} core_http_response.ErrorResponse "Not found"
// @Failure 409 {object} core_http_response.ErrorResponse "Conflict (optimistic lock)"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /teams/{team_id} [put]
func (h *TeamsHTTPHandler) PatchTeam(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var patchTeamRequest PatchTeamRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &patchTeamRequest); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request body")
		return
	}

	teamPatch := teamPatchFromRequest(patchTeamRequest)

	userID, ok := core_http_middleware.UserIDFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(core_errors.ErrUnauthenticated, "failed to get user_id (context)")
		return
	}

	teamID, err := core_http_request.GetUUIDPathParam(r, "team_id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get team_id (path)")
		return
	}

	teamDomain, err := h.teamsService.PatchTeam(ctx, userID, teamID, teamPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch team")
		return
	}

	response := PatchTeamResponse(teamDTOFromDomain(teamDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}
