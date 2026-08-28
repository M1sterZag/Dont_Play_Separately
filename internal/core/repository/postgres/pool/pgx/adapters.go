package core_pgx_pool

import (
	"errors"
	"fmt"

	core_repository "github.com/M1sterZag/Dont_Play_Separately/internal/core/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	pgxViolatesForeignKeyErrorCode = "23503"
	pgxUniqueViolationErrorCode = "23505"
)

type pgxRows struct {
	pgx.Rows
}

func (r pgxRows) Scan(dest ...any) error {
	if err := r.Rows.Scan(dest...); err != nil {
		return mapErrors(err)
	}

	return nil
}

func (r pgxRows) Err() error {
	if err := r.Rows.Err(); err != nil {
		return mapErrors(err)
	}

	return nil
}

type pgxRow struct {
	pgx.Row
}

func (r pgxRow) Scan(dest ...any) error {
	if err := r.Row.Scan(dest...); err != nil {
		return mapErrors(err)
	}

	return nil
}

type pgxCommandTag struct {
	pgconn.CommandTag
}

func mapErrors(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return core_repository.ErrNoRows
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgxViolatesForeignKeyErrorCode {
		return fmt.Errorf("%v: %w", err, core_repository.ErrViolatesForeignKey)
	}

	if errors.As(err, &pgErr) && pgErr.Code == pgxUniqueViolationErrorCode {
		return  fmt.Errorf("%v: %w", err, core_repository.ErrUniqueViolation)
	}

	return fmt.Errorf("%v: %w", err, core_repository.ErrUnknown)
}
