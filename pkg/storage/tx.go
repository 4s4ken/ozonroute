package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Transaction struct {
	db *pgxpool.Pool
	tx pgx.Tx
}

type txKey struct{}

// injectTx injects transaction to context
func injectTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// extractTx extracts transaction from context
func extractTx(ctx context.Context) pgx.Tx {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}
	return nil
}

// Commit метод фиксации блока транзакции
func (t *Transaction) Commit(ctx context.Context) error {
	err := t.tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("transaction.Commit: %w", err)
	}
	return nil
}

// Rollback метод отката блока транзакции
func (t *Transaction) Rollback(ctx context.Context) error {
	err := t.tx.Rollback(ctx)
	if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
		return fmt.Errorf("transaction.Rollback: %w", err)
	}
	return nil
}

// WithinTransaction runs function within transaction
//
// The transaction commits when function were finished without error
func (t *Transaction) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	// begin transaction
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	// run callback
	err = tFunc(injectTx(ctx, tx))
	if err != nil {
		// if error, rollback
		err := tx.Rollback(ctx)
		if err != nil {
			return err
		}
	}

	// if no error, commit
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

// Exec выполняет заданный SQL запрос с параметрами, возвращает ошибку или тэг запроса
func (t *Transaction) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	tx := extractTx(ctx)

	if tx != nil {
		return tx.Exec(ctx, query, args...)
	}

	return t.db.Exec(ctx, query, args...)
}

// Query выполняет запрос и возвращает объект для получения многострочного результата
func (t *Transaction) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	tx := extractTx(ctx)

	if tx != nil {
		return tx.Query(ctx, query, args...)
	}

	return t.db.Query(ctx, query, args...)
}

// QueryRow выполняет запрос и возвращает объект для получения одиночной строки
func (t *Transaction) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	tx := extractTx(ctx)

	if tx != nil {
		return tx.QueryRow(ctx, query, args...)
	}

	return t.db.QueryRow(ctx, query, args...)
}
