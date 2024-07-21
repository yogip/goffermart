-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS withdrawals(
    id         BIGINT PRIMARY KEY,
    order       BIGINT REFERENCES users (id),
    sum         BIGINT NOT NULL DEFAULT 0,
    processed_at  BIGINT NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS withdrawals;
-- +goose StatementEnd
