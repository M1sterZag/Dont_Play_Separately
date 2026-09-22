package games_transport_http

import (
	"fmt"
	"net/http"

	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	core_logger "github.com/M1sterZag/Dont_Play_Separately/internal/core/logger"
	core_http_request "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/request"
	core_http_response "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/response"
)

type SearchGamesResponse []GameDTOResponse

// SearchGames searches games by name.
// @Summary Search games
// @Description Searches the local catalog, falls back to IGDB if empty.
// @Tags games
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Page size"
// @Param offset query int false "Page offset"
// @Success 200 {object} SearchGamesResponse "OK"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /games/search [get]
func (h *GamesHTTPHandler) SearchGames(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	query := core_http_request.GetStringQueryParam(r, "q")
	if query == "" {
		responseHandler.ErrorResponse(fmt.Errorf("query param 'q' is required: %w", core_errors.ErrInvalidArgument), "failed to get search query")
		return
	}

	limit, err := core_http_request.GetIntQueryParam(r, "limit")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get limit query param")
		return
	}

	offset, err := core_http_request.GetIntQueryParam(r, "offset")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get offset query param")
		return
	}

	games, err := h.gamesService.SearchGames(ctx, query, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to search games")
		return
	}

	response := make(SearchGamesResponse, 0, len(games))
	for _, game := range games {
		response = append(response, gameDTOFromDomain(game))
	}

	responseHandler.JSONResponse(response, http.StatusOK)
}
