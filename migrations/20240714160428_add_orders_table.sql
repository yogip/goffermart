-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders(
    id       VARCHAR(255) PRIMARY KEY,
    user_id  BIGINT REFERENCES users (id),
    status   VARCHAR(20) NOT NULL DEFAULT 'NEW',
    created_at  TIMESTAMP DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
