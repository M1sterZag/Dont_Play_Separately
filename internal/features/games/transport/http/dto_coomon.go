package games_transport_http

import "github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"

type GameDTOResponse struct {
	ID      int64   `json:"id"`
	Title   string  `json:"title"`
	Slug    *string `json:"slug"`
	IconURL *string `json:"icon_url"`
}

func gameDTOFromDomain(game domain.Game) GameDTOResponse {
	return GameDTOResponse{
		ID:      game.ID,
		Title:   game.Title,
		Slug:    game.Slug,
		IconURL: game.IconURL,
	}
}
