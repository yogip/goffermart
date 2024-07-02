-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
    id       BIGSERIAL PRIMARY KEY,
    email    VARCHAR(320) UNIQUE NOT NULL,
    password VARCHAR(80) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
