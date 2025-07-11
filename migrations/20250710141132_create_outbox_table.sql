-- +goose Up
-- +goose StatementBegin
CREATE TYPE event_status AS ENUM ('CREATED', 'PROCESSING', 'COMPLETED', 'FAILED');

CREATE TABLE IF NOT EXISTS outbox (
    id UUID PRIMARY KEY,
    payload JSONB NOT NULL,
    status event_status NOT NULL DEFAULT 'CREATED',
    error TEXT,
    attempts INT NOT NULL DEFAULT 0,
    last_attempt_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    sent_at TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS outbox;
-- +goose StatementEnd
