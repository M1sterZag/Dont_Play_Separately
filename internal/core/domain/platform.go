package domain

import "time"

type Platform struct {
	ID           int
	Title        string
	Abbreviation string
	Slug         *string
	IconURL      *string
	Checksum     *string
	UpdatedAt    *int
	SyncedAt     *time.Time
}

func NewPlatform(
	id int,
	title string,
	abbreviation string,
	slug *string,
	iconURL *string,
	checksum *string,
	updatedAt *int,
	syncedAt *time.Time,
) Platform {
	return Platform{
		ID:           id,
		Title:        title,
		Abbreviation: abbreviation,
		Slug:         slug,
		IconURL:      iconURL,
		Checksum:     checksum,
		UpdatedAt:    updatedAt,
		SyncedAt:     syncedAt,
	}
}