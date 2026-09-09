package domain

import "time"

type Platform struct {
	ID           int64
	Title        string
	Abbreviation string
	Slug         string
	IconURL      string
	Checksum     string
	UpdatedAt    int64
	SyncedAt     time.Time
}

func NewPlatform(
	id int64,
	title string,
	abbreviation string,
	slug string,
	iconURL string,
	checksum string,
	updatedAt int64,
	syncedAt time.Time,
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