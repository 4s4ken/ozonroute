package storage

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type DbManager interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Mode int

const (
	ModeWrite Mode = iota
	ModeRead
)

func (m Mode) String() string {
	switch m {
	case ModeWrite:
		return "write"
	case ModeRead:
		return "read"
	default:
		return "unknown"
	}
}

type Client struct {
	write *sql.DB
	read  *sql.DB
}

func New(writeDB, readDB *sql.DB) *Client {
	return &Client{write: writeDB, read: readDB}
}

func (c *Client) Get(mode Mode) DbManager {
	switch mode {
	case ModeRead:
		return &loggingExecutor{inner: c.read}
	case ModeWrite:
		fallthrough
	default:
		return &loggingExecutor{inner: c.write}
	}
}

func (c *Client) Tx(ctx context.Context, fn func(tx DbManager) error) error {
	tx, err := c.write.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	wrappedTx := &loggingExecutor{inner: tx}

	if err := fn(wrappedTx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

type loggingExecutor struct {
	inner DbManager
}

func (l *loggingExecutor) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	start := time.Now()
	res, err := l.inner.ExecContext(ctx, query, args...)
	log.Printf("Query %s duration %s args %v \n ERROR: %v", query, time.Since(start), args, err)
	return res, err
}

func (l *loggingExecutor) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	start := time.Now()
	rows, err := l.inner.QueryContext(ctx, query, args...)
	log.Printf("Query %s duration: %s args %v err: %v", query, time.Since(start), args, err)
	return rows, err
}

func (l *loggingExecutor) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	start := time.Now()
	row := l.inner.QueryRowContext(ctx, query, args...)
	log.Printf("Query %s duration: %s args %v", query, time.Since(start), args)
	return row
}
