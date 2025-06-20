-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders
(
    order_id     UUID PRIMARY KEY        DEFAULT gen_random_uuid(), -- Уникальный идентификатор заказа
    user_id      UUID           NOT NULL default gen_random_uuid(), -- ID пользователя, оформившего заказ
    status       INTEGER        NOT NULL DEFAULT 0,                 -- Текущий статус заказа
    weight       NUMERIC(10, 2) NOT NULL,                           -- Вес заказа
    price        NUMERIC(10, 2) NOT NULL,                           -- Цена заказа
    package_type INTEGER,                                           -- Тип упаковки (nullable, enum)
    expires_at   TIMESTAMPTZ    NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Дата истечения срока хранения заказа
    created_at   TIMESTAMPTZ    NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Время создания записи
    updated_at   TIMESTAMPTZ    NOT NULL DEFAULT CURRENT_TIMESTAMP  -- Время последнего обновления
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd