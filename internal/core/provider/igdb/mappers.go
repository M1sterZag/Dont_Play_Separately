package core_igdb_provider

import (
	"strings"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
)

func (g Game) ToDomain() domain.Game {
	return domain.NewGame(
		g.ID,
		g.Name,
		nonEmptyStr(g.Slug),
		absoluteImageURL(g.Cover),
		nonEmptyStr(g.Checksum),
		nonZeroInt(g.UpdatedAt),
		nil,
	)
}

func (p Platform) ToDomain() domain.Platform {
	return domain.NewPlatform(
		p.ID,
		p.Name,
		p.Abbreviation,
		nonEmptyStr(p.Slug),
		absoluteImageURL(p.PlatformLogo),
		nonEmptyStr(p.Checksum),
		nonZeroInt(p.UpdatedAt),
		nil,
	)
}

func nonEmptyStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func nonZeroInt(n int) *int {
	if n == 0 {
		return nil
	}
	return &n
}

func absoluteImageURL(img *Image) *string {
	if img == nil || strings.TrimSpace(img.URL) == "" {
		return nil
	}
	url := "https:" + img.URL
	return &url
}
