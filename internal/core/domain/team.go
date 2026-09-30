package domain

import (
	"fmt"
	"time"

	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	"github.com/google/uuid"
)

type Team struct {
	ID               uuid.UUID
	Version          int
	OwnerID          uuid.UUID
	GameID           int
	PlatformID       int
	Title            string
	Description      *string
	IsRatingRequired bool
	DesiredRating    *string
	ContactLink      *string
	SlotsTotal       int
	SlotsTaken       int
	CreatedAt        time.Time
	IsActive         bool
}

func NewTeam(
	id uuid.UUID,
	version int,
	ownerID uuid.UUID,
	gameID int,
	platformID int,
	title string,
	description *string,
	isRatingRequired bool,
	desiredRating *string,
	contactLink *string,
	slotsTotal int,
	slotsTaken int,
	createdAt time.Time,
	isActive bool,
) Team {
	return Team{
		ID:               id,
		Version:          version,
		OwnerID:          ownerID,
		GameID:           gameID,
		PlatformID:       platformID,
		Title:            title,
		Description:      description,
		IsRatingRequired: isRatingRequired,
		DesiredRating:    desiredRating,
		ContactLink:      contactLink,
		SlotsTotal:       slotsTotal,
		SlotsTaken:       slotsTaken,
		CreatedAt:        createdAt,
		IsActive:         isActive,
	}
}

func (t *Team) Validate() error {
	titleLen := len([]rune(t.Title))
	if titleLen < 1 || titleLen > 100 {
		return fmt.Errorf("invalid `Title` len: %d: %w", titleLen, core_errors.ErrInvalidArgument)
	}

	if t.Description != nil {
		descriptionLen := len([]rune(*t.Description))
		if descriptionLen < 1 || descriptionLen > 1000 {
			return fmt.Errorf("invalid `Description` len: %d: %w", descriptionLen, core_errors.ErrInvalidArgument)
		}
	}

	if t.IsRatingRequired && t.DesiredRating == nil {
		return fmt.Errorf("`IsRatingRequired` is %t but `DesiredRating` is nil: %w", t.IsRatingRequired, core_errors.ErrInvalidArgument)
	}
	if !t.IsRatingRequired && t.DesiredRating != nil {
		return fmt.Errorf("`IsRatingRequired` is %t but `DesiredRating` is not nil: %w", t.IsRatingRequired, core_errors.ErrInvalidArgument)
	}

	if t.SlotsTotal <= 0 {
		return fmt.Errorf("invalid `SlotsTotal` value: %d: %w", t.SlotsTotal, core_errors.ErrInvalidArgument)
	}

	if t.SlotsTotal < t.SlotsTaken {
		return fmt.Errorf("`SlotsTaken` (%d) exceeds `SlotsTotal` (%d): %w", t.SlotsTaken, t.SlotsTotal, core_errors.ErrInvalidArgument)
	}

	if t.GameID <= 0 {
		return fmt.Errorf("invalid `GameID` value: %d: %w", t.GameID, core_errors.ErrInvalidArgument)
	}

	if t.PlatformID <= 0 {
		return fmt.Errorf("invalid `PlatformID` value: %d: %w", t.PlatformID, core_errors.ErrInvalidArgument)
	}

	return nil
}

func (t *Team) ApplyPatch(patch TeamPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate team patch: %w", err)
	}

	tmp := *t
	if patch.Title.Set {
		tmp.Title = *patch.Title.Value
	}

	if patch.Description.Set {
		tmp.Description = patch.Description.Value
	}

	if patch.IsRatingRequired != nil {
		tmp.IsRatingRequired = *patch.IsRatingRequired
	}

	if patch.DesiredRating.Set {
		tmp.DesiredRating = patch.DesiredRating.Value
	}

	if patch.ContactLink.Set {
		tmp.ContactLink = patch.ContactLink.Value
	}

	if patch.SlotsTotal != nil {
		tmp.SlotsTotal = *patch.SlotsTotal
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched team: %w", err)
	}

	*t = tmp

	return nil
}

type TeamPatch struct {
	Title            Nullable[string]
	Description      Nullable[string]
	IsRatingRequired *bool
	DesiredRating    Nullable[string]
	ContactLink      Nullable[string]
	SlotsTotal       *int
}

func NewTeamPatch(
	title Nullable[string],
	description Nullable[string],
	isRatingRequired *bool,
	desiredRating Nullable[string],
	contactLink Nullable[string],
	slotsTotal *int,
) TeamPatch {
	return TeamPatch{
		Title:            title,
		Description:      description,
		IsRatingRequired: isRatingRequired,
		DesiredRating:    desiredRating,
		ContactLink:      contactLink,
		SlotsTotal:       slotsTotal,
	}
}

func (p *TeamPatch) Validate() error {
	if p.Title.Set && p.Title.Value == nil {
		return fmt.Errorf("`Title` can`t be patched to `NULL`: %w", core_errors.ErrInvalidArgument)
	}

	if p.SlotsTotal != nil && *p.SlotsTotal <= 0 {
		return fmt.Errorf("`SlotsTotal` can`t be patched to '0': %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

type TeamMember struct {
	TeamID   uuid.UUID
	UserID   uuid.UUID
	Nickname string
}

type TeamFilter struct {
	Search     string
	GameID     *int
	PlatformID *int
	Limit      *int
	Offset     *int
}

func (f *TeamFilter) Validate() error {
	if f.Limit != nil && *f.Limit < 0 {
		return fmt.Errorf("`Limit` can`t be bellow zero: %d: %w", *f.Limit, core_errors.ErrInvalidArgument)
	}

	if f.Offset != nil && *f.Offset < 0 {
		return fmt.Errorf("`Offset` can`t be bellow zero: %d: %w", *f.Offset, core_errors.ErrInvalidArgument)
	}

	return nil
}
