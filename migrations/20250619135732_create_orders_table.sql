-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders
(
    id           BIGINT PRIMARY KEY,
    order_id     uuid         NOT NULL,
    user_id      BIGINT         NOT NULL,
    status       INT            NOT NULL,
    weight       NUMERIC(10, 2) NOT NULL,
    price        NUMERIC(10, 2) NOT NULL,
    package_type INT,
    expires_at   TIMESTAMP      NOT NULL,
    created_at   TIMESTAMP      NOT NULL DEFAULT now(),
    updated_at   TIMESTAMP      NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd