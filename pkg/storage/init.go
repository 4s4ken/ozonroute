package storage

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

// NewPgxPool — создает подключение к БД.
func NewPgxPool(ctx context.Context, connString string) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf("failed to create pgxpool: %v", err)
	}

	if err = pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping pgxpool: %v", err)
	}

	return pool
}
