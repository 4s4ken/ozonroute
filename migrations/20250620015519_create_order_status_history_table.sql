-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_status_history
(
    id         BIGSERIAL PRIMARY KEY,                        -- Уникальный ID изменения статуса
    order_id   BIGINT      NOT NULL,                         -- Ссылка на заказ, которому принадлежит история
    status     INT         NOT NULL,                         -- Новый статус заказа после изменения
    changed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP-- Дата/время изменения статуса
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_status_history;
-- +goose StatementEnd
