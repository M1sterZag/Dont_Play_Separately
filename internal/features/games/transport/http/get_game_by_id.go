package games_transport_http

import (
	"net/http"

	core_logger "github.com/M1sterZag/Dont_Play_Separately/internal/core/logger"
	core_http_request "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/request"
	core_http_response "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/response"
)

type GetGameResponse GameDTOResponse

// GetGameByID returns a game by its IGDB ID.
// @Summary Get game by ID
// @Description Returns a catalog game by its numeric ID.
// @Tags games
// @Produce json
// @Param id path int true "Game ID"
// @Success 200 {object} GetGameResponse "OK"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 404 {object} core_http_response.ErrorResponse "Not found"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /games/{id} [get]
func (h *GamesHTTPHandler) GetGameByID(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	gameID, err := core_http_request.GetInt64PathParam(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get game id path param")
		return
	}

	game, err := h.gamesService.GetGameByID(ctx, gameID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get game")
		return
	}

	response := GetGameResponse(gameDTOFromDomain(game))
	responseHandler.JSONResponse(response, http.StatusOK)
}
