package platforms_transport_http

import "github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"

type PlatformDTOResponse struct {
	ID           int64   `json:"id"`
	Title        string  `json:"title"`
	Abbreviation string  `json:"abbreviation"`
	Slug         *string `json:"slug"`
	IconURL      *string `json:"icon_url"`
}

func platformDTOFromDomain(platform domain.Platform) PlatformDTOResponse {
	return PlatformDTOResponse{
		ID:           platform.ID,
		Title:        platform.Title,
		Abbreviation: platform.Abbreviation,
		Slug:         platform.Slug,
		IconURL:      platform.IconURL,
	}
}