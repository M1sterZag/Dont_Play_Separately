package domain

import "time"

type Game struct {
	ID        int
	Title     string
	Slug      *string
	IconURL   *string
	Checksum  *string
	UpdatedAt *int
	SyncedAt  *time.Time
}

func NewGame(
	id int,
	title string,
	slug *string,
	iconURL *string,
	checksum *string,
	updatedAt *int,
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
