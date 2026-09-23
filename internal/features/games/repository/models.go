package games_repository

import (
	"time"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
)

type GameModel struct {
	ID        int64
	Title     string
	Slug      *string
	IconURL   *string
	Checksum  *string
	UpdatedAt *int64
	SyncedAt  *time.Time
}

func GameDomainFromModel(gameModel GameModel) domain.Game {
	return domain.NewGame(
		gameModel.ID,
		gameModel.Title,
		gameModel.Slug,
		gameModel.IconURL,
		gameModel.Checksum,
		gameModel.UpdatedAt,
		gameModel.SyncedAt,
	)
}

func GameDomainsFromModels(gameModels []GameModel) []domain.Game {
	gameDomains := make([]domain.Game, len(gameModels))
	for i, gameModel := range gameModels {
		gameDomains[i] = GameDomainFromModel(gameModel)
	}

	return gameDomains
}
