package platforms_repository

import (
	"time"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
)

type PlatformModel struct {
	ID           int64
	Title        string
	Abbreviation string
	Slug         *string
	IconURL      *string
	Checksum     *string
	UpdatedAt    *int64
	SyncedAt     *time.Time
}

func PlatformDomainFromModel(platformModel PlatformModel) domain.Platform {
	return domain.NewPlatform(
		platformModel.ID,
		platformModel.Title,
		platformModel.Abbreviation,
		platformModel.Slug,
		platformModel.IconURL,
		platformModel.Checksum,
		platformModel.UpdatedAt,
		platformModel.SyncedAt,
	)
}

func PlatformDomainsFromModels(platformModels []PlatformModel) []domain.Platform {
	platformDomains := make([]domain.Platform, len(platformModels))
	for i, platformModel := range platformModels {
		platformDomains[i] = PlatformDomainFromModel(platformModel)
	}

	return platformDomains
}