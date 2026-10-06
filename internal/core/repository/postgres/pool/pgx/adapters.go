package core_pgx_pool

import (
	"context"
	"errors"
	"fmt"

	core_repository "github.com/M1sterZag/Dont_Play_Separately/internal/core/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	pgxViolatesForeignKeyErrorCode = "23503"
	pgxUniqueViolationErrorCode    = "23505"
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

type pgxTx struct {
	pgx.Tx
}

func (t pgxTx) Query(ctx context.Context, sql string, args ...any) (core_repository.Rows, error) {
	rows, err := t.Tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return pgxRows{rows}, nil
}

func (t pgxTx) QueryRow(ctx context.Context, sql string, args ...any) core_repository.Row {
	return pgxRow{t.Tx.QueryRow(ctx, sql, args...)}
}

func (t pgxTx) Exec(ctx context.Context, sql string, args ...any) (core_repository.CommandTag, error) {
	tag, err := t.Tx.Exec(ctx, sql, args...)
	if err != nil {
		return nil, mapErrors(err)
	}

	return pgxCommandTag{tag}, nil
}

func (t pgxTx) Commit(ctx context.Context) error {
	if err := t.Tx.Commit(ctx); err != nil {
		return mapErrors(err)
	}

	return nil
}

func (t pgxTx) Rollback(ctx context.Context) error {
	if err := t.Tx.Rollback(ctx); err != nil {
		return mapErrors(err)
	}

	return nil
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
		return fmt.Errorf("%v: %w", err, core_repository.ErrUniqueViolation)
	}

	return fmt.Errorf("%v: %w", err, core_repository.ErrUnknown)
}
