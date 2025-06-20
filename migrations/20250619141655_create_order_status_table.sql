-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_status
(
    id         BIGSERIAL PRIMARY KEY,
    order_id   BIGINT NOT NULL REFERENCES orders(order_id) ON DELETE CASCADE,
    status     TEXT NOT NULL,
    changed_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_order_status_order_time ON order_status(order_id, changed_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_status;
-- +goose StatementEnd