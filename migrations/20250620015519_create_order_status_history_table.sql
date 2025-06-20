-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_status_history
(
    id         BIGSERIAL PRIMARY KEY,
    order_id   BIGINT    NOT NULL,
    status     INT       NOT NULL,
    changed_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_status_history;
-- +goose StatementEnd
