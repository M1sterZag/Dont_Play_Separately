package core_pgx_pool

import (
	"context"
	"fmt"
	"time"

	core_repository "github.com/M1sterZag/Dont_Play_Separately/internal/core/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	pool      *pgxpool.Pool
	opTimeout time.Duration
}

func NewPool(ctx context.Context, config Config) (*Pool, error) {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	pgxConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("parse pgxconfig: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("create pgxpool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping pgxpool: %w", err)
	}

	return &Pool{
		pool:      pool,
		opTimeout: config.Timeout,
	}, nil
}

func (p *Pool) Query(ctx context.Context, sql string, args ...any) (core_repository.Rows, error) {
	rows, err := p.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return pgxRows{rows}, nil
}

func (p *Pool) QueryRow(ctx context.Context, sql string, args ...any) core_repository.Row {
	return pgxRow{p.pool.QueryRow(ctx, sql, args...)}
}

func (p *Pool) Exec(ctx context.Context, sql string, args ...any) (core_repository.CommandTag, error) {
	tag, err := p.pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return pgxCommandTag{tag}, nil
}

func (p *Pool) Close() {
	p.pool.Close()
}

func (p *Pool) OpTimeout() time.Duration {
	return p.opTimeout
}
