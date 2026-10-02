package teams_transport_http

import (
	"net/http"

	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	core_logger "github.com/M1sterZag/Dont_Play_Separately/internal/core/logger"
	core_http_middleware "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/middleware"
	core_http_request "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/request"
	core_http_response "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/response"
)

// DeleteTeam deletes a team.
// @Summary Delete team
// @Description Deletes a team. Only the owner can delete.
// @Tags teams
// @Param team_id path string true "Team ID" Format(uuid) example(330ca7c9-80bf-4808-a78b-d38cewer87b5)
// @Success 204 "No Content"
// @Failure 401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure 403 {object} core_http_response.ErrorResponse "Forbidden"
// @Failure 404 {object} core_http_response.ErrorResponse "Not found"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /teams/{team_id} [delete]
func (h *TeamsHTTPHandler) DeleteTeam(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

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

	if err := h.teamsService.DeleteTeam(ctx, userID, teamID); err != nil {
		responseHandler.ErrorResponse(err, "failed to delete team")
		return
	}

	responseHandler.NoContentResponse()
}
