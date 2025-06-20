-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_status_history
(
    id         BIGSERIAL PRIMARY KEY,                      -- Уникальный ID изменения статуса
    order_id   UUID      NOT NULL,                         -- Ссылка на заказ, которому принадлежит история
    status     INT         NOT NULL,                         -- Новый статус заказа после изменения
    changed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP-- Дата/время изменения статуса
);

-- +goose StatementEnd
CREATE INDEX CONCURRENTLY idx_status_history_order_changed ON order_status_history (order_id, changed_at DESC);

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_status_history;
-- +goose StatementEnd
