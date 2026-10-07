package teams_transport_http

import (
	"net/http"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	core_logger "github.com/M1sterZag/Dont_Play_Separately/internal/core/logger"
	core_http_request "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/request"
	core_http_response "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/response"
)

type ListTeamsResponse []TeamDTOResponse

// ListTeams returns active teams with optional filtering.
// @Summary List teams
// @Description Returns active teams with optional filtering by search query, game and platform.
// @Tags teams
// @Produce json
// @Param q query string false "Search by team title, game title or platform title/abbreviation"
// @Param game_id query int false "Filter by game ID"
// @Param platform_id query int false "Filter by platform ID"
// @Param limit query int false "Page size (default 20)"
// @Param offset query int false "Page offset (default 0)"
// @Success 200 {object} ListTeamsResponse "OK"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /teams [get]
func (h *TeamsHTTPHandler) ListTeams(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	search := core_http_request.GetStringQueryParam(r, "q")

	gameID, err := core_http_request.GetIntQueryParam(r, "game_id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get game_id (query)")
		return
	}

	platformID, err := core_http_request.GetIntQueryParam(r, "platform_id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get platform_id (query)")
		return
	}

	limit, err := core_http_request.GetIntQueryParam(r, "limit")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get limit (query)")
		return
	}

	offset, err := core_http_request.GetIntQueryParam(r, "offset")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get offset (query)")
		return
	}

	filter := domain.TeamFilter{
		Search:     search,
		GameID:     gameID,
		PlatformID: platformID,
		Limit:      limit,
		Offset:     offset,
	}

	teamDomains, err := h.teamsService.ListTeams(ctx, filter)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get teams")
		return
	}

	response := make(ListTeamsResponse, 0, len(teamDomains))
	for _, team := range teamDomains {
		response = append(response, teamDTOFromDomain(team))
	}

	responseHandler.JSONResponse(response, http.StatusOK)
}
