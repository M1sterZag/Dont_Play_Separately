package games_transport_http

import (
	"context"
	"net/http"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	core_http_server "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/server"
)

type GamesService interface {
	GetGameByID(ctx context.Context, gameID int64) (domain.Game, error)
	SearchGames(ctx context.Context, searchQuery string, limit, offset *int) ([]domain.Game, error)
}

type GamesHTTPHandler struct {
	gamesService GamesService
}

func NewGamesHTTPHandler(gamesService GamesService) *GamesHTTPHandler {
	return &GamesHTTPHandler{
		gamesService: gamesService,
	}
}

func (h *GamesHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method: http.MethodGet,
			Path: "/games/{id}",
			Handler: h.GetGameByID,
		},
		{
			Method: http.MethodGet,
			Path: "/games/search",
			Handler: h.SearchGames,
		},
	}
}
