-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders
(
    order_id     BIGSERIAL PRIMARY KEY NOT NULL,                           -- Уникальный идентификатор заказа
    user_id      BIGINT                NOT NULL,                           -- ID пользователя, оформившего заказ
    status       INTEGER               NOT NULL DEFAULT 0,                 -- Текущий статус заказа
    weight       NUMERIC(10, 2)        NOT NULL,                           -- Вес заказа
    price        NUMERIC(10, 2)        NOT NULL,                           -- Цена заказа
    package_type varchar,                                                  -- Тип упаковки (nullable, enum)
    expires_at   TIMESTAMPTZ           NOT NULL,                           -- Дата истечения срока хранения заказа
    created_at   TIMESTAMPTZ           NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Время создания записи
    updated_at   TIMESTAMPTZ           NOT NULL DEFAULT CURRENT_TIMESTAMP  -- Время последнего обновления
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd