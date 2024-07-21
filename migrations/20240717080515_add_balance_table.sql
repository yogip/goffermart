-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS balance(
    id         BIGINT PRIMARY KEY,
    user_id    BIGINT REFERENCES users (id),
    current    BIGINT NOT NULL DEFAULT 0,
    withdrawn  BIGINT NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS balance;
-- +goose StatementEnd
