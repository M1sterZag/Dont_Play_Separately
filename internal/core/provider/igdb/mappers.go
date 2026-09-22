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
		nonZeroInt64(g.UpdatedAt),
		nil,
	)
}

func nonEmptyStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func nonZeroInt64(n int64) *int64 {
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
