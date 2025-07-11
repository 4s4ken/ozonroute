package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/outbox"
)

type OutboxRepo struct {
	conn *pgxpool.Pool
}

func NewOutboxRepo(conn *pgxpool.Pool) *OutboxRepo {
	return &OutboxRepo{conn: conn}
}

func (r *OutboxRepo) MsgCreated(ctx context.Context, limit int) ([]outbox.ObEvent, error) {
	const query = `
		UPDATE outbox
		SET status = 'PROCESSING'
		WHERE id IN (
			SELECT id FROM outbox
			WHERE 
				status = 'CREATED'
				AND (last_attempt_at IS NULL OR last_attempt_at <= NOW() - INTERVAL '200 milliseconds')
				AND attempts < 3
			ORDER BY created_at
			LIMIT $1
			FOR UPDATE SKIP LOCKED
		)
		RETURNING id, payload, attempts, last_attempt_at, created_at, sent_at, status, error
	`

	rows, err := r.conn.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query for outbox: %v", err)
	}
	defer rows.Close()

	var events []outbox.ObEvent
	for rows.Next() {
		var event outbox.ObEvent
		err = rows.Scan(
			&event.ID,
			&event.Payload,
			&event.Attempts,
			&event.LastAttempted,
			&event.CreatedAt,
			&event.SentAt,
			&event.Status,
			&event.Error,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (r *OutboxRepo) SaveEvent(ctx context.Context, event outbox.Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event payload: %w", err)
	}

	const query = `
		INSERT INTO outbox (id, payload, status, created_at)
		VALUES ($1, $2, 'CREATED', NOW())
	`

	_, err = r.conn.Exec(ctx, query, event.EventID, payload)
	return err
}

func (r *OutboxRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string, reason error) error {
	var query string
	var args []interface{}

	switch status {
	case "PROCESSING":
		return nil

	case "COMPLETED":
		query = `
			UPDATE outbox
			SET status = 'COMPLETED', sent_at = NOW()
			WHERE id = $1
		`
		args = []interface{}{id}

	case "FAILED":
		query = `
			UPDATE outbox
			SET 
				attempts = attempts + 1,
				last_attempt_at = NOW(),
				error = $2,
				sent_at = NOW(),
				status = CASE 
					WHEN attempts > 3 THEN 'FAILED'
					ELSE 'CREATED'
				END
			WHERE id = $1
		`
		errMsg := ""
		if reason != nil {
			errMsg = reason.Error()
		}
		args = []interface{}{id, errMsg}

	default:
		return fmt.Errorf("unsupported status: %s", status)
	}

	_, err := r.conn.Exec(ctx, query, args...)
	return err
}

func (r *OutboxRepo) MsgProcessing(ctx context.Context, limit int) ([]outbox.ObEvent, error) {
	const query = `
		SELECT id, payload, attempts, last_attempt_at, created_at, sent_at, status, error
		FROM outbox
		WHERE status = 'PROCESSING'
		ORDER BY created_at
		LIMIT $1
		FOR UPDATE SKIP LOCKED
	`

	rows, err := r.conn.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query for outbox: %v", err)
	}
	defer rows.Close()

	var events []outbox.ObEvent
	for rows.Next() {
		var event outbox.ObEvent
		err = rows.Scan(
			&event.ID,
			&event.Payload,
			&event.Attempts,
			&event.LastAttempted,
			&event.CreatedAt,
			&event.SentAt,
			&event.Status,
			&event.Error,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}
