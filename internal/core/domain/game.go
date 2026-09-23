package domain

import "time"

type Game struct {
	ID        int64
	Title     string
	Slug      *string
	IconURL   *string
	Checksum  *string
	UpdatedAt *int64
	SyncedAt  *time.Time
}

func NewGame(
	id int64,
	title string,
	slug *string,
	iconURL *string,
	checksum *string,
	updatedAt *int64,
	syncedAt *time.Time,
) Game {
	return Game{
		ID:        id,
		Title:     title,
		Slug:      slug,
		IconURL:   iconURL,
		Checksum:  checksum,
		UpdatedAt: updatedAt,
		SyncedAt:  syncedAt,
	}
}
