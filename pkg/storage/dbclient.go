package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Mode int

const (
	ModeRead Mode = iota
	ModeWrite
)

type DbManager interface {
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
}

type Client struct {
	write *pgxpool.Pool
	read  *pgxpool.Pool
}

func New(write, read *pgxpool.Pool) *Client {
	return &Client{write: write, read: read}
}

func (c *Client) Get(mode Mode) DbManager {
	if mode == ModeRead {
		return c.read
	}
	return c.write
}

func (c *Client) Tx(ctx context.Context, f func(DbManager) error) error {
	tx, err := c.write.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	if err = f(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
