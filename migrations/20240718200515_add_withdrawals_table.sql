-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS withdrawals(
    id            BIGSERIAL PRIMARY KEY,
    user_id       BIGINT REFERENCES users (id),
    order_id      VARCHAR(255) REFERENCES orders(id) UNIQUE,
    sum           BIGINT NOT NULL DEFAULT 0,
    processed_at  TIMESTAMP NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS withdrawals;
-- +goose StatementEnd
