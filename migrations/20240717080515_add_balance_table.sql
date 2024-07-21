-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS balance(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT REFERENCES users (id) UNIQUE,
    current    NUMERIC(8, 2) NOT NULL DEFAULT 0,
    withdrawn  NUMERIC(8, 2) NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS balance;
-- +goose StatementEnd
